package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"reflect"
	"sync"
	"time"
)

type LoadBalancer struct {
	Mutex			             sync.Mutex             // Lock used by load balancer while updating one or more of its attributes
	servers			             []*Server              // list of servers in the server pool
	failureCount 	             map[string]int         // tracks # of times a server has been found unhealthy continously.
	removeAfter		             int                    // # of attempts after which an unhealthy server will be removed from the server pool
	healthCheckInterval          time.Duration          // time interval after which health check needs to be performed by load balancer for each server
	strategy                     LBStrategy       // Can be one of the load balancing strategies 
	connections                  map[string]int         // stores the current number of active connections on each server

}

// Load balancer methods
// This one is used to update the strategy when the user requests it
func (lb *LoadBalancer) setStrategy(newStrategy LBStrategy) {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()

	lb.strategy = newStrategy
	fmt.Printf("Load balancing strategy change to: %T\n", newStrategy)
}


// This one is used for finding the next server using the currently selected strategy
func (lb *LoadBalancer) SelectServer() *Server {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()

	server, err := lb.strategy.SelectServer(lb)

	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	return server
}

// This one is called by user manually when new server needs to be added in the server pool
func (lb *LoadBalancer) addServer(server_addr string) {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()

	u, _ := url.Parse(server_addr)
	server := &Server{URL: u, IsHealthy: true}
	lb.servers = append(lb.servers, server)
	log.Printf("Added new server @ %s", server_addr)
}

// this function is called by load balancer when we need to remove an unhealthy
// server
func (lb *LoadBalancer) removeServer(server_addr string) {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()

	// TODO: Once we remove a server, we also need to ensure that the health check
	// should not happen for that server
	// I think currently it's happening
	for i, server := range lb.servers {
		if server.URL.String() == server_addr {
			lb.servers = append(lb.servers[:i], lb.servers[i+1:]...)         // remove the server from the server pool
			delete(lb.failureCount, server_addr)		                     // delete server entry from the `server-># of failures` mapping
			delete(lb.connections, server_addr)                              // delete server entry from the `server-># of connections` mapping
			log.Printf("Removed server @ %s", server_addr)
			return
		}
	}
}


func (lb *LoadBalancer) healthCheck() {
	for range time.Tick(lb.healthCheckInterval) {
		for _, server := range lb.servers {

			// send each server a Head request
			// this helps us get the health status of the server
			// without downloading the response of the request
			res, err := http.Head(server.URL.String())
			server.Mutex.Lock()
			if err != nil || res.StatusCode != http.StatusOK {
				fmt.Printf("%s is down\n", server.URL)
				server.IsHealthy = false
				lb.failureCount[server.URL.String()]++
				fmt.Printf("Server %s down (%d/%d failures)\n", server.URL.String(), lb.failureCount[server.URL.String()], lb.removeAfter)
			
				if lb.failureCount[server.URL.String()] >= lb.removeAfter {
					lb.removeServer(server.URL.String())
				}
			} else {
				server.IsHealthy = true
				lb.failureCount[server.URL.String()] = 0
			}
			server.Mutex.Unlock()
		}
		
	}
}

type Server struct {
	URL			*url.URL
	IsHealthy	bool
	Mutex		sync.Mutex
}

func (server *Server) ReverseProxy() *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(server.URL)
}


type Config struct {
	Port 					string		`json:"port"`
	HealthCheckInterval		string		`json:"healthCheckInterval"`
	Servers					[]string 	`json:"servers"`
	LbAlgo					string		`json:"lbAlgo"`
	MaxRetries				int			`json:"maxRetries"`
}

func loadConfig(file string) (Config, error) {
	var config Config
	data, err := os.ReadFile(file)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}


func printStructFields(s interface{}) {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	for i := 0; i < val.NumField(); i++ {
		fmt.Printf("%s: %v\n", typ.Field(i).Name, val.Field(i).Interface())
	}
}


func forwardRequest(server *Server, w http.ResponseWriter, r *http.Request) error {
	recorder := &ResponseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
	server.ReverseProxy().ServeHTTP(recorder, r)

	// If response is a server error (5xx), return error
	if recorder.statusCode >= 500 {
		return fmt.Errorf("server returned %d", recorder.statusCode)
	}
	return nil
}

// Custom Response Recorder to capture status code
type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
}

// func (r *ResponseRecorder) WriteHeader(code int) {
// 	r.statusCode = code
// 	// r.ResponseWriter.WriteHeader(code)
// }

func main() {

	// load server config
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err.Error())
	}

	// pretty print the config
	printStructFields(config)

	// get the health check duration from the config
	healthCheckInterval, err := time.ParseDuration(config.HealthCheckInterval)
	if err != nil {
		log.Fatalf("Invalid health check interval: %s", err.Error())
	}

	// initialise backend servers specified as per the config
	var servers []*Server
	for _, serverUrl := range config.Servers {
		u, _ := url.Parse(serverUrl)
		server := &Server{URL: u, IsHealthy: true}
		servers = append(servers, server)
	}

	// Load the load balancer strategy
	strategy, err := NewStrategy(config.LbAlgo)
	if err != nil {
		log.Fatalf("Error: %s", err.Error());
	}

	// Instantiate the load balancer
	lb := LoadBalancer{
						servers: servers,
						failureCount: make(map[string]int),
						removeAfter: 30,
						healthCheckInterval: healthCheckInterval,
						strategy: strategy,
						connections: make(map[string]int),
					}

					
	// Initiate health checks by load balancer in a separate go routine
	go lb.healthCheck()

	maxRetries := config.MaxRetries

	// handler function to handle request to the load balancer
	// insides this handler function, we will call the allocateServer method that will return the server that needs to serve this request
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		
		var lastErr error
		for attempt := 1; attempt <= maxRetries; attempt++ {

			// find the next server using the specified load balancing strategy
			// for now we are implementing round-robin strategy
			server := lb.SelectServer()
		
			if server == nil {
				http.Error(w, "No healthy server available", http.StatusServiceUnavailable)
				return
			}

			// if we found a valid server
			w.Header().Add("X-Forwarded-Server", server.URL.String())

			// try processing the request
			err := forwardRequest(server, w, r)

			if err == nil {
				return
			}
			
			// Log failure & retry
			lastErr = err
			fmt.Printf("Attempt %d failed for server %s: %v. Retrying...\n", attempt, server.URL.String(), err)

			// Optional: Add small delay before retrying (exponential backoff possible)
			time.Sleep(100 * time.Millisecond)

		}

		// If all retries fail, return error response
		fmt.Printf("Request failed after %d retries: %v", maxRetries, lastErr)
		
	})

	// DEMO: add a new server after 2 minutes delay
	// time.AfterFunc(2 * time.Minute, func() {
	// 	log.Printf("Current time: %s", time.Now().Format("2006-01-02 15:04:05"))
	// 	lb.addServer("http://localhost:5900")
	// })


	// // DEMO: remove a server after 3 mins delay
	// time.AfterFunc(3 * time.Minute, func() {
	// 	log.Printf("Current time: %s", time.Now().Format("2006-01-02 15:04:05"))
	// 	lb.removeServer("http://localhost:5001")
	// })

	log.Println("Starting load balancer on port", config.Port)
	err = http.ListenAndServe(config.Port, nil)
	if err != nil {
		log.Fatalf("Error starting load balancer: %s\n", err.Error())
	}

}
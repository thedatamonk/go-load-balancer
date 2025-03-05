package algos

import (
	"sync"
	"errors"
	"math/rand"
)


type StrategyUtils interface {
	SelectServer() (string, error)
	AddServer(server string)
	RemoveServer(server string)
}

// =================================== //
// ========= Round Robin strategy ==== //
// =================================== //

type RoundRobin struct {
	servers []string
	index   int
	mu      sync.Mutex
}

func NewRoundRobin(servers []string) *RoundRobin {
	return &RoundRobin{servers: servers, index: 0}
}

func (rr *RoundRobin) SelectServer() (string, error) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	if len(rr.servers) == 0 {
		return "", errors.New("no servers available")
	}
	server := rr.servers[rr.index]
	rr.index = (rr.index + 1) % len(rr.servers)
	return server, nil
}

func (rr *RoundRobin) AddServer(server string) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	rr.servers = append(rr.servers, server)
}

func (rr *RoundRobin) RemoveServer(server string) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	for i, s := range rr.servers {
		if s == server {
			rr.servers = append(rr.servers[:i], rr.servers[i+1:]...)
			break
		}
	}
}

// =================================== //
// ========= Random strategy ========= //
// =================================== //

type RandomLB struct {
	servers []string
	mu      sync.Mutex
}

func NewRandomLB(servers []string) *RandomLB {
	return &RandomLB{servers: servers}
}

func (r *RandomLB) SelectServer() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.servers) == 0 {
		return "", errors.New("no servers available")
	}
	index := rand.Intn(len(r.servers))
	return r.servers[index], nil
}

func (r *RandomLB) AddServer(server string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.servers = append(r.servers, server)
}

func (r *RandomLB) RemoveServer(server string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, s := range r.servers {
		if s == server {
			r.servers = append(r.servers[:i], r.servers[i+1:]...)
			break
		}
	}
}

// ============================================= //
// ========= Least connection strategy ========= //
// ============================================= //

type LeastConnLB struct {
	servers                 []string
	activeConnections       map[string]int
	mu		                sync.Mutex
}

func NewLeastConnLB(servers []string) *LeastConnLB {
	activeConnections := make(map[string]int)
	for _, server := range servers {
		// initially number of connections for each server is 0
		activeConnections[server] = 0
	}
	return &LeastConnLB{servers: servers, activeConnections: activeConnections}
}


func (lc *LeastConnLB) SelectServer() (string, error) {

	// the actual strategy is implemented here
	lc.mu.Lock()
	defer lc.mu.Unlock()
	if len(lc.servers) == 0 {
		return "", errors.New("no servers available")
	}

	// find server with the least active connections
	var selectedServer string
	minConnections := int(^uint(0) >> 1) // setting minConnections to max int value

	for _, server := range lc.servers {
		if lc.activeConnections[server] < minConnections {
			selectedServer = server
			minConnections = lc.activeConnections[server]
		}
	}

	// Increment the active connection count for the selected server
	lc.activeConnections[selectedServer]++
	return selectedServer, nil
}

func (lc *LeastConnLB) AddServer(server string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.servers = append(lc.servers, server)
	lc.activeConnections[server] = 0
}

func (lc *LeastConnLB) RemoveServer(server string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	for i, s := range lc.servers {
		if s == server {
			lc.servers = append(lc.servers[:i], lc.servers[i+1:]...)
			delete(lc.activeConnections, server)    // remove the entry from the activeConnections map
			break
		}
	}
}

// ================================ //
// ===== LB strategy selector ===== //
// ================================ //

type LBStrategy string

const (
	RoundRobinStrategy              LBStrategy = "round_robin"
	RandomStrategy                  LBStrategy = "random"
	LeastConnectionStrategy         LBStrategy = "least_conn"
)

func NewStrategy(strategy LBStrategy, servers []string) StrategyUtils {
	switch strategy {
	case RoundRobinStrategy:
		return NewRoundRobin(servers)
	case RandomStrategy:
		return NewRandomLB(servers)
	case LeastConnectionStrategy:
		return NewLeastConnLB(servers)
	default:
		return nil
	}
}



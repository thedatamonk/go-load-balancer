package main

import (
	"errors"
	"fmt"
	"math/rand"
)


type LBStrategy interface {
	SelectServer(lb *LoadBalancer) (*Server, error)      // select the next available server based on the strategy
}

// =================================== //
// ========= Round Robin strategy ==== //
// =================================== //

type RoundRobin struct {
	currServerIndex   int
}


func (s *RoundRobin) SelectServer(lb *LoadBalancer) (*Server, error) {
	
	if len(lb.servers) == 0 {
		return nil, errors.New("no servers available")
	}

	server := lb.servers[s.currServerIndex]
	s.currServerIndex = (s.currServerIndex + 1) % len(lb.servers)

	return server, nil
}


// =================================== //
// ========= Random strategy ========= //
// =================================== //

type RandomLB struct {}

func (s *RandomLB) SelectServer(lb *LoadBalancer) (*Server, error) {
	if len(lb.servers) == 0 {
		return nil, errors.New("no servers available")
	}

	index := rand.Intn(len(lb.servers))
	return lb.servers[index], nil
}



// ============================================= //
// ========= Least connection strategy ========= //
// ============================================= //

type LeastConnLB struct {}

func (s *LeastConnLB) SelectServer(lb *LoadBalancer) (*Server, error) {

	if len(lb.servers) == 0 {
		return nil, errors.New("no servers available")
	}
	// find server with the least active connections
	var selectedServer *Server
	minConnections := int(^uint(0) >> 1) // setting minConnections to max int value

	for _, server := range lb.servers {
		server_host := server.URL.Host
		if lb.connections[server_host] < minConnections {
			selectedServer = server
			minConnections = lb.connections[server_host]
		}
	}

	// Increment the active connection count for the selected server
	lb.connections[selectedServer.URL.Host]++
	return selectedServer, nil
}


// ================================ //
// ===== LB strategy selector ===== //
// ================================ //


func NewStrategy(strategyName string) (LBStrategy, error) {
	switch strategyName {
	case "round-robin":
		return &RoundRobin{}, nil
	case "random":
		return &RandomLB{}, nil
	case "least-connections":
		return &LeastConnLB{}, nil
	default:
		return nil, fmt.Errorf("unknown strategy: %s", strategyName)
	}
}



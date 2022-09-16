package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/praveen4g0/loadbalancer/server"
)

var (
	ServerList = []*server.Server{
		server.NewServer("server-1", "http://127.0.0.1:5001"),
		server.NewServer("server-2", "http://127.0.0.1:5002"),
		server.NewServer("server-3", "http://127.0.0.1:5003"),
		server.NewServer("server-4", "http://127.0.0.1:5004"),
		server.NewServer("server-5", "http://127.0.0.1:5005"),
	}
	lastServedIndex = 0
)

func startHealthCheck() {
	s := gocron.NewScheduler(time.Local)
	for _, host := range ServerList {
		_, err := s.Every(2).Seconds().Do(func(s *server.Server) {
			healthy := s.CheckHealth()
			if healthy {
				log.Printf("'%s' is healthy!", s.Name)
			} else {
				log.Printf("'%s' is not healthy", s.Name)
			}
		}, host)
		if err != nil {
			log.Fatalln(err)
		}
	}
	s.StartAsync()
}

func main() {
	http.HandleFunc("/", forwardRequest)
	go startHealthCheck()
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func forwardRequest(res http.ResponseWriter, req *http.Request) {
	server, err := getHealthyServer()
	if err != nil {
		http.Error(res, "Couldn't process request: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	server.ReverseProxy.ServeHTTP(res, req)
}

func getHealthyServer() (*server.Server, error) {
	for i := 0; i < len(ServerList); i++ {
		server := getServer()
		if server.Health {
			return server, nil
		}
	}
	return nil, fmt.Errorf("no healthy hosts")
}

func getServer() *server.Server {
	nextIndex := (lastServedIndex + 1) % len(ServerList)
	server := ServerList[nextIndex]
	lastServedIndex = nextIndex
	return server
}

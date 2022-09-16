package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	servername = flag.String("server", "server-1", "server-name")
	port       = flag.Int("port", 5000, "port")
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(fmt.Sprintf("Hello %s", *servername))
	w.Write([]byte(fmt.Sprintf("Hello from %s", *servername)))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println(fmt.Sprintf("%s, Healthy", *servername))
	w.Write([]byte(fmt.Sprintf("%s, Healthy", *servername)))
}

func main() {
	flag.Parse()
	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	r.HandleFunc("/", defaultHandler).Methods("GET")
	r.HandleFunc("/", healthCheck).Methods("HEAD")
	r.HandleFunc("/", notAllowedHandler)

	log.Println("Application started \n server: "+*servername+"\t port: ", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}

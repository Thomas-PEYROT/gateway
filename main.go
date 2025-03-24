package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func startHTTPServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.RequestURI[1:], "/")
		microserviceName := parts[0]
		query := strings.Join(parts[1:], "/")

		instances, exists := RegisteredMicroservices[microserviceName]
		if !exists || len(instances) == 0 {
			http.Error(w, "Microservice not found", http.StatusNotFound)
			return
		}

		// Extract instances
		keys := make([]string, 0, len(instances))
		for k := range instances {
			keys = append(keys, k)
		}

		// Get a random instance
		selectedInstance := keys[rand.Intn(len(keys))]
		newURL := fmt.Sprintf("http://localhost:%d/%s", instances[selectedInstance].Port, query)
		fmt.Println("Forwarding request to:", newURL)

		// Create a new request
		req, err := http.NewRequest(r.Method, newURL, r.Body)
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		// Copy headers from the original request
		for name, values := range r.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}

		// Execute the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Failed to forward request", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}
		w.WriteHeader(resp.StatusCode)

		// Copy response body
		io.Copy(w, resp.Body)
	})

	fmt.Println("Serveur HTTP démarré sur le port 8080")
	http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), nil)
}

func main() {
	go startHTTPServer()
	godotenv.Load()
	for {
		data, err := FetchMicroservices(fmt.Sprintf("%s/microservices", os.Getenv("DISCOVERY_SERVER_URL")))
		if err != nil {
			fmt.Println("Erreur lors de la récupération des microservices:", err)
			return
		}
		RegisteredMicroservices = data
		log.Println("Refreshed microservices from discovery server")
		time.Sleep(10 * time.Second)
	}

}

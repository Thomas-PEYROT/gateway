package main

import (
	"fmt"
	"github.com/joho/godotenv"
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

		instances := RegisteredMicroservices[microserviceName]

		// Extract instances
		keys := make([]string, 0, len(instances))
		for k := range instances {
			keys = append(keys, k)
		}

		// Get a random instance
		selectedInstance := keys[rand.Intn(len(keys))]
		newURL := fmt.Sprintf("http://localhost:%d/%s", instances[selectedInstance].Port, query)
		fmt.Println(newURL)
		http.Redirect(w, r, newURL, http.StatusTemporaryRedirect)
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

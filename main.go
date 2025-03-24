package main

import (
	"fmt"
	"math/rand"
	"net/http"
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
		fmt.Println("Clé aléatoire :", selectedInstance)
		newURL := fmt.Sprintf("http://localhost:%d/%s", instances[selectedInstance].Port, query)
		fmt.Println(newURL)
		http.Redirect(w, r, newURL, http.StatusTemporaryRedirect)
	})

	fmt.Println("Serveur HTTP démarré sur le port 8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	go startHTTPServer()
	for {
		apiURL := "http://localhost:1234/microservices"
		data, err := FetchMicroservices(apiURL)
		if err != nil {
			fmt.Println("Erreur lors de la récupération des microservices:", err)
			return
		}
		RegisteredMicroservices = data
		time.Sleep(5 * time.Second)
	}

}

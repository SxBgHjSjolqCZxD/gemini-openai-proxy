package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhu327/gemini-openai-proxy/api"
)

func main() {
	// Define a flag for the port
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	// Create a new Gin router
	router := gin.Default()
	api.Register(router)

	// Start the keep-alive routine
	go keepAlive(fmt.Sprintf("http://localhost:%d", *port))

	// Run the server on the specified port
	err := router.Run(fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}
}

func keepAlive(url string) {
	ticker := time.NewTicker(10 * time.Minute) // Ping every 10 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Keep-alive ping failed: %v\n", err)
			} else {
				log.Printf("Keep-alive ping successful: %s\n", resp.Status)
				_ = resp.Body.Close()
			}
		}
	}
}

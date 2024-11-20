package main

import (
	"fmt"
	"log"
	"net/http"
	"urlshortner/cache"
	"urlshortner/ratelimiter"
	"urlshortner/router"

	"github.com/rs/cors"
)

func main() {
	fmt.Println("Hello World")

	cache.InitRedis()
	ratelimiter.InitRedis()

	r := router.Router()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Fatal(http.ListenAndServe(":8080", handler))

}

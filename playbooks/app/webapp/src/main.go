package main

import (
	"log"
	"webapp/src/config"
)

func main() {
	// ルータ起動
	router := initRouter()
	env := config.GetEnv()
	addr := env.ServerHost + ":" + env.ServerPort
	log.Printf("Server starting on http://%s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

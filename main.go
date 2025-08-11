package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nyradhr/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	err = cfg.SetUser("nyradhr")
	if err != nil {
		log.Fatalf("Error setting user name in config file: %v", err)
	}
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}

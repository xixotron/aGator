package main

import (
	"fmt"
	"log"

	"github.com/xixotron/aleyGator/internal/config"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading initial config: %v", err)
	}

	fmt.Printf("Read config: %+v\n", cfg)

	cfg.SetUser("xixotron")

	cfg2, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading modified config: %v", err)
	}

	fmt.Printf("Read modified config: %+v\n", cfg2)
}

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/shashimalcse/tiny-is/internal/config"
	"github.com/shashimalcse/tiny-is/internal/server"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	defaultConfigPath := filepath.Join(cwd, "config", "config.yaml")
	configPath := flag.String("config", defaultConfigPath, "path to config file")
	flag.Parse()
	absConfigPath, err := filepath.Abs(*configPath)
	if err != nil {
		log.Fatalf("Failed to resolve absolute path of config file: %v", err)
	}
	cfg, err := config.LoadConfig(absConfigPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	server.StartServer(cfg)
}

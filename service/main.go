package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml"
)

const configFile = "config.toml"

var defaultEndpoint = "http://thesimpsonsquoteapi.glitch.me:443/quotes"

type Config struct {
	Endpoint string `toml:"endpoint"`
}
type SimpsonsQuote struct {
	Quote              string `json:"quote"`
	Character          string `json:"character"`
	Image              string `json:"image"`
	CharacterDirection string `json:"characterDirection"`
}

func loadConfig() (Config, error) {
	config := Config{}
	snapDir := os.Getenv("SNAP_DATA")
	data, err := os.ReadFile(filepath.Join(snapDir, configFile))
	if err != nil {
		return config, err
	}
	err = toml.Unmarshal(data, &config)
	return config, err
}

func queryAPI(endpoint string) {
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Printf("error querying API: %v", err)
		return
	}
	if resp.StatusCode > 399 {
		log.Fatalf("unexpected response %v - %v", resp.StatusCode, resp.Status)
		return
	}
	defer resp.Body.Close()

	var q []SimpsonsQuote

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
		return
	}

	err = json.Unmarshal(body, &q)
	if err != nil {
		log.Printf("error reading quotes API - %v", err)
		return
	}
	quote := q[0]
	fmt.Printf("\"%s\" - %s\n", quote.Quote, quote.Character)
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Printf("error loading config: %v", err)
		log.Printf("using default endpoint: %v", defaultEndpoint)
		config.Endpoint = defaultEndpoint
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	queryAPI(config.Endpoint)

	for {
		select {
		case <-ticker.C:
			queryAPI(config.Endpoint)
		}
	}
}

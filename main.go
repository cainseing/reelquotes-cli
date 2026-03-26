package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Settings struct {
		Theme     string `json:"theme"`
		Frequency string `json:"frequency"`
		Censored  bool   `json:"censored"`
	}
	FavMovies []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"favMovies"`
}

type APIResponse struct {
	Dialogue []struct {
		Line string `json:"line"`
	} `json:"dialogue"`
	Media struct {
		Title string `json:"title"`
	} `json:"media"`
}

func main() {
	configPath := filepath.Join(os.Getenv("HOME"), ".reelquotes.conf")

	if len(os.Args) > 2 && os.Args[1] == "import" {
		data, err := os.ReadFile(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		_ = os.WriteFile(configPath, data, 0644)
		fmt.Println("Config imported")
		return
	}

	fetchQuote(configPath)
}

func fetchQuote(configPath string) {
	apiURL := "https://api.reelquotes.app/quote?singular=true"

	if file, err := os.Open(configPath); err == nil {
		var cfg Config
		if err := json.NewDecoder(file).Decode(&cfg); err == nil && len(cfg.FavMovies) > 0 {
			randomMovie := cfg.FavMovies[rand.Intn(len(cfg.FavMovies))]
			apiURL = fmt.Sprintf("%s&type=media&identifier=%s", apiURL, randomMovie.ID)
		}
		if cfg.Settings.Censored {
			apiURL = fmt.Sprintf("%s&censored=true", apiURL)
		}
		file.Close()
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("x-reel-quotes", "reel-quotes-cli")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}
	defer resp.Body.Close()

	var data APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || len(data.Dialogue) == 0 {
		return
	}

	fmt.Printf("\n\033[3m\"%s\"\033[0m\n", data.Dialogue[0].Line)
	fmt.Printf("\033[1;34m— %s\033[0m\n\n", data.Media.Title)
}

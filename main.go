package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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
	configPath := filepath.Join(os.Getenv("HOME"), ".reelquotes.json")
	profilePath := getShellProfile()

	title := flag.String("title", "", "Search for a quote from a specific film")
	flag.Parse()

	if flag.NArg() > 0 {
		switch flag.Arg(0) {
		case "import":
			if flag.NArg() < 2 {
				return
			}
			importConfig(flag.Arg(1), configPath)
			return
		case "install":
			manageShell(profilePath, true)
			return
		case "uninstall":
			manageShell(profilePath, false)
			return
		}
	}

	fetchQuote(configPath, *title)
}

func getShellProfile() string {
	home := os.Getenv("HOME")
	shell := os.Getenv("SHELL")

	profile := filepath.Join(home, ".zshrc")

	if strings.Contains(shell, "bash") {
		profile = filepath.Join(home, ".bashrc")
		if _, err := os.Stat(profile); os.IsNotExist(err) {
			profile = filepath.Join(home, ".bash_profile")
		}
	}
	return profile
}

func importConfig(src, dest string) {
	data, err := os.ReadFile(src)
	if err != nil {
		return
	}
	_ = os.WriteFile(dest, data, 0644)
	fmt.Println("Config imported successfully.")
}

func manageShell(profilePath string, install bool) {
	marker := "reelquotes # added by reelquotes-cli"

	content, err := os.ReadFile(profilePath)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	found := false

	for _, line := range lines {
		if strings.Contains(line, "reelquotes") {
			found = true
			if !install {
				continue
			}
		}
		newLines = append(newLines, line)
	}

	if install && !found {
		newLines = append(newLines, marker)
		fmt.Printf("Added reelquotes to %s\n", filepath.Base(profilePath))
	} else if !install {
		fmt.Printf("Removed reelquotes from %s\n", filepath.Base(profilePath))
	} else if install && found {
		fmt.Println("Already installed in your shell profile.")
		return
	}

	_ = os.WriteFile(profilePath, []byte(strings.Join(newLines, "\n")), 0644)
}

func fetchQuote(configPath string, search string) {
	apiURL := "https://api.reelquotes.app/quote?singular=true"

	if search != "" {
		apiURL = fmt.Sprintf("%s&type=media&q=%s", apiURL, url.QueryEscape(search))
	} else if file, err := os.Open(configPath); err == nil {
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

	client := &http.Client{Timeout: 500 * time.Millisecond}
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

# reelquotes-cli

A lightweight, Go-based CLI tool that displays a random film quote every time you open your terminal. 

Built for speed, `reelquotes` integrates directly with your shell profile to keep your workspace inspired by your favorite cinema.

### Installation

#### Via Homebrew (Recommended)

    brew tap cainseing/tap
    brew install reelquotes


#### Manual Build

    go build -o reelquotes main.go
    sudo mv reelquotes /usr/local/bin/

### Setup

#### Enable Startup Quotes
To have a random quote displayed every time you open a new terminal session:

    reelquotes install

*To stop quotes from appearing, simply run `reelquotes uninstall`.*

#### Import Your Favorites
If you have favourites setup from the [ReelQuotes](https://reelquotes.app), export your preferences and import it with the CLI:

    reelquotes import ~/Downloads/settings.json


### Usage

| Command | Description |
 | :--- | :--- |
 | `reelquotes` | Fetches and displays a random quote. |
 | `reel quotes import [path]` | Import your preferences from the web app. |
 | `reelquotes install` | Automatically adds `reelquotes` to your shell profile. |
 | `reelquotes uninstall` | Removes `reelquotes` from your shell profile. |

---

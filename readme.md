# reelquotes-cli

A lightweight Go-based CLI tool that displays a random film quote every time you open your terminal. It supports personalizing quotes based on a curated list of your favorite movies.

## Features

* **Fast & Native:** Written in Go with zero runtime dependencies.
* **Personalized:** Import settings from the https://reelquotes.app web application to receive quotes from your favourite films.

## Installation

Install the `reelquotes` binary to your system path.

### Homebrew

    brew tap cainseing/tap &&
    brew install reelquotes

### Display quote on terminal launch
Add the following to the end of your `~/.zshrc` or `~/.bashrc`

    # Display random movie quote
    reelquotes
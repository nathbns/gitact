package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	// Check if --repos flag is used
	if len(os.Args) >= 3 && os.Args[1] == "--repos" {
		username := strings.TrimSpace(os.Args[2])
		if username == "" {
			fmt.Fprintf(os.Stderr, "error: username can't be empty\n")
			os.Exit(1)
		}
		showPublicRepos(username)
		return
	}

	username := strings.TrimSpace(os.Args[1])

	// flags
	switch username {
	case "-h", "--help":
		showHelp()
		return
	case "-v", "--version":
		showVersion()
		return
	case "--repos":
		fmt.Fprintf(os.Stderr, "error: --repos requires a username\n")
		fmt.Fprintf(os.Stderr, "usage: %s --repos <username>\n", os.Args[0])
		os.Exit(1)
	}

	if username == "" {
		fmt.Fprintf(os.Stderr, "error: username can't be empty\n")
		os.Exit(1)
	}

	// Check rate limit before starting
	if err := checkRateLimit(); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Rate limit warning: %v\n", err)
		fmt.Fprintf(os.Stderr, "💡 Set GITHUB_TOKEN environment variable for higher limits\n\n")
	}

	// init model bubble tea with new modernized UI
	initialModel := NewModel(username)

	p := tea.NewProgram(initialModel, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error during the launch : %v", err)
		os.Exit(1)
	}
}

func showPublicRepos(username string) {
	fmt.Printf("🔍 Fetching public repositories for user: %s\n", username)

	// Fetch public repositories
	publicRepos, err := fetchPublicRepos(username)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error fetching public repositories: %v\n", err)
		os.Exit(1)
	}

	// Display statistics and repositories
	calculatePublicReposStats(publicRepos)
	printPublicRepos(publicRepos)
}

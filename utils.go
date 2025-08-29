package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// format
func formatNumber(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	} else if n >= 1000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}

// score
func getGrade(stats GitHubStats) string {
	if stats.TotalEvents == 0 {
		return "F"
	}
	score := float64(stats.PushEvents)*1.0 +
		float64(stats.PullRequestEvents)*3.0 +
		float64(stats.CreateEvents)*1.0 +
		float64(stats.IssueEvents)*1.5 +
		float64(stats.WatchEvents)*0.5

	switch {
	case score >= 100:
		return "S+"
	case score >= 70:
		return "S"
	case score >= 40:
		return "A+"
	case score >= 25:
		return "A"
	case score >= 15:
		return "B+"
	case score >= 8:
		return "B"
	case score >= 3:
		return "C"
	case score >= 1:
		return "D"
	default:
		return "F"
	}
}

func formatEventShort(event GitHubEvent) string {
	switch event.Type {
	case "PushEvent":
		return fmt.Sprintf("Pushed to %s", event.Repo.Name)
	case "IssuesEvent":
		return fmt.Sprintf("Issue in %s", event.Repo.Name)
	case "WatchEvent":
		return fmt.Sprintf("Starred %s", event.Repo.Name)
	case "ForkEvent":
		return fmt.Sprintf("Forked %s", event.Repo.Name)
	case "CreateEvent":
		return fmt.Sprintf("Created %s", event.Repo.Name)
	case "PullRequestEvent":
		return fmt.Sprintf("PR in %s", event.Repo.Name)
	default:
		return fmt.Sprintf("%s in %s",
			strings.TrimSuffix(event.Type, "Event"), event.Repo.Name)
	}
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else {
			return fmt.Errorf("no utilitary to copy found")
		}
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("OS not supported: %s", runtime.GOOS)
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// Fonctions d'aide et d'information
func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <username>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "   or: %s --repos <username>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Example: %s octocat\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "         %s --repos octocat\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nuse '%s --help' for more informations.\n", os.Args[0])
}

func showHelp() {
	fmt.Printf("GitHub Activity CLI - Modern Edition\n\n")
	fmt.Printf("Description:\n")
	fmt.Printf("Modern interactive CLI to explore GitHub profiles, repositories, and activity.\n")
	fmt.Printf("Built with Charm's Bubbles UI components for a delightful terminal experience.\n\n")
	fmt.Printf("Usage:\n")
	fmt.Printf("  %s <username>        Interactive dashboard with multiple views\n", os.Args[0])
	fmt.Printf("  %s --repos <username> Detailed repository listing\n\n", os.Args[0])
	fmt.Printf("Options:\n")
	fmt.Printf("  -h, --help     Show this help message\n")
	fmt.Printf("  -v, --version  Show version information\n")
	fmt.Printf("  --repos        Display all public repositories with detailed statistics\n\n")
	fmt.Printf("GitHub Token (Recommended):\n")
	fmt.Printf("  Set GITHUB_TOKEN environment variable to avoid rate limits:\n")
	fmt.Printf("  • Without token: 60 requests/hour\n")
	fmt.Printf("  • With token: 5,000 requests/hour\n")
	fmt.Printf("  \n")
	fmt.Printf("  export GITHUB_TOKEN=your_token_here\n")
	fmt.Printf("  %s karpathy\n\n", os.Args[0])
	fmt.Printf("Interactive Dashboard Views:\n")
	fmt.Printf(" Repository List  - Browse repos with search functionality\n")
	fmt.Printf(" Table View       - Detailed tabular data (stars, forks, language)\n")
	fmt.Printf(" Statistics       - Comprehensive stats and insights\n")
	fmt.Printf("  Activity Feed    - Recent GitHub activity timeline\n\n")
	fmt.Printf("Navigation:\n")
	fmt.Printf("  ↑/↓ or j/k    Navigate items\n")
	fmt.Printf("  ←/→ or h/l    Switch between views\n")
	fmt.Printf("  tab           Next view\n")
	fmt.Printf("  /             Search repositories (in list view)\n")
	fmt.Printf("  enter         Select item\n")
	fmt.Printf("  c             Copy git clone command\n")
	fmt.Printf("  x             Copy repository URL\n")
	fmt.Printf("  o             Open repository in browser\n")
	fmt.Printf("  r             Refresh all data\n")
	fmt.Printf("  ?             Toggle help\n")
	fmt.Printf("  q/esc         Quit\n\n")
	fmt.Printf("Examples:\n")
	fmt.Printf("  %s karpathy          # Explore karpathy's ML repositories\n", os.Args[0])
	fmt.Printf("  %s --repos torvalds  # List all of torvalds' projects\n", os.Args[0])
	fmt.Printf("  GITHUB_TOKEN=xxx %s octocat  # With authentication\n\n", os.Args[0])
}

func showVersion() {
	fmt.Printf("gitact CLI v1.0.0\n")
	fmt.Printf("Created with Golang\n")

}

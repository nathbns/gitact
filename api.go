package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func fetchGitHubActivity(username string) ([]GitHubEvent, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating the request: %v", err)
	}

	req.Header.Set("User-Agent", "gh-act-cli/1.0")

	// Add GitHub token if available
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request http error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("user '%s' not found", username)
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http error %d", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture réponse: %v", err)
	}

	var events []GitHubEvent
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("erreur parsing JSON: %v", err)
	}

	return events, nil
}

func calculateStats(events []GitHubEvent) GitHubStats {
	var stats GitHubStats

	for _, event := range events {
		stats.TotalEvents++

		switch event.Type {
		case "PushEvent":
			stats.PushEvents++
		case "IssuesEvent":
			stats.IssueEvents++
		case "WatchEvent":
			stats.WatchEvents++
		case "ForkEvent":
			stats.ForkEvents++
		case "CreateEvent":
			stats.CreateEvents++
		case "DeleteEvent":
			stats.DeleteEvents++
		case "PullRequestEvent":
			stats.PullRequestEvents++
		case "ReleaseEvent":
			stats.ReleaseEvents++
		case "PublicEvent":
			stats.PublicEvents++
		default:
			stats.OtherEvents++
		}
	}
	return stats
}

func getTopRepos(events []GitHubEvent) []RepoInfo {
	repoCount := make(map[string]int)
	reposLastAct := make(map[string]time.Time)

	// counting num of act by repos
	for _, events := range events {
		if events.Repo.Name != "" {
			repoCount[events.Repo.Name]++
			if events.CreatedAt.After(reposLastAct[events.Repo.Name]) {
				reposLastAct[events.Repo.Name] = events.CreatedAt
			}
		}
	}

	// convert to slice
	var repos []RepoInfo
	for name, count := range repoCount {
		repo := RepoInfo{
			Name:         name,
			URL:          fmt.Sprintf("https://github.com/%s", name),
			CloneURL:     fmt.Sprintf("https://github.com/%s.git", name),
			Count:        count,
			LastActivity: reposLastAct[name],
		}
		repos = append(repos, repo)
	}

	// bull sorting
	for i := 0; i < len(repos)-1; i++ {
		for j := 0; j < len(repos)-i-1; j++ {
			if repos[j].Count < repos[j+1].Count {
				repos[j], repos[j+1] = repos[j+1], repos[j]
			}
		}
	}

	return repos
}

func printTopRepo(repos []RepoInfo) {
	fmt.Printf("\n=== Top Repositories by Activity (%d total) ===\n", len(repos))
	for i, repo := range repos {
		fmt.Printf("\n%d. %s\n", i+1, repo.Name)
		fmt.Printf("   📊 Activity Events: %d\n", repo.Count)
		fmt.Printf("   🔗 URL: %s\n", repo.URL)
		fmt.Printf("   📅 Last Activity: %s\n", repo.LastActivity.Format("2006-01-02 15:04:05"))
		if repo.Description != "" {
			fmt.Printf("   📋 Description: %s\n", repo.Description)
		}
	}
}

func fetchPublicRepos(username string) ([]PublicRepo, error) {
	var allRepos []PublicRepo
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("https://api.github.com/users/%s/repos?type=public&sort=stars&direction=desc&per_page=%d&page=%d",
			username, perPage, page)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating the request: %v", err)
		}

		req.Header.Set("User-Agent", "gh-act-cli/1.0")
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		// Add GitHub token if available
		if token := os.Getenv("GITHUB_TOKEN"); token != "" {
			req.Header.Set("Authorization", "token "+token)
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request http error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 404 {
			return nil, fmt.Errorf("user '%s' not found", username)
		} else if resp.StatusCode != 200 {
			return nil, fmt.Errorf("http error %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response: %v", err)
		}

		var repos []PublicRepo
		if err := json.Unmarshal(body, &repos); err != nil {
			return nil, fmt.Errorf("error parsing JSON: %v", err)
		}

		// Si aucun repo n'est retourné, on a atteint la fin
		if len(repos) == 0 {
			break
		}

		// Filter only public repositories and add to collection
		for _, repo := range repos {
			if !repo.Private {
				allRepos = append(allRepos, repo)
			}
		}

		// Si moins de repos que demandé, c'est la dernière page
		if len(repos) < perPage {
			break
		}

		page++
	}

	return allRepos, nil
}

// checkRateLimit checks GitHub API rate limit
func checkRateLimit() error {
	url := "https://api.github.com/rate_limit"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating rate limit request: %v", err)
	}

	req.Header.Set("User-Agent", "gh-act-cli/1.0")

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error checking rate limit: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("rate limit check failed with status: %d", resp.StatusCode)
	}

	var rateLimit struct {
		Resources struct {
			Core struct {
				Limit     int `json:"limit"`
				Remaining int `json:"remaining"`
				Reset     int `json:"reset"`
			} `json:"core"`
		} `json:"resources"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading rate limit response: %v", err)
	}

	if err := json.Unmarshal(body, &rateLimit); err != nil {
		return fmt.Errorf("error parsing rate limit response: %v", err)
	}

	remaining := rateLimit.Resources.Core.Remaining
	limit := rateLimit.Resources.Core.Limit

	if remaining < 10 {
		resetTime := time.Unix(int64(rateLimit.Resources.Core.Reset), 0)
		return fmt.Errorf("rate limit almost exhausted: %d/%d remaining, resets at %v",
			remaining, limit, resetTime.Format("15:04:05"))
	}

	fmt.Printf("🔄 GitHub API Rate Limit: %d/%d requests remaining\n", remaining, limit)
	return nil
}

func printPublicRepos(repos []PublicRepo) {
	fmt.Printf("\n=== Public Repositories (%d total) ===\n", len(repos))

	if len(repos) == 0 {
		fmt.Println("No public repositories found.")
		return
	}

	// Sort repos by stars (descending) to ensure correct order
	for i := 0; i < len(repos)-1; i++ {
		for j := 0; j < len(repos)-i-1; j++ {
			if repos[j].Stars < repos[j+1].Stars {
				repos[j], repos[j+1] = repos[j+1], repos[j]
			}
		}
	}

	for i, repo := range repos {
		fmt.Printf("\n%d. %s\n", i+1, repo.FullName)
		fmt.Printf("   ⭐ Stars: %d | 🍴 Forks: %d\n", repo.Stars, repo.Forks)
		if repo.Language != "" {
			fmt.Printf("   📝 Language: %s\n", repo.Language)
		}
		if repo.Description != "" {
			fmt.Printf("   📋 Description: %s\n", repo.Description)
		}
		fmt.Printf("   🔗 URL: %s\n", repo.URL)
		fmt.Printf("   📅 Created: %s | Updated: %s\n",
			repo.CreatedAt.Format("2006-01-02"),
			repo.UpdatedAt.Format("2006-01-02"))
	}

	// Show summary at the end
	totalStars := 0
	for _, repo := range repos {
		totalStars += repo.Stars
	}
	fmt.Printf("\n📊 Summary: %d repositories with %d total stars\n", len(repos), totalStars)
}

func calculatePublicReposStats(repos []PublicRepo) {
	if len(repos) == 0 {
		fmt.Println("\n=== Public Repository Statistics ===")
		fmt.Println("No public repositories found.")
		return
	}

	totalStars := 0
	totalForks := 0
	languageCount := make(map[string]int)
	mostStarredRepo := repos[0]
	mostForkedRepo := repos[0]

	for _, repo := range repos {
		totalStars += repo.Stars
		totalForks += repo.Forks

		if repo.Language != "" {
			languageCount[repo.Language]++
		}

		if repo.Stars > mostStarredRepo.Stars {
			mostStarredRepo = repo
		}

		if repo.Forks > mostForkedRepo.Forks {
			mostForkedRepo = repo
		}
	}

	fmt.Printf("\n=== Public Repository Statistics ===\n")
	fmt.Printf("📊 Total Repositories: %d\n", len(repos))
	fmt.Printf("⭐ Total Stars: %d\n", totalStars)
	fmt.Printf("🍴 Total Forks: %d\n", totalForks)

	if len(repos) > 0 {
		fmt.Printf("📈 Average Stars per Repository: %.1f\n", float64(totalStars)/float64(len(repos)))
		fmt.Printf("📈 Average Forks per Repository: %.1f\n", float64(totalForks)/float64(len(repos)))
	}

	fmt.Printf("\n🏆 Most Starred Repository: %s (%d stars)\n", mostStarredRepo.FullName, mostStarredRepo.Stars)
	fmt.Printf("🏆 Most Forked Repository: %s (%d forks)\n", mostForkedRepo.FullName, mostForkedRepo.Forks)

	if len(languageCount) > 0 {
		fmt.Printf("\n📝 Programming Languages Used:\n")
		for lang, count := range languageCount {
			fmt.Printf("   - %s: %d repositories\n", lang, count)
		}
	}
}

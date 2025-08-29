package main

import "time"

type GitHubEvent struct {
	Type      string    `json:"type"`
	Actor     Actor     `json:"actor"`
	Repo      Repo      `json:"repo"`
	Payload   Payload   `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type Actor struct {
	Login string `json:"login"`
}

type Repo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Payload struct {
	Action      string       `json:"action,omitempty"`
	RefType     string       `json:"ref_type,omitempty"`
	Ref         string       `json:"ref,omitempty"`
	Commits     []Commit     `json:"commits,omitempty"`
	Issue       *Issue       `json:"issue,omitempty"`
	PullRequest *PullRequest `json:"pull_request,omitempty"`
}

type Commit struct {
	SHA     string `json:"sha"`
	Message string `json:"message"`
}

type Issue struct {
	Title string `json:"title"`
	State string `json:"state"`
}

type PullRequest struct {
	Title string `json:"title"`
	State string `json:"state"`
}

type GitHubStats struct {
	PushEvents        int
	IssueEvents       int
	WatchEvents       int
	ForkEvents        int
	CreateEvents      int
	DeleteEvents      int
	PullRequestEvents int
	ReleaseEvents     int
	PublicEvents      int
	OtherEvents       int
	TotalEvents       int
}

type RepoInfo struct {
	Name         string
	URL          string
	CloneURL     string
	Count        int
	LastActivity time.Time
	Description  string
}

type PublicRepo struct {
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	URL         string    `json:"html_url"`
	CloneURL    string    `json:"clone_url"`
	Stars       int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	Language    string    `json:"language"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Private     bool      `json:"private"`
}

type NotificationMsg struct {
	message   string
	isSuccess bool
}

type ClearNotificationMsg struct{}

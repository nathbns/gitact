package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Key bindings
type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Help    key.Binding
	Quit    key.Binding
	Enter   key.Binding
	Clone   key.Binding
	Copy    key.Binding
	Open    key.Binding
	Search  key.Binding
	Refresh key.Binding
	Tab     key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.Help, k.Quit},
		{k.Enter, k.Clone, k.Copy, k.Open},
		{k.Search, k.Refresh, k.Tab},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "previous tab"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next tab"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Clone: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "clone repo"),
	),
	Copy: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "copy URL"),
	),
	Open: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open in browser"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch view"),
	),
}

// Views
type viewMode int

const (
	repoListView viewMode = iota
	repoTableView
	statsView
	activityView
)

// List item for repositories
type repoItem struct {
	repo PublicRepo
}

func (i repoItem) FilterValue() string { return i.repo.Name }
func (i repoItem) Title() string {
	return fmt.Sprintf("%s ★ %s", i.repo.Name, formatNumber(i.repo.Stars))
}
func (i repoItem) Description() string {
	desc := i.repo.Description
	if desc == "" {
		desc = "No description"
	}
	if len(desc) > 80 {
		desc = desc[:77] + "..."
	}
	return fmt.Sprintf("⑂ %s • %s", formatNumber(i.repo.Forks), desc)
}

// Activity item
type activityItem struct {
	event GitHubEvent
}

func (i activityItem) FilterValue() string { return i.event.Repo.Name }
func (i activityItem) Title() string {
	return formatEventShort(i.event)
}
func (i activityItem) Description() string {
	return i.event.CreatedAt.Format("2006-01-02 15:04")
}

// Model
type Model struct {
	username    string
	events      []GitHubEvent
	repos       []RepoInfo
	publicRepos []PublicRepo
	stats       GitHubStats

	// UI components
	list      list.Model
	table     table.Model
	viewport  viewport.Model
	help      help.Model
	spinner   spinner.Model
	progress  progress.Model
	search    textinput.Model
	paginator paginator.Model

	// State
	currentView  viewMode
	loading      bool
	showHelp     bool
	searchMode   bool
	notification string
	notifSuccess bool
	width        int
	height       int
	ready        bool

	// Data loading state
	reposLoaded  bool
	eventsLoaded bool
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.loadData(),
	)
}

func (m Model) loadData() tea.Cmd {
	return tea.Batch(
		loadReposCmd(m.username),
		loadEventsCmd(m.username),
	)
}

// Commands
type reposLoadedMsg struct {
	repos []PublicRepo
	err   error
}

type eventsLoadedMsg struct {
	events []GitHubEvent
	stats  GitHubStats
	err    error
}

func loadReposCmd(username string) tea.Cmd {
	return func() tea.Msg {
		repos, err := fetchPublicRepos(username)
		return reposLoadedMsg{repos: repos, err: err}
	}
}

func loadEventsCmd(username string) tea.Cmd {
	return func() tea.Msg {
		events, err := fetchGitHubActivity(username)
		if err != nil {
			return eventsLoadedMsg{err: err}
		}
		stats := calculateStats(events)
		return eventsLoadedMsg{events: events, stats: stats, err: nil}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width

		// Update list dimensions with better centering
		headerHeight := 4 // Header takes 3-4 lines
		helpHeight := 3   // Help takes 2-3 lines
		padding := 4      // Left/right padding
		availableHeight := msg.Height - headerHeight - helpHeight - 2

		m.list.SetSize(msg.Width-padding, availableHeight)

		// Update table
		m.updateTableSize()

		// Update viewport with proper sizing
		m.viewport.Width = msg.Width - padding
		m.viewport.Height = availableHeight

		return m, nil

	case reposLoadedMsg:
		m.reposLoaded = true
		if msg.err != nil {
			m.notification = fmt.Sprintf("❌ Error loading repositories: %v", msg.err)
			m.notifSuccess = false
		} else {
			m.publicRepos = msg.repos
			m.updateRepoList()
			m.updateRepoTable()
		}
		m.checkLoadingComplete()
		return m, nil

	case eventsLoadedMsg:
		m.eventsLoaded = true
		if msg.err != nil {
			m.notification = fmt.Sprintf("❌ Error loading activity: %v", msg.err)
			m.notifSuccess = false
		} else {
			m.events = msg.events
			m.stats = msg.stats
			m.updateActivityList()
		}
		m.checkLoadingComplete()
		return m, nil

	case NotificationMsg:
		m.notification = msg.message
		m.notifSuccess = msg.isSuccess
		return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
			return ClearNotificationMsg{}
		})

	case ClearNotificationMsg:
		m.notification = ""
		return m, nil

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		if m.searchMode {
			return m.handleSearchInput(msg)
		}

		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		case key.Matches(msg, keys.Tab):
			m.nextView()
			return m, nil

		case key.Matches(msg, keys.Search):
			if m.currentView == repoListView {
				m.searchMode = true
				m.search.Focus()
				return m, textinput.Blink
			}

		case key.Matches(msg, keys.Refresh):
			m.loading = true
			m.reposLoaded = false
			m.eventsLoaded = false
			m.notification = "Refreshing data..."
			m.notifSuccess = true
			return m, m.loadData()

		case key.Matches(msg, keys.Clone):
			if m.currentView == repoListView && len(m.publicRepos) > 0 {
				selected := m.list.SelectedItem()
				if repoItem, ok := selected.(repoItem); ok {
					return m, m.cloneRepo(repoItem.repo)
				}
			}

		case key.Matches(msg, keys.Copy):
			if m.currentView == repoListView && len(m.publicRepos) > 0 {
				selected := m.list.SelectedItem()
				if repoItem, ok := selected.(repoItem); ok {
					return m, m.copyURL(repoItem.repo)
				}
			}

		case key.Matches(msg, keys.Open):
			if m.currentView == repoListView && len(m.publicRepos) > 0 {
				selected := m.list.SelectedItem()
				if repoItem, ok := selected.(repoItem); ok {
					return m, m.openInBrowser(repoItem.repo)
				}
			}
		}

		// Update current view component
		switch m.currentView {
		case repoListView:
			m.list, cmd = m.list.Update(msg)
		case repoTableView:
			m.table, cmd = m.table.Update(msg)
		case activityView:
			m.list, cmd = m.list.Update(msg)
		case statsView:
			m.viewport, cmd = m.viewport.Update(msg)
		}

		return m, cmd
	}

	// Update spinner if loading
	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) handleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyEsc, tea.KeyCtrlC:
		m.searchMode = false
		m.search.Blur()
		m.search.SetValue("")
		m.updateRepoList() // Reset filter
		return m, nil

	case tea.KeyEnter:
		m.searchMode = false
		m.search.Blur()
		m.filterRepoList(m.search.Value())
		return m, nil
	}

	m.search, cmd = m.search.Update(msg)
	return m, cmd
}

func (m *Model) nextView() {
	switch m.currentView {
	case repoListView:
		m.currentView = repoTableView
	case repoTableView:
		m.currentView = statsView
	case statsView:
		m.currentView = activityView
	case activityView:
		m.currentView = repoListView
	}

	// Update lists based on current view
	switch m.currentView {
	case repoListView:
		m.updateRepoList()
	case activityView:
		m.updateActivityList()
	case statsView:
		m.updateStatsView()
	}
}

func (m *Model) checkLoadingComplete() {
	if m.reposLoaded && m.eventsLoaded {
		m.loading = false
		m.ready = true
	}
}

func (m *Model) updateRepoList() {
	items := make([]list.Item, len(m.publicRepos))
	for i, repo := range m.publicRepos {
		items[i] = repoItem{repo: repo}
	}
	m.list.SetItems(items)
	m.list.Title = fmt.Sprintf("℗ Public Repositories (%d)", len(m.publicRepos))
}

func (m *Model) updateActivityList() {
	items := make([]list.Item, len(m.events))
	for i, event := range m.events {
		items[i] = activityItem{event: event}
	}
	m.list.SetItems(items)
	m.list.Title = fmt.Sprintf("𐧾 Recent Activity (%d events)", len(m.events))
}

func (m *Model) filterRepoList(query string) {
	if query == "" {
		m.updateRepoList()
		return
	}

	var filtered []list.Item
	for _, repo := range m.publicRepos {
		if strings.Contains(strings.ToLower(repo.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(repo.Description), strings.ToLower(query)) {
			filtered = append(filtered, repoItem{repo: repo})
		}
	}
	m.list.SetItems(filtered)
	m.list.Title = fmt.Sprintf("𐧻 Repositories matching '%s' (%d)", query, len(filtered))
}

func (m *Model) updateRepoTable() {
	columns := []table.Column{
		{Title: "Name", Width: 25},
		{Title: "Stars", Width: 8},
		{Title: "Forks", Width: 8},
		{Title: "Language", Width: 12},
		{Title: "Updated", Width: 12},
	}

	var rows []table.Row
	for _, repo := range m.publicRepos {
		lang := repo.Language
		if lang == "" {
			lang = "-"
		}
		rows = append(rows, table.Row{
			repo.Name,
			formatNumber(repo.Stars),
			formatNumber(repo.Forks),
			lang,
			repo.UpdatedAt.Format("2006-01-02"),
		})
	}

	m.table = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(m.height-8),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	m.table.SetStyles(s)
}

func (m *Model) updateTableSize() {
	if m.height > 0 {
		// Recreate table with new height
		columns := m.table.Columns()
		rows := m.table.Rows()
		m.table = table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(m.height-8),
		)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		m.table.SetStyles(s)
	}
}

func (m *Model) updateStatsView() {
	content := m.renderDetailedStats()
	m.viewport.SetContent(content)
}

func (m Model) View() string {
	if !m.ready && m.loading {
		return m.renderLoadingView()
	}

	var content string

	// Header
	header := m.renderHeader()

	// Notification bar
	var notifBar string
	if m.notification != "" {
		notifStyle := lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Padding(0, 1)

		if m.notifSuccess {
			notifBar = notifStyle.Background(lipgloss.Color("28")).
				Foreground(lipgloss.Color("15")).
				Render(m.notification)
		} else {
			notifBar = notifStyle.Background(lipgloss.Color("124")).
				Foreground(lipgloss.Color("15")).
				Render(m.notification)
		}
	}

	// Main content based on current view with proper centering
	switch m.currentView {
	case repoListView:
		content = m.renderRepoListView()
	case repoTableView:
		content = m.renderRepoTableView()
	case statsView:
		content = m.renderStatsView()
	case activityView:
		content = m.renderActivityView()
	}

	// Center the main content
	content = lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Padding(0, 2).
		Render(content)

	// Search bar
	var searchBar string
	if m.searchMode {
		searchBar = m.renderSearchBar()
	}

	// Help
	helpView := m.help.View(keys)

	// Combine all sections
	sections := []string{header}
	if notifBar != "" {
		sections = append(sections, notifBar)
	}
	if searchBar != "" {
		sections = append(sections, searchBar)
	}
	sections = append(sections, content, helpView)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderLoadingView() string {
	content := fmt.Sprintf("\n%s Loading GitHub data for %s...\n\n", m.spinner.View(), m.username)

	if m.reposLoaded {
		content += "Repositories loaded\n"
	} else {
		content += "Loading repositories...\n"
	}

	if m.eventsLoaded {
		content += "Activity loaded\n"
	} else {
		content += "Loading activity...\n"
	}

	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Height(m.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(2).
		Render(content)
}

func (m Model) renderHeader() string {
	title := fmt.Sprintf("GitHub Dashboard - %s", m.username)

	var stats string
	if len(m.publicRepos) > 0 {
		totalStars := 0
		totalForks := 0
		for _, repo := range m.publicRepos {
			totalStars += repo.Stars
			totalForks += repo.Forks
		}
		stats = fmt.Sprintf("® %d repos • ⋆ %s stars • ⑂ %s forks",
			len(m.publicRepos), formatNumber(totalStars), formatNumber(totalForks))
	}

	var viewIndicator string
	switch m.currentView {
	case repoListView:
		viewIndicator = "List View"
	case repoTableView:
		viewIndicator = "Table View"
	case statsView:
		viewIndicator = "Statistics"
	case activityView:
		viewIndicator = "Activity"
	}

	headerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("57")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 2).
		Width(m.width).
		Align(lipgloss.Center)

	headerContent := fmt.Sprintf("%s\n%s\n%s", title, stats, viewIndicator)
	return headerStyle.Render(headerContent)
}

func (m Model) renderRepoListView() string {
	return m.list.View()
}

func (m Model) renderRepoTableView() string {
	return m.table.View()
}

func (m Model) renderStatsView() string {
	return m.viewport.View()
}

func (m Model) renderActivityView() string {
	return m.list.View()
}

func (m Model) renderSearchBar() string {
	searchStyle := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Padding(0, 1).
		Background(lipgloss.Color("235"))

	searchContent := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("Search: ") + m.search.View()

	return searchStyle.Render(searchContent)
}

func (m Model) renderDetailedStats() string {
	var content strings.Builder

	content.WriteString(titleStyle.Render("Detailed Statistics"))
	content.WriteString("\n\n")

	// Repository Statistics
	if len(m.publicRepos) > 0 {
		totalStars := 0
		totalForks := 0
		languageCount := make(map[string]int)

		for _, repo := range m.publicRepos {
			totalStars += repo.Stars
			totalForks += repo.Forks
			if repo.Language != "" {
				languageCount[repo.Language]++
			}
		}

		content.WriteString("® Repository Overview:\n")
		content.WriteString(fmt.Sprintf("   Total Repositories: %d\n", len(m.publicRepos)))
		content.WriteString(fmt.Sprintf("   Total Stars: %s\n", formatNumber(totalStars)))
		content.WriteString(fmt.Sprintf("   Total Forks: %s\n", formatNumber(totalForks)))
		content.WriteString(fmt.Sprintf("   Average Stars: %.1f\n", float64(totalStars)/float64(len(m.publicRepos))))
		content.WriteString("\n")

		// Top repositories
		content.WriteString("Top Repositories by Stars:\n")
		for i, repo := range m.publicRepos {
			if i >= 5 {
				break
			}
			content.WriteString(fmt.Sprintf("   %d. %s - ⋆ %s\n", i+1, repo.Name, formatNumber(repo.Stars)))
		}
		content.WriteString("\n")

		// Languages
		if len(languageCount) > 0 {
			content.WriteString("Programming Languages:\n")
			for lang, count := range languageCount {
				content.WriteString(fmt.Sprintf("   %s: %d repositories\n", lang, count))
			}
			content.WriteString("\n")
		}
	}

	// Activity Statistics
	if len(m.events) > 0 {
		content.WriteString("Activity Statistics:\n")
		content.WriteString(fmt.Sprintf("   Push Events: %d\n", m.stats.PushEvents))
		content.WriteString(fmt.Sprintf("   Pull Request Events: %d\n", m.stats.PullRequestEvents))
		content.WriteString(fmt.Sprintf("   Issue Events: %d\n", m.stats.IssueEvents))
		content.WriteString(fmt.Sprintf("   Create Events: %d\n", m.stats.CreateEvents))
		content.WriteString(fmt.Sprintf("   Watch Events: %d\n", m.stats.WatchEvents))
		content.WriteString(fmt.Sprintf("   Total Events: %d\n", m.stats.TotalEvents))
		content.WriteString(fmt.Sprintf("   Activity Grade: %s\n", getGrade(m.stats)))
	}

	return content.String()
}

// Action commands
func (m Model) cloneRepo(repo PublicRepo) tea.Cmd {
	return func() tea.Msg {
		cloneCmd := fmt.Sprintf("git clone %s", repo.CloneURL)
		if err := copyToClipboard(cloneCmd); err != nil {
			return NotificationMsg{
				message:   fmt.Sprintf("❌ Copy Error: %v", err),
				isSuccess: false,
			}
		}
		return NotificationMsg{
			message:   fmt.Sprintf("Clone command copied: %s", repo.Name),
			isSuccess: true,
		}
	}
}

func (m Model) copyURL(repo PublicRepo) tea.Cmd {
	return func() tea.Msg {
		if err := copyToClipboard(repo.URL); err != nil {
			return NotificationMsg{
				message:   fmt.Sprintf("❌ Copy Error: %v", err),
				isSuccess: false,
			}
		}
		return NotificationMsg{
			message:   fmt.Sprintf("URL copied: %s", repo.Name),
			isSuccess: true,
		}
	}
}

func (m Model) openInBrowser(repo PublicRepo) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd

		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", repo.URL)
		case "linux":
			cmd = exec.Command("xdg-open", repo.URL)
		case "windows":
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", repo.URL)
		default:
			return NotificationMsg{
				message:   "❌ OS not supported for opening browser",
				isSuccess: false,
			}
		}

		if err := cmd.Run(); err != nil {
			return NotificationMsg{
				message:   fmt.Sprintf("❌ Error opening browser: %v", err),
				isSuccess: false,
			}
		}

		return NotificationMsg{
			message:   fmt.Sprintf("Opened in browser: %s", repo.Name),
			isSuccess: true,
		}
	}
}

// Initialize new model with bubbles components
func NewModel(username string) Model {
	// List component with better styling
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Background(lipgloss.Color("57")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Background(lipgloss.Color("57")).
		Foreground(lipgloss.Color("254")).
		Padding(0, 1)

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Title = "Loading repositories..."
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true).
		Padding(0, 2)

	// Table component
	t := table.New()

	// Viewport component
	v := viewport.New(0, 0)
	v.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2)

	// Help component
	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	h.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(lipgloss.Color("237"))

	// Spinner component
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Search input
	ti := textinput.New()
	ti.Placeholder = "Search repositories by name or description..."
	ti.CharLimit = 100
	ti.Width = 50

	return Model{
		username:     username,
		list:         l,
		table:        t,
		viewport:     v,
		help:         h,
		spinner:      s,
		search:       ti,
		currentView:  repoListView,
		loading:      true,
		reposLoaded:  false,
		eventsLoaded: false,
	}
}

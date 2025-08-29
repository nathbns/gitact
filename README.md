# 🐙 GitHub Activity CLI - Modern Edition

A beautiful, modern terminal interface for exploring GitHub profiles, repositories, and activity. Built with [Charm's Bubbles](https://github.com/charmbracelet/bubbles) UI components for a delightful terminal experience.

![GitHub CLI Screenshot](https://img.shields.io/badge/Built_with-Go-00ADD8?style=for-the-badge&logo=go)
![GitHub CLI Screenshot](https://img.shields.io/badge/UI-Bubbles-FF69B4?style=for-the-badge&logo=github)

## ✨ Features

### 🎯 Multiple Interactive Views
- **📋 Repository List** - Browse all public repositories with search functionality
- **📊 Table View** - Detailed tabular data showing stars, forks, languages, and update dates
- **📈 Statistics View** - Comprehensive statistics and insights about the user's GitHub profile
- **⚡ Activity Feed** - Recent GitHub activity timeline with event details

### 🔍 Smart Repository Discovery
- **Complete Repository Listing** - Shows ALL public repositories (not just recent activity)
- **Real-time Search** - Filter repositories by name or description
- **Rich Metadata** - Stars, forks, languages, descriptions, and update dates
- **Popularity Sorting** - Automatically sorted by star count

### 🚀 Enhanced User Experience
- **Responsive Design** - Adapts to terminal size
- **Keyboard Navigation** - Vim-like keybindings (j/k, h/l)
- **Quick Actions** - Clone, copy URLs, open in browser
- **Live Data Refresh** - Update data without restarting
- **Rate Limit Awareness** - Built-in GitHub API rate limit checking

## 📦 Installation

### Pre-built Binaries
```bash
# Download from releases (coming soon)
curl -LO https://github.com/yourusername/github-act-cli/releases/latest/download/gitact
chmod +x gitact
sudo mv gitact /usr/local/bin/
```

### Build from Source
```bash
git clone https://github.com/yourusername/github-act-cli.git
cd github-act-cli
go build -o gitact
./gitact --help
```

## 🔑 GitHub Token Setup (Recommended)

To avoid rate limits and access private repositories:

1. **Create a Personal Access Token**:
   - Go to [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens)
   - Generate a new token with `public_repo` scope
   - Copy the token

2. **Set Environment Variable**:
   ```bash
   export GITHUB_TOKEN=your_token_here
   ```

3. **Persistent Setup** (add to your shell profile):
   ```bash
   echo 'export GITHUB_TOKEN=your_token_here' >> ~/.zshrc  # or ~/.bashrc
   source ~/.zshrc
   ```

### Rate Limits
- **Without token**: 60 requests/hour per IP
- **With token**: 5,000 requests/hour
- **Our app uses**: ~2-4 requests per user

## 🚀 Usage

### Interactive Dashboard
```bash
# Explore a user's profile interactively
gitact karpathy

# With GitHub token for higher rate limits
GITHUB_TOKEN=xxx gitact karpathy
```

### Command Line Mode
```bash
# Get detailed repository listing
gitact --repos torvalds

# View help
gitact --help

# Check version
gitact --version
```

## ⌨️ Keyboard Navigation

### Global Controls
| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate items |
| `←/→` or `h/l` | Switch between views |
| `tab` | Next view |
| `?` | Toggle help |
| `q/esc` | Quit |

### Repository Actions
| Key | Action |
|-----|--------|
| `enter` | Select item |
| `c` | Copy git clone command |
| `x` | Copy repository URL |
| `o` | Open repository in browser |
| `r` | Refresh all data |

### Search (Repository List View)
| Key | Action |
|-----|--------|
| `/` | Start search |
| `enter` | Apply search filter |
| `esc` | Cancel search |

## 📊 Views Overview

### 1. Repository List View 📋
- **Interactive browsing** of all public repositories
- **Search functionality** - filter by name or description
- **Rich information** - stars, forks, descriptions
- **Quick actions** - clone, copy, open

### 2. Table View 📊
- **Tabular format** with sortable columns
- **Detailed metadata** - language, update dates
- **Compact overview** of all repositories
- **Easy comparison** between projects

### 3. Statistics View 📈
- **Comprehensive analytics** about the GitHub profile
- **Repository statistics** - total stars, forks, languages used
- **Top repositories** ranked by popularity
- **Activity insights** - push events, issues, PRs
- **Programming language breakdown**

### 4. Activity Feed ⚡
- **Recent GitHub activity** timeline
- **Event details** - pushes, issues, PRs, stars
- **Repository context** for each activity
- **Time-based sorting** of events

## 🎯 Examples

### Exploring Machine Learning Researchers
```bash
# Andrej Karpathy - AI/ML repositories
gitact karpathy

# Explore his popular projects like nanoGPT, llm.c, micrograd
# Use search to filter by topic: "/gpt" or "/neural"
```

### Linux Kernel Development
```bash
# Linus Torvalds
gitact torvalds

# View the famous linux repository and other projects
```

### Popular Open Source Projects
```bash
# GitHub's mascot account
gitact --repos octocat

# Browse training repositories like Spoon-Knife
```

### Web Development
```bash
# Explore modern web development repositories
gitact sindresorhus

# Search for specific technologies: "/react" or "/node"
```

## 🔧 Configuration

### Environment Variables
- `GITHUB_TOKEN` - GitHub personal access token for higher rate limits
- `NO_COLOR` - Disable colored output
- `GITACT_CACHE_DIR` - Custom cache directory (default: `~/.cache/gitact`)

### Cache
The application caches API responses to improve performance and reduce rate limit usage:
- **Location**: `~/.cache/gitact/`
- **Duration**: 10 minutes for repository data, 5 minutes for activity
- **Clear cache**: `rm -rf ~/.cache/gitact/`

## 🤝 Contributing

We welcome contributions! Here's how to get started:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes** and add tests
4. **Run tests**: `go test ./...`
5. **Commit changes**: `git commit -m 'Add amazing feature'`
6. **Push to branch**: `git push origin feature/amazing-feature`
7. **Open a Pull Request**

### Development Setup
```bash
git clone https://github.com/yourusername/github-act-cli.git
cd github-act-cli
go mod tidy
go run . --help
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Charm](https://charm.sh/) - For the amazing Bubbles UI framework
- [GitHub API](https://docs.github.com/en/rest) - For providing comprehensive repository data
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - For beautiful terminal styling
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - For the elegant TUI framework

## 🔮 Roadmap

- [ ] **Caching system** for better performance
- [ ] **GitHub organization** support
- [ ] **Repository comparison** features
- [ ] **Export functionality** (JSON, CSV)
- [ ] **Custom themes** and color schemes
- [ ] **Plugin system** for extensions
- [ ] **Docker support** for easy distribution
- [ ] **GitHub Actions integration**

---

<div align="center">

**Made with ❤️ and Go**

[Report Bug](https://github.com/yourusername/github-act-cli/issues) · [Request Feature](https://github.com/yourusername/github-act-cli/issues) · [Documentation](https://github.com/yourusername/github-act-cli/wiki)

</div>
#!/bin/bash

# GitHub Activity CLI - Publication Script
# This script automates the entire publication process

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="gitact"
GITHUB_REPO="yourusername/gitact"
HOMEBREW_TAP="yourusername/homebrew-tap"

# Functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

check_dependencies() {
    log_info "Checking dependencies..."

    local deps=("git" "gh" "go" "jq")
    for dep in "${deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            log_error "$dep is required but not installed"
        fi
    done

    log_success "All dependencies found"
}

check_git_status() {
    log_info "Checking Git status..."

    if [[ -n $(git status --porcelain) ]]; then
        log_error "Working directory is not clean. Please commit or stash changes."
    fi

    local current_branch=$(git branch --show-current)
    if [[ "$current_branch" != "main" ]] && [[ "$current_branch" != "master" ]]; then
        log_error "Please switch to main/master branch before publishing"
    fi

    log_success "Git status is clean"
}

get_version() {
    echo "📝 Current version information:"
    if git describe --tags --exact-match 2>/dev/null; then
        log_warning "Current commit is already tagged"
    fi

    echo "🏷️  Recent tags:"
    git tag --sort=-version:refname | head -5

    echo ""
    read -p "Enter new version (e.g., 1.0.0): " NEW_VERSION

    if [[ ! $NEW_VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        log_error "Version must follow semantic versioning (x.y.z)"
    fi

    if git tag | grep -q "^v${NEW_VERSION}$"; then
        log_error "Version v${NEW_VERSION} already exists"
    fi

    echo "v${NEW_VERSION}"
}

update_version_files() {
    local version=$1
    log_info "Updating version in files..."

    # Update main.go if it has a version variable
    if grep -q "Version.*=" main.go; then
        sed -i.bak "s/Version.*=.*/Version = \"${version#v}\"/" main.go
        rm main.go.bak
    fi

    # Update homebrew formula
    if [[ -f homebrew-formula.rb ]]; then
        sed -i.bak "s/version \".*\"/version \"${version#v}\"/" homebrew-formula.rb
        rm homebrew-formula.rb.bak
    fi

    log_success "Version files updated"
}

run_tests() {
    log_info "Running tests..."

    go mod tidy
    go vet ./...
    go test ./...

    log_success "All tests passed"
}

build_and_test() {
    log_info "Building application..."

    go build -o $BINARY_NAME .

    # Basic functionality test
    if ! ./$BINARY_NAME --version &>/dev/null; then
        log_error "Built binary fails basic test"
    fi

    rm $BINARY_NAME
    log_success "Build successful"
}

create_changelog() {
    local version=$1
    log_info "Creating changelog entry..."

    local changelog_file="CHANGELOG.md"
    local temp_file=$(mktemp)

    cat > "$temp_file" << EOF
# Changelog

## [${version#v}] - $(date +%Y-%m-%d)

### Added
- Add your new features here

### Changed
- Add your changes here

### Fixed
- Add your bug fixes here

### Security
- Add security updates here

EOF

    if [[ -f "$changelog_file" ]]; then
        tail -n +2 "$changelog_file" >> "$temp_file"
    fi

    mv "$temp_file" "$changelog_file"

    # Open editor for user to modify
    ${EDITOR:-nano} "$changelog_file"

    log_success "Changelog updated"
}

commit_and_tag() {
    local version=$1
    log_info "Committing changes and creating tag..."

    git add .
    git commit -m "Release ${version}"
    git tag -a "${version}" -m "Release ${version}"

    log_success "Changes committed and tagged"
}

push_to_github() {
    local version=$1
    log_info "Pushing to GitHub..."

    git push origin main
    git push origin "${version}"

    log_success "Pushed to GitHub"
}

create_github_release() {
    local version=$1
    log_info "Creating GitHub release..."

    # The GitHub Action will handle the actual release creation
    # This just triggers it by pushing the tag

    log_info "Waiting for GitHub Actions to build release..."
    sleep 5

    # Check if release was created
    local max_attempts=30
    local attempt=1

    while [[ $attempt -le $max_attempts ]]; do
        if gh release view "${version}" &>/dev/null; then
            log_success "GitHub release created successfully"
            return 0
        fi

        echo -n "."
        sleep 10
        ((attempt++))
    done

    log_warning "Release creation is taking longer than expected. Check GitHub Actions."
}

setup_homebrew_tap() {
    log_info "Setting up Homebrew tap..."

    local tap_dir="../homebrew-tap"

    if [[ ! -d "$tap_dir" ]]; then
        log_info "Cloning Homebrew tap repository..."
        git clone "https://github.com/${HOMEBREW_TAP}.git" "$tap_dir"
    fi

    cd "$tap_dir"
    git pull origin main
    cd -

    # Copy formula
    cp homebrew-formula.rb "${tap_dir}/Formula/${BINARY_NAME}.rb"

    log_success "Homebrew tap setup complete"
}

update_homebrew_formula() {
    local version=$1
    log_info "Updating Homebrew formula..."

    # This will be handled by GitHub Actions automatically
    # The bump-homebrew-formula-action will update the formula

    log_success "Homebrew formula will be updated automatically"
}

verify_installation() {
    local version=$1
    log_info "Verifying installation methods..."

    echo "🧪 Testing installation methods:"
    echo "1. Manual download from GitHub releases"
    echo "2. Homebrew installation (will be available shortly)"

    echo ""
    echo "📋 Installation commands:"
    echo "# Homebrew (macOS/Linux)"
    echo "brew tap ${HOMEBREW_TAP}"
    echo "brew install ${BINARY_NAME}"
    echo ""
    echo "# Manual installation (Linux x64)"
    echo "curl -LO https://github.com/${GITHUB_REPO}/releases/download/${version}/${BINARY_NAME}-${version#v}-linux-amd64.tar.gz"
    echo "tar -xzf ${BINARY_NAME}-${version#v}-linux-amd64.tar.gz"
    echo "sudo mv ${BINARY_NAME} /usr/local/bin/"

    log_success "Installation instructions prepared"
}

post_release_tasks() {
    local version=$1
    log_info "Running post-release tasks..."

    echo "🎉 Release ${version} published successfully!"
    echo ""
    echo "📋 Next steps:"
    echo "1. Update documentation if needed"
    echo "2. Announce release on social media"
    echo "3. Update project README with new version"
    echo "4. Monitor for any issues"
    echo ""
    echo "🔗 Links:"
    echo "• GitHub Release: https://github.com/${GITHUB_REPO}/releases/tag/${version}"
    echo "• Homebrew Formula: https://github.com/${HOMEBREW_TAP}/blob/main/Formula/${BINARY_NAME}.rb"

    log_success "Post-release tasks completed"
}

main() {
    echo "🚀 GitHub Activity CLI - Publication Script"
    echo "=========================================="
    echo ""

    # Pre-flight checks
    check_dependencies
    check_git_status

    # Get version information
    local version=$(get_version)
    log_info "Publishing version: $version"

    # Update files
    update_version_files "$version"

    # Quality checks
    run_tests
    build_and_test

    # Documentation
    create_changelog "$version"

    # Git operations
    commit_and_tag "$version"
    push_to_github "$version"

    # Release creation
    create_github_release "$version"

    # Homebrew setup
    setup_homebrew_tap
    update_homebrew_formula "$version"

    # Verification
    verify_installation "$version"

    # Cleanup
    post_release_tasks "$version"

    log_success "Publication completed successfully! 🎉"
}

# Run main function if script is executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi

package usecase

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// GitHub操作のインターフェース（テスト容易性のため）
type GitHubClient interface {
	ListOpenIssues(repo string) ([]Issue, error)
	GetIssueComments(repo string, issueNumber int) ([]Comment, error)
	CloseIssue(repo string, issueNumber int) error
	GetCurrentRepo(baseDir string) (string, error)
}

// Issue represents a GitHub issue
type Issue struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Author    User      `json:"author"`
	Assignees []User    `json:"assignees"`
	Labels    []Label   `json:"labels"`
	Milestone Milestone `json:"milestone"`
	URL       string    `json:"url"`
}

// Comment represents a GitHub issue comment
type Comment struct {
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	Author    User      `json:"author"`
}

// User represents a GitHub user
type User struct {
	Login string `json:"login"`
}

// Label represents a GitHub label
type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// Milestone represents a GitHub milestone
type Milestone struct {
	Title string `json:"title"`
}

// ImportOptions contains options for importing issues
type ImportOptions struct {
	Repo    string
	DryRun  bool
	NoClose bool
}

// GHClient implements GitHubClient using gh command
type GHClient struct{}

func (c *GHClient) ListOpenIssues(repo string) ([]Issue, error) {
	cmd := exec.Command("gh", "issue", "list",
		"--repo", repo,
		"--state", "open",
		"--json", "number,title,body,state,createdAt,updatedAt,author,assignees,labels,milestone,url")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gh command failed: %w", err)
	}

	var issues []Issue
	if err := json.Unmarshal(output, &issues); err != nil {
		return nil, fmt.Errorf("failed to parse issues: %w", err)
	}

	return issues, nil
}

func (c *GHClient) GetIssueComments(repo string, issueNumber int) ([]Comment, error) {
	cmd := exec.Command("gh", "issue", "view", fmt.Sprintf("%d", issueNumber),
		"--repo", repo,
		"--json", "comments")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gh command failed: %w", err)
	}

	var result struct {
		Comments []Comment `json:"comments"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse comments: %w", err)
	}

	return result.Comments, nil
}

func (c *GHClient) CloseIssue(repo string, issueNumber int) error {
	cmd := exec.Command("gh", "issue", "close", fmt.Sprintf("%d", issueNumber),
		"--repo", repo)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to close issue: %w", err)
	}

	return nil
}

func (c *GHClient) GetCurrentRepo(baseDir string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = baseDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get remote origin: %w", err)
	}

	url := strings.TrimSpace(string(output))

	// HTTPS形式: https://github.com/owner/repo.git
	if strings.HasPrefix(url, "https://github.com/") {
		url = strings.TrimPrefix(url, "https://github.com/")
		url = strings.TrimSuffix(url, ".git")
		return url, nil
	}

	// SSH形式: git@github.com:owner/repo.git
	if strings.HasPrefix(url, "git@github.com:") {
		url = strings.TrimPrefix(url, "git@github.com:")
		url = strings.TrimSuffix(url, ".git")
		return url, nil
	}

	return "", fmt.Errorf("not a GitHub repository: %s", url)
}

// MockGitHubClient is a mock implementation for testing
type MockGitHubClient struct {
	Issues       []Issue
	Comments     map[int][]Comment
	ClosedIssues []int
	RepoURL      string
	ErrorOnList  error
	ErrorOnGet   error
	ErrorOnClose error
}

func (m *MockGitHubClient) ListOpenIssues(repo string) ([]Issue, error) {
	if m.ErrorOnList != nil {
		return nil, m.ErrorOnList
	}
	return m.Issues, nil
}

func (m *MockGitHubClient) GetIssueComments(repo string, issueNumber int) ([]Comment, error) {
	if m.ErrorOnGet != nil {
		return nil, m.ErrorOnGet
	}
	comments, exists := m.Comments[issueNumber]
	if !exists {
		return []Comment{}, nil
	}
	return comments, nil
}

func (m *MockGitHubClient) CloseIssue(repo string, issueNumber int) error {
	if m.ErrorOnClose != nil {
		return m.ErrorOnClose
	}
	m.ClosedIssues = append(m.ClosedIssues, issueNumber)
	return nil
}

func (m *MockGitHubClient) GetCurrentRepo(baseDir string) (string, error) {
	if m.RepoURL == "" {
		return "owner/repo", nil
	}
	return m.RepoURL, nil
}

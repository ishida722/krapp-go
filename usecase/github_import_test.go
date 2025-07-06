package usecase

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ishida722/krapp-go/models"
)

// testConfig implements InboxConfig for testing
type testConfig struct {
	baseDir  string
	inboxDir string
}

func (c *testConfig) GetBaseDir() string {
	return c.baseDir
}

func (c *testConfig) GetInboxDir() string {
	return "inbox"
}

func (c *testConfig) GetInboxTemplate() map[string]any {
	return map[string]any{
		"tags":   []string{},
		"status": "new",
	}
}

func TestImportGitHubIssues(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "krapp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inboxDir := filepath.Join(tempDir, "inbox")
	if err := os.MkdirAll(inboxDir, 0755); err != nil {
		t.Fatalf("Failed to create inbox dir: %v", err)
	}

	// テストデータ
	testIssue := Issue{
		Number:    123,
		Title:     "Test Issue",
		Body:      "This is a test issue body",
		State:     "open",
		CreatedAt: time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		Author:    User{Login: "testuser"},
		Assignees: []User{{Login: "assignee1"}, {Login: "assignee2"}},
		Labels:    []Label{{Name: "bug"}, {Name: "priority-high"}},
		Milestone: Milestone{Title: "v1.2.0"},
		URL:       "https://github.com/owner/repo/issues/123",
	}

	testComment := Comment{
		Body:      "This is a test comment",
		CreatedAt: time.Date(2024, 1, 12, 12, 0, 0, 0, time.UTC),
		Author:    User{Login: "reviewer"},
	}

	tests := []struct {
		name         string
		issues       []Issue
		comments     map[int][]Comment
		options      ImportOptions
		expectFiles  int
		expectClosed int
	}{
		{
			name:   "single issue with comments",
			issues: []Issue{testIssue},
			comments: map[int][]Comment{
				123: {testComment},
			},
			options:      ImportOptions{},
			expectFiles:  1,
			expectClosed: 1,
		},
		{
			name:         "dry run mode",
			issues:       []Issue{testIssue},
			comments:     map[int][]Comment{123: {testComment}},
			options:      ImportOptions{DryRun: true},
			expectFiles:  1,
			expectClosed: 0,
		},
		{
			name:         "no close mode",
			issues:       []Issue{testIssue},
			comments:     map[int][]Comment{123: {testComment}},
			options:      ImportOptions{NoClose: true},
			expectFiles:  1,
			expectClosed: 0,
		},
		{
			name:         "no issues",
			issues:       []Issue{},
			comments:     map[int][]Comment{},
			options:      ImportOptions{},
			expectFiles:  0,
			expectClosed: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用のディレクトリをクリア
			os.RemoveAll(inboxDir)
			os.MkdirAll(inboxDir, 0755)

			mockClient := &MockGitHubClient{
				Issues:       tt.issues,
				Comments:     tt.comments,
				ClosedIssues: []int{},
				RepoURL:      "owner/repo",
			}

			cfg := &testConfig{
				baseDir:  tempDir,
				inboxDir: inboxDir,
			}

			err := ImportGitHubIssues(cfg, mockClient, tt.options)
			if err != nil {
				t.Errorf("ImportGitHubIssues() error = %v", err)
				return
			}

			// ファイル作成数の確認
			files, err := os.ReadDir(inboxDir)
			if err != nil {
				t.Errorf("Failed to read inbox dir: %v", err)
				return
			}
			if len(files) != tt.expectFiles {
				t.Errorf("Expected %d files, got %d", tt.expectFiles, len(files))
			}

			// クローズされたissue数の確認
			if len(mockClient.ClosedIssues) != tt.expectClosed {
				t.Errorf("Expected %d closed issues, got %d", tt.expectClosed, len(mockClient.ClosedIssues))
			}

			// ファイル内容の確認（ファイルがある場合）
			if len(files) > 0 {
				filePath := filepath.Join(inboxDir, files[0].Name())
				note, err := models.LoadNoteFromFile(filePath)
				if err != nil {
					t.Errorf("Failed to load note: %v", err)
					return
				}

				// frontmatterの確認
				if note.FrontMatter["issue_number"] != 123 {
					t.Errorf("Expected issue_number 123, got %v", note.FrontMatter["issue_number"])
				}
				if note.FrontMatter["status"] != "imported" {
					t.Errorf("Expected status 'imported', got %v", note.FrontMatter["status"])
				}

				// タグの確認
				tags, ok := note.FrontMatter["tags"].([]interface{})
				if !ok {
					t.Errorf("Tags should be a slice")
				} else {
					found := false
					for _, tag := range tags {
						if tag == "github-issue" {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected 'github-issue' tag not found")
					}
				}

				// コンテンツの確認
				if !strings.Contains(note.Content, "Test Issue") {
					t.Errorf("Content should contain issue title")
				}
				if !strings.Contains(note.Content, "This is a test issue body") {
					t.Errorf("Content should contain issue body")
				}
			}
		})
	}
}

func TestGenerateIssueFilename(t *testing.T) {
	tests := []struct {
		name     string
		issue    Issue
		expected string
	}{
		{
			name: "normal title",
			issue: Issue{
				Number:    123,
				Title:     "Fix config loading bug",
				CreatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			},
			expected: "2024-01-15-issue-123-fix-config-loading-bug.md",
		},
		{
			name: "title with special characters",
			issue: Issue{
				Number:    456,
				Title:     "Add new feature: user authentication!",
				CreatedAt: time.Date(2024, 2, 20, 15, 45, 0, 0, time.UTC),
			},
			expected: "2024-02-20-issue-456-add-new-feature-user-authentication.md",
		},
		{
			name: "very long title",
			issue: Issue{
				Number:    789,
				Title:     "This is a very long title that should be truncated because it exceeds the maximum length",
				CreatedAt: time.Date(2024, 3, 25, 8, 15, 0, 0, time.UTC),
			},
			expected: "2024-03-25-issue-789-this-is-a-very-long-title-that-should-be-truncated.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateIssueFilename(tt.issue)
			if result != tt.expected {
				t.Errorf("generateIssueFilename() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Simple Title", "simple-title"},
		{"Title with spaces", "title-with-spaces"},
		{"Title!@#$%^&*()with special chars", "title-with-special-chars"},
		{"Multiple---dashes", "multiple-dashes"},
		{"---leading-and-trailing---", "leading-and-trailing"},
		{"CamelCaseTitle", "camelcasetitle"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeFilename(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCreateIssueFrontMatter(t *testing.T) {
	issue := Issue{
		Number:    123,
		Title:     "Test Issue",
		State:     "open",
		CreatedAt: time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		Author:    User{Login: "testuser"},
		Assignees: []User{{Login: "assignee1"}, {Login: "assignee2"}},
		Labels:    []Label{{Name: "bug"}, {Name: "priority-high"}},
		Milestone: Milestone{Title: "v1.2.0"},
		URL:       "https://github.com/owner/repo/issues/123",
	}

	fm := createIssueFrontMatter(issue)

	// 基本フィールドの確認
	if fm["issue_number"] != 123 {
		t.Errorf("Expected issue_number 123, got %v", fm["issue_number"])
	}
	if fm["status"] != "imported" {
		t.Errorf("Expected status 'imported', got %v", fm["status"])
	}
	if fm["issue_url"] != "https://github.com/owner/repo/issues/123" {
		t.Errorf("Expected correct issue URL, got %v", fm["issue_url"])
	}

	// createdの確認（SetCreatedメソッドにより文字列形式で保存される）
	created, ok := fm["created"].(string)
	if !ok {
		t.Errorf("Expected created to be string, got %T", fm["created"])
	} else {
		expectedCreated := issue.CreatedAt.Format("2006-01-02")
		if created != expectedCreated {
			t.Errorf("Expected created %v, got %v", expectedCreated, created)
		}
	}

	// 配列フィールドの確認
	assignees, ok := fm["assignees"].([]string)
	if !ok {
		t.Errorf("Expected assignees to be []string, got %T", fm["assignees"])
	} else {
		expected := []string{"assignee1", "assignee2"}
		if len(assignees) != len(expected) {
			t.Errorf("Expected assignees %v, got %v", expected, assignees)
		}
	}

	labels, ok := fm["labels"].([]string)
	if !ok {
		t.Errorf("Expected labels to be []string, got %T", fm["labels"])
	} else {
		expected := []string{"bug", "priority-high"}
		if len(labels) != len(expected) {
			t.Errorf("Expected labels %v, got %v", expected, labels)
		}
	}

	if fm["milestone"] != "v1.2.0" {
		t.Errorf("Expected milestone 'v1.2.0', got %v", fm["milestone"])
	}
}

func TestGenerateIssueMarkdown(t *testing.T) {
	issue := Issue{
		Number:    123,
		Title:     "Test Issue",
		Body:      "This is a test issue body",
		CreatedAt: time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
		Author:    User{Login: "testuser"},
		Labels:    []Label{{Name: "bug"}},
		Assignees: []User{{Login: "assignee1"}},
		Milestone: Milestone{Title: "v1.2.0"},
	}

	comments := []Comment{
		{
			Body:      "This is a comment",
			CreatedAt: time.Date(2024, 1, 12, 12, 0, 0, 0, time.UTC),
			Author:    User{Login: "reviewer"},
		},
	}

	markdown := generateIssueMarkdown(issue, comments)

	// 基本的な内容の確認
	if !strings.Contains(markdown, "# Issue #123: Test Issue") {
		t.Errorf("Markdown should contain issue title")
	}
	if !strings.Contains(markdown, "**Created by:** @testuser") {
		t.Errorf("Markdown should contain author")
	}
	if !strings.Contains(markdown, "**Labels:** bug") {
		t.Errorf("Markdown should contain labels")
	}
	if !strings.Contains(markdown, "**Assignees:** @assignee1") {
		t.Errorf("Markdown should contain assignees")
	}
	if !strings.Contains(markdown, "**Milestone:** v1.2.0") {
		t.Errorf("Markdown should contain milestone")
	}
	if !strings.Contains(markdown, "This is a test issue body") {
		t.Errorf("Markdown should contain issue body")
	}
	if !strings.Contains(markdown, "## Comments") {
		t.Errorf("Markdown should contain comments section")
	}
	if !strings.Contains(markdown, "This is a comment") {
		t.Errorf("Markdown should contain comment body")
	}
	if !strings.Contains(markdown, "Issue automatically imported") {
		t.Errorf("Markdown should contain footer")
	}
}

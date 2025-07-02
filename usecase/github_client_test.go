package usecase

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGHClientGetCurrentRepo(t *testing.T) {
	tests := []struct {
		name     string
		setupGit func(string) error
		expected string
		wantErr  bool
	}{
		{
			name: "HTTPS GitHub URL",
			setupGit: func(dir string) error {
				return setupGitRepo(dir, "https://github.com/owner/repo.git")
			},
			expected: "owner/repo",
			wantErr:  false,
		},
		{
			name: "SSH GitHub URL",
			setupGit: func(dir string) error {
				return setupGitRepo(dir, "git@github.com:owner/repo.git")
			},
			expected: "owner/repo",
			wantErr:  false,
		},
		{
			name: "HTTPS GitHub URL without .git",
			setupGit: func(dir string) error {
				return setupGitRepo(dir, "https://github.com/owner/repo")
			},
			expected: "owner/repo",
			wantErr:  false,
		},
		{
			name: "non-GitHub URL",
			setupGit: func(dir string) error {
				return setupGitRepo(dir, "https://gitlab.com/owner/repo.git")
			},
			expected: "",
			wantErr:  true,
		},
		{
			name: "no git repository",
			setupGit: func(dir string) error {
				// gitリポジトリを作成しない
				return nil
			},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用の一時ディレクトリを作成
			tempDir, err := os.MkdirTemp("", "krapp-git-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// テスト用のgitリポジトリをセットアップ
			if err := tt.setupGit(tempDir); err != nil {
				t.Fatalf("Failed to setup git repo: %v", err)
			}

			client := &GHClient{}
			result, err := client.GetCurrentRepo(tempDir)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetCurrentRepo() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("GetCurrentRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result != tt.expected {
				t.Errorf("GetCurrentRepo() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// setupGitRepo creates a mock git repository with the specified remote origin URL
func setupGitRepo(dir, originURL string) error {
	gitDir := filepath.Join(dir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(gitDir, "config")
	configContent := `[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
[remote "origin"]
	url = ` + originURL + `
	fetch = +refs/heads/*:refs/remotes/origin/*
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return err
	}

	// HEADファイルも作成（gitが正常に動作するため）
	headPath := filepath.Join(gitDir, "HEAD")
	headContent := "ref: refs/heads/main\n"
	return os.WriteFile(headPath, []byte(headContent), 0644)
}

func TestMockGitHubClient(t *testing.T) {
	// モッククライアントの基本的な動作テスト
	issues := []Issue{
		{Number: 1, Title: "Test Issue 1"},
		{Number: 2, Title: "Test Issue 2"},
	}

	comments := map[int][]Comment{
		1: {{Body: "Comment for issue 1", Author: User{Login: "user1"}}},
		2: {},
	}

	mockClient := &MockGitHubClient{
		Issues:   issues,
		Comments: comments,
		RepoURL:  "test/repo",
	}

	// ListOpenIssues のテスト
	resultIssues, err := mockClient.ListOpenIssues("test/repo")
	if err != nil {
		t.Errorf("ListOpenIssues() error = %v", err)
	}
	if len(resultIssues) != 2 {
		t.Errorf("Expected 2 issues, got %d", len(resultIssues))
	}

	// GetIssueComments のテスト
	resultComments, err := mockClient.GetIssueComments("test/repo", 1)
	if err != nil {
		t.Errorf("GetIssueComments() error = %v", err)
	}
	if len(resultComments) != 1 {
		t.Errorf("Expected 1 comment for issue 1, got %d", len(resultComments))
	}

	resultComments, err = mockClient.GetIssueComments("test/repo", 2)
	if err != nil {
		t.Errorf("GetIssueComments() error = %v", err)
	}
	if len(resultComments) != 0 {
		t.Errorf("Expected 0 comments for issue 2, got %d", len(resultComments))
	}

	// CloseIssue のテスト
	err = mockClient.CloseIssue("test/repo", 1)
	if err != nil {
		t.Errorf("CloseIssue() error = %v", err)
	}
	if len(mockClient.ClosedIssues) != 1 || mockClient.ClosedIssues[0] != 1 {
		t.Errorf("Expected issue 1 to be closed, got %v", mockClient.ClosedIssues)
	}

	// GetCurrentRepo のテスト
	repo, err := mockClient.GetCurrentRepo("/tmp")
	if err != nil {
		t.Errorf("GetCurrentRepo() error = %v", err)
	}
	if repo != "test/repo" {
		t.Errorf("Expected repo 'test/repo', got %v", repo)
	}
}

func TestMockGitHubClientErrors(t *testing.T) {
	// エラーを返すモッククライアントのテスト
	mockClient := &MockGitHubClient{
		ErrorOnList:  ErrTestList,
		ErrorOnGet:   ErrTestGet,
		ErrorOnClose: ErrTestClose,
	}

	// エラーケースのテスト
	_, err := mockClient.ListOpenIssues("test/repo")
	if err != ErrTestList {
		t.Errorf("Expected ErrTestList, got %v", err)
	}

	_, err = mockClient.GetIssueComments("test/repo", 1)
	if err != ErrTestGet {
		t.Errorf("Expected ErrTestGet, got %v", err)
	}

	err = mockClient.CloseIssue("test/repo", 1)
	if err != ErrTestClose {
		t.Errorf("Expected ErrTestClose, got %v", err)
	}
}

// テスト用のエラー定義
var (
	ErrTestList  = &TestError{"test list error"}
	ErrTestGet   = &TestError{"test get error"}
	ErrTestClose = &TestError{"test close error"}
)

type TestError struct {
	msg string
}

func (e *TestError) Error() string {
	return e.msg
}
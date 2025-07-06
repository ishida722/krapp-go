package krapp

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ishida722/krapp-go/config"
	"github.com/ishida722/krapp-go/usecase"
)

func TestImportIssuesCommand(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "krapp-integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用の設定を作成
	testConfig := config.Config{
		BaseDir:      tempDir,
		Inbox:        "inbox",
		DailyNoteDir: "daily",
		Editor:       "vim",
	}

	// inboxディレクトリを作成
	inboxDir := filepath.Join(tempDir, "inbox")
	if err := os.MkdirAll(inboxDir, 0755); err != nil {
		t.Fatalf("Failed to create inbox dir: %v", err)
	}

	// テスト用のissueデータ
	testIssue := usecase.Issue{
		Number:    123,
		Title:     "Integration Test Issue",
		Body:      "This is an integration test issue",
		State:     "open",
		CreatedAt: time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		Author:    usecase.User{Login: "testuser"},
		URL:       "https://github.com/test/repo/issues/123",
	}

	testComment := usecase.Comment{
		Body:      "This is a test comment",
		CreatedAt: time.Date(2024, 1, 12, 12, 0, 0, 0, time.UTC),
		Author:    usecase.User{Login: "reviewer"},
	}

	// モッククライアント
	mockClient := &usecase.MockGitHubClient{
		Issues: []usecase.Issue{testIssue},
		Comments: map[int][]usecase.Comment{
			123: {testComment},
		},
		ClosedIssues: []int{},
		RepoURL:      "test/repo",
	}

	// configAdapterを作成
	adapter := &configAdapter{&testConfig}

	// ImportGitHubIssues関数を直接テスト
	options := usecase.ImportOptions{
		Repo:    "test/repo",
		DryRun:  false,
		NoClose: false,
	}

	err = usecase.ImportGitHubIssues(adapter, mockClient, options)
	if err != nil {
		t.Errorf("ImportGitHubIssues() error = %v", err)
	}

	// 結果の確認
	files, err := os.ReadDir(inboxDir)
	if err != nil {
		t.Errorf("Failed to read inbox dir: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}

	if len(mockClient.ClosedIssues) != 1 {
		t.Errorf("Expected 1 closed issue, got %d", len(mockClient.ClosedIssues))
	}

	if len(mockClient.ClosedIssues) > 0 && mockClient.ClosedIssues[0] != 123 {
		t.Errorf("Expected issue 123 to be closed, got %d", mockClient.ClosedIssues[0])
	}
}

func TestImportIssuesCommandWithDryRun(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "krapp-integration-test-dry-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用の設定を作成
	testConfig := config.Config{
		BaseDir:      tempDir,
		Inbox:        "inbox",
		DailyNoteDir: "daily",
		Editor:       "vim",
	}

	// inboxディレクトリを作成
	inboxDir := filepath.Join(tempDir, "inbox")
	if err := os.MkdirAll(inboxDir, 0755); err != nil {
		t.Fatalf("Failed to create inbox dir: %v", err)
	}

	// テスト用のissueデータ
	testIssue := usecase.Issue{
		Number:    456,
		Title:     "Dry Run Test Issue",
		Body:      "This is a dry run test",
		State:     "open",
		CreatedAt: time.Date(2024, 2, 15, 14, 30, 0, 0, time.UTC),
		Author:    usecase.User{Login: "dryrunuser"},
		URL:       "https://github.com/test/repo/issues/456",
	}

	// モッククライアント
	mockClient := &usecase.MockGitHubClient{
		Issues:       []usecase.Issue{testIssue},
		Comments:     map[int][]usecase.Comment{},
		ClosedIssues: []int{},
		RepoURL:      "test/repo",
	}

	// configAdapterを作成
	adapter := &configAdapter{&testConfig}

	// DryRunオプションでテスト
	options := usecase.ImportOptions{
		Repo:    "test/repo",
		DryRun:  true,
		NoClose: false,
	}

	err = usecase.ImportGitHubIssues(adapter, mockClient, options)
	if err != nil {
		t.Errorf("ImportGitHubIssues() error = %v", err)
	}

	// 結果の確認：ファイルは作成されるがissueはクローズされない
	files, err := os.ReadDir(inboxDir)
	if err != nil {
		t.Errorf("Failed to read inbox dir: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}

	if len(mockClient.ClosedIssues) != 0 {
		t.Errorf("Expected 0 closed issues in dry run, got %d", len(mockClient.ClosedIssues))
	}
}

func TestConfigAdapterImplementsInboxConfig(t *testing.T) {
	// configAdapterがInboxConfigインターフェースを実装していることを確認
	testConfig := config.Config{
		BaseDir: "/test/base",
		Inbox:   "inbox",
	}

	adapter := &configAdapter{&testConfig}

	// InboxConfigインターフェースとして使用できることを確認
	var _ usecase.InboxConfig = adapter

	// メソッドの動作確認
	if adapter.GetBaseDir() != "/test/base" {
		t.Errorf("Expected base dir '/test/base', got %s", adapter.GetBaseDir())
	}

	expectedInboxDir := "inbox"
	if adapter.GetInboxDir() != expectedInboxDir {
		t.Errorf("Expected inbox dir %s, got %s", expectedInboxDir, adapter.GetInboxDir())
	}
}

func TestImportIssuesNoIssues(t *testing.T) {
	// issueが存在しない場合のテスト
	tempDir, err := os.MkdirTemp("", "krapp-no-issues-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testConfig := config.Config{
		BaseDir: tempDir,
		Inbox:   "inbox",
	}

	inboxDir := filepath.Join(tempDir, "inbox")
	if err := os.MkdirAll(inboxDir, 0755); err != nil {
		t.Fatalf("Failed to create inbox dir: %v", err)
	}

	// 空のissueリストを持つモッククライアント
	mockClient := &usecase.MockGitHubClient{
		Issues:       []usecase.Issue{},
		Comments:     map[int][]usecase.Comment{},
		ClosedIssues: []int{},
		RepoURL:      "test/repo",
	}

	adapter := &configAdapter{&testConfig}

	options := usecase.ImportOptions{
		Repo: "test/repo",
	}

	err = usecase.ImportGitHubIssues(adapter, mockClient, options)
	if err != nil {
		t.Errorf("ImportGitHubIssues() should not error with no issues, got: %v", err)
	}

	// ファイルが作成されていないことを確認
	files, err := os.ReadDir(inboxDir)
	if err != nil {
		t.Errorf("Failed to read inbox dir: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files with no issues, got %d", len(files))
	}
}
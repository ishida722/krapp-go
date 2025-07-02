# GitHub Issue インポート機能仕様書

## 概要

krappにGitHub Issueを取得してinboxノートとして保存し、自動クローズする機能を追加する。GitHub CLIツール（gh）を使用してissueを取得し、マークダウン形式で保存する。

## 機能要件

### 1. GitHub Issue取得機能

#### 1.1 基本機能
- オープン状態の全てのissueを取得
- issueの内容をマークダウン形式でinboxノートとして保存
- マークダウン化したissueを自動でクローズ
- issue内のコメントも含めて一つのマークダウンファイルに統合

#### 1.2 依存関係
- GitHub CLI（gh）が必須
- ghコマンドがPATHに存在し、認証済みであることが前提
- ghが利用できない環境でもテストが実行可能な設計

### 2. コマンド仕様

#### 2.1 新規コマンド
```bash
krapp import-issues
# または
krapp ii  # エイリアス
```

#### 2.2 オプション
```bash
# 基本実行（ベースパスのremote originリポジトリのissueを取得）
krapp import-issues

# リポジトリを明示的に指定
krapp import-issues --repo owner/repo-name

# ドライランモード（実際にはクローズしない）
krapp import-issues --dry-run

# クローズせずにインポートのみ
krapp import-issues --no-close
```

#### 2.3 デフォルト動作
- コマンド実行時、まずベースパス（設定の`base_dir`）でGitリポジトリの`remote origin`を確認
- `git remote get-url origin`でリポジトリURLを取得し、GitHub URLの場合はそこからowner/repo形式を抽出
- `--repo`オプション未指定時はこのリポジトリのissueを対象とする
- ベースパスがGitリポジトリでない場合やremote originが設定されていない場合はエラー

### 3. 出力形式

#### 3.1 ファイル名形式
```
notes/inbox/YYYY-MM-DD-issue-{issue_number}-{sanitized_title}.md
```

例: `2024-01-15-issue-123-fix-config-loading-bug.md`

#### 3.2 マークダウン構造
```markdown
---
created: 2024-01-15T10:30:00Z
tags: [github-issue, imported]
status: imported
issue_number: 123
issue_url: https://github.com/owner/repo/issues/123
assignees: [user1, user2]
labels: [bug, priority-high]
milestone: v1.2.0
state: closed
original_created: 2024-01-10T09:00:00Z
original_updated: 2024-01-15T10:30:00Z
---

# Issue #123: Fix config loading bug

**Created by:** @username on 2024-01-10  
**Labels:** bug, priority-high  
**Assignees:** @user1, @user2  
**Milestone:** v1.2.0  

## Description

When loading configuration files, the application fails to merge local and global settings properly.

## Comments

### Comment by @reviewer on 2024-01-12

I can reproduce this issue. It seems to be related to the mergo library configuration.

### Comment by @username on 2024-01-14

Fixed in commit abc1234. Please review the changes.

---
*Issue automatically imported and closed by krapp on 2024-01-15T10:30:00Z*
```

## 技術実装

### 4. アーキテクチャ設計

#### 4.1 新規パッケージ構成
```
usecase/
├── github_import.go          # メインロジック
├── github_import_test.go     # テスト
└── github_client.go          # GitHub CLI インターフェース
```

#### 4.2 インターフェース設計
```go
// GitHub操作のインターフェース（テスト容易性のため）
type GitHubClient interface {
    ListOpenIssues(repo string) ([]Issue, error)
    GetIssueComments(repo string, issueNumber int) ([]Comment, error)
    CloseIssue(repo string, issueNumber int) error
    GetCurrentRepo(baseDir string) (string, error)
}

// 実装（gh コマンドを使用）
type GHClient struct{}

// テスト用モック実装
type MockGitHubClient struct {
    Issues       []Issue
    Comments     map[int][]Comment
    ClosedIssues []int
    RepoURL      string  // テスト用のリポジトリURL
}
```

#### 4.3 データ構造
```go
type Issue struct {
    Number      int       `json:"number"`
    Title       string    `json:"title"`
    Body        string    `json:"body"`
    State       string    `json:"state"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    User        User      `json:"user"`
    Assignees   []User    `json:"assignees"`
    Labels      []Label   `json:"labels"`
    Milestone   Milestone `json:"milestone"`
    HTMLURL     string    `json:"html_url"`
}

type Comment struct {
    Body      string    `json:"body"`
    CreatedAt time.Time `json:"created_at"`
    User      User      `json:"user"`
}

type User struct {
    Login string `json:"login"`
}

type Label struct {
    Name  string `json:"name"`
    Color string `json:"color"`
}

type Milestone struct {
    Title string `json:"title"`
}
```

### 5. 実装詳細

#### 5.1 メインロジック（`usecase/github_import.go`）
```go
func ImportGitHubIssues(cfg InboxConfig, client GitHubClient, options ImportOptions) error {
    // 1. リポジトリ情報取得
    var repo string
    var err error
    
    if options.Repo != "" {
        // --repoオプション指定時はそれを使用
        repo = options.Repo
    } else {
        // デフォルト: ベースディレクトリのremote originから取得
        repo, err = client.GetCurrentRepo(cfg.GetBaseDir())
        if err != nil {
            return fmt.Errorf("failed to get current repository: %w", err)
        }
    }
    
    // 2. オープンissue取得
    issues, err := client.ListOpenIssues(repo)
    if err != nil {
        return fmt.Errorf("failed to list issues: %w", err)
    }
    
    // 3. 各issueを処理
    for _, issue := range issues {
        if err := processIssue(cfg, client, repo, issue, options); err != nil {
            log.Printf("failed to process issue #%d: %v", issue.Number, err)
            continue
        }
    }
    
    return nil
}

func processIssue(cfg InboxConfig, client GitHubClient, repo string, issue Issue, options ImportOptions) error {
    // 1. コメント取得
    comments, err := client.GetIssueComments(repo, issue.Number)
    if err != nil {
        return fmt.Errorf("failed to get comments: %w", err)
    }
    
    // 2. マークダウン生成
    markdown := generateIssueMarkdown(issue, comments)
    
    // 3. ファイル作成
    filename := generateIssueFilename(issue)
    filePath := filepath.Join(cfg.GetInboxDir(), filename)
    
    // 4. frontmatter作成
    fm := createIssueFrontMatter(issue)
    
    // 5. ノート保存
    note, err := models.CreateNewNoteWithFrontMatter(models.NewNoteWithFrontMatter{
        Content:     markdown,
        FilePath:    filePath,
        WriteFile:   true,
        FrontMatter: fm,
    })
    if err != nil {
        return fmt.Errorf("failed to create note: %w", err)
    }
    
    // 6. issue クローズ（オプション）
    if !options.DryRun && !options.NoClose {
        if err := client.CloseIssue(repo, issue.Number); err != nil {
            return fmt.Errorf("failed to close issue: %w", err)
        }
    }
    
    return nil
}
```

#### 5.2 GitHub CLI クライント（`usecase/github_client.go`）
```go
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

func (c *GHClient) GetCurrentRepo(baseDir string) (string, error) {
    // ベースディレクトリでgit remote originを確認
    cmd := exec.Command("git", "remote", "get-url", "origin")
    cmd.Dir = baseDir
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("failed to get remote origin: %w", err)
    }
    
    // GitHub URLからowner/repo形式を抽出
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
```

### 6. コマンド統合

#### 6.1 Cobraコマンド追加（`cmd/krapp/krapp/import_issues.go`）
```go
func importIssuesCmd() *cobra.Command {
    var (
        repo    string
        dryRun  bool
        noClose bool
    )
    
    cmd := &cobra.Command{
        Use:     "import-issues",
        Short:   "Import GitHub issues as inbox notes",
        Aliases: []string{"ii"},
        Run: func(cmd *cobra.Command, args []string) {
            cfg := getConfig()
            client := &usecase.GHClient{}
            
            options := usecase.ImportOptions{
                Repo:    repo,
                DryRun:  dryRun,
                NoClose: noClose,
            }
            
            if err := usecase.ImportGitHubIssues(cfg, client, options); err != nil {
                fmt.Printf("GitHub issueのインポートに失敗しました: %v\n", err)
                os.Exit(1)
            }
            
            fmt.Println("GitHub issueのインポートが完了しました")
        },
    }
    
    cmd.Flags().StringVar(&repo, "repo", "", "Repository (owner/name)")
    cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Don't actually close issues")
    cmd.Flags().BoolVar(&noClose, "no-close", false, "Import issues without closing them")
    
    return cmd
}
```

## テスト戦略

### 7. テスト設計

#### 7.1 単体テスト
```go
func TestImportGitHubIssues(t *testing.T) {
    // モッククライアントを使用
    mockClient := &MockGitHubClient{
        Issues: []Issue{
            {Number: 1, Title: "Test Issue", Body: "Test body"},
        },
        Comments: map[int][]Comment{
            1: {{Body: "Test comment", User: User{Login: "reviewer"}}},
        },
    }
    
    cfg := &testConfig{baseDir: "/tmp/test"}
    options := ImportOptions{}
    
    err := ImportGitHubIssues(cfg, mockClient, options)
    assert.NoError(t, err)
    
    // ファイル作成確認
    // issue クローズ確認
}
```

#### 7.2 統合テスト
- モッククライアントを使用した統合テスト（実際のGitHub APIは呼ばない）
- `gh` コマンドの実行可能性チェック（コマンド存在確認のみ）
- GitリポジトリのURL解析テスト（ローカルの`.git`ディレクトリ使用）

**注意**: GitHub APIを実際に呼び出すe2eテストは実装しない（GitHubリポジトリへの影響を避けるため）

### 8. エラーハンドリング

#### 8.1 想定エラーケース
- `gh` コマンドが見つからない
- GitHub認証エラー（`gh auth status`で確認可能）
- リポジトリアクセス権限エラー
- ネットワークエラー
- ファイル書き込みエラー
- ベースディレクトリがGitリポジトリでない
- remote originが設定されていない
- remote originがGitHubリポジトリでない

#### 8.2 エラー対応
- 各エラーケースで適切なエラーメッセージを表示
- 部分的な失敗でも処理を継続
- ログ出力で詳細な情報を提供

## 運用考慮事項

### 9. 制限事項
- GitHub CLI（gh）が必須
- GitHub APIのレート制限に注意
- 大量のissueがある場合の処理時間

### 10. 将来的な拡張
- フィルタリング機能（ラベル、担当者など）
- バッチサイズ制御
- 進捗表示
- 他のGitプラットフォーム対応（GitLab、Bitbucket）

この仕様に基づいてGitHub Issue インポート機能を実装し、krappの機能を拡張する。
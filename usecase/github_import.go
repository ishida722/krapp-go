package usecase

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ishida722/krapp-go/models"
)

// ImportGitHubIssues imports GitHub issues as inbox notes
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

	if len(issues) == 0 {
		log.Println("No open issues found")
		return nil
	}

	log.Printf("Found %d open issues", len(issues))

	// 3. 各issueを処理
	successCount := 0
	for _, issue := range issues {
		if err := processIssue(cfg, client, repo, issue, options); err != nil {
			log.Printf("failed to process issue #%d: %v", issue.Number, err)
			continue
		}
		successCount++
	}

	log.Printf("Successfully processed %d/%d issues", successCount, len(issues))
	return nil
}

// processIssue processes a single issue
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
	filePath := filepath.Join(cfg.GetBaseDir(), cfg.GetInboxDir(), filename)

	// 4. frontmatter作成（issueの作成日時をcreatedに設定）
	fm := createIssueFrontMatter(issue)

	// 5. ノート保存
	_, err = models.CreateNewNoteWithFrontMatter(models.NewNoteWithFrontMatter{
		Content:     markdown,
		FilePath:    filePath,
		WriteFile:   true,
		FrontMatter: fm,
	})
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	log.Printf("Created note for issue #%d: %s", issue.Number, filename)

	// 6. issue クローズ（オプション）
	if !options.DryRun && !options.NoClose {
		if err := client.CloseIssue(repo, issue.Number); err != nil {
			return fmt.Errorf("failed to close issue: %w", err)
		}
		log.Printf("Closed issue #%d", issue.Number)
	}

	return nil
}

// generateIssueFilename generates a filename for the issue
func generateIssueFilename(issue Issue) string {
	// 日付をYYYY-MM-DD形式で取得
	date := issue.CreatedAt.Format("2006-01-02")

	// タイトルをサニタイズ（ファイル名に使用できない文字を除去）
	sanitizedTitle := sanitizeFilename(issue.Title)

	// 長すぎる場合は切り詰め
	if len(sanitizedTitle) > 50 {
		sanitizedTitle = sanitizedTitle[:50]
	}

	return fmt.Sprintf("%s-issue-%d-%s.md", date, issue.Number, sanitizedTitle)
}

// sanitizeFilename removes characters that are not suitable for filenames
func sanitizeFilename(filename string) string {
	// ファイル名として使用できない文字のみを削除（日本語は保持）
	// Windows/macOS/Linuxで使用できない文字: / \ : * ? " < > |
	reg := regexp.MustCompile(`[/\\:*?"<>|]`)
	filename = reg.ReplaceAllString(filename, "-")

	// 連続するハイフンを一つにまとめる
	reg = regexp.MustCompile(`-+`)
	filename = reg.ReplaceAllString(filename, "-")

	// 先頭と末尾のハイフンを削除
	filename = strings.Trim(filename, "-")

	return filename
}

// createIssueFrontMatter creates frontmatter for the issue
func createIssueFrontMatter(issue Issue) models.FrontMatter {
	fm := models.FrontMatter{}

	// issueの作成日時をcreatedに設定
	fm.SetCreated(issue.CreatedAt)

	// 基本情報
	fm["tags"] = []string{"github-issue", "imported"}
	fm["status"] = "imported"
	fm["issue_number"] = issue.Number
	fm["issue_url"] = issue.URL
	fm["state"] = issue.State
	fm["imported_at"] = time.Now().Format(time.RFC3339)
	fm["original_updated"] = issue.UpdatedAt.Format(time.RFC3339)

	// 担当者
	if len(issue.Assignees) > 0 {
		assignees := make([]string, len(issue.Assignees))
		for i, assignee := range issue.Assignees {
			assignees[i] = assignee.Login
		}
		fm["assignees"] = assignees
	}

	// ラベル
	if len(issue.Labels) > 0 {
		labels := make([]string, len(issue.Labels))
		for i, label := range issue.Labels {
			labels[i] = label.Name
		}
		fm["labels"] = labels
	}

	// マイルストーン
	if issue.Milestone.Title != "" {
		fm["milestone"] = issue.Milestone.Title
	}

	return fm
}

// generateIssueMarkdown generates markdown content for the issue
func generateIssueMarkdown(issue Issue, comments []Comment) string {
	var builder strings.Builder

	// タイトル
	builder.WriteString(fmt.Sprintf("# Issue #%d: %s\n\n", issue.Number, issue.Title))

	// メタ情報
	builder.WriteString(fmt.Sprintf("**Created by:** @%s on %s\n",
		issue.Author.Login, issue.CreatedAt.Format("2006-01-02")))

	if len(issue.Labels) > 0 {
		labels := make([]string, len(issue.Labels))
		for i, label := range issue.Labels {
			labels[i] = label.Name
		}
		builder.WriteString(fmt.Sprintf("**Labels:** %s\n", strings.Join(labels, ", ")))
	}

	if len(issue.Assignees) > 0 {
		assignees := make([]string, len(issue.Assignees))
		for i, assignee := range issue.Assignees {
			assignees[i] = "@" + assignee.Login
		}
		builder.WriteString(fmt.Sprintf("**Assignees:** %s\n", strings.Join(assignees, ", ")))
	}

	if issue.Milestone.Title != "" {
		builder.WriteString(fmt.Sprintf("**Milestone:** %s\n", issue.Milestone.Title))
	}

	builder.WriteString("\n")

	// 本文
	if issue.Body != "" {
		builder.WriteString("## Description\n\n")
		builder.WriteString(issue.Body)
		builder.WriteString("\n\n")
	}

	// コメント
	if len(comments) > 0 {
		builder.WriteString("## Comments\n\n")
		for _, comment := range comments {
			builder.WriteString(fmt.Sprintf("### Comment by @%s on %s\n\n",
				comment.Author.Login, comment.CreatedAt.Format("2006-01-02")))
			builder.WriteString(comment.Body)
			builder.WriteString("\n\n")
		}
	}

	// フッター
	builder.WriteString("---\n")
	builder.WriteString(fmt.Sprintf("*Issue automatically imported and closed by krapp on %s*\n",
		time.Now().Format(time.RFC3339)))

	return builder.String()
}

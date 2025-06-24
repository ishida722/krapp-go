# krapp labeling機能 技術仕様書

## 仕様概要

マークダウンファイルに対する対話型ラベリング機能の詳細技術仕様。既存のFront Matter処理機能を活用し、Clean Architectureの原則に従って実装する。

## コマンド仕様

### 基本構文
```bash
krapp labeling <directory> [flags]
```

### フラグ定義
| フラグ | 短縮形 | 型 | デフォルト | 説明 |
|--------|--------|----|-----------|----|
| `--recursive` | `-r` | bool | false | サブディレクトリを再帰的に処理 |

### 使用例
```bash
# 基本使用
krapp labeling ./notes

# 再帰的処理
krapp labeling ./notes -r
krapp labeling ./notes --recursive
```

## データ構造仕様

### ラベル定義
```go
type LabelType string

const (
    LabelDiary  LabelType = "diary"
    LabelIdea   LabelType = "idea" 
    LabelReview LabelType = "review"
)

var LabelMapping = map[string]LabelType{
    "1": LabelDiary,
    "2": LabelIdea,
    "3": LabelReview,
}
```

### ファイル情報構造
```go
type FileInfo struct {
    Path         string
    RelativePath string
    HasLabel     bool
    CurrentLabel *string
    Preview      string
    Size         int64
    ModTime      time.Time
}

type ProcessingResult struct {
    FilePath string
    Action   string // "labeled", "skipped", "error"
    Label    *LabelType
    Error    error
}
```

### 進行状況構造
```go
type ProgressInfo struct {
    CurrentIndex    int
    TotalFiles      int
    ProcessedFiles  int
    SkippedFiles    int
    LabeledFiles    int
    ErrorFiles      int
    StartTime       time.Time
}
```

## インターフェース仕様

### メインインターフェース
```go
type LabelingService interface {
    // ファイル検索・フィルタリング
    FindTargetFiles(directory string, recursive bool) ([]FileInfo, error)
    
    // ファイル処理
    ProcessFile(filePath string) (*FileInfo, error)
    SetLabel(filePath string, label LabelType) error
    
    // プレビュー生成
    GeneratePreview(content string, maxLines int) string
    
    // 進行状況管理
    UpdateProgress(result ProcessingResult)
    GetProgress() ProgressInfo
}
```

### コマンドインターフェース
```go
type LabelingCommand interface {
    // コマンド実行
    Execute(args []string) error
    
    // ユーザー入力処理
    PromptUser(fileInfo FileInfo, progress ProgressInfo) (UserChoice, error)
    
    // 表示制御
    DisplayFile(fileInfo FileInfo, progress ProgressInfo)
    DisplaySummary(progress ProgressInfo)
}

type UserChoice struct {
    Action string // "label", "skip", "quit"
    Label  *LabelType
}
```

## ファイル処理仕様

### 対象ファイル判定
```go
func IsTargetFile(path string) bool {
    // 拡張子チェック
    if !strings.HasSuffix(strings.ToLower(path), ".md") {
        return false
    }
    
    // 隠しファイル除外
    if strings.HasPrefix(filepath.Base(path), ".") {
        return false
    }
    
    return true
}
```

### ファイル検索アルゴリズム
```go
func FindFiles(directory string, recursive bool) ([]string, error) {
    var files []string
    
    walkFunc := func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        // ディレクトリ処理
        if info.IsDir() {
            if !recursive && path != directory {
                return filepath.SkipDir
            }
            return nil
        }
        
        // ファイル判定
        if IsTargetFile(path) {
            files = append(files, path)
        }
        
        return nil
    }
    
    err := filepath.Walk(directory, walkFunc)
    return files, err
}
```

### Front Matter処理仕様

#### ラベル存在チェック
```go
func HasLabel(content string) (bool, *string) {
    frontMatter, err := extractFrontMatter(content)
    if err != nil {
        return false, nil
    }
    
    if label, exists := frontMatter["label"]; exists {
        if labelStr, ok := label.(string); ok {
            return true, &labelStr
        }
    }
    
    return false, nil
}
```

#### ラベル設定処理
```go
func SetLabel(filePath string, label LabelType) error {
    // ファイル読み込み
    content, err := ioutil.ReadFile(filePath)
    if err != nil {
        return err
    }
    
    // Front Matter更新
    updatedContent, err := updateFrontMatterLabel(string(content), string(label))
    if err != nil {
        return err
    }
    
    // 原子的書き込み
    return writeFileAtomic(filePath, []byte(updatedContent))
}

func updateFrontMatterLabel(content string, label string) (string, error) {
    frontMatter, body, err := parseFrontMatter(content)
    if err != nil {
        return "", err
    }
    
    // ラベル設定
    if frontMatter == nil {
        frontMatter = make(map[string]interface{})
    }
    frontMatter["label"] = label
    
    // Front Matter再構築
    return rebuildContent(frontMatter, body)
}
```

## UI/UX仕様

### 表示フォーマット
```
File: notes/daily/2024-06-24.md (1/5) [2 files skipped]
─────────────────────────────────────────────────────
# 今日の振り返り

今日は新しい機能の設計を考えた。
Clean Architectureの原則に従って...
[Preview truncated - showing first 10 lines]
─────────────────────────────────────────────────────
Labels: 1=diary, 2=idea, 3=review, s=skip, q=quit
Your choice: _
```

### プレビュー生成仕様
```go
const (
    MaxPreviewLines = 20
    MaxLineLength   = 100
    TruncateMarker  = "[Preview truncated - showing first %d lines]"
)

func GeneratePreview(content string, maxLines int) string {
    lines := strings.Split(content, "\n")
    
    // Front Matterをスキップ
    startIndex := 0
    if strings.HasPrefix(content, "---") {
        for i := 1; i < len(lines); i++ {
            if strings.HasPrefix(lines[i], "---") || strings.HasPrefix(lines[i], "...") {
                startIndex = i + 1
                break
            }
        }
    }
    
    // プレビュー行数制限
    endIndex := startIndex + maxLines
    if endIndex >= len(lines) {
        endIndex = len(lines)
    }
    
    previewLines := lines[startIndex:endIndex]
    
    // 行長制限
    for i, line := range previewLines {
        if len(line) > MaxLineLength {
            previewLines[i] = line[:MaxLineLength-3] + "..."
        }
    }
    
    preview := strings.Join(previewLines, "\n")
    
    // 切り詰めマーカー
    if endIndex < len(lines) {
        preview += "\n" + fmt.Sprintf(TruncateMarker, maxLines)
    }
    
    return preview
}
```

### ユーザー入力処理
```go
func PromptUser() (UserChoice, error) {
    fmt.Print("Your choice: ")
    
    reader := bufio.NewReader(os.Stdin)
    input, err := reader.ReadString('\n')
    if err != nil {
        return UserChoice{}, err
    }
    
    input = strings.TrimSpace(strings.ToLower(input))
    
    switch input {
    case "1", "2", "3":
        if label, exists := LabelMapping[input]; exists {
            return UserChoice{Action: "label", Label: &label}, nil
        }
    case "s", "skip":
        return UserChoice{Action: "skip"}, nil
    case "q", "quit":
        return UserChoice{Action: "quit"}, nil
    default:
        return UserChoice{}, fmt.Errorf("invalid input: %s", input)
    }
    
    return UserChoice{}, fmt.Errorf("unexpected input: %s", input)
}
```

## エラーハンドリング仕様

### エラー分類
```go
type LabelingError struct {
    Type    ErrorType
    Message string
    Cause   error
    Context map[string]interface{}
}

type ErrorType int

const (
    ErrorTypeFileSystem ErrorType = iota
    ErrorTypeParsing
    ErrorTypeUserInput
    ErrorTypePermission
    ErrorTypeInternal
)
```

### エラー処理戦略
| エラータイプ | 処理方針 | ユーザー通知 |
|-------------|----------|-------------|
| ディレクトリ不存在 | 即座に終了 | エラーメッセージ表示 |
| ファイル読み込み失敗 | スキップして続行 | 警告表示 |
| Front Matter解析失敗 | スキップして続行 | 警告表示 |
| ファイル書き込み失敗 | 処理中断 | エラーメッセージ表示 |
| 無効なユーザー入力 | 再入力要求 | ヘルプ表示 |

### 原子的ファイル操作
```go
func writeFileAtomic(filename string, data []byte) error {
    tempFile := filename + ".tmp"
    
    // 一時ファイルに書き込み
    if err := ioutil.WriteFile(tempFile, data, 0644); err != nil {
        return err
    }
    
    // 原子的なリネーム
    if err := os.Rename(tempFile, filename); err != nil {
        os.Remove(tempFile)
        return err
    }
    
    return nil
}
```

## パフォーマンス仕様

### 要件
- 1000ファイルの処理: 10秒以内
- メモリ使用量: 100MB以下
- ファイルプレビュー表示: 100ms以内

### 最適化手法
1. **ストリーミング読み込み**: 大きなファイルの部分読み込み
2. **並列ファイル検索**: Goroutineを使用した高速検索
3. **メモリ効率**: ファイル内容の使い回し最小化

### ベンチマーク指標
```go
func BenchmarkFileProcessing(b *testing.B) {
    // 標準的な処理時間計測
}

func BenchmarkPreviewGeneration(b *testing.B) {
    // プレビュー生成時間計測
}
```

## テスト仕様

### テストカバレッジ目標
- **単体テスト**: 90%以上
- **統合テスト**: 主要フロー100%
- **E2Eテスト**: 基本シナリオ100%

### テストデータ構造
```
testdata/
├── valid_files/
│   ├── with_frontmatter.md
│   ├── without_frontmatter.md
│   └── with_label.md
├── invalid_files/
│   ├── corrupted_frontmatter.md
│   └── binary_file.jpg
└── directory_structure/
    ├── subdir1/
    └── subdir2/
```

### テストケース
```go
func TestLabelingService_FindTargetFiles(t *testing.T) {
    tests := []struct {
        name      string
        directory string
        recursive bool
        expected  int
        wantErr   bool
    }{
        {"basic directory", "testdata/valid_files", false, 3, false},
        {"recursive search", "testdata", true, 5, false},
        {"non-existent dir", "nonexistent", false, 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テスト実装
        })
    }
}
```

## セキュリティ仕様

### ファイルアクセス制御
```go
func validatePath(basePath, targetPath string) error {
    // パストラバーサル攻撃防止
    absBase, err := filepath.Abs(basePath)
    if err != nil {
        return err
    }
    
    absTarget, err := filepath.Abs(targetPath)
    if err != nil {
        return err
    }
    
    if !strings.HasPrefix(absTarget, absBase) {
        return fmt.Errorf("path traversal detected: %s", targetPath)
    }
    
    return nil
}
```

### 権限チェック
```go
func checkPermissions(filePath string) error {
    info, err := os.Stat(filePath)
    if err != nil {
        return err
    }
    
    // 読み取り権限チェック
    if info.Mode().Perm()&0400 == 0 {
        return fmt.Errorf("no read permission: %s", filePath)
    }
    
    // 書き込み権限チェック
    if info.Mode().Perm()&0200 == 0 {
        return fmt.Errorf("no write permission: %s", filePath)
    }
    
    return nil
}
```

## 設定仕様

### 将来の設定拡張
```yaml
# .krapp_config.yaml
labeling:
  preview_lines: 20
  max_line_length: 100
  custom_labels:
    - name: "meeting"
      key: "4"
    - name: "task" 
      key: "5"
  auto_backup: true
  confirm_destructive: true
```

### 設定構造体
```go
type LabelingConfig struct {
    PreviewLines   int              `yaml:"preview_lines"`
    MaxLineLength  int              `yaml:"max_line_length"`
    CustomLabels   []CustomLabel    `yaml:"custom_labels"`
    AutoBackup     bool             `yaml:"auto_backup"`
    ConfirmDestruct bool            `yaml:"confirm_destructive"`
}

type CustomLabel struct {
    Name string `yaml:"name"`
    Key  string `yaml:"key"`
}
```

## ログ・監視仕様

### ログレベル
- **INFO**: 処理開始・完了、統計情報
- **WARN**: ファイルスキップ、軽微なエラー
- **ERROR**: 処理中断を伴うエラー
- **DEBUG**: 詳細なフロー情報（開発時のみ）

### メトリクス収集
```go
type Metrics struct {
    ProcessedFiles int
    LabeledFiles   int
    SkippedFiles   int
    ErrorFiles     int
    ProcessingTime time.Duration
}
```

## 互換性仕様

### Goバージョン
- **最小要件**: Go 1.19
- **推奨**: Go 1.21+

### OS互換性
- **サポート**: macOS, Linux
- **制限付きサポート**: Windows（パス処理の違い）

### 既存機能との互換性
- 既存のFront Matter形式を完全保持
- 既存の設定システムとの統合
- 既存コマンドとの名前空間分離

## 性能計測・ベンチマーク

### 測定項目
1. ファイル検索時間（ファイル数別）
2. プレビュー生成時間（ファイルサイズ別）
3. ラベル設定時間（Front Matterサイズ別）
4. メモリ使用量（処理ファイル数別）

### 性能目標
| 項目 | 目標値 |
|------|--------|
| 100ファイル処理 | < 3秒 |
| 1000ファイル処理 | < 15秒 |
| プレビュー生成 | < 50ms |
| メモリ使用量 | < 50MB |
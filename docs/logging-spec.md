# ログ機能仕様書

## 概要

krappにログ機能を追加し、アプリケーションの動作状況、エラー、デバッグ情報を記録する機能を提供する。開発時のデバッグ支援とユーザーサポートの向上を目的とする。

## 機能要件

### 1. ログレベル

#### 1.1 サポートするログレベル
- **DEBUG**: 開発時のデバッグ情報（デフォルトでは出力しない）
- **INFO**: 一般的な情報メッセージ（ノート作成、設定読み込みなど）
- **WARN**: 警告メッセージ（設定の問題、推奨されない使用方法など）
- **ERROR**: エラーメッセージ（処理失敗、ファイルアクセスエラーなど）

#### 1.2 デフォルト設定
- 通常運用: WARN以上を出力
- 開発時: DEBUG以上を出力（環境変数で制御）

### 2. ログ出力先

#### 2.1 標準エラー出力（stderr）
- ユーザー向けメッセージは標準出力（stdout）と分離
- パイプやリダイレクトでの使用を考慮

#### 2.2 ログファイル（オプション）
- 設定で有効化した場合のみファイル出力
- ローテーション機能（サイズまたは日付ベース）
- デフォルト場所: `~/.cache/krapp/logs/krapp.log`（XDG準拠）

### 3. ログフォーマット

#### 3.1 標準フォーマット
```
2024-01-15T10:30:00Z [INFO] create_daily: デイリーノートを作成しました: /path/to/note.md
2024-01-15T10:30:01Z [ERROR] config: 設定ファイルの読み込みに失敗しました: file not found
```

#### 3.2 フォーマット要素
- タイムスタンプ（RFC3339形式、UTC）
- ログレベル
- コンポーネント名（create_daily, config, sync等）
- メッセージ

#### 3.3 開発モード
- ファイル名と行番号を追加
```
2024-01-15T10:30:00Z [DEBUG] create_daily (create_daily.go:45): ディレクトリを作成中: /path/to/dir
```

## 設定

### 4. 設定オプション

#### 4.1 Config構造体への追加
```go
type Config struct {
    // 既存フィールド...
    
    // ログ設定
    LogLevel    string `yaml:"log_level"`     // debug, info, warn, error
    LogToFile   bool   `yaml:"log_to_file"`   // ファイル出力の有効化
    LogFilePath string `yaml:"log_file_path"` // ログファイルパス（空の場合はデフォルト）
    LogMaxSize  int    `yaml:"log_max_size"`  // ログファイル最大サイズ(MB)
    LogMaxFiles int    `yaml:"log_max_files"` // 保持するログファイル数
}
```

#### 4.2 デフォルト設定
```go
var defaultConfig = Config{
    // 既存設定...
    
    LogLevel:    "warn",
    LogToFile:   false,
    LogFilePath: "",      // 空の場合は ~/.cache/krapp/logs/krapp.log
    LogMaxSize:  10,      // 10MB
    LogMaxFiles: 5,       // 5ファイル
}
```

#### 4.3 設定例
```yaml
# ~/.config/krapp/config.yaml
log_level: info
log_to_file: true
log_file_path: "./logs/krapp.log"  # プロジェクト固有のログ
log_max_size: 5
log_max_files: 3
```

### 5. 環境変数

#### 5.1 開発時の制御
- `KRAPP_DEBUG=1`: DEBUGレベルを有効化
- `KRAPP_LOG_LEVEL=debug`: ログレベルを直接指定（設定ファイルより優先）

## 技術実装

### 6. アーキテクチャ設計

#### 6.1 ログパッケージ
```
internal/
└── log/
    ├── logger.go      # メインのロガー実装
    ├── config.go      # ログ設定
    ├── level.go       # ログレベル定義
    └── writer.go      # 出力先管理
```

#### 6.2 インターフェース設計
```go
package log

type Logger interface {
    Debug(component, message string)
    Debugf(component, format string, args ...interface{})
    Info(component, message string)
    Infof(component, format string, args ...interface{})
    Warn(component, message string)
    Warnf(component, format string, args ...interface{})
    Error(component, message string)
    Errorf(component, format string, args ...interface{})
}

type Level int

const (
    LevelDebug Level = iota
    LevelInfo
    LevelWarn
    LevelError
)
```

### 7. 実装詳細

#### 7.1 ロガー初期化（`internal/log/logger.go`）
```go
type logger struct {
    level      Level
    writer     io.Writer
    fileWriter io.WriteCloser
    mutex      sync.Mutex
    debug      bool // 開発モード（ファイル名・行番号表示）
}

func NewLogger(cfg LogConfig) (Logger, error) {
    l := &logger{
        level: parseLevel(cfg.Level),
        debug: cfg.Debug,
    }
    
    // 標準エラー出力は常に有効
    writers := []io.Writer{os.Stderr}
    
    // ファイル出力の設定
    if cfg.ToFile {
        file, err := openLogFile(cfg.FilePath, cfg.MaxSize, cfg.MaxFiles)
        if err != nil {
            return nil, fmt.Errorf("failed to open log file: %w", err)
        }
        l.fileWriter = file
        writers = append(writers, file)
    }
    
    l.writer = io.MultiWriter(writers...)
    return l, nil
}

func (l *logger) log(level Level, component, message string) {
    if level < l.level {
        return
    }
    
    l.mutex.Lock()
    defer l.mutex.Unlock()
    
    timestamp := time.Now().UTC().Format(time.RFC3339)
    levelStr := level.String()
    
    var output string
    if l.debug {
        _, file, line, _ := runtime.Caller(2)
        output = fmt.Sprintf("%s [%s] %s (%s:%d): %s\n", 
            timestamp, levelStr, component, filepath.Base(file), line, message)
    } else {
        output = fmt.Sprintf("%s [%s] %s: %s\n", 
            timestamp, levelStr, component, message)
    }
    
    l.writer.Write([]byte(output))
}
```

#### 7.2 設定との統合（`config/config.go`）
```go
type LogConfig struct {
    Level    string
    ToFile   bool
    FilePath string
    MaxSize  int
    MaxFiles int
    Debug    bool
}

func (c *Config) GetLogConfig() LogConfig {
    logLevel := c.LogLevel
    if envLevel := os.Getenv("KRAPP_LOG_LEVEL"); envLevel != "" {
        logLevel = envLevel
    }
    
    debug := os.Getenv("KRAPP_DEBUG") == "1"
    
    filePath := c.LogFilePath
    if filePath == "" && c.LogToFile {
        // デフォルトパスを生成
        cacheDir, _ := os.UserCacheDir()
        filePath = filepath.Join(cacheDir, "krapp", "logs", "krapp.log")
    }
    
    return LogConfig{
        Level:    logLevel,
        ToFile:   c.LogToFile,
        FilePath: filePath,
        MaxSize:  c.LogMaxSize,
        MaxFiles: c.LogMaxFiles,
        Debug:    debug,
    }
}
```

#### 7.3 アプリケーション統合（`cmd/krapp/main.go`）
```go
func main() {
    // 設定読み込み
    cfg := loadConfig()
    
    // ロガー初期化
    logger, err := log.NewLogger(cfg.GetLogConfig())
    if err != nil {
        fmt.Fprintf(os.Stderr, "ログシステムの初期化に失敗しました: %v\n", err)
        os.Exit(1)
    }
    defer logger.Close()
    
    // グローバルロガーとして設定
    log.SetGlobal(logger)
    
    // CLI実行
    if err := rootCmd.Execute(); err != nil {
        logger.Error("main", fmt.Sprintf("コマンド実行エラー: %v", err))
        os.Exit(1)
    }
}
```

### 8. 使用例

#### 8.1 ユースケースでの使用
```go
// usecase/create_daily.go
func CreateDailyNote(cfg Config, now time.Time) (string, error) {
    log.Info("create_daily", "デイリーノート作成を開始")
    log.Debugf("create_daily", "作成日時: %s", now.Format("2006-01-02"))
    
    // ディレクトリ作成
    dir := cfg.GetDailyNoteDir(now)
    if err := os.MkdirAll(dir, 0755); err != nil {
        log.Errorf("create_daily", "ディレクトリ作成に失敗: %v", err)
        return "", fmt.Errorf("failed to create directory: %w", err)
    }
    log.Debugf("create_daily", "ディレクトリを作成: %s", dir)
    
    // ファイル作成
    filePath := filepath.Join(dir, now.Format("2006-01-02")+".md")
    note, err := models.CreateNewNote(models.NewNote{
        Content:   "",
        FilePath:  filePath,
        WriteFile: true,
    })
    if err != nil {
        log.Errorf("create_daily", "ノート作成に失敗: %v", err)
        return "", fmt.Errorf("failed to create note: %w", err)
    }
    
    log.Infof("create_daily", "デイリーノートを作成しました: %s", filePath)
    return filePath, nil
}
```

#### 8.2 設定読み込みでの使用
```go
// config/config.go
func LoadConfig() (*Config, error) {
    log.Debug("config", "設定ファイルの読み込みを開始")
    
    cfg := defaultConfig
    
    // グローバル設定
    globalPath := getGlobalConfigPath()
    if globalCfg, err := loadConfigFile(globalPath); err == nil {
        log.Infof("config", "グローバル設定を読み込み: %s", globalPath)
        mergo.Merge(&cfg, globalCfg, mergo.WithOverride)
    } else {
        log.Debugf("config", "グローバル設定なし: %s", globalPath)
    }
    
    // ローカル設定
    localPath := "./.krapp_config.yaml"
    if localCfg, err := loadConfigFile(localPath); err == nil {
        log.Infof("config", "ローカル設定を読み込み: %s", localPath)
        mergo.Merge(&cfg, localCfg, mergo.WithOverride)
    } else {
        log.Debugf("config", "ローカル設定なし: %s", localPath)
    }
    
    log.Info("config", "設定読み込み完了")
    return &cfg, nil
}
```

## テスト戦略

### 9. テスト設計

#### 9.1 単体テスト
```go
func TestLogger(t *testing.T) {
    var buf bytes.Buffer
    
    cfg := LogConfig{
        Level:  "info",
        ToFile: false,
        Debug:  false,
    }
    
    logger := &logger{
        level:  LevelInfo,
        writer: &buf,
    }
    
    logger.Info("test", "テストメッセージ")
    
    output := buf.String()
    assert.Contains(t, output, "[INFO]")
    assert.Contains(t, output, "test:")
    assert.Contains(t, output, "テストメッセージ")
}
```

#### 9.2 統合テスト
- 設定ファイルからのログ設定読み込み
- 環境変数によるログレベル制御
- ファイル出力とローテーション
- 複数goroutineからの同時書き込み

### 10. 運用考慮事項

#### 10.1 パフォーマンス
- ログ出力のオーバーヘッドを最小化
- 高頻度なDEBUGログは本番では無効化
- ファイルI/Oの非同期化（将来的な改善）

#### 10.2 セキュリティ
- ログにセンシティブな情報（パスワード、トークン）を記録しない
- ファイルパーミッションの適切な設定
- ログローテーションでの古いファイル削除

#### 10.3 トラブルシューティング
- ログファイルの場所をヘルプメッセージに表示
- 設定問題時のフォールバック動作
- ログ出力エラー時の適切な処理

## 段階的実装

### Phase 1: 基本ログ機能
1. `internal/log`パッケージの実装
2. 標準エラー出力への基本ログ
3. 既存コードへのログ追加（主要な処理のみ）

### Phase 2: 設定統合
1. 設定ファイルとの統合
2. 環境変数サポート
3. ログレベル制御

### Phase 3: ファイル出力
1. ログファイル出力機能
2. ローテーション機能
3. XDGディレクトリ準拠

### Phase 4: 拡張機能
1. 構造化ログ（JSON形式）のサポート
2. 外部ログシステム連携
3. ログ解析ツールの提供

このログ機能により、krappの運用性とメンテナビリティが大幅に向上し、ユーザーサポートと開発効率の改善が期待される。
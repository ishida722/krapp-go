# テンプレート機能仕様書

## 概要

krappにテンプレート機能を追加し、ノート作成時にfrontmatterに事前定義した属性を自動で含める機能を提供する。

## 機能要件

### 1. 設定ファイルでのテンプレート定義

#### 1.1 設定項目の追加

既存の`Config`構造体に以下のフィールドを追加する：

```go
type Config struct {
    // 既存フィールド...
    
    // テンプレート関連
    DailyTemplate map[string]any `yaml:"daily_template"`  // デイリーノート用テンプレート
    InboxTemplate map[string]any `yaml:"inbox_template"`  // インボックスノート用テンプレート
}
```

#### 1.2 設定ファイルの記述例

```yaml
# ~/.krapp_config_global.yaml または .krapp_config_local.yaml

# 既存設定...
base_dir: "./notes"
daily_note_dir: "daily"
inbox_dir: "inbox"
editor: "vim"

# テンプレート設定
daily_template:
  tags: []
  status: "draft"
  priority: null
  
inbox_template:
  tags: []
  status: "new"
  category: "inbox"
  priority: "medium"
```

### 2. ノート作成時の処理変更

#### 2.1 デイリーノート作成（`usecase/create_daily.go`）

`CreateDailyNote`関数を修正し、設定からテンプレートを読み込んでfrontmatterに適用する。

```go
func CreateDailyNote(cfg Config, now time.Time) (string, error) {
    // 既存のディレクトリ作成処理...
    
    // テンプレートから初期frontmatterを作成
    fm := models.FrontMatter{}
    
    // 作成日時を設定
    fm.SetCreated(now)
    
    // テンプレートの属性を追加
    if cfg.DailyTemplate != nil {
        for key, value := range cfg.DailyTemplate {
            fm[key] = value
        }
    }
    
    note, err := models.CreateNewNoteWithFrontMatter(models.NewNote{
        Content:     "",
        FilePath:    filepath.Join(dir, date+".md"),
        WriteFile:   true,
        FrontMatter: fm,
    })
    // ...
}
```

#### 2.2 インボックスノート作成（`usecase/create_inbox.go`）

`CreateInboxNote`関数を修正し、同様にテンプレートを適用する。

```go
func CreateInboxNote(cfg InboxConfig, now time.Time, title string) (string, error) {
    // 既存の処理...
    
    // テンプレートから初期frontmatterを作成
    fm := models.FrontMatter{}
    
    // 作成日時を設定
    fm.SetCreated(now)
    
    // テンプレートの属性を追加
    if cfg.InboxTemplate != nil {
        for key, value := range cfg.InboxTemplate {
            fm[key] = value
        }
    }
    
    note, err := models.CreateNewNoteWithFrontMatter(models.NewNote{
        Content:     "",
        FilePath:    filepath.Join(dir, filename),
        WriteFile:   true,
        FrontMatter: fm,
    })
    // ...
}
```

### 3. モデルの拡張

#### 3.1 `models/note.go`の拡張

`CreateNewNote`関数のバリエーションとして、frontmatterを直接指定できる関数を追加する。

```go
type NewNoteWithFrontMatter struct {
    Content     string
    FilePath    string
    WriteFile   bool
    FrontMatter FrontMatter
}

func CreateNewNoteWithFrontMatter(newNote NewNoteWithFrontMatter) (*Note, error) {
    note := &Note{
        FrontMatter: newNote.FrontMatter,
        Content:     strings.TrimSpace(newNote.Content),
        FilePath:    newNote.FilePath,
    }
    
    if newNote.WriteFile {
        if note.FilePath == "" {
            return note, errors.New("note file path is empty")
        }
        if err := note.SaveToFile(); err != nil {
            return note, fmt.Errorf("failed to save note to file: %w", err)
        }
    }
    
    return note, nil
}
```

### 4. 設定の継承とマージ

#### 4.1 テンプレート設定のマージ

既存の`MergeConfig`関数は`mergo`ライブラリを使用してマップフィールドも適切にマージするため、テンプレート設定も自動的に継承される。

- デフォルト設定 → グローバル設定 → ローカル設定の順でマージ
- ローカル設定でテンプレートを部分的に上書き可能

#### 4.2 デフォルトテンプレート

```go
var defaultConfig = Config{
    // 既存設定...
    
    DailyTemplate: map[string]any{
        "tags": []string{},
    },
    InboxTemplate: map[string]any{
        "tags": []string{},
        "status": "new",
    },
}
```

## 実装詳細

### Phase 1: 基本機能実装

1. `config/config.go`の`Config`構造体にテンプレートフィールドを追加
2. デフォルト設定にテンプレートを追加
3. `models/note.go`に`CreateNewNoteWithFrontMatter`関数を追加
4. `usecase/create_daily.go`と`usecase/create_inbox.go`を修正してテンプレート適用

### Phase 2: 拡張機能

1. 動的値の対応（例：作成日時、曜日など）
2. テンプレート検証機能
3. テンプレート設定のコマンドライン操作

## テスト戦略

### 単体テスト

1. `config_test.go`: テンプレート設定の読み込みとマージ
2. `note_test.go`: テンプレート付きノート作成
3. `create_daily_test.go`: デイリーノートテンプレート適用
4. `create_inbox_test.go`: インボックスノートテンプレート適用

### 統合テスト

1. 設定ファイルからテンプレート読み込み → ノート作成 → frontmatter確認
2. グローバル/ローカル設定のマージ確認

## 互換性

- 既存の設定ファイルには影響なし（テンプレート未設定時はcreatedのみ設定）
- 既存のノート作成フローとの互換性を維持
- 段階的な機能追加により既存機能への影響を最小化

## 使用例

### 設定例

```yaml
# プロジェクト用設定 (.krapp_config_local.yaml)
daily_template:
  tags: ["daily", "work"]
  status: "active"
  project: "krapp-development"
  
inbox_template:
  tags: ["inbox"]
  status: "review"
  priority: "normal"
  project: "krapp-development"
```

### 生成されるノート例

```markdown
---
created: 2024-01-15
tags: [daily, work]
status: active
project: krapp-development
---

# 2024-01-15の日記

```

このテンプレート機能により、ユーザーは自分のワークフローに合わせたノート構造を事前定義でき、一貫性のあるノート管理が可能になる。
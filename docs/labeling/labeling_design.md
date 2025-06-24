# krapp labeling機能 設計書

## 概要

krappにファイルラベリング機能を追加する。指定されたディレクトリ内のマークダウンファイルを順次表示し、ユーザーの入力に基づいてFront Matterにラベルを設定する対話型CLI機能。

## 機能要件

### 基本機能
- 指定ディレクトリ内の`.md`ファイルの一覧表示とラベル付け
- 既存Front Matterの検出とラベル設定状況の確認
- 3種類の定義済みラベル（diary/idea/review）の選択的設定
- スキップ・終了機能による柔軟な処理制御

### 拡張機能
- 再帰的ディレクトリ処理（`-r`オプション）
- 既にラベル付きファイルの自動スキップ
- 進行状況の表示

## アーキテクチャ設計

### Clean Architecture準拠
```
cmd/krapp (Interface Layer)
  └── krapp/labeling.go
      ↓
usecase (Use Case Layer)  
  └── labeling_interactive.go
      ↓
models (Entity Layer)
  └── note.go (既存)
  └── front_matter.go (既存)
```

### 責務分離
- **cmd層**: CLI引数解析、ユーザー入力受付、進行状況表示
- **usecase層**: ビジネスロジック（ファイル検索、ラベル設定処理）
- **models層**: データ構造とドメインルール（既存Front Matter処理を活用）

## データフロー

### 1. 初期化フェーズ
```
ディレクトリ指定 → ファイル検索 → ラベル状況確認 → 処理対象リスト作成
```

### 2. 対話処理フェーズ
```
ファイル表示 → ユーザー入力 → ラベル設定 → 次ファイル → 完了判定
```

### 3. 完了フェーズ
```
統計表示 → 終了
```

## インターフェース設計

### CLI Interface
```go
type LabelingCommand struct {
    Directory string
    Recursive bool
}
```

### Use Case Interface
```go
type LabelingInteractor interface {
    FindTargetFiles(dir string, recursive bool) ([]string, error)
    ProcessFile(filePath string) (*FilePreview, error)
    SetLabel(filePath string, label string) error
    GetProgress() ProgressInfo
}

type FilePreview struct {
    Path     string
    Content  string
    HasLabel bool
}

type ProgressInfo struct {
    Current int
    Total   int
    Skipped int
}
```

## 既存コードとの統合

### 活用する既存機能
1. **models/front_matter.go**: Front Matter解析・編集
2. **models/note.go**: ファイル操作とYAML処理
3. **config/config.go**: 設定管理（将来の拡張用）

### 既存パターンの踏襲
- テスト駆動開発（`*_test.go`ファイル）
- エラーハンドリングパターン
- 設定システムの活用

## UI/UX設計

### 表示レイアウト
```
File: notes/daily/2024-06-24.md (1/5) [2 files skipped]
─────────────────────────────────────────────────────
# 今日の振り返り

今日は新しい機能の設計を考えた。
Clean Architectureの原則に従って...

─────────────────────────────────────────────────────
Labels: 1=diary, 2=idea, 3=review, s=skip, q=quit
Your choice: _
```

### ユーザー体験
- 直感的なキー操作（数字キー + s/q）
- 即座のフィードバック（選択後すぐ次ファイルへ）
- 明確な進行状況表示

## エラーハンドリング戦略

### ファイルシステムエラー
- ディレクトリ不存在 → 明確なエラーメッセージ
- 読み込み権限なし → ファイルスキップして続行
- 書き込み権限なし → エラー表示して中断

### データエラー
- 不正なFront Matter → 警告表示してスキップ
- ファイル破損 → エラーログ出力してスキップ

### ユーザー入力エラー
- 無効なキー入力 → 再入力プロンプト
- 予期しない中断 → 進行状況保存（将来拡張）

## パフォーマンス考慮事項

### メモリ効率
- ファイルのストリーミング読み込み
- 大きなファイルの部分表示（先頭N行）

### レスポンス性
- ファイル検索の並列処理
- プレビュー表示の高速化

### スケーラビリティ
- 大量ファイル（1000+）への対応
- メモリ使用量の制限

## セキュリティ考慮事項

### ファイルアクセス制御
- 指定ディレクトリ外へのアクセス防止
- シンボリックリンク追跡の制限

### データ保護
- バックアップ機能（将来拡張）
- 原子的な書き込み操作

## 今後の拡張性

### Phase 2 機能
- カスタムラベルの定義
- 既存ラベルの変更機能
- バッチ処理モード

### Phase 3 機能
- ラベル統計・分析機能
- フィルタリング・検索機能
- 外部ツール連携

## 実装方針

### 開発順序
1. **Core機能**: 基本的なラベル設定機能
2. **Interactive UI**: 対話型インターフェース
3. **Error Handling**: 包括的なエラー処理
4. **Testing**: 単体・統合テスト
5. **Documentation**: ユーザーガイド

### 品質保証
- TDD（Test-Driven Development）
- 既存テストパターンの踏襲
- Clean Architectureの原則遵守

## 技術的制約・前提条件

### 依存関係
- Go 1.19+
- gopkg.in/yaml.v3（既存依存）
- github.com/spf13/cobra（既存依存）

### 環境要件
- Unix系OS（macOS, Linux）
- ターミナル環境
- ファイルシステム書き込み権限

## リスク分析

### 技術リスク
- **中**: Front Matter破損の可能性 → バックアップ機能で軽減
- **低**: パフォーマンス問題 → ストリーミング処理で対応

### ユーザビリティリスク
- **中**: 操作の複雑さ → シンプルなUI設計で軽減
- **低**: 誤操作 → 確認プロンプトで防止

### 運用リスク
- **低**: データ損失 → 原子的操作で防止
- **低**: 互換性問題 → 既存機能の活用で軽減
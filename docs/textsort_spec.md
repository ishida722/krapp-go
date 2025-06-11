# テキスト仕分けCLIツール 仕様書
2025/06/08

## 概要
指定フォルダ内のテキストファイルにラベルを付与し、ラベルごとに対応するフォルダへ自動で仕分けるCLIツール。

---

## 機能一覧

### 1. ラベル付け機能
- 設定ファイル（例: `labels.yaml`）からラベルリストを読み込む  
  - 例:  
    ```yaml
    labels:
      - key: 0
        name: 日記
        path: diary
      - key: 1
        name: 雑記
        path: notes
    ```
- 指定フォルダ内のテキストファイルを順番に処理
- 1ファイルずつ内容を表示
- ユーザーが対応するキー（0-9）でラベルを選択
- 選択したラベルをyaml front matterとしてファイル先頭に追記
- すでにラベルが付与されているファイルはスキップ
- 全ファイル処理まで繰り返し

### 2. 振り分け機能
- ラベルごとに、対応するパス（サブフォルダ）へファイルを移動

---

## 設定ファイル仕様例（labels.yaml）

```yaml
labels:
  - key: 0
    name: 日記
    path: diary
  - key: 1
    name: 雑記
    path: notes
  - key: 2
    name: 記事
    path: article
```

---

## コマンド例

- ラベル付け:  
  ```
  $ textsort label --dir ./0.Inbox --config labels.yaml
  ```
- 振り分け:  
  ```
  $ textsort sort --dir ./0.Inbox --config labels.yaml
  ```

---

## yaml front matter 例

```markdown
---
label: 日記
---
本文...
```

---

## その他
- 拡張子は`.md`や`.txt`を対象
- 既存のyaml front matterがある場合は上書きしない
- ラベル付け後、即時振り分けも可能（オプション）

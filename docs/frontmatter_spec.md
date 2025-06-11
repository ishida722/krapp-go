# yaml frontmatter編集ユースケース仕様

## 概要
Markdownテキストデータの先頭にあるyaml frontmatter（`---`で囲まれたYAMLブロック）の有無判定・抽出・読み取り・追記・編集・削除などを行う機能を提供する。

---

## 機能一覧

### 1. frontmatterの存在判定
- 入力: Markdownテキスト（string）
- 出力: frontmatterが存在するかどうか（bool）

### 2. frontmatterの抽出
- 入力: Markdownテキスト（string）
- 出力: frontmatter部分のYAMLテキスト（string）、およびfrontmatter以外の本文（string）

### 3. frontmatterのYAML読み取り
- 入力: Markdownテキスト（string）
- 出力: frontmatterをパースしたmap[string]interface{}または構造体

### 4. frontmatterの追記・新規作成
- 入力: Markdownテキスト（string）、追加したいYAMLデータ（mapまたは構造体）
- 出力: frontmatterを追加したMarkdownテキスト（string）

### 5. frontmatterの編集・上書き
- 入力: Markdownテキスト（string）、編集後のYAMLデータ
- 出力: frontmatterを書き換えたMarkdownテキスト（string）

### 6. frontmatterの削除
- 入力: Markdownテキスト（string）
- 出力: frontmatterを除去したMarkdownテキスト（string）

---

## 仕様詳細

- frontmatterはファイル先頭の`---`で始まり、次の`---`または`...`で終わるYAMLブロックとする
- YAMLパースには`gopkg.in/yaml.v3`等を利用
- 本文とfrontmatterの区切りは厳密に行う（先頭以外の`---`は無視）
- 追記・編集時は既存frontmatterがあれば上書き、なければ新規作成
- 削除時はfrontmatter部分のみ除去し、本文はそのまま残す

---

## 例

### 入力例

```markdown
---
title: サンプル
tags: [go, yaml]
---

本文テキスト
```

### 出力例（frontmatter抽出）

- frontmatter:  
  ```yaml
  title: サンプル
  tags: [go, yaml]
  ```
- 本文:  
  ```
  本文テキスト
  ```

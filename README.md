[![Tests](https://github.com/ishida722/krapp-go/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/ishida722/krapp-go/actions/workflows/go.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ishida722/krapp-go?style=flat-square)](https://github.com/ishida722/krapp-go/blob/main/go.mod)


# krapp

krappは、日々のメモやインボックスノートを素早く作成・管理できるGo製CLIツールです。

## 特徴
- デイリーノート（Daily Note）の自動作成
- インボックスノートの作成（タイトル指定可）
- 設定ファイルによる柔軟なカスタマイズ
- お好みのエディタでノートをすぐに編集可能
- シンプルなYAML形式の設定

## インストール

```bash
go install github.com/ishida722/krapp-go/cmd/krapp@HEAD
```

## テストの実施

```sh
go test ./...
```

## 使い方

### 設定ファイル
- ホームディレクトリまたはカレントディレクトリに `.krapp_config.yaml` を作成できます。
- 例:
  ```yaml
  base_dir: "./notes"
  daily_note_dir: "daily"
  inbox_dir: "inbox"
  editor: "nvim"
  ```

### コマンド一覧

- 設定内容の表示
  ```sh
  krapp print-config
  ```

- デイリーノートの作成
  ```sh
  krapp create-daily
  # 省略形
  krapp cd
  # 作成後にエディタで開く
  krapp create-daily --edit
  krapp cd -e
  ```

- インボックスノートの作成
  ```sh
  krapp create-inbox "タイトル"
  # 省略形
  krapp ci "タイトル"
  # 作成後にエディタで開く
  krapp create-inbox "タイトル" --edit
  krapp ci "タイトル" -e
  ```

- バージョン表示
  ```sh
  krapp --version
  krapp version
  ```

## ディレクトリ構成例

```
notes/
  daily/
    2025/
      06/
        2025-06-03.md
  inbox/
    2025-06-03-タイトル.md
```

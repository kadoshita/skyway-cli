# skyway-cli

[![Test](https://github.com/kadoshita/skyway-cli/actions/workflows/test.yml/badge.svg)](https://github.com/kadoshita/skyway-cli/actions/workflows/test.yml)

SkyWay CLIは、SkyWayを利用したアプリケーションの開発を効率化するためのCLIツールです。

## 注意事項

> [!IMPORTANT]
> このCLIツールは公式として提供している機能ではありません。
> あくまでも個人的に開発したものであり、SkyWay公式のサポートは提供していません。
> 利用される場合は、公式のツールではないことをご理解の上でご利用ください。
> もし不明点や不具合があった場合は、GitHubリポジトリのissueとして連絡いただければ、できる限り対応いたします。

## コマンドリファレンス

[![Deploy documents](https://github.com/kadoshita/skyway-cli/actions/workflows/pages.yml/badge.svg)](https://github.com/kadoshita/skyway-cli/actions/workflows/pages.yml)

- [skyway-cliコマンドリファレンス](./docs/skyway-cli.md)

## 動作環境

- Go 1.23以上
  - macOS Sonoma 14.6.1、Go 1.23.2で動作確認をしています

## インストール

1. Go 1.23以上をインストールする
2. skyway-cli をインストールする
```shell
go install github.com/kadoshita/skyway-cli@latest
```
3. `~/go/bin` にパスを通す
```shell
export PATH=$PATH:$HOME/go/bin
```
4. 実行できることを確認する
```shell
skyway-cli --help
```

## 設定ファイル

1. リポジトリ内の `.skyway-cli.sample.yaml` を `~/.skyway-cli.yaml` として配置する
```shell
$ cp .skyway-cli.sample.yaml ~/.skyway-cli.yaml
```
2. SkyWay ConsoleからアプリケーションIDとシークレットキーを取得し、設定ファイル中の `<APP_ID>` と `<SECRET_KEY>` を上書きする

## 設定値の読み込み

- skyway-cliは以下の順で設定値を読み込みます
  1. 設定ファイル
  2. 環境変数
  3. コマンドライン引数
- 例えば、 `channel get` コマンドでappIdを指定する場合は、以下のようになります

```shell
# 設定ファイルのskyway.app_idは `da636fdd-22f1-4721-a43b-8efc0f1707ac`
$ SKYWAY_APP_ID=8d89c0ae-8b95-47f8-b87d-7a0decce6887 skyway-cli channel get --app-id f4d2b0f9-0dba-4abc-bc4b-fb051d66923a
# => appIdとしてf4d2b0f9-0dba-4abc-bc4b-fb051d66923aが使われる
```

## ドキュメントの自動生成

```shell
SKYWAY_CLI_GEN_DOCS=true go run main.go
```

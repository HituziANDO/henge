# henge

[English README](../README.md)

![henge: universal CLI data transformation tool](../readme-image/henge-image.png)

[CyberChef](https://gchq.github.io/CyberChef/) にインスパイアされた、汎用 CLI データ変換ツールです。

エンコード、デコード、ハッシュ、フォーマット、変換、タイムスタンプ変換をターミナルから離れることなく実行できます。`base64` のフラグや `shasum` の構文をもう検索する必要はありません。

```bash
echo "aGVsbG8=" | henge                 # 自動検出: "Base64っぽい" → hello
echo "hello" | henge encode base64      # → aGVsbG8=
echo '{"z":3,"a":1}' | henge format json # → 整形・ソートされた JSON
cat config.yaml | henge convert json    # → YAML から JSON へ変換
echo "hello" | henge hash sha256        # → 2cf24dba5fb0a30e...
```

## なぜ henge？

| 課題 | 従来の方法 | henge なら |
|------|-----------|-----------|
| Base64 デコード | `base64 -d` か `base64 -D`？OS による | `henge decode base64` |
| SHA256 ハッシュ | `shasum -a 256` か `sha256sum`？ | `henge hash sha256` |
| JSON → YAML 変換 | `yq` をインストールして構文を覚える | `henge convert yaml` |
| "このデータは何？" | 複数のツールを手動で試す | `echo "data" \| henge` |
| "この時刻は何？" | `date -d @1735689600`？OS 依存 | `echo "1735689600" \| henge` |
| 画像を Base64 に | `base64 < image.png` + 手動で Data URI 作成 | `henge encode image --file logo.png --data-uri` |
| UNIX タイムスタンプ | `date -d @1735689600` か `date -r`？OS 依存 | `henge time date 1735689600` |

ひとつのツール。ひとつの構文。パイプで動く。

## インストール

### ソースから（Go 1.22+ が必要）

```bash
go install github.com/HituziANDO/henge@latest
```

### リポジトリからビルド

```bash
git clone https://github.com/HituziANDO/henge.git
cd henge
go build -o henge .
```

### ビルド済みバイナリ（goreleaser）

```sh
goreleaser build --snapshot --clean
```

Linux、macOS、Windows（amd64/arm64）向けのビルド済みバイナリは [Releases](https://github.com/HituziANDO/henge/releases) ページから入手できます。

## 使い方

### 自動検出（デフォルト）

引数なしでデータを `henge` にパイプするだけです。フォーマットを自動検出し、最適な変換を適用します:

```bash
echo "1735689600" | henge           # UNIX タイムスタンプ → 2025-01-01T00:00:00Z
echo "aGVsbG8=" | henge              # Base64 → hello
echo '{"a":1}' | henge              # JSON → 整形表示
echo "hello%20world" | henge        # URL エンコード → hello world
echo "68656c6c6f" | henge           # Hex → hello
echo "name: henge" | henge          # YAML → JSON
```

検出優先順位: UNIX タイムスタンプ → JSON → Base64 → YAML → URL エンコーディング → Hex

### エンコード / デコード

```bash
# Base64
echo "hello" | henge encode base64        # → aGVsbG8=
echo "aGVsbG8=" | henge decode base64     # → hello

# URL エンコーディング
echo "hello world" | henge encode url     # → hello+world
echo "hello%20world" | henge decode url   # → hello world

# Hex
echo "hello" | henge encode hex           # → 68656c6c6f
echo "68656c6c6f" | henge decode hex      # → hello
```

### 画像 Base64

画像ファイルを Base64 文字列にエンコード、またはその逆を行います。HTML/CSS への埋め込みや API ペイロードに便利です。

```bash
# 画像 → Base64
henge encode image --file logo.png               # → iVBORw0KGgo...

# 画像 → Data URI（HTML/CSS 埋め込み用）
henge encode image --file logo.png --data-uri    # → data:image/png;base64,iVBORw0KGgo...

# 76 文字で折り返し（メール用）
henge encode image --file photo.jpg --wrap 76

# Base64 をファイルに保存
henge encode image --file logo.png -o encoded.txt

# Base64 → 画像
henge decode image --file encoded.txt -o restored.png

# Data URI → 画像（自動検出）
echo "data:image/png;base64,iVBORw..." | henge decode image -o output.png

# ラウンドトリップ: エンコードしてデコード
henge encode image --file logo.png | henge decode image -o copy.png
```

対応フォーマット: PNG, JPEG, GIF, WebP, BMP, SVG, ICO

### Time（UNIX タイムスタンプ変換）

```bash
# 自動検出: 入力をそのまま渡すだけ
henge time 1735689600                            # → 2025-01-01T00:00:00Z
henge time "2025-01-01T00:00:00Z"                # → 1735689600
echo "1735689600" | henge time                   # → 2025-01-01T00:00:00Z

# 明示的なサブコマンドで制御
henge time unix "2025-01-01T00:00:00Z"           # → 1735689600
henge time date 1735689600                       # → 2025-01-01T00:00:00Z

# ミリ秒タイムスタンプ（自動検出）
henge time date 1735689600000                    # → 2025-01-01T00:00:00Z
henge time unix --millis "2025-01-01T00:00:00Z"  # → 1735689600000

# タイムゾーン
henge time date --timezone Asia/Tokyo 1735689600
# → 2025-01-01T09:00:00+09:00

# 出力フォーマット（プリセットまたは Go レイアウト）
henge time date --format datetime 1735689600     # → 2025-01-01 00:00:00
henge time date --format "2006/01/02" 1735689600 # → 2025/01/01

```

対応入力フォーマット: RFC3339, RFC1123, RFC822, `2006-01-02 15:04:05`, `2006-01-02`, `2006/01/02`, `2006/01/02 15:04:05`

### ハッシュ

```bash
echo -n "hello" | henge hash md5       # → 5d41402abc4b2a76b9719d911017c592
echo -n "hello" | henge hash sha1      # → aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
echo -n "hello" | henge hash sha256    # → 2cf24dba5fb0a30e26e83b2ac5b9e29e...
echo -n "hello" | henge hash sha512    # → 9b71d224bd62f3785d96d46ad3ea3d73...
```

### フォーマット（整形表示）

```bash
# JSON
echo '{"z":3,"a":1,"b":2}' | henge format json
# {
#   "a": 1,
#   "b": 2,
#   "z": 3
# }

# インデント幅の指定
echo '{"a":1}' | henge format json --indent 4

# コンパクト（最小化）
echo '{ "a" : 1 }' | henge format json -c    # → {"a":1}

# YAML
echo "name: henge" | henge format yaml

# XML
echo '<root><a>b</a></root>' | henge format xml
```

### 変換（フォーマット間変換）

```bash
# JSON → YAML
echo '{"name":"henge","version":"1.0"}' | henge convert yaml
# name: henge
# version: "1.0"

# YAML → JSON
echo "name: henge" | henge convert json
# {
#   "name": "henge"
# }

# JSON → TOML
echo '{"name":"henge","version":"1.0"}' | henge convert toml

# CSV → JSON
echo -e "name,age\nAlice,30\nBob,25" | henge convert json
# [
#   {"age":"30","name":"Alice"},
#   {"age":"25","name":"Bob"}
# ]

# 入力フォーマットを明示的に指定（自動検出を上書き）
cat data.txt | henge convert json --from yaml
```

### ファイル入力

```bash
# --file フラグで指定
henge format json --file data.json
henge hash sha256 --file file.txt

# リダイレクト経由
henge convert yaml < config.json

# パイプ経由
cat data.json | henge format json
```

## グローバルフラグ

| フラグ | 短縮形 | 説明 |
|--------|--------|------|
| `--output <file>` | `-o` | stdout の代わりにファイルに出力 |
| `--file <file>` | | ファイルから入力を読み込み |
| `--from <format>` | `-f` | 入力フォーマットを指定（自動検出を上書き） |
| `--compact` | `-c` | コンパクト出力（インデントなし） |
| `--no-newline` | `-n` | 末尾の改行を付加しない |
| `--help` | `-h` | ヘルプを表示 |
| `--version` | `-v` | バージョンを表示 |

## シェル補完

```bash
# Bash
henge completion bash > /etc/bash_completion.d/henge

# Zsh
henge completion zsh > "${fpath[1]}/_henge"

# Fish
henge completion fish > ~/.config/fish/completions/henge.fish
```

## パイプとの親和性

henge は Unix 哲学に従っています。他のツールと組み合わせて使えます:

```bash
# jq でフィールドを抽出し、YAML に変換
curl -s https://api.example.com/data | jq '.config' | henge convert yaml

# ダウンロードしたファイルのハッシュを計算
curl -sL https://example.com/release.tar.gz | henge hash sha256

# JWT ペイロードをデコード（Base64 URL エンコード）
echo "$JWT" | cut -d. -f2 | henge decode base64

# エンコード/デコードの連鎖
echo "hello" | henge encode base64 | henge decode base64
```

## コマンドリファレンス

```
henge [input]                         自動検出して変換
henge auto [input]                    自動検出（明示的なエイリアス）

henge encode base64 [input]           Base64 エンコード
henge encode url [input]              URL パーセントエンコード
henge encode hex [input]              Hex エンコード
henge encode image --file <file>      画像を Base64 に変換 (--data-uri, --wrap)

henge decode base64 [input]           Base64 デコード
henge decode url [input]              URL パーセントデコード
henge decode hex [input]              Hex デコード
henge decode image [input] -o <file>  Base64 を画像ファイルに変換 (--file でファイル入力)

henge hash md5 [input]                MD5 ハッシュ
henge hash sha1 [input]               SHA-1 ハッシュ
henge hash sha256 [input]             SHA-256 ハッシュ
henge hash sha512 [input]             SHA-512 ハッシュ

henge format json [input]             JSON 整形表示
henge format yaml [input]             YAML 整形表示
henge format xml [input]              XML 整形表示

henge convert json [input]            JSON に変換
henge convert yaml [input]            YAML に変換
henge convert toml [input]            TOML に変換

henge time [input]                    自動検出してタイムスタンプ/日時を変換
henge time unix [input]               日時文字列を UNIX タイムスタンプに変換
henge time date [input]               UNIX タイムスタンプを日時文字列に変換
```

## 開発

```bash
# ビルド
go build -o henge .

# テスト実行
go test ./...

# 詳細出力でテスト実行
go test ./... -v
```

### GoReleaser によるクロスプラットフォームビルド

[GoReleaser](https://goreleaser.com/) を使用して複数プラットフォーム向けのバイナリをビルドします。

**対応プラットフォーム:**

| OS | アーキテクチャ | ターゲット |
|----|---------------|-----------|
| macOS | x86_64 | Intel Mac |
| macOS | arm64 | Apple Silicon |
| Linux | x86_64 | Intel / AMD |
| Linux | arm64 | ARM64 |
| Windows | x86_64 | Intel / AMD |

**スナップショットビルド（ローカル、タグなし）:**

```bash
goreleaser build --snapshot --clean
```

バイナリは `dist/` に出力されます。

**リリースビルド（Git タグと `GITHUB_TOKEN` が必要）:**

```bash
git tag v{X.Y.Z}
git push origin v{X.Y.Z}
```

## ライセンス

MIT

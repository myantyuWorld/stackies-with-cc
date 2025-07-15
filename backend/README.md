# Stackies Backend

Go + Echo を使用したバックエンドAPI

## 開発環境構築

### 前提条件
- Go 1.21+
- Air (ホットリロードツール)

### セットアップ

1. 依存関係のインストール
```bash
go mod download
```

2. Air のインストール (初回のみ)
```bash
go install github.com/air-verse/air@latest
```

## 開発

### ホットリロード開発サーバー起動
```bash
make dev
```
または
```bash
air
```

サーバーは `http://localhost:8080` で起動します。

### その他のコマンド

```bash
# プロダクション用ビルド
make build

# テスト実行
make test

# テストカバレッジ
make test-coverage

# クリーンアップ
make clean

# ヘルプ
make help
```

## API エンドポイント

### ヘルスチェック
- `GET /health` - サーバーの状態確認

### 認証 (実装済み)
- `POST /auth/google/login` - Google OAuth認証
- `POST /auth/refresh` - JWTトークンリフレッシュ
- `POST /auth/logout` - ログアウト
- `GET /auth/me` - ユーザー情報取得

## アーキテクチャ

Clean Architecture を採用：

```
backend/
├── domain/          # ドメイン層
│   ├── model/       # エンティティ
│   ├── repository/  # リポジトリインターフェース
│   └── service/     # ドメインサービスインターフェース
├── usecase/         # ユースケース層
├── infra/           # インフラ層
│   ├── external/    # 外部サービス
│   ├── persistence/ # データ永続化
│   └── dto/         # データ転送オブジェクト
├── presentation/    # プレゼンテーション層
│   ├── handler/     # HTTPハンドラー
│   └── middleware/  # ミドルウェア
└── registry/        # 依存性注入
```

## 開発のガイドライン

- TDD (Test-Driven Development) を実践
- 各層にユニットテストを作成
- Clean Architecture の依存関係ルールを遵守
- コミット前にテストとlintを実行
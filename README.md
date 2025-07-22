# Stackies

Google OAuth2.0 + JWT認証を使用したモダンなWebアプリケーション

## プロジェクト構成

```
stack-frontend/
├── frontend/          # Vue 3 + TypeScript + Vite
├── backend/           # Go + Echo + Clean Architecture
└── docs/             # 設計書・ドキュメント
```

## 技術スタック

### フロントエンド
- Vue 3 (Composition API)
- TypeScript
- Vite
- Pinia (状態管理)
- Vue Router
- Vitest (テスト)

### バックエンド
- Go 1.24
- Echo (Webフレームワーク)
- Clean Architecture
- JWT認証
- Google OAuth2.0

## セットアップ

### 1. 環境変数の設定

プロジェクトルートの `env.sample` ファイルをコピーして `.env` ファイルを作成し、実際の値を設定してください：

```bash
cp env.sample .env
```

#### 必要な環境変数

**フロントエンド**
- `VUE_APP_GOOGLE_CLIENT_ID`: Google Cloud Consoleで取得したクライアントID
- `VUE_APP_GOOGLE_REDIRECT_URI`: OAuthリダイレクトURI
- `VUE_APP_API_BASE_URL`: バックエンドAPIのベースURL

**バックエンド**
- `GOOGLE_CLIENT_ID`: Google Cloud Consoleで取得したクライアントID
- `GOOGLE_CLIENT_SECRET`: Google Cloud Consoleで取得したクライアントシークレット
- `JWT_SECRET_KEY`: JWT署名用のシークレットキー
- `DATABASE_URL`: PostgreSQL接続URL
- `REDIS_URL`: Redis接続URL

### 2. フロントエンドのセットアップ

```bash
cd frontend
npm install
npm run dev
```

### 3. バックエンドのセットアップ

```bash
cd backend
go mod download
go run main.go
```

## 開発

### テストの実行

```bash
# フロントエンド
cd frontend
npm run test:unit

# バックエンド
cd backend
go test ./...
```

### 型チェック

```bash
# フロントエンド
cd frontend
npx tsc --noEmit
```

## 環境変数の詳細

詳細な環境変数の説明は以下のファイルを参照してください：

- `frontend/env.sample` - フロントエンド専用の環境変数
- `backend/env.sample` - バックエンド専用の環境変数
- `env.sample` - プロジェクト全体の環境変数

## 認証フロー

1. ユーザーがGoogleログインボタンをクリック
2. Google OAuth2.0認証画面にリダイレクト
3. ユーザーが認証を実行
4. 認証コード付きでコールバックURLにリダイレクト
5. フロントエンドが認証コードをバックエンドに送信
6. バックエンドがGoogleからアクセストークンを取得
7. ユーザー情報を取得・保存
8. JWTトークンを生成してフロントエンドに返却
9. フロントエンドがトークンをローカルストレージに保存

## ライセンス

MIT License

# 認証システム設計書

## 概要
Google OAuth2.0 + JWT認証を使用した認証システムの包括的設計書。フロントエンド・バックエンド間の認証フロー、アーキテクチャ、実装詳細を定義する。

## 認証システム全体アーキテクチャ

```mermaid
graph TB
    subgraph "Frontend (Vue 3 + Traditional Structure)"
        A[Vue App]
        B[GoogleAuth Component]
        C[useUser Composable]
        D[authService]
    end
    
    subgraph "Backend (Go + MVC Structure)"
        F[Auth Handler]
        G[Auth Middleware]
        H[Google Service]
        I[User Repository]
    end
    
    subgraph "External Services"
        L[Google OAuth2.0]
        M[User Database]
    end
    
    A --> B
    B --> C
    C --> D
    D --> F
    F --> G
    F --> H
    F --> I
    H --> L
    I --> M
```

## 認証フロー設計

### 1. ログインフロー

```mermaid
sequenceDiagram
    participant U as User
    participant F as Frontend
    participant B as Backend
    participant G as Google OAuth
    participant DB as Database
    
    U->>F: ログインボタンクリック
    F->>G: Google認証URLリダイレクト
    G->>U: Google認証画面
    U->>G: 認証実行
    G->>F: 認証コード付きリダイレクト
    F->>B: 認証コード送信 (/auth/google/callback)
    B->>G: アクセストークン取得
    G->>B: ユーザー情報取得
    B->>DB: ユーザー情報保存/更新
    B->>B: JWT生成
    B->>F: JWT + ユーザー情報
    F->>F: セッションストレージ保存
    F->>U: ログイン成功
```

### 2. 認証状態確認フロー

```mermaid
sequenceDiagram
    participant F as Frontend
    participant B as Backend
    participant DB as Database
    
    F->>F: JWT取得
    F->>B: API リクエスト + JWT (/api/user/profile)
    B->>B: JWT検証 (Middleware)
    alt JWT有効
        B->>DB: ユーザー情報取得
        DB->>B: ユーザー情報
        B->>F: 認証成功レスポンス
    else JWT無効
        B->>F: 401 Unauthorized
        F->>F: ログイン画面表示
    end
```

### 3. トークンリフレッシュフロー

**現在未実装** - 現在の実装ではリフレッシュトークン機能は含まれていません。JWTの有効期限切れ時は再ログインが必要です。

```mermaid
sequenceDiagram
    participant F as Frontend
    participant B as Backend
    
    F->>F: アクセストークン期限切れ検出
    F->>F: ログイン画面表示
    Note over F,B: 現在はリフレッシュトークン未実装のため<br/>再ログインが必要
```

## バックエンド認証設計

### 1. 実装アーキテクチャ構造 (MVC Pattern)

```mermaid
graph LR
    subgraph "server/handlers"
        A[google_auth.go]
    end
    
    subgraph "server/middleware"
        B[auth.go]
    end
    
    subgraph "server/services"
        C[google_service.go]
    end
    
    subgraph "server/database"
        D[users.go]
        E[database.go]
    end
    
    subgraph "server/models"
        F[user.go]
    end
    
    A --> C
    A --> D
    B --> D
    C --> F
    D --> F
```

### 2. JWT認証設計

```mermaid
stateDiagram-v2
    [*] --> Unauthenticated
    Unauthenticated --> Authenticating: Google OAuth
    Authenticating --> Authenticated: JWT生成成功
    Authenticated --> TokenExpired: JWT期限切れ
    TokenExpired --> Refreshing: リフレッシュトークン使用
    Refreshing --> Authenticated: リフレッシュ成功
    Refreshing --> Unauthenticated: リフレッシュ失敗
    Authenticated --> Unauthenticated: ログアウト
```

### 3. セキュリティ設計

```mermaid
graph TD
    A[リクエスト] --> B{認証ヘッダー存在?}
    B -->|Yes| C[JWT検証]
    B -->|No| D[401 Unauthorized]
    C --> E{JWT有効?}
    E -->|Yes| F[ユーザー情報取得]
    E -->|No| G[401 Unauthorized]
    F --> H[認証成功]
    G --> I[ログイン画面リダイレクト]
```

## フロントエンド認証設計

### 1. 実装フロントエンド構造 (Traditional Vue.js)

```mermaid
graph TB
    subgraph "src/components"
        A[GoogleAuth.vue]
    end
    
    subgraph "src/composables"
        B[useUser.ts]
    end
    
    subgraph "src/services"
        C[authService.ts]
    end
    
    subgraph "src/types"
        D[user.ts]
    end
    
    subgraph "src/router"
        E[index.ts]
    end
    
    A --> B
    B --> C
    B --> D
    A --> E
    C --> D
```

### 2. 実装コンポーネント詳細

```mermaid
graph LR
    subgraph "実装済みファイル"
        A[GoogleAuth.vue]
        B[useUser.ts]
        C[authService.ts]
        D[user.ts]
    end
    
    subgraph "機能"
        E[Google認証]
        F[ユーザー状態管理]
        G[認証API呼び出し]
        H[型定義]
    end
    
    A --> E
    B --> F
    C --> G
    D --> H
```

### 3. 認証状態管理

```mermaid
stateDiagram-v2
    [*] --> Loading: アプリ起動
    Loading --> Authenticated: 有効なJWT
    Loading --> Unauthenticated: 無効なJWT
    Authenticated --> Refreshing: トークン期限切れ
    Refreshing --> Authenticated: リフレッシュ成功
    Refreshing --> Unauthenticated: リフレッシュ失敗
    Unauthenticated --> Authenticating: ログイン開始
    Authenticating --> Authenticated: ログイン成功
    Authenticating --> Unauthenticated: ログイン失敗
```

### 4. 認証ガード設計

```mermaid
graph TD
    A[ルートアクセス] --> B{認証必要?}
    B -->|Yes| C{認証済み?}
    B -->|No| D[アクセス許可]
    C -->|Yes| D
    C -->|No| E[ログイン画面リダイレクト]
    E --> F[ログイン実行]
    F --> G{ログイン成功?}
    G -->|Yes| H[元のページリダイレクト]
    G -->|No| I[エラー表示]
```

## データモデル設計

### 1. バックエンドデータモデル

```mermaid
erDiagram
    USER {
        string id PK
        string email
        string name
        string picture
        timestamp created_at
        timestamp updated_at
    }
    
    AUTH_TOKEN {
        string user_id FK
        string access_token
        string refresh_token
        timestamp expires_at
        timestamp created_at
    }
    
    USER ||--o{ AUTH_TOKEN : has
```

### 2. フロントエンドデータモデル

```mermaid
graph LR
    A[User Interface] --> B[User Model]
    C[Auth State] --> D[Auth Store]
    E[API Response] --> F[Type Definitions]
    
    subgraph "User Model"
        G[id: string]
        H[email: string]
        I[name: string]
        J[picture: string]
    end
    
    subgraph "Auth State"
        K[user: User | null]
        L[accessToken: string | null]
        M[refreshToken: string | null]
        N[isLoading: boolean]
        O[error: string | null]
    end
```

### 3. JWTペイロード設計

```mermaid
graph LR
    A[JWT Header] --> B[JWT Payload]
    B --> C[JWT Signature]
    
    subgraph "JWT Payload"
        D[user_id]
        E[email]
        F[exp]
        G[iat]
        H[iss]
    end
```

## API設計

### 1. 実装済み認証エンドポイント

```mermaid
graph TD
    A[POST /auth/google/login] --> B[Google OAuth URL生成]
    C[GET /auth/google/callback] --> D[認証コード処理]
    E[GET /api/user/profile] --> F[ユーザー情報取得]
    
    B --> G[リダイレクトURL]
    D --> H[JWT + ユーザー情報]
    F --> I[認証済みユーザー情報]
    
    Note1[未実装: リフレッシュエンドポイント]
    Note2[未実装: ログアウトエンドポイント]
```

### 2. レスポンス設計

```mermaid
graph LR
    A[成功レスポンス] --> B[200 OK]
    C[認証エラー] --> D[401 Unauthorized]
    E[権限エラー] --> F[403 Forbidden]
    G[バリデーションエラー] --> H[400 Bad Request]
    I[サーバーエラー] --> J[500 Internal Server Error]
```

## セキュリティ要件設計

### 1. 実装済みセキュリティポリシー

```mermaid
graph TD
    A[セキュリティ要件] --> B[JWT有効期限]
    A --> C[HTTPS推奨]
    A --> D[JWT検証]
    A --> E[Google OAuth2.0]
    
    B --> F[現在の設定を確認要]
    D --> G[ミドルウェアで実装済み]
    E --> H[Google認証フロー実装済み]
    
    Note1[未実装: CSRF対策]
    Note2[未実装: Redis使用]
    Note3[未実装: リフレッシュトークン]
```

### 2. エラーハンドリング設計

```mermaid
flowchart TD
    A[エラー発生] --> B{エラータイプ}
    B -->|認証エラー| C[401 Unauthorized]
    B -->|権限エラー| D[403 Forbidden]
    B -->|バリデーションエラー| E[400 Bad Request]
    B -->|ネットワークエラー| F[500 Internal Server Error]
    
    C --> G[ログイン画面表示]
    D --> H[権限不足メッセージ]
    E --> I[バリデーションエラーメッセージ]
    F --> J[ネットワークエラーメッセージ]
```

## パフォーマンス設計

### 1. 実装済みキャッシュ戦略

```mermaid
graph LR
    A[認証状態] --> B[セッションストレージ]
    C[ユーザー情報] --> D[Vue Composable State]
    
    Note1[未実装: Redisキャッシュ]
    Note2[未実装: メモリキャッシュ]
```

### 2. 最適化戦略

```mermaid
flowchart TD
    A[API リクエスト] --> B{JWTキャッシュ有効?}
    B -->|Yes| C[キャッシュ使用]
    B -->|No| D[JWT検証実行]
    D --> E[キャッシュ更新]
    C --> F[レスポンス返却]
    E --> F
```

## 監視・ログ設計

### 1. 認証ログ設計

```mermaid
graph TD
    A[認証イベント] --> B[ログイン]
    A --> C[ログアウト]
    A --> D[トークンリフレッシュ]
    A --> E[認証失敗]
    A --> F[セキュリティイベント]
    
    B --> G[成功ログ]
    C --> H[ログアウトログ]
    D --> I[リフレッシュログ]
    E --> J[失敗ログ]
    F --> K[セキュリティログ]
```

### 2. メトリクス設計

```mermaid
graph LR
    A[認証メトリクス] --> B[ログイン成功率]
    A --> C[トークンリフレッシュ率]
    A --> D[認証失敗率]
    A --> E[平均レスポンス時間]
    A --> F[セキュリティイベント数]
```

## 実装詳細

### 1. 実装済みバックエンド構造

```
server/
├── main.go                    # エントリーポイント
├── handlers/
│   └── google_auth.go         # Google認証ハンドラー
├── middleware/
│   └── auth.go                # JWT認証ミドルウェア
├── services/
│   └── google_service.go      # Google OAuth サービス
├── database/
│   ├── database.go            # DB接続設定
│   └── users.go               # ユーザーリポジトリ
└── models/
    └── user.go                # ユーザーモデル
```

### 2. 実装済みフロントエンド構造

```
frontend/src/
├── components/
│   └── Google/
│       └── GoogleAuth.vue     # Google認証コンポーネント
├── composables/
│   └── useUser.ts             # ユーザー状態管理
├── services/
│   └── authService.ts         # 認証API サービス
├── types/
│   └── user.ts                # ユーザー型定義
├── router/
│   └── index.ts               # ルーティング設定
└── main.ts                    # アプリエントリーポイント
```

## 環境設定

### 1. 実装済みバックエンド環境変数

```bash
# Google OAuth2.0
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret

# JWT
JWT_SECRET=your-jwt-secret-key

# Database
DATABASE_URL=postgresql://user:password@localhost:5432/dbname

# Server
PORT=8080

# 注意: 以下は未実装
# REDIS_URL (Redis未使用)
# JWT_REFRESH_TOKEN_EXPIRY (リフレッシュトークン未実装)
```

### 2. 実装済みフロントエンド環境変数

```bash
# Google OAuth2.0
VITE_GOOGLE_CLIENT_ID=your-google-client-id

# API
VITE_API_BASE_URL=http://localhost:8080

# 注意: 実装を確認してViteベースの環境変数設定を使用
```

## 実装状況と今後の課題

### ✅ 実装済み機能
1. **バックエンド基盤**
   - ✅ JWT認証基盤
   - ✅ Google OAuth2.0統合
   - ✅ 基本的な認証フロー
   - ✅ データベース設計・実装
   - ✅ 認証ミドルウェア実装

2. **フロントエンド基盤**
   - ✅ Google認証コンポーネント
   - ✅ ユーザー状態管理 (Composable)
   - ✅ 認証API サービス
   - ✅ 基本的な認証フロー

### 🔄 今後の実装課題
1. **セキュリティ強化**
   - ❌ トークンリフレッシュ機能
   - ❌ ログアウト機能
   - ❌ CSRF対策
   - ❌ レート制限実装

2. **パフォーマンス向上**
   - ❌ Redisキャッシュ実装
   - ❌ セッション管理最適化

3. **アーキテクチャ改善**
   - ❌ クリーンアーキテクチャへの移行
   - ❌ FSDアーキテクチャ採用

4. **監視・ログ**
   - ❌ 認証ログ実装
   - ❌ セキュリティ監視

## テスト戦略

### 1. バックエンドテスト

```mermaid
graph TD
    A[Backend Tests] --> B[Unit Tests]
    A --> C[Integration Tests]
    A --> D[E2E Tests]
    
    B --> E[Domain Tests]
    B --> F[Usecase Tests]
    B --> G[Infrastructure Tests]
    
    C --> H[API Tests]
    C --> I[Database Tests]
    
    D --> J[Authentication Flow Tests]
```

### 2. フロントエンドテスト

```mermaid
graph TD
    A[Frontend Tests] --> B[Unit Tests]
    A --> C[Component Tests]
    A --> D[Integration Tests]
    A --> E[E2E Tests]
    
    B --> F[Store Tests]
    B --> G[Composable Tests]
    B --> H[Utility Tests]
    
    C --> I[UI Component Tests]
    C --> J[Auth Component Tests]
    
    D --> K[Auth Flow Tests]
    D --> L[API Integration Tests]
```

### 3. テストカバレッジ目標

```mermaid
pie title テストカバレッジ目標
    "Unit Tests" : 85
    "Integration Tests" : 70
    "E2E Tests" : 50
    "Component Tests" : 80
```

## 運用・保守

### 1. 監視体制

```mermaid
graph LR
    A[認証システム監視] --> B[アプリケーション監視]
    A --> C[インフラ監視]
    A --> D[セキュリティ監視]
    
    B --> E[認証成功率]
    B --> F[レスポンス時間]
    B --> G[エラー率]
    
    C --> H[サーバーリソース]
    C --> I[データベース性能]
    C --> J[Redis性能]
    
    D --> K[不正アクセス検知]
    D --> L[セキュリティログ]
    D --> M[脆弱性スキャン]
```

### 2. 障害対応フロー

```mermaid
flowchart TD
    A[障害検知] --> B{障害レベル}
    B -->|Critical| C[緊急対応]
    B -->|High| D[優先対応]
    B -->|Medium| E[通常対応]
    B -->|Low| F[計画対応]
    
    C --> G[システム停止]
    D --> H[機能制限]
    E --> I[性能劣化]
    F --> J[改善対応]
```

この設計書により、Google OAuth2.0 + JWT認証システムの包括的な実装が可能になります。各フェーズで段階的に機能を実装し、セキュリティとパフォーマンスを確保しながら、保守性の高いシステムを構築できます。 
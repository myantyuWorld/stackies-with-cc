# 認証システム設計書

## 概要
Google OAuth2.0 + JWT認証を使用した認証システムの包括的設計書。フロントエンド・バックエンド間の認証フロー、アーキテクチャ、実装詳細を定義する。

## 認証システム全体アーキテクチャ

```mermaid
graph TB
    subgraph "Frontend (Vue 3 + FSD)"
        A[Vue App]
        B[Auth Composable]
        C[Auth Store]
        D[LoginScreen]
        E[CallbackScreen]
    end
    
    subgraph "Backend (Go + Clean Architecture)"
        F[Auth Handler]
        G[Auth Middleware]
        H[JWT Service]
        I[Google OAuth Service]
        J[Auth Usecase]
        K[User Repository]
    end
    
    subgraph "External Services"
        L[Google OAuth2.0]
        M[User Database]
        N[Redis Cache]
    end
    
    A --> B
    B --> C
    D --> B
    E --> B
    A --> F
    F --> G
    G --> H
    F --> I
    F --> J
    J --> K
    I --> L
    K --> M
    H --> N
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
    participant R as Redis
    
    U->>F: ログインボタンクリック
    F->>G: Google認証URLリダイレクト
    G->>U: Google認証画面
    U->>G: 認証実行
    G->>F: 認証コード付きリダイレクト
    F->>B: 認証コード送信
    B->>G: アクセストークン取得
    G->>B: ユーザー情報取得
    B->>DB: ユーザー情報保存/更新
    B->>B: JWT生成
    B->>R: トークン保存
    B->>F: JWT + ユーザー情報
    F->>F: ローカルストレージ保存
    F->>U: ログイン成功
```

### 2. 認証状態確認フロー

```mermaid
sequenceDiagram
    participant F as Frontend
    participant B as Backend
    participant DB as Database
    participant R as Redis
    
    F->>F: JWT取得
    F->>B: API リクエスト + JWT
    B->>R: トークン検証
    alt JWT有効
        R->>B: 検証成功
        B->>DB: ユーザー情報取得
        DB->>B: ユーザー情報
        B->>F: 認証成功レスポンス
    else JWT無効
        R->>B: 検証失敗
        B->>F: 401 Unauthorized
        F->>F: ログイン画面表示
    end
```

### 3. トークンリフレッシュフロー

```mermaid
sequenceDiagram
    participant F as Frontend
    participant B as Backend
    participant R as Redis
    
    F->>F: アクセストークン期限切れ検出
    F->>B: リフレッシュトークン送信
    B->>R: リフレッシュトークン検証
    alt 有効
        R->>B: 検証成功
        B->>B: 新しいJWT生成
        B->>R: 新しいトークン保存
        B->>F: 新しいJWT
        F->>F: トークン更新
    else 無効
        R->>B: 検証失敗
        B->>F: 401 Unauthorized
        F->>F: ログイン画面表示
    end
```

## バックエンド認証設計

### 1. クリーンアーキテクチャ構造

```mermaid
graph LR
    subgraph "Presentation Layer"
        A[Auth Handler]
        B[Auth Middleware]
    end
    
    subgraph "Usecase Layer"
        C[Auth Usecase]
    end
    
    subgraph "Domain Layer"
        D[User Model]
        E[Auth Repository]
    end
    
    subgraph "Infrastructure Layer"
        F[JWT Service]
        G[Google OAuth Service]
        H[User Repository]
        I[Redis Repository]
    end
    
    A --> C
    B --> C
    C --> D
    C --> E
    C --> F
    C --> G
    E --> H
    E --> I
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

### 1. FSDアーキテクチャ構造

```mermaid
graph TB
    subgraph "App Layer"
        A[App Router]
        B[App Composables]
        C[App Context]
    end
    
    subgraph "Screens Layer"
        D[LoginScreen]
        E[CallbackScreen]
    end
    
    subgraph "Features Layer"
        F[Auth Feature]
        G[Auth Composable]
        H[Auth Store]
        I[Auth API]
    end
    
    subgraph "Shared Layer"
        J[Shared UI]
        K[Shared API]
        L[Shared Constants]
    end
    
    D --> F
    E --> F
    F --> G
    G --> H
    G --> I
    F --> J
    F --> K
    A --> F
    B --> G
```

### 2. Features Layer 詳細設計

```mermaid
graph LR
    subgraph "features/auth/"
        A[api/authAPI.ts]
        B[composables/useAuth.ts]
        C[lib/googleAuth.ts]
        D[lib/errorHandler.ts]
        E[model/authStore.ts]
        F[ui/LoginForm.vue]
        G[ui/UserProfile.vue]
    end
    
    B --> A
    B --> C
    B --> D
    B --> E
    F --> B
    G --> B
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

### 1. 認証エンドポイント

```mermaid
graph TD
    A[POST /api/auth/google/login] --> B[Google OAuth認証]
    C[POST /api/auth/refresh] --> D[JWTリフレッシュ]
    E[POST /api/auth/logout] --> F[ログアウト]
    G[GET /api/auth/me] --> H[ユーザー情報取得]
    
    B --> I[認証成功レスポンス]
    D --> J[新しいJWT]
    F --> K[ログアウト成功]
    H --> L[ユーザー情報]
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

### 1. セキュリティポリシー

```mermaid
graph TD
    A[セキュリティ要件] --> B[JWT有効期限]
    A --> C[HTTPS必須]
    A --> D[CSRF対策]
    A --> E[XSS対策]
    A --> F[SQLインジェクション対策]
    A --> G[Redisセキュリティ]
    
    B --> H[アクセストークン: 1時間]
    B --> I[リフレッシュトークン: 7日]
    G --> J[Redis認証]
    G --> K[Redis暗号化]
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

### 1. キャッシュ戦略

```mermaid
graph LR
    A[ユーザー情報] --> B[メモリキャッシュ]
    C[JWT検証結果] --> D[Redisキャッシュ]
    E[Google OAuth情報] --> F[セッションキャッシュ]
    G[認証状態] --> H[ローカルストレージ]
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

### 1. バックエンド実装構造

```
backend/
├── cmd/
│   └── main.go
├── core/
│   ├── errors/
│   └── logger/
├── domain/
│   ├── model/
│   │   ├── user.go
│   │   └── auth.go
│   └── repository/
│       ├── user_repository.go
│       └── auth_repository.go
├── usecase/
│   └── auth_usecase.go
├── presentation/
│   ├── handler/
│   │   └── auth_handler.go
│   └── middleware/
│       └── auth_middleware.go
├── infra/
│   ├── external/
│   │   ├── google_service.go
│   │   └── jwt_service.go
│   ├── dto/
│   │   └── user_dto.go
│   ├── persistence/
│   │   ├── user_repository.go
│   │   └── auth_repository.go
│   └── db/
│       └── database.go
└── registry/
    └── registry.go
```

### 2. フロントエンド実装構造

```
frontend/
├── app/
│   ├── composables/
│   │   └── useAppAuth.ts
│   ├── routing/
│   │   └── authRoutes.ts
│   └── context/
│       └── AuthContext.tsx
├── screens/
│   ├── LoginScreen/
│   │   └── LoginScreen.vue
│   └── CallbackScreen/
│       └── CallbackScreen.vue
├── features/
│   └── auth/
│       ├── api/
│       │   └── authAPI.ts
│       ├── composables/
│       │   └── useAuth.ts
│       ├── lib/
│       │   ├── googleAuth.ts
│       │   └── errorHandler.ts
│       ├── model/
│       │   └── authStore.ts
│       └── ui/
│           ├── LoginForm.vue
│           └── UserProfile.vue
└── shared/
    ├── api/
    │   └── auth.ts
    ├── auth/
    │   └── authService.ts
    ├── constants/
    │   └── auth.ts
    ├── ui/
    │   ├── Button/
    │   │   └── GoogleLoginButton.vue
    │   ├── LoadingSpinner/
    │   │   └── LoadingSpinner.vue
    │   └── ErrorMessage/
    │       └── ErrorMessage.vue
    └── lib/
        └── utils.ts
```

## 環境設定

### 1. バックエンド環境変数

```bash
# Google OAuth2.0
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret

# JWT
JWT_SECRET_KEY=your-jwt-secret-key
JWT_ACCESS_TOKEN_EXPIRY=1h
JWT_REFRESH_TOKEN_EXPIRY=168h

# Database
DATABASE_URL=postgresql://user:password@localhost:5432/dbname

# Redis
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=your-redis-password

# Server
SERVER_PORT=8080
SERVER_HOST=localhost
```

### 2. フロントエンド環境変数

```bash
# Google OAuth2.0
VUE_APP_GOOGLE_CLIENT_ID=your-google-client-id
VUE_APP_GOOGLE_REDIRECT_URI=http://localhost:3000/auth/callback

# API
VUE_APP_API_BASE_URL=http://localhost:8080

# App
VUE_APP_APP_NAME=Stackies
VUE_APP_APP_VERSION=1.0.0
```

## 実装優先順位

### 1. Phase 1: 基盤実装 (Week 1-2)
1. **バックエンド基盤**
   - JWT認証基盤
   - Google OAuth2.0統合
   - 基本的な認証フロー
   - データベース設計・実装

2. **フロントエンド基盤**
   - Shared Layer (API Client, Constants)
   - Features Layer (Store, API, Composable)
   - Basic UI Components

### 2. Phase 2: 機能実装 (Week 3-4)
1. **バックエンド機能**
   - トークンリフレッシュ機能
   - 認証ミドルウェア実装
   - エラーハンドリング

2. **フロントエンド機能**
   - LoginScreen
   - CallbackScreen
   - Auth Flow Integration

### 3. Phase 3: セキュリティ強化 (Week 5-6)
1. **セキュリティ実装**
   - セキュリティヘッダー実装
   - レート制限実装
   - Redisセキュリティ設定

2. **監視・ログ実装**
   - 認証ログ実装
   - メトリクス収集
   - セキュリティ監視

### 4. Phase 4: 最適化・デプロイ (Week 7-8)
1. **パフォーマンス最適化**
   - キャッシュ戦略実装
   - レスポンス最適化
   - フロントエンド最適化

2. **デプロイメント**
   - 環境設定
   - CI/CD実装
   - 本番デプロイ

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
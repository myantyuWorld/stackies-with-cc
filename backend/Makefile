.PHONY: dev build test clean help

# デフォルトターゲット
help:
	@echo "Available commands:"
	@echo "  dev     - ホットリロードで開発サーバーを起動"
	@echo "  build   - プロダクション用ビルド"
	@echo "  test    - テスト実行"
	@echo "  clean   - ビルド成果物を削除"
	@echo "  help    - このヘルプを表示"

# ホットリロード開発サーバー
dev:
	@echo "🚀 ホットリロード開発サーバーを起動..."
	@/Users/ohbay/go/bin/air

# プロダクション用ビルド
build:
	@echo "🔨 プロダクション用ビルド..."
	@go build -o ./bin/stackies-backend .
	@echo "✅ ビルド完了: ./bin/stackies-backend"

# テスト実行
test:
	@echo "🧪 テスト実行..."
	@go test ./... -v

# テストカバレッジ
test-coverage:
	@echo "📊 テストカバレッジを計測..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ カバレッジレポート: coverage.html"

# ビルド成果物を削除
clean:
	@echo "🧹 クリーンアップ..."
	@rm -rf ./tmp
	@rm -rf ./bin
	@rm -f build-errors.log
	@rm -f coverage.out
	@rm -f coverage.html
	@echo "✅ クリーンアップ完了"

# binディレクトリ作成
bin:
	@mkdir -p bin
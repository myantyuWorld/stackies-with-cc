# Backend Architecture Design Specification

```xml
<architecture>
  <layer name="cmd">
    <description>エントリーポイント。DIや環境変数を取得し、アプリケーションを起動する。</description>
    <dependencies>
      <dependsOn>registry</dependsOn>
      <dependsOn>core</dependsOn>
    </dependencies>
  </layer>

  <layer name="core">
    <description>アプリケーション全体で使用する処理。アプリケーションエラーやロガーなどを実装する。</description>
    <dependencies />
  </layer>

  <layer name="domain">
    <component name="model">
      <description>エンティティ・値オブジェクトを定義する</description>
    </component>
    <component name="repository">
      <description>リポジトリのinterfaceをおく</description>
    </component>
    <dependencies>
      <dependsOn>core</dependsOn>
    </dependencies>
  </layer>

  <layer name="presentation">
    <component name="handler">
      <description>Handlerメソッドを実装する。usecaseを引数として受け取り、リクエストを処理した後usecaseを呼び出す。リクエスト・レスポンスのstructもここに定義する。</description>
    </component>
    <component name="middleware">
      <description>認証やエラーミドルウェアを実装する。</description>
    </component>
    <component name="server">
      <description>APIサーバーの実装。handlerのマッピングやサーバーの起動処理を実装する。</description>
    </component>
    <dependencies>
      <dependsOn>usecase</dependsOn>
      <dependsOn>core</dependsOn>
    </dependencies>
  </layer>

  <layer name="infra">
    <component name="dto">
      <description>DB用のstructを定義する。ドメインモデルからDTOへの変換処理（逆も然り）もここに実装する</description>
    </component>
    <component name="external">
      <description>外部packageの抽象化を行う（redisやfirebaseなど）</description>
    </component>
    <component name="persistence">
      <description>domain/repositoryの実装</description>
    </component>
    <component name="db">
      <description>DBの接続およびマイグレーション処理の実装</description>
    </component>
    <dependencies>
      <dependsOn>domain</dependsOn>
      <dependsOn>core</dependsOn>
    </dependencies>
  </layer>

  <layer name="usecase">
    <description>
      usecaseの実装を行う。interface、input・outputのstructもここに定義する。
      repositoryのinterfaceおよびドメインモデルを呼びユースケースを組み立てる。
    </description>
    <dependencies>
      <dependsOn>domain</dependsOn>
    </dependencies>
  </layer>

  <layer name="registry">
    <description>repositoryおよびusecaseのDIを行う</description>
    <dependencies>
      <dependsOn>infra</dependsOn>
      <dependsOn>usecase</dependsOn>
      <dependsOn>domain</dependsOn>
    </dependencies>
  </layer>
</architecture>
```
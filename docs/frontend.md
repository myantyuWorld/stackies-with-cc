# Frontend Architecture Design Specification

```xml
<architecture>
  <metadata>
    <description>このアーキテクチャは、モノレポにおけるフロントエンドアプリの構成を記述したものである。</description>
    <designPrinciple>FSD（Feature-Sliced Design）を参考にしている。</designPrinciple>
    <structureRule>各レイヤーの下には、lib, model, ui, api などのスライスを配置する。</structureRule>
    <convention>全てのモジュール・コンポーネントは必ず named export を行う。</convention>
  </metadata>

  <layer name="shared">
    <description>アプリケーション共通処理やコンポーネントなどを実装するレイヤー</description>
    <component name="api">
      <description>APIクライアントの初期化を行う。APIエラー型や、openapi.ymlから生成されたTSファイルなどを配置する</description>
    </component>
    <component name="auth">
      <description>Firebaseを使用した認証処理やログイン状態の監視を行う処理を実装</description>
    </component>
    <component name="constants">
      <description>アプリケーション全体で使用する色や各種定数を実装</description>
    </component>
    <component name="ui">
      <description>ButtonやModalのような、アプリケーション共通で使用するUIコンポーネントを実装</description>
    </component>
    <component name="lib">
      <description>画像アップロードや日時変換などのユーティリティ処理を実装</description>
    </component>
  </layer>

  <layer name="screens">
    <description>featuresレイヤーの実装を使用し、画面を実装する</description>
    <dependencies>
      <dependsOn>features</dependsOn>
      <dependsOn>shared</dependsOn>
    </dependencies>
  </layer>

  <layer name="features">
    <description>APIリクエストや、UI用のデータ形式への変換などの処理を行う。UIを提供することもある。</description>
    <dependencies>
      <dependsOn>shared</dependsOn>
    </dependencies>
  </layer>

  <layer name="app">
    <description>ルーティングやルートコンポーネントなどを提供する。Contextなどの実装もここにおく。</description>
    <dependencies>
      <dependsOn>shared</dependsOn>
      <dependsOn>features</dependsOn>
      <dependsOn>screens</dependsOn>
    </dependencies>
  </layer>
</architecture>
```
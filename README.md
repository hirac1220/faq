# FAQサービス

## Docker
### 環境構築(Docker)
事前準備:
* nuxt.config.jsのURL/Portをローカルへ変更

開発環境を起動
```
docker-compose up
```

apiサーバをビルド
```
docker exec -it {docker_id} sh

# /go/src/appで
go build
```
apiサーバを起動
```
./api
```
web開発環境へ移動し起動
```
cd {path}/faq/web
npm run dev
```
### 動作確認(Docker)

ローカル環境でのWebサイト  
[FAQ](http://localhost:3000)

参照
```
curl http://localhost:18080/faq
```
追加
```
curl -X POST -d '{"text":"test"}' http://localhost:18080/faq
```

## Google App Engine
### 環境構築(GAE)

事前準備:
* nuxt.config.jsのURL/PortをGAEへ変更

API: 開発環境へ移動しビルド
```
cd {path}/faq/api
go build
```
Web: 開発環境へ移動しビルド
```
cd {path}/faq/web
npm run build
```
GCPへログインし、認証コード貼り付ける(初回のみ)
```
gcloud auth login
```
GAEへデプロイ(web/api各ディレクトリより)
```
gcloud app deploy --project {PROJECT_ID}
```
### 動作確認(GAE)

Webサイト on GAE  
[FAQ](https://{WEB_SERVICE}-dot-{PROJECT_ID}.appspot.com/)

参照
```
curl https://{API_SERVICE}-dot-{PROJECT_ID}.appspot.com/faq
```
追加
```
curl -X POST -d '{"text":"test"}' https://{API_SERVICE}-dot-{PROJECT_ID}.appspot.com/faq
```
# readme
gingonic + クリーンアーキテクチャー
## run
```
$GO111MODULE=on go run api/server.go 
```

## gae
アカウントの切り替え
```
$ gcloud auth login
```

プロジェクトの設定
```
$ gcloud config set project <project_id>
```

設定の確認
```
$ gcloud config list
```

ローカルで確認
```
$ goapp serve
```

デプロイ
```
$ gcloud app deploy --project <project_id>
```

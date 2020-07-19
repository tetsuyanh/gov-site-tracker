# gov-site-tracker

## development

to test table

```shell
go run main.go
```

## deploy

```shell
gcloud app deploy
gcloud app deploy cron.yaml
```

## ops

```shell
bq show --schema --format=prettyjson lustrous-bus-243613:gov_site.tracking > ./bq/tracking.json
bq mk --table lustrous-bus-243613:gov_site.tracking ./bq/tracking.json
```

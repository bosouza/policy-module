#!/bin/bash


exec-db () {
  mysql -h localhost -P 3306 --protocol=TCP -u root -pmypass
}

deploy-db-docker () {
  docker run --name mariadbtest -e MYSQL_ROOT_PASSWORD=mypass -p 3306:3306 -d mariadb:10.3
}

reset-db () {
  mysql -h localhost -P 3306 --protocol=TCP -u root -pmypass < ./db/reset-db.sql
}

deploy-db-resources () {
  mysql -h localhost -P 3306 --protocol=TCP -u root -pmypass < ./db/db.sql
}

app-run () {
  go run ./cmd/policyserver/main.go
}
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
  mysql -h localhost -P 3306 --protocol=TCP -u root -pmypass < ./db/test-data.sql
}

app-run () {
  go run ./cmd/policyserver/main.go
}

create-policy () {
  curl http://localhost:8180/policy -X POST -i \
    -d '{
      "id": "createPolicies",
      "policyResources" : [
        {
          "resourceId" : "policy",
          "content":"{\"activity\":\"write\"}"
        }
      ]
    }'
}

assign-policy () {
  curl http://localhost:8180/assign/test-user-1/createPolicies -X PUT -i
}

evaluate-policy () {
  curl http://localhost:8180/evaluate -X POST -i \
    -d '{
      "policyCheckId": "allowPolicy",
      "input":{
        "user":"test-user-1",
        "activity":"write"
      }
    }'
}
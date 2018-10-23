#!/bin/bash

curl localhost:8080/topics -XPOST -d'{"name":"test_topic", "engine":"AWS"}' -H'content-type: application/json' -i
curl localhost:8080/subscriptions -XPOST -d'{"name":"test_subscriber", "topic":"test_topic", "endpoint":"http://wequeue:8080/test_subscriber"}' -H'content-type: application/json' -i
curl localhost:8080/messages -XPOST -d'{"topic":"test_topic", "payload":"Payload"}' -H'content-type: application/json' -i

docker ps |grep lambda| awk '{system("docker network connect wequeue-management_api_default " $1 " --link wequeue")}'


# Bondi

This is the event bus of Wenance. it's made for transforming the synchronized api calls between different microservices in an asynchronous, and event oriented comunnication.

## For starting the local environment:
```
docker-compose -p bondi up --build
```

##For running integration tests:
```
docker-compose -p bondi -f docker-compose-integration.yml up --build
```

## API Documentation
http://localhost:8080/swagger/index.html

##Environments

###Dev
http://bondi.dev.fintechpeople.io:81
###Stage
http://bondi.stage.fintechpeople.io:81
###Prod
http://bondi.fintechpeople.io:81

## Create a topic
```
curl -XPOST http://localhost:8080/topics {"name":"topic_name", "engine":"AWS"}
```

## Start your push subscriber app dockerized, in the same network
```
docker run --network=bondi_default --network-alias=$SUBSCRIBER_NAME your:app  
```
## Create a push subscriber
```
curl localhost:8080/subscribers -XPOST -d'{"name":"test_subscriber", "topic":"test_topic", "endpoint":"http://$SUBSCRIBER_NAME:$PORT/test_subscriber", "type":"push"}' -H'content-type: application/json'
```

## Create a pull subscriber
```
curl localhost:8080/subscribers -XPOST -d'{"name":"test_subscriber", "topic":"test_topic", "visibility_timeout":30, "type":"pull"}' -H'content-type: application/json'
```

## Publish a message
```
curl localhost:8080/messages -XPOST -d'{"topic":"test_topic", "payload":{"message":"Hello!!"}}' -H'content-type: application/json'
```
## Consume failed pushed messages from the dead-letter-queue if it's a push subscriber
```
curl "localhost:8080/messages?max_messages=10&subscriber=test_subscriber"
```
## Consume messages from the topic if it's a pull subscriber
```
curl "localhost:8080/messages?max_messages=10&subscriber=test_subscriber"
```

## Delete failed pushed messages from the dead-letter-queue if it's a push subscriber
```
curl -XDELETE localhost:8080/messages -d'{"subscriber":"test_subscriber", "messages": [{"message_id":"1", "delete_token":"x"}]}'
```

## Delete messages from the queue if it's a pull subscriber
```
curl -XDELETE localhost:8080/messages -d'{"subscriber":"test_subscriber", "messages": [{"message_id":"1", "delete_token":"x"}]}'
```


# wequeue-management_api

## For starting the local environment:
```
docker-compose --build
```
## API Documentation
http://localhost:8080/swagger/index.html

## Create a topic:
```
curl -XPOST http://localhost:8080/topics {“name”:”topic_name”, “engine”:”AWS”}
```

## Create subscriber
```
curl localhost:8080/subscribers -XPOST -d'{"name":"test_subscriber", "topic":"test_topic", "endpoint":"http://wequeue:8080/test_subscriber"}' -H'content-type: application/json'
```

## Publish a message
```
curl localhost:8080/messages -XPOST -d'{"topic":"test_topic", "payload":{"message":"Hello!!"}}' -H'content-type: application/json'
```
## Consume Failed messages from the dead-letter-queue
```
curl "localhost:8080/messages?max_messages=10&subscribe=test_subscriber"
```

## Delete failed messages from the dead-letter-queue
```
curl -XDELETE "localhost:8080/messages? -d'{"subscriber":"test_subscriber", "messages": [{"message_id":"1", "delete_token":"x"}]}'
```

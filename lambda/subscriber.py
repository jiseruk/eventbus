from __future__ import print_function
import requests
import os
import json


def handler(event, context):
    print('Loading function')
    # Get the queue
    queueName = os.getenv('queue_name', None)
    endpoint = os.getenv('subscriber_url', None)
    topic = os.getenv('topic', None)
    environment = os.getenv('environment', None)
    print("Received event: " + json.dumps(event, indent=2))
    print("Endpoint:" + endpoint)
    message = json.loads(event['Records'][0]['Sns']['Message'])
    payload = {"payload": message, "topic": topic}
    try:
        r = requests.post(endpoint, json=message, headers={"Content-type": "application/json"})
        print("POST message result: ", r.status_code)
        if r.status_code >= 400:
            raise Exception("Fail posting to " + endpoint + " Response: " + str(r.status_code) + "|" + r.content)
        return message
    except Exception, e:
        # No funciona en Localstack DeadLetterQueue, por lo que encolo a mano en la cola los mensajes no procesados luego 
        # de la primera falla. En prod esto lo maneja AWS Lambda, pero luego de 2 reintentos.
        if environment == "local":
            import boto3
            sqs = boto3.resource('sqs', region_name='us-east-1', endpoint_url='http://localhost:4576')
            queue = sqs.get_queue_by_name(QueueName=queueName)
            queue.send_message(MessageBody=json.dumps(payload))
        else:
            raise e
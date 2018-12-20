from __future__ import print_function
import requests
import os
import json

queueName = os.getenv('queue_name', None)
endpoint = os.getenv('subscriber_url', None)
environment = os.getenv('environment', None)
 
def handler(event, context):
    print('Loading function')
    # Get the queue
    print("Received event: " + json.dumps(event, indent=2))
    print("Endpoint:" + endpoint)
    sns = event['Records'][0]['Sns']
    message = json.loads(sns['Message'])
    headers = sns.get('MessageAttributes', {})
    headers = {k:v["StringValue"] for k, v in headers}
    headers["Content-Type"] = "application/json"
    try:
        r = requests.post(endpoint, json=message, headers=headers)
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
            queue.send_message(MessageBody=json.dumps(event))
        else:
            raise e
from __future__ import print_function
import requests
import os
import json
print('Loading function')

def handler(event, context):
    endpoint = os.getenv('subscriber_url', None)
    topic = os.getenv('topic', None)
    print("Received event: " + json.dumps(event, indent=2))
    print("Endpoint:" + endpoint)
    message = json.loads(event['Records'][0]['Sns']['Message'])
    payload = {"message": message, "topic": topic}
    r = requests.post(endpoint, json=payload, headers={"Content-type": "application/json"})
    print("POST message result: ", r.status_code)
    if r.status_code >= 400:
        raise Exception("Fail posting to " + endpoint + " Response: " + r.status_code + "|" + r.content)
    return message
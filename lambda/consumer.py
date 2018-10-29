from __future__ import print_function
import requests
import os
import json
import base64

print('Loading function')

def handler(event, context):
    endpoint = os.getenv('subscriber_url', None)
    topic = os.getenv('topic', None)
    print("Endpoint:" + endpoint)
    print('Topic:' + topic)
    message = json.loads(base64.b64decode(event['Records'][0]['kinesis']['data']))
    payload = {"payload": message, "topic": topic}
    r = requests.post(endpoint, json=payload, headers={"Content-type": "application/json"})
    print("POST message result: ", r.status_code)
    if r.status_code >= 400:
        raise Exception("Fail posting to " + endpoint + " Response: " + str(r.status_code) + "|" + r.content)
    return message
const axios = require('axios');
/* SNS Event Structure
Links:
https://docs.aws.amazon.com/lambda/latest/dg/retries-on-errors.html
https://docs.aws.amazon.com/lambda/latest/dg/invoking-lambda-function.html#supported-event-source-sns

{
  "Records": [
    {
      "EventVersion": "1.0",
      "EventSubscriptionArn": eventsubscriptionarn,
      "EventSource": "aws:sns",
      "Sns": {
        "SignatureVersion": "1",
        "Timestamp": "1970-01-01T00:00:00.000Z",
        "Signature": "EXAMPLE",
        "SigningCertUrl": "EXAMPLE",
        "MessageId": "95df01b4-ee98-5cb9-9903-4c221d41eb5e",
        "Message": "Hello from SNS!",
        "MessageAttributes": {
          "Test": {
            "Type": "String",
            "Value": "TestString"
          },
          "TestBinary": {
            "Type": "Binary",
            "Value": "TestBinary"
          }
        },
        "Type": "Notification",
        "UnsubscribeUrl": "EXAMPLE",
        "TopicArn": topicarn,
        "Subject": "TestInvoke"
      }
    }
  ]
}
 */
console.log('Loading function');

exports.handler = function(event, context, callback) {

    event.Records.forEach(function(record) {
        // Kinesis data is base64 encoded so decode here
        var payload = new Buffer(record.kinesis.data, 'base64').toString('ascii');
        console.log('Decoded payload:', payload);
        axios.post(process.env.subscriber_url, {
            payload: payload,
            topic: process.env.topic
        }).then((res) => {
            console.log(`statusCode: ${res.statusCode}`);
            if(res.statusCode >= 400){
                var error = new Error("something is wrong");
                callback(error);
            } else {
                console.log(res);
                callback(null, "Success");
            }
        }).catch((error) => {
            console.error(error);
            //context.fail();
            callback(error);
        });
    });
    /*
    var message = event.Records[0].Sns.Message;
    console.log('Message received from SNS:', message);
    console.log('Topic:', process.env.topic);
    console.log('Subscriber url:', process.env.subscriber_url);

    axios.post(process.env.subscriber_url, {
        payload: JSON.parse(message),
        topic: process.env.topic
    }).then((res) => {
        console.log(`statusCode: ${res.statusCode}`);
        if(res.statusCode >= 400){
            var error = new Error("something is wrong");
            callback(error);
        } else {
            console.log(res);
            callback(null, "Success");
        }
    }).catch((error) => {
        console.error(error);
        //context.fail();
        callback(error);
    });*/
};
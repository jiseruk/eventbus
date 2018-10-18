const axios = require('axios');

console.log('Loading function');

exports.handler = function(event, context, callback) {
// console.log('Received event:', JSON.stringify(event, null, 4));

    var message = event.Records[0].Sns.Message;
    console.log('Message received from SNS:', message);
    console.log('Topic:', process.env.topic);
    console.log('Subscriber url:', process.env.subscriber_url);
    axios.post(process.env.subscriber_url, {
        payload: message,
        topic: process.env.topic
    }).then((res) => {
        console.log(`statusCode: ${res.statusCode}`);
        console.log(res);
        callback(null, "Success");
    }).catch((error) => {
        console.error(error);
        callback(null, "Failure");
    });
};
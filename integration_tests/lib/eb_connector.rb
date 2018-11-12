require 'http'

class EBConnector

	attr_reader :last_response
	attr_accessor :host

	def initialize host
		self.host = host
	end	

	# Creates a topic to /topics
	# @param [Hash]  request body {"name" => "TopicName", "engine" => "AWS"}
	# @return [HTTP::Response] {"name"=>"TopicName", "engine"=>"AWS", "resource_id"=>"arn:aws:sns:us-east-1:123456789012:TopicTest1541014797_1541014797", "created_at"=>"2018-10-31T19:40:23.8135555Z"}
	def post_topic body
		url = self.host + "/topics"
		post url, body.to_json
	end

	# Subscribe to given topic
	# @param [Hash] request body {"topic" => topic_name, "name" => subscriber, "endpoint" => endpoint}
	# @return [HTTP::Response]
	def subscribe_to_topic body
		url = self.host + "/subscribers"
		post url, body.to_json
	end

	# Send event to a topic
	# @param [Hash] request body {"topic": "topic_name", "payload": {"some":"thing"}}
	def send_event body
		url = self.host + "/messages"
		post url, body.to_json
	end

	# Returns the messages for a given subscriber
	# @param [String] subscriber = the name of the subscriber
	# @param [String] max_messages = the amount of messages per request
 	# @return [JSON] {  "messages": [{"delete_error": {"code": "string","message": "string"},"delete_token": "string","message_id": "string","payload": {}}], "topic": "string"}
	def messages subscriber=nil, max_messages=nil
		url = self.host + "/messages?"
		url += "&subscriber=#{subscriber}" if subscriber
		url += "&max_messages=#{max_messages}" if max_messages
		get url
	end

	private

	def post url, body
		@last_response = HTTP.post url, body: body
	end

	def get url
		@last_response = HTTP.get url
	end

end

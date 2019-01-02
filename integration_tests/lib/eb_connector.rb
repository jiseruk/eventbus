require 'http'

class EBConnector

	attr_reader :last_response
	attr_accessor :host

	def initialize host
		self.host = host
	end	

	def topics_endpoint
		self.host + "/topics"
	end
	
	def subscription_endpoint
		self.host + "/subscribers"
	end

	def messages_endpoint
		self.host + "/messages"
	end

	def list_topics
		url = self.topics_endpoint + "?"
		# url += "&subscriber=#{subscriber}" if subscriber
		# url += "&max_messages=#{max_messages}" if max_messages
		get url
	end

	def list_subscribers_of_topic topic_name
		url = self.topics_endpoint + "/#{topic_name}/subscribers"
		get url
	end

	def subscriber_data_for subscriber
		url = self.subscription_endpoint + "/#{subscriber}"
		get url
	end

	# Creates a topic to /topics
	# @param [Hash]  request body {"name" => "TopicName", "engine" => "AWS"}
	# @return [HTTP::Response] {"name"=>"TopicName", "engine"=>"AWS", "resource_id"=>"arn:aws:sns:us-east-1:123456789012:TopicTest1541014797_1541014797", "created_at"=>"2018-10-31T19:40:23.8135555Z"}
	def post_topic body
		url = self.topics_endpoint
		post url, body.to_json
	end

	# Subscribe to given topic
	# @param [Hash] request body {"topic" => topic_name, "name" => subscriber, "endpoint" => endpoint}
	# @return [HTTP::Response]
	def subscribe_to_topic body
		url = self.subscription_endpoint
		post url, body.to_json
	end

	# Send event to a topic
	# @param [Hash] request body {"topic": "topic_name", "payload": {"some":"thing"}}
	def send_event body, headers={}
		# url = self.host + "/messages"
		url = self.messages_endpoint
		post url, body.to_json, headers
	end

	# Returns the messages for a given subscriber
	# @param [String] subscriber = the name of the subscriber
	# @param [String] max_messages = the amount of messages per request
 	# @return [JSON] {  "messages": [{"delete_error": {"code": "string","message": "string"},"delete_token": "string","message_id": "string","payload": {}}], "topic": "string"}
	def messages subscriber=nil, max_messages=nil
		# url = self.host + "/messages?"
		url = self.messages_endpoint + "?"
		url += "&subscriber=#{subscriber}" if subscriber
		url += "&max_messages=#{max_messages}" if max_messages
		get url
	end

	def delete_messages body
		url = self.messages_endpoint
		delete url, body.to_json
	end

	private


	def get url
		@last_response = HTTP.get url
	end

	def post url, body, headers={}
		@last_response = HTTP.post url, body: body, headers: headers
	end

	def delete url, body
		@last_response = HTTP.delete url, body: body
	end

end

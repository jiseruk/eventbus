module Topic
	module Methods

		def response
			@response
		end

		def parsed_response
			response.parse
		end

		def post_topic(opts)
			name = opts.delete(:topic_name)
			engine = opts.delete(:engine) || "AWS"
			body = {"name" => name, "engine" => engine}
			@response = $eb_connector.post_topic body
		end

		def success?
			response.status.success?
		end

		def has_topic_id?
			!!parsed_response['resource_id']
		end

		def has_creation_date?
			!!parsed_response['created_at']
		end

		def has_topic_name? name
			parsed_response["name"] == name
		end

		def has_security_token?
			byebug
			!!parsed_response
			true
		end

		def create_topic
			name = random_topic_name
			post_topic topic_name: name
			name
		end

		def create_topics(size)
			@topics = []
			size.times{@topics << create_topic}
		end

		def topics
			@topics
		end

		def random_topic_name
			"TestTopicName-#{timestamp}"
		end

		def random_subscriber_name
			"SubscriberTest#{timestamp}"
		end

		def random_fake_endpoint
			"http://subscriber:9292/returnok/#{random_word}#{timestamp}"
		end

		def unexisting_endpoint
			"http://subscriber:9292/unknown/#{timestamp}"
		end

		def response_message
			parsed_response["message"]
		end

		def status_code
			response.status.code
		end

		def subscribe_to_topic(opts)
			body={}
			body['topic'] 		= opts.delete(:topic_name)
			body['name'] 		= opts.delete(:subscriber) 
			body['endpoint'] 	= opts.delete(:endpoint)
			body['type'] 		= opts.delete(:type)
			body['visibility_timeout'] = opts.delete(:visibility_timeout) if body['type'] == 'pull'
			@response = $eb_connector.subscribe_to_topic body
			parsed_response
		end

		def create_a_push_subscription timeout = 60
			subscriber = @subscriber || random_subscriber_name
			endpoint = @endpoint || random_fake_endpoint
			res = subscribe_to_topic(topic_name: @topic_name, subscriber: subscriber, type: "push", endpoint: endpoint)
			puts res if $debug
		end

		def create_a_pull_subscription timeout = 60
			subscriber = @subscriber || random_subscriber_name
			@timeout = timeout || 30
			res = subscribe_to_topic(topic_name: @topic_name, subscriber: subscriber, type:"pull", visibility_timeout: @timeout)
			puts res if $debug
		end

		def create_a_subscription_of_type type, timeout=nil
			type == 'push' ? create_a_push_subscription(timeout) : create_a_pull_subscription(timeout)
		end

		def create_subscriptions size
			@subscribers = []
			size.times do |i|
				@subscribers << (@subscriber = "subscriber_#{timestamp}")
				puts "Subscriber name: #{@subscriber}" if $debug
				create_a_pull_subscription 5
			end
		end

		def first_time_endpoint
			$configuration.first_time_endpoint
		end

		def subscriber_endpoint
			$configuration.subscriber_endpoint
		end

		def subscribed?
			response.status.success?
			true
		end

		def random_event_for_topic(name=nil, timeout=5)
			payload = {"key" => "value", "sent_by" => "Automated test on #{now.strftime('%Y-%m-%m %H:%M:%S')}"}
			puts "#{payload}"if $debug
			@sent_event = {"topic" => name, "visibility_timeout" => timeout, "payload" => payload}
		end

		def sent_event
			@sent_event
		end

		def send_event(body)
			@response = $eb_connector.send_event(body)
			puts "Sending result: #{parsed_response}"if $debug
		end

		def last_event
			$subscriber_connector.last_event
		end

		def event_transmitted?
  			if last_event["topic"] != sent_event["topic"]
				puts "#{last_event['topic']} != #{sent_event['topic']}"if $debug
				return false
			end
			if last_event["payload"] != sent_event["payload"]
				puts "#{last_event['payload']} != #{sent_event['payload']}"if $debug
				return false
			end
			true
		end

		def ask_for_missing_events(subscriber=nil, max_msg=nil)
			@response = $eb_connector.messages(subscriber,max_msg)
			@response
		end

		def exists_events? sent_event
			messages.select do |msg|
				are_messages_equals?(msg, sent_event)
			end.any?
		end

		def message_has_sent_date?
			!!last_message["message"]["timestamp"]
		end

		def are_messages_equals?(msg1, msg2)
			msg1['message']["topic"]== msg2["topic"] and
			msg1['message']["payload"] == msg2["payload"]
		end

		def messages
			@response.parse['messages']
		end

		def last_message
			messages&.last
		end
		
		def delete_messages subscriber_name=nil, messages_list=nil
			subscriber = subscriber_name || @subscriber
			msgs = messages_list || messages
			# messages is like [{"message_id" => "value", "delete_token" => "value"}]
			body = {"subscriber" => subscriber, "messages" => msgs}
			$eb_connector.delete_messages(body)
		end

		def timestamp
			"#{now.to_f}".gsub(".","")
		end

		def now
			Time.now
		end

		def random_word
			Faker::Lorem.word
		end

	end
end

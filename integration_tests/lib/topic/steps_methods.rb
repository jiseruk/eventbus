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

		def has_topic_name? topic_name
			parsed_response["name"] == topic_name
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
			@response = $eb_connector.subscribe_to_topic body
		end

		def create_a_push_subscription
			@subscriber = random_subscriber_name
			@fake_endpoint = random_fake_endpoint
			subscribe_to_topic(topic_name: @topic_name, subscriber: @subscriber, type: "push", endpoint: @fake_endpoint)
		end

		def create_a_pull_subscription
			@subscriber = random_subscriber_name
			subscribe_to_topic(topic_name: @topic_name, subscriber: @subscriber, type:"pull")
		end

		def create_a_subscription_of_type type
			type == 'push' ? create_a_push_subscription : create_a_pull_subscription
		end

		def first_time_endpoint
			"http://subscriber:9292/ok_only_first_time/#{timestamp}"
		end

		def subscribed?
			response.status.success?
			true
		end

		def random_event_for_topic(topic_name=nil)
			@sent_event = {"topic" => topic_name, "payload" => {"key" => "value", "sent_by" => "Automated test on #{now.strftime('%Y-%m-%m %H:%M:%S')}"}}
		end

		def sent_event
			@sent_event
		end

		def send_event(body)
			@response = $eb_connector.send_event(body)
		end

		def last_event
			$subscriber_connector.last_event
		end

		def event_transmitted?
  			if last_event["topic"] != sent_event["topic"]
				puts "#{last_event['topic']} != #{sent_event['topic']}"
				return false
			end
			if last_event["payload"] != sent_event["payload"]
				puts "#{last_event['payload']} != #{sent_event['payload']}"
				return false
			end
			true
		end

		def ask_for_missing_events(subscriber=nil, max_msg=nil)
			@response = $eb_connector.messages(subscriber,max_msg)
		end

		def exists_events? sent_event
			messages.select do |msg|
				msg["payload"].keys.sort == sent_event.keys.sort
			end.any?
		end

		def message_has_sent_date?
			!!last_message["created_at"]
		end


		def messages
			@response.parse['messages']
		end

		def last_message
			messages&.last
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

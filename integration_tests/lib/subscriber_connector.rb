class SubscriberConnector
	attr_reader :last_response
	attr_accessor :host

	def initialize env = "localhost"
		@host = (env == "docker") ? "subscriber" : "localhost"
		puts "Using host for subscriber: #{@host}"
		puts "events => #{last_event}"
		self
	end

	def last_event
		events.first
	end

	def events
		url = "http://#{@host}:9292" + "/events"
		puts "Using url: #{url}"
		response = HTTP.get url
		response.parse
	end

end

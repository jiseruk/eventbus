require 'clipar'
require 'inrules'

class Configuration

	attr_reader :cli, :env, :host

	def initialize 
		@cli = CLIPar.new
		@env = normalize_environment(@cli.ENVIRONMENT) 
		@host = Inrules.get_params()["hosts"][@env]
		@application_config = Inrules.get_params()["application"]
		puts "
############################################
#
#    Ambiente: #{@env}
#    Host: #{@host}
#    Services: #{@services}
#
############################################"
	end

	def base_url
		@host
	end

	def normalize_environment which
		case which
		when /^dev/
			"dev"
		when /stg|stage|staging/
			"stage"
		when /prod/
			"prod"
		when /local/
			"localhost"
		when /docker/
			"docker"
		else
			"localhost"
		end
	end

	def subscriber_endpoint
		'http://subscriber:9292/events'
	end

	def first_time_endpoint
		"http://subscriber:9292/ok_only_first_time/#{timestamp}"
	end

	def timestamp
		"#{Time.now.to_f}".gsub(".","")
	end

	def timeout
		@application_config["timeout"]
	end

end	


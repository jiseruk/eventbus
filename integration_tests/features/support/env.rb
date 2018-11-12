require 'http'
require 'require_all'
require 'byebug'
require 'syntax'
require 'timeout'
require 'faker'

require_all "#{Dir.pwd}/lib"


$configuration = Configuration.new

$eb_connector = EBConnector.new $configuration.host

$subscriber_connector = SubscriberConnector.new $configuration.env

$execution_id = "#{Time.now.to_i}"


World(Topic::Methods)

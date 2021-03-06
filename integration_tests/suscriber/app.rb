require 'digest'
require "cuba"
require 'json'
require "logger"

# $logger = Logger.new($stdout)

$events = []

$answered = []

Cuba.define do

  on get do

    on root do
      res.status = 404
      output = {'code'=>404, 'message'=>'Not found'}
      res.headers["Content-Type"] = "application/json; charset=utf-8"
      puts "#{output}"
      res.write output.to_json
    end

    on "events" do
      res.headers["Content-Type"] = "application/json; charset=utf-8"
      output = $events.reverse.to_json
      puts "#{output}"
      res.write output
    end

    on true do
      res.status = 404
      output = {'code'=>404, 'message'=>'Not found'}
      puts "#{output}"
      res.headers["Content-Type"] = "application/json; charset=utf-8"
      res.write output.to_json
    end
  end

  on post do
    
    on "events/:any" do |any|
      begin
        json = req.body.read
        body = JSON.parse json
        puts "#{body} to /events/#{any}"
        print "Event received: #{body}"
        $events << body
        # IO.write("#{Dir.pwd}/events.json", $events.to_json)
        res.headers["Content-Type"] = "application/json; charset=utf-8"
        res.status = 201
        res.write json
      rescue => e
        res.status= 500
        res.headers["Content-Type"] = "application/json; charset=utf-8"
        output = {"code" => 500, "message" =>  "An error ocurred #{e}"}
        puts "#{output}"
        res.write output.to_json
      end
    end

    on "events" do
      begin
        json = req.body.read
        body = JSON.parse json
        puts "#{body}"
        print "Event received: #{body}"
        $events << body
        # IO.write("#{Dir.pwd}/events.json", $events.to_json)
        res.headers["Content-Type"] = "application/json; charset=utf-8"
        res.status = 201
        res.write json
      rescue => e
        res.status= 500
        res.headers["Content-Type"] = "application/json; charset=utf-8"
        output = {"code" => 500, "message" =>  "An error ocurred #{e}"}
        puts "#{output}"
        res.write output.to_json
      end
    end

    on "ok_only_first_time/:name" do |value|
      unless $answered.include? value
        code = 201
        message = "Only respond on suscription evaluation"
        puts "#{message}"
        $answered << value
      else # first post recveived
        code = 500
        message = "Only the first time this endpoint respond ok"
        puts "#{message}"
      end
      res.headers["Content-Type"] = "application/json; charset=utf-8"
      res.status = code
      output = {'code'=>code, 'message' => message}
      res.write output.to_json
    end

    on "returnok/:any" do |value|
      res.status = 201
      output = {'code'=>201, 'message'=>"Returning ok to /returnok/#{value}"}
      puts "#{output}"
      res.headers["Content-Type"] = "application/json; charset=utf-8"
      res.write output.to_json
    end

    on "error/:any" do |value|
      code, message = case value
      when "502"
        [400, 'Bad gateway timeout']
      when "400"
        [400, 'Bad request']
      else
        [500, 'Unknown error']
      end
      res.status code
      output = {'code' => code, 'message' => message}
      puts "#{output}"
      res.headers["Content-Type"] = "application/json; charset=utf-8"
      res.write output.to_json
    end

  end
end

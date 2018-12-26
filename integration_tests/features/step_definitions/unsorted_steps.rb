Cuando("pasan {int} segundos") do |int|
  sleep int
end

Dado("se envian una serie de {int} mensajes al tópico") do |int|
  byebug
  @events = []
  int.class
  "#{int}".to_i.times do
  	sent_event = random_event_for_topic(@topic_name)
  	@events << sent_event
  	puts sent_event if $debug
  	send_event(sent_event, security_header)
  end
end

Cuando("uno de los suscriptores borra todos los mensajes") do
  res = ask_for_missing_events(@subscriber, "5")
  puts parsed_response if $debug
  beybug
  delete_messages
end

Entonces("el tópico no debe poseer mensajes para leer") do
  res = ask_for_missing_events(@subscriber, "5")
  byebug
  puts parsed_response
end
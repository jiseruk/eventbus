Cuando("pasan {int} segundos") do |int|
  sleep int
end

Dado("se envian una serie de {int} mensajes al tópico") do |int|
  @events = []
  int.times do
  	sent_event = random_event_for_topic(@topic_name)
  	@events << sent_event
  	puts sent_event if $debug
  	send_event(sent_event, security_header)
  end
end

Cuando("uno de los suscriptores borra todos los mensajes") do
  ask_for_missing_events(@subscriber, "5")
  puts parsed_response if $debug
  delete_messages
end

Entonces("el tópico no debe poseer mensajes para leer") do
  ask_for_missing_events(@subscriber, "5")
  fail "Se encontraron mensajes cuando no debería" if is_there_messages?
end
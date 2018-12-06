Dado("que voy a notificar un evento cualquiera") do
  # do nothing
end

Dado("que voy notificar un evento a un tópico existente") do
  @topic_name = create_topic
end

Dado("que voy a notificar un evento a un tópico inexistente") do
  @sent_event = random_event_for_topic("Unknown-topic-#{Time.now.to_i}")
  send_event(@sent_event)
end

Dado("que estuve sin atender eventos por un tiempo debido a X motivo") do
  # do nothing
end

Dado("que se notificó un evento a un topico que estoy suscripto") do
  @sent_event = random_event_for_topic(@topic_name)
  puts @sent_event if $debug
  send_event(@sent_event)
end

Dado("que se notificó un evento a un topico que estoy suscripto en modo pull con visibilidad de 5 segundos") do
  @sent_event = random_event_for_topic(@topic_name, 5)
  puts @sent_event if $debug
  send_event(@sent_event)
end

Dado("se envía un evento al topico con visibilidad de {int} segundos") do |int|
  pending # Write code here that turns the phrase above into concrete actions
end

Dado("se envía un evento al topico") do
  @sent_event = random_event_for_topic(@topic_name, 5)
  puts @sent_event if $debug
  send_event(@sent_event)
end










Cuando("envío una notificación al tópico") do
  @sent_event = random_event_for_topic(@topic_name)
  send_event(@sent_event)
end

Cuando("envío una notificación al tópico inexistente") do
  @sent_event = random_event_for_topic('unknown-topic')
  send_event(@sent_event)
end

Cuando("envío una notificación sin indicar el tópico") do
  @sent_event = random_event_for_topic()
  send_event(@sent_event)
end

Cuando("envío una notificación sin el payload") do
  @sent_event = random_event_for_topic(@topic_name)
  @sent_event.delete('payload')
  send_event(@sent_event)
end

Cuando("envío una notificación con un payload vacío") do
  @sent_event = random_event_for_topic(@topic_name)
  @sent_event['payload']={}
  send_event(@sent_event)
end

Cuando("envío una notificación con un payload que no es JSON") do
  @sent_event = random_event_for_topic(@topic_name)
  @sent_event['payload']="Un string que no es JSOn"
  send_event(@sent_event)
end

Cuando("consulto los mensajes sin indicar quien soy") do
  ask_for_missing_events()
end

Cuando("consulto los mensajes perdidos al tópico") do
  ask_for_missing_events(@subscriber, "5")
end

Cuando("uno de los suscriptores consulta el mensaje sin hacer mas que leerlo") do
  puts "Subscriber: #{@subscriber}" if $debug
  res = ask_for_missing_events(@subscriber, "1")
  puts parsed_response if $debug
end

Cuando("consulto los mensajes al tópico") do
  puts "Subscriber: #{@subscriber}" if $debug  
  res = ask_for_missing_events(@subscriber, "5")
  puts parsed_response if $debug
end

Cuando("uno de los suscriptores consulta el mensaje y lo borra") do
  puts "Subscriber: #{@subscriber}" if $debug  
  res = ask_for_missing_events(@subscriber, "5")
  puts parsed_response if $debug
  delete_messages
end







Entonces("debo obtener los mensajes existentes") do
  unless exists_events? @sent_event
    fail "No se encontró el mensaje esperado
    Enviado: #{sent_event}
    Esperado: #{messages}"
  end
end

Entonces("debo obtener el mensaje de error que el tópico no existe") do
  expected_msg = "The topic unknown-topic doesn't exist"
  got = response_message
  fail "El mensaje obtenido fué:
  '#{got}'
  Se esperaba:
  #{expected_msg}" unless got == expected_msg
end

Entonces("los sucriptores debe recibir dicho evento") do
  unless event_transmitted?
    fail "No se recibió el evento o el recibido no es lo esperado.
    Enviado: #{sent_event}
    Esperado: #{last_event}"
  end
end

Entonces("debo obtener los mensajes enviados que no pude atender") do
  got = last_message
  fail "No se encontró el mensaje. Obtenido '#{got}':
  Esperado: '#{@sent_event}'" unless are_messages_equals?(got, @sent_event)
end

Entonces("los sucriptores debe poder levantar el mensaje") do
  ask_for_missing_events(@subscriber, "5")
  fail "No se econtró el mensaje" unless exists_events? @sent_event
end

Entonces("los mensajes deben tener la fecha de enviado") do
  fail "No se encontro la fecha de envío del mensaje" unless message_has_sent_date?
end

Entonces("el otro suscriber no puede leer el mensaje") do
  puts "Subscriber: #{@subscriber}" if $debug
  res = ask_for_missing_events(@subscriber, "5")
  puts parsed_response if $debug
  fail "Se obtuvo mensajes cuando no debía ya que los había tomado el primer subscriber" if exists_events? @sent_event
end

Entonces("el otro suscriptor no encontrará el mensaje") do

end

Entonces("el otro suscriber puede leer el mensaje") do
  puts "Subscriber: #{@subscriber}" if $debug
  res = ask_for_missing_events(@subscriber, "5")
  puts parsed_response if $debug
  fail "No se encontró el mensaje esperado" unless exists_events? @sent_event
end
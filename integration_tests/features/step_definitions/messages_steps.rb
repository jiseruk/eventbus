Dado("que voy a notificar un mensaje cualquiera") do
  # do nothing
end

Dado("que voy notificar un mensaje a un tópico existente") do
  @topic_name = create_topic
end

Dado("que voy a notificar un mensaje a un tópico inexistente") do
  @sent_event = random_event_for_topic("Unknown-topic-#{Time.now.to_i}")
  send_event(@sent_event, some_header)
end

Dado("que estuve sin atender mensajes por un tiempo debido a X motivo") do
  # do nothing
end

Dado("que se notificó un mensaje a un topico que estoy suscripto") do
  @sent_event = random_event_for_topic(@topic_name)
  puts @sent_event if $debug
  send_event(@sent_event, security_header)
end

Dado("que se notificó un mensaje a un topico que estoy suscripto en modo pull con visibilidad de 5 segundos") do
  @sent_event = random_event_for_topic(@topic_name, 5)
  puts @sent_event if $debug
  send_event(@sent_event, security_header)
end

Dado("se envía un mensaje al topico con visibilidad de {int} segundos") do |int|
  pending # Write code here that turns the phrase above into concrete actions
end

Dado("se envía un mensaje al topico") do
  @sent_event = random_event_for_topic(@topic_name, 5)
  puts @sent_event if $debug
  send_event(@sent_event, security_header)
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








Cuando("envío una notificación al tópico") do
  @sent_event = random_event_for_topic(@topic_name)
  send_event(@sent_event, security_header)
end

Cuando("envío una notificación al tópico sin pasar el token de seguridad") do
  @sent_event = random_event_for_topic(@topic_name)
  send_event(@sent_event)
end

Cuando("envío una notificación al tópico inexistente") do
  @sent_event = random_event_for_topic('unknown-topic')
  send_event(@sent_event, security_header)
end

Cuando("envío una notificación sin indicar el tópico") do
  @sent_event = random_event_for_topic()
  send_event(@sent_event, security_header)
end

Cuando("envío una notificación sin el payload") do
  @sent_event = random_event_for_topic(@topic_name)
  @sent_event.delete('payload')
  send_event(@sent_event, security_header)
end

Cuando("envío una notificación con un payload vacío") do
  @sent_event = random_event_for_topic(@topic_name)
  @sent_event['payload']={}
  send_event(@sent_event, security_header)
end

Cuando("envío una notificación con un payload que no es JSON") do
  @sent_event = random_event_for_topic(@topic_name)
  @sent_event['payload']="Un string que no es JSOn"
  send_event(@sent_event, security_header)
end

Cuando("consulto los mensajes sin indicar quien soy") do
  ask_for_missing_events()
end

Cuando("consulto los mensajes perdidos al tópico") do
  ask_for_missing_events(@subscriber, "5")
end

Cuando("uno de los suscriptores consulta el mensaje sin hacer mas que leerlo") do
  puts "Subscriber: #{@subscriber}" if $debug
  ask_for_missing_events(@subscriber, "10")
  @first_messages_to_delete = messages
  @first_messages = messages_ids
  puts @first_messages if $debug
  puts parsed_response if $debug
end

Cuando("consulto los mensajes al tópico") do
  puts "Subscriber: #{@subscriber}" if $debug  
  ask_for_missing_events(@subscriber, "5")
  puts parsed_response if $debug
end

Cuando("uno de los suscriptores consulta el mensaje y lo borra") do
  puts "Subscriber: #{@subscriber}" if $debug  
  ask_for_missing_events(@subscriber, "5")
  puts parsed_response if $debug
  delete_messages
end

Cuando("uno de los suscriptores borra todos los mensajes") do
  ask_for_missing_events(@subscriber, "10")
  puts parsed_response if $debug
  delete_messages
end

Cuando("el otro subscriber consulta el mensaje") do
  ask_for_missing_events(@subscriber, "10")
end

Cuando("el primer subscritor intenta borrar el mensaje") do
  delete_messages @subscriber, @first_messages_to_delete
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

Entonces("los sucriptores debe recibir dicho mensaje") do
  unless event_transmitted?
    fail "No se recibió el mensaje o el recibido no es lo esperado.
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
  ask_for_missing_events(@subscriber, "5")
  puts parsed_response if $debug
  fail "Se obtuvo mensajes cuando no debía ya que los había tomado el primer subscriber" if exists_events? @sent_event
end

Entonces("el otro suscriptor no encontrará el mensaje") do
  ask_for_missing_events(@subscriber, "5")
  fail "Se encontró un mensaje cuando se esperaba lo contrario" if is_there_messages?
end

Entonces("el otro suscriber puede leer el mensaje") do
  puts "Subscriber: #{@subscriber}" if $debug
  ask_for_missing_events(@subscriber, "5")
  puts parsed_response if $debug
  fail "No se encontró el mensaje esperado" unless exists_events? @sent_event
end

Entonces("el tópico no debe poseer mensajes para leer") do
  ask_for_missing_events(@subscriber, "10")
  puts parsed_response if $debug
  fail "Se encontraron #{messages_size} mensajes cuando no debería" if is_there_messages?
end

Entonces("el mensaje no puede ser borrado porque el delete token lo tienen el último lector") do
  sleep 5
  ask_for_missing_events(@subscriber, "10")
  actual_messages= messages_ids
  puts actual_messages if $debug
  fail "Se borraron los mensajes cuando no debería ya que el token utilizado había expirado" unless @first_messages == actual_messages
end
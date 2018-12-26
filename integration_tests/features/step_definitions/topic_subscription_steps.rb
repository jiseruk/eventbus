# GIVEN
Dado(/el tópico tienen suscriptores de tipo (push|pull)/) do |type|
  @type = type
  @subscriber = random_subscriber_name
  puts "Subscriber: #{@subscriber}" if $debug
  @endpoint = subscriber_endpoint
  create_a_subscription_of_type type
end

Dado("dos suscriptores del mismo nombre en modo pull al tópico") do
  @subscriber = random_subscriber_name
  puts "Subscriber: #{@subscriber}" if $debug
  @endpoint = subscriber_endpoint
  create_a_subscription_of_type 'pull', 5
end






# WHEN

Dado(/estoy suscripto en modo (push|pull) a dicho topico con un endpoint que solo responde a la suscripcion/) do |type|
  @type = type
  @subscriber = random_subscriber_name
  @endpoint = first_time_endpoint
  create_a_subscription_of_type type, 5
end

Dado("estoy suscripto en modo pull a dicho topico") do
  @subscriber = random_subscriber_name
  create_a_subscription_of_type 'pull', 5
end

Cuando(/me suscribo en modo (push|pull) a un tópico que no existe/) do |type|
  @type = type
  @topic_name = random_topic_name
  @subscriber = random_subscriber_name
  @endpoint = first_time_endpoint
  create_a_subscription_of_type type
end


Cuando("intento suscribirme a un tópico sin pasar el modo de suscripcion") do
	@topic_name = random_topic_name
	@subscriber = random_subscriber_name
	@endpoint = first_time_endpoint
	subscribe_to_topic(topic_name: @topic_name, subscriber: @subscriber, endpoint: @endpoint)
end

Cuando(/me suscribo en modo (push|pull) al tópico correctamente/) do |type|
  @type = type
  create_a_subscription_of_type(type)
end

Cuando(/me suscribo en modo (push|pull) al tópico sin pasar ningún dato/) do |type|
  @type = type
  subscribe_to_topic(topic_name: nil, subscriber: nil, type: nil, endpoint: nil)
end

Cuando(/me suscribo en modo (push|pull) al tópico correctamente con un suscriber que escuchará mensajes/) do |type|
  @type = type
  @subscriber = random_subscriber_name
  subscribe_to_topic(topic_name: @topic_name, subscriber: @subscriber, type: type, endpoint: subscriber_endpoint)
end

Cuando(/intento suscribirme en modo (push|pull) a un tópico sin pasar el nombre del suscriber/) do |type|
  @type = type
  opts = {topic_name: @topic_name, type: type, endpoint: subscriber_endpoint}
  opts[:visibility_timeout]= 10 if type == 'pull'
  subscribe_to_topic(opts)
end

Cuando(/intento suscribirme en modo (push|pull) a un tópico sin pasar el nombre del tópico/) do |type|
  @type = type
  @subscriber = random_subscriber_name
  subscribe_to_topic(topic_name: nil, subscriber: @subscriber, type: type, endpoint: subscriber_endpoint)
end

Cuando(/intento suscribirme en modo (push|pull) a un tópico sin pasar el endpoint de notificaciones/) do |type|
  @type = type
  subscribe_to_topic(topic_name: @topic_name, subscriber: random_subscriber_name, type: type, endpoint: "")
end

Cuando(/intento suscribirme en modo (push|pull) a un tópico con el endpoint sin un formato válidop de url/) do |type|
  @type = type
  @endpoint = 'sarasa'
  @subscriber = random_subscriber_name
  create_a_subscription_of_type(type)
end

Dado(/un suscriber ya suscripto en modo (push|pull) al tópico/) do |type|
  @type = type
  @subscriber = random_subscriber_name
  create_a_subscription_of_type(type)
  @endpoint = parsed_response["endpoint"]
end

Dado(/un suscriber suscripto en modo push al tópico cuyo endpoint debe ser único/) do
  @subscriber = random_subscriber_name
  @endpoint = random_fake_endpoint
  create_a_subscription_of_type('push')
end

Cuando(/intento suscribirme en modo (push|pull) con el mismo nombre de suscriber/) do |type|
  @type = type
  @endpoint = random_fake_endpoint
  create_a_subscription_of_type(type)
end

Cuando("con el mismo nombre de suscriber intento suscribirme en modo pull") do
  create_a_subscription_of_type('pull')
end

Cuando(/intento suscribirme en modo (push|pull) con el mismo enpoint que el suscriber/) do |type|
  @type = type
  @subscriber = random_subscriber_name
  create_a_subscription_of_type(type)
end

Cuando("intento suscribirme en modo pull pasando un endpoint") do 
  @endpoint = random_fake_endpoint
  @subscriber = random_subscriber_name
  create_a_subscription_of_type('pull')
end

Cuando(/intento suscribirme en modo (push|pull) con un endpoint que no responde/) do |type|
  @type = type
  @subscriber = random_subscriber_name
  @endpoint = 'http://localhost:1223/bla'
  create_a_subscription_of_type(type)
end

Cuando(/un subscriber se registra en modo (push|pull) en todos los tópicos/) do |type|
  @type = type
  @subscriber = random_subscriber_name
  @suscriptions = topics.map do |t_name|
    puts "Subscribing with name: #{t_name}" if $debug
    subscribe_to_topic(topic_name: t_name, subscriber: @subscriber, type: type, endpoint: random_fake_endpoint)
    puts "Subscription response: #{parsed_response}" if $debug
    status_code
  end
end






Entonces("debo obtener el mensaje de error de suscriptor existente") do
  expected_msg = "Subscription with name #{@subscriber} already exists"
  got = response_message
  fail "Se esperaba: #{expected_msg}.
  Se obtuvo: #{got}
  Status code; #{response.code}" unless got == expected_msg

end

Entonces("debo obtener el mensaje de error que el vampo endpoint no es permtido para subscripcion pull") do
  expected_msg = "The endpoint field is invalid for pull subscribers"
  got = response_message
  fail "Se esperaba: #{expected_msg}.
  Se obtuvo: #{got}
  Status code; #{response.code}" unless got == expected_msg

end

Entonces("debo obtener el mensaje de error de endpoint existente") do
  expected_msg = "The endpoint #{@endpoint} is used by the subscriber"
  got = response_message
  fail "Se esperaba: #{expected_msg}.
  Se obtuvo: #{got}
  Status code; #{response.code}" unless got.start_with? expected_msg

end

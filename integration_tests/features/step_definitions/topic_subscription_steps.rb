# GIVEN
Dado(/el tópico tienen suscriptores de tipo (push|pull)/) do |type|
  @subscriber = random_subscriber_name
  puts "Subscriber: #{@subscriber}"
  subscribe_to_topic(topic_name: @topic_name, subscriber: @subscriber, type: type, endpoint: $configuration.subscriber_endpoint)
end










# WHEN

Dado(/estoy suscripto en modo (push|pull) a dicho topico con un endpoint que solo responde a la suscripcion/) do |type|
  @subscriber = 'outage'
  @endpoint = first_time_endpoint
  puts "Subscriber: #{@subscriber}"
  puts "Endpoint: #{@endpoint}"
  subscribe_to_topic(topic_name: @topic_name, subscriber: @subscriber, type: type, endpoint: @endpoint)
end

Dado("estoy suscripto en modo pull a dicho topico") do
  @subscriber = 'outage'
  puts "Subscriber: #{@subscriber}"
  subscribe_to_topic(topic_name: @topic_name, subscriber: @subscriber, type: "pull")
end

Cuando(/me suscribo en modo (push|pull) a un tópico que no existe/) do |type|
  @topic_name = random_topic_name
  @subscriber = random_subscriber_name
  @endpoint = first_time_endpoint
  puts "Topic: #{@topic_name}"
  puts "Subscriber: #{@subscriber}"
  puts "Endpoint: #{@endpoint}"
  subscribe_to_topic(topic_name: @topic_name, subscriber: @subscriber, type: type, endpoint: @endpoint)
end


Cuando("intento suscribirme a un tópico sin pasar el modo de suscripcion") do
	@topic_name = random_topic_name
  	@subscriber = random_subscriber_name
  	@endpoint = first_time_endpoint
  	puts "Topic: #{@topic_name}"
  	puts "Subscriber: #{@subscriber}"
  	puts "Endpoint: #{@endpoint}"
  	subscribe_to_topic(topic_name: @topic_name, subscriber: @subscriber, endpoint: @endpoint)
end

Cuando(/me suscribo en modo (push|pull) al tópico correctamente/) do |type|
  create_a_subscription_of_type(type)
end

Cuando(/me suscribo en modo (push|pull) al tópico sin pasar ningún dato/) do |type|
  subscribe_to_topic(topic_name: nil, subscriber: nil, type: nil, endpoint: nil)
end

Cuando(/me suscribo en modo (push|pull) al tópico correctamente con un suscriber que escuchará eventos/) do |type|
  @subscriber = random_subscriber_name
  subscribe_to_topic(topic_name: @topic_name, subscriber: @subscriber, type: type, endpoint: $configuration.suscriber_endpoint)
end

Cuando(/intento suscribirme en modo (push|pull) a un tópico sin pasar el nombre del suscriber/) do |type|
  subscribe_to_topic(topic_name: @topic_name, endpoint: 'http://localhost:9292/suscriber', type: type)
end

Cuando(/intento suscribirme en modo (push|pull) a un tópico sin pasar el endpoint de notificaciones/) do |type|
  subscribe_to_topic(topic_name: @topic_name, subsciber: "", type: type, endpoint: "")
end

Cuando(/intento suscribirme en modo (push|pull) a un tópico con el endpoint sin un formato válidop de url/) do |type|
  subscribe_to_topic(topic_name: @topic_name, subsciber: random_subscriber_name, type: type, endpoint: "sarasa")
end

Dado(/un suscriber ya suscripto en modo (push|pull) al tópico/) do |type|
  create_a_subscription_of_type(type)
end

Cuando(/intento suscribirme en modo (push|pull) con el mismo nombre de suscriber/) do |type|
  @fake_endpoint = random_fake_endpoint
  subscribe_to_topic(topic_name: @topic_name, subscriber: @subscriber, type: type, endpoint: @fake_endpoint)
end

Cuando("con el mismo nombre de suscriber intento suscribirme en modo pull") do
  subscribe_to_topic(topic_name: @topic_name, subscriber: @subscriber, type: "pull")
end

Cuando(/intento suscribirme en modo (push|pull) con el mismo enpoint que el suscriber/) do |type|
  @fake_endpoint = random_fake_endpoint
  subscribe_to_topic(topic_name: @topic_name, subscriber: random_subscriber_name, type: type, endpoint: @fake_endpoint)
end

Cuando("intento suscribirme en modo pull pasando un endpoint") do 
  @fake_endpoint = random_fake_endpoint
  subscribe_to_topic(topic_name: @topic_name, subscriber: random_subscriber_name, type: "pull", endpoint: @fake_endpoint)
end

Cuando(/intento suscribirme en modo (push|pull) con un endpoint que no responde/) do |type|
  @endpoin = 'http://localhost:1223/bla'
  subscribe_to_topic(topic_name: @topic_name, subscriber: @subscriber, type: type, endpoint: @endpoint)
end

Cuando(/un subscriber se registra en modo (push|pull) en todos los tópicos/) do |type|
  @subscriber = random_subscriber_name
  @suscriptions = topics.map do |topic_name|
    subscribe_to_topic(topic_name: topic_name, subscriber: @subscriber, type: type, endpoint: random_fake_endpoint)
    puts @response.parse
    @response.status.code
  end
end









# THEN



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
  expected_msg = "Endpoint already exists on topic"
  got = response_message
  fail "Se esperaba: #{expected_msg}.
  Se obtuvo: #{got}
  Status code; #{response.code}" unless got == expected_msg

end

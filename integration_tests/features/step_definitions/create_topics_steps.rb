Dado("que deseo generar un tópico nuevo") do
  @topic_name = "TopicTest#{$execution_id}_#{Time.now.to_i}"
end

Dado("un tópico existente en el event bus") do
  @topic_name = create_topic
end

Dado("un topico determinado") do
  @topic_name = create_topic
  puts "Topic Name: #{@topic_name}" if $debug
end

Dado("que soy owner de un tópico") do
  @topic_name = create_topic
  puts "Topic Name: #{@topic_name}" if $debug
end

Dado("tres tópicos diferentes") do
  create_topics(3)
  puts "#{topics}" if $debug
end

Dado("un tópico existente") do
  @topic_name = create_topic
  puts "Topic Name: #{@topic_name}" if $debug
end








Cuando("intento crear un tópico con el nombre del existente") do
  post_topic(topic_name: @topic_name)
end

Cuando("ingreso los datos requeridos del tópico") do
  post_topic(topic_name: @topic_name, engine: "AWS")
end

Cuando("ingreso los datos sin pasar el nombre") do
  post_topic(topic_name:"")
end

Cuando("ingreso los datos sin pasar el engine") do
  post_topic(topic_name:@topic_name, engine: "")
end

Cuando("ingreso los datos con un engine que no existe") do
  @engine = "sarasa"
  post_topic(topic_name:@topic_name, engine: @engine)
end












Entonces("debo obtener una respuesta de aceptado") do
  message = response_message
  code = status_code
  fail "Se obtuvo: #{code}
  #{response_message}" unless success?
end

Entonces("debo recibir los datos de id, fecha de creacion y nombre de topico") do
  not_found = []
  not_found << "No se encontró id de creacion" unless has_topic_id?
  not_found << "No se encontró la fecha de creación" unless has_creation_date?
  not_found << "No se encontó el nombre del topico #{@topic_name}" unless has_topic_name? @topic_name
  fail "#{not_found}" unless not_found.empty?
end

Entonces("debo obtener una respuesta de tópico existente") do
  fail "La respuesta no fué al esperada. Se obtuvo #{response_message}" unless response_message == "Topic with name #{@topic_name} already exists"
end

Entonces("debo obtener una respuesa que indique que el nombre es obligatorio") do
  expected_msg = "name: The field is required."
  got = response_message
  fail "Se esperaba el mensaje: '#{expected_msg}'.
  Se obtuvo '#{got}'" unless got == expected_msg
end

Entonces("debo obtener una respuesa que indique que el engine es obligatorio") do
  expected_msg = "engine: The field is required."
  got = response_message
  fail "Se esperaba el mensaje: '#{expected_msg}'.
  Se obtuvo '#{got}'" unless got == expected_msg
end

Entonces("debo obtener una respuesa que indique que el engine no existe") do
  expected_msg = "engine: Must be one of [AWS, AWSStream]."
  got = response_message
  fail "Se esperaba el mensaje: '#{expected_msg}'.
  Se obtuvo '#{got}'" unless got == expected_msg
end

Entonces("la respuesta debe tener los valores {string}") do |fields|
  fields = fields.split(", ").sort
  got = parsed_response.keys.sort
  fail "Se recibió #{got}
  #{parsed_response}" unless got == fields
end

Entonces("debo obtener el mensaje de tópico inexistente") do
  expected_msg = "The topic #{@topic_name} doesn't exist"
  got = response_message
  fail "Se esperaba el mensaje: '#{expected_msg}'.
  Se obtuvo '#{got}'" unless got == expected_msg
end


Entonces("debo recibir una respuesta de suscripción correta") do
  fail "No se recibió la respuesta correcta" unless subscribed?
end

Entonces("debo obtener un status code {int}") do |expected_code|
  received_status_code = status_code
  fail "Se obtuvo el status code #{received_status_code}
  Response message #{response_message}" unless "#{status_code}" == "#{expected_code}"
end

Entonces("debo obtener el mensaje de error {string}") do |expected_msg|
  got = response_message
  fail "Se obtuvo el mensaje '#{got}'
  Status code: #{status_code}" unless got == expected_msg
end

Entonces("Entonces debo obtener el mensaje de error de suscriptor existente") do
  got = response_message
  expected_message = "Subscription with name #{@subscriber} already exists"
  fail "Se experaba: '#{expected_msg}'
  Se obtuvo: '#{got}'" unless got == expected_message
end

Entonces("debo obtener el mensaje de error the endpoint que no responde") do
  expected_msg = "The endpoint #{@endpoint} should return 2xx to a POST HTTP Call, but return error"
  got = response_message
  fail "Se obtuvo el mensaje: '#{got}'
  Se esperaba '#{expected_msg}'" unless got.start_with? expected_msg
end

Entonces("todas las suscripciones deben resultar correctas") do
  fail "No todas las suscripciones fueron correctas.
Se obtuvo los status code siguientes: #{@suscriptions}" unless @suscriptions.uniq == [201]
end


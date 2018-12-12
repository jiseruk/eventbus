Dado("que existen topicos creados en el even bus") do
  @created_topics = []
  3.times do
  	@created_topics << create_topic
  end
end

Cuando("solicito el listado de tópicos") do
  list_topics
end

Entonces("debe aparecer la lista de topicos existentes") do
  fail "No se encontró al menos uno de los tópicos creados #{@created_topics}" unless exist_topics? @created_topics
end

Cuando("consulto los subscribers del topico") do
  list_subscribers_for_topic @topic_name
end

Entonces("debe aparecer el subscriber en la lista del topico") do
  fail "No se encontró el subscriber #{@subscriber}" unless exists_subscriber? @subscriber
end
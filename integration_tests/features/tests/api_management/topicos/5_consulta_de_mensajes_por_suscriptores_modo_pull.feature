# language: es
Característica: Consulta de mensajes por suscriptores en modo pull
	Como un usuario suscripto a un tópico en modo pull
	Cuando consulto los mensajes en el tópico
	Quiero ver los mensajes disponibles

# method: GET
# endpoint: messages?subscriber=<subscriber_name>&max_messages=<int>

Antecedentes: Topico existente con suscriptores en modo pull
	Dado un topico determinado
	Y estoy suscripto en modo pull a dicho topico

Escenario: Consulta de mensajes
	Dado que se notificó un evento a un topico que estoy suscripto en modo pull
	Cuando consulto los mensajes al tópico
	Entonces debo obtener los mensajes existentes


Escenario: Consulta de mensajes sin indicar quien soy
	Cuando consulto los mensajes sin indicar quien soy
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'subscriber and max_messages fields are required'
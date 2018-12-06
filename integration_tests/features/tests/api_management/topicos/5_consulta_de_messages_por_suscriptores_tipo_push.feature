# language: es
Característica: Consulta de mensajes por suscriptores en modo push
	Como un usuario suscripto a un tópico que no pudo atender mensajes
	Cuando consulto mis mensajes perdidos
	Quiero ver los mensajes pendientes

# method: GET
# endpoint: messages?subscriber=<subscriber_name>&max_messages=<int>

Antecedentes: Topico existente con suscriptos
	Dado un topico determinado
	Y estoy suscripto en modo push a dicho topico con un endpoint que solo responde a la suscripcion

Escenario: Consulta de mensajes perdidos
	Dado que estuve sin atender eventos por un tiempo debido a X motivo
	Y que se notificó un evento a un topico que estoy suscripto
	Cuando consulto los mensajes perdidos al tópico
	Entonces debo obtener los mensajes enviados que no pude atender
	Y los mensajes deben tener la fecha de enviado

Escenario: Consulta de mensajes sin indicar quien soy
	Cuando consulto los mensajes sin indicar quien soy
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'max_messages: The field is required; subscriber: The field is required.'
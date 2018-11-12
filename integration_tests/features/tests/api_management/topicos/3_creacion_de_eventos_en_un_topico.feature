# language: es
Característica: Creación de un evento en un tópico
	Como una aplicación que es owner de un tópico
	Cuando envía un evento a dicho tópico
	Quiere enviar la notificación


# method: POST
# endpoint: /messages
# body = 
# 	{
# 		"topic": "topic_name", 
# 		"payload": {}
# 	}


Escenario: Envio de evento a un tópico existente para suscriptores push
	Dado que soy owner de un tópico
	Y el tópico tienen suscriptores de tipo push
	Cuando envío una notificación al tópico
	Entonces los sucriptores debe recibir dicho evento

Escenario: Envio de evento a un tópico existente para suscriptores pull
	Dado que soy owner de un tópico
	Y el tópico tienen suscriptores de tipo pull
	Cuando envío una notificación al tópico
	Entonces los sucriptores debe poder levantar el mensaje

Escenario: Envío de un evento a un tópico inexistente
	Dado que voy a notificar un evento a un tópico inexistente
	Cuando envío una notificación al tópico inexistente
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error que el tópico no existe

Escenario: Envío de un evento sin indicar el tópico
	Dado que voy a notificar un evento cualquiera
	Cuando envío una notificación sin indicar el tópico
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'Topic cannot be null. Topic must be defined'

Escenario: Envío de un evento a un tópico sin indicar el payload
	Dado que voy notificar un evento a un tópico existente
	Cuando envío una notificación sin el payload
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'Payload cannot be null. Payload must be defined'

Escenario: Envío de un evento a un tópico con un payload vacío
	Dado que voy notificar un evento a un tópico existente
	Cuando envío una notificación con un payload vacío
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'Payload cannot be empty. Payload must have content'	

Escenario: Envío de un evento a un tópico con un payload que no es un JSON
	Dado que voy notificar un evento a un tópico existente
	Cuando envío una notificación con un payload que no es JSON
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'Payload should be a json'	



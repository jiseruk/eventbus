# language: es
Característica: Eliminación de mensajes en un tópico
	Como una aplicación que es owner de un tópico
	Cuando existen mensajes enviados
	Quiere que sus mensajes puedan ser eliminados


# method: DELETE
# endpoint: /messages
# body = 
# 	{
#     "topic": "topic_name", 
#     "payload": {}
# 	}


# Escenario: Envio de mensaje a un tópico existente para suscriptores push
# 	Dado que soy owner de un tópico
# 	Y el tópico tienen suscriptores de tipo push
# 	Cuando envío una notificación al tópico
# 	Entonces los sucriptores debe recibir dicho mensaje

# Escenario: Envio de mensaje a un tópico existente para suscriptores pull
# 	Dado que soy owner de un tópico
# 	Y el tópico tienen suscriptores de tipo pull
# 	Cuando envío una notificación al tópico
# 	Entonces los sucriptores debe poder levantar el mensaje

# Escenario: Envio de mensaje a un tópico existente sin indicar el token de seguridad
# 	Dado que soy owner de un tópico
# 	Y el tópico tienen suscriptores de tipo pull
# 	Cuando envío una notificación al tópico sin pasar el token de seguridad
# 	Y debo obtener el mensaje de error 'The X-Publish-Token header is invalid'

# Escenario: Envío de un mensaje a un tópico inexistente
# 	Dado que voy a notificar un mensaje a un tópico inexistente
# 	Cuando envío una notificación al tópico inexistente
# 	Entonces debo obtener un status code 400
# 	Y debo obtener el mensaje de error que el tópico no existe

# Escenario: Envío de un mensaje sin indicar el tópico
# 	Dado que voy a notificar un mensaje cualquiera
# 	Cuando envío una notificación sin indicar el tópico
# 	Entonces debo obtener un status code 400
# 	Y debo obtener el mensaje de error 'topic: The field is required.'

# Escenario: Envío de un mensaje a un tópico sin indicar el payload
# 	Dado que voy notificar un mensaje a un tópico existente
# 	Cuando envío una notificación sin el payload
# 	Entonces debo obtener un status code 400
# 	Y debo obtener el mensaje de error 'payload: The field is required.'

# Escenario: Envío de un mensaje a un tópico con un payload vacío
# 	Dado que voy notificar un mensaje a un tópico existente
# 	Cuando envío una notificación con un payload vacío
# 	Entonces debo obtener un status code 400
# 	Y debo obtener el mensaje de error 'payload: The field is required.'	

# Escenario: Envío de un mensaje a un tópico con un payload que no es un JSON
# 	Dado que voy notificar un mensaje a un tópico existente
# 	Cuando envío una notificación con un payload que no es JSON
# 	Entonces debo obtener un status code 400
# 	Y debo obtener el mensaje de error 'payload: it should be a valid json object.'	



# language: es
@topicos @suscripcion @pull
Característica: Suscripción a tópico de tipo pull
	Como un servicio que quiere enterarse de la novedades de un topico
	Cuando realiza una sucripción a un tópico en modo pull
	Debe quedar suscripto al tópico para luego poder consultar las novedades

# method: POST

# endpoint: /subscriptions

# 	body = {
# 		"topic" : "topic_name", 
# 		"name" : "subscriber name",
#     "type" : "pull",
#     "visibility_timeout": Int Seconds
# 	}


Escenario: Suscripición exitosa a un topico en modo pull
	Dado un tópico existente
	Cuando me suscribo en modo pull al tópico correctamente
	Entonces debo recibir una respuesta de suscripción correta
	Y la respuesta debe tener los valores 'name, topic, type, visibility_timeout'

Escenario: Suscripción a tópico inexistente
	Cuando me suscribo en modo pull a un tópico que no existe
	Entonces debo obtener un status code 404
	Y debo obtener el mensaje de tópico inexistente

Escenario: Suscripción sin datos
	Dado un tópico existente
	Cuando me suscribo en modo pull al tópico sin pasar ningún dato
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'name: The field is required; topic: The field is required; type: The field is required.'

Escenario: Suscripción sin indicar el nombre del tópico
	Dado un tópico existente
	Cuando intento suscribirme en modo pull a un tópico sin pasar el nombre del suscriber
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'name: The field is required.'

Escenario: Suscripción sin indicar el nombre de suscriber
	Dado un tópico existente
	Cuando intento suscribirme en modo pull a un tópico sin pasar el nombre del suscriber
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'name: The field is required.'

Escenario: Suscripción sin indicar el tipo de suscripcion
	Cuando intento suscribirme a un tópico sin pasar el modo de suscripcion
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'type: The field is required.'

Escenario: Suscripción de suscriber con nombre existente
	Dado un tópico existente
	Y un suscriber ya suscripto en modo pull al tópico
	Cuando con el mismo nombre de suscriber intento suscribirme en modo pull
	Entonces debo obtener el mensaje de error de suscriptor existente

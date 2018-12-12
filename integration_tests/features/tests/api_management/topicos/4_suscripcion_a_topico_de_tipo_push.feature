# language: es
@topicos @suscripcion
Característica: Suscripción a tópico de tipo push
	Como un servicio que quiere enterarse de la novedades de un topico
	Cuando realiza una sucripción a un tópico en modo push
	Debe quedar suscripto al tópico listo para recibir las notificaciones

# method: POST

# endpoint: /subscriptions

# 	body = {
#     "topic" : "topic_name", 
#     "name" : "subscriber name",
#     "type" : "push",
#     "endpoint" : "/host:port/something"
# 	}


Escenario: Suscripición exitosa a un topico en modo push
	Dado un tópico existente
	Cuando me suscribo en modo push al tópico correctamente
	Entonces debo recibir una respuesta de suscripción correta
	Y la respuesta debe tener los valores 'endpoint, name, topic, type'

Escenario: Suscripción a tópico inexistente
	Cuando me suscribo en modo push a un tópico que no existe
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de tópico inexistente

Escenario: Suscripción sin datos
	Dado un tópico existente
	Cuando me suscribo en modo push al tópico sin pasar ningún dato
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'name: The field is required; topic: The field is required; type: The field is required.'

Escenario: Suscripción sin indicar el nombre del tópico
	Dado un tópico existente
	Cuando intento suscribirme en modo push a un tópico sin pasar el nombre del tópico
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'topic: The field is required.'

Escenario: Suscripción sin indicar el nombre de suscriber
	Dado un tópico existente
	Cuando intento suscribirme en modo push a un tópico sin pasar el nombre del suscriber
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'name: The field is required.'

Escenario: Suscripción sin indicar el tipo de suscripcion
	Dado un tópico existente
	Cuando intento suscribirme a un tópico sin pasar el modo de suscripcion
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'type: The field is required.'

Escenario: Suscripción de suscriber con nombre existente
	Dado un tópico existente
	Y un suscriber ya suscripto en modo push al tópico
	Cuando intento suscribirme en modo push con el mismo nombre de suscriber
	Entonces debo obtener el mensaje de error de suscriptor existente

@EB-34 @bugs
Escenario: Suscripción de suscriber con endpoint existente
	Dado un tópico existente
	Y un suscriber ya suscripto en modo push al tópico
	Cuando intento suscribirme en modo push con el mismo enpoint que el suscriber
	Entonces debo obtener el mensaje de error de endpoint existente	

Escenario: Suscripción con un endpoint que no responde
	Dado un tópico existente
	Cuando intento suscribirme en modo push con un endpoint que no responde
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error the endpoint que no responde

Escenario: Suscripción sin indicar el endpoint para recibir las notificaciones
	Dado un tópico existente
	Cuando intento suscribirme en modo push a un tópico sin pasar el endpoint de notificaciones
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'endpoint: The field is required.'

Escenario: Suscripción indicando un endpoint sin formato de url
	Dado un tópico existente
	Cuando intento suscribirme en modo push a un tópico con el endpoint sin un formato válidop de url
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'endpoint: must be a valid URL.'

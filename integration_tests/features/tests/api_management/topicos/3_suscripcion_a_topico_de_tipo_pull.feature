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
# 	}


Escenario: Suscripición exitosa a un topico en modo pull
	Dado un tópico existente
	Cuando me suscribo en modo pull al tópico correctamente
	Entonces debo recibir una respuesta de suscripción correta
	Y la respuesta debe tener los valores 'name, endpoint, topic, type'

Escenario: Suscripción a tópico inexistente
	Cuando me suscribo en modo pull a un tópico que no existe
	Entonces debo obtener un status code 404
	Y debo obtener el mensaje de tópico inexistente

Escenario: Suscripción sin datos
	Dado un tópico existente
	Cuando me suscribo en modo pull al tópico sin pasar ningún dato
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'Topic, subscriber and endpoint are required fields'

Escenario: Suscripción sin indicar el nombre del tópico
	Cuando intento suscribirme en modo pull a un tópico sin pasar el nombre del suscriber
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'Name cannot be null'

Escenario: Suscripción sin indicar el nombre de suscriber
	Cuando intento suscribirme en modo pull a un tópico sin pasar el nombre del suscriber
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'Subscriber cannot be null'

Escenario: Suscripción sin indicar el tipo de suscripcion
	Cuando intento suscribirme a un tópico sin pasar el modo de suscripcion
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'Subscription type is required'

Escenario: Suscripción de suscriber con nombre existente
	Dado un tópico existente
	Y un suscriber ya suscripto en modo pull al tópico
	Cuando con el mismo nombre de suscriber intento suscribirme en modo pull
	Entonces debo obtener el mensaje de error de suscriptor existente

Escenario: Suscripción de suscriber en modo pull pasando valor de endpoint
	Dado un tópico existente
	Cuando intento suscribirme en modo pull pasando un endpoint
	Entonces debo obtener el mensaje de error que el vampo endpoint no es permtido para subscripcion pull

Escenario: Suscripción con un endpoint que no responde
	Dado un tópico existente
	Cuando intento suscribirme en modo pull con un endpoint que no responde
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error the endpoint que no responde

Escenario: Suscripción sin indicar el endpoint para recibir las notificaciones
	Cuando intento suscribirme en modo pull a un tópico sin pasar el endpoint de notificaciones
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'Endpoint cannot be null'

Escenario: Suscripción indicando un endpoint sin formato de url
	Dado un tópico existente
	Cuando intento suscribirme en modo pull a un tópico con el endpoint sin un formato válidop de url
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'Endpoint must be a valid URL'

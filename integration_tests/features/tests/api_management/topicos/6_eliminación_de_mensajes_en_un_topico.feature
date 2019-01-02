# language: es
Característica: Eliminación de mensajes en un tópico
	Como una aplicación que es owner de un tópico
	Cuando existen mensajes enviados
	Quiere que sus mensajes puedan ser eliminados


# method: DELETE
# endpoint: /messages
# body = 
# 	{
#     "subscriber": "subscriber_name", 
#     "messages": [{"message_id":"message_id_value", "delete_token":"<delete_token_value>"}]
# 	}

Antecedentes: Topico existente
	Dado un topico determinado

@delete
Escenario: Eliminación de un mensaje de un tópico
	Dado el tópico tiene suscriptores de tipo pull
	Y se envian una serie de 10 mensajes al tópico
	Cuando uno de los suscriptores borra todos los mensajes
	Entonces el tópico no debe poseer mensajes para leer

@delete
Escenario: Sólo el último lector puede borrar un mensaje
	Dado dos suscriptores del mismo nombre en modo pull al tópico
	Y se envian una serie de 10 mensajes al tópico
	Cuando uno de los suscriptores consulta el mensaje sin hacer mas que leerlo
	Y pasan 10 segundos
	* el otro subscriber consulta el mensaje
	Cuando el primer subscritor intenta borrar el mensaje
	Entonces el mensaje no puede ser borrado porque el delete token lo tienen el último lector


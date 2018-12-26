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

Antecedentes: Topico existente
	Dado un topico determinado
	Y el tópico tienen suscriptores de tipo pull
@delete
Escenario: Eliminación de un mensaje de un tópico
	Dado se envian una serie de 10 mensajes al tópico
	Cuando uno de los suscriptores borra todos los mensajes
	Entonces el tópico no debe poseer mensajes para leer
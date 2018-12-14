# language: es
Característica: Consulta de mensajes por suscriptores en modo pull
	Como un usuario suscripto a un tópico en modo pull
	Cuando consulto los mensajes en el tópico
	Quiero ver los mensajes disponibles
	Quien levanta un mensaje con un visibility timeout tendrá dicho tiempo para usar el mensaje
	Y solo el podrá eliminarlo. Pasado ese tiempo y no haber sido borrado, el mensaje estará disponible


# method: GET
# endpoint: messages?subscriber=<subscriber_name>&max_messages=<int>

Antecedentes: Topico existente con suscriptores en modo pull
	Dado un topico determinado
	

Escenario: Consulta de mensajes
	Dado estoy suscripto en modo pull a dicho topico
	Y que se notificó un mensaje a un topico que estoy suscripto en modo pull con visibilidad de 5 segundos
	Cuando consulto los mensajes al tópico
	Entonces debo obtener los mensajes existentes


Escenario: Consulta de mensajes sin indicar quien soy
	Dado estoy suscripto en modo pull a dicho topico
	Cuando consulto los mensajes sin indicar quien soy
	Entonces debo obtener un status code 400
	Y debo obtener el mensaje de error 'max_messages: The field is required; subscriber: The field is required.'

Escenario: Exclusividad de lectura y uso de un mensaje en topic modo pull para un mismo subscriber
 	Dado dos suscriptores del mismo nombre en modo pull al tópico
	Y se envía un mensaje al topico
	Cuando uno de los suscriptores consulta el mensaje sin hacer mas que leerlo
	Entonces el otro suscriber no puede leer el mensaje
	Cuando pasan 5 segundos
	Entonces el otro suscriber puede leer el mensaje

Escenario: Lectura y borrado de un mensaje
 	Dado dos suscriptores del mismo nombre en modo pull al tópico
	Y se envía un mensaje al topico
	Cuando uno de los suscriptores consulta el mensaje y lo borra
	Entonces el otro suscriptor no encontrará el mensaje
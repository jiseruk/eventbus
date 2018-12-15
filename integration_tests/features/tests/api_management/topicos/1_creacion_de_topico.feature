# language: es
@topicos @creacion
Característica: Creación de tópicos
	Una aplicación que desea comunicar novedades
	Cuando quiere informar determinados mensajes
	debe poder crear un tópico

# method: POST

# endpoint: /topics
# 	body = {
# 		"topic": "topic_name", 
# 		"engine":"AWS,
# 		"description":"Some description",
#     "owner":"The owner of the topic"
# 	}

Escenario: Creación correcta de un tópico
	Dado que deseo generar un tópico nuevo
	Cuando ingreso los datos requeridos del tópico
	Entonces debo obtener una respuesta de aceptado
	Y debo recibir los datos de fecha de creacion, nombre de topico, token de seguridad, owner y descripcion

Escenario: Creación de un tópico existente
	Dado un tópico existente en el event bus
	Cuando intento crear un tópico con el nombre del existente
	Entonces debo obtener una respuesta de tópico existente

Escenario: Creación sin pasar nombre de topico
	Dado que deseo generar un tópico nuevo
	Cuando ingreso los datos sin pasar el nombre
	Entonces debo obtener una respuesa que indique que el nombre es obligatorio

Escenario: Creación sin pasar el engine
	Dado que deseo generar un tópico nuevo
	Cuando ingreso los datos sin pasar el engine
	Entonces debo obtener una respuesa que indique que el engine es obligatorio

Escenario: Creación indicando un engine inexistente
	Dado que deseo generar un tópico nuevo
	Cuando ingreso los datos con un engine que no existe
	Entonces debo obtener un status code 400
	Y debo obtener una respuesa que indique que el engine no existe

Escenario: Creación sin pasar descricion
	Dado que deseo generar un tópico nuevo
	Cuando ingreso los datos sin pasar la descripcion
	Entonces debo obtener una respuesa que indique que la descripcion es obligatoria

Escenario: Creacion sin pasar el owner del tópico
	Dado que deseo generar un tópico nuevo
	Cuando ingreso los datos sin el owner del topico
	Entonces debo obtener una respuesa que indique que el owner del topico es obligatorio

Escenario: Creacion sin pasar ningún dato
	Dado que deseo generar un tópico nuevo
	Cuando ingreso la solicitud sin pasar ningún dato
	Entonces debo obtener una respuesa que indique que el nombre, el nombre, el engine, la descripcion y el owner son obligatorio
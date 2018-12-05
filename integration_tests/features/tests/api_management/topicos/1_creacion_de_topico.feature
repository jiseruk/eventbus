# language: es
@topicos @creacion
Característica: Creación de tópicos
	Una aplicación que desea comunicar novedades
	Cuando quiere informar determinados eventos
	debe poder crear un tópico

# method: POST

# endpoint: /topics
# 	body = {
# 		"topic": "topic_name", 
# 		"engine":"AWS
# 	}

Escenario: Creación correcta de un tópico
	Dado que deseo generar un tópico nuevo
	Cuando ingreso los datos requeridos del tópico
	Entonces debo obtener una respuesta de aceptado
	Y debo recibir los datos de id, fecha de creacion y nombre de topico

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


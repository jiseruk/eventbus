# language: es
@topicos @subscribers @listado
Característica: Listado de tópicos y subscribers
	Como un usuarios del event bus
	Cuando solicito el listao de topicos y subscribers
	Quiero obtener la lista de topicos existentes y sus subscribers

# method: GET

# endpoint: /topics

Escenario: Listado de topicos existentes
	Dado que existen topicos creados en el even bus
	Cuando solicito el listado de tópicos
	Entonces debe aparecer la lista de topicos existentes

Escenario: Subscribers de tipo push de un topico
	Dado un tópico existente
	Y un suscriber ya suscripto en modo push al tópico
	Cuando consulto los subscribers del topico
	Entonces debe aparecer el subscriber en la lista del topico
	
Escenario: Subscribers de tipo pull de un topico
	Dado un tópico existente
	Y un suscriber ya suscripto en modo pull al tópico
	Cuando consulto los subscribers del topico
	Entonces debe aparecer el subscriber en la lista del topico

Escenario: Información de un subscriber de tipo pull de un topico
	Dado un tópico existente
	Y un suscriber ya suscripto en modo push al tópico
	Cuando consulto por dicho subscriber en el tópicop
	Entonces debe aparecer el subscriber y la información de name, endpoint, topic, type

Escenario: Información de un subscriber de tipo push de un topico
	Dado un tópico existente
	Y un suscriber ya suscripto en modo pull al tópico
	Cuando consulto por dicho subscriber en el tópicop
	Entonces debe aparecer el subscriber y la información de name, endpoint, topic, type
#########
# Copien este contenido en un archivo llamando 'Dockerfile'. Donde dejen este Dockerfile ejecuten:
#
# 1- buildearlo 
#       $ docker build -t depnot .
# ('imagen_garba-automation' puede ser cualquier nombre que quieran)
#
# 2- Una vez buildeado, pueden ver la imagen con: 
#       $ docker images
#
# 3- Para correr la imagen y entrar al container es: 
#       $ docker run -it -P depnot /bin/bash
#########
FROM ubuntu
FROM ruby:2.4.0

EXPOSE 9292

MAINTAINER Equipo QA Wenance

#Update aptitude with new repo
RUN apt-get update

RUN export TERM=dumb; bash

# update
RUN apt-get update
RUN apt-get install -y 	apt-transport-https

# add gemserver
RUN gem source -a http://gems.fintechpeople.io:9292

#install bundler
RUN gem install bundler

#Clone repo
ADD . /suscriber

WORKDIR /suscriber

RUN bundle install

CMD rackup -p 9292 -o 0.0.0.0
# CMD ["rackup", "-p", "9292"]

# docker run -ti --network wequeue-management_api_default suscriber
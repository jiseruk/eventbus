FROM ubuntu
FROM ruby:2.4.0

MAINTAINER Equipo QA Wenance

# To see colors on terminal
ENV ENVIRONMENT $ENVIRONMENT
RUN export TERM=dumb; bash
RUN apt-get update
RUN apt-get update
RUN gem install bundler
ADD . /integration_tests
WORKDIR /integration_tests
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh ./
RUN chmod +x ./wait-for-it.sh
RUN bundle install

CMD ["./wait-for-it.sh", "localstack:4565", "wequeue:8080", "--", "cucumber"]

# 
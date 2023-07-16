FROM ubuntu:22.04

WORKDIR /service
COPY ./notifications/bin/app .
COPY ./notifications/config.yaml .
COPY ./certs /etc/ssl/certs/

CMD ["./app"]

FROM golang:1.19.3-alpine3.16

EXPOSE 9000

RUN apk update \
  && apk add --no-cache \ 
    build-base
  
RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
COPY ./my_entrypoint.sh /usr/local/bin/my_entrypoint.sh
RUN /bin/chmod +x /usr/local/bin/my_entrypoint.sh

RUN go build cmd/main.go
RUN mv main /usr/local/bin/

CMD ["main"]
ENTRYPOINT ["my_entrypoint.sh"]

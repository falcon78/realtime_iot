FROM golang:1.16-alpine

ENV GO111MODULE=on


RUN apk add alpine-sdk bash git --no-cache \
    && git clone https://github.com/vishnubob/wait-for-it.git /tmp/wait-for-it \
    && mv /tmp/wait-for-it/wait-for-it.sh /usr/local/bin/

COPY . /realtime_iot

WORKDIR /realtime_iot/cmd/server

RUN go get -u github.com/cosmtrek/air

CMD tail -f /dev/null


EXPOSE 31415

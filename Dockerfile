FROM golang:latest AS builder
WORKDIR /home/container
ADD . .
RUN go get -d ./... && \
    go build -o discordimap -i *.go
CMD [ "./discordimap" ]

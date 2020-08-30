FROM golang:latest AS builder
MAINTAINER itsmeadityaagarwal@gmail.com

WORKDIR /go/src/temp
COPY go.mod go.sum ./
COPY . .
RUN go build -v -trimpath -o tempcontrol .
RUN chmod +x tempcontrol
FROM golang:latest
WORKDIR /temp
COPY --from=builder /go/src/temp/tempcontrol .
COPY --from=builder /go/src/temp/startup.sh .

ENTRYPOINT ["./tempcontrol"]


FROM golang:1.22.1

WORKDIR /app

COPY . . 

ENV GO111MODULE='on'
ENV CGO_ENABLED='0'
ENV GOPROXY='https://goproxy.cn,direct'
RUN go mod download
RUN go build -o cloudman cmd/*.go
RUN mv cloudman /app/cloudman
ENV PROTOCOL="http://"
ENV SUB_URL_PREFIX="http://154.13.5.159:9999"
ENV SUB_PORT="2096"
ENV PORT="8088"
ENV MODE="prod"
ENTRYPOINT ["/app/cloudman"]
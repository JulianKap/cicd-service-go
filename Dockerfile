FROM golang:1.22.1 AS builder

#ENV GOPROXY=$GOPROXY
#ENV GO111MODULE=on
#ENV GOSUMDB=off
#
#RUN apt-get update && apt-get install -y \
#    git \
#    openssh-server \
#    ca-certificates \
#    wget

WORKDIR /app

COPY go.* ./

RUN go mod download && \
    go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -o cicd-service-go ./cmd/app

FROM alpine/curl:3.14

WORKDIR /app

COPY --from=builder /app/cicd-service-go /app/cicd-service-go

CMD [ "/app/cicd-service-go" ]
FROM golang:1.22.1 AS builder

WORKDIR /app

COPY go.* ./

RUN go mod download && \
    go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -o cicd-service-go ./cmd/app

FROM docker:26.1.1-dind-alpine3.19

WORKDIR /app

COPY --from=builder /app/cicd-service-go /app/cicd-service-go

CMD [ "/app/cicd-service-go" ]
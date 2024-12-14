FROM golang:1.23.2-alpine3.20

WORKDIR /

RUN ls -la

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

CMD ["go", "run", "/cmd/service/main.go"]
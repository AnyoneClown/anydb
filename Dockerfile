FROM golang:alpine

WORKDIR /app

RUN mkdir -p /root/.anydb && \
    chmod 755 /root/.anydb

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o anydb

EXPOSE 8080

CMD ["./anydb", "run"]

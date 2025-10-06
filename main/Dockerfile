FROM golang:1.24.5-bookworm

WORKDIR /app/

COPY . /app/

RUN go build -tags netgo -ldflags '-s -w' -o app

EXPOSE 8000

ENTRYPOINT [ "./app" ]
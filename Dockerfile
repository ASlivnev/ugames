FROM golang:latest
ENV GO111MODULE=on

RUN mkdir /app
ADD . /app/
WORKDIR /app

COPY cmd/app/ .
RUN go build -o main .

EXPOSE 8088
CMD ["/app/main"]
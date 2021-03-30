FROM golang:1.12-alpine

RUN apk add --no-cache git

WORKDIR .

COPY . .

RUN go get -d -v github.com/gorilla/mux go.mongodb.org/mongo-driver/bson/primitive go.mongodb.org/mongo-driver/mongo go.mongodb.org/mongo-driver/mongo/options gopkg.in/gomail.v2

RUN go build -o main-app .

EXPOSE 8000

CMD ["./main-app"]
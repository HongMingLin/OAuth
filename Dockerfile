FROM golang:1.14
WORKDIR /login
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
EXPOSE 8080
ENTRYPOINT ./login
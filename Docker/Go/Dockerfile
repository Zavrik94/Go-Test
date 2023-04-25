FROM golang:latest

RUN mkdir /app
WORKDIR /app
ADD . /app

RUN go get github.com/githubnemo/CompileDaemon && \
 go install github.com/githubnemo/CompileDaemon

RUN go get github.com/gin-gonic/gin && \
 go install github.com/gin-gonic/gin



ENV PATH="${PATH}:${GOPATH}/bin"

ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main
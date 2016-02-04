FROM golang:1.5.3

ADD . /go/src/shaleapps_todo
RUN go get github.com/mattn/go-sqlite3
RUN go install shaleapps_todo

ENTRYPOINT /go/bin/shaleapps_todo

EXPOSE 8080

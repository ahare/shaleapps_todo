# ShaleApps ToDo Code Sample

This is a simple Go HTTP server that exposes an API for managing a ToDo list.

### Dependencies

The only external dependency is https://github.com/mattn/go-sqlite3.

### Design thought process

I intentionally didn't use an ORM. I agree
with the logic in 
[Golang, ORMs, and why I am still not using one](http://www.hydrogen18.com/blog/golang-orms-and-why-im-still-not-using-one.html)
and believe that ORMs don't fit well with statically typed languages and the benefit
they offer doesn't outweigh the performance hit or the black-box they create in
the app.

I wrote a very simple router (`server/router.go`) to handle routes by path and
HTTP method. Normally, I would have used [httprouter](https://github.com/julienschmidt/httprouter)
but I didn't because I wasn't sure if that would be considered cheating on the
first rule. :) 

### Building

```
$ git clone git@github.com:ahare/shaleapps_todo.git
$ cd shaleapps_todo
$ go get && go build
```

### Running the tests

```
$ go test ./...
?     shaleapps_todo  [no test files]
ok    shaleapps_todo/db 0.018s
ok    shaleapps_todo/server 0.009s
```

### Running the server

```
$ ./shaleapps_todo
ToDo server listening on :8080...
```

### Using the server

Create a new ToDo:

```
$ curl -XPOST -d '{"text": "Pick up milk"}' -H 'Accept: application/json' -i ':8080/todos'
HTTP/1.1 201 Created
Content-Type: application/json
Date: Wed, 03 Feb 2016 18:20:24 GMT
Content-Length: 44

{"id":1,"text":"Pick up milk","done":false}
```

Update a ToDo:

```
$ curl -XPUT -d '{"id":1,"text":"Pick up milk","done":true}' -H 'Accept: application/json' -i ':8080/todos/1'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 03 Feb 2016 18:46:59 GMT
Content-Length: 43

{"id":1,"text":"Pick up milk","done":true}
```

Add several more ToDos:

```
$ curl -XPOST -d '{"text": "Drink milk"}' -H 'Accept: application/json' ':8080/todos'
{"id":2,"text":"Drink milk","done":false}
$ curl -XPOST -d '{"text": "Throw away milk jug"}' -H 'Accept: application/json' ':8080/todos'
{"id":3,"text":"Throw away milk jug","done":false}
$ curl -XPOST -d '{"text": "Wash the car"}' -H 'Accept: application/json' ':8080/todos'
{"id":4,"text":"Wash the car","done":false}
```

Find a ToDo by ID:

```
$ curl -i ':8080/todos/4'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 03 Feb 2016 19:38:14 GMT
Content-Length: 44

{"id":4,"text":"Wash the car","done":false}
```

Find all ToDos that are done:

```
$ curl -i ':8080/todos?done=true'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 03 Feb 2016 19:39:10 GMT
Content-Length: 45

[{"id":1,"text":"Pick up milk","done":true}]
```

Search for ToDos by text: 

```
$ curl -i ':8080/todos?text=milk'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 03 Feb 2016 19:32:44 GMT
Content-Length: 138

[{"id":1,"text":"Pick up milk","done":true},{"id":2,"text":"Drink milk","done":false},{"id":3,"text":"Throw away milk jug","done":false}]
```

Get all ToDos:

```
$ curl -i ':8080/todos?text=milk'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 03 Feb 2016 19:32:47 GMT
Content-Length: 138

[{"id":1,"text":"Pick up milk","done":true},{"id":2,"text":"Drink milk","done":false},{"id":3,"text":"Throw away milk jug","done":false},{"id":4,"text":"Wash the car","done":false}]
```

Delete a ToDo:

```
$ curl -XDELETE -i ':8080/todos/1'
HTTP/1.1 200 OK
Date: Wed, 03 Feb 2016 19:01:18 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8
```

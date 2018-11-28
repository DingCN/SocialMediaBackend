## Use
run
~~~~
go run web.go
~~~~
under directory cmd/local/web/  
run
~~~~
go run backend.go
~~~~
under directory cmd/local/backend/  
web server runs at localhost:8080 by default, backend server runs at localhost:50051 by default
## Test
run
~~~~
go run web_test.go
~~~~
under directory cmd/local/web/
run
~~~~
go run backend_test.go
~~~~
under directory cmd/local/backend/
run
~~~~
go run storage_test.go
~~~~
under directory cmd/local/backend/


## Structure
  When a request comes, it is handled by "pkg/web/web.go";
  
  "web.go" then calls "rpcsend.go" to send rpc to backend server. Backend server receives rpc in "backend.go", it then calls "storage API" to communicate with package storage. 
      
## Note
Server uses memory as storage in Part 1, but out of memory is not handled.  
TODO test: threadsafe

## Reference
  Based on Adam's github code and the [tutorial](astaxie.gitbooks.io/build-web-application-with-golang) 

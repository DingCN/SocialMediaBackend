## Use
run
~~~~
go run web.go
~~~~
under directory cmd/local/web/
  
server runs at localhost:8080 by default

## Structure
  When a request comes, it is handled by "pkg/web/web.go";
  
  "web.go" then calls "operator.go" to get or modify our data structures, which are stored in "model.go".
      
## Note
Server uses memory as storage in Part 1, but out of memory is not handled, since we are using database eventually 

## Reference
  Based on Adam's github code and the [tutorial](astaxie.gitbooks.io/build-web-application-with-golang) 

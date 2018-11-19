## Use
Test file is at pkg/web/web_test.go (all test cleared)
  
Run server by executing cmd/local/web/web.go
  
There are still some bug on frontend in this release, we are still working on it.

  ## Structure
  When a request comes, it is handled by "pkg/web/web.go";
  
  "web.go" then calls "operator.go" to get or modify our data structures, which are stored in "model.go".
      
  ## Reference
  Based on Adam's github code and the [tutorial](astaxie.gitbooks.io/build-web-application-with-golang) 

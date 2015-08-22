package main

import (
    "flag"
    "html/template"
    "log"
    "net/http"
)




func main() {
  http.handleFunc('/', homePageFunc)
  http.handleFunc('/static/', staticHandler)
}

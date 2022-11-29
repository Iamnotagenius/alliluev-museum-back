package main

import (
	"fmt"
	"net/http"
    "github.com/graphql-go/handler"
)

func main() {
    schema, err := init_schema()
    if err != nil {
        fmt.Println(err)
    }
    h := handler.New(&handler.Config{
        Schema: &schema,
        Pretty: true,
        GraphiQL: true,
    })

    http.Handle("/graphql", h)
    err = http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println(err)
    }
}

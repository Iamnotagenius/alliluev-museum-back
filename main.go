package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/graphql-go/handler"
)

func main() {
	schema, err := init_schema()
	if err != nil {
		fmt.Println(err)
	}

	db := NewDB("172.17.0.2:3306")

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", httpDBMiddleware(h, db))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func httpDBMiddleware(next *handler.Handler, db MuseumDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "db", db)
		next.ContextHandler(ctx, w, r)
	})
}

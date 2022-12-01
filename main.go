package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/handler"
)

func main() {
	address := flag.String("database-address", "172.17.0.2:3306", "Host address of database")
	port := flag.String("port", "8080", "Port to serve the api")
	flag.Parse()

	if os.Getenv(UserEnv) == "" {
		log.Fatalln("MYSQL_USER is not set")
	}

	if os.Getenv(PasswdEnv) == "" {
		log.Fatalln("MYSQL_PASSWORD is not set")
	}

	schema, err := init_schema()
	if err != nil {
		fmt.Println(err)
	}

	db := NewDB(*address)

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", httpDBMiddleware(h, db))
	err = http.ListenAndServe(":"+*port, nil)
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

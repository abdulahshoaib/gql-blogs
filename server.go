package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/abdulahshoaib/gql-blogs/graph"
	internalModel "github.com/abdulahshoaib/gql-blogs/internal/model"
	"github.com/abdulahshoaib/gql-blogs/internal/database"
	"github.com/vektah/gqlparser/v2/ast"
	"gorm.io/gorm"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := database.New()
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&internalModel.User{}, &internalModel.Post{}, &internalModel.Comment{})
	if err != nil {
		log.Fatal(err)
	}

	seedUsers(db)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(db)}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func seedUsers(db *gorm.DB) {
	var count int64
	db.Model(&internalModel.User{}).Count(&count)
	if count > 0 {
		return
	}
	users := []internalModel.User{
		{Name: "Alice"},
		{Name: "Bob"},
	}
	for _, u := range users {
		db.Create(&u)
	}
}

package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	resolvers "my_app/internal/server/graphql"
	"time"

	"my_app/internal/config"
	gfql "my_app/internal/graph"
	"my_app/internal/logger"
	"net/http"
	"os"
)

func main() {
	lgr := logger.InitLogger()

	// проверка на наличие файлов конфига
	// если более двух то используется последний
	envfile := ".env"
	if len(os.Args) >= 2 {
		envfile = os.Args[1]
	}

	lgr.Info.Printf("Инициализация конфигурации.\nЧтение файла %s", envfile)
	if err := config.Init(envfile); err != nil {
		lgr.Err.Fatal(err.Error())
	}

	port := os.Getenv("PORT")

	srv := handler.New(gfql.NewExecutableSchema(gfql.Config{Resolvers: &resolvers.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	// добавляем поддержку WebSocket
	srv.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	// добавляем поддержку других транспортов
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	lgr.Info.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	lgr.Err.Fatal(http.ListenAndServe(":"+port, nil))
}

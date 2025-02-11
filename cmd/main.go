package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"my_app/internal/consts"
	"my_app/internal/database"
	"my_app/internal/gateway"
	"my_app/internal/gateway/postgres"
	resolvers "my_app/internal/server/graphql"
	"my_app/internal/service"
	"time"

	"my_app/internal/config"
	"my_app/internal/gateway/in_memory"
	gfql "my_app/internal/graph"
	"my_app/internal/logger"
	"my_app/internal/server/graphql"
	"net/http"
	"os"
)

func main() {
	var gateways *gateway.Gateways

	lgr := logger.InitLogger()

	// проверка на наличие файлов конфига
	// если более двух то используется последний
	envFile := ".env"
	if len(os.Args) >= 2 {
		envFile = os.Args[1]
	}

	lgr.Info.Printf("Инициализация конфигурации.\nЧтение файла %s", envFile)
	if err := config.Init(envFile); err != nil {
		lgr.Err.Fatal(err.Error())
	}

	lgr.Info.Print("Подключение к Postgres.")
	options := database.PostgresInit()
	lgr.Info.Printf("%s %s %s %s %s", options.Name, options.User, options.Password, options.Port, options.Host)
	pgDb, err := database.NewPostgresDB(*options)
	if err != nil {
		lgr.Err.Fatal(err.Error())
	}

	lgr.Info.Print("Создание Шлюза.")
	lgr.Info.Print("USE_IN_MEMORY = ", os.Getenv("USE_IN_MEMORY"))
	if os.Getenv("USE_IN_MEMORY") == "true" {
		posts := in_memory.NewPostsInMemory(consts.PostsPullSize)
		comments := in_memory.NewCommentsInMemory(consts.CommentsPullSize)
		gateways = gateway.NewGateways(posts, comments)
	} else {
		posts := postgres.NewPostsPostgres(pgDb)
		comments := postgres.NewCommentsPostgres(pgDb)
		gateways = gateway.NewGateways(posts, comments)
	}
	lgr.Info.Print("Creating Services.")
	services := service.NewServices(gateways, lgr)

	port := os.Getenv("PORT")
	srv := handler.New(gfql.NewExecutableSchema(gfql.Config{Resolvers: &resolvers.Resolver{
		PostsService:      services.Posts,
		CommentsService:   services.Comments,
		CommentsObservers: graphql.NewCommentsObserver(),
	}}))

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

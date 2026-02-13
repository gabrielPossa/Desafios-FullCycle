package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/gabrielPossa/Desafios-FullCycle/cleanArch/configs"
	"github.com/gabrielPossa/Desafios-FullCycle/cleanArch/internal/event/handler"
	"github.com/gabrielPossa/Desafios-FullCycle/cleanArch/internal/infra/graph"
	"github.com/gabrielPossa/Desafios-FullCycle/cleanArch/internal/infra/grpc/pb"
	"github.com/gabrielPossa/Desafios-FullCycle/cleanArch/internal/infra/grpc/service"
	"github.com/gabrielPossa/Desafios-FullCycle/cleanArch/internal/infra/web/webserver"
	"github.com/gabrielPossa/Desafios-FullCycle/cleanArch/pkg/events"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	loadedConfigs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(loadedConfigs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", loadedConfigs.DBUser, loadedConfigs.DBPassword, loadedConfigs.DBHost, loadedConfigs.DBPort, loadedConfigs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(loadedConfigs.RabbitMQ)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)

	httpServer := webserver.NewWebServer(loadedConfigs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	httpServer.AddHandler("/order", webserver.POST, webOrderHandler.Create)
	httpServer.AddHandler("/order", webserver.GET, webOrderHandler.RetrieveAll)
	fmt.Println("Starting web server on port", loadedConfigs.WebServerPort)
	go httpServer.Start()

	grpcServer := grpc.NewServer()
	OrderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, OrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", loadedConfigs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", loadedConfigs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", loadedConfigs.GraphQLServerPort)
	http.ListenAndServe(":"+loadedConfigs.GraphQLServerPort, nil)
}

func getRabbitMQChannel(c configs.RabbitMQConfig) *amqp.Channel {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", c.User, c.Pass, c.Host, c.Port))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

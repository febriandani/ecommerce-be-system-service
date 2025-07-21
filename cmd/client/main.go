package main

import (
	"fmt"
	"log"

	"github.com/febriandani/ecommerce-be-system-service/cmd/server/handler"
	clientNotification "github.com/febriandani/ecommerce-be-system-service/protogen/golang/proto/system"
	"github.com/spf13/viper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	viper.SetConfigName("/cmd/server/config/app")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", viper.GetString("APP.URL"), viper.GetInt("APP.PORT_SERVER")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error listen server from client : %v", err)
		panic(err)
	}
	client := clientNotification.NewSystemsClient(conn)

	handler := handler.NewHandler(client)
	app := fiber.New()
	app.Use(logger.New())

	system := app.Group("/system",
		basicauth.New(basicauth.Config{
			Users: map[string]string{
				"admin": "secret123",
			},
		}))
	system.Get("/provinces", handler.HandlerSystem_GetProvinces)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", viper.GetInt("APP.PORT_CLIENT"))))
}

package main

import (
	"os"

	"example.com/order-service/db"
	"example.com/order-service/rmq"
	"example.com/order-service/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load .env
	godotenv.Load(".env")

	// set GIN_MODE
	gin.SetMode(os.Getenv("GIN_MODE"))

	// connect to postgres DB
	db.ConnectDB()

	// connect to rabbitmq (rmq folder)
	go rmq.ConnectToRMQ()

	r := gin.Default()

	//localhost:6000/api/v1
	apiV1Route := r.Group("/api/v1")

	//Home Route
	routes.InitHomeRoute(apiV1Route)

	//Order Route
	routes.InitOrderRoute(apiV1Route)

	//allow CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
		// AllowCredentials: true,
	}))

	r.Run(":" + os.Getenv("PORT")) // localhost:6000/api/v1

}

package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/FilipPauco/hospital-webapi/api"
	"github.com/FilipPauco/hospital-webapi/internal/db_service"
	"github.com/FilipPauco/hospital-webapi/internal/hospital_wl"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("HOSPITAL_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("HOSPITAL_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") {
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
	engine.Use(corsMiddleware)

	dbService := db_service.NewMongoService[hospital_wl.Ward](db_service.MongoServiceConfig{})
	defer dbService.Disconnect(context.Background())
	engine.Use(func(ctx *gin.Context) {
		ctx.Set("db_service", dbService)
		ctx.Next()
	})

	handleFunctions := &hospital_wl.ApiHandleFunctions{
		VisitsAPI: hospital_wl.NewVisitsApi(),
		BedsAPI:   hospital_wl.NewBedsApi(),
	}
	hospital_wl.NewRouterWithGinEngine(engine, *handleFunctions)
	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}

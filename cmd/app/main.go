package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "gitlab.com/Nikolay-Yakunin/blog-service/docs"

	"fmt"
)

// @title Blog Service API
// @version 1.0
// @description API сервер для блог-платформы
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @schemes http https
func main() {
    r := gin.Default()
    
    // Load templates
    r.LoadHTMLGlob("templates/*")
    
    // Serve static files
    r.Static("/static", "./static")
    
    // Setup swagger
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    // Routes
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{
            "title": "Blog Service",
        })
    })

	fmt.Print("Go to => 'http://localhost:8080/swagger/index.html' to see the Swagger documentation\n")

    r.Run(":8080")
}

package main

import (
	"github.com/gin-gonic/gin"
	"go-test/routes"
)

func main() {
	r := gin.Default()
	routes.RegisterAPIRoutes(r)
	routes.RegisterAuthRoutes(r)
	routes.RegisterCarsRoutes(r)
	r.Run() // listen and serve on 0.0.0.0:8080
}

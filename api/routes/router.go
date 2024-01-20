package routes

import (
	"github.com/chrpa-jakub/register-api/controller"
	"github.com/gin-gonic/gin"
)

func Init(){
  r := gin.Default()
  r.POST("/api/register", controller.Register)
  r.POST("/api/login", controller.Login)

  r.Run()
}

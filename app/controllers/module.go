package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewUserController),

	fx.Invoke(RegisterUserRoutes),
)

func RegisterUserRoutes(engine *gin.Engine, ctl *UserController) {
	users := engine.Group("/api/v1/users")
	{
		users.POST("", ctl.Create)
		users.GET("", ctl.List)
		users.GET("/:id", ctl.GetByID)
		users.PUT("/:id", ctl.Update)
		users.DELETE("/:id", ctl.Delete)
	}
}

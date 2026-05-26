package controllers

import (
	"net/http"
	"strconv"

	userdto "golang-project-structure/internal/dto/user"
	"golang-project-structure/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	UserService services.IUserService
	Logger      *zap.Logger
}

func NewUserController(userService services.IUserService, logger *zap.Logger) *UserController {
	return &UserController{
		UserService: userService,
		Logger:      logger,
	}
}

func (ctl *UserController) Create(c *gin.Context) {
	var params userdto.CreateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, code, err := ctl.UserService.Create(c.Request.Context(), params)
	if err != nil {
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	c.JSON(code, gin.H{"message": "success", "data": user})
}

func (ctl *UserController) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	user, code, err := ctl.UserService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	c.JSON(code, gin.H{"message": "success", "data": user})
}

func (ctl *UserController) List(c *gin.Context) {
	var params userdto.ListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result, code, err := ctl.UserService.List(c.Request.Context(), params)
	if err != nil {
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	c.JSON(code, gin.H{"message": "success", "data": result.Items, "meta": result.Meta})
}

func (ctl *UserController) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	var params userdto.UpdateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, code, err := ctl.UserService.Update(c.Request.Context(), id, params)
	if err != nil {
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	c.JSON(code, gin.H{"message": "success", "data": user})
}

func (ctl *UserController) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	code, err := ctl.UserService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(code, gin.H{"message": err.Error()})
		return
	}
	c.JSON(code, gin.H{"message": "success"})
}

func parseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

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

// Create godoc
// @Summary  Create a user
// @Tags     users
// @Accept   json
// @Produce  json
// @Param    user  body      userdto.CreateParams  true  "User payload"
// @Success  201   {object}  userdto.UserDTO
// @Failure  400   {object}  map[string]string
// @Router   /users [post]
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

// GetByID godoc
// @Summary  Get a user by ID
// @Tags     users
// @Produce  json
// @Param    id   path      int  true  "User ID"
// @Success  200  {object}  userdto.UserDTO
// @Failure  404  {object}  map[string]string
// @Router   /users/{id} [get]
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

// List godoc
// @Summary  List users
// @Tags     users
// @Produce  json
// @Param    page   query     int  false  "Page number"
// @Param    limit  query     int  false  "Items per page"
// @Success  200    {object}  userdto.ListOutput
// @Router   /users [get]
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

// Update godoc
// @Summary  Update a user
// @Tags     users
// @Accept   json
// @Produce  json
// @Param    id    path      int                true  "User ID"
// @Param    user  body      userdto.UpdateParams  true  "User payload"
// @Success  200   {object}  userdto.UserDTO
// @Failure  400   {object}  map[string]string
// @Failure  404   {object}  map[string]string
// @Router   /users/{id} [put]
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

// Delete godoc
// @Summary  Delete a user
// @Tags     users
// @Produce  json
// @Param    id   path      int  true  "User ID"
// @Success  200  {object}  map[string]string
// @Failure  404  {object}  map[string]string
// @Router   /users/{id} [delete]
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

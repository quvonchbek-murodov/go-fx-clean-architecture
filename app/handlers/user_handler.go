package handlers

import (
	"context"

	userpb "golang-project-structure/genprotos/user"
	userdto "golang-project-structure/internal/dto/user"
	"golang-project-structure/internal/services"

	"go.uber.org/zap"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	UserService services.IUserService
	Logger      *zap.Logger
}

func NewUserHandler(userService services.IUserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		UserService: userService,
		Logger:      logger,
	}
}

func (h *UserHandler) Create(ctx context.Context, req *userpb.CreateRequest) (*userpb.UserResponse, error) {
	user, code, err := h.UserService.Create(ctx, userdto.CreateParams{Name: req.GetName(), Email: req.GetEmail()})
	if err != nil {
		h.Logger.Error("create user", zap.Error(err))
		return &userpb.UserResponse{StatusCode: int32(code), Message: err.Error()}, nil
	}
	return &userpb.UserResponse{StatusCode: int32(code), Message: "success", Data: toProto(user)}, nil
}

func (h *UserHandler) GetByID(ctx context.Context, req *userpb.GetByIDRequest) (*userpb.UserResponse, error) {
	user, code, err := h.UserService.GetByID(ctx, uint(req.GetId()))
	if err != nil {
		return &userpb.UserResponse{StatusCode: int32(code), Message: err.Error()}, nil
	}
	return &userpb.UserResponse{StatusCode: int32(code), Message: "success", Data: toProto(user)}, nil
}

func (h *UserHandler) List(ctx context.Context, req *userpb.ListRequest) (*userpb.ListResponse, error) {
	result, code, err := h.UserService.List(ctx, userdto.ListParams{Page: int(req.GetPage()), Limit: int(req.GetLimit())})
	if err != nil {
		h.Logger.Error("list users", zap.Error(err))
		return &userpb.ListResponse{StatusCode: int32(code), Message: err.Error()}, nil
	}

	data := make([]*userpb.User, 0, len(result.Items))
	for _, u := range result.Items {
		data = append(data, toProto(u))
	}

	return &userpb.ListResponse{
		StatusCode: int32(code),
		Message:    "success",
		Data:       data,
		Meta: &userpb.Meta{
			ItemsPerPage: int32(result.Meta.ItemsPerPage),
			TotalItems:   result.Meta.TotalItems,
			CurrentPage:  int32(result.Meta.CurrentPage),
			LastPage:     int32(result.Meta.LastPage),
		},
	}, nil
}

func (h *UserHandler) Update(ctx context.Context, req *userpb.UpdateRequest) (*userpb.UserResponse, error) {
	user, code, err := h.UserService.Update(ctx, uint(req.GetId()), userdto.UpdateParams{Name: req.GetName(), Email: req.GetEmail()})
	if err != nil {
		return &userpb.UserResponse{StatusCode: int32(code), Message: err.Error()}, nil
	}
	return &userpb.UserResponse{StatusCode: int32(code), Message: "success", Data: toProto(user)}, nil
}

func (h *UserHandler) Delete(ctx context.Context, req *userpb.DeleteRequest) (*userpb.DeleteResponse, error) {
	code, err := h.UserService.Delete(ctx, uint(req.GetId()))
	if err != nil {
		return &userpb.DeleteResponse{StatusCode: int32(code), Message: err.Error()}, nil
	}
	return &userpb.DeleteResponse{StatusCode: int32(code), Message: "success"}, nil
}

func toProto(u *userdto.UserDTO) *userpb.User {
	return &userpb.User{
		Id:        uint64(u.ID),
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

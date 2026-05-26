package services

import (
	"context"
	"errors"
	"net/http"

	"golang-project-structure/internal/dto"
	userdto "golang-project-structure/internal/dto/user"
	"golang-project-structure/internal/models"
	"golang-project-structure/internal/repositories"

	"gorm.io/gorm"
)

type IUserService interface {
	Create(ctx context.Context, params userdto.CreateParams) (*userdto.UserDTO, int, error)
	GetByID(ctx context.Context, id uint) (*userdto.UserDTO, int, error)
	List(ctx context.Context, params userdto.ListParams) (*userdto.ListOutput, int, error)
	Update(ctx context.Context, id uint, params userdto.UpdateParams) (*userdto.UserDTO, int, error)
	Delete(ctx context.Context, id uint) (int, error)
}

type UserService struct {
	UserRepo repositories.IUserRepository
}

func NewUserService(userRepo repositories.IUserRepository) IUserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) Create(ctx context.Context, params userdto.CreateParams) (*userdto.UserDTO, int, error) {
	u := &models.UserModel{Name: params.Name, Email: params.Email}
	if err := s.UserRepo.Create(ctx, u); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return toDTO(u), http.StatusCreated, nil
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*userdto.UserDTO, int, error) {
	u, err := s.UserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.New("user not found")
		}
		return nil, http.StatusInternalServerError, err
	}
	return toDTO(u), http.StatusOK, nil
}

func (s *UserService) List(ctx context.Context, params userdto.ListParams) (*userdto.ListOutput, int, error) {
	params.Normalize()

	users, total, err := s.UserRepo.List(ctx, params)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	items := make([]*userdto.UserDTO, 0, len(users))
	for _, u := range users {
		items = append(items, toDTO(u))
	}

	lastPage := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &userdto.ListOutput{
		Items: items,
		Meta: dto.MetaDTO{
			ItemsPerPage: params.Limit,
			TotalItems:   total,
			CurrentPage:  params.Page,
			LastPage:     lastPage,
		},
	}, http.StatusOK, nil
}

func (s *UserService) Update(ctx context.Context, id uint, params userdto.UpdateParams) (*userdto.UserDTO, int, error) {
	u, err := s.UserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.New("user not found")
		}
		return nil, http.StatusInternalServerError, err
	}

	u.Name = params.Name
	u.Email = params.Email
	if err := s.UserRepo.Update(ctx, u); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return toDTO(u), http.StatusOK, nil
}

func (s *UserService) Delete(ctx context.Context, id uint) (int, error) {
	if _, err := s.UserRepo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, errors.New("user not found")
		}
		return http.StatusInternalServerError, err
	}
	if err := s.UserRepo.Delete(ctx, id); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func toDTO(u *models.UserModel) *userdto.UserDTO {
	return &userdto.UserDTO{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

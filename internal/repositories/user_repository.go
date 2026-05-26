package repositories

import (
	"context"

	"golang-project-structure/internal/dbconnections"
	userdto "golang-project-structure/internal/dto/user"
	"golang-project-structure/internal/models"

	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(ctx context.Context, u *models.UserModel, opts ...QueryOption) error
	GetByID(ctx context.Context, id uint, opts ...QueryOption) (*models.UserModel, error)
	List(ctx context.Context, params userdto.ListParams, opts ...QueryOption) ([]*models.UserModel, int64, error)
	Update(ctx context.Context, u *models.UserModel, opts ...QueryOption) error
	Delete(ctx context.Context, id uint) error
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(dbConn *dbconnections.AppDB) IUserRepository {
	return &UserRepository{DB: dbConn.DB}
}

func (r *UserRepository) Create(ctx context.Context, u *models.UserModel, opts ...QueryOption) error {
	tx := r.DB.WithContext(ctx)
	for _, opt := range opts {
		tx = opt(tx)
	}
	return tx.Create(u).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id uint, opts ...QueryOption) (*models.UserModel, error) {
	tx := r.DB.WithContext(ctx)
	for _, opt := range opts {
		tx = opt(tx)
	}

	var u models.UserModel
	if err := tx.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) List(ctx context.Context, params userdto.ListParams, opts ...QueryOption) ([]*models.UserModel, int64, error) {
	tx := r.DB.WithContext(ctx).Model(&models.UserModel{})
	for _, opt := range opts {
		tx = opt(tx)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []*models.UserModel
	if err := tx.Offset(params.Offset()).Limit(params.Limit).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepository) Update(ctx context.Context, u *models.UserModel, opts ...QueryOption) error {
	tx := r.DB.WithContext(ctx)
	for _, opt := range opts {
		tx = opt(tx)
	}
	return tx.Save(u).Error
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Delete(&models.UserModel{}, id).Error
}

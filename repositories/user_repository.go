package repositories

import (
	"context"
	"database/sql"

	"github.com/peang/amartha-loan-service/models"
	"github.com/uptrace/bun"
)

type UserRepositoryInterface interface {
	Detail(ctx context.Context, id uint) (loan *models.User, err error)
}

type UserRepositoryFilter struct {
	ID *uint
}

type userRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) UserRepositoryInterface {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Detail(ctx context.Context, id uint) (*models.User, error) {
	var loan models.User
	err := r.db.NewSelect().Model(&loan).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &loan, nil
}

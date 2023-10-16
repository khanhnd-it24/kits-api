package repos

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"kits/api/src/common/fault"
	"kits/api/src/core/domains"
	"kits/api/src/infra/postgres/models"
	"time"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) domains.UserRepo {
	return &userRepo{
		db: db,
	}
}

func (u userRepo) FindByUsername(ctx context.Context, username string) (*domains.User, error) {
	var user models.User
	conds := []clause.Expression{
		clause.Eq{
			Column: "username",
			Value:  username,
		},
		clause.Eq{
			Column: "deleted_at",
			Value:  nil,
		},
	}
	if err := u.db.WithContext(ctx).Clauses(conds...).Take(&user).Error; err != nil {
		return nil, fault.DBWrapf(err, "[UserRepo.FindByUsername] failed to find record")
	}
	return user.ToDomain(), nil
}

func (u userRepo) FindById(ctx context.Context, id string) (*domains.User, error) {
	var user models.User
	conds := []clause.Expression{
		clause.Eq{
			Column: "id",
			Value:  id,
		},
		clause.Eq{
			Column: "deleted_at",
			Value:  nil,
		},
	}
	if err := u.db.WithContext(ctx).Clauses(conds...).Take(&user).Error; err != nil {
		return nil, fault.DBWrapf(err, "[UserRepo.FindById] failed to find record")
	}
	return user.ToDomain(), nil
}

func (u userRepo) Save(ctx context.Context, user *domains.UserCreate) (*domains.User, error) {
	userModel := models.UserCreateFromDomain(user)

	if err := u.db.WithContext(ctx).Save(&userModel).Error; err != nil {
		return nil, fault.DBWrapf(err, "[UserRepo.Save] failed to insert user")
	}
	return userModel.ToDomain(), nil
}

func (u userRepo) FindByIdAndDelete(ctx context.Context, id string) error {
	cond := clause.Eq{
		Column: "id",
		Value:  id,
	}
	if err := u.db.WithContext(ctx).Clauses(cond).Update("deleted_at", time.Now()).Error; err != nil {
		return fault.DBWrapf(err, "[UserRepo.Save] failed to insert user")
	}
	return nil
}

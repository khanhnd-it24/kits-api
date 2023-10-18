package repos

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"kits/api/src/common/fault"
	"kits/api/src/core/domains"
	"kits/api/src/infra/postgres/models"
)

type refreshTokenRepo struct {
	db *gorm.DB
}

func NewRefreshTokenRepo(db *gorm.DB) domains.RefreshTokenRepo {
	return &refreshTokenRepo{
		db: db,
	}
}

func (r refreshTokenRepo) Save(ctx context.Context, token *domains.RefreshToken) (*domains.RefreshToken, error) {
	tokenModel := models.RefreshTokenFromDomain(token)

	if err := r.db.WithContext(ctx).Save(&tokenModel).Error; err != nil {
		return nil, fault.DBWrapf(err, "[RefreshTokenRepo.Save] failed to insert refresh token")
	}
	return tokenModel.ToDomain(), nil
}

func (r refreshTokenRepo) FindByToken(ctx context.Context, token string) (*domains.RefreshToken, error) {
	var tokenModel models.RefreshToken
	cond := clause.Eq{
		Column: "token",
		Value:  token,
	}

	if err := r.db.WithContext(ctx).Clauses(cond).Take(&tokenModel).Error; err != nil {
		return nil, fault.DBWrapf(err, "[UserRepo.FindByToken] failed to find record")
	}

	return tokenModel.ToDomain(), nil
}

func (r refreshTokenRepo) FindByUserIdAndDelete(ctx context.Context, userId int64) error {
	//TODO implement me
	panic("implement me")
}

package domains

import (
	"context"
	"time"
)

type Habit struct {
	Id         int64
	UserId     int64
	Name       string
	Order      int
	Icon       string
	ThemeColor string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

type HabitRepo interface {
	Save(ctx context.Context, habit *Habit) (*Habit, error)
	FindByUserId(ctx context.Context, userId int64) ([]*Habit, error)
}

package domains

import "context"

type DailyHabit struct {
	HabitId int64
	Date    string
}

type DailyHabitRepo interface {
	Save(ctx context.Context, dailyHabit *DailyHabit) (*DailyHabit, error)
	FindByHabitId(ctx context.Context, habitId int64) ([]*DailyHabit, error)
}

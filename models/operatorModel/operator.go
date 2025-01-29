package modeloperator

import (
	"time"
)

type Operator struct {
	ID       int64     `json:"id" gorm:"primarykey;autoIncrement"`
	LoginAt  time.Time `json:"login_at"`
	LogoutAt time.Time `json:"logout_at"`
	Operator int       `json:"operator_id"`
}

package quote

import (
	"fmt"
)

type Quote struct {
	UserID  int    `json:"id"`
	Message string `json:"message"`
}

func (u *Quote) String() string {
	return fmt.Sprintf("%s (%d)", u.Message, u.UserID)
}

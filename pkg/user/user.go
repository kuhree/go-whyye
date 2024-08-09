package user

import (
	"fmt"
)

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

func (u *User) String() string {
	return fmt.Sprintf("%s (%d)", u.Name, u.ID)
}

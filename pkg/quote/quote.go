package quote

import (
	"fmt"
	"math/rand"
)

type Quote struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

func (u *Quote) String() string {
	return fmt.Sprint(u.Message)
}

func GetRandom(quotes []Quote) Quote {
	var quote Quote
	randIdx := rand.Intn(len(quotes))
	quote = quotes[randIdx]

	return quote
}

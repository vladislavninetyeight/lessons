package User

import "time"

type Users []User

type User struct {
	ID        string
	CreatedAt time.Time
}

func (u Users) Len() int {
	return len(u)
}

func (u Users) Less(i, j int) bool {
	return u[j].CreatedAt.Unix() < u[i].CreatedAt.Unix()
}

func (u Users) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

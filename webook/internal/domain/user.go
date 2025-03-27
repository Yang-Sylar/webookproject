package domain

import "time"

// User 领域对象，是 DDD 中的 entity，也有人叫 BO (business object)
type User struct {
	Id       int64
	Email    string
	Password string
	Ctime    time.Time

	Nickname string
	Birthday time.Time
	AboutMe  string
}

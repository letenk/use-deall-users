package domain

import "time"

type User struct {
	ID        string    `bson:"_id"`
	Fullname  string    `bson:"fullname"`
	Username  string    `bson:"username"`
	Password  string    `bson:"password"`
	Role      string    `bson:"role"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

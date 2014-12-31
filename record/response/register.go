package response

import "time"

type Register struct {
	Login		string
	Token		string
	Expires		time.Time
}

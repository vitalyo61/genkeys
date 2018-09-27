package model

const (
	CodeNo   = 0
	CodeYes  = 1
	CodeStop = 2
)

type Code struct {
	Number string `bson:"_id"`
	Status int    `bson:"status"`
}

func (c *Code) GetStatus() string {
	switch c.Status {
	case CodeNo:
		return "не выдан"
	case CodeYes:
		return "выдан"
	case CodeStop:
		return "погашен"
	}
	return "неопределен"
}

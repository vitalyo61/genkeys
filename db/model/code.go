package model

type Code struct {
	Number string `bson:"_id"`
	Status int    `bson:"status"`
}

func (c *Code) GetStatus() string {
	switch c.Status {
	case 0:
		return "не выдан"
	case 1:
		return "выдан"
	case 2:
		return "погашен"
	}
	return "неопределен"
}

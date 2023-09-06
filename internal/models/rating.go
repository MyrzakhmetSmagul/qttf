package models

type Rating struct {
	Id         int
	Position   int
	Player     Player
	Rating     int
	UpdateTime string
}

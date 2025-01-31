package model

type Order struct {
	Type CoffeeType `json:"type"`
	Size CoffeeSize `json:"size"`
}

type CoffeeType string

const (
	Cappuccino CoffeeType = "Cappuccino"
	Latte      CoffeeType = "Latte"
)

type CoffeeSize int

const (
	SmallSize  CoffeeSize = 200
	MediumSize CoffeeSize = 300
	LargeSize  CoffeeSize = 400
)

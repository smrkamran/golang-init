package main

import (
	"fmt"
	"math/rand"
)

type Player interface {
	KickBall()
}

type CR7 struct {
	stamina int
	power   int
	SUI     int
}

type FootballPlayer struct {
	stamina int
	power   int
}

func (c CR7) KickBall() {
	shot := c.stamina + c.power*c.SUI
	fmt.Printf("CR7 The player kicks the ball with a power of %d\n", shot)
}

func (f FootballPlayer) KickBall() {
	shot := f.stamina + f.power
	fmt.Printf("The player kicks the ball with a power of %d\n", shot)
}

func main() {
	team := make([]Player, 11)
	for i := 0; i < len(team)-1; i++ {
		team[i] = FootballPlayer{
			stamina: rand.Intn(10),
			power:   rand.Intn(10),
		}
	}

	team[len(team)-1] = CR7{stamina: 10, power: 10, SUI: 2}

	for _, player := range team {
		player.KickBall()
	}
}

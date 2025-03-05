package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type Player struct {
	HP int
}

func NewPlayer(hp int) actor.Producer {
	return func() actor.Receiver {
		return &Player{HP: hp}
	}
}

type TakeDamage struct {
	amount int
}

func (p *Player) Receive(c *actor.Context) {

	switch msg := c.Message().(type) {
	case actor.Started:
		fmt.Println("Player started")
	case actor.Stopped:
		fmt.Println("Player stopped")
	case TakeDamage:
		fmt.Println("Player is taking damage: ", msg.amount)
	}

}

func main() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatal(err)
	}

	pid := e.Spawn(NewPlayer(100), "player", actor.WithID("myuserid69"))

	msg := TakeDamage{amount: 999}

	for i := 0; i < 100; i++ {
		e.Send(pid, msg)
	}

	time.Sleep(time.Second * 2)

	fmt.Println("Process id", pid)
}

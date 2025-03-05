package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type Inventory struct {
	Bottles int
}

type Player struct {
	HP           int
	inventoryPID *actor.PID
}

type DrinkBottle struct {
	amount int
}

func NewPlayer(hp int) actor.Producer {
	return func() actor.Receiver {
		return &Player{HP: hp}
	}
}

func NewInventory(bottles int) actor.Producer {
	return func() actor.Receiver {
		return &Inventory{Bottles: bottles}
	}
}

func (p *Inventory) Receive(c *actor.Context) {

	switch msg := c.Message().(type) {
	case actor.Started:
		_ = msg
		fmt.Println("Inventory started")
		c.Engine().Subscribe(c.PID())
	case actor.Stopped:
		fmt.Println("Inventory stopped")
	case DrinkBottle:
		fmt.Println("Inventory received drink bottle message")
	case MyEvent:
		fmt.Println("Inventory received message: ", msg.foo)
	}

}

func (p *Player) Receive(c *actor.Context) {

	switch msg := c.Message().(type) {
	case actor.Started:
		fmt.Println("Player started")
		p.inventoryPID = c.SpawnChild(NewInventory(1), "inventory")
		c.Engine().Subscribe(c.PID())
	case actor.Stopped:
		fmt.Println("Player stopped")
	case DrinkBottle:
		c.Forward(p.inventoryPID)
	case MyEvent:
		fmt.Println("Player received message: ", msg.foo)
	}

}

type MyEvent struct {
	foo string
}

func main() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatal(err)
	}

	pid := e.Spawn(NewPlayer(100), "player", actor.WithID("myuserid69"))

	e.Send(pid, DrinkBottle{amount: 55})

	time.Sleep(time.Second * 2)

	e.BroadcastEvent(MyEvent{foo: "This is the mssage"})
	time.Sleep(time.Second * 2)
}

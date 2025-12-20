package main

import (
	"fmt"

	evt "github.com/jb-oliveira/fullcycle-secret/pkg/events"
)

func main() {
	newVar := evt.NewEventDispatcher()
	fmt.Println(newVar)

}

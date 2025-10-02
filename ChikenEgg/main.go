package main

import (
	"fmt"
	"time"
)

var chiken = object{"Курица"}
var egg = object{"Яйцо"}
var Done = make(chan string, 1)

func main() {

	go chiken.run()
	go egg.run()

	<-Done
	fmt.Printf("Спор решен! %s появилось первым", <-Done)
}

type object struct {
	name string
}

func (o object) run() {

	for range 10 {
		time.Sleep(time.Millisecond)
		fmt.Printf("%s\n", o.name)
	}
	Done <- fmt.Sprintf(o.name)
}

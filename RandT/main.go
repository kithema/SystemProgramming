package main

import (
	"fmt"
	"time"
)
//в Go отсутвуют приоритеты потоков
var turtle = animal{"Черепаха", 40} 
var rabbit = animal{"Заяц", 40}
var chWin = make(chan string)

func main() {

	go turtle.run()
	go rabbit.run()
	
	fmt.Println(<-chWin)
	fmt.Println("Конец гонки!")
}

type animal struct {
	name     string
	speed int
}

func (a animal) run() {
	
	for i := range 51 {
		time.Sleep(time.Duration(a.speed) * time.Millisecond)
		fmt.Printf("%s пробежал %d метров\n", a.name, i)
	}
	chWin <- fmt.Sprintf("%s прибежал первым", a.name)
}


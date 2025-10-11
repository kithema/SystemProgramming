package main

import (
	"fmt"
	"time"
)
//Создание Курицы и Яйца
var chiken = object{"Курица"}
var egg = object{"Яйцо"}
//Создание буферизованого канала
var Done = make(chan string, 1)

func main() {

	go chiken.run()
	go egg.run()
	//Пропуск первого завершившего свою работу
	<-Done
	//Вывод второго завершившего свою работу
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
	//Отправка имени в канал
	Done <- fmt.Sprintf(o.name)
}


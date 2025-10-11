package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	//Создание переменной стркутурые WaitGroup
    var wg sync.WaitGroup
    //Цикл с зупском 10 горутин
    for i := 0; i < 10; i++ {
		//Добавление горутины в ожидание
        wg.Add(1)
        go func(id int) {
			seconds := rand.Intn(30)
            fmt.Printf("Горутина %d запущена\n", id)
            time.Sleep(time.Second * time.Duration(seconds))//очень сильно работаем
            fmt.Printf("Горутина %d завершена, с множителем %d\n", id, seconds)
			//Сигнал завершения 
            wg.Done()
        }(i)
    }
    //Ожидаем завершение всех горутин
    wg.Wait()
    fmt.Println("Все горутины завершены")
}


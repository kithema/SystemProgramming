package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
    var wg sync.WaitGroup
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
			seconds := rand.Intn(30)
            fmt.Printf("Горутина %d запущена\n", id)
            time.Sleep(time.Second * time.Duration(seconds))//очень сильно работаем
            fmt.Printf("Горутина %d завершена, с множителем %d\n", id, seconds)
            wg.Done()
        }(i)
    }
    
    wg.Wait()
    fmt.Println("Все горутины завершены")
}

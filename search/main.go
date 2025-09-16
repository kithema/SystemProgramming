package main

import (
	"fmt"
	"math/rand"
)
func main(){
	var array [1000]int
	var (
		max = -1
		max2 = -1
		max7 = -1
		max14 = -1
	)
	for i := 0; i < len(array); i++ {
		array[i] = rand.Intn(10001)
	}
	for i := 0; i < len(array); i++{
		if max < array[i]{
			max = array[i]
		}
		if array[i] %2 == 0 && array[i] %14 != 0{
			if max2 < array[i]{
				max2 = array[i]
			}
		}
		if array[i] %7 == 0 && array[i] %14 != 0{
			if max7 < array[i]{
				max7 = array[i]
			}
		}
		if array[i] %14 == 0{
			if max14 < array[i]{
				max14 = array[i]
			}
		}
	}
	if (max * max14) > (max2 * max7){
		fmt.Printf("Произведение числел %d и %d = %d",max , max14, max * max14)
	}else{
		fmt.Printf("Произведение числел %d и %d = %d",max , max14, max2 * max7)
	}
	
}








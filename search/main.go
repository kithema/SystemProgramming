package main

import (
	"fmt"
	"math/rand"
)
func main(){
	//Массив на 1000 элементов
	var array [1000]int
	//присваиваем каждому элементу -1 в случае если его не окажется в массиве
	var (
		max = -1 //максимальнное
		max2 = -1 //максимальное кратное 2, но не 14
		max7 = -1 //максимальное кратное 7, но не 14
		max14 = -1 //максимальное кратное 14
	)
	//заполнение массива случайными числами
	for i := 0; i < len(array); i++ {
		array[i] = rand.Intn(10001)
	}
	
	for i := 0; i < len(array); i++{
		//поиск максимальнного
		if max < array[i]{
			max = array[i]
		}
		//поиск max2
		if array[i] %2 == 0 && array[i] %14 != 0{
			if max2 < array[i]{
				max2 = array[i]
			}
		}
		//поиск max7
		if array[i] %7 == 0 && array[i] %14 != 0{
			if max7 < array[i]{
				max7 = array[i]
			}
		}
		//поиск max14
		if array[i] %14 == 0{
			if max14 < array[i]{
				max14 = array[i]
			}
		}
	}
	//Находим большее из двух произведений
	if (max * max14) > (max2 * max7){
		fmt.Printf("Произведение числел %d и %d = %d",max , max14, max * max14)
	}else{
		fmt.Printf("Произведение числел %d и %d = %d",max , max14, max2 * max7)
	}
	
}









package main

import (
	"fmt"
	"math"
)

func main(){
	firstWay(4,4)
	secondWay(2,16)
	thirdWay(9, 23)
	fourthWay(9, 52)
	fifthWay(4,4)
	fmt.Println(sixthWay(6,10))
}

func firstWay(a, b int){//сумма
	result := 0
	for i := 0; i  <b; i++{
		result += a
	}
	fmt.Println(result)
}
func secondWay(a, b float32){//деление
	fmt.Println( a/(1/b))
}
func thirdWay(a, b float64){//логарифмы
	fmt.Println(int(math.Pow(2, math.Log2(a)+ math.Log2(b))))
}
func fourthWay(a, b float64){//при помощи тангенса
	fmt.Println(math.Round(1 - (a+b)/(math.Tan((math.Atan(a)+ math.Atan(b))))))
}

func fifthWay(a, b int){
    if a == 0 || b == 0 {
        fmt.Println("0")
    }
    result := 0
    for b > 0 {
        // Если текущий бит b установлен, добавляем a с соответствующим сдвигом
        if b&1 == 1 {
            result += a
        }
        a <<= 1// Сдвигаем a влево (умножаем на 2)
        b >>= 1// Сдвигаем b вправо (делим на 2)
    }
	fmt.Println(result)
}

func sixthWay(a,b int) int{
	if (a == 0 || b == 0){
		return 0;
	} 
    if (b == 1){
		return a;
	} 
    if (b == -1){
		 return -a;
	}

	
	halfB := b / 2
    halfProduct := sixthWay(a, halfB)
    remainder := b % 2
    if remainder != 0 {
        if remainder == 1 || remainder == -1 {
            return halfProduct + halfProduct + a
        }
    }
    return halfProduct + halfProduct
}
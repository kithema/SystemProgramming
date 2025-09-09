package main

import (
	"math"
    _ "fmt"
)
func main(){
    
}
func firstWay(a, b int) (int, error) {
    if a == 0 || b == 0 {
        return 0, nil
    }
    result := 0
    sign := 1
    absB := b
    absA := a
    if b < 0 {
        absB = -b
    }
    if a < 0 {
        absA = -a
    }
    for i := 0; i < absB; i++ {
        result += absA
    }
    if ((a < 0 && b > 0)||(b < 0 && a > 0  )){
        sign = -1
    }
    return sign * result, nil
}

func secondWay(a, b float64) (float64, error) {
    if b == 0 {
        return 0, nil
    }
    return a / (1 / b), nil
}
func thirdWay(a, b float64)(float64, error){

    result := math.Round(math.Pow(2, math.Log2(math.Abs(a))+ math.Log2(math.Abs(b))))
    var sign float64 = 1
    if ((a < 0 && b > 0)||(b < 0 && a > 0  )){
        sign = -1
    }
	return sign * result, nil
}
func fourthWay(a, b float64)(float64, error){
    result := (math.Round(1 - (a+b)/(math.Tan((math.Atan(a)+ math.Atan(b))))))
	return result, nil
}

func fifthWay(a, b int)(int, error){
    if a == 0 || b == 0 {
        return 0, nil
    }
    result := 0
    sign := 1
    if ((a < 0 && b > 0)||(b < 0 && a > 0  )){
        sign = -1
    }
    absB := b
    absA := a
    if b < 0 {
        absB = -b
    }
    if a < 0 {
        absA = -a
    }
    for absB > 0 {

        if absB&1 == 1 {
            result += absA
        }
        absA <<= 1
        absB >>= 1
    }
	return sign * result, nil
}

func sixthWay(a,b int) (int, error){
	if (a == 0 || b == 0){
		return 0, nil;
	} 
    if (b == 1){
		return a, nil;
	} 
    if (b == -1){
		 return -a, nil;
	}

	
	halfB := b / 2
    halfProduct, _ := sixthWay(a, halfB)
    remainder := b % 2
    if remainder != 0 {
        if remainder == 1 || remainder == -1 {
            return halfProduct + halfProduct + a, nil
        }
    }
    return halfProduct + halfProduct, nil
}
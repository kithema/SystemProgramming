package main

import "fmt"

func main() {
	var money, price, wrap int

	fmt.Print("Введите количество денег: ")
	fmt.Scan(&money)
	fmt.Print("Введите цену за 1 шоколадку: ")
	fmt.Scan(&price)
	fmt.Print("Введите количество оберток для бесплатной шоколадки: ")
	fmt.Scan(&wrap)

	if price <= 0 || wrap <= 0 {
		fmt.Println("Ошибка: цена и количество оберток должны быть больше 0")
		return
	}

	total := recursiveChocolate(money, price, wrap, 0)
	fmt.Printf("Всего можно получить шоколадок: %d\n", total)
}

func recursiveChocolate(money, price, wrap, wrappers int) int {

	if money < price && wrappers < wrap {
		return 0
	}

	var chocolates int

	if money >= price {
		chocolates = money / price
		money = money % price
		wrappers += chocolates
		fmt.Printf("Куплено %d шоколадок. Оберток: %d\n", chocolates, wrappers)
	}

	if wrappers >= wrap {
		newChocolates := wrappers / wrap
		remainingWrappers := wrappers % wrap
		fmt.Printf("Обменяли обертки на %d шоколадок\n", newChocolates)

		return chocolates + newChocolates + recursiveChocolate(money, price, wrap, remainingWrappers+newChocolates)
	}

	return chocolates
}

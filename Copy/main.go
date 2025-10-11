package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)



func copyFile(src, forCopy string) {
	//Файл для копирования
	file, err := os.Create(src)
	if err != nil {
		fmt.Println("Ошибка при создании: ", err)
	}
	//Запись 1000 строк
	data := "1000 строк тестирования\n"
	for i := 0; i < 1000; i++ {
		file.WriteString(data)
	}
	file.Close()
	//Открытие исходного файла для перемещение указателя
	file, err = os.Open(src)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	//Создания файла куда будем копировать
	copyfile, err := os.Create(forCopy)
	if err != nil {
		fmt.Println("Ошибка при создании: ", err)
	}
	defer copyfile.Close()
	//Копирование
	io.Copy(copyfile, file)
}
func main() {

	//последовательное копирование
	start := time.Now() // Фиксируем начальное время
	copyFile("file.txt", "copyfile.txt")
	copyFile("file1.txt", "copyfile1.txt")
	seconds := time.Since(start) // Вычисляем прошедшее время
	fmt.Printf("Последовательное заняло %v\n", seconds)

	//параллельное копирование
	var wg sync.WaitGroup
	start = time.Now()
	wg.Add(2)
	go func() {
		defer wg.Done()
		copyFile("gofile.txt", "gocopyfile.txt")
	}()
	go func() {
		defer wg.Done()
		copyFile("gofile1.txt", "gocopyfile1.txt")
	}()
	wg.Wait()
	fmt.Printf("Параллельное выполнение: %v\n", time.Since(start))
}

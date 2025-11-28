package main

import (
	"bufio"    
	"fmt"     
	"net"      
	"os"       
	"strings"  
)

// Глобальные константы конфигурации сервера
const (
	IP_ADDR = "127.0.0.1" // Локальный адрес для тестирования
	                       // Измените на "10.254.10.145" для работы в сети
	PORT    = "8189"      // Порт сервера чата
)

func main() {
	// Получение никнейма пользователя
	var nickname string
	
	fmt.Print("Введите никнейм (по умолчанию 'Гость'): ")
	fmt.Scanln(&nickname)
	
	// Установка никнейма по умолчанию при пустом вводе
	if nickname == "" {
		nickname = "Гость"
	}

	// Установка TCP-соединения с сервером
	fmt.Printf("Подключение к %s:%s...\n", IP_ADDR, PORT)
	
	conn, err := net.Dial("tcp", IP_ADDR+":"+PORT)
	if err != nil {
		// Обработка ошибки подключения
		fmt.Printf("Ошибка подключения: %v\n", err)
		os.Exit(1)
	}
	
	// Гарантированное закрытие соединения при выходе из функции
	defer conn.Close()
	
	fmt.Println("Подключение установлено!")
	fmt.Println("Введите сообщения (Enter для отправки, 'exit' для выхода)")
	fmt.Println("────────────────────────────────────")

	// Запуск горутины для асинхронного чтения сообщений от сервера
	go func() {
		// Буферизированный сканер для чтения строк из TCP-соединения
		scanner := bufio.NewScanner(conn)
		
		// Бесконечный цикл чтения входящих сообщений
		for scanner.Scan() {
			// Вывод полученного сообщения в консоль
			fmt.Println(scanner.Text())
		}
		
		// Проверка ошибок чтения после завершения сканирования
		if err := scanner.Err(); err != nil {
			fmt.Printf("\nОшибка чтения: %v\n", err)
		}
	}()

	// Основной цикл чтения и отправки сообщений пользователя
	reader := bufio.NewReader(os.Stdin)
	for {
		// Отображение приглашения для ввода с цветным никнеймом
		fmt.Print("\033[32m" + nickname + ": \033[0m")
		
		// Чтение строки ввода пользователя
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Ошибка ввода: %v\n", err)
			break
		}
		
		// Удаление пробелов в начале и конце сообщения
		msg = strings.TrimSpace(msg)
		
		// Проверка команд выхода
		if msg == "" || msg == "exit" || msg == "quit" {
			fmt.Println("Выход из чата...")
			return
		}
		
		// Формирование полного сообщения с никнеймом
		fullMsg := fmt.Sprintf("%s: %s", nickname, msg)
		
		// Отправка сообщения на сервер
		fmt.Fprintln(conn, fullMsg)
	}
}
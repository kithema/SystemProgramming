package main

import (
	"bufio"    
	"fmt"     
	"net"      
	"strings"  
	"sync"     
)

// Структура сервера чата
type Server struct {
	connections map[net.Conn]*Connection // Карта активных подключений
	mutex       sync.RWMutex             // Мьютекс для безопасного доступа к connections
}

// Структура одного клиентского подключения
type Connection struct {
	conn net.Conn // TCP-соединение с клиентом
	name string   // Имя клиента (по умолчанию "Гость")
}

func main() {
	// Инициализация сервера с пустой картой подключений
	server := &Server{
		connections: make(map[net.Conn]*Connection),
	}

	fmt.Println("Server running on :8189...")
	fmt.Println("Ожидание подключений...")

	// Создание TCP-слушателя на порту 8189
	ln, err := net.Listen("tcp", ":8189")
	if err != nil {
		panic(err) // Критическая ошибка - завершение программы
	}
	defer ln.Close() // Гарантированное закрытие слушателя

	// Основной цикл принятия подключений
	for {
		conn, err := ln.Accept() // Ожидание нового клиента
		if err != nil {
			fmt.Printf("Accept error: %v\n", err)
			continue // Пропускаем ошибку и ждем следующее подключение
		}

		// Создание новой структуры подключения с именем по умолчанию
		connection := &Connection{conn: conn, name: "Гость"}
		
		// Добавление клиента в список активных подключений
		server.addConnection(connection)

		// Запуск горутины для обработки сообщений от этого клиента
		go server.handleConnection(connection)
	}
}

// Добавление нового подключения в список активных клиентов
func (s *Server) addConnection(conn *Connection) {
	s.mutex.Lock()
	s.connections[conn.conn] = conn // Сохранение в карте по указателю на соединение
	s.mutex.Unlock()

	// Уведомление всех клиентов о новом подключении
	msg := fmt.Sprintf("Клиент подключился: %s", conn)
	s.broadcast(msg)
	fmt.Println(msg) 
}

// Удаление отключенного клиента из списка
func (s *Server) removeConnection(conn *Connection) {
	s.mutex.Lock()
	delete(s.connections, conn.conn) // Удаление из карты
	s.mutex.Unlock()

	// Уведомление всех клиентов об отключении
	msg := fmt.Sprintf("Клиент отключился: %s", conn)
	s.broadcast(msg)
	fmt.Println(msg) // Лог для сервера
}

// Обработка сообщений от одного клиента
func (s *Server) handleConnection(conn *Connection) {
	// Гарантированное удаление клиента и закрытие соединения при выходе
	defer s.removeConnection(conn)
	defer conn.conn.Close()

	// Уведомление клиента об успешном подключении
	fmt.Fprintln(conn.conn, "Подключение установлено...")

	// Буферизированный сканер для чтения сообщений построчно
	scanner := bufio.NewScanner(conn.conn)
	
	// Основной цикл обработки сообщений от клиента
	for scanner.Scan() {
		msg := strings.TrimSpace(scanner.Text()) // Удаление пробелов
		
		// Пропуск пустых сообщений
		if msg == "" {
			continue
		}
		
		// Формирование полного сообщения с именем клиента
		fullMsg := fmt.Sprintf("%s: %s", conn.name, msg)
		
		// Рассылка сообщения всем подключенным клиентам
		s.broadcast(fullMsg)
	}

	// Обработка ошибок чтения после завершения цикла
	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка чтения от %s: %v\n", conn, err)
	}
}

// Рассылка сообщения всем активным клиентам
func (s *Server) broadcast(message string) {
	s.mutex.RLock()   // Блокировка только для чтения
	defer s.mutex.RUnlock()

	// Перебор всех активных подключений
	for _, conn := range s.connections {
		fmt.Fprintln(conn.conn, message) // Отправка сообщения клиенту
	}
}

// Реализация метода String() для красивого вывода Connection
func (c Connection) String() string {
	// Возвращает IP:порт клиента в формате "192.168.1.100:54321"
	return fmt.Sprintf("%s:%d",
		c.conn.RemoteAddr().(*net.TCPAddr).IP,
		c.conn.RemoteAddr().(*net.TCPAddr).Port)
}
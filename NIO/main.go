package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	// Определение имен файлов для работы
	sourceFile := "source.txt"       // Исходный файл с данными
	sequentialFile := "sequential.txt" // Файл для последовательного копирования
	concurrentFile := "concurrent.txt" // Файл для параллельного копирования

	// Создаем исходный файл со 100 строками тестовых данных
	if err := createSourceFile(sourceFile); err != nil {
		log.Fatal("Ошибка при создании исходного файла: ", err)
	}

	// Последовательное копирование с замером времени
	startSequential := time.Now() // Засекаем время начала
	if err := sequentialCopy(sourceFile, sequentialFile); err != nil {
		log.Fatal("Ошибка при последовательном копировании: ", err)
	}
	sequentialTime := time.Since(startSequential) // Вычисляем затраченное время
	fmt.Printf("Время последовательного копирования: %v\n", sequentialTime)

	// Параллельное копирование с замером времени
	startConcurrent := time.Now() // Засекаем время начала
	if err := concurrentCopy(sourceFile, concurrentFile); err != nil {
		log.Fatal("Ошибка при параллельном копировании: ", err)
	}
	concurrentTime := time.Since(startConcurrent) // Вычисляем затраченное время
	fmt.Printf("Время параллельного копирования: %v\n", concurrentTime)
}

// createSourceFile создает исходный файл с тестовыми данными
func createSourceFile(filename string) error {
	// Создаем новый файл
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close() // Гарантируем закрытие файла при выходе из функции

	// Создаем буферизованного писателя для эффективной записи
	writer := bufio.NewWriter(file)
	defer writer.Flush() // Гарантируем сброс буфера при выходе

	// Записываем 100 тестовых строк в файл
	for i := 1; i <= 100; i++ {
		line := fmt.Sprintf("Строка %d: это тестовые данные для проверки копирования\n", i)
		if _, err := writer.WriteString(line); err != nil {
			return err
		}
	}
	return nil
}

// sequentialCopy выполняет последовательное копирование файла построчно
func sequentialCopy(source, dest string) error {
	// Открываем исходный файл для чтения
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close() // Гарантируем закрытие

	// Создаем файл назначения
	dst, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dst.Close() // Гарантируем закрытие

	// Создаем сканер для чтения и писатель для записи
	scanner := bufio.NewScanner(src)
	writer := bufio.NewWriter(dst)
	defer writer.Flush() // Гарантируем сброс буфера

	// Последовательно читаем и записываем каждую строку
	for scanner.Scan() {
		line := scanner.Text() + "\n"
		if _, err := writer.WriteString(line); err != nil {
			return err
		}
	}

	// Возвращаем ошибку сканирования, если она возникла
	return scanner.Err()
}

// concurrentCopy выполняет параллельное копирование с использованием горутин
func concurrentCopy(source, dest string) error {
	// Открываем исходный файл и создаем файл назначения
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Создаем канал для передачи строк между горутинами (буфер 10 строк)
	lines := make(chan string, 10)
	var wg sync.WaitGroup    // Для ожидания завершения горутин
	var writeErr error       // Для захвата ошибки записи
	var writeErrMux sync.Mutex // Для безопасного доступа к writeErr

	// Горутина для чтения из файла
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(src)
		
		// Читаем файл построчно и отправляем строки в канал
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines) // Закрываем канал после завершения чтения
		
		// Проверяем ошибки сканирования
		if err := scanner.Err(); err != nil {
			log.Printf("Ошибка чтения: %v", err)
		}
	}()

	// Горутина для записи в файл
	wg.Add(1)
	go func() {
		defer wg.Done()
		writer := bufio.NewWriter(dst)
		defer writer.Flush()
		
		// Читаем строки из канала и записываем в файл
		for line := range lines {
			if _, err := writer.WriteString(line + "\n"); err != nil {
				// Захватываем ошибку записи (только первую)
				writeErrMux.Lock()
				if writeErr == nil {
					writeErr = err
				}
				writeErrMux.Unlock()
				return // Выходим при ошибке
			}
		}
	}()

	// Ожидаем завершения обеих горутин
	wg.Wait()
	return writeErr
}
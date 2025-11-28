package parser

import (
    "bufio"
    "bytes"
    "os"
    "strings"
    "fmt"
    "media-downloader/config"
    "media-downloader/models"
)

// ExtractMediaUrls - главный парсер файла inFile.txt
// Извлекает URL аудиофайлов (.mp3) и обложек (.jpg/.jpeg/.png/.gif)
// Возвращает структуру с разделенными списками URL
func ExtractMediaUrls() (*models.MediaUrls, error) {
    var audioLinks []string  // Список URL аудиофайлов
    var coverLinks []string  // Список URL обложек

    // Открываем входной файл
    file, err := os.Open(config.InputFilePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    fmt.Println("🔍 Чтение inFile.txt (Windows совместимый режим)...")
    
    // Настраиваем сканер для корректной работы 
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines) 
    
    lineNum := 0
    for scanner.Scan() {
        lineNum++
        
        // Убираем завершающие символы
        rawLine := scanner.Bytes()
        line := strings.TrimSpace(string(bytes.TrimRight(rawLine, "\r\n \t")))
        
        // Пропускаем пустые строки
        if line == "" {
            continue
        }
        
        //Исправляем некорректные URL для Windows
        line = fixUrlForWindows(line)
        
        fmt.Printf("📄 Строка %d: [%s]\n", lineNum, line)
        
        //Классифицируем URL
        if strings.Contains(strings.ToLower(line), ".mp3") {
            audioLinks = append(audioLinks, line)
            fmt.Printf("✅ АУДИО #%d: %s\n", len(audioLinks), line)
        } else if isImageUrl(line) {
            coverLinks = append(coverLinks, line)
            fmt.Printf("✅ ОБЛОЖКА #%d: %s\n", len(coverLinks), line)
        }
    }

    // Выводим итоговую статистику
    fmt.Printf("\n🎵 Найдено аудио: %d\n", len(audioLinks))
    fmt.Printf("🖼️ Найдено обложек: %d\n", len(coverLinks))
    
    return &models.MediaUrls{
        AudioUrls: audioLinks,
        CoverUrls: coverLinks,
    }, scanner.Err()
}

// fixUrlForWindows - исправляет типичные ошибки в URL
// Поддерживает исправления:
// 1. https без :// → https://
// 2. Домены без протокола → https://domain
func fixUrlForWindows(urlStr string) string {
    lowerUrl := strings.ToLower(urlStr)

    //Исправление: https без ://
    if strings.HasPrefix(lowerUrl, "https") && 
       !strings.Contains(lowerUrl, "://") {
        
        slashIndex := strings.Index(lowerUrl, "/")
        if slashIndex > 6 {
            return "https://" + urlStr[5:] 
        }
        return "https://" + urlStr
    }

    // Исправление: домены без протокола
    if !strings.HasPrefix(lowerUrl, "http://") && 
       !strings.HasPrefix(lowerUrl, "https://") {
        if strings.Contains(lowerUrl, ".") {
            return "https://" + urlStr
        }
    }
    
    return urlStr
}

// isImageUrl - определяет, является ли URL ссылкой на изображение
// Поддерживаемые форматы: .jpg, .jpeg, .png, .gif
func isImageUrl(url string) bool {
    lower := strings.ToLower(url)
    return strings.Contains(lower, ".jpg") || 
           strings.Contains(lower, ".jpeg") || 
           strings.Contains(lower, ".png") || 
           strings.Contains(lower, ".gif")
}
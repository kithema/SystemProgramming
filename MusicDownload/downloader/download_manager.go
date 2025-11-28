package downloader

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "sync"
    "time"

    "media-downloader/config"
    "media-downloader/parser"
    "media-downloader/utils"
)

// DownloadManager - основной менеджер загрузки файлов
type DownloadManager struct {
    // audioCoverMapping - сопоставление ID аудио и пути к обложке
    audioCoverMapping map[int]string
    // mu - мьютекс для безопасного доступа к mapping
    mu                sync.RWMutex
}

// DownloadResult - результат загрузки
type DownloadResult struct {
    AudioFilesCount   int               // Количество загруженных аудиофайлов
    CoverFilesCount   int               // Количество загруженных обложек
    AudioCoverMapping map[int]string    // Сопоставление ID → путь к обложке
}

// NewDownloadManager - создает новый менеджер загрузки
func NewDownloadManager() *DownloadManager {
    return &DownloadManager{
        audioCoverMapping: make(map[int]string),
    }
}

// ProcessMediaFiles - главный метод загрузки всех медиафайлов
// 1. Парсит inFile.txt
// 2. Запускает параллельную загрузку с семафором
// 3. Создает сопоставление аудио↔обложка
// 4. Пропускает уже существующие файлы
func (dm *DownloadManager) ProcessMediaFiles() DownloadResult {
    fmt.Println("📥 Начинаем загрузку...")
    
    //Парсим входной файл
    urls, err := parser.ExtractMediaUrls()
    if err != nil {
        fmt.Printf("❌ Ошибка парсинга: %v\n", err)
        return DownloadResult{}
    }

    var result DownloadResult
    result.AudioCoverMapping = make(map[int]string)
    
    var wg sync.WaitGroup
    //Семфор для ограничения количества одновременных загрузок
    semaphore := make(chan struct{}, config.MaxDownloadThreads)

    // Запускаем горутину для каждого аудиофайла
    for i, url := range urls.AudioUrls {
        wg.Add(1)
        go func(id int, audioUrl string) {
            defer wg.Done()
            semaphore <- struct{}{} // Захватываем слот семафора
            defer func() { <-semaphore }() // Освобождаем слот
            
            // Загружаем аудиофайл
            audioFilename := fmt.Sprintf("audio_%03d.mp3", id+1)
            audioPath := filepath.Join(config.AudioFolder, audioFilename)
            
            if !utils.FileExists(audioPath) {
                if err := dm.downloadFileWithRetry(audioUrl, audioPath, true); err != nil {
                    fmt.Printf("❌ Аудио %d: %v\n", id+1, err)
                    return
                }
                // Безопасно увеличиваем счетчик
                dm.mu.Lock()
                result.AudioFilesCount++
                dm.mu.Unlock()
            } else {
                fmt.Printf("⏭ Аудио %d уже существует\n", id+1)
            }

            // Создаем начальное сопоставление 
            dm.mu.Lock()
            result.AudioCoverMapping[id+1] = ""
            dm.audioCoverMapping[id+1] = ""
            dm.mu.Unlock()
            
            // Загружаем обложку (если есть соответствующая)
            if i < len(urls.CoverUrls) && urls.CoverUrls[i] != "" {
                coverFilename := fmt.Sprintf("cover_%03d.jpg", id+1)
                coverPath := filepath.Join(config.CoverFolder, coverFilename)
                
                // Обновляем сопоставление
                dm.mu.Lock()
                result.AudioCoverMapping[id+1] = coverPath
                dm.audioCoverMapping[id+1] = coverPath
                dm.mu.Unlock()
                
                if !utils.FileExists(coverPath) {
                    if err := dm.downloadFileWithRetry(urls.CoverUrls[i], coverPath, false); err != nil {
                        fmt.Printf("❌ Обложка %d: %v\n", id+1, err)
                        return
                    }
                    dm.mu.Lock()
                    result.CoverFilesCount++
                    dm.mu.Unlock()
                } else {
                    fmt.Printf("⏭ Обложка %d уже существует\n", id+1)
                }
            }
        }(i, url)
    }
    
    //Ожидаем завершения всех загрузок
    wg.Wait()
    close(semaphore)
    
    fmt.Printf("✅ Загрузка завершена! %d аудио + %d обложек\n", result.AudioFilesCount, result.CoverFilesCount)
    return result
}

// GetAudioCoverMapping - возвращает копию сопоставления аудио↔обложка
func (dm *DownloadManager) GetAudioCoverMapping() map[int]string {
    dm.mu.RLock()
    defer dm.mu.RUnlock()
    // Создаем копию для безопасной передачи
    mapping := make(map[int]string)
    for k, v := range dm.audioCoverMapping {
        mapping[k] = v
    }
    return mapping
}

// downloadFileWithRetry - загружает файл с повторными попытками
// Максимум 3 попытки с увеличивающейся задержкой (2s, 4s, 6s)
func (dm *DownloadManager) downloadFileWithRetry(url, path string, isAudio bool) error {
    // Создаем директорию, если не существует
    if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
        return fmt.Errorf("папка: %v", err)
    }
    
    maxRetries := 3
    for attempt := 1; attempt <= maxRetries; attempt++ {
        if err := dm.downloadFile(url, path, isAudio); err == nil {
            return nil // Успешно!
        }
        
        if attempt < maxRetries {
            waitTime := time.Duration(attempt*2) * time.Second
            fmt.Printf("⏳ Retry %d/%d через %v...\n", attempt, maxRetries, waitTime)
            time.Sleep(waitTime)
        }
    }
    
    return fmt.Errorf("не удалось скачать после %d попыток", maxRetries)
}

// downloadFile - выполняет одну попытку загрузки файла
// Настраивает HTTP-клиент с заголовками браузера
func (dm *DownloadManager) downloadFile(url, path string, isAudio bool) error {
    // HTTP-клиент с таймаутом 60 секунд
    client := &http.Client{
        Timeout: 60 * time.Second,
    }
    
    // Создаем HTTP-запрос
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }

    // Имитируем браузер для обхода ограничений
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
    req.Header.Set("Accept", "*/*")
    req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")
    req.Header.Set("Accept-Encoding", "gzip, deflate, br")
    req.Header.Set("Connection", "keep-alive")
    req.Header.Set("Sec-Fetch-Dest", "audio") 
    
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    // Проверяем HTTP-статус
    if resp.StatusCode != 200 {
        return fmt.Errorf("HTTP %d", resp.StatusCode)
    }

    // Устанавливаем лимит размера файла
    maxSize := int64(100 * 1024 * 1024) // 100MB для аудио
    if !isAudio {
        maxSize = 10 * 1024 * 1024 // 10MB для обложек
    }
    
    // Сохраняем файл с контролем размера
    size, err := dm.saveWithSizeLimit(resp.Body, path, maxSize)
    if err != nil {
        return err
    }
    
    // Выводим информацию о загруженном файле
    if isAudio {
        fmt.Printf("🎵 %s (%s)\n", filepath.Base(path), formatFileSize(size))
    } else {
        fmt.Printf("🖼️ %s (%s)\n", filepath.Base(path), formatFileSize(size))
    }
    
    return nil
}

// saveWithSizeLimit - сохраняет поток в файл с контролем размера
// Автоматически удаляет файл при превышении лимита
func (dm *DownloadManager) saveWithSizeLimit(reader io.Reader, path string, maxSize int64) (int64, error) {
    file, err := os.Create(path)
    if err != nil {
        return 0, err
    }
    defer file.Close()
    
    var written int64
    buf := make([]byte, 64*1024) // 64KB буфер
    
    for {
        n, err := reader.Read(buf)
        if n > 0 {
            written += int64(n)
            // Проверяем превышение лимита
            if written > maxSize {
                os.Remove(path) // Удаляем невалидный файл
                return 0, fmt.Errorf("превышен лимит размера: %d MB", maxSize/(1024*1024))
            }
            if _, err := file.Write(buf[:n]); err != nil {
                os.Remove(path) // Удаляем поврежденный файл
                return 0, err
            }
        }
        if err == io.EOF {
            break // ✅ Завершение потока
        }
        if err != nil {
            os.Remove(path) // Удаляем при ошибке
            return 0, err
        }
    }
    
    return written, nil
}

// formatFileSize - форматирует размер файла
func formatFileSize(bytes int64) string {
    if bytes < 1024 {
        return fmt.Sprintf("%d B", bytes)
    }
    if bytes < 1024*1024 {
        return fmt.Sprintf("%d KB", bytes/1024)
    }
    return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
}
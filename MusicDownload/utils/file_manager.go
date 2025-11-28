package utils

import (
    "os"
    "sync"
    "media-downloader/config"
)

var (
    // downloadCompleted - флаг завершения загрузки (доступен для чтения/записи из разных горутин)
    downloadCompleted bool
    // downloadMu - мьютекс для безопасного доступа к downloadCompleted
    downloadMu        sync.RWMutex
)

// IsDownloadCompleted - безопасно проверяет статус завершения загрузки
// Использует RWMutex для множественного чтения
func IsDownloadCompleted() bool {
    downloadMu.RLock()   // Блокируем только для чтения
    defer downloadMu.RUnlock()
    return downloadCompleted
}

// SetDownloadCompleted - безопасно устанавливает статус завершения загрузки
func SetDownloadCompleted(completed bool) {
    downloadMu.Lock()    // Блокируем для записи
    defer downloadMu.Unlock()
    downloadCompleted = completed
}

// PrepareDirectories - создает необходимые директории для хранения файлов
// Возвращает ошибку, если создание не удалось
func PrepareDirectories() error {
    // Создаем папку для аудиофайлов
    if err := os.MkdirAll(config.AudioFolder, 0755); err != nil {
        return err
    }
    // Создаем папку для обложек
    if err := os.MkdirAll(config.CoverFolder, 0755); err != nil {
        return err
    }
    return nil
}

// FileExists - проверяет существование файла по пути
// Возвращает true, если файл существует и доступен
func FileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}
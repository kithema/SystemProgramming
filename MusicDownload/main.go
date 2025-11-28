package main

import (
    "fmt"
    "sync"
    "time"
    "media-downloader/downloader"
    "media-downloader/player"
    "media-downloader/utils"
)

// 1. Инициализирует необходимые директории
// 2. Запускает загрузку медиафайлов в фоновом режиме
// 3. Запускает интерактивный медиаплеер
// 4. Ожидает завершения загрузки
func main() {
    fmt.Println("🎵 Медиазагрузчик запущен")
    // ✅ Проверяем и создаем необходимые директории
    if err := utils.PrepareDirectories(); err != nil {
        fmt.Printf("❌ Ошибка создания директорий: %v\n", err)
        return
    }

    // 📥 Создаем менеджер загрузки
    dm := downloader.NewDownloadManager()

    // 🟢 Запускаем загрузку медиафайлов в отдельной горутине
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        // Загружаем файлы и выводим итоговую статистику
        result := dm.ProcessMediaFiles()
        displayDownloadSummary(result)
        // Сигнализируем плееру о завершении загрузки
        utils.SetDownloadCompleted(true)
    }()

    // ⏳ Даем загрузке немного времени на старт
    time.Sleep(500 * time.Millisecond)

    // 🎵 Запускаем интерактивный плеер
    coverMap := dm.GetAudioCoverMapping()
    mediaPlayer := player.NewMediaPlayer(coverMap)
    mediaPlayer.StartPlaybackInterface()

    // ⌛ Ожидаем завершения загрузки
    wg.Wait()
    fmt.Println("✅ Работа завершена")
}

// displayDownloadSummary - отображает итоговую статистику загрузки
func displayDownloadSummary(result downloader.DownloadResult) {
    fmt.Println("\n📈 Итоги загрузки:")
    fmt.Printf("🎵 Аудиофайлов: %d\n", result.AudioFilesCount)
    fmt.Printf("🖼️  Обложек: %d\n", result.CoverFilesCount)
    fmt.Printf("🔗 Связей: %d\n", len(result.AudioCoverMapping))
}
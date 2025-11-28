package player

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "sort"
    "strconv"
    "strings"
    "time"
    "media-downloader/config"
    "media-downloader/utils"
)

// MediaPlayer - основная структура медиаплеера
type MediaPlayer struct {
    // coverMap - сопоставление ID трека и пути к обложке
    coverMap map[int]string
    // scanner - сканер для чтения ввода пользователя
    scanner  *bufio.Scanner
    // active - флаг активности плеера
    active   bool
}

// NewMediaPlayer - создает новый экземпляр медиаплеера
func NewMediaPlayer(coverMap map[int]string) *MediaPlayer {
    return &MediaPlayer{
        coverMap: coverMap,
        scanner:  bufio.NewScanner(os.Stdin),
        active:   true,
    }
}

// StartPlaybackInterface - главный цикл медиаплеера
// Работает до тех пор, пока active == true
func (mp *MediaPlayer) StartPlaybackInterface() {
    fmt.Println("\n🎵 Запуск медиаплеера")

    for mp.active {
        // 🔍 Сканируем доступные аудиофайлы
        tracks, err := mp.scanAudioFiles()
        if err != nil {
            fmt.Printf("❌ Ошибка сканирования: %v\n", err)
            continue
        }

        // Если файлов нет, проверяем статус загрузки
        if len(tracks) == 0 {
            if !utils.IsDownloadCompleted() {
                fmt.Println("\n⏳ Аудиофайлы еще загружаются...")
                time.Sleep(2 * time.Second)
                continue
            }
            mp.handleNoFiles() // Обрабатываем отсутствие файлов
            continue
        }

        // Показываем плейлист и обрабатываем выбор пользователя
        mp.displayPlaylist(tracks)
        if !mp.processUserSelection(tracks) {
            break
        }
    }
}

// AudioTrack - структура для представления аудиотрека
type AudioTrack struct {
    File  string // Полный путь к файлу
    Id    int    // Номер трека (из имени файла)
    Size  int64  // Размер файла в байтах
}

// scanAudioFiles - сканирует папку audio/ и возвращает список валидных MP3 файлов
// Фильтрует файлы по:
// - Расширению .mp3
// - Размеру > 1KB (исключаем пустые/поврежденные файлы)
// - Корректному номеру в имени (audio_001.mp3)
func (mp *MediaPlayer) scanAudioFiles() ([]AudioTrack, error) {
    files, err := os.ReadDir(config.AudioFolder)
    if err != nil {
        return nil, err
    }

    var validTracks []AudioTrack
    for _, file := range files {
        // Проверяем, что это не директория и файл имеет расширение .mp3
        if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".mp3") {
            fullPath := filepath.Join(config.AudioFolder, file.Name())
            info, err := file.Info()
            if err == nil && info.Size() > 1024 { // Минимум 1KB
                id := extractAudioNumber(file.Name())
                if id > 0 { // Корректный номер трека
                    validTracks = append(validTracks, AudioTrack{
                        File: fullPath,
                        Id:   id,
                        Size: info.Size(),
                    })
                }
            }
        }
    }

    // 🔢 Сортируем треки по ID (порядковому номеру)
    sort.Slice(validTracks, func(i, j int) bool {
        return validTracks[i].Id < validTracks[j].Id
    })

    return validTracks, nil
}

// displayPlaylist - отображает список доступных треков
func (mp *MediaPlayer) displayPlaylist(tracks []AudioTrack) {
    fmt.Printf("\n📋 Плейлист (%d треков):\n", len(tracks))
    for i, track := range tracks {
        // Проверяем наличие обложки для трека
        coverStatus := mp.checkCoverStatus(track.Id)
        fileSize := formatFileSize(track.Size)
        fmt.Printf("%d. %s %s [%s]\n", i+1, filepath.Base(track.File), coverStatus, fileSize)
    }
    fmt.Println("r - Обновить список")
    fmt.Println("d - Проверить загрузку")
    fmt.Println("0 - Выйти")
}

// processUserSelection - обрабатывает выбор пользователя
// Возвращает false для выхода из плеера
func (mp *MediaPlayer) processUserSelection(tracks []AudioTrack) bool {
    fmt.Print("Выберите трек: ")
    mp.scanner.Scan()
    choice := strings.TrimSpace(mp.scanner.Text())

    switch strings.ToLower(choice) {
    case "0":
        mp.active = false
        return false
    case "r": // Обновить список
        return true
    case "d": // Показать статус загрузки
        mp.checkDownloadStatus()
        return true
    default:
        // Пытаемся распарсить номер трека
        if trackIndex, err := strconv.Atoi(choice); err == nil && trackIndex > 0 && trackIndex <= len(tracks) {
            mp.playSelectedTrack(tracks[trackIndex-1])
        } else {
            fmt.Println("❌ Неверный номер трека!")
        }
    }
    return true
}

// playSelectedTrack - воспроизводит выбранный трек
func (mp *MediaPlayer) playSelectedTrack(track AudioTrack) {
    fmt.Printf("🎵 Воспроизведение: %s\n", filepath.Base(track.File))
    fmt.Printf("📊 Размер файла: %s\n", formatFileSize(track.Size))
    
    // Запускаем аудиофайл в системном плеере
    mp.launchMediaFile(track.File)
    //Предлагаем показать обложку
    mp.proposeCoverDisplay(track.Id)
    
    fmt.Println("Нажмите Enter для продолжения...")
    mp.scanner.Scan()
}

// proposeCoverDisplay - предлагает пользователю открыть обложку трека
func (mp *MediaPlayer) proposeCoverDisplay(trackId int) {
    // Проверяем наличие обложки
    coverPath, exists := mp.coverMap[trackId]
    if !exists {
        fmt.Println("ℹ️  Для этого трека нет связанной обложки")
        return
    }

    // Проверяем существование файла обложки
    if !utils.FileExists(coverPath) {
        fmt.Println("ℹ️  Обложка для этого трека не найдена на диске")
        return
    }

    fmt.Print("🖼️  Показать обложку? (y/n): ")
    mp.scanner.Scan()
    if strings.ToLower(strings.TrimSpace(mp.scanner.Text())) == "y" {
        mp.launchMediaFile(coverPath)
        info, _ := os.Stat(coverPath)
        fmt.Printf("✅ Обложка открыта! (%s)\n", formatFileSize(info.Size()))
    }
}

// launchMediaFile - запускает файл в системном приложении по умолчанию
func (mp *MediaPlayer) launchMediaFile(filePath string) {
    if !utils.FileExists(filePath) {
        fmt.Printf("❌ Файл не существует: %s\n", filePath)
        return
    }
    var cmd *exec.Cmd
    if runtime.GOOS == "windows" {
        // Windows: используем команду start
        cmd = exec.Command("cmd", "/c", "start", "", filePath)
    } else {
        // Linux используем xdg-open
        cmd = exec.Command("xdg-open", filePath)
    }

    if err := cmd.Start(); err != nil {
        fmt.Printf("❌ Ошибка запуска: %v\n", err)
        fmt.Printf("ℹ️  Попробуйте открыть файл вручную: %s\n", filePath)
        return
    }
    fmt.Printf("▶️  Запущен: %s\n", filepath.Base(filePath))
}

// checkCoverStatus - возвращает статус обложки для трека
func (mp *MediaPlayer) checkCoverStatus(trackId int) string {
    coverPath, exists := mp.coverMap[trackId]
    if !exists {
        return "[без обложки]"
    }
    if utils.FileExists(coverPath) {
        return "[обложка]"
    }
    return "[без обложки]"
}

// checkDownloadStatus - показывает текущий статус загрузки
func (mp *MediaPlayer) checkDownloadStatus() {
    // Считаем аудиофайлы
    audioDir, _ := os.ReadDir(config.AudioFolder)
    // Считаем обложки
    coverDir, _ := os.ReadDir(config.CoverFolder)

    var audioCount, coverCount int
    for _, file := range audioDir {
        if strings.HasSuffix(strings.ToLower(file.Name()), ".mp3") {
            audioCount++
        }
    }
    for _, file := range coverDir {
        if !file.IsDir() {
            coverCount++
        }
    }

    fmt.Println("\n📊 Статус загрузки:")
    fmt.Printf("🎵 Аудио файлов: %d\n", audioCount)
    fmt.Printf("🖼️  Обложек: %d\n", coverCount)
    fmt.Printf("✅ Загрузка завершена: %s\n", map[bool]string{true: "Готово", false: "Ожидание"}[utils.IsDownloadCompleted()])
}

// handleNoFiles - обрабатывает ситуацию отсутствия файлов после завершения загрузки
func (mp *MediaPlayer) handleNoFiles() {
    fmt.Println("\n❌ Аудиофайлы не найдены после завершения загрузки!")
    fmt.Println("🔍 Проверьте:")
    fmt.Printf("1. Наличие файла %s\n", config.InputFilePath)
    fmt.Println("2. Доступность интернета")
    fmt.Println("3. Корректность URL в файле")

    fmt.Print("🔄 Повторить проверку? (y/n): ")
    mp.scanner.Scan()
    if strings.ToLower(strings.TrimSpace(mp.scanner.Text())) != "y" {
        mp.active = false
    }
}

// extractAudioNumber - извлекает номер трека из имени файла
// Ожидаемый формат: audio_001.mp3 → 1
// audio_123.mp3 → 123
func extractAudioNumber(filename string) int {
    if !strings.HasPrefix(filename, "audio_") || !strings.HasSuffix(filename, ".mp3") {
        return 0
    }
    // Извлекаем число: audio_001.mp3 → "001"
    numberStr := filename[6 : len(filename)-4]
    if num, err := strconv.Atoi(numberStr); err == nil {
        return num
    }
    return 0
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
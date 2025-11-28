package models

// MediaUrls - структура для хранения списков URL медиафайлов
type MediaUrls struct {
    AudioUrls []string // Список URL аудиофайлов (.mp3)
    CoverUrls []string // Список URL обложек (.jpg/.png/.gif)
}

// TotalUrls - возвращает общее количество URL
func (m *MediaUrls) TotalUrls() int {
    return len(m.AudioUrls) + len(m.CoverUrls)
}
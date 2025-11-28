package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	// Создаем читатель для ввода пользователя
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Путь к .java файлу: ")
	path, _ := reader.ReadString('\n')
	path = strings.TrimSpace(path)

	// Проверяем существование файла
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Файл не найден:", path)
		return
	}
	
	// Проверяем расширение файла
	if !strings.HasSuffix(strings.ToLower(path), ".java") {
		fmt.Println("Нужен .java файл")
		return
	}

	// Запускаем обработку файла
	processFile(path)
}
// processFile - основная функция обработки .java файла
func processFile(filePath string) {
	// Читаем содержимое файла
	data, _ := os.ReadFile(filePath)
	src := string(data)

	// Извлекаем package и imports, отделяем от тела кода
	packageLine, imports, body := extractHeader(src)
	
	// Удаляем комментарии из тела кода
	body = removeComments(body)
	
	// Минифицируем код (убираем лишние пробелы, переносы)
	body = minify(body)
	
	// Переименовываем главный класс в "A"
	body, oldClass, newClass := renameMainClass(body)
	
	// Обфусцируем идентификаторы (переменные, методы)
	body = obfuscateSafeIdentifiers(body)

	// Собираем результат обратно
	var out strings.Builder
	if packageLine != "" {
		out.WriteString(packageLine + "\n")
	}
	if imports != "" {
		out.WriteString(imports + "\n")
	}
	out.WriteString(body)
	result := out.String()

	// Формируем новое имя файла
	newPath := filePath
	if oldClass != "" {
		// Заменяем имя класса в пути файла
		newPath = strings.Replace(filePath, oldClass+".java", newClass+".java", 1)
	} else {
		// Добавляем суффикс _obf если не нашли класс
		newPath = strings.TrimSuffix(filePath, ".java") + "_obf.java"
	}

	// Записываем результат в новый файл
	os.WriteFile(newPath, []byte(result), 0644)

	// Вычисляем степень сжатия
	reduction := 100.0 * (1.0 - float64(len(result))/float64(len(src)))
	fmt.Printf("\nГотово!\n")
	fmt.Printf("   → %s\n", newPath)
	fmt.Printf("   %d → %d байт (сжатие %.1f%%)\n", len(src), len(result), reduction)
}
// extractHeader разделяет Java код на три части: package, imports и тело
// Возвращает: строка package, блок imports, основное тело кода
func extractHeader(code string) (packageLine, imports, body string) {
	lines := strings.Split(code, "\n")
	var pkg string
	var imp []string
	var bodyLines []string
	inHeader := true // Флаг нахождения в "шапке" файла

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Пропускаем пустые строки и однострочные комментарии в шапке
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			if inHeader {
				continue
			}
		}

		if inHeader {
			// Находим строку package 
			if strings.HasPrefix(trimmed, "package ") && pkg == "" {
				pkg = line + "\n"
				continue
			}
			// Находим все import 
			if strings.HasPrefix(trimmed, "import ") {
				imp = append(imp, line)
				continue
			}
			// Первая не-package, не-import строка завершает шапку
			inHeader = false
		}
		bodyLines = append(bodyLines, line)
	}

	// Собираем imports в одну строку
	if len(imp) > 0 {
		imports = strings.Join(imp, "\n") + "\n"
	}
	body = strings.Join(bodyLines, "\n")
	return
}

// removeComments удаляет все комментарии из Java кода
// Удаляет: многострочные /* */ комментарии и однострочные // комментарии
func removeComments(code string) string {
	// Удаляем многострочные комментарии
	code = regexp.MustCompile(`(?s)/\*.*?\*/`).ReplaceAllString(code, "")
	// Удаляем однострочные комментарии
	code = regexp.MustCompile(`//[^\n]*`).ReplaceAllString(code, "")
	return code
}

// minify удаляет лишние пробельные символы для уменьшения размера кода
// Сохраняет пробелы внутри строковых литералов
func minify(code string) string {
	var result strings.Builder
	inString := false  // Флаг нахождения внутри строки
	quote := rune(0)   // Тип кавычки (' или ") которыми открыта строка

	// Проходим по каждому символу исходного кода
	for _, r := range code {
		if r == '"' || r == '\'' {
			if !inString {
				// Начало строки
				inString = true
				quote = r
			} else if r == quote && result.Len() > 0 && result.String()[result.Len()-1] != '\\' {
				// Конец строки (если кавычка не экранирована)
				inString = false
			}
		}

		// Вне строк: заменяем пробельные символы на одиночные пробелы
		if !inString && (r == ' ' || r == '\t' || r == '\n' || r == '\r') {
			// Добавляем только один пробел вместо серии пробелов
			if result.Len() == 0 || result.String()[result.Len()-1] != ' ' {
				result.WriteRune(' ')
			}
		} else {
			// Внутри строк или для не-пробельных символов: сохраняем как есть
			result.WriteRune(r)
		}
	}

	code = result.String()

	// Убираем лишние пробелы вокруг операторов и разделителей
	ops := []string{"{", "}", "(", ")", "[", "]", ";", ",", "=", "!", ">", "<", "+", "-", "*", "/", "&", "|", ":", "?", "%"}
	for _, op := range ops {
		// Убираем пробелы с обеих сторон оператора
		re := regexp.MustCompile(`\s+` + regexp.QuoteMeta(op) + `\s+`)
		code = re.ReplaceAllString(code, op)
		// Убираем пробелы только слева
		re2 := regexp.MustCompile(`\s+` + regexp.QuoteMeta(op))
		code = re2.ReplaceAllString(code, op)
		// Убираем пробелы только справа
		re3 := regexp.MustCompile(regexp.QuoteMeta(op) + `\s+`)
		code = re3.ReplaceAllString(code, op)
	}

	return strings.TrimSpace(code)
}

// renameMainClass находит и переименовывает главный класс в "A"
// Возвращает: измененный код, старое имя класса, новое имя класса
func renameMainClass(code string) (string, string, string) {
	// Паттерны для поиска объявления класса
	patterns := []string{
		`public\s+class\s+([A-Za-z_][A-Za-z0-9_]*)`, // public class ClassName
		`class\s+([A-Za-z_][A-Za-z0-9_]*)`,          // class ClassName
	}
	
	var old string
	// Ищем первое совпадение с любым из паттернов
	for _, p := range patterns {
		if m := regexp.MustCompile(p).FindStringSubmatch(code); len(m) > 1 {
			old = m[1]
			break
		}
	}
	
	// Если класс не найден, возвращаем исходный код
	if old == "" {
		return code, "", ""
	}
	
	// Переименовываем класс в "A"
	newName := "A"
	// Заменяем все вхождения имени класса (с учетом границ слов)
	code = regexp.MustCompile(`\b` + regexp.QuoteMeta(old) + `\b`).ReplaceAllString(code, newName)
	return code, old, newName
}

// obfuscateSafeIdentifiers заменяет имена переменных и методов на короткие
// Не трогает ключевые слова Java и системные идентификаторы
func obfuscateSafeIdentifiers(code string) string {
	// Множество Java ключевых слов и системных идентификаторов
	keywords := map[string]bool{
		"abstract":true, "assert":true, "boolean":true, "break":true, "byte":true,
		"case":true, "catch":true, "char":true, "class":true, "continue":true,
		"default":true, "do":true, "double":true, "else":true, "enum":true,
		"extends":true, "final":true, "finally":true, "float":true, "for":true,
		"if":true, "implements":true, "import":true, "instanceof":true, "int":true,
		"interface":true, "long":true, "native":true, "new":true, "package":true,
		"private":true, "protected":true, "public":true, "return":true,
		"short":true, "static":true, "super":true, "switch":true,
		"synchronized":true, "this":true, "throw":true, "throws":true,
		"try":true, "void":true, "while":true, "true":true, "false":true, "null":true,
		"System":true, "out":true, "println":true, "printf":true,
    	"String":true, "List":true, "ArrayList":true, "Override":true,
	}

	// Защищенные идентификаторы 
	protected := map[string]bool{
		"System":true, "out":true, "in":true, "err":true,
		"println":true, "print":true, "printf":true,
		"Math":true, "String":true, "List":true, "ArrayList":true,
		"Integer":true, "Long":true, "Collections":true, "Arrays":true,
	}

	// Находим все идентификаторы в коде
	reIdent := regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`)
	allIds := reIdent.FindAllString(code, -1)

	// Собираем кандидатов на замену (исключая ключевые и защищенные слова)
	candidates := make(map[string]struct{})
	for _, id := range allIds {
		if len(id) <= 1 {
			continue // Не заменяем слишком короткие имена
		}
		if keywords[id] || protected[id] {
			continue // Пропускаем ключевые и защищенные слова
		}
		// Не трогаем идентификаторы 
		if regexp.MustCompile(`\b`+regexp.QuoteMeta(id)+`\s*\.`).MatchString(code) ||
		   regexp.MustCompile(`\.\s*`+regexp.QuoteMeta(id)+`\b`).MatchString(code) {
			continue
		}
		candidates[id] = struct{}{}
	}

	// Создаем старых имен на новые короткие
	mapping := make(map[string]string)
	idx := 0
	for id := range candidates {
		mapping[id] = shortName(idx)
		idx++
	}

	// Применяем замены ко всему коду
	result := code
	for old, new := range mapping {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(old) + `\b`)
		result = re.ReplaceAllString(result, new)
	}
	return result
}


// shortName генерирует короткие имена по порядку: a, b, c, ..., z, aa, ab, ...
// n - порядковый номер (0-based)
func shortName(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	if n < 26 {
		// Однобуквенные имена для первых 26 идентификаторов
		return string(letters[n])
	}
	// Двухбуквенные имена для остальных
	return string(letters[(n/26)-1]) + string(letters[n%26])
}
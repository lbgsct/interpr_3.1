package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"regexp"
	"strconv"
)

func main() {
	// Инициализация стека для отслеживания переменных в текущей области видимости
	var scopeStack []map[string]int

	// Открытие файла
	file, err := os.Open("/home/sergey/micro/3/3.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Создание регулярных выражений для поиска строк
	osc := regexp.MustCompile(`{`)                    // Поиск начала области видимости
	csc := regexp.MustCompile(`}`)                    // Поиск конца области видимости
	show := regexp.MustCompile(`ShowVar;`)            // Поиск строки для вывода переменных
	re := regexp.MustCompile(`(\S+)\s*=\s*(\d+);`)    // Поиск строк вида "переменная = значение;"

	// Создание нового сканера для чтения файла построчно
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case osc.MatchString(line):
			// Начало новой области видимости
			scopeStack = append(scopeStack, make(map[string]int))
		case csc.MatchString(line):
			// Конец текущей области видимости
			if len(scopeStack) > 0 {
				scopeStack = scopeStack[:len(scopeStack)-1]
			}
		case show.MatchString(line):
			// Вывод всех видимых переменных из всех уровней области видимости
			fmt.Println("Variables:")
			visibleVars := make(map[string]int)
			for i := len(scopeStack) - 1; i >= 0; i-- {
				for k, v := range scopeStack[i] {
					if _, exists := visibleVars[k]; !exists {
						visibleVars[k] = v
					}
				}
			}
			// Сортируем переменные по именам
			keys := make([]string, 0, len(visibleVars))
			for k := range visibleVars {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Printf("%s = %d\n", k, visibleVars[k])
			}
		case re.MatchString(line):
			// Обработка строки вида "переменная = значение;"
			matches := re.FindStringSubmatch(line)
			if matches != nil {
				key := matches[1]
				number, err := strconv.Atoi(matches[2])
				if err != nil {
					log.Printf("Ошибка преобразования числа: %v", err)
					continue
				}

				if len(scopeStack) > 0 {
					// Добавление переменной в текущую область видимости
					scopeStack[len(scopeStack)-1][key] = number
				}
			}
		}
	}

	// Проверка на ошибки во время сканирования
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
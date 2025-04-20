package sitemap

import (
	"fmt"
	"strings"
)

var translitMap = map[rune]string{
	'Ё': "YO", 'Й': "I", 'Ц': "TS", 'У': "U", 'К': "K", 'Е': "E", 'Н': "N", 'Г': "G", 'Ш': "SH", 'Щ': "SCH", 'З': "Z", 'Х': "H", 'Ъ': "'",
	'ё': "yo", 'й': "i", 'ц': "ts", 'у': "u", 'к': "k", 'е': "e", 'н': "n", 'г': "g", 'ш': "sh", 'щ': "sch", 'з': "z", 'х': "h", 'ъ': "'",
	'Ф': "F", 'Ы': "I", 'В': "V", 'А': "A", 'П': "P", 'Р': "R", 'О': "O", 'Л': "L", 'Д': "D", 'Ж': "ZH", 'Э': "E",
	'ф': "f", 'ы': "i", 'в': "v", 'а': "a", 'п': "p", 'р': "r", 'о': "o", 'л': "l", 'д': "d", 'ж': "zh", 'э': "e",
	'Я': "Ya", 'Ч': "CH", 'С': "S", 'М': "M", 'И': "I", 'Т': "T", 'Ь': "'", 'Б': "B", 'Ю': "YU",
	'я': "ya", 'ч': "ch", 'с': "s", 'м': "m", 'и': "i", 'т': "t", 'ь': "'", 'б': "b", 'ю': "yu",
}

func transliterate(word string) string {
	var result []rune
	for _, char := range word {
		if val, ok := translitMap[char]; ok {
			result = append(result, []rune(val)...)
		} else {
			result = append(result, char)
		}
	}
	return string(result)
}

func generatePathValue(title string, id int64) string {
	itemName := transliterate(title)
	return fmt.Sprintf("%s_%d", strings.ReplaceAll(itemName, " ", "_"), id)
}

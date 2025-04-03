package recomendation

var skipGenres = map[string]struct{}{
	"Ранний доступ":             {},
	"С поддержкой контроллеров": {},
	"macOS":           {},
	"SteamOS и Linux": {},
	"Демоверсии":      {},
	"Бесплатные":      {},
}

var mapGenreNameToID = map[string]string{
	"Бесплатные":           "Free to Play",
	"Экшены":               "Action",
	"Приключенческие игры": "Adventure",
	"Стратегии":            "Strategy",
	"Ролевые игры":         "RPG",
	"Инди":                 "Indie",
	"Многопользовательские игры": "Massively Multiplayer",
	"Казуальные игры":            "Casual",
	"Симуляторы":                 "Simulation",
	"Гонки":                      "Racing",
	"Спорт":                      "Sports",
	"Бухгалтерия":                "Accounting",
	"Работа со звуком":           "Audio Production",
	"Образование":                "Education",
	"Обработка фото":             "Photo Editing",
	"Обучение работе с ПО":       "Software Training",
	"Утилиты":                    "Utilities",
	"Создание видео":             "Video Production",
	"Веб-разработка":             "Web Publishing",
	"Ранний доступ":              "Early Access",
	"С поддержкой контроллеров":  "Controller support",
	"macOS":           "Mac OS X",
	"SteamOS и Linux": "Linux",
	"Демоверсии":      "genre_demos",
}

package steam

// Возвращает не сервисном формате,
// то есть копейки не конвертирует в последние два знака числа
func SteamCalcPrice(amount int64) int64 {
	// Фикс - 30 рублей
	fix := amount + 30
	// Динамическая цена - 14 %
	return int64(float64(fix) * 1.14)
}

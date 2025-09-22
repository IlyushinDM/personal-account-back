package utils

import "strconv"

// ParseUserID преобразует строку в uint64.
// Используется для извлечения ID пользователя из refresh-токена.
func ParseUserID(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

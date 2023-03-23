package pkg

// CheckPassword проверка пароля: 8 символов минимум, наличие знаков арефметических операций
func CheckPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	for _, r := range password {
		if r == '+' || r == '-' || r == '*' || r == '/' {
			return true
		}
	}

	return false
}

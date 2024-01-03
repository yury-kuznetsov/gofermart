package validation

// IsValidLuhn проверяет число по алгоритму Луна
func IsValidLuhn(idNumber string) bool {
	if idNumber == "" {
		return false
	}

	length := len(idNumber)
	sum := int(idNumber[length-1] - '0')
	odd := length % 2
	for i := 0; i < length-1; i++ {
		digit := int(idNumber[i] - '0')
		if i%2 == odd {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}

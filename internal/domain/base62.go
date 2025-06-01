package domain

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func EncodeBase62(num int64) string {
	if num == 0 {
		return string(alphabet[0])
	}
	result := make([]byte, 0)
	for num > 0 {
		rem := num % 62
		num = num / 62
		result = append([]byte{alphabet[rem]}, result...)
	}
	return string(result)
}

func DecodeBase62(code string) int64 {
	var num int64
	for _, c := range code {
		pos := int64(-1)
		for i, a := range alphabet {
			if a == c {
				pos = int64(i)
				break
			}
		}
		if pos == -1 {
			return -1
		}
		num = num*62 + pos
	}
	return num
}

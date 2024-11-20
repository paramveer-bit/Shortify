package helper

const base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

func ToBase64(num int64) string {
	if num == 0 {
		return "A" // 'A' represents 0 in base 64
	}

	result := ""
	for num > 0 {
		remainder := num % 64
		result = string(base64Chars[remainder]) + result
		num = num / 64
	}
	return result
}

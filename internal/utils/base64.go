package utils

import "encoding/base64"

func EncodeHTML(html string) string {
	return base64.StdEncoding.EncodeToString([]byte(html))
}

// DecodeHTML 將Base64字符串解碼為HTML
func DecodeHTML(encodedHTML string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encodedHTML)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

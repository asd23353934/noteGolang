package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode/utf8"
)

// EncodeRawData 將原始數據編碼為安全的 UTF-8 字符串
func EncodeRawData(input string) string {
	// 將輸入視為原始字節
	rawBytes := []byte(input)

	// 使用 Base64 編碼
	encoded := base64.StdEncoding.EncodeToString(rawBytes)

	// 添加一個標記，表示這是編碼後的數據
	return "BASE64:" + encoded
}

// DecodeRawData 將編碼後的字符串解碼回原始數據
func DecodeRawData(input string) (string, error) {
	// 檢查是否是我們編碼的數據
	if !strings.HasPrefix(input, "BASE64:") {
		return input, nil // 如果不是，直接返回原始輸入
	}

	// 移除前綴
	encodedData := strings.TrimPrefix(input, "BASE64:")

	// 解碼 Base64
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", err
	}

	// 將解碼後的字節轉換回字符串
	return string(decodedBytes), nil
}

func PrepareHTMLForCassandra(html string) string {
	// 1. 移除或转义可能导致问题的字符
	html = strings.ReplaceAll(html, "'", "''") // 转义单引号

	// 2. 确保文本是有效的UTF-8
	if !utf8.ValidString(html) {
		html = strings.ToValidUTF8(html, "")
	}

	// 3. 如果内容过长，可以考虑截断
	maxLength := 65535 // Cassandra text类型的最大长度
	if len(html) > maxLength {
		html = html[:maxLength]
	}

	return html
}

func ValidateAndFixUTF8(s string) (string, error) {
	var buf bytes.Buffer
	var diagnostics bytes.Buffer

	for i, r := range s {
		if r == utf8.RuneError {
			runeBytes := []byte(s[i:])
			r, size := utf8.DecodeRune(runeBytes)
			if r == utf8.RuneError {
				diagnostics.WriteString(fmt.Sprintf("Invalid UTF-8 sequence at position %d: %s\n", i, hex.EncodeToString(runeBytes[:size])))
				buf.WriteRune('�') // 替换为 Unicode 替换字符
			} else {
				buf.WriteRune(r)
			}
		} else {
			buf.WriteRune(r)
		}
	}

	if diagnostics.Len() > 0 {
		return buf.String(), fmt.Errorf("UTF-8 validation issues found:\n%s", diagnostics.String())
	}

	return buf.String(), nil
}

func SanitizeForCassandra(s string) string {
	// 步骤 1: 移除或替换可能导致问题的字符
	s = strings.Map(func(r rune) rune {
		if r > 127 || r < 32 {
			return -1 // 移除非 ASCII 字符和控制字符
		}
		return r
	}, s)

	// 步骤 2: 转义单引号和反斜杠
	s = strings.ReplaceAll(s, "'", "''")
	s = strings.ReplaceAll(s, "\\", "\\\\")

	// 步骤 3: Base64 编码
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func DecodeFromCassandra(s string) (string, error) {
	// 步骤 1: Base64 解码
	decodedBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("Base64 解码失败: %v", err)
	}

	// 步骤 2: 转换为字符串
	decoded := string(decodedBytes)

	// 步骤 3: 还原转义字符
	decoded = strings.ReplaceAll(decoded, "''", "'")
	decoded = strings.ReplaceAll(decoded, "\\\\", "\\")

	return decoded, nil
}

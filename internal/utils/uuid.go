package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateUUID() (string, error) {
	uuidV4, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return uuidV4.String(), err
}

func GenerateUUIDToken(length ...int) string {
	u, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	count := 7
	if len(length) != 0 {
		count = length[0]
	}

	return strings.Replace(u.String(), "-", "", -1)[:count]
}

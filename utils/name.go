package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomName(length int) string {
	// Caracteres permitidos
	chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Calcula o tamanho da string de caracteres permitidos
	charsLength := big.NewInt(int64(len(chars)))

	// Cria um slice para armazenar os bytes aleatórios
	randomBytes := make([]byte, length)

	// Gera bytes aleatórios usando crypto/rand
	for i := 0; i < length; i++ {
		randomByte, err := rand.Int(rand.Reader, charsLength)
		if err != nil {
			panic(err)
		}

		// Converte o byte gerado para um índice válido no conjunto de caracteres
		randomBytes[i] = chars[randomByte.Int64()]
	}

	// Converte o slice de bytes para uma string
	randomName := string(randomBytes)

	return randomName
}

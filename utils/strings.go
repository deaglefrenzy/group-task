package utils

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateString(numWords int) string {
	if numWords <= 0 {
		return ""
	}

	loremIpsum := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

	words := strings.Fields(loremIpsum)
	if len(words) == 0 {
		return ""
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	generatedWords := make([]string, numWords)

	for i := 0; i < numWords; i++ {
		randomIndex := r.Intn(len(words))
		generatedWords[i] = words[randomIndex]
	}

	return strings.Join(generatedWords, " ")
}

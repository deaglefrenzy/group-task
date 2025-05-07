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

	adjectives := []string{
		"quick", "fast", "slow", "smart", "sharp", "dull", "bright", "dark",
		"light", "heavy", "soft", "hard", "loud", "quiet", "new", "old",
		"young", "fresh", "stale", "clean", "dirty", "rich", "poor", "fun",
		"boring", "hot", "cold", "warm", "cool", "big", "small", "short",
		"tall", "long", "wide", "narrow", "happy", "sad", "good", "bad",
		"nice", "mean", "kind", "cruel", "brave", "shy", "calm", "angry",
		"lucky", "neat", "messy", "rare", "common", "simple", "complex",
		"strong", "weak", "tough", "gentle", "sweet", "sour", "bitter",
		"salty", "smooth", "rough", "clear", "fuzzy", "wild", "tame",
		"major", "minor", "huge", "tiny", "plump", "lean", "thick", "thin",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	generatedWords := make([]string, numWords)

	for i := 0; i < numWords; i++ {
		randomIndex := r.Intn(len(adjectives))
		generatedWords[i] = adjectives[randomIndex]
	}

	return strings.Join(generatedWords, " ")
}

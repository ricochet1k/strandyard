package idgen

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

const (
	// Base36 alphabet (0-9, a-z)
	base36Chars = "0123456789abcdefghijklmnopqrstuvwxyz"
	tokenLength = 6
)

// GenerateToken creates a cryptographically secure 4-character base36 token
func GenerateToken() (string, error) {
	token := make([]byte, tokenLength)
	max := big.NewInt(int64(len(base36Chars)))

	for i := 0; i < tokenLength; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		token[i] = base36Chars[n.Int64()]
	}

	return string(token), nil
}

// Slugify converts a title to a URL-safe slug
// Example: "Implement Parser" -> "implement-parser"
func Slugify(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace non-alphanumeric characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Collapse multiple hyphens
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Limit length to reasonable size (e.g., 50 chars)
	if len(slug) > 50 {
		slug = slug[:50]
		// Trim trailing hyphen if we cut in the middle
		slug = strings.TrimRight(slug, "-")
	}

	return slug
}

// GenerateID creates a complete task ID with prefix, token, and slug
// Example: GenerateID("T", "Implement Parser") -> "T3k7x-implement-parser"
func GenerateID(prefix, title string) (string, error) {
	token, err := GenerateToken()
	if err != nil {
		return "", err
	}

	slug := Slugify(title)
	if slug == "" {
		return "", fmt.Errorf("invalid title: produces empty slug")
	}

	return fmt.Sprintf("%s%s-%s", prefix, token, slug), nil
}

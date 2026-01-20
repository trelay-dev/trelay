package slug

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"strings"

	"github.com/aftaab/trelay/internal/core/domain"
)

const (
	// Base62 charset for URL-safe slugs
	charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	MinSlugLength = 4
	MaxSlugLength = 32
	DefaultLength = 6
)

var (
	// validSlugPattern matches alphanumeric slugs with optional hyphens/underscores
	validSlugPattern = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9_-]*[a-zA-Z0-9])?$|^[a-zA-Z0-9]$`)

	// reservedSlugs are paths that cannot be used as link slugs
	reservedSlugs = map[string]bool{
		"api":      true,
		"admin":    true,
		"login":    true,
		"logout":   true,
		"register": true,
		"health":   true,
		"healthz":  true,
		"metrics":  true,
		"static":   true,
		"assets":   true,
		"favicon":  true,
	}
)

// Generator handles slug generation and validation.
type Generator struct {
	length int
}

// NewGenerator creates a new slug generator with the specified default length.
func NewGenerator(length int) *Generator {
	if length < MinSlugLength {
		length = MinSlugLength
	}
	if length > MaxSlugLength {
		length = MaxSlugLength
	}
	return &Generator{length: length}
}

// Generate creates a new random slug using Base62 encoding.
func (g *Generator) Generate() (string, error) {
	return generateRandomSlug(g.length)
}

// GenerateWithLength creates a random slug of specific length.
func (g *Generator) GenerateWithLength(length int) (string, error) {
	if length < MinSlugLength {
		length = MinSlugLength
	}
	if length > MaxSlugLength {
		length = MaxSlugLength
	}
	return generateRandomSlug(length)
}

func generateRandomSlug(length int) (string, error) {
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

// Validate checks if a custom slug is valid.
func (g *Generator) Validate(slug string) error {
	if len(slug) < MinSlugLength {
		return domain.ErrSlugTooShort
	}
	if len(slug) > MaxSlugLength {
		return domain.ErrSlugTooLong
	}
	if !validSlugPattern.MatchString(slug) {
		return domain.ErrSlugInvalid
	}
	if IsReserved(slug) {
		return domain.ErrSlugInvalid
	}
	return nil
}

// Normalize cleans up a slug (lowercase, trim spaces).
func Normalize(slug string) string {
	return strings.ToLower(strings.TrimSpace(slug))
}

// IsReserved checks if a slug is reserved for system use.
func IsReserved(slug string) bool {
	return reservedSlugs[strings.ToLower(slug)]
}

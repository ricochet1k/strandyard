package idgen

import (
	"regexp"
	"strings"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	// Test basic token generation
	token, err := GenerateToken()
	if err != nil {
		t.Fatalf("GenerateToken() returned error: %v", err)
	}

	// Check length
	if len(token) != tokenLength {
		t.Errorf("GenerateToken() = %q, want length %d, got %d", token, tokenLength, len(token))
	}

	// Check that all characters are valid base36
	for _, c := range token {
		if !strings.ContainsRune(base36Chars, c) {
			t.Errorf("GenerateToken() = %q contains invalid character %c", token, c)
		}
	}
}

func TestGenerateTokenUniqueness(t *testing.T) {
	// Generate many tokens and check for uniqueness
	seen := make(map[string]bool)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		token, err := GenerateToken()
		if err != nil {
			t.Fatalf("GenerateToken() iteration %d returned error: %v", i, err)
		}

		if seen[token] {
			t.Errorf("GenerateToken() generated duplicate token %q after %d iterations", token, i)
		}
		seen[token] = true
	}

	// We should have generated 1000 unique tokens
	if len(seen) != iterations {
		t.Errorf("Generated %d unique tokens, want %d", len(seen), iterations)
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Implement Parser", "implement-parser"},
		{"Add Unit Tests", "add-unit-tests"},
		{"Update CLI Commands", "update-cli-commands"},
		{"Fix Bug #123", "fix-bug-123"},
		{"Multi   Space   Test", "multi-space-test"},
		{"Trailing hyphens---", "trailing-hyphens"},
		{"---Leading hyphens", "leading-hyphens"},
		{"Special!@#$%^&*()Chars", "special-chars"},
		{"Already-slugified", "already-slugified"},
		{"CamelCaseTitle", "camelcasetitle"},
		{"", ""},
		{"   spaces   ", "spaces"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Slugify(tt.input)
			if got != tt.want {
				t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSlugifyLength(t *testing.T) {
	// Test that very long titles are truncated
	longTitle := strings.Repeat("word ", 20) // "word word word ..."
	slug := Slugify(longTitle)

	if len(slug) > 50 {
		t.Errorf("Slugify() produced slug of length %d, want <= 50", len(slug))
	}

	// Should not end with hyphen
	if strings.HasSuffix(slug, "-") {
		t.Errorf("Slugify() produced slug ending with hyphen: %q", slug)
	}
}

func TestGenerateID(t *testing.T) {
	tests := []struct {
		prefix string
		title  string
		wantRe string // regex pattern to match
	}{
		{"T", "Implement Parser", `^T[0-9a-z]{4,6}-implement-parser$`},
		{"E", "Epic Task", `^E[0-9a-z]{4,6}-epic-task$`},
		{"D", "Design Doc", `^D[0-9a-z]{4,6}-design-doc$`},
	}

	for _, tt := range tests {
		t.Run(tt.prefix+"-"+tt.title, func(t *testing.T) {
			id, err := GenerateID(tt.prefix, tt.title)
			if err != nil {
				t.Fatalf("GenerateID(%q, %q) returned error: %v", tt.prefix, tt.title, err)
			}

			matched, err := regexp.MatchString(tt.wantRe, id)
			if err != nil {
				t.Fatalf("regexp.MatchString failed: %v", err)
			}

			if !matched {
				t.Errorf("GenerateID(%q, %q) = %q, want to match pattern %q", tt.prefix, tt.title, id, tt.wantRe)
			}

			// Check that it starts with prefix
			if !strings.HasPrefix(id, tt.prefix) {
				t.Errorf("GenerateID(%q, %q) = %q, want prefix %q", tt.prefix, tt.title, id, tt.prefix)
			}
		})
	}
}

func TestGenerateIDEmptyTitle(t *testing.T) {
	// Test that empty or whitespace-only titles return error
	tests := []string{"", "   ", "!!!"}

	for _, title := range tests {
		_, err := GenerateID("T", title)
		if err == nil {
			t.Errorf("GenerateID(%q, %q) should return error for invalid title", "T", title)
		}
	}
}

func TestGenerateIDUniqueness(t *testing.T) {
	// Generate multiple IDs with same prefix and title
	// They should all be unique due to random token
	seen := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		id, err := GenerateID("T", "Test Task")
		if err != nil {
			t.Fatalf("GenerateID() iteration %d returned error: %v", i, err)
		}

		if seen[id] {
			t.Errorf("GenerateID() generated duplicate ID %q after %d iterations", id, i)
		}
		seen[id] = true
	}

	// All IDs should be unique
	if len(seen) != iterations {
		t.Errorf("Generated %d unique IDs, want %d", len(seen), iterations)
	}
}

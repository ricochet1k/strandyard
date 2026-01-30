package task

import (
	"fmt"
	"hash/fnv"
	"strings"
	"testing"
)

func testToken(parts ...string) string {
	h := fnv.New32a()
	for _, part := range parts {
		_, _ = h.Write([]byte(part))
	}
	return fmt.Sprintf("%08x", h.Sum32())[:6]
}

func testRoleName(t *testing.T, suffix string) string {
	name := strings.TrimSpace(t.Name())
	return "role-" + testToken("role", name, suffix)
}

func testTypeName(t *testing.T, suffix string) string {
	name := strings.TrimSpace(t.Name())
	return "type-" + testToken("type", name, suffix)
}

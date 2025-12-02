package hash


import (
	"testing"
)

func TestHashFunctions(t *testing.T) {
	data := []byte("Hello, World!")
	hash := Hash64(data)
	t.Logf("Hash64: %v", hash)
	t.Logf("Hash64: %x", hash)
}

func TestHashFile(t *testing.T) {
	filePath := "test.md"
	hash, err := HashFile(filePath)
	if err != nil {
		t.Fatalf("Failed to hash file: %v", err)
	}
	t.Logf("MD5 Hash of file %s: %v", filePath, hash)
}
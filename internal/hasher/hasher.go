package hasher

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

// MD5Hash returns the MD5 hex digest of the input.
func MD5Hash(input string) string {
	sum := md5.Sum([]byte(input))
	return fmt.Sprintf("%x", sum)
}

// SHA1Hash returns the SHA-1 hex digest of the input.
func SHA1Hash(input string) string {
	sum := sha1.Sum([]byte(input))
	return fmt.Sprintf("%x", sum)
}

// SHA256Hash returns the SHA-256 hex digest of the input.
func SHA256Hash(input string) string {
	sum := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", sum)
}

// SHA512Hash returns the SHA-512 hex digest of the input.
func SHA512Hash(input string) string {
	sum := sha512.Sum512([]byte(input))
	return fmt.Sprintf("%x", sum)
}

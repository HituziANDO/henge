package cmd_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	// Build the binary once before all tests
	dir, err := os.MkdirTemp("", "henge-test-*")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}
	defer os.RemoveAll(dir)

	binaryPath = filepath.Join(dir, "henge")
	build := exec.Command("go", "build", "-o", binaryPath, ".")
	build.Dir = filepath.Join(mustGetWd(), "..")
	build.Env = append(os.Environ(), "CGO_ENABLED=0")
	if out, err := build.CombinedOutput(); err != nil {
		panic("failed to build henge: " + err.Error() + "\n" + string(out))
	}

	os.Exit(m.Run())
}

func mustGetWd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}

// runHenge executes the henge binary with the given args and stdin input.
func runHenge(t *testing.T, stdinInput string, args ...string) (string, string, error) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	cmd.Stdin = strings.NewReader(stdinInput)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// --- Encode tests ---

func TestEncodeBase64(t *testing.T) {
	stdout, _, err := runHenge(t, "hello", "encode", "base64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "aGVsbG8=" {
		t.Errorf("encode base64: got %q, want %q", got, "aGVsbG8=")
	}
}

func TestEncodeBase64WithArgs(t *testing.T) {
	stdout, _, err := runHenge(t, "", "encode", "base64", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "aGVsbG8=" {
		t.Errorf("encode base64 with arg: got %q, want %q", got, "aGVsbG8=")
	}
}

// --- Decode tests ---

func TestDecodeBase64(t *testing.T) {
	stdout, _, err := runHenge(t, "aGVsbG8=", "decode", "base64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "hello" {
		t.Errorf("decode base64: got %q, want %q", got, "hello")
	}
}

func TestDecodeBase64WithArgs(t *testing.T) {
	stdout, _, err := runHenge(t, "", "decode", "base64", "aGVsbG8=")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "hello" {
		t.Errorf("decode base64 with arg: got %q, want %q", got, "hello")
	}
}

// --- Hash tests ---

func TestHashSHA256(t *testing.T) {
	stdout, _, err := runHenge(t, "hello", "hash", "sha256")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	expected := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if got != expected {
		t.Errorf("hash sha256: got %q, want %q", got, expected)
	}
}

func TestHashMD5(t *testing.T) {
	stdout, _, err := runHenge(t, "hello", "hash", "md5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	expected := "5d41402abc4b2a76b9719d911017c592"
	if got != expected {
		t.Errorf("hash md5: got %q, want %q", got, expected)
	}
}

func TestHashSHA1(t *testing.T) {
	stdout, _, err := runHenge(t, "hello", "hash", "sha1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	expected := "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"
	if got != expected {
		t.Errorf("hash sha1: got %q, want %q", got, expected)
	}
}

func TestHashSHA512(t *testing.T) {
	stdout, _, err := runHenge(t, "hello", "hash", "sha512")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	expected := "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043"
	if got != expected {
		t.Errorf("hash sha512: got %q, want %q", got, expected)
	}
}

// --- Pipe round-trip tests ---

func TestPipeEncodeDecodePoundTrip(t *testing.T) {
	// Encode then decode should return the original input
	encOut, _, err := runHenge(t, "hello world", "encode", "base64")
	if err != nil {
		t.Fatalf("encode step failed: %v", err)
	}

	decOut, _, err := runHenge(t, encOut, "decode", "base64")
	if err != nil {
		t.Fatalf("decode step failed: %v", err)
	}
	got := strings.TrimSpace(decOut)
	if got != "hello world" {
		t.Errorf("round-trip encode/decode: got %q, want %q", got, "hello world")
	}
}

func TestPipeHashConsistency(t *testing.T) {
	// Hashing the same input twice via pipe should produce the same result
	out1, _, err := runHenge(t, "test input", "hash", "sha256")
	if err != nil {
		t.Fatalf("first hash failed: %v", err)
	}
	out2, _, err := runHenge(t, "test input", "hash", "sha256")
	if err != nil {
		t.Fatalf("second hash failed: %v", err)
	}
	if out1 != out2 {
		t.Errorf("hash consistency: %q != %q", out1, out2)
	}
}

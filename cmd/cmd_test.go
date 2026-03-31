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

// --- Image encode/decode tests ---

// Minimal 1x1 pixel PNG
var testPNGData = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
	0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, 0xde, 0x00, 0x00, 0x00,
	0x0c, 0x49, 0x44, 0x41, 0x54, 0x08, 0xd7, 0x63, 0xf8, 0xcf, 0xc0, 0x00,
	0x00, 0x00, 0x02, 0x00, 0x01, 0xe2, 0x21, 0xbc, 0x33, 0x00, 0x00, 0x00,
	0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
}

// writeTempPNG creates a temporary PNG file and returns its path.
func writeTempPNG(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.png")
	if err := os.WriteFile(p, testPNGData, 0644); err != nil {
		t.Fatalf("failed to write temp PNG: %v", err)
	}
	return p
}

func TestEncodeImage(t *testing.T) {
	pngFile := writeTempPNG(t)

	stdout, _, err := runHenge(t, "", "encode", "image", pngFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := strings.TrimSpace(stdout)
	want := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAADElEQVQI12P4z8AAAAACAAHiIbwzAAAAAElFTkSuQmCC"
	if got != want {
		t.Errorf("encode image: got %q, want %q", got, want)
	}
}

func TestEncodeImageDataURI(t *testing.T) {
	pngFile := writeTempPNG(t)

	stdout, _, err := runHenge(t, "", "encode", "image", pngFile, "--data-uri")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := strings.TrimSpace(stdout)
	wantPrefix := "data:image/png;base64,"
	if !strings.HasPrefix(got, wantPrefix) {
		t.Errorf("encode image --data-uri: output %q does not start with %q", got, wantPrefix)
	}
}

func TestEncodeImageWrap(t *testing.T) {
	pngFile := writeTempPNG(t)

	stdout, _, err := runHenge(t, "", "encode", "image", pngFile, "--wrap", "20")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := strings.TrimSpace(stdout)
	lines := strings.Split(got, "\n")
	for i, line := range lines {
		// All lines except the last should be exactly 20 chars
		if i < len(lines)-1 && len(line) != 20 {
			t.Errorf("line %d has length %d, want 20: %q", i, len(line), line)
		}
	}
	if len(lines) < 2 {
		t.Errorf("expected wrapped output to have multiple lines, got %d", len(lines))
	}
}

func TestDecodeImage(t *testing.T) {
	// Prepare base64 of test PNG
	b64 := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAADElEQVQI12P4z8AAAAACAAHiIbwzAAAAAElFTkSuQmCC"
	dir := t.TempDir()
	outPath := filepath.Join(dir, "restored.png")

	_, _, err := runHenge(t, "", "decode", "image", b64, "-o", outPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !bytes.Equal(got, testPNGData) {
		t.Errorf("decoded image does not match original PNG data")
	}
}

func TestDecodeImageDataURI(t *testing.T) {
	b64 := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAADElEQVQI12P4z8AAAAACAAHiIbwzAAAAAElFTkSuQmCC"
	dataURI := "data:image/png;base64," + b64
	dir := t.TempDir()
	outPath := filepath.Join(dir, "restored.png")

	_, _, err := runHenge(t, "", "decode", "image", dataURI, "-o", outPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !bytes.Equal(got, testPNGData) {
		t.Errorf("decoded data URI image does not match original PNG data")
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	pngFile := writeTempPNG(t)

	// Step 1: encode image to base64
	stdout, _, err := runHenge(t, "", "encode", "image", pngFile)
	if err != nil {
		t.Fatalf("encode step failed: %v", err)
	}
	encoded := strings.TrimSpace(stdout)

	// Step 2: decode base64 back to image
	dir := t.TempDir()
	outPath := filepath.Join(dir, "roundtrip.png")
	_, _, err = runHenge(t, "", "decode", "image", encoded, "-o", outPath)
	if err != nil {
		t.Fatalf("decode step failed: %v", err)
	}

	// Step 3: compare with original
	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !bytes.Equal(got, testPNGData) {
		t.Errorf("round-trip encode/decode: decoded file does not match original")
	}
}

// --- Time conversion tests ---

func TestTimeUnix(t *testing.T) {
	stdout, _, err := runHenge(t, "", "time", "unix", "2025-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "1735689600" {
		t.Errorf("time unix: got %q, want %q", got, "1735689600")
	}
}

func TestTimeUnixFromStdin(t *testing.T) {
	stdout, _, err := runHenge(t, "2025-01-01T00:00:00Z", "time", "unix")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "1735689600" {
		t.Errorf("time unix from stdin: got %q, want %q", got, "1735689600")
	}
}

func TestTimeUnixMillis(t *testing.T) {
	stdout, _, err := runHenge(t, "", "time", "unix", "--millis", "2025-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "1735689600000" {
		t.Errorf("time unix --millis: got %q, want %q", got, "1735689600000")
	}
}

func TestTimeUnixWithTimezone(t *testing.T) {
	stdout, _, err := runHenge(t, "", "time", "unix", "--timezone", "Asia/Tokyo", "2025-01-01 09:00:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "1735689600" {
		t.Errorf("time unix with timezone: got %q, want %q", got, "1735689600")
	}
}

func TestTimeDate(t *testing.T) {
	stdout, _, err := runHenge(t, "", "time", "date", "1735689600")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "2025-01-01T00:00:00Z" {
		t.Errorf("time date: got %q, want %q", got, "2025-01-01T00:00:00Z")
	}
}

func TestTimeDateFromStdin(t *testing.T) {
	stdout, _, err := runHenge(t, "1735689600", "time", "date")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "2025-01-01T00:00:00Z" {
		t.Errorf("time date from stdin: got %q, want %q", got, "2025-01-01T00:00:00Z")
	}
}

func TestTimeDateWithTimezone(t *testing.T) {
	stdout, _, err := runHenge(t, "", "time", "date", "--timezone", "Asia/Tokyo", "1735689600")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "2025-01-01T09:00:00+09:00" {
		t.Errorf("time date with timezone: got %q, want %q", got, "2025-01-01T09:00:00+09:00")
	}
}

func TestTimeDateWithFormat(t *testing.T) {
	stdout, _, err := runHenge(t, "", "time", "date", "--format", "date", "1735689600")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "2025-01-01" {
		t.Errorf("time date with format: got %q, want %q", got, "2025-01-01")
	}
}

func TestTimeDateMillisAutoDetect(t *testing.T) {
	stdout, _, err := runHenge(t, "", "time", "date", "1735689600000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "2025-01-01T00:00:00Z" {
		t.Errorf("time date millis auto-detect: got %q, want %q", got, "2025-01-01T00:00:00Z")
	}
}

func TestTimeAutoTimestamp(t *testing.T) {
	stdout, _, err := runHenge(t, "", "time", "1735689600")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "2025-01-01T00:00:00Z" {
		t.Errorf("time auto (timestamp): got %q, want %q", got, "2025-01-01T00:00:00Z")
	}
}

func TestTimeAutoDate(t *testing.T) {
	stdout, _, err := runHenge(t, "", "time", "2025-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "1735689600" {
		t.Errorf("time auto (date): got %q, want %q", got, "1735689600")
	}
}

func TestTimeAutoFromStdin(t *testing.T) {
	stdout, _, err := runHenge(t, "1735689600", "time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "2025-01-01T00:00:00Z" {
		t.Errorf("time auto from stdin: got %q, want %q", got, "2025-01-01T00:00:00Z")
	}
}

func TestTimeAutoMillis(t *testing.T) {
	stdout, _, err := runHenge(t, "", "time", "1735689600000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(stdout)
	if got != "2025-01-01T00:00:00Z" {
		t.Errorf("time auto (millis): got %q, want %q", got, "2025-01-01T00:00:00Z")
	}
}

func TestTimeRoundTrip(t *testing.T) {
	// Step 1: convert date string to unix timestamp
	stdout1, _, err := runHenge(t, "", "time", "unix", "2025-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("unix step failed: %v", err)
	}
	unixStr := strings.TrimSpace(stdout1)

	// Step 2: convert unix timestamp back to date string
	stdout2, _, err := runHenge(t, "", "time", "date", unixStr)
	if err != nil {
		t.Fatalf("date step failed: %v", err)
	}
	got := strings.TrimSpace(stdout2)

	// Step 3: verify round-trip
	if got != "2025-01-01T00:00:00Z" {
		t.Errorf("round-trip time unix/date: got %q, want %q", got, "2025-01-01T00:00:00Z")
	}
}

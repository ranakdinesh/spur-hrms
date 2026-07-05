package pbkdf2verifier

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"golang.org/x/crypto/pbkdf2"
)

func TestVerifyDollarSHA256Base64(t *testing.T) {
	hash := pbkdf2.Key([]byte("secret"), []byte("tenant-salt"), 12000, 32, sha256.New)
	encoded := "pbkdf2_sha256$12000$tenant-salt$" + base64.StdEncoding.EncodeToString(hash)

	ok, err := Verify("secret", encoded)
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
	if !ok {
		t.Fatal("Verify returned false for a valid SHA-256 PBKDF2 hash")
	}
}

func TestVerifyColonSHA1Hex(t *testing.T) {
	hash := pbkdf2.Key([]byte("secret"), []byte("salt"), 8000, 20, sha1.New)
	encoded := "pbkdf2:sha1:8000:salt:" + hex.EncodeToString(hash)

	ok, err := Verify("secret", encoded)
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
	if !ok {
		t.Fatal("Verify returned false for a valid SHA-1 PBKDF2 hash")
	}
}

func TestVerifyWrongPassword(t *testing.T) {
	hash := pbkdf2.Key([]byte("secret"), []byte("salt"), 10000, 32, sha256.New)
	encoded := "pbkdf2_sha256$10000$salt$" + base64.StdEncoding.EncodeToString(hash)

	ok, err := Verify("wrong", encoded)
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
	if ok {
		t.Fatal("Verify returned true for an invalid password")
	}
}

func TestVerifyRejectsUnsupportedAlgorithm(t *testing.T) {
	_, err := Verify("secret", "pbkdf2_md5$10000$salt$abcdef")
	if !errors.Is(err, ErrUnsupportedAlgorithm) {
		t.Fatalf("expected ErrUnsupportedAlgorithm, got %v", err)
	}
}

func TestVerifyRejectsInvalidIterations(t *testing.T) {
	_, err := Verify("secret", "pbkdf2_sha256$0$salt$abcdef")
	if !errors.Is(err, ErrInvalidIterationCount) {
		t.Fatalf("expected ErrInvalidIterationCount, got %v", err)
	}
}

func TestVerifyRejectsEmptyPassword(t *testing.T) {
	_, err := Verify(" ", "pbkdf2_sha256$10000$salt$abcdef")
	if !errors.Is(err, ErrInvalidPassword) {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}
}

func TestVerifyTable(t *testing.T) {
	sha256Hash := pbkdf2.Key([]byte("secret"), []byte("salt"), 10000, 32, sha256.New)
	sha1Hash := pbkdf2.Key([]byte("secret"), []byte("salt"), 8000, 20, sha1.New)
	tests := []struct {
		name    string
		encoded string
		want    bool
		wantErr error
	}{
		{"dollar sha256 base64", "pbkdf2_sha256$10000$salt$" + base64.StdEncoding.EncodeToString(sha256Hash), true, nil},
		{"colon sha256 hex", "pbkdf2:sha256:10000:salt:" + hex.EncodeToString(sha256Hash), true, nil},
		{"colon sha1 uppercase hex", "pbkdf2:sha1:8000:salt:" + strings.ToUpper(hex.EncodeToString(sha1Hash)), true, nil},
		{"wrong password", "pbkdf2_sha256$10000$salt$" + base64.StdEncoding.EncodeToString(sha256Hash), false, nil},
		{"bad format", "pbkdf2_sha256:10000:salt", false, ErrInvalidLegacyHash},
		{"bad hash encoding", "pbkdf2_sha256$10000$salt$not-valid-encoding", false, ErrInvalidLegacyHash},
		{"bad iterations", "pbkdf2_sha256$abc$salt$" + base64.StdEncoding.EncodeToString(sha256Hash), false, ErrInvalidIterationCount},
		{"unsupported", "pbkdf2_md5$10000$salt$" + base64.StdEncoding.EncodeToString(sha256Hash), false, ErrUnsupportedAlgorithm},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password := "secret"
			if tt.name == "wrong password" {
				password = "wrong"
			}
			got, err := Verify(password, tt.encoded)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseLegacyPBKDF2Hash(t *testing.T) {
	hash := pbkdf2.Key([]byte("secret"), []byte("salt"), 10000, 32, sha256.New)
	tests := []struct {
		name           string
		encoded        string
		wantAlgorithm  string
		wantIterations int
		wantSalt       string
	}{
		{"dollar format", "pbkdf2_sha256$10000$salt$" + base64.StdEncoding.EncodeToString(hash), "pbkdf2_sha256", 10000, "salt"},
		{"colon format", "pbkdf2:sha256:12000:tenant-salt:" + hex.EncodeToString(hash), "pbkdf2_sha256", 12000, "tenant-salt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.encoded)
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}
			if got.Algorithm != tt.wantAlgorithm || got.Iterations != tt.wantIterations || string(got.Salt) != tt.wantSalt {
				t.Fatalf("unexpected parsed hash: %#v", got)
			}
			if len(got.Hash) == 0 {
				t.Fatal("expected decoded hash bytes")
			}
		})
	}
}

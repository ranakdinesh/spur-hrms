package pbkdf2verifier

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

var (
	ErrInvalidPassword       = errors.New("password is required")
	ErrInvalidLegacyHash     = errors.New("legacy pbkdf2 hash is invalid")
	ErrUnsupportedAlgorithm  = errors.New("legacy pbkdf2 algorithm is unsupported")
	ErrInvalidIterationCount = errors.New("legacy pbkdf2 iteration count is invalid")
)

type ParsedHash struct {
	Algorithm  string
	Iterations int
	Salt       []byte
	Hash       []byte
}

func Verify(password string, encoded string) (bool, error) {
	if strings.TrimSpace(password) == "" {
		return false, ErrInvalidPassword
	}
	parsed, err := Parse(encoded)
	if err != nil {
		return false, err
	}
	hashFunc, err := hashForAlgorithm(parsed.Algorithm)
	if err != nil {
		return false, err
	}
	actual := pbkdf2.Key([]byte(password), parsed.Salt, parsed.Iterations, len(parsed.Hash), hashFunc)
	return subtle.ConstantTimeCompare(actual, parsed.Hash) == 1, nil
}

func Parse(encoded string) (*ParsedHash, error) {
	clean := strings.TrimSpace(encoded)
	if clean == "" {
		return nil, ErrInvalidLegacyHash
	}
	if strings.Contains(clean, "$") {
		return parseDollarFormat(clean)
	}
	if strings.Contains(clean, ":") {
		return parseColonFormat(clean)
	}
	return nil, ErrInvalidLegacyHash
}

func parseDollarFormat(value string) (*ParsedHash, error) {
	parts := strings.Split(value, "$")
	if len(parts) != 4 {
		return nil, ErrInvalidLegacyHash
	}
	algorithm := normalizeAlgorithm(parts[0])
	iterations, err := parseIterations(parts[1])
	if err != nil {
		return nil, err
	}
	hashBytes, err := decodeHash(parts[3])
	if err != nil {
		return nil, err
	}
	if _, err := hashForAlgorithm(algorithm); err != nil {
		return nil, err
	}
	return &ParsedHash{Algorithm: algorithm, Iterations: iterations, Salt: []byte(parts[2]), Hash: hashBytes}, nil
}

func parseColonFormat(value string) (*ParsedHash, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 5 || !strings.EqualFold(parts[0], "pbkdf2") {
		return nil, ErrInvalidLegacyHash
	}
	algorithm := normalizeAlgorithm("pbkdf2_" + parts[1])
	iterations, err := parseIterations(parts[2])
	if err != nil {
		return nil, err
	}
	hashBytes, err := decodeHash(parts[4])
	if err != nil {
		return nil, err
	}
	if _, err := hashForAlgorithm(algorithm); err != nil {
		return nil, err
	}
	return &ParsedHash{Algorithm: algorithm, Iterations: iterations, Salt: []byte(parts[3]), Hash: hashBytes}, nil
}

func parseIterations(value string) (int, error) {
	iterations, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || iterations <= 0 {
		return 0, ErrInvalidIterationCount
	}
	return iterations, nil
}

func decodeHash(value string) ([]byte, error) {
	clean := strings.TrimSpace(value)
	if clean == "" {
		return nil, ErrInvalidLegacyHash
	}
	if isHexDigest(clean) {
		decoded, err := hex.DecodeString(clean)
		if err != nil {
			return nil, fmt.Errorf("%w: hash encoding", ErrInvalidLegacyHash)
		}
		return decoded, nil
	}
	if decoded, err := base64.StdEncoding.DecodeString(clean); err == nil {
		return decoded, nil
	}
	if decoded, err := hex.DecodeString(clean); err == nil {
		return decoded, nil
	}
	return nil, fmt.Errorf("%w: hash encoding", ErrInvalidLegacyHash)
}

func isHexDigest(value string) bool {
	if len(value) == 0 || len(value)%2 != 0 {
		return false
	}
	for _, char := range value {
		if (char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F') {
			continue
		}
		return false
	}
	return true
}

func normalizeAlgorithm(value string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(value), "-", "_"))
}

func hashForAlgorithm(algorithm string) (func() hash.Hash, error) {
	switch normalizeAlgorithm(algorithm) {
	case "pbkdf2_sha256", "sha256":
		return sha256.New, nil
	case "pbkdf2_sha1", "sha1":
		return sha1.New, nil
	default:
		return nil, ErrUnsupportedAlgorithm
	}
}

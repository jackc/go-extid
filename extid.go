package extid

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
)

type Type struct {
	prefix string
	b      cipher.Block
}

func NewType(prefix string, key []byte) (*Type, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &Type{
		prefix: prefix + "_",
		b:      b,
	}, nil
}

func (t *Type) Encode(n int64) string {
	var plaintext [16]byte
	var cyphertext [16]byte

	binary.BigEndian.PutUint64(plaintext[:], uint64(n))
	t.b.Encrypt(cyphertext[:], plaintext[:])
	return t.prefix + hex.EncodeToString(cyphertext[:])
}

func (t *Type) Decode(s string) (int64, error) {
	if !strings.HasPrefix(s, t.prefix) {
		return 0, errors.New("invalid prefix")
	}

	s = s[len(t.prefix):]
	if len(s) != 32 {
		return 0, errors.New("invalid length")
	}

	var plaintext [16]byte
	var cyphertext [16]byte

	_, err := hex.Decode(cyphertext[:], []byte(s))
	if err != nil {
		return 0, errors.New("invalid data")
	}

	t.b.Decrypt(plaintext[:], cyphertext[:])

	n := binary.BigEndian.Uint64(plaintext[:])
	return int64(n), nil
}

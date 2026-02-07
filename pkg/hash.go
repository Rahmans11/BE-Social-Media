package pkg

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type HashConfig struct {
	Memory  uint32
	Time    uint32
	Thread  uint8
	KeyLen  uint32
	SaltLen uint32
}

func NewHashConfig(m, t, kl, sl uint32, th uint8) *HashConfig {
	return &HashConfig{
		Memory:  m,
		Time:    t,
		Thread:  th,
		KeyLen:  kl,
		SaltLen: sl,
	}
}

func (h *HashConfig) UseRecommended() {
	h.Memory = 64 * 1024
	h.Time = 2
	h.Thread = 1
	h.KeyLen = 32
	h.SaltLen = 16
}

func (h *HashConfig) GenSalt() ([]byte, error) {
	salt := make([]byte, h.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func (h *HashConfig) GenHash(password string) (string, error) {
	salt, err := h.GenSalt()
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, h.Time, h.Memory, h.Thread, h.KeyLen) // $jenisKey$versiKey$konfigurasi(memory, time, parallelism)$salt$hash
	version := argon2.Version
	saltStr := base64.RawStdEncoding.EncodeToString(salt)
	hashStr := base64.RawStdEncoding.EncodeToString(hash)

	hashedPwd := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", version, h.Memory, h.Time, h.Thread, saltStr, hashStr)
	return hashedPwd, nil
}

func (h *HashConfig) ComparePwdAndHash(pwd string, hashedPwd string) (bool, error) {
	result := strings.Split(hashedPwd, "$")

	if len(result) != 6 {
		return false, errors.New("invalid hash format")
	}

	if result[1] != "argon2id" {
		return false, errors.New("invalid hash")
	}

	var Version int
	if _, err := fmt.Sscanf(result[2], "v=%d", &Version); err != nil {
		return false, errors.New("invalid version format")
	}
	if Version != argon2.Version {
		return false, errors.New("invalid version")
	}

	if _, err := fmt.Sscanf(result[3], "m=%d,t=%d,p=%d", &h.Memory, &h.Time, &h.Thread); err != nil {
		return false, errors.New("invalid config format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(result[4])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(result[5])
	if err != nil {
		return false, err
	}
	h.KeyLen = uint32(len(hash))

	hp := argon2.IDKey([]byte(pwd), salt, h.Time, h.Memory, h.Thread, h.KeyLen)

	if subtle.ConstantTimeCompare(hp, hash) == 0 {
		return false, nil
	}
	return true, nil
}

package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog/log"
)

var (
	vaultOnce sync.Once
	vault     *LunarVault
)

type LunarVault struct {
	key []byte
}

func InitializeVault() {
	vaultOnce.Do(func() {
		v, err := createVault()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize cosmic vault")
		}
		vault = v
	})
}

func createVault() (*LunarVault, error) {
	configDir, err := os.UserConfigDir() //replace os.Getenv("HOME") with os.UserConfigDir() for cross-platform 
	if err != nil {                      // return prper error when directory cant be dtermined
		return nil, err
	}
	keyPath := filepath.Join(configDir, "svcmgr", "vault.key")
	if err := os.MkdirAll(filepath.Dir(keyPath), 0700); err != nil {
		return nil, err
	}

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		newKey := make([]byte, 32)
		if _, err := rand.Read(newKey); err != nil {
			return nil, err
		}
		if err := os.WriteFile(keyPath, newKey, 0400); err != nil {
			return nil, err
		}
	}

	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	return &LunarVault{key: keyData}, nil
}

func EncryptConfig(data []byte) ([]byte, error) {
	InitializeVault()                                //vault not initialized
	block, err := aes.NewCipher(vault.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func DecryptConfig(data []byte) ([]byte, error) {
	InitializeVault()                                //guard as EncryptConfig vault 	
	block, err := aes.NewCipher(vault.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, errors.New("invalid cosmic ciphertext")
	}

	return gcm.Open(nil, data[:gcm.NonceSize()], data[gcm.NonceSize():], nil)
}

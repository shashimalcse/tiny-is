package security

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

type KeyPair struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

type KeyManager struct {
	keyPairs map[string]*KeyPair
	mutex    sync.RWMutex
}

func NewKeyManager() *KeyManager {
	return &KeyManager{
		keyPairs: make(map[string]*KeyPair),
	}
}

func (km *KeyManager) LoadKeys(keyDir string) error {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	keyDir = filepath.Join(cwd, keyDir)
	files, err := os.ReadDir(keyDir)
	if err != nil {
		return fmt.Errorf("failed to read key directory: %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) == ".pem" {
			keyName := strings.TrimSuffix(file.Name(), ".pem")
			err := km.loadKeyPair(keyName, filepath.Join(keyDir, file.Name()))
			if err != nil {
				log.Printf("Failed to load key pair %s: %v", keyName, err)
			}
		}
	}
	return nil
}

func (km *KeyManager) loadKeyPair(name, filePath string) error {
	pemData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read key file: %v", err)
	}
	privateKey, publicKey, err := parsePrivateKey(pemData)
	if err != nil {
		return fmt.Errorf("failed to parse key: %v", err)
	}
	km.mutex.Lock()
	defer km.mutex.Unlock()

	km.keyPairs[name] = &KeyPair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	return nil
}

func (km *KeyManager) GetKeyPair(name string) (*KeyPair, error) {
	km.mutex.RLock()
	defer km.mutex.RUnlock()
	keyPair, exists := km.keyPairs[name]
	if !exists {
		return nil, fmt.Errorf("key pair not found: %s", name)
	}
	return keyPair, nil
}

func parsePrivateKey(data []byte) (ed25519.PrivateKey, ed25519.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, nil, fmt.Errorf("failed to decode PEM block")
	}

	var privateKey interface{}
	var err error

	switch block.Type {
	case "OPENSSH PRIVATE KEY":
		privateKey, err = ssh.ParseRawPrivateKey(data)
	case "PRIVATE KEY":
		privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	default:
		return nil, nil, fmt.Errorf("unsupported key type: %s", block.Type)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	ed25519PrivateKey, ok := privateKey.(ed25519.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("key is not an Ed25519 private key")
	}

	publicKey := ed25519PrivateKey.Public().(ed25519.PublicKey)

	return ed25519PrivateKey, publicKey, nil
}

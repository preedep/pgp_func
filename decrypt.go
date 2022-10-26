package main

import (
	"bytes"
	"compress/gzip"
	_ "crypto/sha256"
	"errors"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	_ "golang.org/x/crypto/ripemd160"
	"io"
)

func Decrypt(entity *openpgp.Entity, encrypted []byte, passphraseBytes []byte) ([]byte, error) {

	err := entity.PrivateKey.Decrypt(passphraseBytes)
	if err != nil {
		return nil, err
	}
	for _, subkey := range entity.Subkeys {
		err := subkey.PrivateKey.Decrypt(passphraseBytes)
		if err != nil {
			return nil, err
		}
	}
	//////
	// Decode message
	block, err := armor.Decode(bytes.NewReader(encrypted))
	if err != nil {
		return []byte{}, fmt.Errorf("Error decoding: %v", err)
	}
	if block.Type != "Message" {
		return []byte{}, errors.New("Invalid message type")
	}

	// Decrypt message
	entityList := openpgp.EntityList{entity}
	messageReader, err := openpgp.ReadMessage(block.Body, entityList, nil, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("Error reading message: %v", err)
	}
	read, err := io.ReadAll(messageReader.UnverifiedBody)
	if err != nil {
		return []byte{}, fmt.Errorf("Error reading unverified body: %v", err)
	}

	// Uncompress message
	reader := bytes.NewReader(read)
	uncompressed, err := gzip.NewReader(reader)
	if err != nil {
		return []byte{}, fmt.Errorf("Error initializing gzip reader: %v", err)
	}
	defer func(uncompressed *gzip.Reader) {
		err := uncompressed.Close()
		if err != nil {
			panic(err)
		}
	}(uncompressed)

	out, err := io.ReadAll(uncompressed)
	if err != nil {
		return []byte{}, err
	}

	// Return output - an unencoded, unencrypted, and uncompressed message
	return out, nil
}

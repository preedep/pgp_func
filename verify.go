package main

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"io"
)

func Verify(publicKeyEntity *openpgp.Entity, message []byte, signature []byte) error {
	sig, err := decodeSignature(signature)
	if err != nil {
		return err
	}
	hash := sig.Hash.New()
	messageReader := bytes.NewReader(message)
	_, err = io.Copy(hash, messageReader)
	if err != nil {
		return err
	}

	err = publicKeyEntity.PrimaryKey.VerifySignature(hash, sig)
	if err != nil {
		return err
	}
	return nil
}

func decodeSignature(signature []byte) (*packet.Signature, error) {
	signatureReader := bytes.NewReader(signature)
	block, err := armor.Decode(signatureReader)
	if err != nil {
		return nil, fmt.Errorf("Error decoding OpenPGP Armor: %s", err)
	}

	if block.Type != openpgp.SignatureType {
		return nil, errors.New("Invalid signature file")
	}

	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		return nil, err
	}

	sig, ok := pkt.(*packet.Signature)
	if !ok {
		return nil, errors.New("Error parsing signature")
	}
	return sig, nil
}

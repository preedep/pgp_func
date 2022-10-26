package main

import (
	"bytes"
	"golang.org/x/crypto/openpgp"
)

func Sign(entity *openpgp.Entity, message []byte) ([]byte, error) {
	writer := new(bytes.Buffer)
	reader := bytes.NewReader(message)
	err := openpgp.ArmoredDetachSign(writer, entity, reader, nil)
	if err != nil {
		return []byte{}, err
	}
	return writer.Bytes(), nil
}

package main

import (
	"crypto/rand"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
)

func loadOrCreateIdentity() (crypto.PrivKey, error) {

	if _, err := os.Stat("identity.key"); err == nil {

		data, err := os.ReadFile("identity.key")

		if err != nil {
			return nil, err
		}

		return crypto.UnmarshalPrivateKey(data)
	}

	priv, _, err := crypto.GenerateKeyPairWithReader(
		crypto.Ed25519,
		-1,
		rand.Reader,
	)

	if err != nil {
		return nil, err
	}

	data, err := crypto.MarshalPrivateKey(priv)

	if err != nil {
		return nil, err
	}

	err = os.WriteFile("identity.key", data, 0600)

	if err != nil {
		return nil, err
	}

	return priv, nil
}

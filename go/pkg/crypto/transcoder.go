package crypto

import (
	"bytes"
	"io"

	"filippo.io/age"
)

type Transcoder func([]byte) ([]byte, error)

func NopTranscoder(in []byte) ([]byte, error) {
	return in, nil
}

func AgeEncoder(recipients ...age.Recipient) Transcoder {
	return func(plaintext []byte) ([]byte, error) {
		b := &bytes.Buffer{}

		w, err := age.Encrypt(b, recipients...)
		if err != nil {
			return nil, err
		}

		if _, err := w.Write(plaintext); err != nil {
			return nil, err
		}

		if err := w.Close(); err != nil {
			return nil, err
		}

		return b.Bytes(), nil
	}
}

func AgeDecoder(identities ...age.Identity) Transcoder {
	return func(ciphertext []byte) ([]byte, error) {
		r, err := age.Decrypt(bytes.NewReader(ciphertext), identities...)
		if err != nil {
			return nil, err
		}

		return io.ReadAll(r)
	}
}

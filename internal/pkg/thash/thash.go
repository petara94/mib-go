package thash

import (
	"encoding/base64"
	"fmt"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/pkg/errors"
	"strings"
)

const SepEnc = "<!>"
const SepDec = "<?>"

type CryptoFunc = func(string, string) (string, error)

func EncryptWithPass(message, pass string) (hash string, err error) {
	var pgp *crypto.PGPMessage

	if pgp, err = crypto.EncryptMessageWithPassword(crypto.NewPlainMessageFromString(message), []byte(pass)); err != nil {
		return "", errors.Wrap(err, "thash: unable to encrypt message with password")
	}
	hash = base64.StdEncoding.EncodeToString(pgp.Data)

	return
}

func DecryptWithPass(hash, pass string) (message string, err error) {
	var pgp = &crypto.PGPMessage{}
	var msg *crypto.PlainMessage

	pgp.Data, err = base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return "", err
	}

	if msg, err = crypto.DecryptMessageWithPassword(pgp, []byte(pass)); err != nil {
		return "", errors.Wrap(err, "thash: unable to decrypt message with password")
	}

	return msg.GetString(), nil
}

func EncryptWithSeparator(text string, pass string) (string, error) {
	parts := strings.Split(text, SepEnc)

	if len(parts)%2 == 0 || len(parts) == 1 {
		return "", fmt.Errorf(`separator "%s" not used or used with wrong, use like this: %sword/hash%s`, SepEnc, SepEnc, SepEnc)
	}

	res := ""
	for i := 0; i < len(parts)-1; i += 2 {
		cryptoRes, err := EncryptWithPass(parts[i+1], pass)
		if err != nil {
			return "", err
		}

		res += parts[i] + SepDec + cryptoRes + SepDec
	}
	res += parts[len(parts)-1]

	return res, nil
}

func DecryptWithSeparator(text string, pass string) (string, error) {
	parts := strings.Split(text, SepDec)

	if len(parts)%2 == 0 || len(parts) == 1 {
		return "", fmt.Errorf(`separator "%s" not used or used with wrong, use like this: %sword/hash%s`, SepDec, SepDec, SepDec)
	}

	res := ""
	for i := 0; i < len(parts)-1; i += 2 {
		cryptoRes, err := DecryptWithPass(parts[i+1], pass)
		if err != nil {
			return "", err
		}

		res += parts[i] + cryptoRes
	}
	res += parts[len(parts)-1]

	return res, nil
}

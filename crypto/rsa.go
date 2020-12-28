package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
)

func ReadPublicRsaKey() (*rsa.PublicKey, error) {
	pemString, readErr := ioutil.ReadFile(os.Getenv("RSA_PUB_PATH"))
	if readErr != nil {
		return nil, readErr
	}
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(pemString)
	if err != nil {
		return nil, err
	}
	// convert ssh public key to rsa public key
	parsedCryptoKey := pubKey.(ssh.CryptoPublicKey)
	pubCrypto := parsedCryptoKey.CryptoPublicKey()
	return pubCrypto.(*rsa.PublicKey), nil
}

func ReadPrivateRsaKey() (*rsa.PrivateKey, error) {
	pemString, readErr := ioutil.ReadFile(os.Getenv("RSA_PATH"))
	if readErr != nil {
		return nil, readErr
	}
	block, _ := pem.Decode(pemString)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

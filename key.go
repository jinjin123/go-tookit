package lib

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	mathrand "math/rand"
	"github.com/lonng/nanoserver/pkg/errutil"

)


//GenRSAKey gen a rsa key pair, the bit size is 512
func GenRSAKey() (privateKey, publicKey string, err error) {
	//public gen the private key
	privKey, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return "", "", err
	}

	derStream := x509.MarshalPKCS1PrivateKey(privKey)
	privateKey = base64.StdEncoding.EncodeToString(derStream)

	//gen the public key
	pubKey := &privKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", "", err
	}
	publicKey = base64.StdEncoding.EncodeToString(derPkix)
	return privateKey, publicKey, nil
}

// RSAEncrypt encrypt data by rsa
func RSAEncrypt(plain []byte, pubKey string) ([]byte, error) {
	buf, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}
	p, err := x509.ParsePKIXPublicKey(buf)
	if err != nil {
		return nil, err
	}
	if pub, ok := p.(*rsa.PublicKey); ok {
		return rsa.EncryptPKCS1v15(rand.Reader, pub, plain) //RSA算法加密
	}
	return nil, errutil.ErrIllegalParameter
}

// RsaDecrypt decrypt data by rsa
func RSADecrypt(cipher []byte, privKey string) ([]byte, error) {
	if cipher == nil {
		return nil, errutil.ErrIllegalParameter
	}
	buf, err := base64.StdEncoding.DecodeString(privKey)
	if err != nil {
		return nil, err
	}
	priv, err := x509.ParsePKCS1PrivateKey(buf)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipher) //RSA解密算法
}

func pemParse(data []byte, pemType string) ([]byte, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("No PEM block found")
	}
	if pemType != "" && block.Type != pemType {
		return nil, fmt.Errorf("Key's type is '%s', expected '%s'", block.Type, pemType)
	}
	return block.Bytes, nil
}

func ParsePrivateKey(data []byte) (*rsa.PrivateKey, error) {
	pemData, err := pemParse(data, "RSA PRIVATE KEY")
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PrivateKey(pemData)
}

func LoadPrivateKey(privKeyPath string) (*rsa.PrivateKey, error) {
	certPEMBlock, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return nil, err
	}

	return ParsePrivateKey(certPEMBlock)
}

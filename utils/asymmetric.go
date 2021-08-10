package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto"
	//"crypto/sha256"
	//"fmt"
)


//Generate ECDSA keys
func GenerateECDSA() (*ecdsa.PrivateKey, crypto.PublicKey, error){
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	publickey := privateKey.Public()
	return privateKey, publickey, nil
}


//Sign with ECDSA private key
func SignECDSA(privatekey *ecdsa.PrivateKey, hash []byte) ([]byte, error){
	sig, err := ecdsa.SignASN1(rand.Reader, privatekey, hash)
	if err != nil {
		return nil, err
	}
	return sig, nil
}


//Verify ECDSA signature
func VerifyECDSA(publickey *ecdsa.PublicKey, msghash, sig []byte) (bool){
	isvalid := ecdsa.VerifyASN1(publickey, msghash, sig)
	return isvalid
}



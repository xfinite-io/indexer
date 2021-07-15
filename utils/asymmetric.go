package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)


//Generate ECDSA keys
func GenerateECDSA(){
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	publickey := privatekey.Public()
	return privatekey, publickey, nil
}


//Sign with ECDSA private key
func SignECDSA(privatekey *PrivateKey, hash []byte){
	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash)
	if err != nil {
		return nil, err
	}
	return sig, nil
}


//Verify ECDSA signature
func VerifyECDSA(publickey *PublicKey, msghash, sig []byte){
	isvalid := ecdsa.VerifyASN1(&privateKey.PublicKey, msghash, sig)
	return isvalid
}



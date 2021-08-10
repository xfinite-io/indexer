package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

var HttpClient = &http.Client{}

//SecretResponse is a type of secret response
type SecretResponse struct {
	Data string `json:"data"`
}

var ErrKeyNotFound error =  errors.New("KEY_NOT_FOUND")
var ErrDecodeErr error =  errors.New("ERROR_DECODING")
var ErrSetSecret  = errors.New("CANT_SET_SECRET")

//SetSecret Request
type SetSecretRequest struct{
	Blockchain string `json:"blockchain"`
	Key string `json:"key"`
	Secret string `json:"secret"`
}


//GetSecret gets the secret of the user which is asked
func GetSecret( blockchain string, key string)( *SecretResponse, error ){
	//KMS HOST
	var kmsHost = os.Getenv("KMS_HOST")
	res, err := http.Get( fmt.Sprintf(kmsHost+"/secret?blockchain=%s&key=%s",blockchain,key) )
	if err != nil {
		return &SecretResponse{}, ErrKeyNotFound
	}
	defer res.Body.Close()
	var secretresp SecretResponse
	if err := json.NewDecoder(res.Body).Decode(&secretresp); err!=nil{
		return &SecretResponse{}, ErrDecodeErr
	}
	return &secretresp,nil
}

////SetSecret sets the secret of the user which is asked
func SetSecret( blockchain string, key string, secret string )(bool, error){
	//KMS HOST
	var kmsHost = os.Getenv("KMS_HOST")
	setData := SetSecretRequest{
		Blockchain: blockchain,
		Key: key,
		Secret: secret,
	}
	byte,_ := json.Marshal(setData)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(kmsHost+"/secret"), bytes.NewBuffer(byte))
	if err != nil{
		return false, ErrSetSecret
	}
	resp, err := HttpClient.Do(req)
	if err != nil{
		log.Fatalf("Error %s\n", err )
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}
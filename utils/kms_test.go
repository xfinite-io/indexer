package utils

import (
	"os"
	"testing"
)

func TestGetSecret(t *testing.T) {
	 os.Setenv("KMS_HOST","http://104.211.218.220:3000")
	 notExistingKey := "FANTOKEN_111_publicKey"

	 data, err := GetSecret( "algo", notExistingKey )
	 if err == nil {
	 	t.Errorf("Not existing key in KMS should throw an error , err=%s\n", err)
	 }
	 if data.Data != "" {
		t.Errorf("Error")
	 }

	 //should be able to fetch the details of the key from the kms which exists
	existingKey := "FANTOKEN_XYZ"

	data, err = GetSecret( "algo", existingKey )
	if err != nil {
		t.Errorf("Existing key in KMS should not throw an error , err=%s\n", err)
	}
	if data.Data == "" {
		t.Errorf("Error")
	}
}


func TestSetSecret(t *testing.T) {
	os.Setenv("KMS_HOST","http://104.211.218.220:3000")
	notExistingCurrently := "FANTOKEN_XYZ"
	secretToSet := "Xfinite"

	data, err := SetSecret("algo", notExistingCurrently, secretToSet)
	if err != nil {
		t.Errorf("Should be able to set secret in kms")
	}
	if data == false {
		t.Errorf("Response should be true")
	}
}
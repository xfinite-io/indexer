package utils

import (
	"testing"
	"github.com/twinj/uuid"
)

func TestCreateToken(t *testing.T) {
	invalid_userid := uint64(1)
	td, err := CreateToken(invalid_userid)
	if err == nil {
		t.Errorf("CreateToken with invalid userid should throw an error , err=%s\n", err)
	} 
	if td.AccessToken != "" || td.RefreshToken != "" || td.AccessUuid != "" || td.RefreshUuid != "" {
		t.Errorf("Error")
	}

	valid_userid := int64(1)
	td, err := CreateToken(valid_userid)
	if err != nil {
		t.Errorf("CreateToken with valid userid should not throw an error , err=%s\n", err)
	}
	if td.AccessToken == "" || td.RefreshToken == "" || td.AccessUuid == "" || td.RefreshUuid == "" {
		t.Errorf("Invalid token created, %v\n", td)
	}
}


func TestCreateAuth(t *Testing.T) {
	valid_userid := int64(1)
	var td TokenDetails
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = td.AccessUuid + "++" + strconv.Itoa(int(userid))

	td.AccessToken = "accesstoken"
	td.RefreshToken = "refreshtoken"

	err := CreateAuth(valid_userid, td)
	if err != nil {
		t.Errorf("CreateAuth with valid userid should not throw an error , err=%s\n", err)
	}
}
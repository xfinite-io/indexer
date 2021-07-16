package utils

import (
	"testing"
	"github.com/twinj/uuid"
	"net/http"
	"net/http/httptest"
	"io"
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


func TestCreateAuth(t *testing.T) {
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


func TestExtractToken(t *testing.T) {
	user_id = int64(1)
	real_token, err := CreateToken(user_id)
	if err != nil {
		t.Errorf("error")
	}
	err := CreateAuth(user_id, real_token)
	if err != nil {
		t.Errorf("error")
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		out := ExtractToken(r)
		io.WriteString(w, out)
	}

	req := httptest.NewRequest("GET", "http://localhost:4000/testextracttoken", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	client, err := http.Client{}
	request, err := http.NewRequest("GET", "http://localhost:4000/testextracttoken", nil)
	request.Header.Set("Authorization", "Bearer " + real_token.AccessToken)
	if err != nil {
		panic(err)
	}
	_, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	if body == "" {
		t.Errorf("ExtractToken should be able to get token from header")
	}
}


func TestVerifyToken(t *testing.T) {
	user_id = int64(1)
	real_token, err := CreateToken(user_id)
	if err != nil {
		t.Errorf("error")
	}
	err := CreateAuth(user_id, real_token)
	if err != nil {
		t.Errorf("error")
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		_, err := VerifyToken(r)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("VerifyToken should not throw an error, err=%s\n", err)))
		}
		w.Write([]byte(""))
	}

	req := httptest.NewRequest("GET", "http://localhost:4000/testverifytoken/1", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	client, err := http.Client{}
	request, err := http.NewRequest("GET", "http://localhost:4000/testverifytoken/1", nil)
	request.Header.Set("Authorization", "Bearer " + real_token.AccessToken)
	if err != nil {
		panic(err)
	}
	_, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	if body != "" {
		t.Errorf(body)
	}

	fake_token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	handler2 := func(w http.ResponseWriter, r *http.Request) {
		_, err := VerifyToken(r)
		if err == nil {
			w.Write([]byte("VerifyToken should throw an error when accesstoken is fake"))
		}
		w.Write([]byte(""))
	}

	req2 := httptest.NewRequest("GET", "http://localhost:4000/testverifytoken/2", nil)
	w2 := httptest.NewRecorder()
	handler2(w, req)
	client2, err := http.Client{}
	request2, err := http.NewRequest("GET", "http://localhost:4000/testverifytoken/2", nil)
	request2.Header.Set("Authorization", "Bearer " + fake_token)
	if err != nil {
		panic(err)
	}
	_, err := client2.Do(request)
	if err != nil {
		panic(err)
	}

	resp2 := w2.Result()
	body2, _ := io.ReadAll(resp2.Body)
	if body2 != "" {
		t.Errorf(body2)
	}
}

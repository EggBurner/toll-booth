package handler

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterUser(t *testing.T) {

	req := httptest.NewRequest(http.MethodPost, "/register", nil)

	creds := base64.StdEncoding.EncodeToString([]byte("myEmail:myPassword"))
	req.Header.Set("Authorization", "Basic "+creds)

	w := httptest.NewRecorder()
	HandleRegisterUser(w, req)

	desiredCode := 200

	if desiredCode != w.Code {
		t.Errorf("\nCouldnt register. desired code: %d, got %d", desiredCode, w.Code)
	}
	expectedMessage := []byte("request receivedmyEmailmyPassword")

	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("\nerror while rgistering. \nexpected output: %s\ngot: %s", expectedMessage, w.Body)
	}
}

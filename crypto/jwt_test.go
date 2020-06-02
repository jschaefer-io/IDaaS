package crypto

import (
	"net/http"
	"strings"
	"testing"
)

// Check that a huge number of jwt generations
// yield a unique jwt every time given
// the same secret and claims
func TestNewJWT(t *testing.T) {
	tokens := map[string]Token{}
	data := map[string]interface{}{}
	for i := 0; i < 10000; i++ {
		secret := NewToken()
		_, token, err := NewJWT(secret, data)
		if err != nil {
			t.Errorf("Error occured for token %s: %s", secret, err.Error())
		}
		if o, ok := tokens[token]; ok {
			t.Errorf("The following secrets created the same jwt: %s and %s", o, secret)
		}
		tokens[token] = secret
	}
}

// Test that a huge number of jwt can be
// validated correctly
func TestCheckJWT(t *testing.T) {
	baseId := 1234
	secret := NewToken()
	data := map[string]interface{}{
		"id": baseId,
	}
	for i := 0; i < 10000; i++ {
		_, token, err := NewJWT(secret, data)
		if err != nil {
			t.Errorf("Token creation should not return an error: %s", err)
		}
		_, testedToken, err := CheckJWT(token, "id", func(id int) (token Token, err error) {
			if id != baseId {
				t.Errorf("Expected id %d but got %d", baseId, id)
			}
			return secret, nil
		})

		if err != nil {
			t.Errorf("Token checking should not return an error: %s", err)
		}
		if !testedToken.Valid {
			t.Error("Token is should be valid")
		}
	}
}

// Tests, that the jwt can be extracted from
// an http.Request instance correctly
func TestExtractJWT(t *testing.T) {
	r, err := http.NewRequest("GET", "localhost", strings.NewReader(""))
	if err != nil {
		t.Error(err)
	}
	testData := []map[string]interface{}{
		{
			"value": "",
			"error": true,
			"token": "",
		},
		{
			"value": "Bearer",
			"error": true,
			"token": "",
		},
		{
			"value": "Bearer x",
			"error": false,
			"token": "x",
		},
		{
			"value": "Bearer abc6",
			"error": false,
			"token": "abc6",
		},
		{
			"value": "Bearer x f",
			"error": true,
			"token": "",
		},
		{
			"value": "x y",
			"error": true,
			"token": "",
		},
	}
	for _, d := range testData {
		value := d["value"].(string)
		isError := d["error"].(bool)
		token := d["token"].(string)
		r.Header.Set("Authorization", value)
		reqToken, err := ExtractJWT(r)
		if (err != nil) != isError {
			t.Errorf("Expected Error = \"%v\" but got \"%v\" for value \"%s\"", err != nil, isError, value)
			continue
		}
		if !isError && token != reqToken {
			t.Errorf("Expected token %s but got %s for value \"%s\"", token, reqToken, value)
		}
	}
}

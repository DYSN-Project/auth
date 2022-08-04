package jwt

/*
import (
	"github.com/DYSN-Project/auth/internal/pkg/jwt/claims"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testAuthTable = []struct {
	claims jwt.Claims
	key    string
}{
	{
		claims: claims.NewAuthClaims(123, 600),
		key:    "testTest123",
	},
	{
		claims: claims.NewAuthClaims(3452, 200),
		key:    "qwerty123",
	},
}

var testRegisterTable = []struct {
	claims jwt.Claims
	key    string
}{
	{
		claims: claims.NewRegisterClaims("test@mail.ru", "testTest123", 600),
		key:    "AFhjs324djsclk",
	},
	{
		claims: claims.NewRegisterClaims("auth-mks@mail.ru", "Qwertyqwerty123", 100),
		key:    "JEdsfJj23djei",
	},
}

func TestGenerateToken(t *testing.T) {
	jwtSrv := NewJwtService()
	for _, testCase := range testAuthTable {
		t.Logf("generating auth token ...")
		token, err := jwtSrv.GenerateToken(testCase.claims, testCase.key)
		if err != nil {
			t.Error(err)
		}
		assert.IsType(t, "string", token, "incorrect auth token")
	}

	for _, testCase := range testRegisterTable {
		t.Logf("generating register token ...")
		token, err := jwtSrv.GenerateToken(testCase.claims, testCase.key)
		if err != nil {
			t.Error(err)
		}
		assert.IsType(t, "string", token, "incorrect register token")
	}
}

func TestParseToken(t *testing.T) {
	jwtSrv := NewJwtService()
	for _, testCase := range testRegisterTable {
		t.Logf("generating register token ...")
		token, err := jwtSrv.GenerateToken(testCase.claims, testCase.key)
		if err != nil {
			t.Error(err)
		}
		t.Logf("parsing register token ...")
		parseTkn, err := jwtSrv.ParseRegisterToken(token, testCase.key)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, parseTkn, testCase.claims)
	}

	for _, testCase := range testAuthTable {
		t.Logf("generating auth token ...")
		token, err := jwtSrv.GenerateToken(testCase.claims, testCase.key)
		if err != nil {
			t.Error(err)
		}
		t.Logf("parsing auth token ...")
		parseTkn, err := jwtSrv.ParseAuthToken(token, testCase.key)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, parseTkn, testCase.claims)
	}
}

func TestVerify(t *testing.T) {
	jwtSrv := NewJwtService()
	for _, testCase := range testRegisterTable {
		t.Logf("generating register token ...")
		token, err := jwtSrv.GenerateToken(testCase.claims, testCase.key)
		if err != nil {
			t.Error(err)
		}
		t.Logf("verify register token ...")
		isValid, err := jwtSrv.Verify(token, testCase.key)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, true, isValid)
	}

	for _, testCase := range testAuthTable {
		t.Logf("generating register token ...")
		token, err := jwtSrv.GenerateToken(testCase.claims, testCase.key)
		if err != nil {
			t.Error(err)
		}
		t.Logf("verify register token ...")
		isValid, err := jwtSrv.Verify(token, testCase.key)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, true, isValid)
	}
}
*/

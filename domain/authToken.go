package domain

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/manishdangi98/banking-lib/errs"
	"github.com/manishdangi98/banking-lib/logger"
)

type AuthToken struct {
	token *jwt.Token
}

func (t AuthToken) NewAccessToken() (string, *errs.AppError) {
	signedString, err := t.token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		logger.Error("Failed while siging access token:" + err.Error())
		return "", errs.NewUnexpectedError("cannot generate access token")
	}
	return signedString, nil
}

func NewAuthToken(claims Claims) AuthToken {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return AuthToken{token: token}
}

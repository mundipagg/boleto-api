package stonebank

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var (
	signKey *rsa.PrivateKey
)

const (
	privKeyPath = "../OpenBank.pem"
)

func generateJWT() (string, error) {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return "", err
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return "", err
	}

	now := time.Now()

	atClaims := jwt.MapClaims{}

	atClaims["exp"] = now.Add(15 * time.Minute).Unix()
	atClaims["nbf"] = now.Unix()
	atClaims["aud"] = "https://sandbox-accounts.openbank.stone.com.br/auth/realms/stone_bank"
	atClaims["realm"] = "stone_bank"
	atClaims["sub"] = "x"
	atClaims["clientId"] = "x"
	atClaims["iat"] = now.Unix()
	atClaims["jti"] = generateJTIFromTime(now)

	at := jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims)

	token, err := at.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func generateJTIFromTime(t time.Time) string {
	id, _ := uuid.NewUUID()
	nowStr := t.Format("2006-01-02T15:04:05.000Z")
	nowStr = strings.ReplaceAll(nowStr, "-", "")
	nowStr = strings.ReplaceAll(nowStr, "T", "")
	nowStr = strings.ReplaceAll(nowStr, ":", "")
	nowStr = strings.ReplaceAll(nowStr, ".", "")

	return fmt.Sprintf("%s.%s", nowStr[:17], id.String()[:7])
}

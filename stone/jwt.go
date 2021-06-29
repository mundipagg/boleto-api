package stone

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/mundipagg/boleto-api/config"
)

var (
	signKey *rsa.PrivateKey
)

const (
	privKeyPath = "../OpenBank.pem"
	StoneRealm  = "stone"
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
	atClaims["exp"] = now.Add(time.Duration(config.Get().StoneTokenDurationInMinutes) * time.Minute).Unix()
	atClaims["nbf"] = now.Unix()
	atClaims["aud"] = config.Get().StoneAudience
	atClaims["realm"] = StoneRealm
	atClaims["sub"] = config.Get().StoneClientID
	atClaims["clientId"] = config.Get().StoneClientID
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

	removable := []string{"-", "T", ":", "."}
	for _, ch := range removable {
		nowStr = strings.ReplaceAll(nowStr, ch, "")
	}

	return fmt.Sprintf("%s.%s", nowStr[:17], id.String()[:7])
}

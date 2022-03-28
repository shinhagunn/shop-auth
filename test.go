package main

import (
	"database/sql"
	"log"

	GJWT "github.com/golang-jwt/jwt"
	"github.com/zsmartex/pkg/jwt"
)

var hmacSampleSecret []byte

func main() {
	ks, _ := jwt.LoadOrGenerateKeys("./private.key", "./public.key")

	tokenString, err := jwt.ForgeToken("UID56415146", "huuhadzz@sczxdc.czx", "admin", sql.NullString{}, 3, ks.PrivateKey, GJWT.MapClaims{})

	log.Println(tokenString, err)

	j, err := jwt.ParseAndValidate(tokenString, ks.PublicKey)

	log.Println(j)
}

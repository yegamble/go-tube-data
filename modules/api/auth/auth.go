package auth

import (
	"crypto/rand"
	"encoding/hex"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	jwt "github.com/golang-jwt/jwt"
	"github.com/twinj/uuid"
	"os"
	"time"
)

type TokenDetails struct {
	AccessToken    string
	RefreshToken   string
	AccessUuid     string
	RefreshUuid    string
	AtExpires      int64
	RtExpires      int64
	CookieHTTPOnly bool
	CookieSameSite string
	KeyGenerator   func() string
}

func AuthRequired() fiber.Handler {
	return jwtware.New(jwtware.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
		SigningKey: []byte(os.Getenv("ACCESS_SECRET")),
	})
}

//func Session(h http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		session, err := store.Get(r, sessionName)
//		if err != nil {
//			log.WithError(err).Error("bad session")
//			http.SetCookie(w, &http.Cookie{Name: sessionName, MaxAge: -1, Path: "/"})
//			return
//		}
//
//		r = r.WithContext(context.WithValue(r.Context(), "session", session))
//		h(w, r)
//	}
//}

func CreateJWTToken(userid uint64) (*TokenDetails, error) {
	td := &TokenDetails{}

	var err error

	//Creating Access Token
	os.Setenv("ACCESS_SECRET", os.Getenv("ACCESS_SECRET"))
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = uuid.NewV4()
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = uuid.NewV4().String()
	rtClaims["user_id"] = userid
	rtClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func GenerateSessionToken(length int) string {

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

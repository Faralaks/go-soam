package tools

import (
	"context"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	p "go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TokenData struct {
	Uid       string
	Owner     string
	Status    string
	CreatedAt time.Time
	ExpireAt  time.Time
}

func (t *TokenData) Save() error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	_, err := TokensCol.InsertOne(ctx, bson.M{"_id": t.Uid, "owner": t.Owner, "status": t.Status, "createdAt": t.CreatedAt})
	if err != nil {
		return err
	}
	return nil
}

func CreateTokens(uid string, status string) (string, string, error) {
	atd := &TokenData{
		Owner:     uid,
		Status:    status,
		CreatedAt: time.Now().UTC(),
	}
	atd.ExpireAt = atd.CreatedAt.Add(Config.ATLifeTime)

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["owner"] = atd.Owner
	atClaims["status"] = atd.Status
	atClaims["exp"] = atd.ExpireAt.Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	signedAt, err := at.SignedString(Config.AccessSecret)
	if err != nil {
		return "", "", err
	}

	rtd := &TokenData{
		Uid:       p.NewObjectID().Hex(),
		Owner:     uid,
		Status:    status,
		CreatedAt: time.Now().UTC(),
	}
	rtd.ExpireAt = rtd.CreatedAt.Add(Config.RTLifeTime)

	rtClaims := jwt.MapClaims{}
	rtClaims["uid"] = rtd.Uid
	rtClaims["owner"] = rtd.Owner
	rtClaims["status"] = rtd.Status
	rtClaims["exp"] = rtd.ExpireAt.Unix()
	rt := jwt.NewWithClaims(jwt.SigningMethodHS512, rtClaims)
	signedRt, err := rt.SignedString(Config.RefreshSecret)
	if err != nil {
		return "", "", err
	}
	err = rtd.Save()
	if err != nil {
		return "", "", err
	}

	return signedAt, signedRt, nil
}

func ExtractToken(tokenString string, Key []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return Key, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

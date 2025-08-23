package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	mg "go.mongodb.org/mongo-driver/mongo"

	"i18n-flow/internal/model"
	"i18n-flow/internal/pkg/mongo"
)

func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coll := mongo.GetCollection("i18nflow", "users")
	err := coll.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mg.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func CreateUser(user model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coll := mongo.GetCollection("i18nflow", "users")
	_, err := coll.InsertOne(ctx, user)
	return err
}

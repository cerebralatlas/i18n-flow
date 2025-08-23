package mongo

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
)

// Init 初始化 MongoDB 连接
func Init(uri string) *mongo.Client {
	clientOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			log.Fatalf("MongoDB connection error: %v", err)
		}

		// 检查连接
		if err = client.Ping(ctx, nil); err != nil {
			log.Fatalf("MongoDB ping error: %v", err)
		}

		log.Println("✅ MongoDB connected")
		clientInstance = client
	})
	return clientInstance
}

// GetCollection 获取集合
func GetCollection(dbName, collName string) *mongo.Collection {
	return clientInstance.Database(dbName).Collection(collName)
}

// EnsureUserIndexes 确保 users 集合的索引存在
func EnsureUserIndexes() {
	coll := GetCollection("i18nflow", "users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"username": 1, // 按 username 升序
		},
		Options: options.Index().SetUnique(true), // 唯一索引
	}

	_, err := coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatalf("Failed to create users index: %v", err)
	}

	log.Println("✅ Users collection unique index ensured")
}

func EnsureIndexes() {
	EnsureUserIndexes()
}

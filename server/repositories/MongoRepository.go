package repositories

import "go.mongodb.org/mongo-driver/mongo"

type MongoOptions struct {
	Database *mongo.Database
	Recreate bool
}

type MongoRepository struct {
	db       *mongo.Database
	recreate bool
}

func NewMongoRepository(opts MongoOptions) *MongoRepository {
	return &MongoRepository{
		db:       opts.Database,
		recreate: opts.Recreate,
	}
}

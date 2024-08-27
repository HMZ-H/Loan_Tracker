// log_repository.go
package Repository

import (
	"Loan_manager/Domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogRepository interface {
	GetAllLogs() ([]Domain.Log, error)
	CreateLog(log Domain.Log) error
}

type logRepository struct {
	collection *mongo.Collection
}

func NewLogRepository(collection *mongo.Collection) LogRepository {
	return &logRepository{collection: collection}
}

func (lr *logRepository) GetAllLogs() ([]Domain.Log, error) {
	var logs []Domain.Log
	cursor, err := lr.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var log Domain.Log
		if err := cursor.Decode(&log); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (lr *logRepository) CreateLog(log Domain.Log) error {
	_, err := lr.collection.InsertOne(context.TODO(), log)
	return err
}

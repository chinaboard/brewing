package collection

import (
	"github.com/chinaboard/brewing/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection interface {
	Add(string, any) error
	Update(string, any) error
	Get(string) (any, error)
	List() (any, error)
	Del(string) error
}

type TaskCollection struct {
	mc *mongo.Collection
}

func NewTaskCollection(databaseName string) (*TaskCollection, error) {
	r, e := newDbClient()
	if e != nil {
		return nil, e
	}
	return &TaskCollection{mc: r.Database(databaseName).Collection("task")}, nil
}

func (tc *TaskCollection) Add(uniqueId string, task any) error {
	opt := options.Update().SetUpsert(true)

	filter := bson.M{
		"uniqueId": uniqueId,
	}

	update := bson.M{
		"$set": task,
	}

	_, err := tc.mc.
		UpdateOne(nil, filter, update, opt)

	return err
}

func (tc *TaskCollection) Update(uniqueId string, job any) error {
	return tc.Add(uniqueId, job)
}

func (tc *TaskCollection) Get(uniqueId string) (any, error) {
	filter := bson.M{
		"uniqueId": uniqueId,
	}
	var v model.Task
	err := tc.mc.FindOne(nil, filter).Decode(&v)
	return &v, err
}

func (tc *TaskCollection) List() (any, error) {
	cursor, err := tc.mc.Find(nil, bson.M{"exitCode": 0})
	if err != nil {
		return nil, err
	}
	var result []model.Task
	err = cursor.All(nil, &result)
	return &result, err
}

func (tc *TaskCollection) Del(uniqueId string) error {
	filter := bson.M{
		"uniqueId": uniqueId,
	}
	_, e := tc.mc.DeleteOne(nil, filter)
	return e
}

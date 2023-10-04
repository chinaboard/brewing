package collection

import (
	"github.com/chinaboard/brewing/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Collection interface {
	Add(any) error
	Update(any) error
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

func (tc *TaskCollection) Add(task any) error {
	opt := options.Update().SetUpsert(true)

	tv := task.(*model.Task)
	tv.UpdateAt = time.Now()

	filter := bson.M{
		"uniqueId": tv.UniqueId,
	}

	update := bson.M{
		"$set": tv,
	}

	_, err := tc.mc.UpdateOne(nil, filter, update, opt)

	return err
}

func (tc *TaskCollection) Update(job any) error {
	return tc.Add(job)
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

package collection

import (
	"github.com/chinaboard/brewing/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AsrCollection struct {
	mc *mongo.Collection
}

func NewAsrCollection(databaseName string) (*AsrCollection, error) {
	r, e := newDbClient()
	if e != nil {
		return nil, e
	}
	return &AsrCollection{
		mc: r.Database(databaseName).Collection("asr"),
	}, nil
}

func (ac *AsrCollection) Add(uniqueId string, asrModel any) error {
	opt := options.Update().SetUpsert(true)

	filter := bson.M{
		"uniqueId": uniqueId,
	}

	update := bson.M{
		"$set": asrModel,
	}

	_, err := ac.mc.
		UpdateOne(nil, filter, update, opt)

	return err
}

func (ac *AsrCollection) Update(uniqueId string, asrModel any) error {
	return ac.Add(uniqueId, asrModel)
}

func (ac *AsrCollection) Get(uniqueId string) (any, error) {
	filter := bson.M{
		"uniqueId": uniqueId,
	}
	var v model.AsrReponse
	err := ac.mc.FindOne(nil, filter).Decode(&v)
	return &v, err
}

func (ac *AsrCollection) List() (any, error) {
	cursor, err := ac.mc.Find(nil, bson.M{"content": nil})
	if err != nil {
		return nil, err
	}
	var result []model.AsrReponse
	err = cursor.All(nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}

func (ac *AsrCollection) Del(uniqueId string) error {
	filter := bson.M{
		"uniqueId": uniqueId,
	}
	_, e := ac.mc.DeleteOne(nil, filter)
	return e
}

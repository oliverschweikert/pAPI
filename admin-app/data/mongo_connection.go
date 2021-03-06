package data

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	PASSWORD_PREFIX = "mongodb+srv://db_admin:"
	PASSWORD_SUFFIX = "@apidata.ino8ejr.mongodb.net/?retryWrites=true&w=majority"
	DATABASE_NAME   = "papi_db"
)

type MongoConnection struct {
	passwordPrefix, passwordSuffix, password string
	client                                   *mongo.Client
	ctx                                      context.Context
	cancel                                   context.CancelFunc
}

func NewConnection() MongoConnection {
	return MongoConnection{PASSWORD_PREFIX, PASSWORD_SUFFIX, "", nil, nil, nil}
}

func (m *MongoConnection) FullURI() string {
	return m.passwordPrefix + m.password + m.passwordSuffix
}

func (m *MongoConnection) TestPassword(password string) error {
	m.password = password
	err := m.connect()
	if err != nil {
		return errors.New(m.FullURI())
	}
	defer m.close()
	err = m.client.Ping(m.ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoConnection) connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.FullURI()))
	if err != nil {
		cancel()
		return err
	}
	m.ctx = ctx
	m.cancel = cancel
	m.client = client
	return nil
}

func (m *MongoConnection) close() {
	defer m.cancel()
	defer func() {
		if err := m.client.Disconnect(m.ctx); err != nil {
			log.Fatal(err)
		}
	}()
}

func (m *MongoConnection) GetAllCollections() ([]*MongoCollection, error) {
	err := m.connect()
	if err != nil {
		return []*MongoCollection{}, err
	}
	defer m.close()

	collectionNames, err := m.client.Database(DATABASE_NAME).ListCollectionNames(m.ctx, bson.D{})
	if err != nil {
		return []*MongoCollection{}, err
	}
	var collections []*MongoCollection
	for _, n := range collectionNames {
		collections = append(collections, NewMongoCollection(n))
	}

	return collections, nil
}

func (m *MongoConnection) GetCategories() ([]*Category, error) {
	err := m.connect()
	if err != nil {
		return []*Category{}, err
	}
	defer m.close()

	cursor, err := m.client.Database(DATABASE_NAME).Collection("categories").Find(m.ctx, bson.D{})
	if err != nil {
		return []*Category{}, err
	}
	var categories []*Category
	for cursor.Next(m.ctx) {
		document := cursor.Current
		bsonData, err := bson.Marshal(document)
		if err != nil {
			return []*Category{}, err
		}
		var category Category
		err = bson.Unmarshal(bsonData, &category)
		if err != nil {
			return []*Category{}, err
		}
		categories = append(categories, &category)
	}
	return categories, nil
}

func (m *MongoConnection) UpdateCategory(categoryID primitive.ObjectID, category bson.D) error {
	err := m.connect()
	if err != nil {
		return err
	}
	defer m.close()

	updateFilter := bson.D{{Key: "_id", Value: categoryID}}
	result, err := m.client.Database(DATABASE_NAME).Collection("categories").ReplaceOne(m.ctx, updateFilter, category)
	if err != nil {
		return err
	}
	if result.ModifiedCount != 1 {
		return errors.New("error: a single entry wasn't modified")
	}
	return nil
}

func (m *MongoConnection) AddCategory(category bson.D) error {
	err := m.connect()
	if err != nil {
		return err
	}
	defer m.close()

	result, err := m.client.Database(DATABASE_NAME).Collection("categories").InsertOne(m.ctx, category)
	if err != nil {
		return err
	}
	if result.InsertedID == nil {
		return errors.New("error: Something went terribly wrong")
	}
	return nil
}

func (m *MongoConnection) DeleteCategory(id primitive.ObjectID) error {
	err := m.connect()
	if err != nil {
		return err
	}
	defer m.close()

	delFilter := bson.D{{Key: "_id", Value: id}}
	result, err := m.client.Database(DATABASE_NAME).Collection("categories").DeleteOne(m.ctx, delFilter)
	if err != nil {
		return err
	}

	if result.DeletedCount != 1 {
		return errors.New("error: Something has gone terribly wrong")
	}

	return nil
}

func (m *MongoConnection) GetAPIDatas() ([]*APIData, error) {
	err := m.connect()
	if err != nil {
		return []*APIData{}, err
	}
	defer m.close()

	cursor, err := m.client.Database(DATABASE_NAME).Collection("apiData").Find(m.ctx, bson.D{})
	if err != nil {
		return []*APIData{}, err
	}
	var apiDatas []*APIData
	for cursor.Next(m.ctx) {
		document := cursor.Current
		bsonData, err := bson.Marshal(document)
		if err != nil {
			return []*APIData{}, err
		}
		var apiData APIData
		err = bson.Unmarshal(bsonData, &apiData)
		if err != nil {
			return []*APIData{}, err
		}
		apiDatas = append(apiDatas, &apiData)
	}
	return apiDatas, nil
}

func (m *MongoConnection) AddCategoryToAPI(categoryHex string, apiID primitive.ObjectID) error {
	oid, err := primitive.ObjectIDFromHex(categoryHex)
	if err != nil {
		return err
	}
	err = m.connect()
	if err != nil {
		return err
	}
	defer m.close()

	searchFilter := bson.D{{Key: "_id", Value: oid}}
	exists := m.client.Database(DATABASE_NAME).Collection("categories").FindOne(m.ctx, searchFilter)
	if err != nil {
		return err
	}
	if exists.Err() != nil {
		return exists.Err()
	}
	updateFilter := bson.D{{Key: "_id", Value: apiID}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "categories", Value: categoryHex}}}}
	result, err := m.client.Database(DATABASE_NAME).Collection("apiData").UpdateOne(m.ctx, updateFilter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount != 1 {
		return errors.New("something went terribly wrong with AddCategoryToAPI")
	}

	return nil
}

func (m *MongoConnection) RemoveCategoryFromAPI(categoryHex string, apiID primitive.ObjectID) error {
	oid, err := primitive.ObjectIDFromHex(categoryHex)
	if err != nil {
		return err
	}
	err = m.connect()
	if err != nil {
		return err
	}
	defer m.close()

	searchFilter := bson.D{{Key: "_id", Value: oid}}
	exists := m.client.Database(DATABASE_NAME).Collection("categories").FindOne(m.ctx, searchFilter)
	if err != nil {
		return err
	}
	if exists.Err() != nil {
		return exists.Err()
	}
	updateFilter := bson.D{{Key: "_id", Value: apiID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "categories", Value: categoryHex}}}}
	result, err := m.client.Database(DATABASE_NAME).Collection("apiData").UpdateOne(m.ctx, updateFilter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount != 1 {
		return errors.New("something went terribly wrong with RemoveCategoryFromAPI")
	}

	return nil
}

func (m *MongoConnection) AddRequestToAPI(request string, apiID primitive.ObjectID) error {
	err := m.connect()
	if err != nil {
		return err
	}
	defer m.close()

	newRequest := Request{
		URL:        request,
		LastUpdate: time.Unix(0.0, 0.0),
		Response:   &map[string]interface{}{},
	}
	updateFilter := bson.D{{Key: "_id", Value: apiID}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "requests", Value: newRequest}}}}
	result, err := m.client.Database(DATABASE_NAME).Collection("apiData").UpdateOne(m.ctx, updateFilter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount != 1 {
		return errors.New("something went terribly wrong with AddCategoryToAPI")
	}

	return nil
}

func (m *MongoConnection) RemoveRequestFromAPI(requestURL string, apiID primitive.ObjectID) error {
	err := m.connect()
	if err != nil {
		return err
	}
	defer m.close()

	updateFilter := bson.D{{Key: "_id", Value: apiID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "requests", Value: bson.D{{Key: "request", Value: requestURL}}}}}}
	result, err := m.client.Database(DATABASE_NAME).Collection("apiData").UpdateOne(m.ctx, updateFilter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount != 1 {
		return errors.New("something went terribly wrong with RemoveCategoryFromAPI")
	}

	return nil
}

func (m *MongoConnection) UpdateAPIData(categoryID primitive.ObjectID, category bson.D) error {
	err := m.connect()
	if err != nil {
		return err
	}
	defer m.close()

	updateFilter := bson.D{{Key: "_id", Value: categoryID}}
	result, err := m.client.Database(DATABASE_NAME).Collection("apiData").ReplaceOne(m.ctx, updateFilter, category)
	if err != nil {
		return err
	}
	if result.ModifiedCount != 1 {
		return errors.New("error: replaced more than one entry")
	}
	return nil
}

func (m *MongoConnection) AddAPIData(apiData bson.D) error {
	err := m.connect()
	if err != nil {
		return err
	}
	defer m.close()

	result, err := m.client.Database(DATABASE_NAME).Collection("apiData").InsertOne(m.ctx, apiData)
	if err != nil {
		return err
	}
	if result.InsertedID == nil {
		return errors.New("error: Something went terribly wrong with AddAPIData")
	}
	return nil
}

func (m *MongoConnection) DeleteAPIData(id primitive.ObjectID) error {
	err := m.connect()
	if err != nil {
		return err
	}
	defer m.close()

	delFilter := bson.D{{Key: "_id", Value: id}}
	result, err := m.client.Database(DATABASE_NAME).Collection("apiData").DeleteOne(m.ctx, delFilter)
	if err != nil {
		return err
	}

	if result.DeletedCount != 1 {
		return errors.New("error: Something has gone wrong with DeleteAPIData")
	}

	return nil
}

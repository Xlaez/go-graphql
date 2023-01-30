package database

// TODO: change all [log.fatal] to 500 requests

import (
	"context"
	"log"
	"time"

	"github.com/Xlaez/go-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// TODO:  remeber to put into env file
var connection_str string = "mongodb://root:password123@localhost:6000"

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI(connection_str))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Fatal(err)
	}
	return &DB{
		client: client,
	}
}

func (db *DB) GetJob(id string) *model.JobListing {
	job_col := db.client.Database("graph").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	var job_listig model.JobListing
	if err := job_col.FindOne(ctx, filter).Decode(&job_listig); err != nil {
		log.Fatal(err)
	}
	return &job_listig
}

func (db *DB) GetJobs() []*model.JobListing {
	job_col := db.client.Database("graph").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var job_listig []*model.JobListing
	cursor, err := job_col.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	if err := cursor.All(context.TODO(), &job_listig); err != nil {
		panic(err)
	}

	return job_listig
}

func (db *DB) CreateJobListing(jobInfo model.CreateJobListingInput) *model.JobListing {
	job_col := db.client.Database("graph").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	data, err := job_col.InsertOne(ctx, jobInfo)
	if err != nil {
		log.Fatal(err)
	}
	inserted_id := data.InsertedID.(primitive.ObjectID).Hex()
	returnJobListing := model.JobListing{ID: inserted_id, Title: jobInfo.Title, Description: jobInfo.Description, Company: jobInfo.Company, URL: jobInfo.URL}

	return &returnJobListing
}

func (db *DB) UpdateJobListing(job_id string, jobInfo model.UpdateJobListingInput) *model.JobListing {
	job_col := db.client.Database("graph").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	updateJobInfo := bson.M{}

	if jobInfo.Title != nil {
		updateJobInfo["title"] = jobInfo.Title
	}
	if jobInfo.Description != nil {
		updateJobInfo["description"] = jobInfo.Description
	}
	if jobInfo.URL != nil {
		updateJobInfo["url"] = jobInfo.URL
	}

	_id, _ := primitive.ObjectIDFromHex(job_id)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateJobInfo}

	results := job_col.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var job_listig model.JobListing
	if err := results.Decode(&job_listig); err != nil {
		log.Fatal(err)
	}
	return &job_listig
}

func (db *DB) DeleteJobListing(job_id string) *model.DeleteJobResponse {
	job_col := db.client.Database("graph").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(job_id)
	filter := bson.M{"_id": _id}
	_, err := job_col.DeleteOne(ctx, filter)

	if err != nil {
		log.Fatal(err)
	}

	return &model.DeleteJobResponse{
		DeleteJobID: job_id,
	}
}

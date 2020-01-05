package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func Register(user User) (User, error) {
	collection := Client.Database("nc-user").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, error := collection.InsertOne(ctx, user)
	if error != nil {
		return User{}, error
	}
	newID, error := createId("userID")
	if error != nil {
		log.Println(error)
		return User{}, error
	}
	filter := bson.M{"_id": result.InsertedID}
	update := bson.M{"$set": bson.M{"id": newID}}
	// upsert := true
	// options := options.UpdateOptions{Upsert: &upsert}
	resultUpdate := collection.FindOneAndUpdate(ctx, filter, update)
	if resultUpdate.Err() != nil {
		return User{}, resultUpdate.Err()
	}
	var userResult User
	resultUpdate.Decode(&userResult)
	return userResult, nil
}

func createId(name string) (int, error) {
	collection := Client.Database("nc-user").Collection("counters")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var counter Counters
	filter := bson.M{"_id": name}
	update := bson.M{"$inc": bson.M{"id": 1}}
	upsert := true
	options := options.UpdateOptions{Upsert: &upsert}
	_, error := collection.UpdateOne(ctx, filter, update, &options)
	if error != nil {
		return 0, nil
	}
	result := collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return 0, result.Err()
	}
	error = result.Decode(&counter)
	if error != nil {
		log.Println(error)
		return 0, nil
	}
	log.Println(counter)

	return counter.ID, nil
}

// func UpdateStudentByID(id string, student Student) (Student, error) {
// 	collection := Client.Database("nc-user").Collection("users")
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	idInt, error := strconv.Atoi(id)
// 	if error != nil {
// 		return Student{}, error
// 	}
// 	var result Student
// 	filter := bson.M{"id": idInt}
// 	update := bson.M{"$set": student}
// 	error = collection.FindOneAndUpdate(ctx, filter, update).Decode(&result)
// 	if error != nil {
// 		return Student{}, error
// 	}
// 	return result, nil
// }

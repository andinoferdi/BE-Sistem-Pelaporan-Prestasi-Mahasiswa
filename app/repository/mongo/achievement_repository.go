package repository

import (
	"context"
	"time"

	model "sistem-pelaporan-prestasi-mahasiswa/app/model/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateAchievement(db *mongo.Database, achievement model.Achievement) (*model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection("achievements")

	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()

	result, err := collection.InsertOne(ctx, achievement)
	if err != nil {
		return nil, err
	}

	achievement.ID = result.InsertedID.(primitive.ObjectID)
	return &achievement, nil
}

func GetAchievementByID(db *mongo.Database, id string) (*model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	collection := db.Collection("achievements")
	var achievement model.Achievement

	err = collection.FindOne(ctx, bson.M{
		"_id":      objectID,
		"deletedAt": bson.M{"$exists": false},
	}).Decode(&achievement)
	if err != nil {
		return nil, err
	}

	return &achievement, nil
}

func UpdateAchievement(db *mongo.Database, id string, req model.UpdateAchievementRequest) (*model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	collection := db.Collection("achievements")

	update := bson.M{
		"updatedAt": time.Now(),
	}

	if req.AchievementType != "" {
		update["achievementType"] = req.AchievementType
	}
	if req.Title != "" {
		update["title"] = req.Title
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.Details != nil {
		update["details"] = req.Details
	}
	if req.Attachments != nil {
		update["attachments"] = req.Attachments
	}
	if req.Tags != nil {
		update["tags"] = req.Tags
	}
	if req.Points != nil {
		update["points"] = *req.Points
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, err
	}

	return GetAchievementByID(db, id)
}

func DeleteAchievement(db *mongo.Database, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := db.Collection("achievements")
	now := time.Now()
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"deletedAt": now,
			"updatedAt": now,
		}},
	)
	return err
}

func GetAchievementsByStudentID(db *mongo.Database, studentID string) ([]model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection("achievements")
	cursor, err := collection.Find(ctx, bson.M{
		"studentId": studentID,
		"deletedAt": bson.M{"$exists": false},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []model.Achievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

func GetAchievementsByIDs(db *mongo.Database, ids []string) ([]model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		objectIDs = append(objectIDs, objectID)
	}

	if len(objectIDs) == 0 {
		return []model.Achievement{}, nil
	}

	collection := db.Collection("achievements")
	cursor, err := collection.Find(ctx, bson.M{
		"_id":      bson.M{"$in": objectIDs},
		"deletedAt": bson.M{"$exists": false},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []model.Achievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

func AddAttachmentToAchievement(db *mongo.Database, id string, attachment model.Attachment) (*model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	collection := db.Collection("achievements")

	_, err = collection.UpdateOne(
		ctx,
		bson.M{
			"_id":      objectID,
			"deletedAt": bson.M{"$exists": false},
		},
		bson.M{
			"$push": bson.M{"attachments": attachment},
			"$set":  bson.M{"updatedAt": time.Now()},
		},
	)
	if err != nil {
		return nil, err
	}

	return GetAchievementByID(db, id)
}


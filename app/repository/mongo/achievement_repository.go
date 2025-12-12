package repository

import (
	"context"
	"time"

	model "sistem-pelaporan-prestasi-mahasiswa/app/model/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IAchievementRepository interface {
	CreateAchievement(ctx context.Context, achievement *model.Achievement) (*model.Achievement, error)
	GetAchievementByID(ctx context.Context, id string) (*model.Achievement, error)
	UpdateAchievement(ctx context.Context, id string, req model.UpdateAchievementRequest) (*model.Achievement, error)
	DeleteAchievement(ctx context.Context, id string) error
	GetAchievementsByStudentID(ctx context.Context, studentID string) ([]model.Achievement, error)
	GetAchievementsByIDs(ctx context.Context, ids []string) ([]model.Achievement, error)
	AddAttachmentToAchievement(ctx context.Context, id string, attachment model.Attachment) (*model.Achievement, error)
}

type AchievementRepository struct {
	collection *mongo.Collection
}

func NewAchievementRepository(db *mongo.Database) IAchievementRepository {
	return &AchievementRepository{
		collection: db.Collection("achievements"),
	}
}

func (r *AchievementRepository) CreateAchievement(ctx context.Context, achievement *model.Achievement) (*model.Achievement, error) {
	achievement.ID = primitive.NilObjectID
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, achievement)
	if err != nil {
		return nil, err
	}

	achievement.ID = result.InsertedID.(primitive.ObjectID)
	return achievement, nil
}

func (r *AchievementRepository) GetAchievementByID(ctx context.Context, id string) (*model.Achievement, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var achievement model.Achievement
	filter := bson.M{
		"_id":      objectID,
		"deletedAt": bson.M{"$exists": false},
	}
	err = r.collection.FindOne(ctx, filter).Decode(&achievement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &achievement, nil
}

func (r *AchievementRepository) UpdateAchievement(ctx context.Context, id string, req model.UpdateAchievementRequest) (*model.Achievement, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

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

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, err
	}

	return r.GetAchievementByID(ctx, id)
}

func (r *AchievementRepository) DeleteAchievement(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"deletedAt": now,
			"updatedAt": now,
		}},
	)
	return err
}

func (r *AchievementRepository) GetAchievementsByStudentID(ctx context.Context, studentID string) ([]model.Achievement, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
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

func (r *AchievementRepository) GetAchievementsByIDs(ctx context.Context, ids []string) ([]model.Achievement, error) {
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

	cursor, err := r.collection.Find(ctx, bson.M{
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

func (r *AchievementRepository) AddAttachmentToAchievement(ctx context.Context, id string, attachment model.Attachment) (*model.Achievement, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	_, err = r.collection.UpdateOne(
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

	return r.GetAchievementByID(ctx, id)
}

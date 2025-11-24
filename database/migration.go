package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RunMigrations(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	achievementsCollection := db.Collection("achievements")

	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "student_id", Value: 1}},
			Options: options.Index().SetName("idx_student_id"),
		},
		{
			Keys: bson.D{{Key: "achievement_type", Value: 1}},
			Options: options.Index().SetName("idx_achievement_type"),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_created_at"),
		},
		{
			Keys: bson.D{
				{Key: "title", Value: "text"},
				{Key: "description", Value: "text"},
			},
			Options: options.Index().SetName("idx_text_search"),
		},
	}

	_, err := achievementsCollection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	fmt.Println("MongoDB migrations completed successfully")
	return nil
}


package repository_test

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	modelmongo "sistem-pelaporan-prestasi-mahasiswa/app/model/mongo"
	repositorymongo "sistem-pelaporan-prestasi-mahasiswa/app/repository/mongo"
)

func setupTestMongoDB(t *testing.T) *mongo.Database {
	mongoURI := "mongodb://localhost:27017"
	if mongoURI == "" {
		t.Skip("MongoDB not available for testing")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		t.Skipf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		t.Skipf("Failed to ping MongoDB: %v", err)
	}

	return client.Database("test_achievements")
}

func TestAchievementRepository_CreateAchievement_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestMongoDB(t)
	defer db.Client().Disconnect(context.Background())

	repo := repositorymongo.NewAchievementRepository(db)
	ctx := context.Background()

	achievement := &modelmongo.Achievement{
		StudentID:       "550e8400-e29b-41d4-a716-446655440000",
		AchievementType: "academic",
		Title:           "Test Achievement",
		Description:     "Test Description",
		Points:          100,
	}

	result, err := repo.CreateAchievement(ctx, achievement)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected achievement, got nil")
	}

	if result.ID == primitive.NilObjectID {
		t.Error("Expected non-nil ID, got NilObjectID")
	}

	if result.Title != achievement.Title {
		t.Errorf("Expected Title %s, got %s", achievement.Title, result.Title)
	}

	cleanupTestData(t, db, result.ID.Hex())
}

func TestAchievementRepository_GetAchievementByID_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestMongoDB(t)
	defer db.Client().Disconnect(context.Background())

	repo := repositorymongo.NewAchievementRepository(db)
	ctx := context.Background()

	achievement := &modelmongo.Achievement{
		StudentID:       "550e8400-e29b-41d4-a716-446655440000",
		AchievementType: "academic",
		Title:           "Test Achievement",
		Description:     "Test Description",
		Points:          100,
	}

	created, err := repo.CreateAchievement(ctx, achievement)
	if err != nil {
		t.Fatalf("Failed to create test achievement: %v", err)
	}

	result, err := repo.GetAchievementByID(ctx, created.ID.Hex())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected achievement, got nil")
	}

	if result.ID != created.ID {
		t.Errorf("Expected ID %s, got %s", created.ID.Hex(), result.ID.Hex())
	}

	cleanupTestData(t, db, created.ID.Hex())
}

func TestAchievementRepository_GetAchievementByID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestMongoDB(t)
	defer db.Client().Disconnect(context.Background())

	repo := repositorymongo.NewAchievementRepository(db)
	ctx := context.Background()

	nonExistentID := primitive.NewObjectID().Hex()

	result, err := repo.GetAchievementByID(ctx, nonExistentID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestAchievementRepository_DeleteAchievement_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestMongoDB(t)
	defer db.Client().Disconnect(context.Background())

	repo := repositorymongo.NewAchievementRepository(db)
	ctx := context.Background()

	achievement := &modelmongo.Achievement{
		StudentID:       "550e8400-e29b-41d4-a716-446655440000",
		AchievementType: "academic",
		Title:           "Test Achievement",
		Description:     "Test Description",
		Points:          100,
	}

	created, err := repo.CreateAchievement(ctx, achievement)
	if err != nil {
		t.Fatalf("Failed to create test achievement: %v", err)
	}

	err = repo.DeleteAchievement(ctx, created.ID.Hex())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result, err := repo.GetAchievementByID(ctx, created.ID.Hex())
	if err != nil {
		t.Fatalf("Expected no error when checking deleted achievement, got %v", err)
	}

	if result != nil {
		t.Error("Expected nil result after deletion, got non-nil")
	}
}

func cleanupTestData(t *testing.T, db *mongo.Database, id string) {
	ctx := context.Background()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Logf("Failed to parse ID for cleanup: %v", err)
		return
	}

	_, err = db.Collection("achievements").DeleteOne(ctx, primitive.M{"_id": objectID})
	if err != nil {
		t.Logf("Failed to cleanup test data: %v", err)
	}
}

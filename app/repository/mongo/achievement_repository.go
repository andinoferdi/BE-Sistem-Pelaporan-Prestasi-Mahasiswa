package repository

// #1 proses: import library yang diperlukan untuk MongoDB, context, dan time
import (
	"context"
	"time"

	model "sistem-pelaporan-prestasi-mahasiswa/app/model/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// #2 proses: struct untuk hasil aggregasi top students berdasarkan points
type TopStudentResult struct {
	StudentID        string `bson:"_id" json:"student_id"`
	TotalPoints      int    `bson:"totalPoints" json:"total_points"`
	AchievementCount int    `bson:"count" json:"achievement_count"`
}

// #3 proses: definisikan interface untuk operasi database achievement di MongoDB
type IAchievementRepository interface {
	CreateAchievement(ctx context.Context, achievement *model.Achievement) (*model.Achievement, error)
	GetAchievementByID(ctx context.Context, id string) (*model.Achievement, error)
	UpdateAchievement(ctx context.Context, id string, req model.UpdateAchievementRequest) (*model.Achievement, error)
	DeleteAchievement(ctx context.Context, id string) error
	GetAchievementsByStudentID(ctx context.Context, studentID string) ([]model.Achievement, error)
	GetAchievementsByIDs(ctx context.Context, ids []string) ([]model.Achievement, error)
	AddAttachmentToAchievement(ctx context.Context, id string, attachment model.Attachment) (*model.Achievement, error)
	GetAchievementsByType(ctx context.Context) (map[string]int, error)
	GetCompetitionLevelDistribution(ctx context.Context) (map[string]int, error)
	GetTopStudentsByPoints(ctx context.Context, limit int) ([]TopStudentResult, error)
}

// #4 proses: struct repository untuk operasi database achievement di MongoDB
type AchievementRepository struct {
	collection *mongo.Collection
}

// #5 proses: constructor untuk membuat instance AchievementRepository baru
func NewAchievementRepository(db *mongo.Database) IAchievementRepository {
	return &AchievementRepository{
		collection: db.Collection("achievements"),
	}
}

// #6 proses: buat achievement baru di MongoDB
func (r *AchievementRepository) CreateAchievement(ctx context.Context, achievement *model.Achievement) (*model.Achievement, error) {
	// #6a proses: reset ID dan set timestamp sebelum insert
	achievement.ID = primitive.NilObjectID
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()

	// #6b proses: insert achievement ke collection dan ambil ID yang baru dibuat
	result, err := r.collection.InsertOne(ctx, achievement)
	if err != nil {
		return nil, err
	}

	achievement.ID = result.InsertedID.(primitive.ObjectID)
	return achievement, nil
}

// #7 proses: ambil achievement berdasarkan ID
func (r *AchievementRepository) GetAchievementByID(ctx context.Context, id string) (*model.Achievement, error) {
	// #7a proses: convert string ID jadi ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// #7b proses: buat filter untuk cari achievement yang belum dihapus
	var achievement model.Achievement
	filter := bson.M{
		"_id":       objectID,
		"deletedAt": bson.M{"$exists": false},
	}
	// #7c proses: cari di collection dan decode hasil ke struct
	err = r.collection.FindOne(ctx, filter).Decode(&achievement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &achievement, nil
}

// #8 proses: update achievement yang sudah ada, partial update
func (r *AchievementRepository) UpdateAchievement(ctx context.Context, id string, req model.UpdateAchievementRequest) (*model.Achievement, error) {
	// #8a proses: convert string ID jadi ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// #8b proses: buat update document, mulai dari updatedAt
	update := bson.M{
		"updatedAt": time.Now(),
	}

	// #8c proses: tambahkan field ke update document jika ada nilainya
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

	// #8d proses: update document di MongoDB dengan $set operator
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, err
	}

	// #8e proses: ambil achievement yang sudah diupdate untuk return
	return r.GetAchievementByID(ctx, id)
}

// #9 proses: soft delete achievement dengan set deletedAt
func (r *AchievementRepository) DeleteAchievement(ctx context.Context, id string) error {
	// #9a proses: convert string ID jadi ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// #9b proses: update achievement dengan set deletedAt untuk soft delete
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

// #10 proses: ambil semua achievement milik student tertentu
func (r *AchievementRepository) GetAchievementsByStudentID(ctx context.Context, studentID string) ([]model.Achievement, error) {
	// #10a proses: query untuk cari achievement dengan filter studentId dan belum dihapus
	cursor, err := r.collection.Find(ctx, bson.M{
		"studentId": studentID,
		"deletedAt": bson.M{"$exists": false},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// #10b proses: decode semua hasil ke slice achievements
	var achievements []model.Achievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

// #11 proses: ambil beberapa achievement sekaligus berdasarkan list ID
func (r *AchievementRepository) GetAchievementsByIDs(ctx context.Context, ids []string) ([]model.Achievement, error) {
	// #11a proses: convert semua string ID jadi ObjectID
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		objectIDs = append(objectIDs, objectID)
	}

	// #11b proses: jika tidak ada ID valid, return slice kosong
	if len(objectIDs) == 0 {
		return []model.Achievement{}, nil
	}

	// #11c proses: query dengan $in untuk ambil multiple documents
	cursor, err := r.collection.Find(ctx, bson.M{
		"_id":       bson.M{"$in": objectIDs},
		"deletedAt": bson.M{"$exists": false},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// #11d proses: decode semua hasil ke slice achievements
	var achievements []model.Achievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

// #12 proses: tambahkan attachment ke achievement yang sudah ada
func (r *AchievementRepository) AddAttachmentToAchievement(ctx context.Context, id string, attachment model.Attachment) (*model.Achievement, error) {
	// #12a proses: convert string ID jadi ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// #12b proses: update achievement dengan $push untuk tambahkan attachment ke array
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{
			"_id":       objectID,
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

	// #12c proses: ambil achievement yang sudah diupdate untuk return
	return r.GetAchievementByID(ctx, id)
}

// #13 proses: ambil statistik jumlah achievement per tipe menggunakan aggregation
func (r *AchievementRepository) GetAchievementsByType(ctx context.Context) (map[string]int, error) {
	// #13a proses: buat aggregation pipeline untuk group by achievementType
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"deletedAt": bson.M{"$exists": false},
			},
		},
		{
			"$group": bson.M{
				"_id":   "$achievementType",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	// #13b proses: eksekusi aggregation pipeline
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// #13c proses: loop hasil aggregation dan masukkan ke map result
	result := make(map[string]int)
	for cursor.Next(ctx) {
		var item struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&item); err != nil {
			continue
		}
		result[item.ID] = item.Count
	}

	return result, nil
}

// #14 proses: ambil distribusi level kompetisi dari achievement tipe competition
func (r *AchievementRepository) GetCompetitionLevelDistribution(ctx context.Context) (map[string]int, error) {
	// #14a proses: buat aggregation pipeline untuk filter competition dan group by level
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"achievementType":          "competition",
				"deletedAt":                bson.M{"$exists": false},
				"details.competitionLevel": bson.M{"$exists": true, "$ne": nil},
			},
		},
		{
			"$group": bson.M{
				"_id":   "$details.competitionLevel",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	// #14b proses: eksekusi aggregation pipeline
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// #14c proses: loop hasil aggregation dan masukkan ke map result, skip yang ID kosong
	result := make(map[string]int)
	for cursor.Next(ctx) {
		var item struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&item); err != nil {
			continue
		}
		if item.ID != "" {
			result[item.ID] = item.Count
		}
	}

	return result, nil
}

// #15 proses: ambil ranking top students berdasarkan total points menggunakan aggregation
func (r *AchievementRepository) GetTopStudentsByPoints(ctx context.Context, limit int) ([]TopStudentResult, error) {
	// #15a proses: validasi limit, set default 10 jika tidak valid
	if limit <= 0 {
		limit = 10
	}

	// #15b proses: buat aggregation pipeline untuk sum points per student, sort, dan limit
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"deletedAt": bson.M{"$exists": false},
			},
		},
		{
			"$group": bson.M{
				"_id":         "$studentId",
				"totalPoints": bson.M{"$sum": "$points"},
				"count":       bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"totalPoints": -1},
		},
		{
			"$limit": limit,
		},
	}

	// #15c proses: eksekusi aggregation pipeline
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// #15d proses: decode semua hasil ke slice results
	var results []TopStudentResult
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

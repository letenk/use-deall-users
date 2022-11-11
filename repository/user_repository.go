package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/letenk/use_deal_user/models/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	GetAll() ([]domain.User, error)
	GetOne(id string) (domain.User, error)
	GetOneByUsername(username string) (domain.User, error)
	Insert(user domain.User) (string, error)
	Update(user domain.User) (bool, error)
	Delete(id string) (bool, error)
}

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetAll() ([]domain.User, error) {
	// Create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// Cancel context after all process ends
	defer cancel()

	var users []domain.User

	// Set option find sort with `desc`
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := r.db.Collection("users").Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return users, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item domain.User

		err := cursor.Decode(&item)
		if err != nil {
			return users, err
		}

		users = append(users, item)
	}

	return users, nil
}

func (r *userRepository) GetOne(id string) (domain.User, error) {
	// Create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// Cancel context after all process ends
	defer cancel()

	var user domain.User

	// Find one by id
	err := r.db.Collection("users").FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) GetOneByUsername(username string) (domain.User, error) {
	// Create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// Cancel context after all process ends
	defer cancel()

	var user domain.User

	// Find one by id
	err := r.db.Collection("users").FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) Insert(user domain.User) (string, error) {
	// Create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// Cancel context after all process ends
	defer cancel()

	// Create new uuid
	id := uuid.NewString()

	// Create new object
	newUser := domain.User{
		ID:        id,
		Fullname:  user.Fullname,
		Username:  user.Username,
		Password:  user.Password,
		Role:      user.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert into db
	res, err := r.db.Collection("users").InsertOne(ctx, newUser)
	if err != nil {
		return "", err
	}

	// Convert interface to string
	return res.InsertedID.(string), nil
}

func (r *userRepository) Update(user domain.User) (bool, error) {
	// Create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// Cancel context after all process ends
	defer cancel()

	userID := bson.M{"_id": user.ID}
	updateData := bson.D{
		{"$set", bson.D{
			{"fullname", user.Fullname},
			{"password", user.Password},
			{"role", user.Role},
			{"updated_at", time.Now()},
		}}}

	// Update
	result, err := r.db.Collection("users").UpdateOne(ctx, userID, updateData)
	if err != nil {
		return false, err
	}

	// If result.ModifiedCount not same 1 (unsuccess updated)
	if result.ModifiedCount != 1 {
		return false, errors.New("Update failed")
	}

	return true, nil
}

func (r *userRepository) Delete(id string) (bool, error) {
	// Create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// Cancel context after all process ends
	defer cancel()

	// userID will use as a selector
	userID := bson.M{"_id": id}
	result, err := r.db.Collection("users").DeleteOne(ctx, userID)
	if err != nil {
		return false, err
	}

	// If result.DeletedCount not same 1 (unsuccess deleted)
	if result.DeletedCount != 1 {
		return false, errors.New("Delete failed")
	}

	return true, nil
}

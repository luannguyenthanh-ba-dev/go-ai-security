package repository

import (
	"context"
	"time"

	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/domain"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/shared"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// MongoDB implementation of user repository

type mongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(collection *mongo.Collection) UserRepository {
	return &mongoUserRepository{collection: collection}
}

// Mongo - CreateUser creates a new user in the database and returns the created user
func (r *mongoUserRepository) CreateNew(ctx context.Context, user *domain.UserEntity) (*domain.UserEntity, error) {
	if user == nil {
		zap.L().Error("user is nil", zap.Any("user", user))
		return nil, domain.ErrUserInvalidInput
	}

	// Set user id and timestamps
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now().UnixMilli()
	user.UpdatedAt = time.Now().UnixMilli()

	// Set default values for the user if not provided
	if user.Role == "" {
		user.Role = shared.RoleUser
	}
	if user.Gender == 0 {
		user.Gender = shared.GenderUnknown
	}

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		zap.L().Error("error inserting user", zap.Error(err))
		// Wrap infra error before returning to usecase
		return nil, domain.ErrUserInternalServerError
	}

	return user, nil
}

// Mongo - FindAUserByFilters finds a user by filters
func (r *mongoUserRepository) FindOneByFilters(ctx context.Context, filters UserFilters) (*domain.UserEntity, error) {
	filter := primitive.D{}

	// Add ID filter if provided
	if filters.ID != nil {
		filter = append(filter, primitive.E{Key: "_id", Value: *filters.ID})
	}
	// Add Username filter if provided
	if filters.Username != nil {
		filter = append(filter, primitive.E{Key: "username", Value: *filters.Username})
	}
	// Add Email filter if provided
	if filters.Email != nil {
		filter = append(filter, primitive.E{Key: "email", Value: *filters.Email})
	}
	// Add Phone filter if provided
	if filters.Phone != nil {
		filter = append(filter, primitive.E{Key: "phone",
			Value: primitive.E{Key: "$regex",
				Value: primitive.Regex{Pattern: *filters.Phone, Options: "i"},
			},
		}) // case insensitive regex
	}
	// if filters.FromTime != nil {
	// 	filter = append(filter, primitive.E{Key: "created_at", Value: primitive.E{Key: "$gte", Value: *filters.FromTime}})
	// }
	// if filters.ToTime != nil {
	// 	filter = append(filter, primitive.E{Key: "created_at", Value: primitive.E{Key: "$lte", Value: *filters.ToTime}})
	// }

	// Find one user by filters
	user := &domain.UserEntity{} // Use pointer because we want to return the user by reference
	err := r.collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		zap.L().Error("error finding user by filters", zap.Error(err))
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *mongoUserRepository) UpdateBasicInfoByID(ctx context.Context, userID *primitive.ObjectID, updateData *domain.UserEntity) (bool, error) {
	if updateData == nil {
		zap.L().Error("update data is nil", zap.Any("update data", updateData))
		return false, domain.ErrUserInvalidInput
	}

	if userID == nil {
		zap.L().Error("user id is nil to update", zap.Any("user id", userID))
		return false, domain.ErrUserInvalidID
	}

	// Build update document - only include fields that have values (not empty)
	update := primitive.D{}

	// Only update Name if provided (non-empty)
	if updateData.Name != "" {
		update = append(update, primitive.E{Key: "name", Value: updateData.Name})
	}

	// Only update Phone if provided (non-empty)
	if updateData.Phone != "" {
		update = append(update, primitive.E{Key: "phone", Value: updateData.Phone})
	}

	// Only update Address if provided (non-empty)
	if updateData.Address != "" {
		update = append(update, primitive.E{Key: "address", Value: updateData.Address})
	}

	// Only update Gender if provided (non-zero)
	if updateData.Gender.IsValid() {
		update = append(update, primitive.E{Key: "gender", Value: updateData.Gender})
	}

	// Only update Role if provided (non-empty)
	if updateData.Role.IsValid() {
		update = append(update, primitive.E{Key: "role", Value: updateData.Role})
	}

	// Only update Avatar if provided (non-empty)
	if updateData.Avatar != "" {
		update = append(update, primitive.E{Key: "avatar", Value: updateData.Avatar})
	}

	// Only update password if provided (non-empty)
	if updateData.Password != "" {
		update = append(update, primitive.E{Key: "password", Value: updateData.Password})
	}

	// If no fields to update, return early
	if len(update) == 0 {
		zap.L().Warn("no fields to update", zap.String("user_id", userID.Hex()))
		return false, nil
	}

	// Always update updated_at timestamp
	update = append(update, primitive.E{Key: "updated_at", Value: time.Now().UnixMilli()})

	// Execute update
	result, err := r.collection.UpdateOne(ctx, primitive.D{{Key: "_id", Value: *userID}}, primitive.D{{Key: "$set", Value: update}})
	if err != nil {
		zap.L().Error("error updating user", zap.Error(err), zap.String("user_id", userID.Hex()))
		return false, domain.ErrUserInternalServerError
	}

	// Check if any document was updated
	if result.MatchedCount == 0 {
		zap.L().Warn("user not found for update", zap.String("user_id", userID.Hex()))
		return false, domain.ErrUserNotFound
	}

	return result.ModifiedCount > 0, nil
}

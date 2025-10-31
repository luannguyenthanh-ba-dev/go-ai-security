package repository

import (
	"context"
	"time"

	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/shared"
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/domain"
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
func (r *mongoUserRepository) CreateUser(ctx context.Context, user *domain.UserEntity) (*domain.UserEntity, error) {
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
func (r *mongoUserRepository) FindAUserByFilters(ctx context.Context, filters UserFilters) (*domain.UserEntity, error) {
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
		zap.L().Error("error finding user by filters", zap.Error(err))
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

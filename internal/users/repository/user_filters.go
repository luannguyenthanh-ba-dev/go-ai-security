package repository

import (
	"github.com/luannguyenthanh-ba-dev/go-ai-security/internal/shared"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User filters

type UserFilters struct {
	ID        *primitive.ObjectID
	Username  *string
	Email     *string
	Phone     *string
	Address   *string
	Gender    *shared.Gender
	Role      *shared.Role
	FromTime  *int64
	ToTime    *int64
	SortBy    *string
	SortOrder *string
	Limit     *int
	Offset    *int
}

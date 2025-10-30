package shared

// Shared roles

type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// Write a function to check if the role is valid
// IsValid is a method of the Role type that returns a boolean to check if the role is valid
func (r Role) IsValid() bool {
	switch r {
	case RoleSuperAdmin, RoleAdmin, RoleUser:
		return true
	default:
		return false
	}
}

type Gender int

const (
	GenderMale Gender = 1
	GenderFemale Gender = 2
	GenderUnknown Gender = 3
)

// Write a function to check if the gender is valid
// IsValid is a method of the Gender type that returns a boolean to check if the gender is valid
func (g Gender) IsValid() bool {
	switch g {
	case GenderMale, GenderFemale, GenderUnknown:
		return true
	default:
		return false
	}
}

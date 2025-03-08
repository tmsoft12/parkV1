package modelsuser

import (
	"gorm.io/gorm"
)

// RoleType defines a custom type for user roles.
type RoleType string

// User represents a user in the system.
// @Description User structure for registration and authentication.
type User struct {
	*gorm.Model
	Id        int      `json:"id" example:"1" format:"int64"`
	Username  string   `json:"username" example:"johndoe"`
	Firstname string   `json:"firstname" example:"John"`
	Lastname  string   `json:"lastname" example:"Doe"`
	Password  string   `json:"password" example:"hashed_password" format:"string"`
	IsActive  bool     `json:"isActive" example:"true"`
	Role      RoleType `json:"role" example:"admin"`
	ParkNo    *string  `json:"park_no" example:"P123"`
}

type UserRes struct {
	Id        int      `json:"id" example:"1" format:"int64" description:"Unique identifier of the user"`
	Username  string   `json:"username" example:"johndoe" description:"Username of the user"`
	Firstname string   `json:"firstname" example:"John" description:"First name of the user"`
	Lastname  string   `json:"lastname" example:"Doe" description:"Last name of the user"`
	IsActive  bool     `json:"isActive" example:"true" description:"Indicates if the user is active"`
	Role      RoleType `json:"role" example:"operator" description:"Role assigned to the user" enums:"admin,operator,accountant"`
	ParkNo    *string  `json:"park_no" example:"P123" description:"Assigned parking number of the user"`
}

const (
	AdminRole      RoleType = "admin"
	OperatorRole   RoleType = "operator"
	AccountantRole RoleType = "accountant"
)

type MacUser struct {
	Id          int    `json:"id"`
	MacUsername string `json:"macusername"`
	MacPassword string `json:"macpassword"`
}

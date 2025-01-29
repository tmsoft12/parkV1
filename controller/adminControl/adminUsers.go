package admincontrol

import (
	"math"
	"park/database"
	modelsuser "park/models/modelsUser"
	"park/util"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a new user
// @Summary Create a new user
// @Description Creates a new user in the database
// @Tags Users
// @Accept json
// @Produce json
// @Param user body modelsuser.User true "User data"
// @Success 200 {object} modelsuser.User
// @Failure 400 {string} string "Can not parse"
// @Failure 500 {string} string "Can not create"
// @Router /api/v1/users [post]
func CreateUser(c *fiber.Ctx) error {
	var user modelsuser.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON("Can not parse")
	}
	user.IsActive = false
	var existingUser modelsuser.User
	if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"message": "Username already exists"})
	}
	if len(user.Password) < 8 {
		return c.Status(400).JSON(fiber.Map{"message": "Password must be at least 8 characters long"})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Error hashing password"})
	}
	user.Password = string(hashedPassword)
	if !util.IsValidRole(user.Role) {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid role provided"})
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON("can not create")
	}
	return c.Status(200).JSON(user)
}

// GetAllUsers retrieves all users with pagination
// @Summary Get all users
// @Description Retrieves a list of users with pagination support
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {string} string "Can not retrieve users"
// @Router /api/v1/users [get]
func GetAllUsers(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	var users []modelsuser.User
	var totalUsers int64

	if err := database.DB.Model(&modelsuser.User{}).Count(&totalUsers).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Can not retrieve users"})
	}

	if err := database.DB.Limit(limit).Offset(offset).Order("id DESC").Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Can not retrieve users"})
	}

	var userRes []modelsuser.UserRes
	for _, user := range users {
		userRes = append(userRes, modelsuser.UserRes{
			Id:        user.Id,
			Username:  user.Username,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			IsActive:  user.IsActive,
			Role:      user.Role,
			ParkNo:    user.ParkNo,
		})
	}

	totalPages := int(math.Ceil(float64(totalUsers) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1

	// Return response
	return c.Status(200).JSON(fiber.Map{
		"users":      userRes,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
		"hasNext":    hasNext,
		"hasPrev":    hasPrev,
	})
}

// UserGetByID retrieves a user by ID
// @Summary Get user by ID
// @Description Retrieves a user from the database using their unique ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} modelsuser.UserRes
// @Router /api/v1/users/{id} [get]
func UserGetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var user modelsuser.User

	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	userRes := modelsuser.UserRes{
		Id:        user.Id,
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		IsActive:  user.IsActive,
		Role:      user.Role,
		ParkNo:    user.ParkNo,
	}

	return c.Status(200).JSON(userRes)
}

// UserUpdate updates user data based on the provided fields
// @Summary Update user fields based on the provided data
// @Description Updates a user's data (isActive, username, firstname, lastname, etc.) in the database based on the input provided
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Param user body modelsuser.User true "User data to update"
// @Success 200 {object} modelsuser.UserRes
// @Failure 400 {string} string "Invalid user data"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Error updating user"
// @Router /api/v1/users/{id} [put]
func UserUpdate(c *fiber.Ctx) error {
	// Kullanıcı ID'sini alıyoruz
	id := c.Params("id")
	var user modelsuser.User

	// Veritabanında kullanıcıyı buluyoruz
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Kullanıcıdan gelen veriyi alıyoruz
	var updatedUser modelsuser.User
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid user data",
		})
	}

	// Güncellenmiş alanları kontrol ediyoruz ve yalnızca mevcut veriyi güncelliyoruz
	if updatedUser.Username != "" {
		user.Username = updatedUser.Username
	}
	if updatedUser.Firstname != "" {
		user.Firstname = updatedUser.Firstname
	}
	if updatedUser.Lastname != "" {
		user.Lastname = updatedUser.Lastname
	}
	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Error hashing password",
			})
		}
		user.Password = string(hashedPassword)
	}
	if updatedUser.IsActive != user.IsActive {
		user.IsActive = updatedUser.IsActive
	}
	if updatedUser.Role != "" {
		if !util.IsValidRole(updatedUser.Role) {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid role",
			})
		}
		user.Role = updatedUser.Role
	}
	if updatedUser.ParkNo != nil {
		user.ParkNo = updatedUser.ParkNo
	}

	// Güncellenmiş kullanıcıyı veritabanına kaydediyoruz
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error updating user",
		})
	}

	// Güncellenmiş kullanıcıyı döndürüyoruz
	userRes := modelsuser.UserRes{
		Id:        user.Id,
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		IsActive:  user.IsActive,
		Role:      user.Role,
		ParkNo:    user.ParkNo,
	}

	return c.Status(200).JSON(userRes)
}

// UserDelete deletes a user by ID
// @Summary Delete a user by ID
// @Description Deletes a user's information from the database using their unique ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {string} string "User deleted successfully"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Error deleting user"
// @Router /api/v1/users/{id} [delete]
func UserDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	var user modelsuser.User

	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error deleting user",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// GetOperator retrieves all users with the role "operator"
// @Summary Get all operators
// @Description Retrieves a list of users who have the role "operator" in descending order by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} map[string]interface{} "List of operators with pagination metadata"
// @Failure 404 {object} map[string]interface{} "No operators found"
// @Failure 500 {object} map[string]interface{} "Error retrieving users with operator role"
// @Router /api/v1/user/operators [get]
func GetOperator(c *fiber.Ctx) error {
	role := "operator"
	var users []modelsuser.User
	var operators []modelsuser.UserRes

	// Pagination params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	if err := database.DB.Where("role = ?", role).Order("id DESC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error retrieving users with operator role",
		})
	}

	// Total count of operators to calculate totalPages
	var totalCount int64
	database.DB.Model(&modelsuser.User{}).Where("role = ?", role).Count(&totalCount)

	// Calculating pagination details
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1

	// Map the results to the operators slice
	for _, user := range users {
		operators = append(operators, modelsuser.UserRes{
			Id:        user.Id,
			Username:  user.Username,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			IsActive:  user.IsActive,
			Role:      user.Role,
			ParkNo:    user.ParkNo,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"operators":  operators,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
		"hasNext":    hasNext,
		"hasPrev":    hasPrev,
	})
}

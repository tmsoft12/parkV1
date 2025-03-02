package usercontrol

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"park/controller/realtime"
	"park/database"
	modelsuser "park/models/modelsUser"
	"park/util"
)

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Parkno   string `json:"parkno" `
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	IsActive  bool      `json:"is_active"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// @Summary      Register User
// @Description  Creates a new user and stores their hashed password. Example: { "username": "newUser", "password": "password123", "firstname": "John", "lastname": "Doe", "role": "admin" }
// @Tags         Auth
// @Accept       json
// @Produce      json
//
//	@Param        user body modelsuser.User true "User Registration Data" {
//	  "username": "newUser",
//	  "password": "password123",
//	  "firstname": "John",
//	  "lastname": "Doe",
//	  "role": "admin"
//	}
//
// @Success      201 {object} map[string]string "message: User Created"
// @Failure      400 {object} map[string]string "message: Bad Request"
// @Failure      400 {object} map[string]string "message: Username already exists"
// @Failure      400 {object} map[string]string "message: Password must be at least 8 characters long"
// @Failure      500 {object} map[string]string "message: Internal Server Error"
// @Router       /api/v1/auth/register [post]
func Register(c *fiber.Ctx) error {
	var user modelsuser.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Bad Request", "error": err.Error()})
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

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Internal Server Error", "error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User Created"})
}

// @Summary      Login User
// @Description   { "username": "Dowran", "password": "12345678", "parkno": "P4" }
// @Tags         Auth
// @Accept       json
// @Produce      json
//
//	@Param        credentials body LoginInput true "User Login Data" {
//	  "username": "Dowran",
//	  "password": "12345678",
//	  "parkno": "P4-1"
//	}
//
// @Success      200 {object} map[string]string "message: Login successful"
// @Failure      400 {object} map[string]string "message: Invalid request body"
// @Failure      401 {object} map[string]string "message: Invalid username or password"
// @Failure      500 {object} map[string]string "message: Internal Server Error"
// @Router       /api/v1/auth/login [post]
func Login(c *fiber.Ctx) error {
	var loginInput struct {
		LoginInput
		ParkNo string `json:"parkno" validate:"required"`
	}
	if err := c.BodyParser(&loginInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	var user modelsuser.User
	if err := database.DB.Where("username = ?", loginInput.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	if !user.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Account is not active",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	token, err := util.CreateJWT(user.Id, user.Username, user.Role, loginInput.ParkNo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating JWT",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		HTTPOnly: true,
		SameSite: "Strict",
		MaxAge:   86400,
	})
	util.LoginMath(user.Username, string(user.Role), loginInput.ParkNo)
	role := user.Role
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"role":    role,
	})
}

// @Summary      Logout User
// @Description  Ends the session of a logged-in user by deleting the JWT token cookie.
// @Tags         Auth
// @Produce      json
// @Success      200 {object} map[string]string "message: Logout successful"
// @Failure      500 {object} map[string]string "message: Internal Server Error"
// @Router      /api/v1/auth/logout [post]
func Logout(c *fiber.Ctx) error {
	userIDVal := c.Locals("username")
	roleVal := c.Locals("role")
	parkno := c.Locals("parkno")
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		HTTPOnly: true,
		SameSite: "Strict",
		Expires:  time.Now().Add(-time.Hour),
		MaxAge:   -1,
	})
	role, ok := roleVal.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error - Invalid user role type",
		})
	}

	Username, ok := userIDVal.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error - Invalid user type",
		})
	}

	total_payment := 0

	if role == string(modelsuser.OperatorRole) {
		var err error
		total_payment, err = util.CalculateV2(Username, role)
		if err != nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Failed to calculate total payment",
				"error":   err.Error(),
			})
		}
	}
	if err := realtime.ResetParkingCount(parkno.(string)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Logout successful but failed to reset parking count",
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":       "Logout successful",
		"total_payment": total_payment,
	})
}

// @Summary      Get current user information
// @Description  Retrieves the current user's username, role, and user ID from the JWT token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{} "Returns user information"
// @Failure      400 {object} map[string]string "message: Bad Request - Missing data from middleware"
// @Failure      401 {object} map[string]string "message: Unauthorized - Invalid token"
// @Failure      500 {object} map[string]string "message: Internal Server Error - Missing data from middleware"
// @Router       /api/v1/auth/me [get]
func Me(c *fiber.Ctx) error {
	usernameVal := c.Locals("username")
	roleVal := c.Locals("role")
	userIDVal := c.Locals("user_id")
	parkno := c.Locals("parkno")

	if usernameVal == nil || roleVal == nil || userIDVal == nil || parkno == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error - Missing data from middleware",
		})
	}

	username, ok := usernameVal.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error - Invalid username type",
		})
	}
	park, ok := parkno.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error - Invalid parkno type",
		})
	}
	role, ok := roleVal.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error - Invalid role type",
		})
	}

	userID, ok := userIDVal.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error - Invalid user ID type",
		})
	}

	return c.JSON(fiber.Map{
		"username": username,
		"role":     role,
		"user_id":  userID,
		"parkno":   park,
	})
}

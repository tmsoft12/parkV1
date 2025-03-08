package camfix

import (
	"fmt"
	"park/database"
	"park/models/camera"
	modelsuser "park/models/modelsUser"
	"park/util"

	"github.com/gofiber/fiber/v2"
)

// Camfix godoc
// @Summary Create a New CamFix
// @Description Creates a new cam and saves it to the database.
// @Tags CamFix
// @Accept json
// @Produce json
// @Param cam body camera.CamFix true "Cam details to be created"
// @Success 201 {object} camera.CamFix "Successfully created"
// @Router /api/v1/addcam [post]
func AddCam(c *fiber.Ctx) error {
	var fix camera.CamFix
	if err := c.BodyParser(&fix); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Can not parse data",
		})
	}

	var existing camera.CamFix
	if err := database.DB.Where("channel_name = ?", fix.ChannelName).First(&existing).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{
			"message": fmt.Sprintf("A camera with ChannelName '%s' already exists", fix.ChannelName),
		})
	}

	if err := database.DB.Create(&fix).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Can not create",
		})
	}

	fmt.Println(fix)
	return c.Status(201).JSON(fix)
}

// Camfix godoc
// @Summary Create a New CamFix
// @Description Creates a new cam and saves it to the database.
// @Tags CamFix
// @Accept json
// @Produce json
// @Param cam body camera.CamFix true "Cam details to be created"
// @Success 201 {object} camera.CamFix "Successfully created"
// @Failure 400 {object} Response "Invalid data"
// @Failure 409 {object} Response "Camera already exists"
// @Failure 500 {object} Response "Internal server error"
// @Router /api/v1/addcam [post]
func GetAllCams(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var total int64
	if err := database.DB.Model(&camera.CamFix{}).Count(&total).Error; err != nil {
		return c.Status(500).JSON(ErrorResponse{Error: "Cannot fetch total count"})
	}

	var cams []camera.CamFix
	if err := database.DB.Order("id DESC").Limit(limit).Offset(offset).Find(&cams).Error; err != nil {
		return c.Status(500).JSON(ErrorResponse{Error: "Cannot fetch cameras"})
	}

	totalPages := (int(total) + limit - 1) / limit
	hasNext := page < totalPages
	hasPrev := page > 1

	response := fiber.Map{
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
		"hasNext":    hasNext,
		"hasPrev":    hasPrev,
		"data":       cams,
	}

	return c.JSON(response)
}

// UpdateChannelIdsByChannelName godoc
// @Summary Update ChannelIds by ChannelName
// @Description Updates ChannelId for all CamFix records matching the provided ChannelName(s).
// @Tags CamFix
// @Accept json
// @Produce json
// @Param updates body []map[string]string true "List of ChannelName and ChannelId pairs to update"
// @Router /api/v1/update-channel-ids [put]
func UpdateChannelIdsByChannelName(c *fiber.Ctx) error {
	var updates []map[string]string
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Can not parse data",
		})
	}

	for _, update := range updates {
		channelName, nameExists := update["ChannelName"]
		channelId, idExists := update["ChannelId"]

		if !nameExists || !idExists {
			return c.Status(400).JSON(fiber.Map{
				"message": "Each update must contain 'ChannelName' and 'ChannelId'",
			})
		}

		result := database.DB.Model(&camera.CamFix{}).
			Where("channel_name = ?", channelName).
			Update("channel_id", channelId)

		if result.Error != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Can not update ChannelId for ChannelName: " + channelName,
			})
		}

		if result.RowsAffected == 0 {
			return c.Status(404).JSON(fiber.Map{
				"message": "No records found for ChannelName: " + channelName,
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "ChannelIds updated successfully",
		"updated": len(updates),
	})
}

// UpdateMacUser godoc
// @Summary Update an existing MacUser
// @Description Update an existing MacUser's details such as MacUsername and MacPassword
// @Tags CamFix
// @Accept json
// @Produce json
// @Param macuser body modelsuser.MacUser true "MacUser object"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/updatemac [put]
func UpdateMacUser(c *fiber.Ctx) error {
	var macuser modelsuser.MacUser

	if err := c.BodyParser(&macuser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	id := 1

	var existingUser modelsuser.MacUser
	err := database.DB.Where("id = ?", id).First(&existingUser).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	existingUser.MacUsername = macuser.MacUsername
	existingUser.MacPassword = macuser.MacPassword

	err = database.DB.Save(&existingUser).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"user":    existingUser,
	})
}

// UpdateCamera updates only the 'Type' field of a specific camera
// @Summary Update Camera Type
// @Description Update only the Type of a camera by its ID
// @Tags CamFix
// @Accept json
// @Produce json
// @Param id path int true "Camera ID"
// @Param body body UpdateCameraTypeRequest true "New camera type"
// @Success 200 {object} camera.CamFix
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/type/{id} [patch]
func UpdateCamera(c *fiber.Ctx) error {
	id := c.Params("id")
	var datacam camera.CamFix

	// Kamerayı veritabanında bul
	if err := database.DB.Where("id = ?", id).First(&datacam).Error; err != nil {
		return c.Status(404).JSON(ErrorResponse{Error: "Camera not found"})
	}

	// Gelen veriyi parse et
	var updateReq UpdateCameraTypeRequest
	if err := c.BodyParser(&updateReq); err != nil {
		return c.Status(400).JSON(ErrorResponse{Error: "Invalid request body"})
	}

	// Kamera tipi doğrulaması
	newType := camera.CameraType(updateReq.Type) // String'i CameraType'e dönüştür
	if newType == "" || !util.IsValidCamera(newType) {
		return c.Status(400).JSON(ErrorResponse{Error: "Invalid camera type"})
	}

	// Sadece Type alanını güncelle
	if err := database.DB.Model(&datacam).Update("type", newType).Error; err != nil {
		return c.Status(500).JSON(ErrorResponse{Error: "Internal server error while updating camera"})
	}

	return c.Status(200).JSON(datacam)
}

type UpdateCameraTypeRequest struct {
	Type string `json:"type" validate:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
type Response struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string             `json:"message"`
	User    modelsuser.MacUser `json:"user"`
}
type DeleteResponse struct {
	Message string `json:"message"`
}

// DeleteCam godoc
// @Summary Delete a CamFix by ID
// @Description Deletes a camera from the database using its ID
// @Tags CamFix
// @Accept json
// @Produce json
// @Param id path int true "ID of the camera to delete"
// @Success 200 {object} Response "Successfully deleted"
// @Failure 404 {object} Response "Camera not found"
// @Failure 500 {object} Response "Internal server error"
// @Router /api/v1/deletecam/{id} [delete]
func DeleteCam(c *fiber.Ctx) error {
	id := c.Params("id")

	var cam camera.CamFix
	if err := database.DB.Where("id = ?", id).First(&cam).Error; err != nil {
		return c.Status(404).JSON(Response{
			Message: fmt.Sprintf("Camera with ID '%s' not found", id),
		})
	}

	if err := database.DB.Delete(&cam).Error; err != nil {
		return c.Status(500).JSON(Response{
			Message: "Internal server error while deleting camera",
		})
	}

	return c.Status(200).JSON(Response{
		Message: fmt.Sprintf("Camera with ID '%s' successfully deleted", id),
	})
}

package camfix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

	// fmt.Println(fix)
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

	if err := database.DB.Where("id = ?", id).First(&datacam).Error; err != nil {
		return c.Status(404).JSON(ErrorResponse{Error: "Camera not found"})
	}

	var updateReq UpdateCameraTypeRequest
	if err := c.BodyParser(&updateReq); err != nil {
		return c.Status(400).JSON(ErrorResponse{Error: "Invalid request body"})
	}

	newType := camera.CameraType(updateReq.Type) // String'i CameraType'e dönüştür
	if newType == "" || !util.IsValidCamera(newType) {
		return c.Status(400).JSON(ErrorResponse{Error: "Invalid camera type"})
	}

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

type CameraType string

const (
	Inside  CameraType = "inside"
	Outside CameraType = "outside"
)

type ConfigResponse struct {
	Channels []struct {
		Id   string `json:"Id"`
		Name string `json:"Name"`
	} `json:"Channels"`
}

// SyncCamFixWithConfig godoc
// @Summary Sync CamFix records with config data
// @Description Fetches data from config endpoint and synchronizes CamFix records: creates new ones, updates existing ones, and deletes obsolete ones.
// @Tags CamFix
// @Accept json
// @Produce json
// // @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/sync-camfix [get]
func SyncCamFixWithConfig(c *fiber.Ctx) error {
	var user modelsuser.MacUser
	if err := database.DB.First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to fetch MacUser data: " + err.Error(),
		})
	}
	macroscope := os.Getenv("MACROSCOP_URL")
	url := fmt.Sprintf("http://%s/configex?login=%s&password=%s&responsetype=json", macroscope, user.MacUsername, user.MacPassword)
	fmt.Println(url)
	resp, err := http.Get(url)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to fetch config data: " + err.Error(),
		})
	}
	defer resp.Body.Close()

	var config ConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to decode config data: " + err.Error(),
		})
	}

	configChannels := make(map[string]string)
	for _, channel := range config.Channels {
		configChannels[channel.Name] = channel.Id
	}

	var existingCams []camera.CamFix
	if err := database.DB.Find(&existingCams).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to fetch existing cameras: " + err.Error(),
		})
	}

	for _, cam := range existingCams {
		if _, exists := configChannels[cam.ChannelName]; !exists {
			if err := database.DB.Delete(&cam).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
					Error: fmt.Sprintf("Failed to delete camera with ChannelName '%s': %v", cam.ChannelName, err),
				})
			}
		}
	}

	for channelName, channelId := range configChannels {
		var existingCam camera.CamFix
		err := database.DB.Where("channel_name = ?", channelName).First(&existingCam).Error
		if err != nil {
			newCam := camera.CamFix{
				ChannelName: channelName,
				ChannelId:   channelId,
				Type:        camera.CameraType(Outside),
			}
			if err := database.DB.Create(&newCam).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
					Error: fmt.Sprintf("Failed to create camera with ChannelName '%s': %v", channelName, err),
				})
			}
		} else {
			if existingCam.ChannelId != channelId {
				existingCam.ChannelId = channelId
				if err := database.DB.Save(&existingCam).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
						Error: fmt.Sprintf("Failed to update camera with ChannelName '%s': %v", channelName, err),
					})
				}
			}
		}
	}

	var updatedCams []camera.CamFix
	if err := database.DB.Find(&updatedCams).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to fetch updated cameras: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "CamFix records synchronized successfully",
		"data":    updatedCams,
	})
}

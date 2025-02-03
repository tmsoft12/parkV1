package tarifcontrol

import (
	"encoding/json"
	resmodel "park/controller/getdata/resModel"
	"park/database"
	"park/models/tarif"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CreateTarif godoc
// @Summary Create a New Tarif
// @Description Creates a new tarif and saves it to the database.
// @Tags Tarif
// @Accept json
// @Produce json
// @Param tarif body tarif.Tarif true "Tarif details to be created"
// @Success 201 {object} tarif.Tarif "Successfully created"
// @Failure 400 {object} resmodel.ErrorResponse "Invalid request data"
// @Failure 500 {object} resmodel.ErrorResponse "Failed to save data to the database"
// @Router /api/v1/accountant/tarif [post]
var TimeFormat = "2006-01-02 15:04:05"

// Tarif struct with custom time unmarshaling
type Tarif struct {
	Id         int       `json:"id"`
	Plate      string    `json:"plate"`
	Name       string    `json:"name"`
	Start_time time.Time `json:"start_time"`
	End_time   time.Time `json:"end_time"`
	Price      int       `json:"price"`
}

func (t *Tarif) UnmarshalJSON(data []byte) error {
	type Alias Tarif
	aux := &struct {
		Start_time string `json:"start_time"`
		End_time   string `json:"end_time"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	startTime, err := time.Parse(TimeFormat, aux.Start_time)
	if err != nil {
		return err
	}
	t.Start_time = startTime

	endTime, err := time.Parse(TimeFormat, aux.End_time)
	if err != nil {
		return err
	}
	t.End_time = endTime

	return nil
}

// CreateTarif godoc
// @Summary Create a New Tarif
// @Description Creates a new tarif and saves it to the database.
// @Tags Tarif
// @Accept json
// @Produce json
// @Param tarif body tarif.Tarif true "Tarif details to be created"
// @Success 201 {object} tarif.Tarif "Successfully created"
// @Failure 400 {object} resmodel.ErrorResponse "Invalid request data"
// @Failure 500 {object} resmodel.ErrorResponse "Failed to save data to the database"
// @Router /api/v1/accountant/tarif [post]
func CreateTarif(c *fiber.Ctx) error {
	var tarif Tarif

	if err := c.BodyParser(&tarif); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resmodel.ErrorResponse{
			Error:   "Failed to parse request body",
			Details: err.Error(),
		})
	}
	if err := database.DB.Create(&tarif).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(resmodel.ErrorResponse{
			Error:   "Failed to save data to the database",
			Details: err.Error(),
		})
	}
	return c.Status(201).JSON(tarif)
}

// DeleteTarif godoc
// @Summary Delete Tarif
// @Description Deletes a tarif by its ID.
// @Tags Tarif
// @Param id path int true "ID of the tarif to delete"
// @Success 200 {string} string "Tarif successfully deleted"
// @Failure 400 {object} resmodel.ErrorResponse "Invalid ID format"
// @Failure 404 {object} resmodel.ErrorResponse "Tarif not found"
// @Failure 500 {object} resmodel.ErrorResponse "Database error"
// @Router /api/v1/accountant/tarif/{id} [delete]
func DeleteTarif(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resmodel.ErrorResponse{
			Error:   "Invalid ID format",
			Details: "ID must be a number",
		})
	}

	var tarif tarif.Tarif
	if err := database.DB.First(&tarif, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(resmodel.ErrorResponse{
			Error:   "Tarif not found",
			Details: err.Error(),
		})
	}

	if err := database.DB.Delete(&tarif).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(resmodel.ErrorResponse{
			Error:   "Failed to delete tarif",
			Details: err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Tarif successfully deleted"})
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"totalPages"`
	HasNext    bool        `json:"hasNext"`
	HasPrev    bool        `json:"hasPrev"`
}

// GetAllTarif godoc
// @Summary Get all Tarifs with pagination
// @Description Retrieves all tarifs from the database with pagination support.
// @Tags Tarif
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} PaginatedResponse "List of tarifs"
// @Failure 500 {object} resmodel.ErrorResponse "Database error"
// @Router /api/v1/accountant/tarif [get]
func GetAllTarif(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	var tarifs []tarif.Tarif
	var totalCount int64

	if err := database.DB.Model(&tarif.Tarif{}).Count(&totalCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(resmodel.ErrorResponse{
			Error:   "Failed to get total count of tarifs",
			Details: err.Error(),
		})
	}

	if err := database.DB.Offset((page - 1) * limit).Limit(limit).Find(&tarifs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(resmodel.ErrorResponse{
			Error:   "Failed to retrieve tarifs",
			Details: err.Error(),
		})
	}

	totalPages := int(totalCount / int64(limit))
	if totalCount%int64(limit) != 0 {
		totalPages++
	}
	hasNext := page < totalPages
	hasPrev := page > 1

	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Data:       tarifs,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	})
}

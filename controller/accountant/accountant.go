package accountant

import (
	"database/sql"
	"fmt"
	"math"
	"park/database"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Operator struct {
	UserID       int     `json:"user_id"`
	TotalPayment float64 `json:"total_payment"`
	Username     string  `json:"username"`
	ParkNo       string  `json:"park_no"`
	IsActive     bool    `json:"is_active"`
}
type GetOperatorResponse struct {
	Operators  []Operator `json:"operators"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalCount int64      `json:"total_count"`
	HasNext    bool       `json:"has_next"`
	HasPrev    bool       `json:"has_prev"`
}
type Results struct {
	UserID       int     `json:"user_id"`
	TotalPayment float64 `json:"total_payment"`
	Username     string  `json:"username"`
	ParkNo       string  `json:"park_no"`
	Total        int     `json:"total"`
	IsActive     bool    `json:"is_active"`
}
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

// GetOperators godoc
// @Summary Get list of operators
// @Description Get list of operators with pagination
// @Tags Accountant
// @Accept  json
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(5)
// @Success 200 {object} GetOperatorResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/accountant [get]
func GetOperators(c *fiber.Ctx) error {
	var totalCount int64
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "5")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid page number",
		})
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid limit number",
		})
	}

	offset := (page - 1) * limit

	var results []Results
	firstDayOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	lastDayOfMonth := firstDayOfMonth.Add(time.Hour * 24 * 30)
	fmt.Printf("now %v", firstDayOfMonth)
	err_qr := database.DB.Raw(`
        WITH  operators AS (SELECT
        users.id AS user_id,
        SUM(car_models.total_payment) AS total_payment,
        users.username,
        users.park_no,
				users.is_active
        FROM 
        car_models 
    INNER JOIN 
        users
    ON 
        CAST(car_models.user_id AS bigint)= users.id 
                WHERE 
        car_models.user_id ~ '^\d+$'
        AND users.role = 'operator'
				AND car_models.start_time >= ?
				AND car_models.end_time <= ?
GROUP BY
 users.id LIMIT ? OFFSET ?) 
 SELECT operators.*, COUNT(*) AS total from operators GROUP BY operators.user_id, operators.total_payment, operators.username, operators.park_no, operators.is_active
	`, firstDayOfMonth.Format("2006-01-02 15:04:05.999999-07:00"), lastDayOfMonth.Format("2006-01-02 15:04:05.999999-07:00"), limit, offset).Scan(&results).Error
	if err_qr != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_qr.Error(),
		})
	}

	if len(results) == 0 {
		results = []Results{}
	}

	var operator []Operator
	for _, res := range results {
		optr := Operator{
			UserID:       res.UserID,
			TotalPayment: res.TotalPayment,
			Username:     res.Username,
			ParkNo:       res.ParkNo,
			IsActive:     res.IsActive,
		}
		totalCount = int64(res.Total)
		operator = append(operator, optr)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1

	return c.Status(200).JSON(GetOperatorResponse{
		Operators:  operator,
		Page:       page,
		Limit:      limit,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
		TotalCount: totalCount,
	})
}

type ResultById struct {
	StartAt      time.Time `json:"start_at"`
	EndAt        time.Time `json:"end_at"`
	OperatorId   int       `json:"opertor_id"`
	TotalPayment int       `json:"total_payment"`
}

// GetUserByIdPayments godoc
// @Summary Get payments for a specific operator by ID
// @Description Retrieves payments for a specific operator within their login/logout periods
// @Tags Accountant
// @Accept  json
// @Produce  json
// @Param id path int true "Operator ID"
// @Success 200 {array} ResultById
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/accountant/{id} [get]
func GetUserByIdPayments(c *fiber.Ctx) error {
	operatorId := c.Params("id")
	var operatorResult []ResultById
	database.DB.Raw(`WITH acc AS (
	SELECT 
		oprt.login_at AS start_at, 
		oprt.logout_at AS end_at,
		oprt.operator AS operator_id,
		car.start_time,
		car.end_time,
		car.total_payment
	FROM operators AS oprt 
	LEFT JOIN car_models AS car 
	ON CAST(car.user_id AS bigint) = oprt.operator 
	WHERE car.user_id ~ '^\d+$'
	AND oprt.logout_at IS NOT NULL
	AND CAST(car.start_time AS TIMESTAMP) >= oprt.login_at
	AND CAST(car.end_time AS TIMESTAMP) <= oprt.logout_at
	AND oprt.operator = @operatorId
), 
opt AS (
	SELECT 
		oprt.login_at AS start_at, 
		oprt.logout_at AS end_at,
		oprt.operator AS operator_id 
	FROM operators as oprt 
	WHERE oprt.logout_at IS NOT NULL
	AND oprt.operator = @operatorId
),
elect AS(
	SELECT 
		SUM(acc.total_payment) AS total_payment, 
		acc.start_at, 
		acc.end_at, 
		acc.operator_id 
	FROM acc 
	GROUP BY acc.start_at, acc.end_at, acc.operator_id
)
SELECT 
	opt.*, 
	COALESCE(elect.total_payment, 0) AS total_payment 
FROM elect 
RIGHT JOIN opt ON elect.operator_id = opt.operator_id 
AND elect.start_at = opt.start_at 
AND elect.end_at = opt.end_at;`, sql.Named("operatorId", operatorId)).Scan(&operatorResult)
	return c.JSON(operatorResult)
}

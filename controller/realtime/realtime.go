package realtime

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var (
	parkingCounts = make(map[string]int)
	clients       = make(map[*websocket.Conn]bool)
	clientsMutex  sync.Mutex
)

type UpdateRequest struct {
	Total  int    `json:"total_payment"`
	ParkNo string `json:"parkno"`
}

func ResetParkingCount(parkNo string) error {
	if parkNo == "" {
		return fmt.Errorf("park number cannot be empty")
	}

	parkingCounts[parkNo] = 0

	fmt.Printf("Reset Park %s to 0\n", parkNo)
	fmt.Println("All parking counts:", parkingCounts)

	broadcastCount()

	return nil
}

func broadcastCount() {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	countData := parkingCounts
	for client := range clients {
		if err := client.WriteJSON(countData); err != nil {
			fmt.Printf("Error sending to client: %v\n", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// UpdateCount godoc
// @Summary Update the count value for a specific parking number
// @Description Adds the provided total to the existing count for the specified park number
// @Accept json
// @Produce json
// @Param request body UpdateRequest true "Total value to add and park number"
// @Success 200 {object} map[string]interface{} "Updated total value and park number"
// @Failure 400 {object} map[string]string "Error message"
// @Router /api/v1/update/count [put]
func UpdateCount(c *fiber.Ctx) error {
	var data UpdateRequest
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if data.ParkNo == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Park number is required"})
	}
	if data.Total < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Total cannot be negative"})
	}

	parkingCounts[data.ParkNo] += data.Total
	broadcastCount()

	return c.JSON(fiber.Map{
		"total_payment": parkingCounts[data.ParkNo],
		"park_no":       data.ParkNo,
	})
}

// GetAllCounts godoc
// @Summary Establish WebSocket connection for parking counts
// @Description Provides real-time updates of parking counts via WebSocket
// @Produce json
// @Success 101 {object} nil "WebSocket upgrade"
// @Router /api/v1/update/count [get]
func GetAllCounts(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func HandleWebSocketCount(c *websocket.Conn) {
	clientsMutex.Lock()
	clients[c] = true
	clientsMutex.Unlock()

	err := c.WriteJSON(parkingCounts)
	if err != nil {
		fmt.Printf("Error sending initial counts: %v\n", err)
		c.Close()
		return
	}

	defer func() {
		clientsMutex.Lock()
		delete(clients, c)
		clientsMutex.Unlock()
		c.Close()
	}()

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}

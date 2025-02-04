package util

import (
	"log"
	"sync"

	"park/database"

	"github.com/bits-and-blooms/bloom/v3"
)

var (
	vipFilter *bloom.BloomFilter
	vipMutex  sync.Mutex
)

func LoadVIPPlates() {
	vipMutex.Lock()
	defer vipMutex.Unlock()

	vipFilter = nil
	vipFilter = bloom.NewWithEstimates(1000000, 0.01)
	var plates []string
	err := database.DB.Raw("SELECT plate FROM tarifs").Scan(&plates).Error
	if err != nil {
		log.Println("Error occurred while loading VIP plates:", err)
		return
	}

	for _, plate := range plates {
		vipFilter.AddString(plate)
	}
	log.Println("VIP plates loaded successfully!")
}

func IsVIPPlate(plate string) bool {
	if plate == "" {
		return false
	}

	vipMutex.Lock()
	defer vipMutex.Unlock()

	if vipFilter == nil {
		log.Println("VIP Bloom Filter has not been loaded yet!")
		return false
	}

	return vipFilter.TestString(plate)
}

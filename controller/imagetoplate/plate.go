package imagetoplate

import (
	"log"
	modelscar "park/models/modelsCar"
	"path/filepath"
	"regexp"

	"github.com/fsnotify/fsnotify"
	"gorm.io/gorm"
)

func extractCarNumberFromFilename(filename string) string {
	re := regexp.MustCompile(`[A-Z0-9]+`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

func WatchDirectory(dir string, db *gorm.DB) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					imageUrl := filepath.Base(event.Name)
					carNumber := extractCarNumberFromFilename(imageUrl)
					if carNumber != "" {
						updateLatestRecordByCarNumber(db, carNumber, imageUrl)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()
	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
func updateLatestRecordByCarNumber(db *gorm.DB, carNumber, imageUrl string) {
	var latestCarStay modelscar.Car_Model

	if err := db.Where("car_number = ?", carNumber).Order("id DESC").First(&latestCarStay).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No records found for CarNumber:", carNumber, "- Skipping update.")
			return
		}
		log.Println("Error retrieving the latest record:", err)
		return
	}

	latestCarStay.Image_Url = imageUrl
	if err := db.Save(&latestCarStay).Error; err != nil {
		log.Println("Failed to update the latest record for CarNumber:", carNumber, "Error:", err)
	} else {
		log.Println("Latest record updated for CarNumber:", carNumber, "with ImageUrl:", imageUrl)
	}
}

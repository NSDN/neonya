package plate

import (
	"log"

	"gorm.io/gorm"
)

func InitPlateList(db *gorm.DB) {
	plates, err := GetPlates(db)

	if err != nil {
		log.Fatal(err)
	}

	if len(plates) > 0 {
		return
	}

	result := db.Create(&DefaultPlates)

	if result.Error != nil {
		log.Fatal(result.Error)
	}
}

func GetPlates(db *gorm.DB) ([]Plate, error) {
	var plates []Plate
	result := db.Find(&plates)

	if result.Error != nil {
		return nil, result.Error
	}

	return plates, nil
}

package board

import (
	"log"

	"gorm.io/gorm"
)

func InitBoardList(db *gorm.DB) {
	boards, err := GetBoards(db)

	if err != nil {
		log.Fatal(err)
	}

	if len(boards) > 0 {
		return
	}

	result := db.Create(&DefaultBoards)

	if result.Error != nil {
		log.Fatal(result.Error)
	}
}

func GetBoards(db *gorm.DB) ([]Board, error) {
	var boards []Board
	result := db.Find(&boards)

	if result.Error != nil {
		return nil, result.Error
	}

	return boards, nil
}

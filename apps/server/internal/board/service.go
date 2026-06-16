package board

import (
	"gorm.io/gorm"
)

func GetBoards(db *gorm.DB) ([]Board, error) {
	var boards []Board
	result := db.Find(&boards)

	if result.Error != nil {
		return nil, result.Error
	}

	return boards, nil
}

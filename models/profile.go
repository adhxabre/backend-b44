package models

import "time"

type Profile struct {
	ID        int                  `json:"id"`
	Phone     string               `json:"phone" gorm:"type: varchar(50)"`
	Gender    string               `json:"gender" gorm:"type: varchar(100)"`
	Address   string               `json:"address"  gorm:"type: varchar(255)"`
	UserID    int                  `json:"user_id"`
	User      UsersProfileResponse `json:"user"`
	CreatedAt time.Time            `json:"-"`
	UpdatedAt time.Time            `json:"-"`
}

type ProfileResponse struct {
	Phone   string `json:"phone"`
	Gender  string `json:"gender"`
	Address string `json:"address"`
	UserID  int    `json:"-"`
}

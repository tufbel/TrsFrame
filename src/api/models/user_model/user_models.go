package user_model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Uid      string `json:"uid" gorm:"primaryKey; type:char(27)"`     // 用户ID
	Name     string `json:"name" form:"name" gorm:"size:64"`          // 名称
	Email    string `json:"email" form:"email" gorm:"index; size:64"` // 邮箱
	Password string `json:"password" form:"password" gorm:"size:64"`  // 密码

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"-" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (*User) TableName() string {
	// 指定表名
	return "users"
}

// schema
type CreateUserSchema struct {
	Name     string `form:"name" binding:"required" required:"true"`     // 名称
	Email    string `form:"email" binding:"required" required:"true"`    // 邮箱
	Password string `form:"password" binding:"required" required:"true"` // 密码
}

type UpdateUserSchema struct {
	Name     string `form:"name"`     // 用户名
	Email    string `form:"email"`    // 邮箱
	Password string `form:"password"` // 密码
}

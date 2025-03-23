package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	// 存秒数 / 存纳秒数 / 存毫秒数
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	return dao.db.WithContext(ctx).Create(&u).Error // 必须要取指针, gorm规定
}

// User 直接对应于数据库表结构，有些人叫做 entity，有些人叫做 model，也有人叫做 PO (persistent object)
type User struct {
	Id       int64  `gorm:"primaryKey, autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}

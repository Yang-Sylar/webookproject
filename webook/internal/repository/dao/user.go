package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate = errors.New("邮箱冲突，手机号码冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDAO interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	Insert(ctx context.Context, u User) error
	UpdateById(ctx context.Context, u User) error
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *GORMUserDAO {
	return &GORMUserDAO{
		db: db,
	}
}
func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	// 存秒数 / 存纳秒数 / 存毫秒数
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now

	err := dao.db.WithContext(ctx).Create(&u).Error // 必须要取指针, gorm规定

	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueIndexErrNo uint16 = 1062
		if mysqlErr.Number == uniqueIndexErrNo {
			// 邮箱冲突，手机号码重复
			return ErrUserDuplicate
		}
	}

	return err
}

func (dao *GORMUserDAO) UpdateById(ctx context.Context, u User) error {
	return dao.db.WithContext(ctx).Model(&u).Where("id = ?", u.Id).
		Updates(map[string]any{
			"Utime":    time.Now().UnixMilli(),
			"Nickname": u.Nickname,
			"Birthday": u.Birthday,
			"AboutMe":  u.AboutMe,
		}).Error
}

// User 直接对应于数据库表结构，有些人叫做 entity，有些人叫做 model，也有人叫做 PO (persistent object)
type User struct {
	Id       int64          `gorm:"primaryKey, autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	Ctime    int64 // 创建时间
	Utime    int64 // 更新时间
	// 唯一索引允许有多个空值
	// 但不能有多个 ""
	Phone sql.NullString `gorm:"unique"`
	// Profile
	Nickname string `gorm:"type=varchar(128)"`
	Birthday int64
	AboutMe  string `gorm:"type=varchar(4096)"`
}

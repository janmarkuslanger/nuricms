package base

import (
	"gorm.io/gorm"
)

type QueryOption func(*gorm.DB) *gorm.DB

type CRUDRepository[T any] interface {
	Create(entity *T) error
	Save(entity *T) error
	Delete(entity *T) error
	FindByID(id uint, opts ...QueryOption) (*T, error)
	List(page, pageSize int, opts ...QueryOption) ([]T, int64, error)
}

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

func (r *BaseRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *BaseRepository[T]) Save(entity *T) error {
	return r.db.Save(entity).Error
}

func (r *BaseRepository[T]) Delete(entity *T) error {
	return r.db.Delete(entity).Error
}

func (r *BaseRepository[T]) FindByID(id uint, opts ...QueryOption) (*T, error) {
	db := r.db
	for _, o := range opts {
		db = o(db)
	}
	var entity T
	err := db.First(&entity, id).Error
	return &entity, err
}

func (r *BaseRepository[T]) List(page, pageSize int, opts ...QueryOption) ([]T, int64, error) {
	db := r.db
	for _, o := range opts {
		db = o(db)
	}
	var total int64
	err := db.Model(new(T)).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	var items []T
	err = db.Offset(offset).Limit(pageSize).Find(&items).Error
	return items, total, err
}

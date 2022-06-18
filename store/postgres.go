package store

import (
	"context"
	"log"
	"os"

	"go-inventory/errors"
	"go-inventory/objects"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type pg struct {
	db *gorm.DB
}

// NewPostgresStockStore returns a postgres implementation of Stock store
func NewPostgresStockStore(conn string) IStockStore {
	// create database connection
	db, err := gorm.Open(postgres.Open(conn),
		&gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "", log.LstdFlags),
				logger.Config{
					LogLevel: logger.Info,
					Colorful: true,
				},
			),
		},
	)
	if err != nil {
		panic("Enable to connect to database: " + err.Error())
	}
	if err := db.AutoMigrate(&objects.Stock{}); err != nil {
		panic("Enable to migrate database: " + err.Error())
	}
	// return store implementation
	return &pg{db: db}
}

func (p *pg) Get(ctx context.Context, in *objects.GetRequest) (*objects.Stock, error) {
	evt := &objects.Stock{}
	// take event where id == uid from database
	err := p.db.WithContext(ctx).Take(evt, "id = ?", in.ID).Error
	if err == gorm.ErrRecordNotFound {
		// not found
		return nil, errors.ErrStockNotFound
	}
	return evt, err
}

func (p *pg) List(ctx context.Context, in *objects.ListRequest) ([]*objects.Stock, error) {
	if in.Limit == 0 || in.Limit > objects.MaxListLimit {
		in.Limit = objects.MaxListLimit
	}
	query := p.db.WithContext(ctx).Limit(in.Limit)
	if in.After != "" {
		query = query.Where("id > ?", in.After)
	}
	if in.Name != "" {
		query = query.Where("name ilike ?", "%"+in.Name+"%")
	}
	list := make([]*objects.Stock, 0, in.Limit)
	err := query.Order("id").Find(&list).Error
	return list, err
}

func (p *pg) Create(ctx context.Context, in *objects.CreateRequest) error {
	if in.Stock == nil {
		return errors.ErrObjectIsRequired
	}
	in.Stock.ID = GenerateUniqueID()
	in.Stock.CreatedOn = p.db.NowFunc()
	return p.db.WithContext(ctx).
		Create(in.Stock).
		Error
}

func (p *pg) UpdateDetails(ctx context.Context, in *objects.UpdateDetailsRequest) error {
	evt := &objects.Stock{
		ID:           in.ID,
		Name:         in.Name,
		Price:        in.Price,
		Availability: in.Availability,
		IsActive:     in.IsActive,
		UpdatedOn:    p.db.NowFunc(),
	}
	log.Println(evt)
	return p.db.WithContext(ctx).Model(evt).
		Select("id", "name", "price", "availability", "is_active", "updated_on").
		Updates(evt).
		Error
}

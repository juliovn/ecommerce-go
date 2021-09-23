package models

import "github.com/jinzhu/gorm"

type Services struct {
	Item ItemService
	User UserService
	db   *gorm.DB
}

// NewServices is a helper function that will load all services defined on struct to facilitate coding
func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)

	// construct services
	return &Services{
		User: NewUserService(db),
		Item: NewItemService(db),
		db:   db,
	}, nil
}

// Close closes the database connection
func (s *Services) Close() error {
	return s.db.Close()
}

// AutoMigrate will automatically attempt to migrate all tables
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Item{}).Error
}

// DestructiveReset drops all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Item{}).Error

	if err != nil {
		return err
	}

	return s.AutoMigrate()
}

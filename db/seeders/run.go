package seeders

import "gorm.io/gorm"

func Run(db *gorm.DB) error {
	seeders := []Seeder{
		UserSeeder{},
	}

	for _, seeder := range seeders {
		if err := seeder.Seed(db); err != nil {
			return err
		}
	}

	return nil
}

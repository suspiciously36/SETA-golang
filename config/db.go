package config

import (
	"github.com/seta-namnv-6798/go-apis/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=goapi_user password=goapi_password dbname=goapi_db port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	// Auto-migrate all models in the correct order
	// 1. First create base tables without foreign key dependencies
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Team{})

	// 2. Then create junction/relationship tables
	db.AutoMigrate(&models.TeamMember{})
	db.AutoMigrate(&models.TeamManager{})

	// 3. Create content tables
	db.AutoMigrate(&models.Folder{})
	db.AutoMigrate(&models.Note{})

	// 4. Finally create sharing tables
	db.AutoMigrate(&models.FolderShare{})
	db.AutoMigrate(&models.NoteShare{})

	DB = db
}

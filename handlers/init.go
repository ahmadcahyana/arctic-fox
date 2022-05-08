package handlers

import (
	"fmt"
	"github.com/harranali/authority"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ahmadcahyana/arctic-fox/config"
	"github.com/ahmadcahyana/arctic-fox/datatransfers"
	"github.com/ahmadcahyana/arctic-fox/models"
)

var Handler HandlerFunc

type HandlerFunc interface {
	AuthenticateUser(credentials datatransfers.UserLogin) (token string, err error)
	RegisterUser(credentials datatransfers.UserSignup) (err error)

	RetrieveUser(username string) (user models.User, err error)
	UpdateUser(id uint, user datatransfers.UserUpdate) (err error)
}

type module struct {
	db *dbEntity
}

type dbEntity struct {
	conn      *gorm.DB
	userOrmer models.UserOrmer
}

func InitializeHandler() (err error) {
	// Initialize DB
	var db *gorm.DB
	db, err = gorm.Open(postgres.Open(
		fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			config.AppConfig.DBHost, config.AppConfig.DBPort, config.AppConfig.DBDatabase,
			config.AppConfig.DBUsername, config.AppConfig.DBPassword),
	), &gorm.Config{})
	if err != nil {
		log.Println("[INIT] failed connecting to PostgreSQL")
		return
	}
	log.Println("[INIT] connected to PostgreSQL")
	errMigration := db.AutoMigrate(&models.User{})
	if errMigration != nil {
		log.Println("[INIT] failed run some migrations")
		return
	}
	auth := authority.New(authority.Options{
		TablesPrefix: "authority_",
		DB:           db,
	})
	errAuthority := auth.CreateRole("super_admin")
	if errAuthority != nil {
		log.Println("[INIT] failed create Role annd Permission")
		return
	}
	errAuthority = auth.CreatePermission("create")
	if errAuthority != nil {
		log.Println("[INIT] failed create Role annd Permission")
		return
	}
	errAuthority = auth.CreatePermission("read")
	if errAuthority != nil {
		log.Println("[INIT] failed create Role annd Permission")
		return
	}
	errAuthority = auth.CreatePermission("update")
	if errAuthority != nil {
		log.Println("[INIT] failed create Role annd Permission")
		return
	}
	errAuthority = auth.CreatePermission("delete")
	if errAuthority != nil {
		log.Println("[INIT] failed create Role annd Permission")
		return
	}
	errAuthority = auth.AssignPermissions("super_admin", []string{
		"create",
		"read",
		"update",
		"delete",
	})
	if errAuthority != nil {
		log.Println("[INIT] failed create Role annd Permission")
		return
	}
	errAuthority = auth.AssignRole(1, "super_admin")
	if errAuthority != nil {
		log.Println("[INIT] failed create Role annd Permission")
		return
	}
	_, errAuthority = auth.CheckRole(1, "super_admin")
	if errAuthority != nil {
		log.Println("[INIT] failed create Role annd Permission")
		return
	}
	_, errAuthority = auth.CheckPermission(1, "create")
	if errAuthority != nil {
		log.Println("[INIT] failed create Role annd Permission")
		return
	}
	_, errAuthority = auth.CheckRolePermission("super_admin", "create")
	if errAuthority != nil {
		log.Println("[INIT] failed create Role annd Permission")
		return
	}

	// Compose handler modules
	Handler = &module{
		db: &dbEntity{
			conn:      db,
			userOrmer: models.NewUserOrmer(db),
		},
	}
	return
}

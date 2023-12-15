package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string
}

type UserRepository interface {
	CreateUser(data interface{}) error
	GetUserByID(id uint, result interface{}) error
	UpdateUser(data interface{}) error
	DeleteUser(data interface{}) error
}

type DBRepository struct {
	DB *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{DB: db}
}

func (repo *DBRepository) CreateUser(data interface{}) error {
	return repo.DB.Create(data).Error
}

func (repo *DBRepository) GetUserByID(id uint, result interface{}) error {
	return repo.DB.First(result, id).Error
}

func (repo *DBRepository) UpdateUser(data interface{}) error {
	return repo.DB.Save(data).Error
}

func (repo *DBRepository) DeleteUser(data interface{}) error {
	return repo.DB.Delete(data).Error
}

func main() {
	dsn := "postgresql://root:@localhost:26257/defaultdb?sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&User{})

	userRepository := NewDBRepository(db)

	app := fiber.New()

	app.Post("/users", func(c *fiber.Ctx) error {
		newUser := new(User)
		if err := c.BodyParser(newUser); err != nil {
			return err
		}
		if err := userRepository.CreateUser(newUser); err != nil {
			return err
		}
		return c.JSON(newUser)
	})

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			return err
		}
		var user User
		if err := userRepository.GetUserByID(uint(id), &user); err != nil {
			return err
		}
		return c.JSON(user)
	})

	app.Put("/users/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			return err
		}
		var updatedUser User
		if err := c.BodyParser(&updatedUser); err != nil {
			return err
		}
		updatedUser.ID = uint(id) // Set the ID to update the specific user
		if err := userRepository.UpdateUser(&updatedUser); err != nil {
			return err
		}
		return c.JSON(updatedUser)
	})

	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			return err
		}
		var user User
		if err := userRepository.GetUserByID(uint(id), &user); err != nil {
			return err
		}
		if err := userRepository.DeleteUser(&user); err != nil {
			return err
		}
		return c.SendString("User deleted")
	})

	// Start Fiber server
	app.Listen(":3000")
}

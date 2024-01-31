package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type Network struct {
	SubnetID int    `json:"id"`
	ParentID interface{}   `json:"parent_id,omitempty"`
	Name     string `json:"name"`
	Subnet   string `json:"subnet"`
}

func getNetworks(c *fiber.Ctx, db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM networks")
	if handleDBError(c, err) {
		return err
	}
	defer rows.Close()

	var networks []Network
	for rows.Next() {
		var net Network
		if err := rows.Scan(&net.SubnetID, &net.ParentID, &net.Name, &net.Subnet); handleDBError(c, err) {
			return err
		}
		networks = append(networks, net)
	}
	return c.Status(fiber.StatusOK).JSON(networks)
}

func createNetwork(c *fiber.Ctx, db *sql.DB) error {
	newNetwork := Network{}
	if err := c.BodyParser(&newNetwork); handleError(c, err, fiber.StatusBadRequest) {
		return err
	}

	parentID := c.Params("parent_id")
	if parentID != "" {
		parentIDInt, err := strconv.Atoi(parentID)
		if handleError(c, err, fiber.StatusBadRequest) {
			return err
		}
		newNetwork.ParentID = &parentIDInt
	}

	var subnetID int
	err := db.QueryRow("INSERT INTO networks (name, subnet, parent_id) VALUES ($1, $2, $3) RETURNING subnet_id",
		newNetwork.Name, newNetwork.Subnet, newNetwork.ParentID).Scan(&subnetID)
	if handleDBError(c, err) {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": subnetID, "subnet": newNetwork.Subnet})
}

func handleError(c *fiber.Ctx, err error, statusCode int) bool {
	if err != nil {
		log.Printf("Error: %v", err)
		c.Status(statusCode).SendString(err.Error())
		return true
	}
	return false
}

func handleDBError(c *fiber.Ctx, err error) bool {
	return handleError(c, err, fiber.StatusInternalServerError)
}

func main() {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dsn)
	if handleError(nil, err, 0) {
		log.Fatal(err)
	}

	app := fiber.New()
	app.Static("/", "./web")

	api := app.Group("/api")
	api.Get("/networks", func(c *fiber.Ctx) error {
		return getNetworks(c, db)
	})
	api.Post("/networks", func(c *fiber.Ctx) error {
		return createNetwork(c, db)
	})

	log.Fatal(app.Listen(":3000"))
}
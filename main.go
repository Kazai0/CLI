package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func main() {
	var rootCommand = &cobra.Command{}
	var projectName, projectPath string

	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create boilerplate for a new project with Fiber and PostgreSQL",
		Run: func(cmd *cobra.Command, args []string) {
			if projectName == "" {
				fmt.Println("You must supply a project name.")
				return
			}
			if projectPath == "" {
				fmt.Println("You must supply a project path.")
				return
			}
			fmt.Println("Creating project...")

			globalPath := filepath.Join(projectPath, projectName)

			if _, err := os.Stat(globalPath); err == nil {
				fmt.Println("Project directory already exists.")
				return
			}
			if err := os.Mkdir(globalPath, os.ModePerm); err != nil {
				log.Fatal(err)
			}

			// Initialize Go module
			startGo := exec.Command("go", "mod", "init", projectName)
			startGo.Dir = globalPath
			startGo.Stdout = os.Stdout
			startGo.Stderr = os.Stderr
			err := startGo.Run()
			if err != nil {
				log.Fatal(err)
			}

			// Install dependencies
			dependencies := []string{
				"github.com/gofiber/fiber/v2",
				"github.com/rs/zerolog",
				"github.com/jmoiron/sqlx",
				"github.com/lib/pq",
			}
			for _, dep := range dependencies {
				installCmd := exec.Command("go", "get", dep)
				installCmd.Dir = globalPath
				installCmd.Stdout = os.Stdout
				installCmd.Stderr = os.Stderr
				if err := installCmd.Run(); err != nil {
					log.Fatal(err)
				}
			}

			// Create directories
			directories := []string{
				"cmd", "internal/handler", "configs",
			}
			for _, dir := range directories {
				if err := os.MkdirAll(filepath.Join(globalPath, dir), os.ModePerm); err != nil {
					log.Fatal(err)
				}
			}

			// Write files
			files := map[string]func(string) error{
				"cmd/main.go":                 WriteMainFile,
				"internal/handler/handler.go": WriteHandlerFile,
			}

			//path to docker-compose
			dockerComposePath := filepath.Join(globalPath, "docker-compose.yml")

			// create docker-compose
			if err := WriteDockerComposeFile(dockerComposePath); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Docker Compose file created at", dockerComposePath)
			for filePath, writer := range files {
				fullPath := filepath.Join(globalPath, filePath)
				if err := writer(fullPath); err != nil {
					log.Fatal(err)
				}
			}

		},
	}

	cmd.Flags().StringVarP(&projectName, "name", "n", "", "Name of the project")
	cmd.Flags().StringVarP(&projectPath, "path", "p", "", "Path where the project will be created")

	rootCommand.AddCommand(cmd)
	rootCommand.Execute()
}

func WriteMainFile(mainPath string) error {
	packageContent := []byte(`package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Initialize database
	db, err := sqlx.Connect("postgres", "user=postgres password=secret dbname=mydb sslmode=disable")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to the database")
	}
	defer db.Close()

	// Initialize Fiber app
	app := fiber.New()

	// Add routes
	app.Get("/health", func(c *fiber.Ctx) error {
		fmt.Println("Health check")
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/simple", func(c *fiber.Ctx) error {
		return c.SendString("Hello from Fiber!")
	})

	// Graceful shutdown
	go func() {
		if err := app.Listen(":3000"); err != nil {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error().Err(err).Msg("Error during shutdown")
	}
}
`)

	return os.WriteFile(mainPath, packageContent, 0644)
}

func WriteHandlerFile(handlerPath string) error {
	packageContent := []byte(`package handler

import "github.com/gofiber/fiber/v2"

// SimpleHandler demonstrates a simple GET endpoint
func SimpleHandler(c *fiber.Ctx) error {
	return c.SendString("Hello from handler!")
}
`)

	return os.WriteFile(handlerPath, packageContent, 0644)
}
func WriteDockerComposeFile(dockerComposePath string) error {
	composeContent := []byte(`version: '3.8'

services:
  postgres:
    image: postgres:15
    container_name: postgres_container
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: mydb
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
`)

	return os.WriteFile(dockerComposePath, composeContent, 0644)
}

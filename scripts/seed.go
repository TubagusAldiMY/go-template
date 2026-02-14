package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/TubagusAldiMY/go-template/internal/infrastructure/config"
	"github.com/TubagusAldiMY/go-template/pkg/crypto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to database successfully")

	// Create password hasher
	hasher := crypto.NewPasswordHasher(bcrypt.DefaultCost)

	// Seed admin user
	adminPassword, err := hasher.Hash("Admin123!")
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	adminID := uuid.New().String()
	now := time.Now()

	_, err = pool.Exec(context.Background(), `
		INSERT INTO users (id, email, username, password, full_name, role, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (email) DO NOTHING
	`, adminID, "admin@example.com", "admin", adminPassword, "System Administrator", "admin", "active", now, now)

	if err != nil {
		log.Printf("Warning: Failed to seed admin user (might already exist): %v", err)
	} else {
		log.Println("✓ Admin user seeded successfully")
		log.Println("  Email: admin@example.com")
		log.Println("  Password: Admin123!")
	}

	// Seed test user
	userPassword, err := hasher.Hash("User123!")
	if err != nil {
		log.Fatalf("Failed to hash user password: %v", err)
	}

	userID := uuid.New().String()

	_, err = pool.Exec(context.Background(), `
		INSERT INTO users (id, email, username, password, full_name, role, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (email) DO NOTHING
	`, userID, "user@example.com", "testuser", userPassword, "Test User", "user", "active", now, now)

	if err != nil {
		log.Printf("Warning: Failed to seed test user (might already exist): %v", err)
	} else {
		log.Println("✓ Test user seeded successfully")
		log.Println("  Email: user@example.com")
		log.Println("  Password: User123!")
	}

	log.Println("\n✅ Database seeding completed!")
}

package migration

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
	"gorm.io/gorm"
)

type Runner struct {
	db     *gorm.DB
	sqlDB  *sql.DB
	logger *logger.Logger
}

type Migration struct {
	ID        int
	Name      string
	Content   string
	AppliedAt *time.Time
}

func NewRunner(db *gorm.DB, sqlDB *sql.DB, logger *logger.Logger) *Runner {
	return &Runner{
		db:     db,
		sqlDB:  sqlDB,
		logger: logger,
	}
}

// RunMigrations executes all pending migrations from the migration directory
func (r *Runner) RunMigrations(ctx context.Context, migrationDir string) error {
	// Create migrations table if it doesn't exist
	if err := r.createMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	migrationFiles, err := r.getMigrationFiles(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Get applied migrations from database
	appliedMigrations, err := r.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Execute pending migrations
	for _, file := range migrationFiles {
		if !strings.HasSuffix(file, ".up.sql") {
			continue // Skip down migrations and other files
		}

		migrationID := r.extractMigrationID(file)
		if migrationID == 0 {
			r.logger.Warn(ctx, "Skipping file with invalid migration ID", "file", file)
			continue
		}

		// Check if migration is already applied
		if r.isMigrationApplied(migrationID, appliedMigrations) {
			r.logger.Debug(ctx, "Migration already applied", "id", migrationID, "file", file)
			continue
		}

		// Read and execute migration
		content, err := os.ReadFile(filepath.Join(migrationDir, file))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		r.logger.Info(ctx, "Applying migration", "id", migrationID, "file", file)

		if err := r.executeMigration(ctx, migrationID, file, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		r.logger.Info(ctx, "Migration applied successfully", "id", migrationID, "file", file)
	}

	r.logger.Info(ctx, "All migrations completed successfully")
	return nil
}

// createMigrationsTable creates the migrations tracking table
func (r *Runner) createMigrationsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id INT NOT NULL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE KEY unique_migration (id)
		)
	`

	if err := r.sqlDB.QueryRowContext(ctx, query).Err(); err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}

// getMigrationFiles returns a sorted list of migration files
func (r *Runner) getMigrationFiles(migrationDir string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(migrationDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			files = append(files, d.Name())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort files to ensure they run in order
	sort.Strings(files)
	return files, nil
}

// getAppliedMigrations returns a map of applied migration IDs
func (r *Runner) getAppliedMigrations(ctx context.Context) (map[int]bool, error) {
	appliedMigrations := make(map[int]bool)

	rows, err := r.sqlDB.QueryContext(ctx, "SELECT id FROM migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		appliedMigrations[id] = true
	}

	return appliedMigrations, rows.Err()
}

// extractMigrationID extracts the numeric ID from migration filename
func (r *Runner) extractMigrationID(filename string) int {
	// Extract number from filename like "001_create_users_table.up.sql"
	parts := strings.Split(filename, "_")
	if len(parts) == 0 {
		return 0
	}

	idStr := parts[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0
	}

	return id
}

// isMigrationApplied checks if a migration has already been applied
func (r *Runner) isMigrationApplied(id int, appliedMigrations map[int]bool) bool {
	return appliedMigrations[id]
}

// executeMigration executes a single migration and records it
func (r *Runner) executeMigration(ctx context.Context, id int, name, content string) error {
	// Start transaction
	tx, err := r.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute migration content
	// Split by semicolon to handle multiple statements
	statements := strings.Split(content, ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("failed to execute statement: %s, error: %w", stmt, err)
		}
	}

	// Record migration as applied
	if _, err := tx.ExecContext(ctx, "INSERT INTO migrations (id, name) VALUES (?, ?)", id, name); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	return nil
}

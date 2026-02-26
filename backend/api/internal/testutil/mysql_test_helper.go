// MySQL test helper - provides isolated database for each test
package testutil

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"dmh/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLTestConfig holds configuration for MySQL test databases
type MySQLTestConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	BaseDB   string // Base database to connect (usually 'dmh' or 'mysql')
}

// GetMySQLTestConfig returns MySQL configuration from environment or defaults
func GetMySQLTestConfig() *MySQLTestConfig {
	return &MySQLTestConfig{
		Host:     getEnv("MYSQL_TEST_HOST", "127.0.0.1"),
		Port:     getEnv("MYSQL_TEST_PORT", "3306"),
		User:     getEnv("MYSQL_TEST_USER", "root"),
		Password: getEnv("MYSQL_TEST_PASSWORD", "Admin168"),
		BaseDB:   getEnv("MYSQL_TEST_DB", "dmh"),
	}
}

// SetupMySQLTestDB creates an isolated MySQL database for a test
// Returns the database connection and the test database name
func SetupMySQLTestDB(t *testing.T) (*gorm.DB, string) {
	config := GetMySQLTestConfig()

	// Generate unique test database name
	testDBName := generateTestDBName(t)

	// Connect to base database to create test database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&multiStatements=true",
		config.User, config.Password, config.Host, config.Port, config.BaseDB)

	adminDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to MySQL admin database: %v", err)
	}

	// Create test database
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", testDBName)
	if err := adminDB.Exec(createDBSQL).Error; err != nil {
		t.Fatalf("Failed to create test database %s: %v", testDBName, err)
	}

	// Close admin connection
	sqlDB, _ := adminDB.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	// Connect to test database
	testDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.User, config.Password, config.Host, config.Port, testDBName)

	db, err := gorm.Open(mysql.Open(testDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database %s: %v", testDBName, err)
	}

	// Configure connection pool
	sqlDB, err = db.DB()
	if err != nil {
		t.Fatalf("Failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	// Auto migrate all models
	if err := migrateTestSchema(db); err != nil {
		t.Fatalf("Failed to migrate test schema: %v", err)
	}

	// Register cleanup
	t.Cleanup(func() {
		CleanupMySQLTestDB(t, testDBName)
	})

	return db, testDBName
}

// CleanupMySQLTestDB drops the test database
func CleanupMySQLTestDB(t *testing.T, dbName string) {
	config := GetMySQLTestConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.BaseDB)

	adminDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Logf("Warning: Failed to connect for cleanup: %v", err)
		return
	}

	// Drop test database
	dropSQL := fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)
	if err := adminDB.Exec(dropSQL).Error; err != nil {
		t.Logf("Warning: Failed to drop test database %s: %v", dbName, err)
	}

	sqlDB, _ := adminDB.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

// migrateTestSchema creates all necessary tables and indexes for testing
func migrateTestSchema(db *gorm.DB) error {
	// Create users table first (many tests depend on it)
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			username VARCHAR(50) NOT NULL,
			password VARCHAR(255) NOT NULL,
			phone VARCHAR(20) NOT NULL,
			email VARCHAR(100),
			avatar VARCHAR(255),
			real_name VARCHAR(50),
			role VARCHAR(50) NOT NULL DEFAULT 'participant',
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			login_attempts INT DEFAULT 0,
			locked_until DATETIME,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE INDEX idx_username (username),
			UNIQUE INDEX idx_phone (phone),
			INDEX idx_role (role),
			INDEX idx_status (status)
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create roles table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS roles (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(50) NOT NULL UNIQUE,
			description VARCHAR(255),
			permissions JSON,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create roles table: %w", err)
	}

	// Create brands table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS brands (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(100) NOT NULL,
			logo VARCHAR(255),
			description TEXT,
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create brands table: %w", err)
	}

	// Try AutoMigrate for remaining models (non-critical if some fail)
	models := []interface{}{
		&model.Role{},
		&model.Permission{},
		&model.UserRole{},
		&model.RolePermission{},
		&model.UserBrand{},
		&model.Menu{},
		&model.Member{},
		&model.Order{},
		&model.VerificationRecord{},
		&model.Distributor{},
		&model.DistributorLevelReward{},
		&model.DistributorReward{},
		&model.UserFeedback{},
		&model.Withdrawal{},
		&model.PosterTemplate{},
		&model.PasswordPolicy{},
		&model.AuditLog{},
		&model.SyncLog{},
		&model.PageConfig{},
	}

	for _, m := range models {
		db.AutoMigrate(m)
	}

	// Create campaigns table
	db.Exec(`
		CREATE TABLE IF NOT EXISTS campaigns (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			brand_id BIGINT NOT NULL,
			name VARCHAR(200) NOT NULL,
			description TEXT,
			form_fields JSON,
			reward_rule DECIMAL(10,2) NOT NULL DEFAULT 0.00,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			enable_distribution TINYINT(1) NOT NULL DEFAULT 0,
			distribution_level INT NOT NULL DEFAULT 1,
			distribution_rewards JSON,
			payment_config JSON,
			poster_template_id BIGINT DEFAULT 1,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			INDEX idx_brand_id (brand_id),
			INDEX idx_status (status)
		)
	`)

	// Create unique index for order duplicate guard
	if err := db.Exec("ALTER TABLE orders ADD UNIQUE KEY uk_orders_campaign_phone (campaign_id, phone)").Error; err != nil {
		errMsg := strings.ToLower(err.Error())
		if !strings.Contains(errMsg, "duplicate key name") {
			return fmt.Errorf("failed to create unique index uk_orders_campaign_phone: %w", err)
		}
	}

	return nil
}

// generateTestDBName creates a unique database name for the test
func generateTestDBName(t *testing.T) string {
	// Sanitize test name for MySQL identifier
	testName := strings.ReplaceAll(t.Name(), "/", "_")
	testName = strings.ReplaceAll(testName, "-", "_")
	testName = strings.ReplaceAll(testName, " ", "_")

	maxNameLen := 25
	if len(testName) > maxNameLen {
		testName = testName[:20] + fmt.Sprintf("_%d", rand.Intn(1000))
	}

	timestamp := time.Now().Format("150405")
	suffix := rand.Intn(1000)

	return fmt.Sprintf("t_%s_%s_%03d", testName, timestamp, suffix)
}

// getEnv is defined in env.go

// SkipIfNoMySQL skips the test if MySQL is not available
func SkipIfNoMySQL(t *testing.T) {
	t.Helper()
	config := GetMySQLTestConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.BaseDB)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("SKIP_REASON:%s | MySQL not available: %v", SkipReasonMySQLUnavailable, err)
		return
	}

	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

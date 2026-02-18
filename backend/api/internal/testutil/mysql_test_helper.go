// MySQL test helper - provides isolated database for each test
package testutil

import (
	"fmt"
	"math/rand"
	"os"
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
	// Auto migrate models
	models := []interface{}{
		&model.Campaign{},
		&model.Order{},
		&model.VerificationRecord{},
		&model.Distributor{},
		&model.DistributorLevelReward{},
		&model.DistributorReward{},
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.Brand{},
		&model.Member{},
		&model.UserFeedback{},
		&model.Withdrawal{},
		&model.PosterTemplate{},
		&model.PasswordPolicy{},
		&model.AuditLog{},
		&model.SyncLog{},
	}

	for _, m := range models {
		if err := db.AutoMigrate(m); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", m, err)
		}
	}

	// Create unique index for order duplicate guard (matching production)
	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS uk_orders_campaign_phone ON orders(campaign_id, phone)").Error; err != nil {
		return fmt.Errorf("failed to create unique index: %w", err)
	}

	return nil
}

// generateTestDBName creates a unique database name for the test
func generateTestDBName(t *testing.T) string {
	// Sanitize test name for MySQL identifier
	testName := strings.ReplaceAll(t.Name(), "/", "_")
	testName = strings.ReplaceAll(testName, "-", "_")
	testName = strings.ReplaceAll(testName, " ", "_")

	// Truncate if too long (MySQL identifier limit is 64 chars)
	if len(testName) > 40 {
		testName = testName[:40]
	}

	// Add timestamp and random suffix
	timestamp := time.Now().Format("20060102_150405")
	suffix := rand.Intn(10000)

	return fmt.Sprintf("test_%s_%s_%04d", testName, timestamp, suffix)
}

// getEnv returns environment variable value or default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SkipIfNoMySQL skips the test if MySQL is not available
func SkipIfNoMySQL(t *testing.T) {
	config := GetMySQLTestConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.BaseDB)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("MySQL not available, skipping test: %v", err)
		return
	}

	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

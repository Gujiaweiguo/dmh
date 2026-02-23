package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"dmh/api/internal/config"
	"dmh/api/internal/handler"
	"dmh/api/internal/svc"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/dmh-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	applyEnvOverrides(&c)

	// 配置 CORS - 允许所有来源
	server := rest.MustNewServer(c.RestConf, rest.WithCors("*"))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)

	// 在注册其他路由之前，先注册静态文件路由
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/api/v1/posters/:filename",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			pathParts := strings.Split(r.URL.Path, "/")

			filename := pathParts[len(pathParts)-1]


			if filename == "" {
				http.Error(w, "Invalid filename", http.StatusBadRequest)
				return
			}

			filePath := filepath.Join("/opt/data/posters", filename)

			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "image/png")
			w.Header().Set("Cache-Control", "public, max-age=86400")
			http.ServeFile(w, r, filePath)
		}),
	})

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

func applyEnvOverrides(c *config.Config) {
	if v := strings.TrimSpace(os.Getenv("APP_HOST")); v != "" {
		c.Host = v
	}
	if v := strings.TrimSpace(os.Getenv("APP_PORT")); v != "" {
		if port, err := strconv.Atoi(v); err == nil && port > 0 {
			c.Port = port
		}
	}
	if v := strings.TrimSpace(os.Getenv("JWT_SECRET")); v != "" {
		c.Auth.AccessSecret = v
	}

	dbHost := strings.TrimSpace(os.Getenv("DB_HOST"))
	dbPort := strings.TrimSpace(os.Getenv("DB_PORT"))
	dbUser := strings.TrimSpace(os.Getenv("DB_USER"))
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := strings.TrimSpace(os.Getenv("DB_NAME"))

	if dbHost == "" && dbPort == "" && dbUser == "" && dbPassword == "" && dbName == "" {
		return
	}

	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbUser == "" {
		dbUser = "root"
	}
	if dbName == "" {
		dbName = "dmh"
	}

	cfg := mysqlDriver.NewConfig()
	cfg.User = dbUser
	cfg.Passwd = dbPassword
	cfg.Net = "tcp"
	cfg.Addr = dbHost + ":" + dbPort
	cfg.DBName = dbName
	cfg.ParseTime = true
	cfg.Params = map[string]string{
		"charset": "utf8mb4",
		"loc":     "Local",
	}

	c.Mysql.DataSource = cfg.FormatDSN()
}

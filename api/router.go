package router

import (
	"fmt"
	"log"

	general "rlp-member-service/api/http"
	"rlp-member-service/api/http/middleware"
	"rlp-member-service/config"
	"rlp-member-service/model"
	"rlp-member-service/system"

	"github.com/gin-gonic/gin"
)

// Option type and global slice for router modifications.
type Option func(*gin.RouterGroup)

var options = []Option{}

var endpointList []map[string]string

func Include(opts ...Option) {
	options = append(options, opts...)
}

func Init() *gin.Engine {
	// Include additional routers
	Include(general.Routers)

	db := system.GetDb()
	if err := model.MigrateAuditLog(db); err != nil {
		log.Fatalf("audit log migration: %v", err)
	}
	if err := model.MigrateUser(db); err != nil {
		log.Fatalf("user migration: %v", err)
	}
	if err := model.MigrateUserSession(db); err != nil {
		log.Fatalf("user session migration: %v", err)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.AuditLogger(db))

	apiGroup := r.Group("/api")
	for _, opt := range options {
		opt(apiGroup)
	}

	// Capture routes but exclude the HandlerFunc to avoid JSON marshalling errors.
	routes := r.Routes()
	for _, route := range routes {
		endpointList = append(endpointList, map[string]string{
			"method": route.Method,
			"path":   route.Path,
		})
	}

	r.Run(fmt.Sprintf(":%d", config.GetConfig().Http.Port))
	return r
}

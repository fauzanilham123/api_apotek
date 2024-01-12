package routes

import (
	"api_apotek/controllers"
	"api_apotek/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Convert c.Handler to a gin.HandlerFunc
	corsMiddleware := func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}

	// Use the converted middleware in gin
	r.Use(corsMiddleware)

	// Buat objek limiter
	var limiter = rate.NewLimiter(rate.Limit(10), 1) // Contoh: 10 permintaan per detik

	// Tambahkan middleware rate limiting ke router utama
	r.Use(func(c *gin.Context) {
		// Gunakan limiter untuk memeriksa rate limiting
		if limiter.Allow() == false {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	})

	// set db to gin context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	r.POST("/register", controllers.Register)
    r.POST("/login", controllers.Login)
	r.GET("/logactivity/", controllers.GetAllLogActivity)
	r.GET("/user/", controllers.GetAllUser)
	r.GET("/obat/", controllers.GetAllObat)
	r.GET("/obat/:id", controllers.GetObatById)

	Obat := r.Group("/obat")
    Obat.Use(middlewares.JwtAuthMiddleware()) //use jwt
    Obat.POST("/", controllers.CreateObat)
    Obat.PATCH("/:id", controllers.UpdateObat)
    Obat.DELETE("/:id", controllers.DeleteObat)

	return r
}
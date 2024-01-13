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

	// Menangani rute yang tidak ditemukan
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})

	r.POST("/register", controllers.Register)
    r.POST("/login", controllers.Login)
	r.GET("/logactivity/", controllers.GetAllLogActivity)
	r.GET("/user/", controllers.GetAllUser)
	r.GET("/user/:id", controllers.GetUserById)
	r.GET("/obat/", controllers.GetAllObat)
	r.GET("/obat/:id", controllers.GetObatById)
	r.GET("/recipe/", controllers.GetAllRecipe)
	r.GET("/recipe/:id", controllers.GetRecipeById)
	r.GET("/transaction/", controllers.GetAllTransaction)
	r.GET("/transaction/:id", controllers.GetTransactionById)

	obat := r.Group("/obat")
    obat.Use(middlewares.JwtAuthMiddleware()) //use jwt
    obat.POST("/", controllers.CreateObat)
    obat.PATCH("/:id", controllers.UpdateObat)
    obat.DELETE("/:id", controllers.DeleteObat)
	
	recipe := r.Group("/recipe")
    recipe.Use(middlewares.JwtAuthMiddleware()) //use jwt
    recipe.POST("/", controllers.CreateRecipe)
    recipe.PATCH("/:id", controllers.UpdateRecipe)
    recipe.DELETE("/:id", controllers.DeleteRecipe)
	
	user := r.Group("/user")
    user.Use(middlewares.JwtAuthMiddleware()) //use jwt
    user.POST("/", controllers.CreateUser)
    user.PATCH("/:id", controllers.UpdateUser)
    user.DELETE("/:id", controllers.DeleteUser)
	
	transaction := r.Group("/transaction")
    transaction.Use(middlewares.JwtAuthMiddleware()) //use jwt
    transaction.POST("/", controllers.CreateTransaction)
    transaction.PATCH("/:id", controllers.UpdateTransaction)
    transaction.DELETE("/:id", controllers.DeleteTransaction)

	return r
}
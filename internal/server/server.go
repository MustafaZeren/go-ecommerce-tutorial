package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mustafazeren/go-ecommerce-course/internal/config"
	"github.com/mustafazeren/go-ecommerce-course/internal/services"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	config         *config.Config
	db             *gorm.DB
	logger         *zerolog.Logger
	authService    *services.AuthService
	productService *services.ProductService
	userService    *services.UserService
}

func New(
	c *config.Config,
	db *gorm.DB,
	logger *zerolog.Logger,
	authService *services.AuthService,
	productService *services.ProductService,
	userService *services.UserService,
) *Server {
	return &Server{
		config:         c,
		db:             db,
		logger:         logger,
		authService:    authService,
		productService: productService,
		userService:    userService,
	}
}
func (s *Server) SetupRoutes() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(s.corsMiddleware())

	r.GET("/health", s.healthCheck)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		auth.POST("/register", s.register)
		auth.POST("/login", s.login)
		auth.POST("/refresh", s.refreshToken)
		auth.POST("/logout", s.logout)

		protected := api.Group("/")
		protected.Use(s.authMiddleware())
		{
			users := protected.Group("/users")
			{
				userRoute := users
				userRoute.GET("/profile", s.getProfile)
				userRoute.PUT("/profile", s.updateProfile)
			}
			categories := protected.Group("/categories")
			{
				categoryRoute := categories
				categoryRoute.POST("/", s.adminMiddleware(), s.createCategory)
				categoryRoute.PUT("/:id", s.adminMiddleware(), s.updateCategory)
				categoryRoute.DELETE("/:id", s.adminMiddleware(), s.deleteCategory)
			}
			products := protected.Group("/products")
			{
				productRoute := products
				productRoute.POST("/", s.adminMiddleware(), s.createProduct)
				productRoute.PUT("/:id", s.adminMiddleware(), s.updateProduct)
				productRoute.DELETE("/:id", s.adminMiddleware(), s.deleteProduct)
			}
		}

		// public routes
		api.GET("/categories", s.getCategories)
		api.GET("/products", s.getProducts)
		api.GET("/products/:id", s.getProduct)

		admin := api.Group("/admin")
		// It is written for using adminMiddleware function
		admin.Use(s.adminMiddleware())
	}

	return r
}
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

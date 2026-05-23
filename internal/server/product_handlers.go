package server

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mustafazeren/go-ecommerce-course/internal/dto"
	"github.com/mustafazeren/go-ecommerce-course/internal/utils"
	"gorm.io/gorm"
)

func (s *Server) createCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	category, err := s.productService.CreateCategory(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create category", err)
		return
	}
	utils.CreatedResponse(c, "Category created successfully", category)
}
func (s *Server) getCategories(c *gin.Context) {
	categories, err := s.productService.GetCategories()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get categories", err)
	}
	utils.SuccessResponse(c, "Category retrieved successfully", categories)
}
func (s *Server) updateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid id", err)
		return
	}
	var req dto.UpdateCategoryRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	category, err := s.productService.UpdateCategory(id, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update category", err)
		return
	}
	utils.SuccessResponse(c, "Category updated successfully", category)
}
func (s *Server) deleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid id", err)
		return
	}
	err = s.productService.DeleteCategory(id)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete category", err)
		return
	}
	utils.SuccessResponse(c, "Category deleted successfully", nil)
}
func (s *Server) createProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	product, err := s.productService.CreateProduct(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create product", err)
		return
	}
	utils.CreatedResponse(c, "Product created successfully", product)
}
func (s *Server) getProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	products, meta, err := s.productService.GetProducts(page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get products", err)
		return
	}
	utils.PaginatedSuccessResponse(c, "Products received successfully", products, *meta)
}
func (s *Server) getProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid id", err)
		return
	}
	product, err := s.productService.GetProduct(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundResponse(c, "Product not found")
		}
		utils.InternalServerErrorResponse(c, "Failed to get product", err)
		return
	}
	utils.SuccessResponse(c, "Product retrieved successfully", product)
}
func (s *Server) updateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid id", err)
		return
	}
	var req dto.UpdateProductRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	product, err := s.productService.UpdateProduct(uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update product", err)
		return
	}
	utils.SuccessResponse(c, "Product updated successfully", product)
}
func (s *Server) deleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid id", err)
		return
	}
	err = s.productService.DeleteProduct(id)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete product", err)
	}
	utils.SuccessResponse(c, "Product deleted successfully", nil)
}

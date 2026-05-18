package services

import (
	"github.com/mustafazeren/go-ecommerce-course/internal/dto"
	"github.com/mustafazeren/go-ecommerce-course/internal/models"
	"github.com/mustafazeren/go-ecommerce-course/internal/utils"
	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{
		db: db,
	}
}

func (s *ProductService) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.db.Create(&category).Error; err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
	}, nil
}
func (s *ProductService) GetCategories() ([]*dto.CategoryResponse, error) {
	var categories []*models.Category
	if err := s.db.Where("is_active = ?", true).Find(&categories).Error; err != nil {
		return nil, err
	}

	result := make([]*dto.CategoryResponse, len(categories))
	for i, category := range categories {
		result[i] = &dto.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			IsActive:    category.IsActive,
		}
	}

	return result, nil
}
func (s *ProductService) UpdateCategory(id uint64, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	var category models.Category
	if err := s.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	category.Name = req.Name
	category.Description = req.Description
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	if err := s.db.Save(&category).Error; err != nil {
		return nil, err
	}
	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
	}, nil
}
func (s *ProductService) DeleteCategory(id uint64) error {
	return s.db.Delete(&models.Category{}, id).Error
}
func (s *ProductService) CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.SKU,
	}
	if err := s.db.Create(&product).Error; err != nil {
		return nil, err
	}
	return s.GetProduct(product.ID)
}
func (s *ProductService) GetProducts(page, limit int) ([]*dto.ProductResponse, *utils.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 5
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit
	var products []*models.Product
	var total int64

	if err := s.db.Model(&models.Product{}).
		Where("is_active = ?", true).
		Count(&total).Error; err != nil {
		return nil, nil, err
	}
	if err := s.db.Preload("Category").Preload("Images").
		Where("is_active = ?", true).
		Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, nil, err
	}

	response := make([]*dto.ProductResponse, len(products))
	for i := range products {
		response[i] = s.convertToProductResponse(products[i])
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
	return response, meta, nil
}
func (s *ProductService) GetProduct(id uint) (*dto.ProductResponse, error) {
	var product models.Product
	if err := s.db.Preload("Category").Preload("Images").Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}
	return s.convertToProductResponse(&product), nil
}

func (s *ProductService) UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	var product models.Product
	if err := s.db.First(&product, id).Error; err != nil {
		return nil, err
	}

	product.CategoryID = req.CategoryID
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := s.db.Save(&product).Error; err != nil {
		return nil, err
	}

	return s.GetProduct(id)
}
func (s *ProductService) convertToProductResponse(product *models.Product) *dto.ProductResponse {
	images := make([]dto.ProductImageResponse, len(product.Images))
	for i := range product.Images {
		images[i] = dto.ProductImageResponse{
			ID:        product.Images[i].ID,
			URL:       product.Images[i].URL,
			AltText:   product.Images[i].AltText,
			IsPrimary: product.Images[i].IsPrimary,
		}
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		IsActive:    product.IsActive,
		Images:      images,
		Price:       product.Price,
		Stock:       product.Stock,
		SKU:         product.SKU,
		CategoryID:  product.CategoryID,
		Category: dto.CategoryResponse{
			ID:          product.Category.ID,
			Name:        product.Category.Name,
			Description: product.Category.Description,
			IsActive:    product.Category.IsActive,
		},
	}
}

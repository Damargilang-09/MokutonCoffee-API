package controllers

import (
	"MokutonCoffee-API/config"
	"MokutonCoffee-API/dto/common"
	"MokutonCoffee-API/dto/request"
	"MokutonCoffee-API/dto/response"
	"MokutonCoffee-API/models"
	"MokutonCoffee-API/utils"
	"context"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/imagekit-developer/imagekit-go/v2"
)

func CreateProduct(c *gin.Context) {

	var input request.CreateProductRequest
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request input",
			"error":   err.Error(),
		})
		return
	}

	//Upload Image

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Image file is required",
		})
		return
	}
	defer file.Close()

	ik := utils.InitImageKit()
	uploadRes, errUpload := ik.Files.Upload(
		context.Background(),
		imagekit.FileUploadParams{
			File:     file,
			FileName: header.Filename,
		})

	if errUpload != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload image",
		})
		return
	}

	product := models.Product {
		SKU:         input.SKU,
		Name:        input.Name,
		Image:       uploadRes.URL,
		ImageId:     uploadRes.FileID,
		Description: input.Description,
		Stock:       *input.Stock,
		Price:       *input.Price,
	}

	productCreated := config.DB.Create(&product).Error

	if productCreated != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create product",
		})
		return
	}

	result := response.ProductResponse{
		ID:          uuid.UUID(product.ID),
		SKU:         product.SKU,
		Name:        product.Name,
		Image:       product.Image,
		Description: product.Description,
		Stock:       product.Stock,
		Price:       product.Price,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "New product created",
		"data":    result,
	})
}

func GetAllProducts(c *gin.Context) {
	var query common.PaginationQuery
	c.ShouldBindQuery(&query)

	var products []models.Product

	var productsData = config.DB.Scopes(query.Paginate()).Find(&products).Error
	if productsData != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve products",
		})
		return
	}

	var total int64
	if err := config.DB.Model(&models.Product{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to count products",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Products retrieved successfully",
		"page":    query.Page,
		"limit":   query.Limit,
		"total":   total,
		"data":    products,
	})
}

func GetProductById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID format",
		})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail product retrieved successfully",
		"data":    product,
	})
}

func UpdateProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID format",
		})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	var input request.UpdateProductRequest

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()

		//upload new image
		ik := utils.InitImageKit()

		fileName := header.Filename
		uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
			File:     file,
			FileName: fileName,
		})

		if errUpload == nil {
			// hapus gambar lama sebelum mengganti
			if product.ImageId != "" {
				if errDel := ik.Files.Delete(context.Background(), product.ImageId); errDel != nil {
					log.Printf("ImageKit deletion failed: %v", errDel)
				}
			}

			//upload image
			product.Image = uploadRes.URL
			product.ImageId = uploadRes.FileID
		}
	}

	if sku := input.SKU; sku != "" {
		product.SKU = sku
	}

	if name := input.Name; name != "" {
		product.Name = name
	}

	if description := input.Description; description != "" {
		product.Description = description
	}

	if stock := input.Stock; stock != nil {
		product.Stock = *stock
	}

	if price := input.Price; price != nil {
		product.Price = *price
	}

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update product",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"data":    product,
	})
}

func DeleteProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID format",
		})
		return
	}

	var product models.Product

	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	// Delete image from ImageKit
	if product.ImageId != "" {
		ik := utils.InitImageKit()

		err := ik.Files.Delete(context.Background(), product.ImageId)
		if err != nil {
			log.Printf("ImageKit deletion failed: %v", err)
		}
	}

	// Delete product from database
	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}

package controllers

import (
	"MokutonCoffee-API/config"
	"MokutonCoffee-API/models"
	"MokutonCoffee-API/utils"
	"context"
	"log"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/imagekit-developer/imagekit-go/v2"
)

func CreateProduct(c *gin.Context) {
	// 1. Ambil dan validasi input teks terlebih dahulu sebelum upload file
	sku := c.PostForm("sku")
	name := c.PostForm("name")
	description := c.PostForm("description")

	if sku == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "SKU and Name are required",
		})
		return
	}

	stock, err := strconv.Atoi(c.PostForm("stock"))
	if err != nil || stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Stock must be a valid positive number",
		})
		return
	}

	price, err := strconv.ParseInt(c.PostForm("price"), 10, 64)
	if err != nil || price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Price must be a valid number greater than 0",
		})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Image file is required",
		})
		return
	}
	defer file.Close()

	ik := utils.InitImageKit()
	uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
		File:     file,
		FileName: header.Filename,
	})

	if errUpload != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload image",
		})
		return
	}

	product := models.Product{
		SKU:         sku,
		Name:        name,
		Image:       uploadRes.URL,
		ImageId:     uploadRes.FileID,
		Description: description,
		Stock:       stock,
		Price:       price,
	}

	productCreated := config.DB.Create(&product).Error

	if productCreated != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "New product created",
		"data":    product,
	})
}

func GetAllProducts(c *gin.Context) {
	var query models.PaginationQuery
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
			//delete image
			if product.ImageId == "" {
				ik.Files.Delete(context.Background(), product.ImageId)
			}

			//upload image
			product.Image = uploadRes.URL
			product.ImageId = uploadRes.FileID
		}
	}

	if sku := c.PostForm("sku"); sku != "" {
		product.SKU = sku
	}

	if name := c.PostForm("name"); name != "" {
		product.Name = name
	}

	if description := c.PostForm("description"); description != "" {
		product.Description = description
	}

	stockStr := c.PostForm("stock")
	if stockStr != "" {
		stock, err := strconv.Atoi(stockStr)
		if err != nil || stock <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Stock must be a valid non-negative integer",
			})
			return
		}
		product.Stock = stock
	}

	priceStr := c.PostForm("price")
	if priceStr != "" {
		price, err := strconv.Atoi(priceStr)
		if err != nil || price < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Stock must be a valid non-negative integer",
			})
			return
		}
		product.Price = int64(price)
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

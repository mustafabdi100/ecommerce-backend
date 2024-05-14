package routes

import (
	"database/sql"
	"ecommerce-backend/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/api/categories", getCategories)
	router.POST("/categories", createCategory)
	router.GET("/api/categories/:id", getCategoryByID)
	router.PUT("/api/categories/:id", updateCategory)
	router.DELETE("/api/categories/:id", deleteCategory)
	router.GET("/api/products/:id", getProductByID)
	router.GET("/api/products", getProductsByCategory)
	router.POST("/api/products", createProduct)
	router.PUT("/api/products/:id", updateProduct)
	router.DELETE("/api/products/:id", deleteProduct)
}

func getCategories(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, name FROM categories")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var categories []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	for rows.Next() {
		var category struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		categories = append(categories, category)
	}

	c.JSON(200, categories)
}

// createCategory creates a new category
func createCategory(c *gin.Context) {
	var category struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use db.QueryRow instead of db.Exec
	var id int
	err := db.DB.QueryRow("INSERT INTO categories (name) VALUES ($1) RETURNING id", category.Name).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":   id,
		"name": category.Name,
	})
}

func getCategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var category struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	err = db.DB.QueryRow("SELECT id, name FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Name)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, category)
}

func updateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var category struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = db.DB.Exec("UPDATE categories SET name = $1 WHERE id = $2", category.Name, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

func deleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	_, err = db.DB.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

func getProducts(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, name, description, price, image_url, category_id FROM products")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var products []struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		ImageURL    string  `json:"image_url"`
		CategoryID  int     `json:"category_id"`
	}

	for rows.Next() {
		var product struct {
			ID          int     `json:"id"`
			Name        string  `json:"name"`
			Description string  `json:"description"`
			Price       float64 `json:"price"`
			ImageURL    string  `json:"image_url"`
			CategoryID  int     `json:"category_id"`
		}
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.ImageURL, &product.CategoryID); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		products = append(products, product)
	}

	c.JSON(200, products)
}

func getProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		ImageURL    string  `json:"image_url"`
		CategoryID  int     `json:"category_id"`
	}

	err = db.DB.QueryRow("SELECT id, name, description, price, image_url, category_id FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.ImageURL, &product.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func getProductsByCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Query("category"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	rows, err := db.DB.Query("SELECT id, name, description, price, image_url FROM products WHERE category_id = $1", categoryID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var products []struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		ImageURL    string  `json:"image_url"`
	}

	for rows.Next() {
		var product struct {
			ID          int     `json:"id"`
			Name        string  `json:"name"`
			Description string  `json:"description"`
			Price       float64 `json:"price"`
			ImageURL    string  `json:"image_url"`
		}
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.ImageURL); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		products = append(products, product)
	}

	c.JSON(200, products)
}

// createProduct creates a new product
func createProduct(c *gin.Context) {
	var product struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Price       float64 `json:"price" binding:"required"`
		ImageURL    string  `json:"image_url" binding:"required"`
		CategoryID  int     `json:"category_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use db.QueryRow instead of db.Exec
	var id int
	err := db.DB.QueryRow("INSERT INTO products (name, description, price, image_url, category_id) VALUES ($1, $2, $3, $4, $5) RETURNING id", product.Name, product.Description, product.Price, product.ImageURL, product.CategoryID).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          id,
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"image_url":   product.ImageURL,
		"category_id": product.CategoryID,
	})
}

func updateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		ImageURL    string  `json:"image_url"`
		CategoryID  int     `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = db.DB.Exec("UPDATE products SET name = $1, description = $2, price = $3, image_url = $4, category_id = $5 WHERE id = $6",
		product.Name, product.Description, product.Price, product.ImageURL, product.CategoryID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func deleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	_, err = db.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

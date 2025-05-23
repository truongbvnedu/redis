package controllers

import (
	"encoding/json"
	"go-mvc-demo/config"
	"go-mvc-demo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

func (ctl *PostController) CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&post)
	c.JSON(http.StatusOK, post)
}

func (ctl *PostController) GetPosts(c *gin.Context) {
	cacheKey := "all_posts"
	var posts []models.Post

	val, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		c.Data(http.StatusOK, "application/json", []byte(val))
		return
	}

	config.DB.Preload("Comments").Find(&posts)

	jsonData, err := json.Marshal(posts)
	if err == nil {
		config.RedisClient.Set(config.Ctx, cacheKey, jsonData, 0)
	}

	c.JSON(http.StatusOK, posts)
}
func (ctl *PostController) UpdatePost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var updatedData models.Post
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.Title = updatedData.Title
	post.Content = updatedData.Content

	if err := config.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully", "post": post})
}
func (ctl *PostController) DeletePost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err := config.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
func (ctl *PostController) GetOne(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := config.DB.Preload("Comments").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

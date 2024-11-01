/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
*/
package web

import (
	"strconv"

	"github.com/AnyoneClown/anydb/config"
	"github.com/AnyoneClown/anydb/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct{}

type ConfigInput struct {
	ConfigName string `json:"configName" binding:"required"`
	Driver     string `json:"driver" binding:"required,oneof=postgres cockroachdb"`
	Host       string `json:"host" binding:"required"`
	Port       string `json:"port" binding:"required,port"`
	User       string `json:"user" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Database   string `json:"database" binding:"required"`
}

// Custom validator for port
func portValidator(fl validator.FieldLevel) bool {
	port := fl.Field().String()
	if p, err := strconv.Atoi(port); err == nil {
		return p > 0 && p <= 65535
	}
	return false
}

// GET /api/configs
func (h *Handler) GetConfigs(c *gin.Context) {
	configs, err := utils.LoadConfigs(config.ConfigFile)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to load existing configurations"})
		return
	}

	c.JSON(200, gin.H{"configs": configs})
}

// GET /api/configs/:id
func (h *Handler) GetConfig(c *gin.Context) {
	configID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Log.Error("Invalid UUID format", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid UUID format"})
		return
	}

	config, err := utils.GetConfigByID(configID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Configuration not found"})
		return
	}

	c.JSON(200, gin.H{"config": config})
}

// POST /api/configs
func (h *Handler) CreateConfig(c *gin.Context) {
	var input ConfigInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Error("Invalid input", zap.Error(err))
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make(map[string]string)
		for _, fieldError := range validationErrors {
			errorMessages[fieldError.Field()] = fieldError.Error()
		}
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	newConfig := config.DBConfig{
		ID:         uuid.New(),
		ConfigName: input.ConfigName,
		Driver:     input.Driver,
		Host:       input.Host,
		Port:       input.Port,
		User:       input.User,
		Password:   input.Password,
		Database:   input.Database,
	}

	// Load existing configurations
	configs, err := utils.LoadConfigs(config.ConfigFile)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to load existing configurations"})
		return
	}

	// Append the new configuration
	configs = append(configs, newConfig)

	// Save the updated configurations back to the file
	if err := utils.SaveConfigs(configs, config.ConfigFile); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save new configuration"})
		return
	}

	c.JSON(201, gin.H{"message": "Configuration created successfully", "config": newConfig})
}

// DELETE /api/configs/:id
func (h *Handler) DeleteConfig(c *gin.Context) {
	configID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Log.Error("Invalid UUID format", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid UUID format"})
		return
	}

	configToDelete, err := utils.GetConfigByID(configID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Configuration not found"})
		return
	}

	configs, err := utils.LoadConfigs(config.ConfigFile)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to load existing configurations"})
		return
	}

	// Remove the configuration from the slice
	for i, cfg := range configs {
		if cfg.ID == configToDelete.ID {
			configs = append(configs[:i], configs[i+1:]...)
			break
		}
	}

	// Save the updated configurations back to the file
	if err := utils.SaveConfigs(configs, config.ConfigFile); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save new configuration"})
		return
	}

	c.JSON(200, gin.H{"message": "Configuration deleted successfully"})
}
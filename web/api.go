package web

import (
	"github.com/AnyoneClown/anydb/config"
	"github.com/AnyoneClown/anydb/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct{}

// GET /api/configs
func (h *Handler) GetConfigs(c *gin.Context) {
	configs, err := utils.LoadConfigs(config.ConfigFile)
	if err != nil {
		utils.Log.Error("Failed to load configs from config file", zap.Error(err))
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
	var input config.ConfigInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Error("Invalid input", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// Validate inputs
	if err := utils.ValidateConfig(input); err != nil {
		utils.Log.Error("Invalid input", zap.Error(err))
		c.JSON(400, gin.H{"error": "Input didn't pass validation"})
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
		utils.Log.Error("Failed to load configs from config file", zap.Error(err))
		c.JSON(500, gin.H{"error": "Failed to load existing configurations"})
		return
	}

	// Append the new configuration
	configs = append(configs, newConfig)

	// Save the updated configurations back to the file
	if err := utils.SaveConfigs(configs, config.ConfigFile); err != nil {
		utils.Log.Error("Failed to save new configuration", zap.Error(err))
		c.JSON(500, gin.H{"error": "Failed to save new configuration"})
		return
	}

	c.JSON(201, gin.H{"message": "Configuration created successfully", "config": newConfig})
}

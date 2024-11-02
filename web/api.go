package web

import (
	"net/http"
	"os"
	"strconv"

	"github.com/AnyoneClown/anydb/config"
	"github.com/AnyoneClown/anydb/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// Handler struct to group all handler methods
type Handler struct{}

// ConfigInput struct for binding JSON input
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

// ErrorResponse struct for consistent error responses
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse struct for consistent success responses
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Helper function to handle errors
func handleError(c *gin.Context, status int, err error, message string) {
	utils.Log.Error(message, zap.Error(err))
	c.JSON(status, ErrorResponse{Error: message})
}

// GET /api/configs
func (h *Handler) GetConfigs(c *gin.Context) {
	configs, err := utils.LoadConfigs(config.ConfigFile)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "Failed to load existing configurations")
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Message: "Configurations retrieved successfully", Data: configs})
}

// GET /api/configs/:id
func (h *Handler) GetConfig(c *gin.Context) {
	configID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, err, "Invalid UUID format")
		return
	}

	config, err := utils.GetConfigByID(configID)
	if err != nil {
		handleError(c, http.StatusNotFound, err, "Configuration not found")
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Configuration retrieved successfully", Data: config})
}

// POST /api/configs
func (h *Handler) CreateConfig(c *gin.Context) {
	var input ConfigInput
	if err := c.ShouldBindJSON(&input); err != nil {
		handleError(c, http.StatusBadRequest, err, "Invalid input")
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

	configs, err := utils.LoadConfigs(config.ConfigFile)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "Failed to load existing configurations")
		return
	}

	configs = append(configs, newConfig)
	if err := utils.SaveConfigs(configs, config.ConfigFile); err != nil {
		handleError(c, http.StatusInternalServerError, err, "Failed to save new configuration")
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{Message: "Configuration created successfully", Data: newConfig})
}

// PUT /api/configs/:id
func (h *Handler) UpdateConfig(c *gin.Context) {
	var input ConfigInput
	if err := c.ShouldBindJSON(&input); err != nil {
		handleError(c, http.StatusBadRequest, err, "Invalid input")
		return
	}

	configID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, err, "Invalid UUID format")
		return
	}

	configToUpdate, err := utils.GetConfigByID(configID)
	if err != nil {
		handleError(c, http.StatusNotFound, err, "Configuration not found")
		return
	}

	configs, err := utils.LoadConfigs(config.ConfigFile)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "Failed to load existing configurations")
		return
	}

	for i, cfg := range configs {
		if cfg.ID == configToUpdate.ID {
			configs[i] = config.DBConfig{
				ID:         configID,
				ConfigName: input.ConfigName,
				Driver:     input.Driver,
				Host:       input.Host,
				Port:       input.Port,
				User:       input.User,
				Password:   input.Password,
				Database:   input.Database,
			}
			break
		}
	}

	if err := utils.SaveConfigs(configs, config.ConfigFile); err != nil {
		handleError(c, http.StatusInternalServerError, err, "Failed to save updated configuration")
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Configuration updated successfully"})
}

// DELETE /api/configs/:id
func (h *Handler) DeleteConfig(c *gin.Context) {
	configID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, err, "Invalid UUID format")
		return
	}

	configToDelete, err := utils.GetConfigByID(configID)
	if err != nil {
		handleError(c, http.StatusNotFound, err, "Configuration not found")
		return
	}

	configs, err := utils.LoadConfigs(config.ConfigFile)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "Failed to load existing configurations")
		return
	}

	for i, cfg := range configs {
		if cfg.ID == configToDelete.ID {
			configs = append(configs[:i], configs[i+1:]...)
			break
		}
	}

	if err := utils.SaveConfigs(configs, config.ConfigFile); err != nil {
		handleError(c, http.StatusInternalServerError, err, "Failed to save updated configurations")
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Configuration deleted successfully"})
}

// POST /api/configs/select/:id
func (h *Handler) SelectConfig(c *gin.Context) {
	configID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, err, "Invalid UUID format")
		return
	}

	configs, err := utils.LoadConfigs(config.ConfigFile)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "Failed to load existing configurations")
		return
	}

	var selectedConfig config.DBConfig
	var configExists bool
	for _, cfg := range configs {
		if cfg.ID == configID {
			selectedConfig = cfg
			configExists = true
			break
		}
	}

	if !configExists {
		handleError(c, http.StatusNotFound, err, "Configuration not found")
		return
	}

	data, err := yaml.Marshal(selectedConfig)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err, "Failed to marshal selected config")
		return
	}

	if err := os.WriteFile(config.DefaultConfigFile, data, 0644); err != nil {
		handleError(c, http.StatusInternalServerError, err, "Failed to write default config file")
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Configuration selected successfully", Data: selectedConfig})
}

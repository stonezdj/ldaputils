package handlers

import (
	"fmt"
	"github.com/goharbor/ldaputils/dao"
	"github.com/goharbor/ldaputils/dao/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type H map[string]interface{}

// GetConfigs endpoint
func GetConfigs(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		configs := dao.DAO.Query(db)
		result := &models.LdapConfigCollection{
			Items: configs,
		}
		return c.JSON(http.StatusOK, result)
	}
}

// PutConfigs endpoint
func PutConfig(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := &models.LdapConfig{}
		c.Bind(cfg)
		dao.DAO.Add(db, cfg)
		return c.JSON(http.StatusCreated, H{"created": cfg.ID})
	}
}

// DeleteTask endpoint
func DeleteConfig(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		dao.DAO.Delete(db, id)
		return c.JSON(http.StatusOK, H{
			"deleted": id})
	}
}

func TestingConfig(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.QueryParam("id"))
		//get configurations and run the ldap ping testing, send out the message
		fmt.Printf("the testing id is %v\n", id)
		return c.JSON(http.StatusOK, models.LDAPTestResult{Success: true, Message: "Test passed"})
	}
}

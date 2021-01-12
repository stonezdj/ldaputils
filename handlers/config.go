package handlers

import (
	"fmt"
	"github.com/goharbor/harbor/src/common/utils/ldap"
	"github.com/goharbor/ldaputils/dao"
	"github.com/goharbor/ldaputils/dao/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type H map[string]interface{}
type Test struct {
	Type string `json:"type"`
}

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
		result := &models.LDAPTestResult{}
		id, _ := strconv.Atoi(c.Param("id"))
		test := &Test{}
		err := c.Bind(test)
		if err != nil {
			return c.JSON(http.StatusOK, result.FailWithError(err))
		}
		//get configurations and run the ldap ping testing, send out the message
		cfg := dao.DAO.Get(db, id)
		if cfg == nil {
			return c.JSON(http.StatusOK, models.LDAPTestResult{Success: false,
				Message: []string{fmt.Sprintf("Configure not found: %v", id)}})
		}

		result = Ping(cfg)
		fmt.Printf("the testing id is %v\n", id)
		return c.JSON(http.StatusOK, result)
	}
}

func Ping(ldapCfg *models.LdapConfig) *models.LDAPTestResult {
	result := &models.LDAPTestResult{}
	fmt.Printf("Start to ping LDAP server: %v\n", ldapCfg.LdapURL)
	err := ldap.ConnectionTestWithAllConfig(ldapCfg.LdapConf, ldapCfg.LdapGroupConf)
	if err != nil {
		return result.Fail().WithMsg(fmt.Sprintf("Error at connection test, %+v", err))
	}
	return result.Suc().WithMsg("Ping test success")
}

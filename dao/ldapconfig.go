package dao

import (
	"github.com/goharbor/ldaputils/dao/models"
	"gorm.io/gorm"
)

var DAO = &LdapConfigDAO{}

type LdapConfigDAO struct {
}

func (l *LdapConfigDAO) Query(db *gorm.DB) []models.LdapConfig {
	result := make([]models.LdapConfig, 0)
	db.Limit(20).Find(&result)
	return result
}

func (l *LdapConfigDAO) Get(db *gorm.DB, id int) *models.LdapConfig {
	cfg := &models.LdapConfig{}
	db.First(cfg, id)
	if cfg.ID > 0 {
		return cfg
	}
	return nil // not found
}

func (l *LdapConfigDAO) Add(db *gorm.DB, config *models.LdapConfig) {
	db.Create(config)
	return
}

func (l *LdapConfigDAO) Delete(db *gorm.DB, id int) {
	db.Delete(&models.LdapConfig{}, id)
	return
}

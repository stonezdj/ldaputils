package dao

import (
	"fmt"
	"github.com/goharbor/ldaputils/dao/models"
	"github.com/golangplus/testing/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"testing"
)

var database *gorm.DB

func TestMain(m *testing.M) {
	if database == nil {
		db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			fmt.Println("failed to open database")
			return
		}
		database = db
		database.AutoMigrate(&models.LdapConfig{})
		os.Exit(m.Run())
	}
}

func TestAddAndDelete(t *testing.T) {
	config := &models.LdapConfig{
		LdapURL:                      "10.193.28.58",
		LdapSearchDn:                 "cn=admin,dc=example,dc=com",
		LdapSearchPassword:           "admin",
		LdapBaseDn:                   "dc=example,dc=com",
		LdapFilter:                   "",
		LdapUID:                      "cn",
		LdapScope:                    2,
		LdapConnectionTimeout:        30,
		LdapVerifyCert:               false,
		LdapGroupBaseDN:              "dc=example,dc=com",
		LdapGroupFilter:              "class=groupOfNames",
		LdapGroupNameAttribute:       "cn",
		LdapGroupSearchScope:         2,
		LdapGroupAdminDN:             "cn=harbor_users,dc=example,dc=com",
		LdapGroupMembershipAttribute: "memberof",
	}
	DAO.Add(database, config)
	results := DAO.Query(database)
	assert.True(t, "must greater than zero", len(results) > 0)
	for _, v := range results {
		DAO.Delete(database, int(v.ID))
	}
	results2 := DAO.Query(database)
	assert.Equal(t, "all should be cleaned", len(results2), 0)
}

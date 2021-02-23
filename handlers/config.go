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
	Type  string `json:"type"`
	User  string `json:"user"`
	Group string `json:"group"`
}

var TestingHandlerMap = map[string]TestingHandler{
	"ping":              Ping,
	"search_user":       SearchUser,
	"test_group_config": TestGroupConfig,
	"test_group_admin":  TestGroupAdminConfig,
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
		h, ok := TestingHandlerMap[test.Type]
		if !ok {
			return c.JSON(http.StatusOK, result.FailWithError(fmt.Errorf("The testing command %v is not found", test.Type)))
		}
		result = h(cfg, test)
		fmt.Printf("the testing id is %v\n", id)
		return c.JSON(http.StatusOK, result)
	}
}

type TestingHandler func(ldapCfg *models.LdapConfig, test *Test) *models.LDAPTestResult

func Ping(ldapCfg *models.LdapConfig, test *Test) *models.LDAPTestResult {
	result := &models.LDAPTestResult{}
	fmt.Printf("Start to ping LDAP server: %v\n", ldapCfg.LdapURL)
	err := ldap.ConnectionTestWithAllConfig(ldapCfg.LdapConf, ldapCfg.LdapGroupConf)
	if err != nil {
		return result.Fail().WithMsg(fmt.Sprintf("Error at connection test, %+v", err))
	}
	return result.Suc().WithMsg("Ping test success")
}

func TestGroupConfig(ldapCfg *models.LdapConfig, test *Test) *models.LDAPTestResult {
	ret := &models.LDAPTestResult{}
	fmt.Printf("Start to test LDAP group config\n")
	session, err := ldap.CreateWithAllConfig(ldapCfg.LdapConf, ldapCfg.LdapGroupConf)
	if err != nil {
		return ret.FailWithError(err)
	}
	session.Open()
	defer session.Close()

	fmt.Println("Verify LDAP group configurations")

	if len(ldapCfg.LdapGroupConf.LdapGroupBaseDN) == 0 {
		ret.Fail().WithMsg("LDAP group DN is not configured")
	}

	fmt.Println("Trying to search group in current search conditions.")
	if len(test.Group) == 0 {
		return ret.Fail().WithMsg("Need to provide LDAP group name")
	}
	groups, err := session.SearchGroupByName(test.Group)
	if err != nil {
		return ret.FailWithError(err)
	}
	if len(groups) == 0 {
		return ret.Fail().WithMsg("No LDAP group found!")
	} else {
		return ret.Suc().WithMsg(fmt.Sprintf("Found %v groups in current condition.", len(groups)))
	}

}

func TestGroupAdminConfig(ldapCfg *models.LdapConfig, test *Test) *models.LDAPTestResult {
	ret := &models.LDAPTestResult{}
	fmt.Printf("Start to test LDAP group config\n")
	session, err := ldap.CreateWithAllConfig(ldapCfg.LdapConf, ldapCfg.LdapGroupConf)
	if err != nil {
		return ret.FailWithError(err)
	}
	session.Open()
	defer session.Close()

	if len(ldapCfg.LdapGroupConf.LdapGroupAdminDN) == 0 {
		return ret.Fail().WithMsg("The LDAP group admin DN is not configured.")
	}
	ret.WithMsg(fmt.Sprintf("Trying to search the group with admin privileges, ldap group admin dn: %v\n", ldapCfg.LdapGroupConf.LdapGroupAdminDN))

	groups, err := session.SearchGroupByDN(ldapCfg.LdapGroupConf.LdapGroupAdminDN)
	if err != nil {
		return ret.FailWithError(err)
	}
	ret.WithMsg(fmt.Sprintf("Found %v groups with admin privileges.", len(groups)))

	if len(groups) > 0 {
		fmt.Printf("Trying to find users in this group: %v\n", ldapCfg.LdapGroupConf.LdapGroupAdminDN)
		count := 0
		userList, err := session.SearchUser("")
		if err != nil {
			ret.FailWithError(err)
		}
		for _, user := range userList {
			if stringInSlice(ldapCfg.LdapGroupConf.LdapGroupAdminDN, user.GroupDNList) {
				ret.WithMsg(fmt.Sprintf("username: %v, groupDNList=%+v", user.Username, user.GroupDNList))
				count++
			}
		}
		ret.WithMsg(fmt.Sprintf("Found %v users in this group", count))
		return ret.Suc()
	}
	return ret.Fail().WithMsg("No admin group found!")

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func SearchUser(ldapCfg *models.LdapConfig, test *Test) *models.LDAPTestResult {
	ret := &models.LDAPTestResult{}
	fmt.Printf("Start to ping LDAP server: %v\n", ldapCfg.LdapURL)
	session, err := ldap.CreateWithAllConfig(ldapCfg.LdapConf, ldapCfg.LdapGroupConf)
	if err != nil {
		return ret.FailWithError(err)
	}
	session.Open()
	defer session.Close()

	fmt.Println("Start to search users")
	result, err := session.SearchUser(test.User)
	if err != nil {
		return ret.FailWithError(err)
	}
	if len(result) == 0 {
		return ret.Fail().WithMsg(fmt.Sprintf("No LDAP user found in current search conditions."))
	} else {
		ret.Suc().WithMsg(fmt.Sprintf("Found %d LDAP users in current search conditions", len(result)))
	}
	if len(test.User) == 0 {
		return ret.Fail().WithMsg("Please provide username to search")
	}
	fmt.Printf("Trying to find user %v\n")
	singleUser, err := session.SearchUser(test.User)
	if err != nil {
		return ret.FailWithError(err)
	}
	if len(singleUser) == 0 {
		return ret.Fail().WithMsg(fmt.Sprintf("The user %v is not found!", test.User))
	} else {
		ret.WithMsg(fmt.Sprintf("User %v found!", test.User))
		if len(singleUser[0].GroupDNList) == 0 {
			ret.WithMsg("Current user is not in any ldap group.")
		} else {
			for _, dn := range singleUser[0].GroupDNList {
				ret.WithMsg(fmt.Sprintf("User in the group with dn: [%v] OnboardGroup: %v", dn, SearchGroup(session, dn)))
			}
		}
	}
	return ret
}

func SearchGroup(session *ldap.Session, groupDN string) bool {
	groups, err := session.SearchGroupByDN(groupDN)
	if err != nil {
		fmt.Printf("Failed to search group, error found %v", err)
		return false
	}
	if len(groups) == 0 {
		return false
	}
	return true
}

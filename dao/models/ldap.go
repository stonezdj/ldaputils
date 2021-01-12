package models

import (
	"fmt"
	"github.com/goharbor/harbor/src/common/models"
	"gorm.io/gorm"
)

type LdapConfig struct {
	gorm.Model
	Name string `json:"name"`
	models.LdapConf
	//LdapURL               string `json:"ldap_url"`
	//LdapSearchDn          string `json:"ldap_search_dn"`
	//LdapSearchPassword    string `json:"ldap_search_password"`
	//LdapBaseDn            string `json:"ldap_base_dn"`
	//LdapFilter            string `json:"ldap_filter"`
	//LdapUID               string `json:"ldap_uid"`
	//LdapScope             int    `json:"ldap_scope"`
	//LdapConnectionTimeout int    `json:"ldap_connection_timeout"`
	//LdapVerifyCert        bool   `json:"ldap_verify_cert"`
	models.LdapGroupConf
	//LdapGroupBaseDN              string `json:"ldap_group_base_dn,omitempty"`
	//LdapGroupFilter              string `json:"ldap_group_filter,omitempty"`
	//LdapGroupNameAttribute       string `json:"ldap_group_name_attribute,omitempty"`
	//LdapGroupSearchScope         int    `json:"ldap_group_search_scope"`
	//LdapGroupAdminDN             string `json:"ldap_group_admin_dn,omitempty"`
	//LdapGroupMembershipAttribute string `json:"ldap_group_membership_attribute,omitempty"`
}

type LdapConfigCollection struct {
	Items []LdapConfig `json:"items"`
}

type LDAPTestResult struct {
	Success bool     `json:"success"`
	Message []string `json:"message"`
}

func (l *LDAPTestResult) WithMsg(msg string) *LDAPTestResult {
	l.Message = append(l.Message, msg)
	return l
}
func (l *LDAPTestResult) Suc() *LDAPTestResult {
	l.Success = true
	return l
}
func (l *LDAPTestResult) Fail() *LDAPTestResult {
	l.Success = false
	return l
}
func (l *LDAPTestResult) FailWithError(err error) *LDAPTestResult {
	return l.Fail().WithMsg(fmt.Sprintf("Failed with error %v", err))
}

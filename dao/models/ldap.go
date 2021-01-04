package models

import "gorm.io/gorm"

type LdapConfig struct {
	gorm.Model
	Name                  string `json:"name"`
	LdapURL               string `json:"ldap_url"`
	LdapSearchDn          string `json:"ldap_search_dn"`
	LdapSearchPassword    string `json:"ldap_search_password"`
	LdapBaseDn            string `json:"ldap_base_dn"`
	LdapFilter            string `json:"ldap_filter"`
	LdapUID               string `json:"ldap_uid"`
	LdapScope             int    `json:"ldap_scope"`
	LdapConnectionTimeout int    `json:"ldap_connection_timeout"`
	LdapVerifyCert        bool   `json:"ldap_verify_cert"`

	LdapGroupBaseDN              string `json:"ldap_group_base_dn,omitempty"`
	LdapGroupFilter              string `json:"ldap_group_filter,omitempty"`
	LdapGroupNameAttribute       string `json:"ldap_group_name_attribute,omitempty"`
	LdapGroupSearchScope         int    `json:"ldap_group_search_scope"`
	LdapGroupAdminDN             string `json:"ldap_group_admin_dn,omitempty"`
	LdapGroupMembershipAttribute string `json:"ldap_group_membership_attribute,omitempty"`
}

type LdapConfigCollection struct {
	Items []LdapConfig `json:"items"`
}

type LDAPTestResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

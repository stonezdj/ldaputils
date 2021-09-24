package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/goharbor/harbor/src/common/utils/ldap"

	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/lib/log"
)

//LDAPConfigAll ...
type LDAPConfigAll struct {
	LDAPConf      models.LdapConf      `json:"ldap_conf,omitempty"`
	LDAPGroupConf models.LdapGroupConf `json:"ldap_group_conf,omitempty"`
}

var jsonFileContent = `{
    "ldap_conf":
    {
    "ldap_url":"10.193.23.3",
    "ldap_search_dn":"cn=admin,dc=example,dc=com",
    "ldap_search_password":"admin",
    "ldap_base_dn":"dc=example,dc=com",
    "ldap_filter":"",
    "ldap_uid":"cn",
    "ldap_scope": 2,
    "ldap_verify_cert": false
    },
    "ldap_group_conf":
    {
    "ldap_group_base_dn":"dc=example,dc=com",
    "ldap_group_filter": "objectClass=groupOfNames",
    "ldap_group_name_attribute":"cn",
    "ldap_group_search_scope":2,
    "ldap_group_admin_dn":"",
    "ldap_group_membership_attribute": "memberof"
    }
}`

func genLdapJson() {
	f, err := os.Create("ldap.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	f.Write([]byte(jsonFileContent))
}

func main() {
	webserver := flag.Bool("webserver", false, "Run ldap utils as a web server")
	configJSON := flag.String("config", "ldap.json", "LDAP json file")
	username := flag.String("username", "mike_0", "search this username in LDAP")
	testName := flag.String("test", "basic", "run the test, can be basic, admin, admingroup, bind etc")
	password := flag.String("password", "", "The password of the ldap user")

	flag.Parse()

	if *webserver {
		WebServer()
		return
	}

	jsonFile, err := os.Open(*configJSON)
	if os.IsNotExist(err) {
		fmt.Println("The ldap.json file doesn't exist, create a new one")
		genLdapJson()
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var ldapConfigAll LDAPConfigAll
	json.Unmarshal([]byte(byteValue), &ldapConfigAll)

	//log.Debugf("ldapConfigAll %#v\n", ldapConfigAll)

	fmt.Println("================================================")
	fmt.Println("Verify basic LDAP information")

	session, err := ldap.CreateWithAllConfig(ldapConfigAll.LDAPConf, ldapConfigAll.LDAPGroupConf)
	CheckError(err)
	session.Open()
	defer session.Close()
	err = session.Bind(ldapConfigAll.LDAPConf.LdapSearchDn, ldapConfigAll.LDAPConf.LdapSearchPassword)
	CheckError(err)

	switch *testName {
	case "basic":
		if ok := SearchUser(session, username); !ok {
			fmt.Println("test failed")
			return
		}
	case "group":
		if ok := VerifyGroupConfig(ldapConfigAll, session); !ok {
			fmt.Println("test failed")
			return
		}
	case "admingroup":
		if ok := VerifyAdminGroupConfig(ldapConfigAll, session); !ok {
			fmt.Println("test failed")
			return
		}
	case "ping":
		if ok := Ping(ldapConfigAll, err); !ok {
			fmt.Println("ping test failed!")
		}
	case "bind":
		fmt.Println("bind success!")

	case "login":
		if ok := Login(session, *username, *password); !ok {
			fmt.Println("failed to login user!")
		}
		fmt.Println("Login user success!")
	}
}

func VerifyGroupConfig(ldapConfigAll LDAPConfigAll, session *ldap.Session) bool {
	fmt.Println("Verify LDAP group configurations")
	if len(ldapConfigAll.LDAPGroupConf.LdapGroupBaseDN) > 0 {
		fmt.Println("Trying to search group in current search conditions.")
		groups, err := session.SearchGroupByName("harbor_group_0")
		if err != nil {
			fmt.Println(err)
			return false
		}
		if len(groups) == 0 {
			DumpResult("No user group found!")
		} else {
			DumpResult(fmt.Sprintf("Found %v groups in current condition.", len(groups)))
		}
	}
	return true
}
func CheckError(err error) {
	if err != nil {
		fmt.Printf("Error %v,\n", err)
		os.Exit(1)
	}
}

func VerifyAdminGroupConfig(ldapConfigAll LDAPConfigAll, session *ldap.Session) bool {
	if len(ldapConfigAll.LDAPGroupConf.LdapGroupAdminDN) > 0 {
		fmt.Printf("Trying to search the group with admin privileges, ldap group admin dn: %v\n", ldapConfigAll.LDAPGroupConf.LdapGroupAdminDN)
		groups, err := session.SearchGroupByDN(ldapConfigAll.LDAPGroupConf.LdapGroupAdminDN)
		if err != nil {
			fmt.Println(err)
			return false
		}

		DumpResult(fmt.Sprintf("Found %v groups with admin privileges.", len(groups)))
		if len(groups) > 0 {
			fmt.Printf("Trying to find users in this group: %v\n", ldapConfigAll.LDAPGroupConf.LdapGroupAdminDN)
			count := 0
			userList, err := session.SearchUser("")
			if err != nil {
				fmt.Println(err)
			}
			for _, user := range userList {
				log.Debugf("username: %v, groupDNList=%+v", user.Username, user.GroupDNList)
				if stringInSlice(ldapConfigAll.LDAPGroupConf.LdapGroupAdminDN, user.GroupDNList) {
					count++
				}
			}
			DumpResult(fmt.Sprintf("Found %v users in this group", count))

		}

	}
	fmt.Println("================================================")
	return true
}

func SearchUser(session *ldap.Session, username *string) bool {
	fmt.Println("Start to search users")
	result, err := session.SearchUser("mike_0")
	if err != nil {
		DumpResult(fmt.Sprintf("Failed to search LDAP, error %v", err))
		return false
	}
	if len(result) == 0 {
		DumpResult(fmt.Sprintf("No LDAP user found in current search conditions."))
	} else {
		DumpResult(fmt.Sprintf("Found %d LDAP users in current search conditions", len(result)))
	}
	if len(*username) > 0 {
		fmt.Printf("Trying to find user %v\n", *username)
		singleUser, err := session.SearchUser(*username)
		if err != nil {
			fmt.Println(err)
			return false
		}
		if len(singleUser) == 0 {
			DumpResult(fmt.Sprintf("The user %v is not found!", *username))
		} else {
			DumpResult(fmt.Sprintf("User %v found!", *username))
			if len(singleUser[0].Email) == 0 {
				DumpResult("User email is empty!")
			} else {
				DumpResult(fmt.Sprintf("User email is %v", singleUser[0].Email))
			}
			if len(singleUser[0].GroupDNList) == 0 {
				DumpResult("Current user is not in any ldap group.")
			} else {
				for _, dn := range singleUser[0].GroupDNList {
					DumpResult(fmt.Sprintf("User in the group with dn: [%v] OnboardGroup: %v", dn, SearchGroup(session, dn)))
				}
			}
		}
	}
	fmt.Println("================================================")
	return true
}

func Ping(ldapConfigAll LDAPConfigAll, err error) bool {
	fmt.Printf("Start to ping LDAP server: %v\n", ldapConfigAll.LDAPConf.LdapURL)
	err = ldap.ConnectionTestWithAllConfig(ldapConfigAll.LDAPConf, ldapConfigAll.LDAPGroupConf)
	if err != nil {
		fmt.Printf("Error at connection test, %+v\n", err)
		return false
	}
	DumpResult("Success to ping LDAP server")
	return true
}

func SearchGroup(session *ldap.Session, groupDN string) bool {
	groups, err := session.SearchGroupByDN(groupDN)
	CheckError(err)
	if len(groups) == 0 {
		return false
	}
	return true
}

func Login(session *ldap.Session, username, password string) bool {
	ldapUsers, err := session.SearchUser(username)
	CheckError(err)
	if len(ldapUsers) == 0 {
		log.Warningf("Not found an entry.")
		return false
	} else if len(ldapUsers) != 1 {
		log.Warningf("Found more than one entry.")
		return false
	}
	log.Debugf("Found ldap user %+v", ldapUsers[0])

	dn := ldapUsers[0].DN
	if err = session.Bind(dn, password); err != nil {
		log.Warningf("Failed to bind user, username: %s, dn: %s, error: %v", username, dn, err)
		return false
	}
	return true
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// DumpResult ...
func DumpResult(msg string) {
	fmt.Println("==> " + msg)
}

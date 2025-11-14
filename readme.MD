# LDAP utilities


## Usage

Query LDAP user
```
docker run -it firstfloor/ldaputils:2.0 bash
>ldaputils
The ldap.json file doesn't exist, create a new one

# update the json with your ldap configurations
>vi ldap.json

# Search the ldap user with this command
>ldaputils -username mike02
================================================
Verify basic LDAP information
Start to search users
==> No LDAP user found in current search conditions.
Trying to find user mike02
==> User mike02 found!
==> User email is empty!
==> User in the group with dn: [cn=harbor_users,ou=groups,dc=example,dc=com] OnboardGroup: true
================================================
```

Launch: 
./ldaputils -webserver

Go to url:
http://localhost:8080/
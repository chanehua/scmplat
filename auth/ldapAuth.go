package auth

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego"
	"gopkg.in/ldap.v2"
)

func LdapAuth(uname, pwd string) error {
	//获取ldap认证基本信息
	ldapSection, err := beego.AppConfig.GetSection("ldap")
	if err != nil {
		return err
	}
	ldapdomain := ldapSection["ldapdomain"]
	ldapport := ldapSection["ldapport"]
	binduame := ldapSection["binduame"]
	bindpwd := ldapSection["bindpwd"]
	basedn := ldapSection["basedn"]

	//connet ldap server
	l, err := ldap.Dial(`tcp`, fmt.Sprintf("%s:%s", ldapdomain, ldapport))
	if err != nil {
		return err
	}
	defer l.Close()

	// First bind with a read only user
	err = l.Bind(binduame, bindpwd)
	if err != nil {
		return err
	}

	// Search for the given username
	filterStr := fmt.Sprintf("(&(objectClass=organizationalPerson)(mail=%s@*))", uname)
	searchRequest := ldap.NewSearchRequest(
		basedn, ldap.ScopeSingleLevel,
		ldap.NeverDerefAliases, 0, 0, false,
		filterStr,
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return err
	}
	if len(sr.Entries) != 1 {
		err = errors.New("User does not exist or too many entries returned")
		return err
	}

	userdn := sr.Entries[0].DN

	// Bind as the user to verify their password
	err = l.Bind(userdn, pwd)
	if err != nil {
		return err
	}

	// Rebind as the read only user for any futher queries
	err = l.Bind(binduame, bindpwd)
	if err != nil {
		return err
	}

	return nil
}

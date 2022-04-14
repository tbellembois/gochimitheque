package ldap

import (
	"crypto/tls"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/tbellembois/gochimitheque/logger"
)

type LDAPConnection struct {
	IsEnabled bool
	l         *ldap.Conn
}

type LDAPSearchResult struct {
	NbResults int
	R         *ldap.SearchResult
}

var (
	// LDAPServerURL ldap(s)://host.
	LDAPServerURL string
	// LDAPServerUsername CN=adminread,OU=FOO,OU=COM,OU=users,DC=foo,DC=com.
	LDAPServerUsername string
	// LDAPServerPassword pAsswRd.
	LDAPServerPassword string
	// LDAPGroupSearchBaseDN.
	LDAPGroupSearchBaseDN string
	// LDAPGroupSearchFilter.
	LDAPGroupSearchFilter string
	// LDAPUserSearchBaseDN.
	LDAPUserSearchBaseDN string
	// LDAPUserSearchFilter.
	LDAPUserSearchFilter string
)

func Connect() (l *LDAPConnection, err error) {
	l = &LDAPConnection{
		IsEnabled: false,
	}

	if LDAPServerURL == "" {
		return
	}

	// LDAP connection.
	if l.l, err = ldap.DialURL(LDAPServerURL, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true})); err != nil {
		return
	}

	// Admin bind.
	if LDAPServerUsername != "" {
		if LDAPServerPassword != "" {
			if err = l.l.Bind(LDAPServerUsername, LDAPServerPassword); err != nil {
				return
			}
		} else {
			if err = l.l.UnauthenticatedBind(LDAPServerUsername); err != nil {
				return
			}
		}
	}

	l.IsEnabled = true

	return
}

func TestSearchUser(email string) (result *LDAPSearchResult, err error) {

	var ldapConnection *LDAPConnection

	if ldapConnection, err = Connect(); err != nil {
		return
	}

	result, err = ldapConnection.SearchUser(email)

	return
}

func (conn *LDAPConnection) SearchUser(email string) (result *LDAPSearchResult, err error) {
	result = &LDAPSearchResult{}

	searchRequest := ldap.NewSearchRequest(
		LDAPUserSearchBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 100, 0, false,
		fmt.Sprintf(LDAPUserSearchFilter, email),
		[]string{"*"},
		nil,
	)

	if result.R, err = conn.l.Search(searchRequest); err != nil {
		logger.Log.Error(err)
		return
	}

	result.NbResults = len(result.R.Entries)

	return
}

func TestSearchGroup(partofname string) (result *LDAPSearchResult, err error) {

	var ldapConnection *LDAPConnection

	if ldapConnection, err = Connect(); err != nil {
		return
	}

	result, err = ldapConnection.SearchGroup(partofname)

	return
}

func (conn *LDAPConnection) SearchGroup(partofname string) (result *LDAPSearchResult, err error) {
	result = &LDAPSearchResult{}

	searchRequest := ldap.NewSearchRequest(
		LDAPGroupSearchBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 100, 0, false,
		fmt.Sprintf(LDAPGroupSearchFilter, partofname),
		[]string{"*"},
		nil,
	)

	if result.R, err = conn.l.Search(searchRequest); err != nil {
		logger.Log.Error(err)
		return
	}

	result.NbResults = len(result.R.Entries)

	return
}

func (conn *LDAPConnection) Bind(userdn string, password string) (err error) {
	err = conn.l.Bind(userdn, password)

	return
}

package mldap

import (
	"crypto/tls"
	"gopkg.in/ldap.v2"
	"fmt"
)

type LdapService struct {

}

func NewAdModel() *LdapService {
	return &LdapService{}
}

const (
	URL        = "ad-ds.bluetrans.cn"
	PORT       = 389
	USER_NAME   = "bluetrans"
	PASSWORD  = "bluetrans2021"
	BASE_DN     = "dc=bluetrans,dc=local"
)

func (m *LdapService) CheckLogin(name, password string) error {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", URL, PORT))
	if err != nil {
		return err
	}

	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return  err
	}

	// 绑定用于管理的用户
	err = l.Bind(USER_NAME, PASSWORD)
	if err != nil {
		return  err
	}

	// 查询
	sql := ldap.NewSearchRequest(BASE_DN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		// "(&(objectClass=organizationalPerson))",   // 查询所有人
		fmt.Sprintf("(&(objectCategory=Person)(sAMAccountName=%s))", name), //查询指定人
		[]string{"cn", "sAMAccountName"}, nil)

	cur, err := l.Search(sql)
	if err != nil {
		return err
	}

	if len(cur.Entries) == 0 {
		err = fmt.Errorf("%s does not exist", name)
		return err
	}

	if len(cur.Entries) > 1 {
		err = fmt.Errorf("exist multiple %s", name)
		return  err
	}

	userdn := cur.Entries[0].DN
	// 用户密码校验，一条对应的dn记录的密码校验
	err = l.Bind(userdn, password)
	if err != nil {
		return err
	}

	return nil
}
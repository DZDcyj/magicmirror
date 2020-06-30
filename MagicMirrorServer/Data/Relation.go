package Data

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

const CONFIRMED = 1
const UNCONFIRMED = 0

type Relation struct {
	Id         int
	SourceName string
	SourceAcct string
	TargetName string
	TargetAcct string
	Status     int // 状态，表示是否已经确认添加
}

type UserProtocol struct {
	Nickname string
	Account  string
}

func QueryAvailableFriendsByAccount(account string, users *[]*User) (length int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	length, err = qs.Filter("account__icontains", account).All(users)
	return
}

func QueryAvailableFriendsByNickname(nickname string, users *[]*User) (length int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	length, err = qs.Filter("nickname__icontains", nickname).All(users)
	return
}

func QueryRelationWithTwoAccounts(sourceAcct, targetAcct string) (result bool, err error) {
	result = false
	o := orm.NewOrm()
	qs := o.QueryTable("relation")
	var relations []*Relation
	length, err := qs.Filter("source_acct", sourceAcct).Filter("target_acct", targetAcct).All(&relations)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if length > 0 {
		result = true
	}
	return
}

func SendFriendRequest(sourceAcct, targetAcct, sourceName, targetName string) (status int, info string) {
	var users []*User
	length, err := QueryUsersByAccount(sourceAcct, &users)
	if err != nil {
		fmt.Println(err.Error())
		return 0, "Query error occurred"
	}
	if length == 0 {
		return 0, "Source user not found"
	}
	if users[0].Nickname != sourceName {
		return 0, "Source user account and nickname not matched"
	}
	length, err = QueryUsersByAccount(targetAcct, &users)
	if err != nil {
		fmt.Println(err.Error())
		return 0, "Query error occurred"
	}
	if length == 0 {
		return 0, "Target user not found"
	}
	if users[0].Nickname != targetName {
		return 0, "Target user account and nickname not matched"
	}
	result, _ := QueryRelationWithTwoAccounts(sourceAcct, targetAcct)
	if result {
		return 0, "Relation has existed"
	}
	forward := Relation{
		Id:         time.Now().Nanosecond(),
		SourceName: sourceName,
		SourceAcct: sourceAcct,
		TargetName: targetName,
		TargetAcct: targetAcct,
		Status:     UNCONFIRMED,
	}
	o := orm.NewOrm()
	_, _ = o.Insert(&forward)
	backward := Relation{
		Id:         time.Now().Nanosecond(),
		SourceName: targetName,
		SourceAcct: targetAcct,
		TargetName: sourceName,
		TargetAcct: sourceAcct,
		Status:     UNCONFIRMED,
	}
	_, _ = o.Insert(&backward)
	return 1, "Add unconfirmed relation successfully"
}

func ReceiveFriendRequests(account string, users *[]*Relation) {
	o := orm.NewOrm()
	qs := o.QueryTable("relation")
	_, _ = qs.Filter("source_acct", account).Filter("status", UNCONFIRMED).All(users)
}

func AcceptRequest(sourceAcct, targetAcct string) {
	o := orm.NewOrm()
	res, _ := QueryRelationWithTwoAccounts(sourceAcct, targetAcct)
	if !res {
		return
	}
	qs := o.QueryTable("relation")
	var forward []*Relation
	_, _ = qs.Filter("source_acct", sourceAcct).Filter("target_acct", targetAcct).All(&forward)
	for _, relation := range forward {
		relation.Status = CONFIRMED
		_, _ = o.Update(relation)
	}
	_, _ = qs.Filter("target_acct", sourceAcct).Filter("source_acct", targetAcct).All(&forward)
	for _, relation := range forward {
		relation.Status = CONFIRMED
		_, _ = o.Update(relation)
	}
}

func GetFriendListByAccount(account string) (protocols []*UserProtocol) {
	o := orm.NewOrm()
	qs := o.QueryTable("relation")
	var relations []*Relation
	n, err := qs.Filter("source_acct", account).Filter("status", CONFIRMED).All(&relations)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if n > 0 {
		for _, relation := range relations {
			protocol := UserProtocol{
				Nickname: relation.TargetName,
				Account:  relation.TargetAcct,
			}
			protocols = append(protocols, &protocol)
		}
	}
	return
}

func DeleteRelation(sourceAcct, targetAcct string) {
	o := orm.NewOrm()
	qs := o.QueryTable("relation")
	var relations []*Relation
	n, err := qs.Filter("source_acct", sourceAcct).Filter("target_acct", targetAcct).All(&relations)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if n > 0 {
		for _, relation := range relations {
			_, _ = o.Delete(relation)
		}
	}
	n, err = qs.Filter("source_acct", targetAcct).Filter("target_acct", sourceAcct).All(&relations)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if n > 0 {
		for _, relation := range relations {
			_, _ = o.Delete(relation)
		}
	}
}

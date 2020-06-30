package Data

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
)

const MALE = 0
const FEMALE = 1

type User struct {
	Id       int
	Account  string
	Password string
	Nickname string
	Gender   int
	Birthday string
	Age      int
	PhotoId  int
}

func AddUser(user *User) error {
	o := orm.NewOrm()
	_, err := o.Insert(user)
	return err
}

func (user *User) ToString() string {
	var gender string
	if user.Gender == MALE {
		gender = "Male"
	} else if user.Gender == FEMALE {
		gender = "Female"
	} else {
		gender = "Unknown"
	}
	return "Id: " + strconv.Itoa(user.Id) + "\nAccount: " + user.Account + "\nPassword: " + user.Password + "\nNickname: " + user.Nickname + "\nGender: " + gender + "\nBirthday: " + user.Birthday + "\nAge: " + strconv.Itoa(user.Age)
}

func QueryUserById(Id int) (result User, err error) {
	o := orm.NewOrm()
	result.Id = Id
	err = o.Read(&result)
	return
}

func QueryUserByAccount(account string) User {
	var users []*User
	n, err := QueryUsersByAccount(account, &users)
	if err != nil {
		fmt.Println(err.Error())
		return User{}
	}
	if n > 0 {
		return *users[0]
	}
	return User{}
}

func QueryUsersByAccount(account string, users *[]*User) (length int64, err error) {
	o := orm.NewOrm()
	var user User
	qs := o.QueryTable(user)
	length, err = qs.Filter("account", account).All(users)
	return
}

func DeleteUserById(Id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&User{Id: Id})
	return err
}

func UpdateUserById(Id int, user *User) error {
	if Id != user.Id {
		return errors.New("id not match")
	}
	o := orm.NewOrm()
	_, err := o.Update(user)
	return err
}

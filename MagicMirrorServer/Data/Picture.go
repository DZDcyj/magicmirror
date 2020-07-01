package Data

import "github.com/astaxie/beego/orm"

type Picture struct {
	Id       int
	Url      string
	ParentId int
}

func GetPicturesById(parentId int, result *[]*Picture) {
	o := orm.NewOrm()
	qs := o.QueryTable("picture")
	_, _ = qs.Filter("parent_id", parentId).All(result)
}

func AddPicture(picture *Picture) {
	o := orm.NewOrm()
	_, _ = o.Insert(picture)
}


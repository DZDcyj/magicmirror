package Data

import "github.com/astaxie/beego/orm"

type Post struct {
	Id      int
	Content string
	Time    string
}

func GetAllPosts(posts *[]*Post) {
	o := orm.NewOrm()
	qs := o.QueryTable("post")
	_, _ = qs.All(posts)
}

func QueryPostFromId(id int) (result Post, err error) {
	result.Id = id
	o := orm.NewOrm()
	err = o.Read(&result)
	return
}

func PublishPost(post *Post) error {
	o := orm.NewOrm()
	_, err := o.Insert(post)
	return err
}

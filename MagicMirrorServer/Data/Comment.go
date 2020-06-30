package Data

import "github.com/astaxie/beego/orm"

type Comment struct {
	Id      int
	Content string
	PostId  int
	Time    string
}

func PublishComment(comment *Comment) error {
	o := orm.NewOrm()
	_, err := o.Insert(comment)
	return err
}

func GetCommentsByPostId(postId int, comments *[]*Comment) {
	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	_, _ = qs.Filter("post_id", postId).All(comments)
}

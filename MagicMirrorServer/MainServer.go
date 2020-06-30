package main

import (
	. "./Data"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"log"
	"net/http"
	"strconv"
)
import _ "github.com/go-sql-driver/mysql"

// 这里修改相应的数据库信息，包括用户名密码等
// 在使用前务必在 mysql 中创建对应名称的数据库
const DatabaseUsername = "chin"
const DatabasePassword = "gerenyinsi"
const DatabaseAddress = "localhost"
const DatabasePort = "3306"
const DatabaseName = "magicmirror"

const Port = "50000"

const SUCCESS = 1
const FAILURE = 0

func init() {
	dataSource := DatabaseUsername + ":" + DatabasePassword + "@tcp(" + DatabaseAddress + ":" + DatabasePort + ")" + "/" + DatabaseName + "?charset=utf8"
	err := orm.RegisterDataBase("default", "mysql", dataSource, 30)
	if err != nil {
		panic(err)
	}
	orm.RegisterModel(new(Picture), new(User), new(Relation), new(Post), new(Comment))
	_ = orm.RunSyncdb("default", false, true)
}

type Response struct {
	Status int
	Info   string
}

type Picture struct {
	Id       int
	Path     string
	ParentId string
}

// 根据提供的账号/昵称（至多一项）返回对应的用户列表
func dealWithUserListRequest(w http.ResponseWriter, r *http.Request) {
	var protocols []*User
	_ = r.ParseForm()
	acct := r.Form.Get("account")
	name := r.Form.Get("name")
	resp := Response{
		Status: FAILURE,
	}
	if acct == "" && name == "" {
		resp.Info = "Account or name can't be both empty"
		sendResponse(w, resp)
		return
	}
	if acct != "" && name != "" {
		resp.Info = "Can only use one method to get list"
		sendResponse(w, resp)
		return
	}
	if acct != "" {
		_, _ = QueryAvailableFriendsByAccount(acct, &protocols)
	} else {
		_, _ = QueryAvailableFriendsByNickname(name, &protocols)
	}
	res, _ := json.Marshal(protocols)
	_, _ = fmt.Fprint(w, string(res))
}

// 处理登录
func dealWithLogin(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	resp := Response{
		Status: SUCCESS,
		Info:   "Login successfully",
	}
	acct := r.Form.Get("account")
	if acct == "" {
		resp.Status = FAILURE
		resp.Info = "Account empty"
		sendResponse(w, resp)
		return
	}
	user := QueryUserByAccount(acct)
	password := r.Form.Get("password")
	if password == "" {
		resp.Status = FAILURE
		resp.Info = "Password empty"
		sendResponse(w, resp)
		return
	}
	if user.Password != password {
		resp.Status = FAILURE
		resp.Info = "Invalid account or password"
	}
	sendResponse(w, resp)
}

// 发送对应的反馈体
func sendResponse(w http.ResponseWriter, response Response) {
	res, _ := json.Marshal(response)
	_, _ = fmt.Fprint(w, string(res))
}

// 处理注册
func dealWithRegister(w http.ResponseWriter, r *http.Request) {
	var user User
	response := Response{
		Status: SUCCESS,
		Info:   "Registered successfully",
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		_ = r.Body.Close()
		response.Status = FAILURE
		response.Info = "Error occurred in json decoder"
		sendResponse(w, response)
		log.Fatal(err)
	}
	var users []*User
	n, err := QueryUsersByAccount(user.Account, &users)
	if err != nil {
		response.Status = FAILURE
		response.Info = "Error occurred in query user"
		sendResponse(w, response)
		log.Fatal(err)
	}
	if n > 0 {
		response.Status = FAILURE
		response.Info = "This account has been registered"
		sendResponse(w, response)
		return
	}
	_ = AddUser(&user)
	sendResponse(w, response)
}

// 添加好友关系请求
func addFriendRequest(w http.ResponseWriter, r *http.Request) {
	var relation Relation
	response := Response{
		Status: SUCCESS,
		Info:   "Added relation successfully",
	}
	if err := json.NewDecoder(r.Body).Decode(&relation); err != nil {
		_ = r.Body.Close()
		response.Status = FAILURE
		response.Info = "Error occurred in json decoder"
		sendResponse(w, response)
		log.Fatal(err)
	}
	response.Status, response.Info = SendFriendRequest(relation.SourceAcct, relation.TargetAcct, relation.SourceName, relation.TargetName)
	sendResponse(w, response)
}

// 获取好友请求列表
func getFriendRequests(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	resp := Response{
		Status: FAILURE,
	}
	acct := r.Form.Get("account")
	if acct == "" {
		resp.Info = "Account empty"
		sendResponse(w, resp)
		return
	}
	var relations []*Relation
	ReceiveFriendRequests(acct, &relations)
	res, err := json.Marshal(relations)
	if err != nil {
		resp.Info = "Marshal error"
		sendResponse(w, resp)
		return
	}
	_, _ = fmt.Fprint(w, string(res))
}

// 接受好友申请
func dealWithAcceptation(w http.ResponseWriter, r *http.Request) {
	var relation Relation
	response := Response{
		Status: FAILURE,
	}
	if err := json.NewDecoder(r.Body).Decode(&relation); err != nil {
		_ = r.Body.Close()
		response.Info = "Error occurred in json decoder"
		sendResponse(w, response)
		log.Fatal(err)
	}
	AcceptRequest(relation.SourceAcct, relation.TargetAcct)
}

// 删除关系（可用于拒绝申请/删除好友）
func dealWithRefusal(w http.ResponseWriter, r *http.Request) {
	var relation Relation
	response := Response{
		Status: FAILURE,
	}
	if err := json.NewDecoder(r.Body).Decode(&relation); err != nil {
		_ = r.Body.Close()
		response.Info = "Error occurred in json decoder"
		sendResponse(w, response)
		log.Fatal(err)
	}
	DeleteRelation(relation.SourceAcct, relation.TargetAcct)
}

func publishPost(w http.ResponseWriter, r *http.Request) {
	var post Post
	response := Response{
		Status: FAILURE,
	}
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		_ = r.Body.Close()
		response.Info = "Error occurred in json decoder"
		sendResponse(w, response)
		log.Fatal(err)
	}
	_ = PublishPost(&post)
}

func commitComment(w http.ResponseWriter, r *http.Request) {
	var comment Comment
	response := Response{
		Status: FAILURE,
	}
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		_ = r.Body.Close()
		response.Info = "Error occurred in json decoder"
		sendResponse(w, response)
		log.Fatal(err)
	}
	_ = PublishComment(&comment)
}

// 发送所有的帖子
func getAllPost(w http.ResponseWriter, _ *http.Request) {
	var posts []*Post
	GetAllPosts(&posts)
	res, _ := json.Marshal(posts)
	_, _ = fmt.Fprint(w, string(res))
}

// 获取对应账号的用户信息
func getUserMessage(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	account := r.Form.Get("account")
	if account == "" {
		resp := Response{
			Status: FAILURE,
			Info:   "Account empty",
		}
		sendResponse(w, resp)
		return
	}
	user := QueryUserByAccount(account)
	res, _ := json.Marshal(user)
	_, _ = fmt.Fprint(w, string(res))
}

// 获取对应账号好友列表
func getFriendList(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	acct := r.Form.Get("account")
	if acct == "" {
		resp := Response{
			Status: FAILURE,
			Info:   "Account empty",
		}
		sendResponse(w, resp)
		return
	}
	friends := GetFriendListByAccount(acct)
	if friends == nil {
		friends = []*UserProtocol{}
	}
	result, _ := json.Marshal(friends)
	_, _ = fmt.Fprintf(w, string(result))
}

// 修改对应账号用户信息
func changeUserMessage(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	acct := r.Form.Get("account")
	resp := Response{
		Status: SUCCESS,
		Info:   "Changed successfully",
	}
	if acct == "" {
		resp.Status = FAILURE
		resp.Info = "Account empty"
		sendResponse(w, resp)
		return
	}
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		_ = r.Body.Close()
		resp.Info = "Error occurred in json decoder"
		sendResponse(w, resp)
		log.Fatal(err)
	}
	_ = UpdateUserById(user.Id, &user)
	sendResponse(w, resp)
}

// 获取对应的帖子
func getPostById(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	strId := r.Form.Get("postid")
	if strId == "" {
		resp := Response{
			Status: FAILURE,
			Info:   "Post id empty",
		}
		sendResponse(w, resp)
		return
	}
	postId, _ := strconv.Atoi(strId)
	post, _ := QueryPostFromId(postId)
	res, _ := json.Marshal(post)
	_, _ = fmt.Fprint(w, string(res))
}

// 获取对应帖子的评论列表
func getCommentsByPostId(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	strId := r.Form.Get("postid")
	if strId == "" {
		resp := Response{
			Status: FAILURE,
			Info:   "Post id empty",
		}
		sendResponse(w, resp)
		return
	}
	postId, _ := strconv.Atoi(strId)
	var comments []*Comment
	GetCommentsByPostId(postId, &comments)
	res, _ := json.Marshal(comments)
	_, _ = fmt.Fprint(w, string(res))
}

func main() {
	fmt.Println("Starting server...")
	http.HandleFunc("/getAllPosts", getAllPost)
	http.HandleFunc("/getUserMessage", getUserMessage)
	http.HandleFunc("/register", dealWithRegister)
	http.HandleFunc("/login", dealWithLogin)
	http.HandleFunc("/queryUsers", dealWithUserListRequest)
	http.HandleFunc("/addFriendRequest", addFriendRequest)
	http.HandleFunc("/getFriendRequests", getFriendRequests)
	http.HandleFunc("/acceptFriendRequest", dealWithAcceptation)
	http.HandleFunc("/removeFriendRelation", dealWithRefusal)
	http.HandleFunc("/getFriendList", getFriendList)
	http.HandleFunc("/changeUserMessage", changeUserMessage)
	http.HandleFunc("/getComments", getCommentsByPostId)
	http.HandleFunc("/getPost", getPostById)
	http.HandleFunc("/publishPost", publishPost)
	http.HandleFunc("/commitComment", commitComment)

	err := http.ListenAndServe(":"+Port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

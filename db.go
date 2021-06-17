package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db, err = sql.Open("sqlite3", "./forum.db")

type User struct {
	Id           int
	Username     string
	Email        string
	Password     string
	Avatar       string
	RegisterDate string
	Ban          int
	Rank         int
}

type UserInfo struct {
	Id           int
	Username     string
	Email        string
	Avatar       string
	RegisterDate string
	Ban          int
	Rank         int
	Badges       []Badges
	Permissions  Ranks
	IsAuthorized bool
}

type Ranks struct {
	Id               int
	Name             string
	AdminPanelAccess int
	Ban              int
	Deban            int
	DeletePost       int
	DeleteUser       int
	ModifyCategory   int
	TicketAccess     int
	ViewClosedPost   int
	DeleteComment    int
	BadgeAttribution int
	ModifyRank       int
}

type CategoriesUsersRanks struct {
	Id         int
	UserId     int
	CategoryId int
	Category   Category
	UserInfo   UserInfo
}

type AllCategoriesUserRanks struct {
	Id       int
	Category []Category
	UserInfo UserInfo
}

type Reactions struct {
	Id    int
	Title string
	Icon  string
}

type Notifications struct {
	Id               int
	NotificationType int
	Subject          string
	Action           string
	FromId           int
	ToId             int
}

type AccessTokens struct {
	Id     int
	Token  string
	Date   string
	UserId int
}

type Logs struct {
	Id           int
	Subject      string
	LogsType     int
	Action       string
	CreationDate string
	UserId       int
}

type Tickets struct {
	Id           int
	Title        string
	Subject      string
	CreationDate string
	Status       int
	Creator      UserInfo
}

type TicketsMessages struct {
	Id       int
	Comment  string
	Creator  UserInfo
	UserId   int
	TicketId int
}

type Promotes struct {
	Id         int
	EndDate    string
	PostId     int
	CommentId  int
	CategoryId int
}

// Init all the db tables
func InitTable() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//_, error := db.Exec("DROP TABLE IF EXISTS users;DROP TABLE IF EXISTS  categories;DROP TABLE IF EXISTS  posts;DROP TABLE IF EXISTS  comments;DROP TABLE IF EXISTS  ranks;DROP TABLE IF EXISTS  badges;DROP TABLE IF EXISTS  user_badges;DROP TABLE IF EXISTS  categories_users_ranks;DROP TABLE IF EXISTS  logs;DROP TABLE IF EXISTS  notifications;DROP TABLE IF EXISTS  promotes;DROP TABLE IF EXISTS  reactions;DROP TABLE IF EXISTS  votes;DROP TABLE IF EXISTS  ticket_messages;DROP TABLE IF EXISTS  tickets;DROP TABLE IF EXISTS users_badges;")
	//if error != nil {
	//	log.Fatal(error)
	//}
	createTable := `
	-- badges definition
	
	CREATE TABLE IF NOT EXISTS "badges" (
		"id"    INTEGER,
		"name"    TEXT DEFAULT '',
		"icon"    TEXT DEFAULT '',
		"type"    INTEGER DEFAULT 0,
		"value"    INTEGER DEFAULT 0,
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	
	-- ranks definition
	
	CREATE TABLE IF NOT EXISTS "ranks" (
		"id"    INTEGER,
		"name"    TEXT DEFAULT '',
		"adminPanelAccess"    INTEGER DEFAULT 0,
		"ban"    INTEGER DEFAULT 0,
		"deban"    INTEGER DEFAULT 0,
		"deletePost"    INTEGER DEFAULT 0,
		"deleteUser"    INTEGER DEFAULT 0,
		"modifyCategory"    INTEGER DEFAULT 0,
		"ticketAccess"    INTEGER DEFAULT 0,
		"viewClosedPost"    INTEGER DEFAULT 0,
		"deleteComment"    INTEGER DEFAULT 0,
		"badgeAttribution"    INTEGER DEFAULT 0,
		"modifyRank"    INTEGER DEFAULT 0,
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	
	
	-- reactions definition
	
	CREATE TABLE IF NOT EXISTS  "reactions" (
		"id"	INTEGER,
		"title"	TEXT DEFAULT '',
		"icon"	TEXT DEFAULT '',
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	
	
	-- categories definition
	
	CREATE TABLE IF NOT EXISTS  "categories" (
		"id"	INTEGER,
		"title"	TEXT DEFAULT '',
		"description"	TEXT DEFAULT '',
		"icon"	TEXT DEFAULT '',
		"creationDate"	TEXT DEFAULT '',
		"parentId"	INTEGER  DEFAULT 0,
		FOREIGN KEY("parentId") REFERENCES "categories"("id"),
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	
	
	-- users definition
	
	CREATE TABLE IF NOT EXISTS  "users" (
		"id"	INTEGER,
		"username"	TEXT DEFAULT '',
		"email"	TEXT DEFAULT '',
		"password"	TEXT DEFAULT '',
		"avatar"	TEXT DEFAULT '',
		"registerDate"	TEXT DEFAULT '',
		"ban"	INTEGER DEFAULT 0,
		"rank"	INTEGER DEFAULT 1,
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("rank") REFERENCES "ranks"("id")
	);
	
	
	-- users_badges definition
	
	CREATE TABLE IF NOT EXISTS  "users_badges" (
		"id"	INTEGER,
		"badgeId"	INTEGER DEFAULT 0,
		"userId"	INTEGER DEFAULT 0,
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("userId") REFERENCES "users"("id"),
		FOREIGN KEY("badgeId") REFERENCES "badges"("id")
	);
	
	
	-- access_tokens definition
	
	CREATE TABLE IF NOT EXISTS  "access_tokens" (
		"id"	INTEGER,
		"token"	TEXT DEFAULT '',
		"date"	TEXT DEFAULT '',
		"userId"	INTEGER DEFAULT 0,
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("userId") REFERENCES "users"("id")
	);
	
	
	-- categories_users_ranks definition
	
	CREATE TABLE IF NOT EXISTS  "categories_users_ranks" (
		"id"	INTEGER,
		"userId"	INTEGER DEFAULT 0,
		"categoryId"	INTEGER DEFAULT 0,
		FOREIGN KEY("categoryId") REFERENCES "categories"("id"),
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("userId") REFERENCES "users"("id")
	);
	
	
	-- logs definition
	
	CREATE TABLE IF NOT EXISTS  "logs" (
		"id"	INTEGER,
		"subject"	TEXT DEFAULT '',
		"type"	INTEGER DEFAULT 0,
		"action"	TEXT DEFAULT '',
		"creationDate"	TEXT DEFAULT '',
		"userId"	INTEGER DEFAULT 0,
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("userId") REFERENCES "users"("id")
	);
	
	
	-- notifications definition
	
	CREATE TABLE IF NOT EXISTS  "notifications" (
		"id"	INTEGER,
		"type"	INTEGER DEFAULT 0,
		"subject"	TEXT DEFAULT '',
		"action"	TEXT DEFAULT '',
		"fromId"	INTEGER DEFAULT 0,
		"toId"	INTEGER  DEFAULT 0,
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("toId") REFERENCES "users"("id"),
		FOREIGN KEY("fromId") REFERENCES "users"("id")
	);
	
	
	-- posts definition
	
	CREATE TABLE IF NOT EXISTS  "posts" (
		"id"	INTEGER,
		"title"	TEXT DEFAULT '',
		"subject"	TEXT DEFAULT '',
		"creationDate"	TEXT DEFAULT '',
		"status"	INTEGER DEFAULT 0,
		"categoryId"	INTEGER DEFAULT 0,
		"userId"	INTEGER DEFAULT 0,
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("categoryId") REFERENCES "categories"("id"),
		FOREIGN KEY("userId") REFERENCES "users"("id")
	);
	
	
	-- tickets definition
	
	CREATE TABLE IF NOT EXISTS  "tickets" (
		"id"	INTEGER,
		"title"	TEXT DEFAULT '',
		"subject"	TEXT DEFAULT '',
		"creationDate"	TEXT DEFAULT '',
		"status"	INTEGER DEFAULT 0,
		"userId"	INTEGER DEFAULT 0,
		FOREIGN KEY("userId") REFERENCES "users"("id"),
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	
	
	-- comments definition
	
	CREATE TABLE IF NOT EXISTS  "comments" (
		"id"	INTEGER,
		"comment"	TEXT DEFAULT '',
		"creationDate"	TEXT DEFAULT '',
		"userId"	INTEGER DEFAULT 0,
		"postId"	INTEGER DEFAULT 0,
		"parentId"	INTEGER DEFAULT 0,
		"status"	INTEGER DEFAULT 0,
		FOREIGN KEY("userId") REFERENCES "users"("id"),
		FOREIGN KEY("postId") REFERENCES "posts"("id"),
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	
	
	-- promotes definition
	
	CREATE TABLE IF NOT EXISTS  "promotes" (
		"id"	INTEGER,
		"endDate"	TEXT DEFAULT '',
		"postId"	INTEGER DEFAULT 0,
		"commentId"	INTEGER DEFAULT 0,
		"categoryId"	INTEGER DEFAULT 0,
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("commentId") REFERENCES "comments"("id"),
		FOREIGN KEY("postId") REFERENCES "posts"("id")
	);
	
	
	-- ticket_messages definition
	
	CREATE TABLE IF NOT EXISTS  "ticket_messages" (
		"id"	INTEGER,
		"comment"	TEXT DEFAULT '',
		"userId"	INTEGER DEFAULT 0,
		"ticketId"	INTEGER DEFAULT 0,
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("ticketId") REFERENCES "tickets"("id"),
		FOREIGN KEY("userId") REFERENCES "users"("id")
	);
	
	
	-- votes definition
	
	CREATE TABLE IF NOT EXISTS "votes" (
		"id"	INTEGER,
		"type"	INTEGER DEFAULT 0,
		"postId"	INTEGER DEFAULT 0,
		"commentId"	INTEGER DEFAULT 0,
		"userId"	INTEGER DEFAULT 0,
		FOREIGN KEY("commentId") REFERENCES "comments"("id"),
		FOREIGN KEY("userId") REFERENCES "users"("id"),
		FOREIGN KEY("postId") REFERENCES "posts"("id"),
		FOREIGN KEY("type") REFERENCES "reactions"("id"),
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	`
	_, err := db.Exec(createTable)
	if err != nil {
		fmt.Print(err)
	}

}

// Generate the badges
func InitBadges() {
	InsertBadge("+10votes done", "https://image.flaticon.com/icons/png/512/1579/1579472.png", 1, 1)
	InsertBadge("1st post created", "https://image.flaticon.com/icons/png/512/1312/1312313.png", 2, 1)
	InsertBadge("+5 posts created", "https://image.flaticon.com/icons/png/512/1579/1579458.png", 2, 1)
	InsertBadge("+10votes on your post", "https://image.flaticon.com/icons/png/512/1029/1029183.png", 1, 1)
}

func GetUserByCookie(cookie string) int {
	var userId int
	row := db.QueryRow("SELECT userId FROM access_tokens WHERE token = ?", cookie)
	err := row.Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Zero rows found")
		} else {
			fmt.Print(err)
		}
	}
	return userId
}
func GetRank(id int) Ranks {
	row := db.QueryRow("SELECT * FROM ranks WHERE id = ?", id)
	var Temp Ranks
	err := row.Scan(&Temp.Id, &Temp.Name, &Temp.AdminPanelAccess, &Temp.Ban, &Temp.Deban, &Temp.DeletePost, &Temp.DeleteUser, &Temp.ModifyCategory, &Temp.TicketAccess, &Temp.ViewClosedPost, &Temp.DeleteComment, &Temp.BadgeAttribution, &Temp.ModifyRank)
	if err != nil {
		if err == sql.ErrNoRows {
			return Temp
		} else {
			fmt.Print(err)
		}
	}
	return Temp
}

func GetRanks() []Ranks {
	request, err := db.Prepare("SELECT id FROM ranks")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := request.Query()
	if err != nil {
		log.Fatal(err)
	}
	var TempTable []Ranks
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				return TempTable
			} else {
				log.Fatal(err)
			}
		}
		TempTable = append(TempTable, GetRank(id))
	}
	return TempTable
}

func GetUser(id int) User {
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	var Temp User
	err := row.Scan(&Temp.Id, &Temp.Username, &Temp.Email, &Temp.Password, &Temp.Avatar, &Temp.RegisterDate, &Temp.Ban, &Temp.Rank)
	if err != nil {
		if err == sql.ErrNoRows {
			return Temp
		} else {
			fmt.Print(err)
		}
	}
	return Temp
}

func GetUserIdByEmail(email string) int {
	row := db.QueryRow("SELECT id FROM users WHERE email = ?", email)
	var id int
	err := row.Scan(&id)
	if err != nil {
		fmt.Println(err)
	}
	return id
}

func GetUserInfo(id int) UserInfo {
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	var Temp UserInfo
	var notuse string
	err := row.Scan(&Temp.Id, &Temp.Username, &Temp.Email, &notuse, &Temp.Avatar, &Temp.RegisterDate, &Temp.Ban, &Temp.Rank)
	Temp.Permissions = GetRank(Temp.Rank)
	Temp.Badges = GetUserBadges(Temp.Id)
	Temp.RegisterDate = Temp.RegisterDate[:10]
	if err != nil {
		if err == sql.ErrNoRows {
			return Temp
		} else {
			fmt.Print(err)
		}
	}
	return Temp
}

func GetUsersInfo() []UserInfo {
	request, err := db.Prepare("SELECT id FROM users")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := request.Query()
	if err != nil {
		log.Fatal(err)
	}
	var TempTable []UserInfo
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				return TempTable
			} else {
				log.Fatal(err)
			}
		}
		TempTable = append(TempTable, GetUserInfo(id))
	}
	return TempTable
}

func GetNumberOfPosts(id int) int {
	row := db.QueryRow("SELECT COUNT(*) FROM posts WHERE categoryId = ?", id)
	var count int
	err := row.Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		} else {
			fmt.Print(err)
		}
	}
	return count
}

func GetCategoriesUsersRanksForUser(id int) []CategoriesUsersRanks {
	request, err := db.Prepare("SELECT * FROM categories_users_ranks WHERE userId = $1")
	if err != nil {
		fmt.Print(err)
	}
	rows, err := request.Query(id)
	if err != nil {
		fmt.Print(err)
	}
	var TempTable []CategoriesUsersRanks
	var Temp CategoriesUsersRanks
	for rows.Next() {
		err := rows.Scan(&Temp.Id, &Temp.UserId, &Temp.CategoryId)
		Temp.UserInfo = GetUserInfo(Temp.Id)
		Temp.Category = GetCategory(Temp.CategoryId)
		if err != nil {
			if err == sql.ErrNoRows {
				return TempTable
			} else {
				fmt.Print(err)
			}
		}
		TempTable = append(TempTable, Temp)
	}
	return TempTable
}

func GetALLCategoriesUsersRanksForUser(id int) []Category {
	request, err := db.Prepare("SELECT * FROM categories_users_ranks WHERE userId = $1")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := request.Query(id)
	if err != nil {
		log.Fatal(err)
	}
	var Temp CategoriesUsersRanks
	var Categories []Category
	for rows.Next() {
		err := rows.Scan(&Temp.Id, &Temp.UserId, &Temp.CategoryId)
		Categories = append(Categories, GetCategory(Temp.CategoryId))
		if err != nil {
			if err == sql.ErrNoRows {
				return Categories
			} else {
				log.Fatal(err)
			}
		}
	}
	return Categories
}

func GetIdTableSubCategoryRankForUser(id int) []int {
	request, err := db.Prepare("SELECT * FROM categories_users_ranks WHERE userId = $1")
	var table []int
	if err != nil {
		log.Fatal(err)
	}
	rows, err := request.Query(id)
	if err != nil {
		log.Fatal(err)
	}
	var Temp CategoriesUsersRanks
	for rows.Next() {
		err := rows.Scan(&Temp.Id, &Temp.UserId, &Temp.CategoryId)
		table = append(table, Temp.CategoryId)
		if err != nil {
			if err == sql.ErrNoRows {
				return table
			} else {
				log.Fatal(err)
			}
		}
	}
	return table
}

func GeAllCategoriesUsersRanks() []AllCategoriesUserRanks {
	request, err := db.Prepare("SELECT * FROM users WHERE rank = 2")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := request.Query()
	if err != nil {
		log.Fatal(err)
	}
	var CategoriesUserTable []AllCategoriesUserRanks
	var CategoriesUser AllCategoriesUserRanks
	var TempUserInfo UserInfo
	var notuse string
	for rows.Next() {
		err := rows.Scan(&TempUserInfo.Id, &TempUserInfo.Username, &TempUserInfo.Email, &notuse, &TempUserInfo.Avatar, &TempUserInfo.RegisterDate, &TempUserInfo.Ban, &TempUserInfo.Rank)
		if err != nil {
			if err == sql.ErrNoRows {
				return CategoriesUserTable
			} else {
				log.Fatal(err)
			}
		}
		TempUserInfo.Permissions = GetRank(TempUserInfo.Rank)
		TempUserInfo.RegisterDate = TempUserInfo.RegisterDate[5:10] + "-" + TempUserInfo.RegisterDate[0:4]
		CategoriesUser.UserInfo = TempUserInfo
		CategoriesUser.Id = TempUserInfo.Id
		CategoriesUser.Category = GetALLCategoriesUsersRanksForUser(TempUserInfo.Id)
		CategoriesUserTable = append(CategoriesUserTable, CategoriesUser)
	}
	return CategoriesUserTable
}

// func GetNotification(id int) Notifications {
// 	row := db.QueryRow("SELECT * FROM notifications WHERE id = ?", id)
// 	var Temp Notifications
// 	err := row.Scan(&Temp.Id, &Temp.NotificationType, &Temp.Subject, &Temp.Action, &Temp.FromId, &Temp.ToId)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return Temp
// 		} else {
// 			fmt.Print(err)
// 		}
// 	}
// 	return Temp
// }

// func GetNotifications(id int) []Notifications {
// 	request, err := db.Prepare("SELECT id FROM notifications WHERE toId = $1")
// 	if err != nil {
// 		fmt.Print(err)
// 	}
// 	rows, err := request.Query(id)
// 	if err != nil {
// 		fmt.Print(err)
// 	}
// 	var TempTable []Notifications
// 	for rows.Next() {
// 		var id int
// 		err := rows.Scan(&id)
// 		if err != nil {
// 			if err == sql.ErrNoRows {
// 				return TempTable
// 			} else {
// 				fmt.Print(err)
// 			}
// 		}
// 		TempTable = append(TempTable, GetNotification(id))
// 	}
// 	return TempTable
// }

func GetReaction(id int) Reactions {
	row := db.QueryRow("SELECT * FROM reactions WHERE id = ?", id)
	var Temp Reactions
	err := row.Scan(&Temp.Id, &Temp.Title, &Temp.Icon)
	if err != nil {
		if err == sql.ErrNoRows {
			return Temp
		} else {
			fmt.Print(err)
		}
	}
	return Temp
}

// func GetPromote(id int) Reactions {
// 	row := db.QueryRow("SELECT * FROM reactions WHERE id = ?", id)
// 	var Temp Reactions
// 	err := row.Scan(&Temp.Id, &Temp.Title, &Temp.Icon)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return Temp
// 		} else {
// 			fmt.Print(err)
// 		}
// 	}
// 	return Temp
// }

func GetPromoteByPostId(id int) Promotes {
	row := db.QueryRow("SELECT * FROM promotes WHERE postId = ?", id)
	var Temp Promotes
	err := row.Scan(&Temp.Id, &Temp.EndDate, &Temp.PostId, &Temp.CommentId)
	if err != nil {
		if err == sql.ErrNoRows {
			return Temp
		} else {
			fmt.Print(err)
		}
	}
	return Temp
}

func GetPromoteByCommentId(id int) Promotes {
	row := db.QueryRow("SELECT * FROM promotes WHERE commentId = ?", id)
	var Temp Promotes
	err := row.Scan(&Temp.Id, &Temp.EndDate, &Temp.PostId, &Temp.CommentId)
	if err != nil {
		if err == sql.ErrNoRows {
			return Temp
		} else {
			fmt.Print(err)
		}
	}
	return Temp
}

func GetAccessToken(id int) AccessTokens {
	row := db.QueryRow("SELECT * FROM access_tokens WHERE id = ?", id)
	var Temp AccessTokens
	err := row.Scan(&Temp.Id, &Temp.Token, &Temp.Date, &Temp.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return Temp
		} else {
			fmt.Print(err)
		}
	}
	return Temp
}

func GetAccessTokenByUser(id int) AccessTokens {
	row := db.QueryRow("SELECT * FROM access_tokens WHERE userId = ?", id)
	var Temp AccessTokens
	err := row.Scan(&Temp.Id, &Temp.Token, &Temp.Date, &Temp.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return Temp
		} else {
			fmt.Print(err)
		}
	}
	return Temp
}

func GetAccessTokenByToken(token string) AccessTokens {
	row := db.QueryRow("SELECT * FROM access_tokens WHERE token = ?", token)
	var Temp AccessTokens
	err := row.Scan(&Temp.Id, &Temp.Token, &Temp.Date, &Temp.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return Temp
		} else {
			fmt.Print(err)
		}
	}
	return Temp
}

// func GetLog(id int) Logs {
// 	row := db.QueryRow("SELECT * FROM logs WHERE id = ?", id)
// 	var Temp Logs
// 	err := row.Scan(&Temp.Id, &Temp.Subject, &Temp.LogsType, &Temp.Action, &Temp.CreationDate, &Temp.UserId)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return Temp
// 		} else {
// 			fmt.Print(err)
// 		}
// 	}
// 	return Temp
// }

// func GetLogByType(logsType int) []Logs {
// 	request, err := db.Prepare("SELECT id FROM logs WHERE type = $1")
// 	if err != nil {
// 		fmt.Print(err)
// 	}
// 	rows, err := request.Query(logsType)
// 	if err != nil {
// 		fmt.Print(err)
// 	}
// 	var TempTable []Logs
// 	for rows.Next() {
// 		var id int
// 		err := rows.Scan(&id)
// 		if err != nil {
// 			if err == sql.ErrNoRows {
// 				return TempTable
// 			} else {
// 				fmt.Print(err)
// 			}
// 		}
// 		TempTable = append(TempTable, GetLog(id))
// 	}
// 	return TempTable
// }

// func GetLogByDate(date string) []Logs {
// 	request, err := db.Prepare("SELECT id FROM logs WHERE creationDate = $1")
// 	if err != nil {
// 		fmt.Print(err)
// 	}
// 	rows, err := request.Query(date)
// 	if err != nil {
// 		fmt.Print(err)
// 	}
// 	var TempTable []Logs
// 	for rows.Next() {
// 		var id int
// 		err := rows.Scan(&id)
// 		if err != nil {
// 			if err == sql.ErrNoRows {
// 				return TempTable
// 			} else {
// 				fmt.Print(err)
// 			}
// 		}
// 		TempTable = append(TempTable, GetLog(id))
// 	}
// 	return TempTable
// }

func GetTicket(id int) Tickets {
	row := db.QueryRow("SELECT * FROM tickets WHERE id = ?", id)
	var Temp Tickets
	var TempId int
	err := row.Scan(&Temp.Id, &Temp.Title, &Temp.Subject, &Temp.CreationDate, &Temp.Status, &TempId)
	Temp.Creator = GetUserInfo(TempId)
	if err != nil {
		if err == sql.ErrNoRows {
			return Temp
		} else {
			fmt.Print(err)
		}
	}
	return Temp
}

func GetTicketMessage(id int) TicketsMessages {
	row := db.QueryRow("SELECT * FROM ticket_messages WHERE id = ?", id)
	var Temp TicketsMessages
	err := row.Scan(&Temp.Id, &Temp.Comment, &Temp.UserId, &Temp.TicketId)
	Temp.Creator = GetUserInfo(Temp.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return Temp
		} else {
			fmt.Print(err)
		}
	}
	return Temp
}

func GetTicketMessages(id int) []TicketsMessages {
	request, err := db.Prepare("SELECT id FROM ticket_messages WHERE ticketId = $1")
	if err != nil {
		fmt.Print(err)
	}
	rows, err := request.Query(id)
	if err != nil {
		fmt.Print(err)
	}
	var TempTable []TicketsMessages
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				return TempTable
			} else {
				fmt.Print(err)
			}
		}
		TempTable = append(TempTable, GetTicketMessage(id))
	}
	return TempTable
}

func GetTickets() []Tickets {
	request, err := db.Prepare("SELECT id FROM tickets ORDER BY status")
	if err != nil {
		fmt.Print(err)
	}
	rows, err := request.Query()
	if err != nil {
		fmt.Print(err)
	}
	var TempTable []Tickets
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				return TempTable
			} else {
				fmt.Print(err)
			}
		}
		TempTable = append(TempTable, GetTicket(id))
	}
	return TempTable
}

func GetUserTickets(userId int) []Tickets {
	request, err := db.Prepare("SELECT id FROM tickets WHERE userId = $1")
	if err != nil {
		fmt.Print(err)
	}
	rows, err := request.Query(userId)
	if err != nil {
		fmt.Print(err)
	}
	var TempTable []Tickets
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				return TempTable
			} else {
				fmt.Print(err)
			}
		}
		TempTable = append(TempTable, GetTicket(id))
	}
	return TempTable
}

func InsertAccessToken(token string, userId int) {
	_, err := db.Exec("INSERT INTO access_tokens (token,date,userId) VALUES (?,?,?)", token, time.Now().Add(24*time.Hour).Format("02-01-2006 15:04:05"), userId)
	if err != nil {
		fmt.Print(err)
	}
}

func InsertCategorieUserRank(userId, categoryId int) {
	_, err := db.Exec("INSERT INTO categories_users_ranks (userId,categoryId) VALUES (?,?)", userId, categoryId)
	if err != nil {
		fmt.Print(err)
	}
}

func InsertLog(subject, action string, Type, userId int) {
	_, err := db.Exec("INSERT INTO logs (subject, action,creationDate,type, userId) VALUES (?,?,?,?,?)", subject, action, time.Now().Format("02-01-2006 15:04:05"), Type, userId)
	if err != nil {
		fmt.Print(err)
	}
}

func InsertNotification(subject, action string, Type, fromId, toId int) {
	_, err := db.Exec("INSERT INTO notifications (subject,action,type,fromId,toId) VALUES (?,?,?,?,?)", subject, action, Type, fromId, toId)
	if err != nil {
		fmt.Print(err)
	}
}

func InsertPromote(endDate string, PostId, CommentId, CategoryId int) {
	_, err := db.Exec("INSERT INTO promotes (endDate,postId,commentId,categoryId) VALUES (?,?,?,?)", endDate, PostId, CommentId, CategoryId)
	if err != nil {
		fmt.Print(err)
	}
}

func InsertReaction(title, icon string) {
	_, err := db.Exec("INSERT INTO reactions (title,icon) VALUES (?,?)", title, icon)
	if err != nil {
		fmt.Print(err)
	}
}

func InsertTicketMessage(comment string, userId, ticketId int) {
	_, err := db.Exec("INSERT INTO ticket_messages (comment,userId,ticketId) VALUES (?,?,?)", comment, userId, ticketId)
	if err != nil {
		fmt.Print(err)
	}
}

func InsertTicket(title, subject string, status int, userId int) {
	_, err := db.Exec("INSERT INTO tickets (title,subject,creationDate, status, userId) VALUES (?,?,?,?,?)", title, subject, time.Now().Format("02-01-2006 15:04:05"), status, userId)
	if err != nil {
		fmt.Print(err)
	}

}

func CloseTicket(id int) {
	request := "UPDATE tickets SET status = 1 WHERE id = $1"
	_, err = db.Exec(request, id)
	if err != nil {
		fmt.Print("error when update")
	}
}

func InsertUser(username, email, password, avatar string) (string, bool) { // to do avec verif email bonne semantique, username pas vide, mdp min len 8 , puis existe deja ou pas pour email et username
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	var error string
	if !re.MatchString(email) {
		return "L'email rentré est invalide", false
	}
	if len(username) < 4 {
		return "Pseudo invalide ou trop court", false
	}
	if len(password) < 8 {
		return "Mot de passe invalide ou trop court", false
	}

	// verification que email et username n'existe pas dans la db, retour de la reponse à insérer dans le html en fonction du résultat.
	emailExist, usernameExist := CheckEmailExist(email), CheckusernameExist(username)
	if !emailExist && !usernameExist {
		password, err = HashPassword(password)
		if err != nil {
			fmt.Print(err)
		}
		if avatar == "" {
			rand.Seed(time.Now().UnixNano())
			avatar = strconv.Itoa(rand.Intn(50-1) + 1)
		}
		_, err := db.Exec("INSERT INTO users (username, email, password, avatar,registerDate) VALUES ($1,$2,$3,$4,$5)", username, email, password, avatar, time.Now().Format("02-01-2006 15:04:05"))
		if err != nil {
			fmt.Print(err)
		}
		return "Inscription reussite", true
		//success = "Votre inscription a bien été effectué! Vous allez être redirigé vers la page de connexion"
	} else if emailExist && !usernameExist {
		error = "Email has already been taken, please retry"
	} else if usernameExist && !emailExist {
		error = "username has already been taken, please retry"
	} else {
		error = "username and email has already been takens,please retry"
	}
	return error, false
}

func InsertUserBadge(badgeId, userId int) {
	if !CheckUserBadgeExist(userId, badgeId) {
		_, err := db.Exec("INSERT INTO users_badges (badgeId,userId) VALUES (?,?)", badgeId, userId)
		if err != nil {
			fmt.Print(err)
		}
	}
}

func RemoveAccessToken(token string) {
	_, err := db.Exec("DELETE FROM access_tokens WHERE token = ?", token)
	if err != nil {
		fmt.Print(err)
	}
}

func ModifyEmail(value string, id int) {
	request := "UPDATE users SET email = $1 WHERE id = $2"
	_, err = db.Exec(request, value, id)
	if err != nil {
		fmt.Print("error when update: ", err)
	}
}

func DeleteUserAll(id int) {
	sqlStatement := `
	DELETE FROM access_tokens
	WHERE userId = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}
	sqlStatement = `
	DELETE FROM categories_users_ranks
	WHERE userId = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}

	sqlStatement = `
	UPDATE comments SET comment = "Commentaire Supprimé"
	WHERE parentId = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}

	sqlStatement = `
	DELETE FROM comments
	WHERE userId = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}

	sqlStatement = `
	DELETE FROM notifications
	WHERE toId = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}

	sqlStatement = `
	DELETE FROM posts
	WHERE userId = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}

	sqlStatement = `
	DELETE FROM ticket_messages
	WHERE userId = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}

	sqlStatement = `
	DELETE FROM tickets
	WHERE userId = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}

	sqlStatement = `
	DELETE FROM users
	WHERE id = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}

	sqlStatement = `
	DELETE FROM users_badges
	WHERE userId = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}

	sqlStatement = `
	DELETE FROM votes
	WHERE userId = $1;`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}

}

func AddCategoryRank(id int, idCategory int) {
	_, err := db.Exec("DELETE FROM categories_users_ranks WHERE userId = ? AND categoryId = ?", id, idCategory)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO categories_users_ranks (userId,categoryId) VALUES (?,?)", id, idCategory)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteCategoryRank(id int, idCategory int) {
	_, err := db.Exec("DELETE FROM categories_users_ranks WHERE userId = ? AND categoryId = ?", id, idCategory)
	if err != nil {
		log.Fatal(err)
	}
}

func ModifyPassword(value string, id int) {
	request := "UPDATE users SET password = $1 WHERE id = $2"
	value, err := HashPassword(value)
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec(request, value, id)
	if err != nil {
		fmt.Print("error when update: ", err)
	}
}

func ModifyUsername(value string, id int) {
	request := "UPDATE users SET username = $1 WHERE id = $2"
	_, err = db.Exec(request, value, id)
	if err != nil {
		fmt.Print("error when update: ", err)
	}
}

func DeleteProfil(userId int) {
	_, err = db.Exec("UPDATE users SET username = 'DELETED ACCOUNT',email = 'deleted', ban = 3 WHERE id = $1", userId)
	if err != nil {
		fmt.Print(err)
	}

}

func ModifyRank(userId int, roleId int) {
	if roleId > 0 && roleId <= 3 {
		_, err := db.Exec("UPDATE users SET rank = $1 WHERE id = $2", roleId, userId)
		if err != nil {
			fmt.Print(err)
		}
	} else {
		return
	}

}

func ModifyBan(value int, id int) {
	request := "UPDATE users SET ban = $1 WHERE id = $2"
	_, err = db.Exec(request, value, id)
	if err != nil {
		fmt.Print("error when update: ", err)
	}
}

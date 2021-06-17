package main

import (
	"database/sql"
	"fmt"
	"time"
)

type Post struct {
	Id               int
	Title            string
	Subject          string
	IdCreator        int
	BadgesCreator    []Badges
	RegisterDate     string
	UserName         string
	UserRank         string
	Avatar           string
	CreationDate     string
	Status           int
	CategoryId       int
	CategoryTitle    string
	UserId           int
	Comment          []Comment
	NumberOfComments int
	Votes            int
	LikedByUser      bool
	DislikedByUser   bool
}

func GetPost(id int) Post {
	row := db.QueryRow("SELECT * FROM posts WHERE id = ?", id)
	var Temp Post
	err := row.Scan(&Temp.Id, &Temp.Title, &Temp.Subject, &Temp.CreationDate, &Temp.Status, &Temp.CategoryId, &Temp.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return Temp
		} else {
			fmt.Println(err)
		}
	}
	Temp.Avatar = NumberToPpIcon(GetUser(Temp.UserId).Avatar)
	return Temp
}

func GetPosts() []Post {
	rows, err := db.Query("SELECT id FROM posts")
	if err != nil {
		fmt.Print(err)
	}
	var TempTable []Post
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
		TempTable = append(TempTable, GetPost(id))
	}
	return TempTable
}

func GetPostsByUser(id int) []Post {
	request, err := db.Prepare("SELECT id FROM posts WHERE userId = $1")
	if err != nil {
		fmt.Print(err)
	}
	rows, err := request.Query(id)
	if err != nil {
		fmt.Print(err)
	}
	var TempTable []Post
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
		TempTable = append(TempTable, GetPost(id))
	}
	return TempTable
}

func GetPostsByCategory(id int) []Post {
	request, err := db.Prepare("SELECT id FROM posts WHERE categoryId = $1 AND status = 0")
	if err != nil {
		fmt.Print(err)
	}
	rows, err := request.Query(id)
	if err != nil {
		fmt.Print(err)
	}
	var TempTable []Post
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
		TempTable = append(TempTable, GetPost(id))
	}
	return TempTable
}

func InsertPost(title, subject string, status, categoryId, userId int) {
	_, err := db.Exec("INSERT INTO posts (title,subject,creationDate,status,categoryId,userId) VALUES (?,?,?,?,?,?)", title, subject, time.Now().Format("02-01-2006 15:04:05"), status, categoryId, userId)
	if err != nil {
		fmt.Print(err)
	}
}

func CheckPostExist(id int) bool {
	row := db.QueryRow("SELECT title FROM posts WHERE id= ?", id) // db correspond a la database ouverte dans db.go
	temp := ""
	err = row.Scan(&temp)
	return temp != ""
}

func CheckPostOwnedByUser(postid, userid int) bool {
	row := db.QueryRow("SELECT id FROM posts WHERE userId= ? AND id = ?", userid, postid) // db correspond a la database ouverte dans db.go
	temp := ""
	err = row.Scan(&temp)
	return temp != ""
}

func CheckCommentOwnedByUser(commentid, userid int) bool {
	row := db.QueryRow("SELECT id FROM comments WHERE userId= $1 AND id = $2", userid, commentid) // db correspond a la database ouverte dans db.go
	temp := ""
	err = row.Scan(&temp)
	return temp != ""
}

func DeletePost(postId, userId int) {

	_, err = db.Exec("UPDATE posts SET status = 1 WHERE id = $1 AND userId = $2", postId, userId)

	if err != nil {
		fmt.Println(err)
	}

}

func DeletePostModo(postId int) {
	_, err = db.Exec("UPDATE posts SET status = 1 WHERE id = $1", postId)
	if err != nil {
		fmt.Println(err)
	}
}

func ModifyPost(id int, subject, justification string) {
	_, err = db.Exec("UPDATE posts SET subject = $1 WHERE id = $2", subject+" "+justification, id)
	if err != nil {
		fmt.Print("error when update: ", err)
	}

}

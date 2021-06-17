package main

import (
	"database/sql"
	"fmt"
	"time"
)

type Comment struct {
	Id                 int
	Comment            string
	CreationDate       string
	RegisterDate       string
	UserId             int
	Status             int
	BadgesCreator      []Badges
	UserName           string
	UserRank           string
	Avatar             string
	PostId             int
	ParentId           int
	ParentComment      string
	ParentUserId       int
	ParentCreationDate string
	Votes              int
	LikedByUser        bool
	DislikedByUser     bool
	SubComment         *Comment
}

func GetComment(id int) Comment {
	row := db.QueryRow("SELECT * FROM comments WHERE id = ?", id)
	var Temp Comment
	err := row.Scan(&Temp.Id, &Temp.Comment, &Temp.CreationDate, &Temp.UserId, &Temp.PostId, &Temp.ParentId, &Temp.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return Temp
		} else {
			fmt.Print(err)
		}
	}
	return Temp
}

func GetCommentsForPost(id int) []Comment {
	request, err := db.Prepare("SELECT id FROM comments WHERE postId = $1 AND STATUS = 0")
	if err != nil {
		fmt.Print(err)
	}
	rows, err := request.Query(id)
	if err != nil {
		fmt.Print(err)
	}
	var TempTable []Comment
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
		TempTable = append(TempTable, GetComment(id))
	}
	return TempTable
}

func GetComments() []Comment {
	rows, err := db.Query("SELECT id FROM comments")
	if err != nil {
		fmt.Print(err)
	}
	var TempTable []Comment
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
		TempTable = append(TempTable, GetComment(id))
	}
	return TempTable
}

func GetCommentsByUser(id int) []Comment {
	request, err := db.Prepare("SELECT id FROM comments WHERE userId = $1")
	if err != nil {
		fmt.Print(err)
	}
	rows, err := request.Query(id)
	if err != nil {
		fmt.Print(err)
	}
	var TempTable []Comment
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
		TempTable = append(TempTable, GetComment(id))
	}
	return TempTable
}

func InsertComment(comment string, userId, postId, parentId int) {
	_, err := db.Exec("INSERT INTO comments (comment,creationDate,userId,postId,parentId) VALUES (?,?,?,?,?)", comment, time.Now().Format("02-01-2006 15:04:05"), userId, postId, parentId)
	if err != nil {
		fmt.Print(err)
	}
}

func DeleteComment(postId, userId int) {
	_, err = db.Exec("UPDATE posts SET status = 1 WHERE id = $1 AND userId = $2", postId, userId)
	if err != nil {
		fmt.Println(err)
	}
}

func ModifyComment(id int, justification string) {
	_, err = db.Exec("UPDATE comments SET comment = $1 WHERE id = $2", justification, id)
	if err != nil {
		fmt.Print("error when update: ", err)
	}
}

func DeleteCommentModo(id int) {
	_, err = db.Exec("UPDATE comments SET status = 1 WHERE id = $1", id)
	if err != nil {
		fmt.Println(err)
	}
}

func CheckCommentInPostExist(commentId int) bool {
	return CheckPostExist(GetComment(commentId).PostId) // db correspond a la database ouverte dans db.go
}

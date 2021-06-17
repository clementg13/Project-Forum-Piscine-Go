package main

import (
	"database/sql"
	"fmt"
)

type Votes struct {
	Id        int
	VoteType  int
	PostId    int
	CommentId int
	UserId    int
}

func GetVote(id int) Votes {
	row := db.QueryRow("SELECT * FROM votes WHERE id = ?", id)
	var Temp Votes
	err := row.Scan(&Temp.Id, &Temp.VoteType, &Temp.PostId, &Temp.CommentId, &Temp.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return Temp
		} else {
			fmt.Print(err)
		}
	}
	return Temp
}

func GetVotesByPost(id int) []Votes {
	request, err := db.Prepare("SELECT id FROM votes WHERE postId = $1")
	if err != nil {
		fmt.Print(err)
	}
	rows, err := request.Query(id)
	if err != nil {
		fmt.Print(err)
	}
	var TempTable []Votes
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
		TempTable = append(TempTable, GetVote(id))
	}
	return TempTable
}

func GetVotesByComment(id int) []Votes {
	request, err := db.Prepare("SELECT id FROM votes WHERE commentId = $1")
	if err != nil {
		fmt.Print(err)
	}
	rows, err := request.Query(id)
	if err != nil {
		fmt.Print(err)
	}
	var TempTable []Votes
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
		TempTable = append(TempTable, GetVote(id))
	}
	return TempTable
}

func GetVotesByUser(id int) []Votes {
	request, err := db.Prepare("SELECT id FROM votes WHERE userId = $1")
	if err != nil {
		fmt.Print(err)
	}
	rows, err := request.Query(id)
	if err != nil {
		fmt.Print(err)
	}
	var TempTable []Votes
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
		TempTable = append(TempTable, GetVote(id))
	}
	return TempTable
}

func InsertVote(Type, postId, commentId, userId int) {
	_, err := db.Exec("INSERT INTO votes (Type, postId, commentId, userId ) VALUES (?,?,?,?)", Type, postId, commentId, userId)
	if err != nil {
		fmt.Print(err)
	}
}

// votes counter, return the score with the likes and dislikes
// with an []votes from a post or comment.
func CountVotes(votesArray []Votes) int {
	var count int
	for _, vote := range votesArray {
		if vote.VoteType == 1 {
			count++
		} else {
			count--
		}
	}
	return count
}

func DeleteVote(userId, postId, commentId int) {
	_, err := db.Exec("DELETE FROM votes WHERE postId = ? OR commentId = ? AND userId = ?;", postId, commentId, userId)
	if err != nil {
		fmt.Print(err)
	}
}

func CheckLikePostByUser(userid, postid int) bool {
	var Type int
	sqlStatement := `SELECT type FROM votes WHERE userId=$1 AND postId=$2`
	row := db.QueryRow(sqlStatement, userid, postid)
	err := row.Scan(&Type)
	if err != nil {
		if err == sql.ErrNoRows {
			// fmt.Println("Zero rows found")
		} else {
			panic(err)
		}
	}
	return Type == 1
}

func CheckDislikePostByUser(userid, postid int) bool {
	var Type int
	sqlStatement := `SELECT type FROM votes WHERE userId=$1 AND postId=$2`
	row := db.QueryRow(sqlStatement, userid, postid)
	err := row.Scan(&Type)
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			panic(err)
		}
	}
	return Type == 2
}

func CheckLikeCommentByUser(userid, postid int) bool {
	var Type int
	sqlStatement := `SELECT type FROM votes WHERE userId=$1 AND commentId=$2`
	row := db.QueryRow(sqlStatement, userid, postid)
	err := row.Scan(&Type)
	if err != nil {
		if err == sql.ErrNoRows {
			// fmt.Println("Zero rows found")
		} else {
			panic(err)
		}
	}
	return Type == 1
}

func CheckDislikeCommentByUser(userid, postid int) bool {
	var Type int
	sqlStatement := `SELECT type FROM votes WHERE userId=$1 AND commentId=$2`
	row := db.QueryRow(sqlStatement, userid, postid)
	err := row.Scan(&Type)
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			panic(err)
		}
	}
	return Type == 2
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

type UsersBadges struct {
	Id      int
	BadgeId int
	UserId  int
}

type Badges struct {
	Id    int
	Name  string
	Icon  string
	Type  int
	Value int
}

// Check if the user do complete some badges conditions and give it
// to him if it is.
func BadgesEligibilityChecker(userId int) {
	// user = GetUser(userId)
	if len(GetVotesByUser(userId)) >= 10 {
		InsertUserBadge(GetBadgeId("+10votes done"), userId)
	}

	if len(GetPostsByUser(userId)) >= 1 {
		InsertUserBadge(GetBadgeId("1st post created"), userId)
	}

	if len(GetPostsByUser(userId)) >= 5 {
		InsertUserBadge(GetBadgeId("+5 posts created"), userId)
	}

	usersPostsArray := GetPostsByUser(userId)
	count := 1
	for _, y := range usersPostsArray {
		if count == 1 {
			if len(GetVotesByPost(y.Id)) >= 10 {
				count++
				InsertUserBadge(GetBadgeId("+10votes on your post"), userId)
			}
		} else {
			break
		}
	}
}

func InsertBadge(name, icon string, Type, value int) {
	_, err := db.Exec("INSERT INTO badges (name,icon,type,value) VALUES (?,?,?,?)", name, icon, Type, value)
	if err != nil {
		fmt.Print(err)
	}
}

// return true if the badge do exist for the user.
func CheckUserBadgeExist(userId, badgeId int) bool {
	row := db.QueryRow("SELECT id FROM users_badges WHERE userId = $1 AND badgeId = $2", userId, badgeId) // db correspond a la database ouverte dans db.go
	temp := ""
	err := row.Scan(&temp)
	if err != nil {
		fmt.Print(err)
	}
	return temp != ""
}

func GetBadgeId(name string) int {
	var temp string
	row := db.QueryRow("SELECT id FROM badges WHERE name = ?", name) // db correspond a la database ouverte dans db.go
	err := row.Scan(&temp)
	if err != nil {
		fmt.Print(err)
	}
	result, err := strconv.Atoi(temp)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func GetUserBadges(userId int) []Badges {
	var badges_array []Badges
	request, err := db.Prepare("SELECT * FROM users_badges WHERE userId = $1 LIMIT 3")
	if err != nil {
		fmt.Print(err)
	}
	rows, err := request.Query(userId)
	if err != nil {
		fmt.Print(err)
	}
	for rows.Next() {
		var id, badgeId, userId int
		err := rows.Scan(&id, &badgeId, &userId)

		if err != nil {
			if err == sql.ErrNoRows {
				return badges_array
			} else {

				fmt.Print(err)
			}
		}
		badge := GetBadgeById(badgeId)

		badges_array = append(badges_array, badge)
	}
	return badges_array
}

func GetBadgeById(badgeId int) Badges {
	var badge Badges
	request, err := db.Prepare("SELECT * FROM badges WHERE id = $1")
	if err != nil {
		fmt.Print(err)
	}
	rows, err := request.Query(badgeId)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var Id, Type, Value int
		var Name, Icon string
		err = rows.Scan(&Id, &Name, &Icon, &Type, &Value)
		if err != nil {
			fmt.Println("erreur=>", err)
		}
		badge.Icon = Icon
		badge.Id = Id
		badge.Value = Value
		badge.Name = Name
		badge.Type = Type
	}
	return badge
}
func GetBadges() []Badges {
	request, err := db.Prepare("SELECT id FROM badges")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := request.Query()
	if err != nil {
		log.Fatal(err)
	}
	var TempTable []Badges
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
		TempTable = append(TempTable, GetBadgeById(id))
	}
	return TempTable
}

func AddBadge(value int, id int) {
	request := "INSERT INTO users_badges(badgeId, userId) VALUES ($1,$2)"
	_, err = db.Exec(request, value, id)
	if err != nil {
		fmt.Print("error when update: ", err)
	}
}

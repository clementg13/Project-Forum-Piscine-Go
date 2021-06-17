package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Category struct {
	Id           int
	Title        string
	Description  string
	Icon         string
	CreationDate string
	ParentId     int
	SubCategory  []SubCategory
	IsAuthorized bool
}

type SubCategory struct {
	Id           int
	Title        string
	Description  string
	Icon         string
	CreationDate string
	NumberOfPost int
	ParentId     int
	IsAuthorized bool
}

func GetCategory(id int) Category {
	row := db.QueryRow("SELECT * FROM categories WHERE id = ?", id)
	var Temp Category
	err := row.Scan(&Temp.Id, &Temp.Title, &Temp.Description, &Temp.Icon, &Temp.CreationDate, &Temp.ParentId)
	if Temp.ParentId == 0 {
		Temp.SubCategory = GetSubCategories(Temp.Id)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return Temp
		} else {
			fmt.Print(err)
		}
	}
	return Temp
}

func GetCategories() []Category {
	rows, err := db.Query("SELECT id FROM categories WHERE parentId = 0")
	if err != nil {
		fmt.Print(err)
	}
	var TempTable []Category
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
		TempTable = append(TempTable, GetCategory(id))
	}
	return TempTable
}

func GetAuthorizedCategories(userId int) []Category {
	rows, err := db.Query("SELECT id FROM categories_users_ranks WHERE userId = $1", userId)
	if err != nil {
		fmt.Print(err)
	}
	var TempTable []Category
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
		TempTable = append(TempTable, GetCategory(id))
	}
	return TempTable
}

func GetAllCategories() []Category {
	rows, err := db.Query("SELECT id FROM categories")
	if err != nil {
		log.Fatal(err)
	}
	var TempTable []Category
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
		TempTable = append(TempTable, GetCategory(id))
	}
	return TempTable
}

func InsertCategory(title, description, icon string, parentId int) {
	_, err := db.Exec("INSERT INTO categories (title, description,icon,creationDate,parentId ) VALUES (?,?,?,?,?)", title, description, icon, time.Now().Format("02-01-2006 15:04:05"), parentId)
	if err != nil {
		fmt.Print(err)
	}
}

func GetSubCategory(id int) SubCategory {
	row := db.QueryRow("SELECT * FROM categories WHERE id = ?", id)
	var Temp SubCategory
	err := row.Scan(&Temp.Id, &Temp.Title, &Temp.Description, &Temp.Icon, &Temp.CreationDate, &Temp.ParentId)
	if err != nil {
		if err == sql.ErrNoRows {
			return Temp
		} else {
			fmt.Print(err)
		}
	}
	Temp.NumberOfPost = GetNumberOfPosts(Temp.Id)
	return Temp
}

func GetSubCategories(id int) []SubCategory {
	request, err := db.Prepare("SELECT id FROM categories WHERE parentId = $1")
	if err != nil {
		fmt.Print(err)
	}
	rows, err := request.Query(id)
	if err != nil {
		fmt.Print(err)
	}
	var TempTable []SubCategory
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
		TempTable = append(TempTable, GetSubCategory(id))
	}
	return TempTable
}

func CheckCategoryExist(id int) bool {
	row := db.QueryRow("SELECT title FROM categories WHERE id= ?", id) // db correspond a la database ouverte dans db.go
	temp := ""
	err = row.Scan(&temp)
	if err != nil {
		fmt.Print(err)
	}
	return temp != ""
}

// return true if the user has the category authorization
func CheckCategoryAuthorization(userId, categoryId int) bool {
	row := db.QueryRow("SELECT id FROM categories_users_ranks WHERE userId = $1 AND categoryId = $2", userId, categoryId) // db correspond a la database ouverte dans db.go
	temp := ""
	err = row.Scan(&temp)
	return temp != ""
}

func CheckIsPrimaryCategory(categoryid int) bool {
	row := db.QueryRow("SELECT id FROM categories WHERE id = $1 and parentId = 0", categoryid) // db correspond a la database ouverte dans db.go
	temp := ""
	err = row.Scan(&temp)
	return temp != ""
}

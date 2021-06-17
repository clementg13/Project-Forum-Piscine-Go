package main

import (
	"database/sql"
	"log"
)

// return the number of posts
func StatsPostNumber() int {
	var stat int
	row := db.QueryRow("SELECT COUNT(*) FROM posts")
	err := row.Scan(&stat)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		} else {
			log.Fatal(err)
		}
	}
	return stat
}

// return the number of posts for a given time interval
func StatsDatePostNumber(first string, second string) int {
	firstDay := first[8:]
	firstMonth := first[5:7]
	firstYear := first[:4]

	if second == "" {
		second = first
	}
	secondDay := second[8:]
	secondMonth := second[5:7]
	secondYear := second[:4]

	var count int
	request, err := db.Prepare("SELECT creationDate FROM posts")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := request.Query()
	if err != nil {
		log.Fatal(err)
	}
	var tempDate string
	for rows.Next() {
		err := rows.Scan(&tempDate)
		if err != nil {
			if err == sql.ErrNoRows {
				return count
			} else {
				log.Fatal(err)
			}
		}
		dbDay := tempDate[:2]
		dbMonth := tempDate[3:5]
		dbYears := tempDate[6:10]
		if dbDay >= firstDay && dbDay <= secondDay {
			if dbMonth >= firstMonth && dbMonth <= secondMonth {
				if dbYears >= firstYear && dbYears <= secondYear {
					count++
				}
			}
		}
	}
	return count
}

// return the number of comments
func StatsCommentNumber() int {
	var stat int
	row := db.QueryRow("SELECT COUNT(*) FROM comments")
	err := row.Scan(&stat)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		} else {
			log.Fatal(err)
		}
	}
	return stat
}

// return the number of comment for a given time interval
func StatsDateCommentNumber(first string, second string) int {
	firstDay := first[8:]
	firstMonth := first[5:7]
	firstYear := first[:4]

	if second == "" {
		second = first
	}
	secondDay := second[8:]
	secondMonth := second[5:7]
	secondYear := second[:4]

	var count int
	request, err := db.Prepare("SELECT creationDate FROM comments")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := request.Query()
	if err != nil {
		log.Fatal(err)
	}
	var tempDate string
	for rows.Next() {
		err := rows.Scan(&tempDate)
		if err != nil {
			if err == sql.ErrNoRows {
				return count
			} else {
				log.Fatal(err)
			}
		}
		dbDay := tempDate[:2]
		dbMonth := tempDate[3:5]
		dbYears := tempDate[6:10]
		if dbDay >= firstDay && dbDay <= secondDay {
			if dbMonth >= firstMonth && dbMonth <= secondMonth {
				if dbYears >= firstYear && dbYears <= secondYear {
					count++
				}
			}
		}
	}
	return count
}

// return the number of categories
func StatsCategoryNumber() int {
	var stat int
	row := db.QueryRow("SELECT COUNT(*) FROM categories WHERE parentId = 0")
	err := row.Scan(&stat)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		} else {
			log.Fatal(err)
		}
	}
	return stat
}

// return the number of catgegories for a given time interval
func StatsDateCategoryNumber(first string, second string) int {
	firstDay := first[8:]
	firstMonth := first[5:7]
	firstYear := first[:4]

	if second == "" {
		second = first
	}
	secondDay := second[8:]
	secondMonth := second[5:7]
	secondYear := second[:4]

	var count int
	request, err := db.Prepare("SELECT creationDate FROM categories WHERE parentId = 0")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := request.Query()
	if err != nil {
		log.Fatal(err)
	}
	var tempDate string
	for rows.Next() {
		err := rows.Scan(&tempDate)
		if err != nil {
			if err == sql.ErrNoRows {
				return count
			} else {
				log.Fatal(err)
			}
		}
		dbDay := tempDate[:2]
		dbMonth := tempDate[3:5]
		dbYears := tempDate[6:10]
		if dbDay >= firstDay && dbDay <= secondDay {
			if dbMonth >= firstMonth && dbMonth <= secondMonth {
				if dbYears >= firstYear && dbYears <= secondYear {
					count++
				}
			}
		}
	}
	return count
}

// return the number of subcategories
func StatsSubCategoryNumber() int {
	var stat int
	row := db.QueryRow("SELECT COUNT(*) FROM categories WHERE parentId != 0")
	err := row.Scan(&stat)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		} else {
			log.Fatal(err)
		}
	}
	return stat
}

// return the number of sub categories for a given time interval
func StatsDateSubCategoryNumber(first string, second string) int {
	firstDay := first[8:]
	firstMonth := first[5:7]
	firstYear := first[:4]

	if second == "" {
		second = first
	}
	secondDay := second[8:]
	secondMonth := second[5:7]
	secondYear := second[:4]

	var count int
	request, err := db.Prepare("SELECT creationDate FROM categories WHERE parentId != 0")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := request.Query()
	if err != nil {
		log.Fatal(err)
	}
	var tempDate string
	for rows.Next() {
		err := rows.Scan(&tempDate)
		if err != nil {
			if err == sql.ErrNoRows {
				return count
			} else {
				log.Fatal(err)
			}
		}
		dbDay := tempDate[:2]
		dbMonth := tempDate[3:5]
		dbYears := tempDate[6:10]
		if dbDay >= firstDay && dbDay <= secondDay {
			if dbMonth >= firstMonth && dbMonth <= secondMonth {
				if dbYears >= firstYear && dbYears <= secondYear {
					count++
				}
			}
		}
	}
	return count
}

// return the number of votes
func StatsVotes() int {
	var statLike int
	row := db.QueryRow("SELECT COUNT(*) FROM votes WHERE type = 1")
	err := row.Scan(&statLike)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		} else {
			log.Fatal(err)
		}
	}
	var statDislike int
	row = db.QueryRow("SELECT COUNT(*) FROM votes WHERE type = 2")
	err = row.Scan(&statDislike)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		} else {
			log.Fatal(err)
		}
	}

	return statLike - statDislike
}

package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func CheckEmailExist(email string) bool {
	row := db.QueryRow("SELECT email FROM users WHERE email= ?", email) // db correspond a la database ouverte dans db.go
	temp := ""
	err := row.Scan(&temp)
	fmt.Println(temp)
	if err != nil {
		fmt.Println(err)
	}
	return temp != ""
}

func GetHashedPassword(email string) string {
	row := db.QueryRow("SELECT password FROM users WHERE email= ?", email) // db correspond a la database ouverte dans db.go
	temp := ""

	row.Scan(&temp)
	if err != nil {
		fmt.Print(err)
	}

	return temp
}

func CheckusernameExist(username string) bool {
	row := db.QueryRow("SELECT username FROM users WHERE username= ?", username) // db correspond a la database ouverte dans db.go
	temp := ""
	row.Scan(&temp)
	if err != nil {
		fmt.Print(err)
	}
	return temp != ""
}

// check ban => regarde dans la db (table ban) si l'utilisateur (id) est ban
func CheckBan(id int) bool {
	row := db.QueryRow("SELECT ban FROM user WHERE id= ?", id) // db correspond a la database ouverte dans db.go
	temp := ""
	row.Scan(&temp)
	if err != nil {
		fmt.Print(err)
	}
	return temp != "0"
}

func LoginValidation(email string, password string) (bool, int) {
	if CheckPasswordHash(password, GetHashedPassword(email)) {
		row := db.QueryRow("SELECT id FROM users WHERE email = ?", email)
		var id int
		err := row.Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				return false, 0
			} else {
				fmt.Print(err)
			}
		}
		return true, id
	} else {
		return false, 0
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SendRegistrationData(email string, password string, username string) { // envoie des données d'inscription dans la db.
	userData, err := db.Prepare(`INSERT INTO user (email, password,username,profilpicture) VALUES (?,?,?,?)`)
	if err != nil {
		fmt.Print(err)
	}
	rand.Seed(time.Now().UnixNano())
	profilpicture_int := rand.Intn(50-1) + 1                             // génération d'un nombre entre 1 et 50 pour la génération aléatoire de photo de profil par defaut
	_, err = userData.Exec(email, password, username, profilpicture_int) // exécution et envoi des données dans la db.
	if err != nil {
		fmt.Print(err)
	}
}

func CheckToken(w http.ResponseWriter, r *http.Request, token string) int {
	row := db.QueryRow("SELECT * FROM access_tokens WHERE token = ?", token)
	var Temp AccessTokens
	err := row.Scan(&Temp.Id, &Temp.Token, &Temp.Date, &Temp.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		} else {
			fmt.Print(err)
		}
	}
	date := Temp.Date
	dateFormated, err := time.Parse("02-01-2006 15:04:05", date)
	timenow := time.Now().Format("02-01-2006 15:04:05")
	timenowFormated, err := time.Parse("02-01-2006 15:04:05", timenow)
	if err != nil {
		_, err = db.Exec("DELETE FROM access_tokens WHERE token = $1", Temp.Token)
		if err != nil {
			fmt.Print(err)
		}
		DeleteCookie(w, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return 0
	}
	after := dateFormated.After(timenowFormated)
	if after {
		return Temp.UserId
	}
	_, err = db.Exec("DELETE FROM access_tokens WHERE token = $1", token)
	if err != nil {
		fmt.Print(err)
	}
	DeleteCookie(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return 0
}

func GetUserCookieInfo(id int) UserInfo {
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	var Temp User
	var TempUserInfo UserInfo
	err := row.Scan(&Temp.Id, &Temp.Username, &Temp.Email, &Temp.Password, &Temp.Avatar, &Temp.RegisterDate, &Temp.Ban, &Temp.Rank)
	if err != nil {
		if err == sql.ErrNoRows {
			return TempUserInfo
		} else {
			fmt.Print(err)
		}
	}
	TempUserInfo.Id = Temp.Id
	TempUserInfo.Username = Temp.Username
	TempUserInfo.Email = Temp.Email
	TempUserInfo.Ban = Temp.Ban
	TempUserInfo.Avatar = Temp.Avatar
	TempUserInfo.Rank = Temp.Rank
	TempUserInfo.Ban = Temp.Ban
	TempUserInfo.RegisterDate = Temp.RegisterDate
	TempUserInfo.Permissions = GetRank(Temp.Rank)
	return TempUserInfo
}

func CheckEmailIsValid(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(email) {
		return true
	} else {
		return false
	}
}

//SendMail for password recovery with the SMTP
func SendMail(email, newpassword string) {

	// Sender data.
	from := "whalumassistance@gmail.com"
	password := "#admin@1234#"

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	message := []byte(newpassword)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}

var (
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*"
	numberSet      = "0123456789"
	allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
)

func ModifyPasswordRandom(email string) string {
	rand.Seed(time.Now().Unix())
	minSpecialChar := 1
	minNum := 1
	minUpperCase := 1
	passwordLength := 8
	password := generatePassword(passwordLength, minSpecialChar, minNum, minUpperCase)
	fmt.Println(password)

	minSpecialChar = 2
	minNum = 2
	minUpperCase = 2
	passwordLength = 20
	password = generatePassword(passwordLength, minSpecialChar, minNum, minUpperCase)
	ModifyPassword(password, GetUserIdByEmail(email))
	return password
}

func generatePassword(passwordLength, minSpecialChar, minNum, minUpperCase int) string {
	var password strings.Builder

	//Set special character
	for i := 0; i < minSpecialChar; i++ {
		random := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

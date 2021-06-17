package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var templates map[string]*template.Template

// Parse all files in page.gohtml and .common
// for the template execution.
func InitTemplate() error {
	templates = make(map[string]*template.Template)

	files, err := filepath.Glob(filepath.Join("./template", "*.page.gohtml"))
	if err != nil {
		return err
	}
	for _, file := range files {
		name := filepath.Base(file)
		file, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		tmpl, err := template.New(name).Parse(string(file))
		if err != nil {
			return err
		}
		tmpl, err = tmpl.ParseGlob(filepath.Join("./template", "*.layout.gohtml"))
		if err != nil {
			return err
		}
		templates[tmpl.Name()] = tmpl
	}

	return nil

}

// first page of the website
// Take all the categories and subcategories to display them
// with the index.page.gohtml template
// all the Index Page data are collected inside by calling other functions
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	type IndexPage struct {
		TableCategory []Category
		IsConnected   bool
		User          UserInfo
	}
	if r.URL.Path != "/" && r.URL.Path != "/index" && r.URL.Path != "/index.html" {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	var indexStruct IndexPage
	connected, token := CheckCookie(r)
	userId := CheckToken(w, r, token)
	if connected {
		userId := CheckToken(w, r, token)
		CheckCookieValidity(r, userId)
		if userId != 0 {
			indexStruct.User = GetUserCookieInfo(userId)
			BadgesEligibilityChecker(userId)
			if indexStruct.User.Ban > 0 {
				DeleteCookie(w, r)
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			indexStruct.IsConnected = true
		}
		indexStruct.User.Avatar = NumberToPpIcon(GetUser(indexStruct.User.Id).Avatar)
	} else {
		indexStruct.IsConnected = false
	}
	indexStruct.TableCategory = GetCategories()
	if connected {
		for i := range indexStruct.TableCategory {
			indexStruct.TableCategory[i].IsAuthorized = CheckCategoryAuthorization(userId, indexStruct.TableCategory[i].Id)
			for x := range indexStruct.TableCategory[i].SubCategory {
				indexStruct.TableCategory[i].SubCategory[x].IsAuthorized = CheckCategoryAuthorization(userId, indexStruct.TableCategory[i].SubCategory[x].Id)
			}
		}
	}
	err := templates["index.page.gohtml"].Execute(w, indexStruct)
	if err != nil {
		log.Fatalln("template didn't execute: ", err)
	}

}

// Login Page Handler checking the cookie validity if
// user is connected, else, if the login informations are corrects,
// create a new UUID and insert it as a new cookie.
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	type LoginPage struct {
		ErrorMessage   string
		SuccessMessage string
		Email          string
		IsConnected    bool
		User           UserInfo
	}

	var pageStruct LoginPage

	connected, token := CheckCookie(r)

	if connected {
		userId := CheckToken(w, r, token)
		CheckCookieValidity(r, userId)
		if userId != 0 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

	if r.Method != http.MethodPost {
		err = templates["login.page.gohtml"].Execute(w, nil)
		if err != nil {
			log.Fatalln("template didn't execute: ", err)
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		email, password := r.FormValue("email"), r.FormValue("password")
		pageStruct.Email = email
		if email == "" || password == "" {
			pageStruct.ErrorMessage = "Merci de Remplir tout les Champs"
		} else {
			loginIsValid, id := LoginValidation(email, password)
			if loginIsValid {
				pageStruct.SuccessMessage = "Connexion réussie, vous allez être redirigé sur la page d'accueil"
				cookieExist, _ := CheckCookie(r)
				if !cookieExist {
					cookieUuid := CreateUniqueCookie(w)
					InsertAccessToken(cookieUuid, id)
					pageStruct.IsConnected = true
					http.Redirect(w, r, "/?success='connexion réussite'", http.StatusSeeOther)
					return
				} else {
					err, cookieUuid := DeleteCookie(w, r)
					if err == nil {
						InsertAccessToken(cookieUuid, id)
						pageStruct.IsConnected = true
						http.Redirect(w, r, "/?success='connexion réussite'", http.StatusSeeOther)
						return
					} else {
						pageStruct.ErrorMessage = "Erreur lors de la création du cookie 2"
						pageStruct.SuccessMessage = ""
						pageStruct.IsConnected = false
					}
				}
			} else {
				pageStruct.ErrorMessage = "Informations de Connexion Invalide"
			}
		}

		err = templates["login.page.gohtml"].Execute(w, pageStruct)
		if err != nil {
			log.Fatalln("template didn't execute: ", err)
		}
	}

}

// Delete Handler erasing the actual user login cookie
func DisconnectHandler(w http.ResponseWriter, r *http.Request) {
	connected, _ := CheckCookie(r)
	if connected {
		DeleteCookie(w, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// register handler, insert a new user in the database if
// email and pseudo do not exist in the db + email is correct
// and password is more then +8 chars.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	type RegisterPage struct {
		ErrorMessage   string
		SuccessMessage string
		Email          string
		Pseudo         string
		IsConnected    bool
		User           UserInfo
	}

	var pageStruct RegisterPage
	connected, token := CheckCookie(r)
	if connected {
		userId := CheckToken(w, r, token)
		CheckCookieValidity(r, userId)
		if userId != 0 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	if r.Method != http.MethodPost {
		err = templates["register.page.gohtml"].Execute(w, nil)
		if err != nil {
			log.Fatalln("template didn't execute: ", err)
		}
	} else {
		email, pseudo, password := r.FormValue("email"), r.FormValue("pseudo"), r.FormValue("password")
		errorMessage, success := InsertUser(pseudo, email, password, "")
		if success {
			pageStruct.SuccessMessage = "Votre inscription a bien été effectué! Vous allez être redirigé vers la page de connexion"
			http.Redirect(w, r, "/login?success='inscription réussite'", http.StatusSeeOther)
			return
		} else {
			pageStruct.Email = email
			pageStruct.Pseudo = pseudo
			pageStruct.ErrorMessage = errorMessage
		}
		err = templates["register.page.gohtml"].Execute(w, pageStruct)
		if err != nil {
			log.Fatalln("template didn't execute: ", err)
		}
	}

}

// Admin ticket handler is returning the struct
//  with every users tickets informations
func AdminTicketHandler(w http.ResponseWriter, r *http.Request) {
	type AdminPage struct {
		ErrorMessage   string
		SuccessMessage string
		IsConnected    bool
		Page           string
		Tickets        []Tickets
		User           UserInfo
	}
	var pageStruct AdminPage
	connected, token := CheckCookie(r)
	if connected {
		userId := CheckToken(w, r, token)
		if userId != 0 {
			pageStruct.User = GetUserCookieInfo(userId)
			pageStruct.User.Avatar = NumberToPpIcon(pageStruct.User.Avatar)
			pageStruct.IsConnected = true
		}
	} else {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	rank := pageStruct.User.Permissions

	if rank.AdminPanelAccess != 1 {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	if rank.TicketAccess != 1 {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	pageStruct.Page = "ticket"
	pageStruct.Tickets = GetTickets()
	err = templates["admin_tickets.page.gohtml"].Execute(w, pageStruct)
	if err != nil {
		log.Fatalln("template didn't execute: ", err)
	}
}

// show a single ticket conversations witb users messages.
func AdminTicketViewHandler(w http.ResponseWriter, r *http.Request) {
	type AdminPage struct {
		ErrorMessage   string
		SuccessMessage string
		IsConnected    bool
		Page           string
		Ticket         Tickets
		TicketMessages []TicketsMessages
		User           UserInfo
	}
	var pageStruct AdminPage
	connected, token := CheckCookie(r)
	if connected {
		userId := CheckToken(w, r, token)
		if userId != 0 {
			pageStruct.User = GetUserCookieInfo(userId)
			pageStruct.User.Avatar = NumberToPpIcon(pageStruct.User.Avatar)
			pageStruct.IsConnected = true
			BadgesEligibilityChecker(userId)
		}
		pageStruct.User.Avatar = NumberToPpIcon(GetUser(pageStruct.User.Id).Avatar)

	} else {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	rank := pageStruct.User.Permissions

	if rank.AdminPanelAccess != 1 {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	if rank.TicketAccess != 1 {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}

	if ticketId, ok := r.URL.Query()["id"]; ok {
		if rank.TicketAccess != 1 {
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
			return
		}
		ticket, err := strconv.Atoi(ticketId[0])
		if err != nil {
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
			return
		}
		tempTicket := GetTicket(ticket)
		if !reflect.ValueOf(tempTicket).IsZero() {
			if r.Method == http.MethodPost {
				message := r.FormValue("sendticket-message")
				closeMessage := r.FormValue("closeticket-message")
				if len(message) > 1 {
					InsertTicketMessage(message, pageStruct.User.Id, ticket)
				} else if len(closeMessage) > 1 {
					closeMessage = "Fermeture du Ticket, Raison: " + closeMessage
					InsertTicketMessage(closeMessage, pageStruct.User.Id, ticket)
					CloseTicket(ticket)
					http.Redirect(w, r, "/admin_ticketpage?id="+strconv.Itoa(ticket), http.StatusSeeOther)
				}
			}

			pageStruct.Ticket = tempTicket
			pageStruct.Page = "ticketPage"
			pageStruct.TicketMessages = GetTicketMessages(ticket)
			pageStruct.User.Avatar = NumberToPpIcon(GetUser(pageStruct.User.Id).Avatar)
			err = templates["admin_ticketview.page.gohtml"].Execute(w, pageStruct)
			if err != nil {
				log.Fatalln("template didn't execute: ", err)
			}
			return
		} else {
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
			return
		}
	}
	err = templates["404.page.gohtml"].Execute(w, basePage)
	w.WriteHeader(404)
}

// Users infos manager by the admin of the website.
// can ban and modify users infos.
func AdminGestionUsersHandler(w http.ResponseWriter, r *http.Request) {
	type AdminPage struct {
		ErrorMessage   string
		SuccessMessage string
		IsConnected    bool
		Page           string
		User           UserInfo
		Users          []UserInfo
		Permissions    []Ranks
		Badges         []Badges
	}
	var pageStruct AdminPage
	connected, token := CheckCookie(r)
	if connected {
		userId := CheckToken(w, r, token)
		if userId != 0 {
			pageStruct.User = GetUserCookieInfo(userId)
			pageStruct.User.Avatar = NumberToPpIcon(pageStruct.User.Avatar)
			pageStruct.IsConnected = true
		} else {
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
			return
		}
	} else {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	rank := pageStruct.User.Permissions

	if rank.AdminPanelAccess != 1 {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	if r.Method == http.MethodPost {
		idtemp := r.FormValue("id")
		name := r.FormValue("name")
		bantemp := r.FormValue("ban")
		email := r.FormValue("email")
		roletemp := r.FormValue("role")
		badgestemp := r.FormValue("badges")
		var id int
		var ban int
		var role int
		var badges int
		if len(idtemp) > 0 {
			id, err = strconv.Atoi(idtemp)
			if err != nil {
				fmt.Print("ban ", err)
				err = templates["404.page.gohtml"].Execute(w, basePage)
				w.WriteHeader(404)
				return
			}
		}
		if len(bantemp) > 0 {
			ban, err = strconv.Atoi(bantemp)
			if err != nil {
				fmt.Print("ban ", err)
				err = templates["404.page.gohtml"].Execute(w, basePage)
				w.WriteHeader(404)
				return
			}
		}
		if len(roletemp) > 0 {
			role, err = strconv.Atoi(roletemp)
			if err != nil {
				fmt.Print("role ", err)
				w.WriteHeader(404)
				return
			}
		}
		if len(badgestemp) > 0 && len(badgestemp) <= 2 {
			badges, err = strconv.Atoi(badgestemp)
			if err != nil {
				fmt.Print("badges ", err)
				err = templates["404.page.gohtml"].Execute(w, basePage)
				w.WriteHeader(404)
				return
			}
		}
		if id > 0 {
			if len(name) > 1 {
				pageStruct.SuccessMessage = "Action effectué avec succès"
				ModifyUsername(name, id)
			}
			if ban >= 0 && ban != 100 {
				if ban == 69 {
					pageStruct.SuccessMessage = "Action effectué avec succès"
					DeleteUserAll(id)
				}
				if ban == 3 {
					pageStruct.SuccessMessage = "Action effectué avec succès"
					DeleteProfil(id)
				}
				if ban == 1 || ban == 0 {
					pageStruct.SuccessMessage = "Action effectué avec succès"
					ModifyBan(ban, id)
				}
			}
			if badges > 0 {
				pageStruct.SuccessMessage = "Action effectué avec succès"
				AddBadge(badges, id)
			}

			if len(email) > 1 {
				pageStruct.SuccessMessage = "Action effectué avec succès"
				ModifyEmail(email, id)
			}
			if role > 0 {
				pageStruct.SuccessMessage = "Action effectué avec succès"
				ModifyRank(id, role)
			}
		}
		if len(pageStruct.SuccessMessage) < 1 {
			pageStruct.ErrorMessage = "Une erreur c'est produite"
		}

		pageStruct.Page = "gestionusers"
		pageStruct.Users = GetUsersInfo()
		pageStruct.Permissions = GetRanks()
		pageStruct.Badges = GetBadges()
		err = templates["admin_gestionusers.page.gohtml"].Execute(w, pageStruct)
		if err != nil {
			log.Fatalln("template didn't execute: ", err)
		}
		return

	} else {
		pageStruct.Page = "gestionusers"
		pageStruct.Users = GetUsersInfo()
		pageStruct.Permissions = GetRanks()
		pageStruct.Badges = GetBadges()
		err = templates["admin_gestionusers.page.gohtml"].Execute(w, pageStruct)
		if err != nil {
			log.Fatalln("template didn't execute: ", err)
		}
		return
	}
}

// Handler for promotion or unpromotion of the users/modo
// (for the admin only)
func AdminCategoriesPermission(w http.ResponseWriter, r *http.Request) {
	type AdminPage struct {
		ErrorMessage   string
		SuccessMessage string
		IsConnected    bool
		Page           string
		User           UserInfo
		Users          []UserInfo
		UserCategory   []AllCategoriesUserRanks
		Category       []Category
	}
	var pageStruct AdminPage
	connected, token := CheckCookie(r)
	if connected {
		userId := CheckToken(w, r, token)
		if userId != 0 {
			pageStruct.User = GetUserCookieInfo(userId)
			pageStruct.User.Avatar = NumberToPpIcon(pageStruct.User.Avatar)
			pageStruct.IsConnected = true
		} else {
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
			return
		}
	} else {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	rank := pageStruct.User.Permissions

	if rank.AdminPanelAccess != 1 {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	if rank.ModifyRank != 1 {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		idtemp := r.FormValue("id")
		deleteaddtemp := r.FormValue("deleteadd")
		categoryFormTemp := r.FormValue("category")
		var id int
		var deleteadd int
		var categoryId int
		if len(idtemp) > 0 {
			id, err = strconv.Atoi(idtemp)
			if err != nil {
				fmt.Print("ban ", err)
				err = templates["404.page.gohtml"].Execute(w, basePage)
				w.WriteHeader(404)
				return
			}
		}
		if len(deleteaddtemp) > 0 {
			deleteadd, err = strconv.Atoi(deleteaddtemp)
			if err != nil {
				fmt.Print("ban ", err)
				err = templates["404.page.gohtml"].Execute(w, basePage)
				w.WriteHeader(404)
				return
			}
		}
		if len(categoryFormTemp) > 0 {
			categoryId, err = strconv.Atoi(categoryFormTemp)
			if err != nil {
				fmt.Print("ban ", err)
				err = templates["404.page.gohtml"].Execute(w, basePage)
				w.WriteHeader(404)
				return
			}
		}
		if id > 0 {
			if deleteadd == 1 {
				if CheckIsPrimaryCategory(categoryId) {
					subcategories := GetSubCategories(categoryId)
					for _, i := range subcategories {
						AddCategoryRank(id, i.Id)
					}
					pageStruct.SuccessMessage = "Action effectué avec succès"
					AddCategoryRank(id, categoryId)
				} else {
					pageStruct.SuccessMessage = "Action effectué avec succès"
					AddCategoryRank(id, categoryId)
				}
			}
			if deleteadd == 2 {
				pageStruct.SuccessMessage = "Action effectué avec succès"
				DeleteCategoryRank(id, categoryId)
			}
		}
		if len(pageStruct.SuccessMessage) < 1 {
			pageStruct.ErrorMessage = "Une erreur c'est produite"
		}

		pageStruct.Page = "categoriespermissions"
		pageStruct.Users = GetUsersInfo()
		pageStruct.Category = GetAllCategories()
		pageStruct.UserCategory = GeAllCategoriesUsersRanks()

		err = templates["admin_categoriespermissions.page.gohtml"].Execute(w, pageStruct)
		if err != nil {
			log.Fatalln("template didn't execute: ", err)
		}
		return

	} else {
		pageStruct.Page = "categoriespermissions"
		pageStruct.Users = GetUsersInfo()
		pageStruct.Category = GetAllCategories()
		pageStruct.UserCategory = GeAllCategoriesUsersRanks()

		err = templates["admin_categoriespermissions.page.gohtml"].Execute(w, pageStruct)
		if err != nil {
			log.Fatalln("template didn't execute: ", err)
		}
		return
	}

}

// Fill the struct to see evert post & comment
// in my category(ies)
// (for modo and admin)
func AdminMyFeedHandler(w http.ResponseWriter, r *http.Request) {
	type AdminPage struct {
		ErrorMessage   string
		SuccessMessage string
		IsConnected    bool
		Page           string
		Categories     []Category
		Comments       []Comment
		User           UserInfo
		Posts          []Post
	}
	var pageStruct AdminPage
	connected, token := CheckCookie(r)
	if connected {
		userId := CheckToken(w, r, token)
		if userId != 0 {
			pageStruct.User = GetUserCookieInfo(userId)
			pageStruct.User.Avatar = NumberToPpIcon(pageStruct.User.Avatar)
			pageStruct.IsConnected = true
		}
	} else {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	rank := pageStruct.User.Permissions

	if rank.AdminPanelAccess != 1 {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	if rank.Id != 2 {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	pageStruct.Page = "adminmyfeed"
	pageStruct.Categories = GetALLCategoriesUsersRanksForUser(pageStruct.User.Id)
	tempCategoryTableId := GetIdTableSubCategoryRankForUser(pageStruct.User.Id)
	for _, e := range tempCategoryTableId {
		for _, e2 := range GetPostsByCategory(e) {
			pageStruct.Posts = append(pageStruct.Posts, e2)
		}
	}
	for _, e := range tempCategoryTableId {
		for _, e2 := range GetCommentsForPost(e) {
			pageStruct.Comments = append(pageStruct.Comments, e2)
		}
	}
	err = templates["admin_myfeed.page.gohtml"].Execute(w, pageStruct)
	if err != nil {
		log.Fatalln("template didn't execute: ", err)
	}
	return
}

// Handler filling the structs for statistics
// of the webite in a certain time interval
func AdminStatsHandler(w http.ResponseWriter, r *http.Request) {
	type AdminPage struct {
		ErrorMessage      string
		SuccessMessage    string
		IsConnected       bool
		Page              string
		User              UserInfo
		CommentNumber     int
		CategoryNumber    int
		SubCategoryNumber int
		PostNumber        int
		LikeNumber        int
		Dislike           int
		VotesStats        int
	}
	var pageStruct AdminPage
	connected, token := CheckCookie(r)
	if connected {
		userId := CheckToken(w, r, token)
		if userId != 0 {
			pageStruct.User = GetUserCookieInfo(userId)
			pageStruct.IsConnected = true
			pageStruct.User.Avatar = NumberToPpIcon(pageStruct.User.Avatar)
		}
	} else {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	rank := pageStruct.User.Permissions

	if rank.AdminPanelAccess != 1 {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}

	if r.Method == http.MethodPost {
		startdate := r.FormValue("startdate")
		endate := r.FormValue("endate")
		pageStruct.PostNumber = StatsDatePostNumber(startdate, endate)
		pageStruct.CommentNumber = StatsDateCommentNumber(startdate, endate)
		pageStruct.CategoryNumber = StatsDateCategoryNumber(startdate, endate)
		pageStruct.SubCategoryNumber = StatsDateSubCategoryNumber(startdate, endate)
		pageStruct.VotesStats = StatsVotes()

	} else {
		pageStruct.PostNumber = StatsPostNumber()
		pageStruct.CommentNumber = StatsCommentNumber()
		pageStruct.CategoryNumber = StatsCategoryNumber()
		pageStruct.SubCategoryNumber = StatsSubCategoryNumber()
		pageStruct.VotesStats = StatsVotes()
	}

	pageStruct.Page = "adminstats"
	err = templates["admin_stats.page.gohtml"].Execute(w, pageStruct)
	if err != nil {
		log.Fatalln("template didn't execute: ", err)
	}
	return
}

// Check to cookie validity each time we check the 'connected' users.
//  Return false if the cookie has been modified, or is expired.
func CheckCookieValidity(r *http.Request, userId int) bool {
	var cookie, err = r.Cookie("session")
	if err == nil {
		return GetAccessTokenByToken(cookie.Value).UserId == userId
	} else {
		fmt.Println(err)
		return false
	}
}

// Check if there is a cookie in the http request
// return true if there is one in + return it.
func CheckCookie(r *http.Request) (bool, string) {
	var cookie, err = r.Cookie("session")
	if err == nil {
		cookieValue := cookie.Value
		return true, cookieValue
	}
	return false, ""
}

// Create a new cookie with the help of
// UUID and return it.
func CreateUniqueCookie(w http.ResponseWriter) string {
	var cookieUuid string
	cookieUuid = uuid.New().String()
	cookie := &http.Cookie{
		Name:     "session",
		Value:    cookieUuid,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, cookie)
	return cookieUuid
}

// delete the cookie of the connected user.
// (using it with disconnection)
func DeleteCookie(w http.ResponseWriter, r *http.Request) (error, string) {
	var cookie, err = r.Cookie("session")
	if err != nil {
		return err, ""
	}
	cookieValue := cookie.Value
	var cookieUuid string
	cookie, err = r.Cookie("session")
	if err != nil {
		return err, ""
	} else {
		cookieUuid = uuid.New().String()
		cookie = &http.Cookie{
			Name:     "session",
			Value:    cookieUuid,
			Secure:   true,
			HttpOnly: true,
			Path:     "/",
			MaxAge:   -1,
		}
		http.SetCookie(w, cookie)
	}
	RemoveAccessToken(cookieValue)
	return err, cookieValue

}

// Like handler managing the like/dislike from a post
// and de/incrementing with insert functions
// the corresponding votes from this one.
func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err = r.ParseForm()
		if err != nil {
			fmt.Println(err)
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
			return
		}
		like, err := strconv.Atoi(r.Form.Get("like"))
		if err != nil {
			fmt.Println(err)
		}
		dislike, err := strconv.Atoi(r.Form.Get("dislike"))
		if err != nil {
			fmt.Println(err)
		}

		userId, err := strconv.Atoi(r.Form.Get("userId"))
		if err != nil {
			fmt.Println(err)
		}
		commentId, err := strconv.Atoi(r.Form.Get("commentId"))
		if err != nil {
			fmt.Println(err)
		}
		postId, err := strconv.Atoi(r.Form.Get("postId"))
		if err != nil {
			fmt.Println(err)
		}
		if CheckCookieValidity(r, userId) {
			if like == 1 && dislike == 0 {
				InsertVote(1, postId, commentId, userId)
			} else if like == 0 && dislike == 1 {
				InsertVote(2, postId, commentId, userId)
			} else if like == -1 && dislike == 0 {
				DeleteVote(userId, postId, commentId)
			} else if like == -1 && dislike == 1 {
				DeleteVote(userId, postId, commentId)
				InsertVote(2, postId, commentId, userId)
			} else if like == 0 && dislike == -1 {
				DeleteVote(userId, postId, commentId)
			} else if like == 1 && dislike == -1 {
				DeleteVote(userId, postId, commentId)
				InsertVote(1, postId, commentId, userId)
			} else {
				err = templates["404.page.gohtml"].Execute(w, basePage)
				w.WriteHeader(404)
			}
		} else {
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
		}

	} else {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
}

// Fill the struct with the post list of the corresponding
// cateogry + insert the promoted post/comment if there is one
// promoted in the category.
func PostPresentationHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/presentationpost" {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	if r.Method == http.MethodPost {
		err = r.ParseForm()
		if err != nil {
			fmt.Println(err)
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
			return
		}
		// parse form + get form values + InsertPost (in db)
		userId, err := strconv.Atoi(r.Form.Get("userId"))
		if err != nil {
			fmt.Println(err)
		}
		if CheckCookieValidity(r, userId) {
			categoryId, err := strconv.Atoi(r.Form.Get("categoryId"))
			if err != nil {
				fmt.Println(err)
			}
			if CheckCategoryExist(categoryId) {
				title, content := r.Form.Get("title"), r.Form.Get("content")
				InsertPost(title, content, 0, categoryId, userId)
			}
		}
	}
	query := r.URL.Query()
	postId, present := query["id"]
	if present && len(postId) == 1 {
		if _, err := strconv.ParseInt(postId[0], 10, 64); err == nil {
			CategoryId, err := strconv.Atoi(postId[0])
			if err != nil {
				fmt.Print(err)
			}
			if CheckCategoryExist(CategoryId) {
				page := DisplayPostList(CategoryId)
				connected, token := CheckCookie(r)
				if connected {
					page.IsConnected = true
					userId := CheckToken(w, r, token)
					if userId != 0 {
						page.User = GetUserCookieInfo(userId)
						BadgesEligibilityChecker(userId)
					}
					if CheckCategoryAuthorization(userId, page.TablePost[0].CategoryId) {
						page.IsAuthorized = true
					}
					page.User.Avatar = NumberToPpIcon(GetUser(page.User.Id).Avatar)

				}
				if !CheckPromoteExpiration(CategoryId) { // true if not expired
					page.Promote.IsPromoted = false
				}
				err = templates["presentationpost.page.gohtml"].Execute(w, page)
				if err != nil {
					fmt.Println("template didn't execute: ", err)
					// http.Redirect(w, r, r.URL.Path[0:], http.StatusSeeOther)

				}
			} else {
				fmt.Print("\nouste vilain client qui change les url et rajoute des parametres (҂‾ ▵‾)︻デ═一")
				err = templates["404.page.gohtml"].Execute(w, basePage)
				w.WriteHeader(404)
				return
			}
		} else {
			fmt.Print("\nvilain client qui met un string en id ٩ (╬ʘ益ʘ╬) ۶")
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
		}

	} else { // il faudrait aussi check si le post existe ou pas pour ne pas rien afficher ou chercher quelque chose qui n'existe pas et faire crash le serveur
		fmt.Print("\nouste vilain client qui change les url et rajoute des parametres (҂‾ ▵‾)︻デ═一")
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
	}
}

// Fill the struct for a
// single post page with evert comment related to
func PostHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	postId, present := query["id"]
	type SinglePostPage struct {
		IsConnected  bool
		IsAuthorized bool
		Page         string
		User         UserInfo
		Post         Post
	}
	var singlePostPage SinglePostPage
	connected, token := CheckCookie(r)
	userId := 0
	if connected {
		singlePostPage.IsConnected = true
		userId = CheckToken(w, r, token)
		if userId != 0 {
			singlePostPage.User = GetUserCookieInfo(userId)
			BadgesEligibilityChecker(userId)
		}
	} else {
		singlePostPage.IsConnected = false
	}
	if present && len(postId) == 1 {
		if _, err := strconv.ParseInt(postId[0], 10, 64); err == nil {
			PostId, err := strconv.Atoi(postId[0])
			if CheckPostExist(PostId) { // ne marche pas encore
				if err != nil {
					fmt.Print(err)
				}
				if r.Method == http.MethodPost {
					err = r.ParseForm()
					if err != nil {
						fmt.Println(err)
						err = templates["404.page.gohtml"].Execute(w, basePage)
						w.WriteHeader(404)
						return
					}
					// parse form + get form values + InsertPost (in db)
					postIdCheck, err := strconv.Atoi(r.Form.Get("postId"))
					if err != nil {
						fmt.Println(err)
					}
					if postIdCheck != PostId {
						err = templates["404.page.gohtml"].Execute(w, basePage)
						w.WriteHeader(404)
						return
					}
					parentId, err := strconv.Atoi(r.Form.Get("commentId"))
					if err != nil {
						fmt.Println(err)
					}
					content := r.Form.Get("content")
					if singlePostPage.IsConnected {
						InsertComment(content, singlePostPage.User.Id, PostId, parentId)
					}

				}
			}
			singlePostPage.Post = DisplayPost(PostId)
			if singlePostPage.IsConnected {
				singlePostPage.User.Avatar = NumberToPpIcon(GetUser(singlePostPage.User.Id).Avatar)

				singlePostPage.Post.LikedByUser = CheckLikePostByUser(singlePostPage.User.Id, PostId)
				singlePostPage.Post.DislikedByUser = CheckDislikePostByUser(singlePostPage.User.Id, PostId)
				for x := 0; x < len(singlePostPage.Post.Comment); x++ {
					singlePostPage.Post.Comment[x].DislikedByUser = CheckDislikeCommentByUser(singlePostPage.User.Id, singlePostPage.Post.Comment[x].Id)
					singlePostPage.Post.Comment[x].LikedByUser = CheckLikeCommentByUser(singlePostPage.User.Id, singlePostPage.Post.Comment[x].Id)
				}
			}
			if CheckCategoryAuthorization(userId, GetPost(PostId).CategoryId) {
				singlePostPage.IsAuthorized = true
			}
			err = templates["post.page.gohtml"].Execute(w, singlePostPage)
			if err != nil {
				log.Fatalln("template didn't execute: ", err)
			}
		} else {
			fmt.Print("\nouste vilain client qui change les url et rajoute des parametres (҂‾ ▵‾)︻デ═一")
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
		}

	} else { // il faudrait aussi check si le post existe ou pas pour ne pas rien afficher ou chercher quelque chose qui n'existe pas et faire crash le serveur
		fmt.Print("\nvilain client qui met un string en id ٩ (╬ʘ益ʘ╬) ۶")
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
	}

}

// Handler for all the profil info page, the post, comment list,
// and the profil infos/
func ProfilClientHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/profil" && r.URL.Path != "/Profil" {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	type profilPage struct {
		IsConnected bool
		User        UserInfo
		ErrorMsg    string
		Allposts    []Post
		AllComments []Comment
	}
	var ProfilPage profilPage
	connected, token := CheckCookie(r)
	if connected {
		ProfilPage.IsConnected = true
		userId := CheckToken(w, r, token)
		if userId != 0 {
			ProfilPage.User = GetUserCookieInfo(userId)
			BadgesEligibilityChecker(userId)
			ProfilPage.User.Avatar = NumberToPpIcon(ProfilPage.User.Avatar)
			ProfilPage.Allposts = GetPostsByUser(userId)
			ProfilPage.AllComments = GetCommentsByUser(userId)
			for x := range ProfilPage.Allposts {
				ProfilPage.Allposts[x].RedeemComsInPost()
			}
			for x := range ProfilPage.AllComments {
				ProfilPage.AllComments[x].RedeemParentCommentAndOwnLikesData()
			}
		}
		if r.Method != http.MethodPost {

			err = templates["profil.page.gohtml"].Execute(w, ProfilPage)
			if err != nil {
				log.Fatalln("template didn't execute: ", err)
			}
		} else {
			err := r.ParseForm()
			if err != nil {
				fmt.Println(err)
			}
			email, username, password, passwordconfirmation := r.FormValue("email"), r.FormValue("username"), r.FormValue("password"), r.FormValue("passwordconfirmation")
			if CheckEmailIsValid(email) {
				ModifyEmail(email, ProfilPage.User.Id)
			} else {
				ProfilPage.ErrorMsg = "Email is uncorrect."
			}
			if !CheckusernameExist(username) {
				ModifyUsername(username, ProfilPage.User.Id)
			} else {
				ProfilPage.ErrorMsg = "Username has already been taken."

			}
			fmt.Print(password, passwordconfirmation)
			if password == passwordconfirmation && password != "" {
				ModifyPassword(password, ProfilPage.User.Id)
			} else if password != passwordconfirmation && password != "" {
				ProfilPage.ErrorMsg = "Passwords do not match, please try again."
			}
			ProfilPage.User = GetUserCookieInfo(userId)
			BadgesEligibilityChecker(userId)
			ProfilPage.User.Avatar = NumberToPpIcon(ProfilPage.User.Avatar)

			err = templates["profil.page.gohtml"].Execute(w, ProfilPage)
			if err != nil {
				log.Fatalln("template didn't execute: ", err)
			}
		}
	} else {
		ProfilPage.IsConnected = false
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404) // Accesss denied page if possible
		return
	}

}

// If user did forgot his password, generate a random password,
// insert it in the db, and send an email with the corresponding
// password to the user.
func PasswordRecoveryHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/passwordrecovery" && r.URL.Path != "/PasswordRecovery" && r.URL.Path != "/Passwordrecovery" {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}
	type passwordRecoveryPage struct {
		SuccessMessage string
		ErrorMessage   string
		IsConnected    bool
	}
	var passwordPage passwordRecoveryPage
	if r.Method != http.MethodPost {
		err = templates["passwordrecovery.page.gohtml"].Execute(w, passwordPage)
		if err != nil {
			log.Fatalln("template didn't execute: ", err)
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		email := r.FormValue("email")
		if CheckEmailIsValid(email) && CheckEmailExist(email) {
			newpassword := ModifyPasswordRandom(email)
			SendMail(email, newpassword)
			passwordPage.SuccessMessage = "Recovery email has been sent!"
			err = templates["passwordrecovery.page.gohtml"].Execute(w, passwordPage)
			if err != nil {
				log.Fatalln("template didn't execute: ", err)
			}
		} else {
			passwordPage.ErrorMessage = "Uncorrect email / Not an existing account."
			err = templates["passwordrecovery.page.gohtml"].Execute(w, passwordPage)

		}
	}

}

// function to delete a post if he is owned by the connected user.
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err = r.ParseForm()
		if err != nil {
			fmt.Println(err)
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
			return
		}
		id, err := strconv.Atoi(r.Form.Get("id"))
		if err != nil {
			fmt.Println(err)
		}
		Type := r.FormValue("type")
		if err != nil {
			fmt.Println(err)
		}
		if Type != "post" {
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
			return
		}
		test_bis := false
		test := false
		connected, token := CheckCookie(r)
		if connected {
			userId := CheckToken(w, r, token)
			test_bis = CheckCookieValidity(r, userId)
			if test_bis {
				test = CheckPostOwnedByUser(id, userId)
				if test {
					DeletePost(id, userId)
				} else {
					err = templates["404.page.gohtml"].Execute(w, basePage)
					w.WriteHeader(404)
				}
			} else {
				err = templates["404.page.gohtml"].Execute(w, basePage)
				w.WriteHeader(404)
				return
			}
		}

	} else {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}

}

// Function to delete the post if the connected user is a moderator
// with the corresponding category authorization.
func DeletePostModoHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		err = r.ParseForm()
		if err != nil {
			fmt.Println(err)
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
			return
		}

		id, err := strconv.Atoi(r.Form.Get("id"))
		if err != nil {
			fmt.Println(err)
		}

		justification := r.FormValue("justification")
		if justification == "" {
			return
		}

		test_bis := false
		connected, token := CheckCookie(r)

		if connected {
			userId := CheckToken(w, r, token)
			justification = "       ⚠  DELETED BY: " + GetUser(userId).Username + "REASON: " + justification + "⚠     "
			test_bis = CheckCookieValidity(r, userId)
			if test_bis {

				authorized := CheckCategoryAuthorization(userId, GetPost(id).CategoryId)
				if authorized {
					fmt.Print("test")
					ModifyPost(id, GetPost(id).Subject, justification)
					DeletePostModo(id)

				} else {
					w.WriteHeader(403) // unauthorized
				}
			} else {
				err = templates["404.page.gohtml"].Execute(w, basePage)
				w.WriteHeader(404)
				return
			}
		}

	} else {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}

}

// Delete the comment if the connected user is a moderator
// with the corresponding category authorization.
func DeleteCommentModoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err = r.ParseForm()
		if err != nil {
			fmt.Println(err)
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
			return
		}

		id, _ := strconv.Atoi(r.FormValue("id"))
		justification := r.FormValue("justification")
		if justification == "" {
			return
		}
		test_bis := false
		connected, token := CheckCookie(r)
		if connected {
			userId := CheckToken(w, r, token)
			justification = "       ⚠  DELETED BY: " + GetUser(userId).Username + "REASON: " + justification + "⚠     "
			test_bis = CheckCookieValidity(r, userId)
			if test_bis {
				authorized := CheckCategoryAuthorization(userId, GetPost(GetComment(id).PostId).CategoryId)
				if authorized {
					fmt.Print("salut")
					ModifyComment(id, justification)
					DeleteCommentModo(id)
				} else {
					w.WriteHeader(403) // unauthorized
				}
			} else {
				err = templates["404.page.gohtml"].Execute(w, basePage)
				w.WriteHeader(404)
				return
			}
		}
	} else {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	}

}

func ProfilDeleteHandler(w http.ResponseWriter, r *http.Request) {
	type IndexPage struct {
		TableCategory []Category
		IsConnected   bool
		User          UserInfo
	}
	connected, token := CheckCookie(r)
	if connected {
		userId := CheckToken(w, r, token)
		if CheckCookieValidity(r, userId) {
			DeleteProfil(userId)
			DeleteCookie(w, r)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Handler to see the tickets from the client side
func TicketsHandler(w http.ResponseWriter, r *http.Request) {
	type PageTickets struct {
		IsConnected    bool
		Page           string
		Tickets        []Tickets
		Ticket         Tickets
		TicketMessages []TicketsMessages
		User           UserInfo
	}
	var TicketsPage PageTickets
	connected, token := CheckCookie(r)
	if connected {
		TicketsPage.IsConnected = true
		userId := CheckToken(w, r, token)
		if userId != 0 {
			TicketsPage.User = GetUserCookieInfo(userId)
			BadgesEligibilityChecker(userId)
			TicketsPage.User.Avatar = NumberToPpIcon(TicketsPage.User.Avatar)
			TicketsPage.Tickets = GetUserTickets(userId)
			TicketsPage.Page = "ticket"
			err = templates["userticket.page.gohtml"].Execute(w, TicketsPage)

		}
	}
}

// Show the ticket conversation for the client side.
func TicketsMessageHandler(w http.ResponseWriter, r *http.Request) {
	type PageTickets struct {
		IsConnected    bool
		Page           string
		Tickets        []Tickets
		Ticket         Tickets
		TicketMessages []TicketsMessages
		User           UserInfo
	}
	var TicketsPage PageTickets
	query := r.URL.Query()
	postId, present := query["ticketpage"]
	if r.Method == "POST" {
		connected, token := CheckCookie(r)
		if connected {
			TicketsPage.IsConnected = true
			userId := CheckToken(w, r, token)

			if userId != 0 {
				if ticketId, ok := r.URL.Query()["ticketpage"]; ok {
					ticket, err := strconv.Atoi(ticketId[0])
					if err != nil {
						err = templates["404.page.gohtml"].Execute(w, basePage)
						w.WriteHeader(404)
						return
					}
					tempTicket := GetTicket(ticket)
					if !reflect.ValueOf(tempTicket).IsZero() {
						if r.Method == http.MethodPost {
							message := r.FormValue("sendticket-message")
							if len(message) > 1 {
								InsertTicketMessage(message, userId, ticket)
							}
						}
					}
				}
			}
		}
	}
	if present && len(postId) == 1 {
		if _, err := strconv.ParseInt(postId[0], 10, 64); err == nil {
			TicketId, err := strconv.Atoi(postId[0])
			if err != nil {
				fmt.Print(err)
			}
			connected, token := CheckCookie(r)
			if connected {
				TicketsPage.IsConnected = true
				userId := CheckToken(w, r, token)
				if userId != 0 {
					TicketsPage.User = GetUserCookieInfo(userId)
					BadgesEligibilityChecker(userId)
					TicketsPage.User.Avatar = NumberToPpIcon(TicketsPage.User.Avatar)
					TicketsPage.Tickets = GetUserTickets(userId)
					TicketsPage.Page = "ticketsmessage"
					TicketsPage.TicketMessages = GetTicketMessages(TicketId)
					for i := range TicketsPage.TicketMessages {
						TicketsPage.TicketMessages[i].Creator.Username = GetUser(TicketsPage.TicketMessages[i].UserId).Username
					}
					err = templates["userticketmessage.page.gohtml"].Execute(w, TicketsPage)
					if err != nil {
						fmt.Print("template didnt execute =>", err)
					}

				}
			}
		}
	}
}

// Search handler with showing at first every post
// and comment of the site, and if u make a search with,
// return only the corresponding ones.
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	type SearchStruct struct {
		IsConnected bool
		Page        string
		User        UserInfo
		AllPosts    []Post
		AllComments []Comment
	}
	var searchstruct SearchStruct
	connected, token := CheckCookie(r)
	if connected {
		searchstruct.IsConnected = true
		userId := CheckToken(w, r, token)
		if userId != 0 {
			searchstruct.User = GetUserCookieInfo(userId)
			BadgesEligibilityChecker(userId)
		}
		searchstruct.User.Avatar = NumberToPpIcon(GetUser(searchstruct.User.Id).Avatar)

	}

	searchstruct.AllPosts = GetPosts()
	for i, post := range searchstruct.AllPosts {
		searchstruct.AllPosts[i].Comment = GetCommentsForPost(post.Id)
		searchstruct.AllPosts[i].NumberOfComments = len(GetCommentsForPost(post.Id))
		searchstruct.AllPosts[i].Votes = CountVotes(GetVotesByPost(post.Id))
		searchstruct.AllPosts[i].CategoryTitle = GetCategory(post.CategoryId).Title
		searchstruct.AllPosts[i].UserName = GetUser(post.UserId).Username
		searchstruct.AllPosts[i].Avatar = NumberToPpIcon(GetUser(post.UserId).Avatar)
		searchstruct.AllPosts[i].BadgesCreator = GetUserBadges(post.UserId)
	}
	searchstruct.AllComments = GetComments()
	for I := range searchstruct.AllComments {
		if searchstruct.AllComments[I].ParentId != 0 { // si post parent on recup le commentaire, l'id du createur, et la date de creation
			searchstruct.AllComments[I].ParentComment, searchstruct.AllComments[I].ParentUserId, searchstruct.AllComments[I].ParentCreationDate = GetComment(searchstruct.AllComments[I].ParentId).Comment, GetComment(searchstruct.AllComments[I].ParentId).UserId, GetComment(searchstruct.AllComments[I].ParentId).CreationDate
		}
		searchstruct.AllComments[I].Votes = CountVotes(GetVotesByComment(searchstruct.AllComments[I].Id))
		searchstruct.AllComments[I].UserName = GetUser(searchstruct.AllComments[I].UserId).Username
		searchstruct.AllComments[I].CreationDate = searchstruct.AllComments[I].CreationDate[0:16]
		searchstruct.AllComments[I].Avatar = NumberToPpIcon(GetUser(searchstruct.AllComments[I].UserId).Avatar)
		RegisterDate := GetUser(searchstruct.AllComments[I].UserId).RegisterDate
		searchstruct.AllComments[I].RegisterDate = RegisterDate[5:10] + "-" + RegisterDate[0:4]
		searchstruct.AllComments[I].UserRank = RankIntToRankString(GetUser(searchstruct.AllComments[I].UserId).Rank)
		searchstruct.AllComments[I].BadgesCreator = GetUserBadges(searchstruct.AllComments[I].UserId)
	}
	var NewSearchStruct SearchStruct
	if connected {
		NewSearchStruct.IsConnected = true
		userId := CheckToken(w, r, token)
		if userId != 0 {
			NewSearchStruct.User = GetUserCookieInfo(userId)
			BadgesEligibilityChecker(userId)
		}
		NewSearchStruct.User.Avatar = NumberToPpIcon(GetUser(NewSearchStruct.User.Id).Avatar)

	}
	if searchParameter, ok := r.URL.Query()["research"]; ok {
		for i := range searchstruct.AllPosts {
			if strings.Contains(strings.ToLower(searchstruct.AllPosts[i].Subject), strings.ToLower(searchParameter[0])) || strings.Contains(strings.ToLower(searchstruct.AllPosts[i].Title), strings.ToLower(searchParameter[0])) {
				NewSearchStruct.AllPosts = append(NewSearchStruct.AllPosts, searchstruct.AllPosts[i])
			}
		}
		for i := range searchstruct.AllComments {
			if strings.Contains(strings.ToLower(searchstruct.AllComments[i].Comment), strings.ToLower(searchParameter[0])) {
				NewSearchStruct.AllComments = append(NewSearchStruct.AllComments, searchstruct.AllComments[i])
			}
		}
		err = templates["search.page.gohtml"].Execute(w, NewSearchStruct)
		if err != nil {
			fmt.Print("template didnt execute =>", err)
		}
	} else {

		err = templates["search.page.gohtml"].Execute(w, searchstruct)
		if err != nil {
			fmt.Print("template didnt execute =>", err)
		}
	}

}

// Handler for creation of a new ticket from the client side.
func NewTicketHandler(w http.ResponseWriter, r *http.Request) {
	type NewTicket struct {
		IsConnected bool
		Page        string
		User        UserInfo
	}

	var newticket NewTicket
	connected, token := CheckCookie(r)
	if connected {
		newticket.IsConnected = true
		userId := CheckToken(w, r, token)
		if userId != 0 {
			newticket.User = GetUserCookieInfo(userId)
			BadgesEligibilityChecker(userId)
		}
		newticket.User.Avatar = NumberToPpIcon(GetUser(newticket.User.Id).Avatar)

	} else {
		w.WriteHeader(403)
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		subject := r.FormValue("subject")
		if title != "" && subject != "" {
			InsertTicket(title, subject, 0, newticket.User.Id)
			http.Redirect(w, r, "/tickets", http.StatusSeeOther)
		} else {
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
		}

	}
	err = templates["newticket.page.gohtml"].Execute(w, newticket)
	if err != nil {
		fmt.Print("template didnt execute =>", err)
	}

}

// For promoting a post in a category.
// (only modo and admin can do this)
func PromotePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	} else {
		err = r.ParseForm()
		if err != nil {
			fmt.Println(err)
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
			return
		}
		connected, token := CheckCookie(r)
		if connected {
			userId := CheckToken(w, r, token)
			if CheckCookieValidity(r, userId) {
				postId, err := strconv.Atoi(r.Form.Get("postId"))
				if err != nil {
					fmt.Println(err)
				}
				endDate := r.Form.Get("expiration")
				if CheckPromoteExist(GetPost(postId).CategoryId) {
					DeletePromote(GetPost(postId).CategoryId)
				}
				if CheckPostExist(postId) {
					InsertPromote(endDate, postId, 0, GetPost(postId).CategoryId)
				}
			}
		}
	}
}

// For promoting a comment in a category.
// (only modo and admin can do this)
func PromoteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err = templates["404.page.gohtml"].Execute(w, basePage)
		w.WriteHeader(404)
		return
	} else {
		err = r.ParseForm()
		if err != nil {
			fmt.Println(err)
			err = templates["404.page.gohtml"].Execute(w, basePage)
			w.WriteHeader(404)
			return
		}
		connected, token := CheckCookie(r)
		if connected {
			userId := CheckToken(w, r, token)
			if CheckCookieValidity(r, userId) {
				commentId, err := strconv.Atoi(r.Form.Get("commentId"))
				if err != nil {
					fmt.Println(err)
				}
				endDate := r.Form.Get("expiration")
				if CheckPromoteExist(GetPost(GetComment(commentId).PostId).CategoryId) {
					DeletePromote(GetPost(GetComment(commentId).PostId).CategoryId)
				}
				if CheckCommentInPostExist(commentId) {
					InsertPromote(endDate, 0, commentId, GetPost(GetComment(commentId).PostId).CategoryId)
				}
			}
		}
	}
}

// Check if a promotion already exist in the category.
// The code using it delete the one already existing
// If it return a true
func CheckPromoteExist(CategoryId int) bool {
	row := db.QueryRow("SELECT id FROM promotes WHERE categoryId= ?", CategoryId) // db correspond a la database ouverte dans db.go
	temp := ""
	err := row.Scan(&temp)
	if err != nil {
		fmt.Println(err)
	}
	return temp != ""
}

func DeletePromote(CategoryId int) {
	db.Exec("DELETE FROM promotes WHERE categoryId = ?", CategoryId)
	// pas de handle pour err car il envoie des erreurs chaques premiers promotes de categories
}

// return if there is a propmoted content (true false),
// and the type of it (1=>post,2=>comment)
// + the corresponding id.
func GetPromotedContent(categoryId int) (bool, int, int) {
	row := db.QueryRow("SELECT * FROM promotes WHERE categoryId= ?", categoryId) // db correspond a la database ouverte dans db.go
	var temp Promotes
	err := row.Scan(&temp.Id, &temp.EndDate, &temp.PostId, &temp.CommentId, &temp.CategoryId)
	if err != nil {
		fmt.Println(err)
	}
	if temp.Id != 0 {

		if temp.CommentId == 0 {
			return true, 1, temp.PostId
		} else {
			return true, 2, temp.CommentId
		}
	} else {
		return false, 0, 0
	}

}

// Check if the end date of the promotion in the category is expired.
// true if expired.
func CheckPromoteExpiration(categoryId int) bool {
	row := db.QueryRow("SELECT endDate FROM promotes WHERE categoryId = ?", categoryId)
	var endDate string
	err := row.Scan(&endDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			fmt.Print(err)
		}
	}
	fmt.Print(endDate)
	enddate, err := time.Parse("02-01-2006 15:04:05", endDate)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(enddate)
	return true // il faut parse comme il faut pour comparer time.now et enddate
}

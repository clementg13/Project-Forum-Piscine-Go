package main

import (
	"fmt"
	"log"
	"net/http"
)

// All the initialization + handle function
func main() {
	defer db.Close()
	// little but cute message é_è
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nBy Clément Grosieux & Despres Antoine\n")
	fmt.Printf("Ctrl+c for turn off the server\n")
	fmt.Printf("\nServer start on https://localhost:3123\n")
	//InitTable()
	//InitBadges()
	//InitTableVariable()
	err := InitTemplate()
	if err != nil {
		fmt.Println(err)
	}
	IconsArray()

	// include the static path for js and css
	http.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(http.Dir("statics"))))
	// include the page
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/disconnect", DisconnectHandler)

	http.HandleFunc("/admin", AdminTicketHandler)
	http.HandleFunc("/admin_tickets", AdminTicketHandler)
	http.HandleFunc("/admin_ticketpage", AdminTicketViewHandler)
	http.HandleFunc("/admin_gestionusers", AdminGestionUsersHandler)
	http.HandleFunc("/admin_categoriespermissions", AdminCategoriesPermission)
	http.HandleFunc("/admin_myfeed", AdminMyFeedHandler)
	http.HandleFunc("/admin_stats", AdminStatsHandler)

	http.HandleFunc("/presentationpost", PostPresentationHandler)
	http.HandleFunc("/post", PostHandler)
	http.HandleFunc("/sendlike", LikeHandler)
	http.HandleFunc("/profil", ProfilClientHandler)
	http.HandleFunc("/profildeletion", ProfilDeleteHandler)
	http.HandleFunc("/passwordrecovery", PasswordRecoveryHandler)
	http.HandleFunc("/deletepost", DeletePostHandler)
	http.HandleFunc("/deletepostmodo", DeletePostModoHandler)
	http.HandleFunc("/deletecommentmodo", DeleteCommentModoHandler)
	http.HandleFunc("/promotepost", PromotePostHandler)
	http.HandleFunc("/promotecomment", PromoteCommentHandler)
	http.HandleFunc("/tickets", TicketsHandler)
	http.HandleFunc("/newticket", NewTicketHandler)
	http.HandleFunc("/ticketsmessage", TicketsMessageHandler)
	http.HandleFunc("/search", SearchHandler)
	// put the server on a port
	http.ListenAndServeTLS(":3123", "https-server.crt", "https-server.key", nil)

}

type BasePage struct {
	ErrorMessage   string
	SuccessMessage string
	IsConnected    bool
	Page           string
	User           UserInfo
}

var basePage BasePage

package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func InitTableVariable() {
	userFakeData, err := db.Prepare(`INSERT INTO users (username,email,password,avatar,registerDate,ban,rank) VALUES (?,?,?,?,?,?,?)`)
	if err != nil {
		fmt.Println(err)
	}

	primaryCategoryFakeData, err := db.Prepare(`INSERT INTO categories (title,description,creationDate) VALUES (?,?,?)`)
	if err != nil {
		fmt.Print(err)
	}

	secondaryCategoryFakeData, err := db.Prepare(`INSERT INTO categories (title,description,icon,creationDate,ParentId) VALUES (?,?,?,?,?)`)
	if err != nil {
		fmt.Print(err)
	}
	postFakeData, err := db.Prepare(`INSERT INTO posts (title,subject,creationDate,categoryId,userId) VALUES (?,?,?,?,?)`)
	if err != nil {
		fmt.Print(err)
	}
	commentFakeData, err := db.Prepare(`INSERT INTO comments (comment,creationDate,userId,postId) VALUES (?,?,?,?)`)
	if err != nil {
		fmt.Println(err)
	}
	rankFakeData, err := db.Prepare(`INSERT INTO ranks (name,adminPanelAccess,ban,deban,deletePost,deleteUser,modifyCategory,ticketAccess,viewClosedPost,deleteComment,badgeAttribution,modifyRank) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		fmt.Println(err)
	}
	user_badgesFakeData, err := db.Prepare(`INSERT INTO users_badges (badgeId,userId) VALUES (?,?)`)
	if err != nil {
		fmt.Print(err)
	}

	categoriesRanks_badgesFakeData, err := db.Prepare(`INSERT INTO categories_users_ranks (userId,categoryId) VALUES (?,?)`)
	if err != nil {
		log.Fatal(err)
	}

	encrypt, err := HashPassword("admin123")
	if err != nil {
		log.Fatal(err)
	}
	_, err = userFakeData.Exec("Admin", "admin@gmail.com", encrypt, "4", time.Now().Format("02-01-2006 15:04:05"), 0, 3)
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(time.Now().UnixNano())
	profilpicture_int := strconv.Itoa(rand.Intn(50-1) + 1)

	_, err = userFakeData.Exec("Antoine", "AntoineTest@gmail.com", encrypt, profilpicture_int, time.Now().Format("02-01-2006 15:04:05"), 0, 2)
	if err != nil {
		fmt.Print(err)
	}

	profilpicture_int = strconv.Itoa(rand.Intn(50-1) + 1)

	_, err = userFakeData.Exec("M.ban", "BanTest@gmail.com", encrypt, profilpicture_int, time.Now().Format("02-01-2006 15:04:05"), 1, 1)
	if err != nil {
		fmt.Print(err)
	}
	profilpicture_int = strconv.Itoa(rand.Intn(50-1) + 1)

	_, err = userFakeData.Exec("M.toutLeMonde", "EverybodyTest@gmail.com", encrypt, profilpicture_int, time.Now().Format("02-01-2006 15:04:05"), 0, 1)
	if err != nil {
		fmt.Print(err)
	}

	_, err = primaryCategoryFakeData.Exec("Informatique", "blablabla Informatique useless qu'ils disaient", time.Now().Format("02-01-2006 15:04:05"))
	if err != nil {
		fmt.Print(err)
	}
	_, err = secondaryCategoryFakeData.Exec("Cybersecurité", "purée la securité", "fas fa-shield-alt", time.Now().Format("02-01-2006 15:04:05"), 1)
	if err != nil {
		fmt.Print(err)
	}
	_, err = secondaryCategoryFakeData.Exec("Software", "le soft c'est rigolo", "fas fa-laptop-code", time.Now().Format("02-01-2006 15:04:05"), 1)
	if err != nil {
		fmt.Print(err)
	}
	_, err = secondaryCategoryFakeData.Exec("Hardware", "bof le hard", "fas fa-microchip", time.Now().Format("02-01-2006 15:04:05"), 1)
	if err != nil {
		fmt.Print(err)
	}
	_, err = primaryCategoryFakeData.Exec("Sport", "le sport", time.Now().Format("02-01-2006 15:04:05"))
	if err != nil {
		fmt.Print(err)
	}
	_, err = secondaryCategoryFakeData.Exec("Yoga", "saint", "fas fa-yin-yang", time.Now().Format("02-01-2006 15:04:05"), 5)
	if err != nil {
		fmt.Print(err)
	}
	_, err = secondaryCategoryFakeData.Exec("Foot", "le soft c'est rigolo", "far fa-futbol", time.Now().Format("02-01-2006 15:04:05"), 5)
	if err != nil {
		fmt.Print(err)
	}
	_, err = secondaryCategoryFakeData.Exec("Fitness", "bof le hard", "fas fa-heartbeat", time.Now().Format("02-01-2006 15:04:05"), 5)
	if err != nil {
		fmt.Print(err)
	}
	_, err = postFakeData.Exec("Le saucisson, la nouvelle alimentation des ordinateurs", "Nouveau projet revolutionnaire permettant d'alimenter avec des saucissons votre machine.", time.Now().Format("02-01-2006 15:04:05"), 4, 1)
	if err != nil {
		fmt.Print(err)
	}
	_, err = postFakeData.Exec("Qu'est-ce que le Lorem Ipsum?", "Le Lorem Ipsum est simplement du faux texte employé dans la composition et la mise en page avant impression. ", time.Now().Format("02-01-2006 15:04:05"), 3, 2)
	if err != nil {
		fmt.Print(err)
	}
	_, err = postFakeData.Exec("Le Marathon le danger du dos", "Faut pas courir ca fait bobo au dos ET aux mollets :cc", time.Now().Format("02-01-2006 15:04:05"), 8, 2)
	if err != nil {
		fmt.Print(err)
	}
	_, err = postFakeData.Exec("Le Foot c'est de la merde", "Et oui c'est surcoté :@", time.Now().Format("02-01-2006 15:04:05"), 7, 1)
	if err != nil {
		fmt.Print(err)
	}
	_, err = postFakeData.Exec("OVH se fait hacker", "Après le barbecue géant, OVH est maintenant un véritable nid à virus, quelle malchance.", time.Now().Format("02-01-2006 15:04:05"), 2, 1)
	if err != nil {
		fmt.Print(err)
	}
	_, err = postFakeData.Exec("UP to stream bon lecteur ", "Tres pratique", time.Now().Format("02-01-2006 15:04:05"), 2, 3)
	if err != nil {
		fmt.Print(err)
	}
	_, err = postFakeData.Exec("Nouveau logiciel anti cookies", "On l'appelle le lance flammes", time.Now().Format("02-01-2006 15:04:05"), 2, 2)
	if err != nil {
		fmt.Print(err)
	}
	_, err = postFakeData.Exec("C'est sympa quand meme le dev", "Hop le forum est plié", time.Now().Format("02-01-2006 15:04:05"), 2, 1)
	if err != nil {
		fmt.Print(err)
	}
	_, err = postFakeData.Exec("Mediter c'est bon pour la santé ( la votre et celle des autres )", "Mediter reduit le taux de mortalité et la chance de tuer quelqu'un de 20%", time.Now().Format("02-01-2006 15:04:05"), 6, 4)
	if err != nil {
		fmt.Print(err)
	}
	_, err = postFakeData.Exec("Lorem Ipsum le retour", "Donec molestie, magna ut luctus ultrices, tellus arcu nonummy velit, sit amet pulvinar elit justo et mauris. In pede. Maecenas euismod elit eu erat. Aliquam augue wisi, facilisis congue, suscipit in, adipiscing et, ante. In justo. Cras lobortis neque ac ipsum. Nunc fermentum massa at ante. ", time.Now().Format("02-01-2006 15:04:05"), 3, 2)
	if err != nil {
		fmt.Print(err)
	}
	_, err = postFakeData.Exec("Neymar l'andouille", "Il fait que simuler mdrrr", time.Now().Format("02-01-2006 15:04:05"), 7, 2)
	if err != nil {
		fmt.Print(err)
	}
	_, err = postFakeData.Exec("La cybersecurité c'est cool", "c'est en vogue.", time.Now().Format("02-01-2006 15:04:05"), 2, 2)
	if err != nil {
		fmt.Print(err)
	}
	_, err = commentFakeData.Exec("Maecenas1 mi massa, fermentum eu, venenatis et, cursus id, ipsum. Morbi vehicula justo faucibus mauris. Donec non neque. Fusce id mi ut neque tincidunt posuere. Suspendisse quis enim. Cras porttitor. Sed quis velit.", time.Now().Format("02-01-2006 15:04:05"), 2, 1)
	if err != nil {
		fmt.Print(err)
	}
	_, err = commentFakeData.Exec("Maecenas2 mi massa, fermentum eu, venenatis et, cursus id, ipsum. Morbi vehicula justo faucibus mauris. Donec non neque. Fusce id mi ut neque tincidunt posuere. Suspendisse quis enim. Cras porttitor. Sed quis velit.", time.Now().Format("02-01-2006 15:04:05"), 1, 2)
	if err != nil {
		fmt.Print(err)
	}
	_, err = commentFakeData.Exec("Maecenas3 mi massa, fermentum eu, venenatis et, cursus id, ipsum. Morbi vehicula justo faucibus mauris. Donec non neque. Fusce id mi ut neque tincidunt posuere. Suspendisse quis enim. Cras porttitor. Sed quis velit.", time.Now().Format("02-01-2006 15:04:05"), 2, 3)
	if err != nil {
		fmt.Print(err)
	}
	_, err = commentFakeData.Exec("Maecenas4 mi massa, fermentum eu, venenatis et, cursus id, ipsum. Morbi vehicula justo faucibus mauris. Donec non neque. Fusce id mi ut neque tincidunt posuere. Suspendisse quis enim. Cras porttitor. Sed quis velit.", time.Now().Format("02-01-2006 15:04:05"), 3, 4)
	if err != nil {
		fmt.Print(err)
	}
	_, err = commentFakeData.Exec("Maecenas5 mi massa, fermentum eu, venenatis et, cursus id, ipsum. Morbi vehicula justo faucibus mauris. Donec non neque. Fusce id mi ut neque tincidunt posuere. Suspendisse quis enim. Cras porttitor. Sed quis velit.", time.Now().Format("02-01-2006 15:04:05"), 1, 5)
	if err != nil {
		fmt.Print(err)
	}
	_, err = commentFakeData.Exec("Maecenas6 mi massa, fermentum eu, venenatis et, cursus id, ipsum. Morbi vehicula justo faucibus mauris. Donec non neque. Fusce id mi ut neque tincidunt posuere. Suspendisse quis enim. Cras porttitor. Sed quis velit.", time.Now().Format("02-01-2006 15:04:05"), 2, 6)
	if err != nil {
		fmt.Print(err)
	}
	_, err = commentFakeData.Exec("Maecenas7 mi massa, fermentum eu, venenatis et, cursus id, ipsum. Morbi vehicula justo faucibus mauris. Donec non neque. Fusce id mi ut neque tincidunt posuere. Suspendisse quis enim. Cras porttitor. Sed quis velit.", time.Now().Format("02-01-2006 15:04:05"), 3, 7)
	if err != nil {
		fmt.Print(err)
	}
	_, err = commentFakeData.Exec("Maecenas8 mi massa, fermentum eu, venenatis et, cursus id, ipsum. Morbi vehicula justo faucibus mauris. Donec non neque. Fusce id mi ut neque tincidunt posuere. Suspendisse quis enim. Cras porttitor. Sed quis velit.", time.Now().Format("02-01-2006 15:04:05"), 4, 8)
	if err != nil {
		fmt.Print(err)
	}
	_, err = commentFakeData.Exec("Maecenas9 mi massa, fermentum eu, venenatis et, cursus id, ipsum. Morbi vehicula justo faucibus mauris. Donec non neque. Fusce id mi ut neque tincidunt posuere. Suspendisse quis enim. Cras porttitor. Sed quis velit.", time.Now().Format("02-01-2006 15:04:05"), 1, 9)
	if err != nil {
		fmt.Print(err)
	}
	_, err = commentFakeData.Exec("Maecenas10 mi massa, fermentum eu, venenatis et, cursus id, ipsum. Morbi vehicula justo faucibus mauris. Donec non neque. Fusce id mi ut neque tincidunt posuere. Suspendisse quis enim. Cras porttitor. Sed quis velit.", time.Now().Format("02-01-2006 15:04:05"), 2, 1)
	if err != nil {
		fmt.Print(err)
	}
	_, err = commentFakeData.Exec("Maecenas11 mi massa, fermentum eu, venenatis et, cursus id, ipsum. Morbi vehicula justo faucibus mauris. Donec non neque. Fusce id mi ut neque tincidunt posuere. Suspendisse quis enim. Cras porttitor. Sed quis velit.", time.Now().Format("02-01-2006 15:04:05"), 1, 2)
	if err != nil {
		fmt.Print(err)
	}
	_, err = commentFakeData.Exec("Maecenas12 mi massa, fermentum eu, venenatis et, cursus id, ipsum. Morbi vehicula justo faucibus mauris. Donec non neque. Fusce id mi ut neque tincidunt posuere. Suspendisse quis enim. Cras porttitor. Sed quis velit.", time.Now().Format("02-01-2006 15:04:05"), 2, 2)
	if err != nil {
		fmt.Print(err)
	}
	_, err = commentFakeData.Exec("Maecenas13 mi massa, fermentum eu, venenatis et, cursus id, ipsum. Morbi vehicula justo faucibus mauris. Donec non neque. Fusce id mi ut neque tincidunt posuere. Suspendisse quis enim. Cras porttitor. Sed quis velit.", time.Now().Format("02-01-2006 15:04:05"), 1, 7)
	if err != nil {
		fmt.Print(err)
	}
	_, err = commentFakeData.Exec("Maecenas14 mi massa, fermentum eu, venenatis et, cursus id, ipsum. Morbi vehicula justo faucibus mauris. Donec non neque. Fusce id mi ut neque tincidunt posuere. Suspendisse quis enim. Cras porttitor. Sed quis velit.", time.Now().Format("02-01-2006 15:04:05"), 4, 6)
	if err != nil {
		fmt.Print(err)
	}

	_, err = rankFakeData.Exec("user", 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	if err != nil {
		fmt.Print(err)
	}
	_, err = rankFakeData.Exec("moderator", 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0)
	if err != nil {
		fmt.Print(err)
	}
	_, err = rankFakeData.Exec("administrator", 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	if err != nil {
		fmt.Print(err)
	}
	_, err = user_badgesFakeData.Exec(1, 2)
	if err != nil {
		fmt.Print(err)
	}
	_, err = user_badgesFakeData.Exec(3, 1)
	if err != nil {
		fmt.Print(err)
	}
	_, err = user_badgesFakeData.Exec(3, 2)
	if err != nil {
		fmt.Print(err)
	}
	_, err = user_badgesFakeData.Exec(2, 2)
	if err != nil {
		fmt.Print(err)
	}
	_, err = user_badgesFakeData.Exec(2, 1)
	if err != nil {
		fmt.Print(err)
	}

	_, err = user_badgesFakeData.Exec(2, 3)
	if err != nil {
		fmt.Print(err)
	}
	_, err = user_badgesFakeData.Exec(1, 4)
	if err != nil {
		fmt.Print(err)
	}

	_, err = categoriesRanks_badgesFakeData.Exec(2, 1)
	if err != nil {
		log.Fatal(err)
	}
	_, err = categoriesRanks_badgesFakeData.Exec(2, 6)
	if err != nil {
		log.Fatal(err)
	}
	_, err = categoriesRanks_badgesFakeData.Exec(2, 4)
	if err != nil {
		log.Fatal(err)
	}
	_, err = categoriesRanks_badgesFakeData.Exec(2, 5)
	if err != nil {
		log.Fatal(err)
	}
	_, err = categoriesRanks_badgesFakeData.Exec(2, 3)
	if err != nil {
		log.Fatal(err)
	}
	_, err = categoriesRanks_badgesFakeData.Exec(2, 2)
	if err != nil {
		log.Fatal(err)
	}
	InsertComment("Comment1 of comment", 1, 3, 10)
	InsertComment("Comment2 of comment", 2, 1, 12)
	InsertComment("Comment3 of comment", 3, 3, 11)
	InsertComment("Comment4 of comment", 1, 3, 13)
	InsertComment("Comment5 of comment", 2, 1, 14)
	InsertComment("Comment6 of comment", 3, 2, 15)
	InsertComment("Comment7 of comment", 2, 6, 5)
	InsertComment("Comment8 of comment", 1, 2, 2)
	InsertComment("Comment9 of comment", 3, 3, 5)
	InsertComment("Comment10 of comment", 1, 3, 2)
	InsertComment("Comment11 of comment", 2, 3, 6)
	InsertComment("Comment12 of comment", 2, 3, 1)
	InsertComment("Comment13 of comment", 3, 3, 2)
	InsertComment("Comment14 of comment", 1, 3, 3)
	InsertComment("Comment15 of comment", 2, 1, 4)
	InsertComment("Comment16 of comment", 3, 1, 5)
	InsertComment("Comment17 of comment", 2, 6, 6)
	InsertComment("Comment18 of comment", 1, 2, 7)
	InsertComment("Comment19 of comment", 3, 3, 8)
	InsertComment("Comment20 of comment", 1, 3, 9)
	/*IF NEED MORE DATA*/
	// rand.Seed(time.Now().UnixNano())
	// count := 0
	// for count <= 20 {
	// 	InsertPost("Phasellus placerat vulputate quam. ", "Lorem ipsum dolor sit amet. Et totam consequatur aut molestiae magni ut omnis tempore. Qui maiores fuga rem magni doloribus ut voluptate harum aut assumenda galisum aut ipsa minus. Qui suscipit sapiente cum recusandae quae id rerum quos sed odit illum. Et nulla odio ut nemo dignissimos ut facilis veritatis.", 0, rand.Intn(5-2)+2, rand.Intn(4-1)+1)
	// 	InsertPost("Phasellus placerat vulputate quam. ", "Lorem ipsum dolor sit amet. Et totam consequatur aut molestiae magni ut omnis tempore. Qui maiores fuga rem magni doloribus ut voluptate harum aut assumenda galisum aut ipsa minus. Qui suscipit sapiente cum recusandae quae id rerum quos sed odit illum. Et nulla odio ut nemo dignissimos ut facilis veritatis.", 0, rand.Intn(9-6)+6, rand.Intn(4-1)+1)
	// 	count++
	// }
	// count = 0
	// for count <= 20 {
	// 	InsertVote(1, rand.Intn(52-1)+1, 0, rand.Intn(4-1)+1)
	// 	InsertVote(1, 0, rand.Intn(34-1)+1, rand.Intn(4-1)+1)
	// 	count++
	// }
	InsertTicket("Un titre", "fzefezfezfezfezfze", 0, 1)
	InsertTicket("Un titre2", "fzefezfezfezfezfze", 1, 1)
	InsertTicket("Un titre3", "fzefezfezfezfezfze", 0, 1)
	InsertTicket("Un titre4", "fzefezfezfezfezfze", 1, 2)
	InsertTicket("Un titre5", "fzefezfezfezfezfze", 0, 1)
	InsertTicket("Un titre", "fzefezfezfezfezfze", 0, 2)
	InsertTicket("Un titre2", "fzefezfezfezfezfze", 1, 2)
	InsertTicket("Un titre3", "fzefezfezfezfezfze", 0, 3)
	InsertTicket("Un titre4", "fzefezfezfezfezfze", 1, 4)
	InsertTicket("Un titre5", "fzefezfezfezfezfze", 0, 2)

	InsertTicketMessage("Ceci est le com f", 1, 1)
	InsertTicketMessage("Ceci est le com fez", 2, 2)
	InsertTicketMessage("Ceci est le com fezg", 2, 1)
	InsertTicketMessage("Ceci est le com htr", 1, 3)
	InsertTicketMessage("Ceci est le com qsf", 2, 3)
	InsertTicketMessage("Ceci est le com hrtt", 1, 3)
	InsertTicketMessage("Ceci est le com leng", 2, 1)
	InsertTicketMessage("Ceci est le com aof", 1, 4)
	InsertTicketMessage("Ceci est le com qjfjr", 2, 5)
	InsertCategorieUserRank(1, 1)
	InsertCategorieUserRank(2, 2)
	InsertCategorieUserRank(2, 3)
	InsertCategorieUserRank(1, 3)
	InsertCategorieUserRank(1, 4)
	InsertCategorieUserRank(2, 6)
	InsertCategorieUserRank(1, 4)

}

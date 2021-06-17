package main

/***************************/
/*Display        Categories*/
/***************************/

type CategoryPage struct {
	TableCategory []Category
	IsConnected   bool
	User          UserInfo
}

type PostListPage struct {
	IsConnected  bool
	Page         string
	User         UserInfo
	IsAuthorized bool
	Promote      PromoteContent
	TablePost    []Post
}

type PromoteContent struct {
	IsPromoted bool
	Type       int
	Post       Post
	Comment    Comment
}

// Fill the data in the category page struct for index handler.
func DisplayCategory() CategoryPage {
	var page CategoryPage
	page.TableCategory = GetCategories() // voir category.go
	return page
}

/**************************/
/*Display PostPresentation*/
/**************************/

//Fill the data in the postlistpage struct for postpresentation handler.
func DisplayPostList(categoryId int) PostListPage {
	var page PostListPage
	page.TablePost = GetPostsByCategory(categoryId)
	page.RedeemComsVotesCatnameInPostList()
	var promotecontentid int
	page.Promote.IsPromoted, page.Promote.Type, promotecontentid = GetPromotedContent(categoryId)
	if page.Promote.IsPromoted {
		if page.Promote.Type == 1 {
			page.Promote.Post = GetPost(promotecontentid)
			page.Promote.Post.RedeemComsInPost()
		} else {
			page.Promote.Comment = GetComment(promotecontentid)
			page.Promote.Comment.RedeemParentCommentAndOwnLikesData()
		}
	}
	return page
}

// Insert coms + addiotionnal informations in postlistpage (post presentation) struct
func (u *PostListPage) RedeemComsVotesCatnameInPostList() {
	for i, post := range u.TablePost {
		u.TablePost[i].Comment = GetCommentsForPost(post.Id)
		u.TablePost[i].NumberOfComments = len(GetCommentsForPost(post.Id))
		u.TablePost[i].Votes = CountVotes(GetVotesByPost(post.Id))
		u.TablePost[i].CategoryTitle = GetCategory(post.CategoryId).Title
		u.TablePost[i].UserName = GetUser(post.UserId).Username
		u.TablePost[i].Avatar = NumberToPpIcon(GetUser(post.UserId).Avatar)
		u.TablePost[i].BadgesCreator = GetUserBadges(post.UserId)
	}
}

/**************************/
/*Display             Post*/
/**************************/

// Fill the post  in the post struct for post handler.
func DisplayPost(postId int) Post {
	page := GetPost(postId)
	page.RedeemComsInPost()
	return page
}

// Insert of coms+additionnal informations in post struct
func (u *Post) RedeemComsInPost() {
	u.Comment = GetCommentsForPost(u.Id)       // recup du tableau de commentaires
	u.NumberOfComments = len(u.Comment)        // recup le nombre de commentaires
	u.Votes = CountVotes(GetVotesByPost(u.Id)) // recup nombre de likes
	u.UserName = GetUser(u.UserId).Username
	u.Avatar = NumberToPpIcon(GetUser(u.UserId).Avatar)
	u.CreationDate = u.CreationDate[0:16]
	RegisterDate := GetUser(u.UserId).RegisterDate
	u.RegisterDate = RegisterDate[0:10]
	u.UserRank = RankIntToRankString(GetUser(u.UserId).Rank)
	u.BadgesCreator = GetUserBadges(u.UserId)
	for i := range u.Comment {
		u.Comment[i].RedeemParentCommentAndOwnLikesData() // recup des données du post parent + les likes associés au commentaire
	}
}

// Insert of informations in comment struct
func (u *Comment) RedeemParentCommentAndOwnLikesData() {
	if u.ParentId != 0 { // si post parent on recup le commentaire, l'id du createur, et la date de creation
		u.ParentComment, u.ParentUserId, u.ParentCreationDate = GetComment(u.ParentId).Comment, GetComment(u.ParentId).UserId, GetComment(u.ParentId).CreationDate
	}
	u.Votes = CountVotes(GetVotesByComment(u.Id))
	u.UserName = GetUser(u.UserId).Username
	u.CreationDate = u.CreationDate[0:16]
	u.Avatar = NumberToPpIcon(GetUser(u.UserId).Avatar)
	RegisterDate := GetUser(u.UserId).RegisterDate
	u.RegisterDate = RegisterDate[0:10]
	u.UserRank = RankIntToRankString(GetUser(u.UserId).Rank)
	u.BadgesCreator = GetUserBadges(u.UserId)

}

// return the translation of rank in int to string to display it.
func RankIntToRankString(rank int) string {
	if rank == 1 {
		return "User"
	} else if rank == 2 {
		return "Moderator"
	} else if rank == 3 {
		return "Administrator"
	} else {
		return "error"
	}
}

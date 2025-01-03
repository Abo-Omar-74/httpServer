package model

import (
	"strings"
	"time"

	"github.com/Abo-Omar-74/httpServer/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsPremium bool `json:"is_premium"`
}

func DatabaseUserToUser(dbUser database.User) User{
	return User{
		ID : dbUser.ID , 
		CreatedAt: dbUser.CreatedAt , 
		UpdatedAt: dbUser.UpdatedAt , 
		Email: dbUser.Email , 
		IsPremium: dbUser.IsPremium,
	}
}

type Post struct{
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
  Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

func DatabasePostToPost(dbPost database.Post) Post{
	return Post{
		dbPost.ID,
		dbPost.CreatedAt,
		dbPost.UpdatedAt,
		dbPost.Body,
		dbPost.UserID,
	}
}


func DatabasePostsToPosts (dbPosts []database.Post , sortParam string) []Post{
	var posts []Post
	if strings.ToUpper(sortParam) == "DESC"{
		for i:=len(dbPosts)-1;i>=0;i--{
			posts = append(posts , DatabasePostToPost(dbPosts[i]))
		}
	}else {
		for _ , post := range dbPosts{
			posts = append(posts , DatabasePostToPost(post))
		}
	}
	return posts
}
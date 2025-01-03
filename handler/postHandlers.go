package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Abo-Omar-74/httpServer/helper"
	"github.com/Abo-Omar-74/httpServer/internal/database"
	"github.com/Abo-Omar-74/httpServer/model"
	"github.com/google/uuid"
)

type postResponse struct{
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
  Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

// PostHandler creates a new post for the user.
func (h *Handler)PostHandler(w http.ResponseWriter , r *http.Request , jwtUserID uuid.UUID){
	type parameters struct{
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	var params parameters
  
	if r.Method != http.MethodPost{
		helper.RespondWithError(w,http.StatusMethodNotAllowed , "Only POST requests are allowed")
		return
	}
	


	body , err := io.ReadAll(r.Body)
	if err != nil{
		helper.RespondWithError(w,http.StatusBadRequest , "Failed to read request body")
		return
  }
	err = json.Unmarshal(body , &params)
	if err != nil{
		helper.RespondWithError(w,http.StatusBadRequest , "Invalid JSON payload")
    return
  }

  // decoder := json.NewDecoder(r.Body)
  // err := decoder.Decode(&post)

	_ , err = h.Cfg.Db.FindUserByID(r.Context() , params.UserID)

	if err != nil {
		helper.RespondWithError(w,http.StatusNotFound , "User not found")
		return
	}

	dbPost , err := h.Cfg.Db.CreatePost(r.Context() , database.CreatePostParams{Body: params.Body , UserID: params.UserID})

	if err != nil {
		helper.RespondWithError(w , http.StatusInternalServerError , "Failed to create post")
		return
	}
	helper.RespondWithJSON(w,http.StatusOK , model.DatabasePostToPost(dbPost))
}


// GetPostsHandler retrieves posts, optionally filtered by author ID.
func (h *Handler) GetPostsHandler(w http.ResponseWriter , r *http.Request){
	if r.Method != http.MethodGet{
		helper.RespondWithError(w,http.StatusMethodNotAllowed , "Only GET requests are allowed")
		return
	}
	
	authorIdStr := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort")

	if sortParam == ""{
		sortParam = "ASC"
	}
	sortParam = strings.ToUpper(sortParam)
	if sortParam != "ASC" && sortParam != "DESC"{
		helper.RespondWithError(w , http.StatusBadRequest , "Sort parameter must be either 'ASC' or 'DESC'")
		return
	}

	var dbPosts []database.Post
	var err error


	if authorIdStr == ""{
		dbPosts , err = h.Cfg.Db.GetAllPosts(r.Context())
	}else {
		var authorId uuid.UUID
		authorId , err = uuid.Parse(authorIdStr)
		if err != nil{
			helper.RespondWithError(w , http.StatusBadRequest , "Invalid request parameters.")
			return
		}
		dbPosts , err = h.Cfg.Db.GetPostsByAuthorID(r.Context() , authorId)
	}
	if err != nil{
		helper.RespondWithError(w , http.StatusInternalServerError , "Failed to fetch posts.")
		return
	}
	
	posts := model.DatabasePostsToPosts(dbPosts , sortParam)
	
	helper.RespondWithJSON(w,http.StatusOK , posts)
}


// GetPostByIDHandler retrieves a post by its ID from the database.
func (h *Handler) GetPostByIDHandler(w http.ResponseWriter , r *http.Request){

	if r.Method != http.MethodGet{
		helper.RespondWithError(w,http.StatusMethodNotAllowed , "Only GET requests are allowed")
		return
	}
	
	idString := r.PathValue("postID")
	if idString  == ""{
		helper.RespondWithError(w,http.StatusBadRequest , "Missing or invalid 'postID' in the request.")
		return
	}

	id , err := uuid.Parse(idString)
	if err != nil {
		helper.RespondWithError(w,http.StatusBadRequest , "Invalid request parameters.")
		return 
	}

	post , err := h.Cfg.Db.GetPost(r.Context() , id)
	if err != nil{
		if errors.Is(err , sql.ErrNoRows){
			helper.RespondWithError(w , http.StatusNotFound , "Post not found.")
		}else {
			helper.RespondWithError(w, http.StatusInternalServerError, "An unexpected error occurred.")
		}
		return
  }

	helper.RespondWithJSON(w,http.StatusOK , postResponse{
		post.ID , 
		post.CreatedAt , 
		post.UpdatedAt, 
		post.Body,
		post.UserID,
	})
}

// DeletePostHandler deletes a post by ID, ensuring the user is authorized.
func (h *Handler) DeletePostHandler(w http.ResponseWriter , r *http.Request , jwtUserID uuid.UUID){
	
	if r.Method != http.MethodDelete{
		helper.RespondWithError(w,http.StatusMethodNotAllowed , "Only DELETE requests are supported.")
		return
	}
	idString := r.PathValue("postID")
	if idString  == ""{
		helper.RespondWithError(w,http.StatusBadRequest , "Post ID is required.")
		return
	}

	id , err := uuid.Parse(idString)
	if err != nil {
		helper.RespondWithError(w,http.StatusBadRequest , "Invalid request parameters.")
		return 
	}
	post, err := h.Cfg.Db.GetPost(r.Context() , id)
	if err != nil{
		helper.RespondWithError(w,http.StatusNotFound , "Post not found.")
		return
  }
	if jwtUserID != post.UserID{
		helper.RespondWithError(w,http.StatusForbidden ,"You are not allowed to delete this post.")
		return
	}
	_ , err = h.Cfg.Db.DeletePost(r.Context() , post.ID)
	if err != nil{
		helper.RespondWithError(w , http.StatusInternalServerError , "Unable to delete post.")
		return 
	}
	helper.RespondWithJSON(w,http.StatusNoContent , nil)
}

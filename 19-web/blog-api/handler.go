package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type postHandler struct{ service *PostService }

func newPostHandler(service *PostService) *postHandler {
	return &postHandler{service: service}
}

func (h *postHandler) list(c *gin.Context) {
	posts, err := h.service.ListPosts(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, "query posts")
		return
	}
	response := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		response = append(response, toPostResponse(post))
	}
	c.JSON(http.StatusOK, response)
}

func (h *postHandler) create(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "title and body are required")
		return
	}
	post, err := h.service.CreatePost(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidPost) {
			writeError(c, http.StatusBadRequest, "title and body are required")
			return
		}
		writeError(c, http.StatusInternalServerError, "create post")
		return
	}
	c.JSON(http.StatusCreated, toPostResponse(post))
}

func writeError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"github.com/lacolaco/activitypub.lacolaco.net/auth"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
)

type UserRepository interface {
	FindByLocalID(ctx context.Context, localID string) (*model.LocalUser, error)
}

type RelationshipUsecase interface {
	Follow(r *http.Request, uid model.UID, whom string) error
	Unfollow(r *http.Request, uid model.UID, whom string) error
}

type SearchUsecase interface {
	SearchPerson(ctx context.Context, id string) (*ap.Person, error)
}

type service struct {
	userRepo            UserRepository
	relationshipUsecase RelationshipUsecase
	searchUsecase       SearchUsecase
}

func New(userRepo UserRepository, relationshipUsecase RelationshipUsecase, searchUsecase SearchUsecase) *service {
	return &service{userRepo: userRepo, relationshipUsecase: relationshipUsecase, searchUsecase: searchUsecase}
}

func (s *service) RegisterRoutes(r *gin.Engine) {
	apiRoutes := r.Group("/api")
	apiRoutes.GET("/ping", s.ping)
	apiRoutes.GET("/users/show/:id", s.showUser)
	apiRoutes.GET("/search/person/:id", auth.AssertAuthenticated(), s.searchPerson)
	apiRoutes.POST("/following/create", auth.AssertAuthenticated(), s.followUser)
	apiRoutes.POST("/following/delete", auth.AssertAuthenticated(), s.unfollowUser)
}

func (s *service) ping(c *gin.Context) {
	c.JSON(200, gin.H{})
}

type showUserResp struct {
	User *model.LocalUser `json:"user"`
}

func (s *service) showUser(c *gin.Context) {
	localID := c.Param("id")
	if localID == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "id is required"})
		return
	}
	user, err := s.userRepo.FindByLocalID(c.Request.Context(), localID)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, showUserResp{User: user})
}

type searchPersonResp struct {
	Person *ap.Person `json:"person"`
}

func (s *service) searchPerson(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "id is required"})
		return
	}
	person, err := s.searchUsecase.SearchPerson(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, searchPersonResp{Person: person})
}

type followUserRequest struct {
	ID string `json:"id"`
}

func (s *service) followUser(c *gin.Context) {
	req := followUserRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	uid := auth.UIDFromContext(c.Request.Context())
	if err := s.relationshipUsecase.Follow(c.Request, uid, req.ID); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{})
}

type unfollowUserRequest struct {
	ID string `json:"id"`
}

func (s *service) unfollowUser(c *gin.Context) {
	req := unfollowUserRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	uid := auth.UIDFromContext(c.Request.Context())
	if err := s.relationshipUsecase.Unfollow(c.Request, uid, req.ID); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{})
}

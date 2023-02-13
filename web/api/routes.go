package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"github.com/lacolaco/activitypub.lacolaco.net/auth"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"github.com/lacolaco/activitypub.lacolaco.net/web/util"
	"github.com/lacolaco/activitypub.lacolaco.net/webfinger"
)

type UserRepository interface {
	FindByUID(ctx context.Context, id string) (*model.LocalUser, error)
	UpsertFollowing(ctx context.Context, user *model.LocalUser, following *model.Following) error
	DeleteFollowing(ctx context.Context, user *model.LocalUser, whom string) error
}

type JobRepository interface {
	Add(ctx context.Context, job *model.Job) error
}

type service struct {
	userRepo UserRepository
}

func New(userRepo UserRepository) *service {
	return &service{userRepo: userRepo}
}

func (s *service) Register(r *gin.Engine) {
	apiRoutes := r.Group("/api")
	apiRoutes.GET("/ping", auth.AssertAuthenticated(), s.ping)
	apiRoutes.GET("/users/search", auth.AssertAuthenticated(), s.searchUser)
	apiRoutes.POST("/users/follow", auth.AssertAuthenticated(), s.followUser)
	apiRoutes.POST("/users/unfollow", auth.AssertAuthenticated(), s.unfollowUser)
}

func (s *service) ping(c *gin.Context) {
	c.JSON(200, gin.H{})
}

type searchUserResponse struct {
	User *ap.Person `json:"user"`
}

func (s *service) searchUser(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "id is required"})
		return
	}
	personURI, err := webfinger.ResolveAccountURI(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	if personURI == "" {
		c.JSON(200, searchUserResponse{User: nil})
		return
	}
	person, err := ap.GetPerson(c.Request.Context(), personURI)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, person)
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
	currentUser := auth.CurrentUserFromContext(c.Request.Context())
	actor := ap.NewPerson(currentUser, util.GetBaseURI(c))
	whom, err := ap.GetPerson(c.Request.Context(), req.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := ap.FollowPerson(c.Request.Context(), actor, whom); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	following := model.NewFollowing(whom.GetID().String(), model.AttemptStatusPending)
	if err := s.userRepo.UpsertFollowing(c.Request.Context(), currentUser, following); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	currentUser := auth.CurrentUserFromContext(c.Request.Context())
	actor := ap.NewPerson(currentUser, util.GetBaseURI(c))
	whom, err := ap.GetPerson(c.Request.Context(), req.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := ap.UnfollowPerson(c.Request.Context(), actor, whom); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := s.userRepo.DeleteFollowing(c.Request.Context(), currentUser, whom.GetID().String()); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
	c.JSON(http.StatusAccepted, gin.H{})
}

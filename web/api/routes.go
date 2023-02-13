package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	goap "github.com/go-ap/activitypub"
	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"github.com/lacolaco/activitypub.lacolaco.net/auth"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"github.com/lacolaco/activitypub.lacolaco.net/web/util"
	"github.com/lacolaco/activitypub.lacolaco.net/webfinger"
)

type UserRepository interface {
	FindByLocalID(ctx context.Context, localID string) (*model.LocalUser, error)
	UpsertFollowing(ctx context.Context, user *model.LocalUser, following *model.Following) error
	DeleteFollowing(ctx context.Context, user *model.LocalUser, whom string) error
}

type service struct {
	userRepo UserRepository
}

func New(userRepo UserRepository) *service {
	return &service{userRepo: userRepo}
}

func (s *service) RegisterRoutes(r *gin.Engine) {
	apiRoutes := r.Group("/api")
	apiRoutes.GET("/ping", s.ping)
	apiRoutes.GET("/users/show/:id", s.showUser)
	apiRoutes.GET("/users/search/:id", auth.AssertAuthenticated(), s.searchUser)
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

type searchUserResp struct {
	User *goap.Person `json:"user"`
}

func (s *service) searchUser(c *gin.Context) {
	id := c.Param("id")
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
		c.JSON(200, searchUserResp{User: nil})
		return
	}
	person, err := ap.GetPerson(c.Request.Context(), personURI)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, searchUserResp{User: person})
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

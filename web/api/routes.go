package api

import (
	"context"
	"net/http"
	"strings"

	fbauth "firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/lacolaco/activitypub.lacolaco.net/ap"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
	"github.com/lacolaco/activitypub.lacolaco.net/web/middleware/auth"
	"github.com/lacolaco/activitypub.lacolaco.net/web/utils"
	"github.com/lacolaco/activitypub.lacolaco.net/webfinger"
)

type UserRepository interface {
	FindByUID(ctx context.Context, id string) (*model.LocalUser, error)
}

type JobRepository interface {
	Add(ctx context.Context, job *model.Job) error
}

type service struct {
	authClient *fbauth.Client
	userRepo   UserRepository
	jobRepo    JobRepository
}

func New(authClient *fbauth.Client, userRepo UserRepository, jobRepo JobRepository) *service {
	return &service{authClient: authClient, userRepo: userRepo, jobRepo: jobRepo}
}

func (s *service) Register(r *gin.Engine) {

	apiRoutes := r.Group("/api", auth.Authenticate(s.authClient, s.userRepo.FindByUID))
	apiRoutes.GET("/ping", auth.AssertAuthenticated(), s.ping)
	apiRoutes.GET("/users/search", auth.AssertAuthenticated(), s.searchUser)
	apiRoutes.POST("/users/follow", auth.AssertAuthenticated(), s.followUser)
	apiRoutes.POST("/users/unfollow", auth.AssertAuthenticated(), s.unfollowUser)
}

func (s *service) ping(c *gin.Context) {
	if !strings.HasPrefix(c.Request.Header.Get("Authorization"), "Bearer ") {
		c.AbortWithStatus(401)
		return
	}
	idToken := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
	_, err := s.authClient.VerifyIDToken(c.Request.Context(), idToken)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}
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
	currentUser := auth.FromContext(c.Request.Context())
	actor := ap.NewPerson(currentUser, utils.GetBaseURI(c))
	job, err := ap.FollowPerson(c.Request.Context(), actor, req.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := s.jobRepo.Add(c.Request.Context(), job); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{})
}

func (s *service) unfollowUser(c *gin.Context) {
	c.AbortWithStatus(501)
}

package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
)

const (
	SendFriendRequestError        = "SendFriendRequestError Handler: "
	UpdateFriendshipError         = "UpdateFriendshipError Handler: "
	GetFriendsError               = "GetFriendsError Handler: "
	GetFriendStatsError           = "GetFriendStatsError Handler: "
	GetPendingFriendRequestsError = "GetPendingFriendRequestsError Handler: "
)

type friendshipService interface {
	SendFriendRequest(userID, friendID uint) (*request.FriendshipOutputDTO, error)
	UpdateFriendshipStatus(id uint, status string) error
	GetFriends(userID uint) ([]*request.UserOutputDTO, error)
	GetFriendsStatsSortedByCompletedGoals(userID uint) (*request.UserAndFriendsStatsResponse, error)
	GetPendingFriendRequests(userID uint) ([]*request.FriendshipOutputDTO, error)
}

type FriendshipHandler struct {
	friendshipService friendshipService
}

func NewFriendshipHandler(friendshipService friendshipService) *FriendshipHandler {
	return &FriendshipHandler{
		friendshipService: friendshipService,
	}
}

func (f *FriendshipHandler) AddRoutes(router *echo.Router) {
	router.Add(echo.POST, "users/:userID/friends", f.sendFriendRequest)
	router.Add(echo.PUT, "friendships/:id/status", f.updateFriendshipStatus)
	router.Add(echo.GET, "users/:userID/friends", f.getFriends)
	router.Add(echo.GET, "users/:userID/friends/stats", f.getFriendStats)
	router.Add(echo.GET, "users/:userID/friend-requests/pending", f.getPendingFriendRequests)
}

func (f *FriendshipHandler) sendFriendRequest(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(SendFriendRequestError+"invalid user ID: "+err.Error()).Error())
	}

	r := new(request.SendFriendRequest)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(SendFriendRequestError+"invalid request body: "+err.Error()).Error())
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(SendFriendRequestError+err.Error()).Error())
	}

	friendshipOutput, err := f.friendshipService.SendFriendRequest(uint(userID), r.FriendID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"layer": SendFriendRequestError,
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, friendshipOutput)
}

func (f *FriendshipHandler) updateFriendshipStatus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(UpdateFriendshipError+"invalid friendship ID: "+err.Error()).Error())
	}

	r := new(request.UpdateFriendshipStatus)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(UpdateFriendshipError+"invalid request body: "+err.Error()).Error())
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(UpdateFriendshipError+err.Error()).Error())
	}

	err = f.friendshipService.UpdateFriendshipStatus(uint(id), r.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"layer": UpdateFriendshipError,
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (f *FriendshipHandler) getFriends(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(GetFriendsError+"invalid user ID: "+err.Error()).Error())
	}

	friendsOutput, err := f.friendshipService.GetFriends(uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"layer": GetFriendsError,
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, friendsOutput)
}

func (f *FriendshipHandler) getFriendStats(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(GetFriendStatsError+"invalid user ID: "+err.Error()).Error())
	}

	response, err := f.friendshipService.GetFriendsStatsSortedByCompletedGoals(uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"layer": GetFriendStatsError,
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (f *FriendshipHandler) getPendingFriendRequests(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(GetPendingFriendRequestsError+"invalid user ID: "+err.Error()).Error())
	}

	pendingRequestsOutput, err := f.friendshipService.GetPendingFriendRequests(uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"layer": GetPendingFriendRequestsError,
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, pendingRequestsOutput)
}

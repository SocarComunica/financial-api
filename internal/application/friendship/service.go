package friendship

import (
	"errors"
	"sort"

	"github.com/socarcomunica/financial-api/internal/adapters/producer/http/request"
	"github.com/socarcomunica/financial-api/internal/domain"
)

// Helper function to convert domain.User to request.UserOutputDTO
func toUserOutputDTO(user *domain.User) *request.UserOutputDTO {
	if user == nil {
		return nil
	}
	return &request.UserOutputDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

// UserGetter defines an interface for fetching a user by ID.
// This is to avoid circular dependencies if Database needs to fetch users.
// Alternatively, the Database interface could have a GetUserByID method.
type UserGetter interface {
	GetUserByID(id uint) (*domain.User, error)
}

const (
	CreateFriendshipServiceError         = "CreateFriendship Service Error: "
	UpdateFriendshipServiceError         = "UpdateFriendship Service Error: "
	GetFriendsServiceError               = "GetFriends Service Error: "
	GetFriendStatsServiceError           = "GetFriendStats Service Error: "
	GetPendingFriendRequestsServiceError = "GetPendingFriendRequests Service Error: "
	UserServiceError                     = "User Service Error: " // For fetching user details
)

type FriendshipDatabase interface {
	CreateFriendship(user1ID, user2ID uint) (*domain.Friendship, error)
	UpdateFriendshipStatus(id uint, status domain.FriendshipStatus) error
	GetFriends(userID uint) ([]*domain.User, error)                     // Returns domain users, service will map
	GetFriendStats(userID uint) (*request.FriendStats, error)           // Already returns desired stats struct (now with email)
	GetPendingFriendRequests(userID uint) ([]*domain.Friendship, error) // Returns domain friendships, service will map
	GetUserByID(id uint) (*domain.User, error)                          // Added to fetch user details
}

type Service struct {
	Database FriendshipDatabase
	// UserGetter UserGetter // If using a separate user service/getter
}

func NewFriendshipService(database FriendshipDatabase) *Service {
	return &Service{
		Database: database,
	}
}

func (s *Service) SendFriendRequest(userID, friendID uint) (*request.FriendshipOutputDTO, error) {
	if userID == friendID {
		return nil, errors.New(CreateFriendshipServiceError + "cannot send friend request to yourself")
	}

	friendship, err := s.Database.CreateFriendship(userID, friendID)
	if err != nil {
		return nil, errors.New(CreateFriendshipServiceError + err.Error())
	}

	user1, err := s.Database.GetUserByID(friendship.User1ID)
	if err != nil {
		return nil, errors.New(UserServiceError + "failed to get details for user1: " + err.Error())
	}
	user2, err := s.Database.GetUserByID(friendship.User2ID)
	if err != nil {
		return nil, errors.New(UserServiceError + "failed to get details for user2: " + err.Error())
	}

	return &request.FriendshipOutputDTO{
		ID:        friendship.ID,
		User1:     toUserOutputDTO(user1),
		User2:     toUserOutputDTO(user2),
		Status:    string(friendship.Status),
		CreatedAt: friendship.CreatedAt,
		UpdatedAt: friendship.UpdatedAt,
	}, nil
}

func (s *Service) UpdateFriendshipStatus(id uint, status string) error {
	friendshipStatus := domain.FriendshipStatus(status)
	if friendshipStatus != domain.FriendshipStatusAccepted && friendshipStatus != domain.FriendshipStatusRejected {
		return errors.New(UpdateFriendshipServiceError + "invalid status")
	}

	err := s.Database.UpdateFriendshipStatus(id, friendshipStatus)
	if err != nil {
		return errors.New(UpdateFriendshipServiceError + err.Error())
	}
	return nil
}

func (s *Service) GetFriends(userID uint) ([]*request.UserOutputDTO, error) {
	domainFriends, err := s.Database.GetFriends(userID)
	if err != nil {
		return nil, errors.New(GetFriendsServiceError + err.Error())
	}

	outputDTOs := make([]*request.UserOutputDTO, len(domainFriends))
	for i, friend := range domainFriends {
		outputDTOs[i] = toUserOutputDTO(friend)
	}
	if outputDTOs == nil {
		return []*request.UserOutputDTO{}, nil
	}
	return outputDTOs, nil
}

func (s *Service) GetFriendStats(userID uint) (*request.FriendStats, error) {
	stats, err := s.Database.GetFriendStats(userID)
	if err != nil {
		return nil, errors.New(GetFriendStatsServiceError + err.Error())
	}
	return stats, nil
}

func (s *Service) GetFriendsStatsSortedByCompletedGoals(userID uint) (*request.UserAndFriendsStatsResponse, error) {
	userStats, err := s.Database.GetFriendStats(userID)
	if err != nil {
		return nil, errors.New(GetFriendStatsServiceError + "failed to get user stats: " + err.Error())
	}

	domainFriends, err := s.Database.GetFriends(userID)
	if err != nil {
		return nil, errors.New(GetFriendsServiceError + "failed to get friends for stats: " + err.Error())
	}

	var friendsStatsList []*request.FriendStats
	if len(domainFriends) > 0 {
		for _, friend := range domainFriends {
			stats, err := s.Database.GetFriendStats(friend.ID)
			if err != nil {
				return nil, errors.New(GetFriendStatsServiceError + "failed to get stats for friend " + friend.Username + ": " + err.Error())
			}
			if stats != nil {
				friendsStatsList = append(friendsStatsList, stats)
			}
		}
		sort.Slice(friendsStatsList, func(i, j int) bool {
			return friendsStatsList[i].CompletedGoals > friendsStatsList[j].CompletedGoals
		})
	} else {
		friendsStatsList = []*request.FriendStats{}
	}

	return &request.UserAndFriendsStatsResponse{
		UserStats:    userStats,
		FriendsStats: friendsStatsList,
	}, nil
}

func (s *Service) GetPendingFriendRequests(userID uint) ([]*request.FriendshipOutputDTO, error) {
	domainRequests, err := s.Database.GetPendingFriendRequests(userID)
	if err != nil {
		return nil, errors.New(GetPendingFriendRequestsServiceError + err.Error())
	}

	outputDTOs := make([]*request.FriendshipOutputDTO, 0, len(domainRequests))
	for _, fr := range domainRequests {
		var senderDTO *request.UserOutputDTO
		// Assuming fr.User1 is of type domain.User (not a pointer) and populated by Joins if User1ID > 0
		if fr.User1ID != 0 && fr.User1.ID == fr.User1ID { // Check if User1 was preloaded and its ID matches User1ID
			senderDTO = toUserOutputDTO(&fr.User1) // Pass address of fr.User1
		} else if fr.User1ID != 0 { // Fallback: if User1 was not preloaded but User1ID exists
			sender, err_db := s.Database.GetUserByID(fr.User1ID)
			if err_db != nil {
				// Log this event
				continue
			}
			senderDTO = toUserOutputDTO(sender)
		} else {
			// User1ID is zero, data integrity issue or unexpected case
			continue
		}

		// User2 is the recipient
		recipient, err_db := s.Database.GetUserByID(fr.User2ID)
		if err_db != nil {
			// Log error
			continue
		}

		outputDTOs = append(outputDTOs, &request.FriendshipOutputDTO{
			ID:        fr.ID,
			User1:     senderDTO,
			User2:     toUserOutputDTO(recipient),
			Status:    string(fr.Status),
			CreatedAt: fr.CreatedAt,
			UpdatedAt: fr.UpdatedAt,
		})
	}
	return outputDTOs, nil
}

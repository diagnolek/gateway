package grpc

import (
	"context"
	"errors"
	"gateway/pkg/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

func DefaultGRPCUserServer(service user.UserService) UserServiceServer {
	return &userGrpcServer{
		service:       service,
		subscriptions: map[uint]*UserService_SubscribeToChatServer{},
	}
}

type userGrpcServer struct {
	service       user.UserService
	subscriptions map[uint]*UserService_SubscribeToChatServer
}

func (u *userGrpcServer) CreateUser(ctx context.Context, request *CreateUserRequest) (*CreateUserResponse, error) {

	createUser := user.User{}
	createUser.Name = request.Name
	createUser.Email = request.Email
	createUser.Lastname = request.LastName
	createUser.Password = request.Password

	err := u.service.CreateUser(createUser)
	if err != nil {
		return nil, err
	}
	return &CreateUserResponse{}, nil
}

func (u *userGrpcServer) LoginUser(ctx context.Context, request *LoginUserRequest) (*LoginUserResponse, error) {
	accessToken, err := u.service.Login(request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	response := LoginUserResponse{
		AccessToken: accessToken,
	}
	return &response, nil
}

func (u *userGrpcServer) GetUserInfo(ctx context.Context, request *GetUserInfoRequest) (*GetUserInfoResponse, error) {
	userId, err := u.authorize(ctx)
	if err != nil {
		return nil, err
	}
	userInfo, err := u.service.GetUserInfo(userId)
	if err != nil {
		return nil, err
	}
	response := GetUserInfoResponse{
		Name:     userInfo.Name,
		LastName: userInfo.Lastname,
		Email:    userInfo.Email,
	}
	return &response, nil
}

func (u *userGrpcServer) authorize(ctx context.Context) (uint, error) {
	data, success := metadata.FromIncomingContext(ctx)
	unauthorized := errors.New("error-unauthorized")
	if !success {
		return 0, unauthorized
	}
	accessToken := data.Get("authorization")
	if len(accessToken) == 0 || accessToken[0] == "" {
		return 0, unauthorized
	}
	userId, err2 := u.service.Authorize(accessToken[0])
	if err2 != nil {
		return 0, err2
	}

	return userId, nil
}

// Chat
func (u *userGrpcServer) SendMessage(ctx context.Context, request *SendMessageRequest) (*emptypb.Empty, error) {
	userId, err := u.authorize(ctx)
	if err != nil {
		return nil, err
	}

	toSubscription := u.subscriptions[uint(request.To)]

	message := Message{
		From: int64(userId),
		To:   request.To,
		Text: request.Text,
		Date: time.Now().Unix(),
	}

	messages := StreamMessage_Messages{
		Messages: &Messages{
			Messages: []*Message{&message},
		},
	}

	streamMessage := &StreamMessage{
		Body: &messages,
	}

	fromSubscription := u.subscriptions[userId]

	if toSubscription != nil {
		_ = (*toSubscription).Send(streamMessage)
	}
	if fromSubscription != nil && userId != uint(request.To) {
		_ = (*fromSubscription).Send(streamMessage)
	}

	userMessage := user.Message{
		From: uint(message.From),
		To:   uint(message.To),
		Date: int(message.Date),
		Text: message.Text,
	}

	err = u.service.SendMessage(userMessage)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (u *userGrpcServer) SubscribeToChat(empty *emptypb.Empty, subscription grpc.ServerStreamingServer[StreamMessage]) error {
	userId, err := u.authorize(subscription.Context())
	if err != nil {
		return err
	}
	if _, ok := u.subscriptions[userId]; ok {
		return errors.New("error-already-subscribed")
	}
	u.subscriptions[userId] = &subscription

	_ = u.sendListAllAvailableUsers()

	<-subscription.Context().Done()
	delete(u.subscriptions, userId)

	return nil
}

func (u *userGrpcServer) GetMessageHistory(ctx context.Context, request *GetMessageHistoryRequest) (*Messages, error) {
	fromUserId, err := u.authorize(ctx)
	if err != nil {
		return nil, err
	}
	toUserId := request.UserId

	messages, err := u.service.GetMessageHistory(fromUserId, uint(toUserId))
	if err != nil {
		return nil, err
	}
	response := Messages{}
	for _, message := range messages {
		response.Messages = append(response.Messages, &Message{
			From: int64(message.From),
			To:   int64(message.To),
			Date: int64(message.Date),
			Text: message.Text,
		})
	}

	return &response, nil
}

func (u *userGrpcServer) mustEmbedUnimplementedUserServiceServer() {

}

func (u *userGrpcServer) getListOfAvailableUsers() []*User {
	users := []*User{}
	for key := range u.subscriptions {
		userInfo, err := u.service.GetUserInfo(key)
		if err != nil {
			continue
		}
		mappedUser := User{
			UserId:   int64(key),
			Name:     userInfo.Name,
			LastName: userInfo.Lastname,
		}
		users = append(users, &mappedUser)
	}
	return users
}

func (u *userGrpcServer) sendListAllAvailableUsers() error {
	users := u.getListOfAvailableUsers()
	message := StreamMessage_List{
		List: &UserList{
			User: users,
		},
	}
	for _, subscription := range u.subscriptions {
		_ = (*subscription).Send(&StreamMessage{
			Body: &message,
		})

	}
	return nil
}

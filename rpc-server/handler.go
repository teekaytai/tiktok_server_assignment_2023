package main

import (
	"context"
	"fmt"
	"github.com/teekaytai/tiktok_server_assignment_2023/rpc-server/kitex_gen/rpc"
	"strings"
	"time"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct {
	db DbClient
}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	if err := validateSendRequest(req); err != nil {
		return nil, err
	}

	message := &Message{
		Sender:    req.Message.GetSender(),
		Text:      req.Message.GetText(),
		Timestamp: time.Now().Unix(),
	}

	roomId, err := chatToRoomId(req.Message.GetChat())
	if err != nil {
		return nil, err
	}

	err = s.db.SaveMessage(ctx, roomId, message)
	if err != nil {
		return nil, err
	}

	resp := rpc.NewSendResponse()
	resp.Code, resp.Msg = 0, "success"
	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	roomId, err := chatToRoomId(req.GetChat())
	if err != nil {
		return nil, err
	}

	start := req.GetCursor()
	limit := req.GetLimit()
	end := start + int64(limit) // Get limit + 1 messages to check if more messages available

	messages, err := s.db.GetMessagesByRoomId(ctx, roomId, start, end, req.GetReverse())
	if err != nil {
		return nil, err
	}

	respMessages := make([]*rpc.Message, 0)
	hasMore := false
	var nextCursor int64 = 0
	for i, msg := range messages {
		if int32(i) >= limit {
			hasMore = true
			nextCursor = end
			break
		}
		respMsg := &rpc.Message{
			Chat:     req.GetChat(),
			Text:     msg.Text,
			Sender:   msg.Sender,
			SendTime: msg.Timestamp,
		}
		respMessages = append(respMessages, respMsg)
	}

	resp := &rpc.PullResponse{
		Messages:   respMessages,
		Code:       0,
		Msg:        "success",
		HasMore:    &hasMore,
		NextCursor: &nextCursor,
	}

	return resp, nil
}

func validateSendRequest(req *rpc.SendRequest) error {
	users := strings.Split(req.Message.GetChat(), ":")
	if len(users) != 2 {
		err := fmt.Errorf("invalid chat ID '%s', should be in the format of user1:user2", req.Message.GetChat())
		return err
	}

	sender := req.Message.GetSender()
	if sender != users[0] && sender != users[1] {
		err := fmt.Errorf("sender '%s' is not in the chat room", sender)
		return err
	}

	return nil
}

func chatToRoomId(chat string) (string, error) {
	var roomId string

	users := strings.Split(strings.ToLower(chat), ":")
	if len(users) != 2 {
		err := fmt.Errorf("invalid chat ID '%s', should be in the format user1:user2", chat)
		return "", err
	}

	user1, user2 := users[0], users[1]
	if user1 > user2 {
		user1, user2 = user2, user1
	}
	roomId = fmt.Sprintf("%s:%s", user1, user2)

	return roomId, nil
}

package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/teekaytai/tiktok_server_assignment_2023/rpc-server/kitex_gen/rpc"
)

func TestIMServiceImpl_Send(t *testing.T) {
	type args struct {
		ctx context.Context
		req *rpc.SendRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success1",
			args: args{
				ctx: context.Background(),
				req: &rpc.SendRequest{
					Message: &rpc.Message{
						Chat:     "user1:user2",
						Text:     "Test text",
						Sender:   "user1",
						SendTime: time.Now().Unix(),
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "success2",
			args: args{
				ctx: context.Background(),
				req: &rpc.SendRequest{
					Message: &rpc.Message{
						Chat:     "user3:user4",
						Text:     "Blah blah blah",
						Sender:   "user4",
						SendTime: time.Now().Unix(),
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "invalidChatIdError",
			args: args{
				ctx: context.Background(),
				req: &rpc.SendRequest{
					Message: &rpc.Message{
						Chat:     "user1:user2:user3",
						Text:     "Invalid chat ID",
						Sender:   "user1",
						SendTime: time.Now().Unix(),
					},
				},
			},
			wantErr: fmt.Errorf("invalid chat ID 'user1:user2:user3', should be in the format user1:user2"),
		},
		{
			name: "senderNotInRoomError",
			args: args{
				ctx: context.Background(),
				req: &rpc.SendRequest{
					Message: &rpc.Message{
						Chat:     "user1:user2",
						Text:     "Sender not in room",
						Sender:   "userX",
						SendTime: time.Now().Unix(),
					},
				},
			},
			wantErr: fmt.Errorf("sender 'userX' is not in the chat room"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := &MockDbClient{}
			err := mockDb.InitClient(tt.args.ctx, "", "")
			s := &IMServiceImpl{mockDb}
			got, err := s.Send(tt.args.ctx, tt.args.req)
			if tt.wantErr == nil {
				assert.Nil(t, err)
			} else {
				assert.True(t, err != nil && err.Error() == tt.wantErr.Error())
			}
			assert.NotNil(t, got)
		})
	}
}

func TestIMServiceImpl_Pull(t *testing.T) {
	type args struct {
		ctx context.Context
		req *rpc.PullRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success1",
			args: args{
				ctx: context.Background(),
				req: &rpc.PullRequest{
					Chat:    "user1:user2",
					Cursor:  0,
					Limit:   10,
					Reverse: new(bool),
				},
			},
			wantErr: nil,
		},
		{
			name: "success2",
			args: args{
				ctx: context.Background(),
				req: &rpc.PullRequest{
					Chat:   "user1:user2",
					Cursor: 0,
					Limit:  1,
				},
			},
			wantErr: nil,
		},
		{
			name: "success3",
			args: args{
				ctx: context.Background(),
				req: &rpc.PullRequest{
					Chat:   "user1:user2",
					Cursor: 1,
					Limit:  2,
				},
			},
			wantErr: nil,
		},
		{
			name: "invalidChatIdError",
			args: args{
				ctx: context.Background(),
				req: &rpc.PullRequest{
					Chat:   "user1:user2:user3",
					Cursor: 0,
					Limit:  10,
				},
			},
			wantErr: fmt.Errorf("invalid chat ID 'user1:user2:user3', should be in the format user1:user2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := &MockDbClient{}
			err := mockDb.InitClient(tt.args.ctx, "", "")
			message := &Message{
				Sender:    "sender",
				Text:      "text",
				Timestamp: time.Now().Unix(),
			}
			mockDb.messagesSaved = []*Message{message, message, message}

			s := &IMServiceImpl{mockDb}
			got, err := s.Pull(tt.args.ctx, tt.args.req)
			if tt.wantErr == nil {
				assert.Nil(t, err)
			} else {
				assert.True(t, err != nil && err.Error() == tt.wantErr.Error())
			}
			assert.NotNil(t, got)
		})
	}
}

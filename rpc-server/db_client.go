package main

import "context"

type DbClient interface {
	InitClient(ctx context.Context, address, password string) error
	SaveMessage(ctx context.Context, roomId string, message *Message) error
	GetMessagesByRoomId(ctx context.Context, roomId string, start, end int64, reverse bool) ([]*Message, error)
}

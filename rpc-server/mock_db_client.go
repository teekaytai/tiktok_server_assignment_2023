package main

import "context"

type MockDbClient struct {
	messagesSaved []*Message
}

func (c MockDbClient) InitClient(_ context.Context, _, _ string) error {
	return nil
}

func (c MockDbClient) SaveMessage(_ context.Context, _ string, message *Message) error {
	c.messagesSaved = append(c.messagesSaved, message)
	return nil
}

func (c MockDbClient) GetMessagesByRoomId(_ context.Context, _ string, start, end int64, reverse bool) ([]*Message, error) {
	stop := int64(len(c.messagesSaved) - 1)
	if end < stop {
		stop = end
	}
	if reverse {
		res := make([]*Message, 0)
		for i := stop; i >= start; i-- {
			res = append(res, c.messagesSaved[i])
		}
		return res, nil
	} else {
		return c.messagesSaved[start : stop+1], nil
	}
}

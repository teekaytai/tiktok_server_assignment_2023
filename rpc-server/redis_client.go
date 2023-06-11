package main

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func (c *RedisClient) InitClient(ctx context.Context, address, password string) error {
	r := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	if err := r.Ping(ctx).Err(); err != nil {
		return err
	}

	c.client = r
	return nil
}

func (c *RedisClient) SaveMessage(ctx context.Context, roomId string, message *Message) error {
	msgJson, err := json.Marshal(message)
	if err != nil {
		return err
	}

	member := &redis.Z{
		Score:  float64(message.Timestamp),
		Member: msgJson,
	}

	_, err = c.client.ZAdd(ctx, roomId, *member).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisClient) GetMessagesByRoomId(ctx context.Context, roomId string, start, end int64, reverse bool) ([]*Message, error) {
	var (
		rawMessages []string
		messages    []*Message
		err         error
	)

	if reverse {
		// Messages in reverse chronological order
		rawMessages, err = c.client.ZRevRange(ctx, roomId, start, end).Result()
		if err != nil {
			return nil, err
		}
	} else {
		// Messages in chronological order
		rawMessages, err = c.client.ZRange(ctx, roomId, start, end).Result()
		if err != nil {
			return nil, err
		}
	}

	for _, msgJson := range rawMessages {
		msg := &Message{}
		err := json.Unmarshal([]byte(msgJson), msg)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

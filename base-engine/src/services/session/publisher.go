/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package session

import (
	"context"
	"errors"
	"sync"

	"github.com/gofrs/uuid"
	"base-engine/gen"
)

// Publisher 是 GraphQL Subscription 使用的会话事件发布契约。
type Publisher interface {
	Subscribe(ctx context.Context, sessionID string) (<-chan *gen.SessionEvent, error)
	PublishSession(sessionID string, event *gen.SessionEvent)
}

// LocalPublisher 为单 Engine 实例维护按 Session 隔离的订阅集合。
type LocalPublisher struct {
	mu          sync.Mutex
	subscribers map[string]map[string]chan *gen.SessionEvent
}

// NewPublisher 创建进程内会话事件发布器。
func NewPublisher() *LocalPublisher {
	return &LocalPublisher{subscribers: make(map[string]map[string]chan *gen.SessionEvent)}
}

// Subscribe 注册一个缓冲为一的浏览器标签订阅。
func (p *LocalPublisher) Subscribe(ctx context.Context, sessionID string) (<-chan *gen.SessionEvent, error) {
	if sessionID == "" {
		return nil, errors.New("SESSION_ID_REQUIRED")
	}
	subscriberID := uuid.Must(uuid.NewV4()).String()
	channel := make(chan *gen.SessionEvent, 1)
	p.mu.Lock()
	if p.subscribers[sessionID] == nil {
		p.subscribers[sessionID] = make(map[string]chan *gen.SessionEvent)
	}
	p.subscribers[sessionID][subscriberID] = channel
	p.mu.Unlock()
	go func() {
		<-ctx.Done()
		p.unsubscribe(sessionID, subscriberID)
	}()
	return channel, nil
}

// PublishSession 向所有标签非阻塞发布，满缓冲即合并重复终态事件。
func (p *LocalPublisher) PublishSession(sessionID string, event *gen.SessionEvent) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, channel := range p.subscribers[sessionID] {
		select {
		case channel <- event:
		default:
		}
	}
}

func (p *LocalPublisher) unsubscribe(sessionID, subscriberID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	sessionSubscribers := p.subscribers[sessionID]
	channel, exists := sessionSubscribers[subscriberID]
	if !exists {
		return
	}
	delete(sessionSubscribers, subscriberID)
	if len(sessionSubscribers) == 0 {
		delete(p.subscribers, sessionID)
	}
	close(channel)
}

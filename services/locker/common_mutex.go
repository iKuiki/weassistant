package locker

import (
	"fmt"
	redisLock "github.com/bsm/redis-lock"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"time"
)

// CommonLockerService 通用单对象锁服务
type CommonLockerService interface {
	ObtainLock(target interface{}) (locker Locker)
}

type registerLockerService struct {
	client  *redis.Client
	lockKey string
}

// MustNewCommonLockerService 务必创建一个通用单对象锁服务
func MustNewCommonLockerService(client *redis.Client, lockKey string) CommonLockerService {
	serv, err := NewCommonLockerService(client, lockKey)
	if err != nil {
		panic(errors.WithStack(err))
	}
	return serv
}

// NewCommonLockerService 创建一个通用单对象锁服务
func NewCommonLockerService(client *redis.Client, lockKey string) (serv CommonLockerService, err error) {
	if client == nil {
		err = errors.New("Redis client is nil")
		return
	}
	if lockKey == "" {
		err = errors.New("lockKey is empty")
		return
	}
	serv = &registerLockerService{
		client:  client,
		lockKey: lockKey,
	}
	return
}

// ObtainCommonLock 生成通用单对象锁
func (locker *registerLockerService) ObtainLock(target interface{}) (lock Locker) {
	lockName := fmt.Sprintf(locker.lockKey, target)
	return &CommonMutex{
		lock: redisLock.New(locker.client, lockName, &redisLock.Options{
			RetryCount:  10,
			RetryDelay:  200 * time.Millisecond,
			LockTimeout: 10 * time.Second,
		}),
	}
}

// CommonMutex 注册锁
type CommonMutex struct {
	lock *redisLock.Locker
}

// Lock 给锁上锁
func (mutex *CommonMutex) Lock() (err error) {
	if ok, e := mutex.lock.Lock(); e != nil {
		err = errors.New("Lock error: " + e.Error())
	} else if !ok {
		err = errors.New("Lock fail")
	}
	return
}

// Unlock 给锁解锁
func (mutex *CommonMutex) Unlock() (err error) {
	err = mutex.lock.Unlock()
	return
}

package locker

// Locker 分布式锁
type Locker interface {
	Lock() (err error)
	Unlock() (err error)
}

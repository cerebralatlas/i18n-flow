package service

import "sync"

// CacheMutexManager 缓存互斥锁管理器
// 用于防止缓存击穿的场景，通过互斥锁确保同一时刻只有一个请求去查询数据库
type CacheMutexManager struct {
	mutexes sync.Map
}

// NewCacheMutexManager 创建缓存互斥锁管理器
func NewCacheMutexManager() *CacheMutexManager {
	return &CacheMutexManager{}
}

// GetMutex 获取指定键的互斥锁，用于防止缓存击穿
// 如果锁不存在则创建一个新的锁
func (m *CacheMutexManager) GetMutex(key string) *sync.Mutex {
	if mutex, exists := m.mutexes.Load(key); exists {
		return mutex.(*sync.Mutex)
	}

	mutex := &sync.Mutex{}
	actual, loaded := m.mutexes.LoadOrStore(key, mutex)
	if loaded {
		// 已经有其他协程创建了锁，返回已有的
		return actual.(*sync.Mutex)
	}
	return mutex
}

// RemoveMutex 移除指定键的互斥锁
// 在操作完成后调用，以避免内存泄漏
func (m *CacheMutexManager) RemoveMutex(key string) {
	m.mutexes.Delete(key)
}

// WithLock 使用互斥锁执行操作的便捷方法
// 自动处理加锁、执行、解锁和移除锁的完整流程
func (m *CacheMutexManager) WithLock(key string, fn func()) {
	mutex := m.GetMutex(key)
	mutex.Lock()
	defer func() {
		mutex.Unlock()
		m.RemoveMutex(key)
	}()
	fn()
}

package ratelimit

import (
    "sync"
    "time"
)

// PlatformLimiter 平台限流器
type PlatformLimiter struct {
    // 不同平台的限流配置
    limits map[string]*RateLimit
    mu     sync.RWMutex
}

type RateLimit struct {
    MaxRequests int           // 最大请求数
    Window      time.Duration // 时间窗口
    Tokens      chan struct{} // 令牌桶
    Queue       chan *Task    // 任务队列
}

type Task struct {
    ID       string
    Platform string
    UserID   string
    CreateAt time.Time
    Done     chan struct{}
}

// 初始化限流器
func NewPlatformLimiter() *PlatformLimiter {
    pl := &PlatformLimiter{
        limits: make(map[string]*RateLimit),
    }
    
    // 配置各平台限制
    pl.limits["netease"] = &RateLimit{
        MaxRequests: 1000,
        Window:     time.Minute,
        Tokens:     make(chan struct{}, 1000),
        Queue:      make(chan *Task, 5000),
    }
    
    pl.limits["qqmusic"] = &RateLimit{
        MaxRequests: 500,
        Window:     time.Minute,
        Tokens:     make(chan struct{}, 500),
        Queue:      make(chan *Task, 2500),
    }
    
    // 启动令牌生成器
    for platform, limit := range pl.limits {
        go pl.tokenGenerator(platform)
    }
    
    return pl
}

// 生成令牌
func (pl *PlatformLimiter) tokenGenerator(platform string) {
    limit := pl.limits[platform]
    ticker := time.NewTicker(limit.Window / time.Duration(limit.MaxRequests))
    defer ticker.Stop()

    for range ticker.C {
        select {
        case limit.Tokens <- struct{}{}:
        default:
        }
    }
} 
package ratelimit

import (
    "time"
    "errors"
)

// TaskProcessor 任务处理器
type TaskProcessor struct {
    limiter *PlatformLimiter
}

func NewTaskProcessor(limiter *PlatformLimiter) *TaskProcessor {
    tp := &TaskProcessor{
        limiter: limiter,
    }
    
    // 启动任务处理
    for platform := range limiter.limits {
        go tp.processQueue(platform)
    }
    
    return tp
}

// 提交任务
func (tp *TaskProcessor) SubmitTask(platform, userID string) (*Task, error) {
    limit, ok := tp.limiter.limits[platform]
    if !ok {
        return nil, errors.New("unsupported platform")
    }
    
    task := &Task{
        ID:       generateID(),
        Platform: platform,
        UserID:   userID,
        CreateAt: time.Now(),
        Done:     make(chan struct{}),
    }
    
    // 尝试加入队列
    select {
    case limit.Queue <- task:
        return task, nil
    default:
        return nil, errors.New("queue is full")
    }
}

// 处理队列中的任务
func (tp *TaskProcessor) processQueue(platform string) {
    limit := tp.limiter.limits[platform]
    
    for task := range limit.Queue {
        // 等待令牌
        <-limit.Tokens
        
        // 处理任务
        go func(t *Task) {
            defer close(t.Done)
            // 实际的任务处理逻辑
        }(task)
    }
} 
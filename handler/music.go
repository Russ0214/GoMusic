package handler

import (
	"GoMusic/pkg/ratelimit"
	"GoMusic/pkg/websocket"
	"github.com/gin-gonic/gin"
)

type MusicHandler struct {
	processor *ratelimit.TaskProcessor
	notifier  *websocket.Notifier
}

func (h *MusicHandler) HandleMusicList(c *gin.Context) {
	platform := c.Query("platform")
	userID := c.GetString("user_id")
	
	// 提交任务
	task, err := h.processor.SubmitTask(platform, userID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	// 发送初始状态
	h.notifier.SendUpdate(userID, websocket.Message{
		Type:    "queue",
		TaskID:  task.ID,
		Message: "Task queued",
		ETA:     calculateETA(task),
	})
	
	// 等待任务完成
	<-task.Done
	
	// 发送完成通知
	h.notifier.SendUpdate(userID, websocket.Message{
		Type:    "complete",
		TaskID:  task.ID,
		Message: "Task completed",
	})
	
	// 返回结果
	c.JSON(200, gin.H{"status": "success"})
}

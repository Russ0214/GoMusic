<template>
  <div class="music-list">
    <!-- 任务状态显示 -->
    <div v-if="status" class="status-panel">
      <el-alert
        :title="status.message"
        :type="getAlertType(status.type)"
        :description="getDescription(status)"
        show-icon
      />
    </div>
    
    <!-- 进度条 -->
    <el-progress 
      v-if="status && status.type === 'queue'"
      :percentage="getProgress()"
      :format="formatProgress"
    />
  </div>
</template>

<script>
export default {
  data() {
    return {
      status: null,
      ws: null
    }
  },
  
  methods: {
    initWebSocket() {
      this.ws = new WebSocket('ws://your-domain/ws')
      this.ws.onmessage = this.handleMessage
    },
    
    handleMessage(event) {
      const data = JSON.parse(event.data)
      this.status = data
    },
    
    getAlertType(type) {
      switch(type) {
        case 'queue': return 'info'
        case 'processing': return 'warning'
        case 'complete': return 'success'
        default: return 'info'
      }
    },
    
    getDescription(status) {
      if (status.eta) {
        return `预计剩余时间: ${status.eta} 秒`
      }
      return ''
    },
    
    getProgress() {
      // 根据ETA计算进度
      if (!this.status || !this.status.eta) return 0
      const maxETA = 300 // 假设最大等待时间5分钟
      return Math.max(0, Math.min(100, (1 - this.status.eta / maxETA) * 100))
    },
    
    formatProgress(percentage) {
      return this.status.eta ? `${this.status.eta}s` : ''
    }
  },
  
  mounted() {
    this.initWebSocket()
  },
  
  beforeDestroy() {
    if (this.ws) {
      this.ws.close()
    }
  }
}
</script> 
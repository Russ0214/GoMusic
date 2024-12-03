package main

import (
	"fmt"
	"net/http"
	"os"

	"GoMusic/handler"
	"GoMusic/misc/log"
)

func main() {
	r := handler.NewRouter()
	
	// 获取 Vercel 分配的端口或使用默认端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Errorf("fail to run server: %v", err)
		panic(err)
	}
}

package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	go testlog()

	// 打开日志文件
	t, err := tail.TailFile("test.log", tail.Config{Follow: true})
	if err != nil {
		log.Fatal(err)
	}

	// HTTP服务器
	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		// 设置Content-Type为text/event-stream
		w.Header().Set("Content-Type", "text/event-stream")
		// 设置缓存控制头，禁用缓存
		w.Header().Set("Cache-Control", "no-cache")
		// 设置连接保持活动
		w.Header().Set("Connection", "keep-alive")

		// 实时推送日志
		for line := range t.Lines {
			// 生成一个日志事件
			event := fmt.Sprintf("data: %s\n\n", line.Text)

			// 将事件发送到客户端
			_, err := w.Write([]byte(event))
			if err != nil {
				// 发送失败，客户端可能已关闭连接
				fmt.Println("Client disconnected.")
				return
			}

			// 刷新响应流
			w.(http.Flusher).Flush()
		}
	})

	// 启动HTTP服务器
	fmt.Println("Server listening on :8899")
	log.Fatal(http.ListenAndServe(":8899", nil))
}

func testlog() {
	f, err := os.OpenFile("test.log", os.O_TRUNC|os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	for {
		f.WriteString("test log " + time.Now().String() + "\n")
		f.Sync()
		time.Sleep(time.Second * 2)
	}
}

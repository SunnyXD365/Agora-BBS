package main

import (
	"log"
	"os"

	wf "agora-backend/internal/workflow"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	temporalHost := os.Getenv("TEMPORAL_HOST")
	if temporalHost == "" {
		temporalHost = "agora-temporal:7233"
	}

	c, err := client.Dial(client.Options{HostPort: temporalHost})
	if err != nil {
		log.Fatalf("Worker 连接 Temporal 失败: %v", err)
	}
	defer c.Close()

	w := worker.New(c, wf.TaskQueueName, worker.Options{})

	// 注册 Workflow 和 Activity
	w.RegisterWorkflow(wf.TopicAuditWorkflow)
	w.RegisterActivity(wf.AuditActivity)

	log.Println("Worker 已启动，正在监听任务队列:", wf.TaskQueueName)
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("Worker 异常退出: %v", err)
	}
}

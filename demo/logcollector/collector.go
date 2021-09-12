package logcollector

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/37acoder/goqueue"
)

type logIdCtxKey struct{}

func getLogIdFromContext(ctx context.Context) (string, bool) {
	logId := ctx.Value(logIdCtxKey{})
	if logId == nil {
		return "", false
	}
	if realLogId, ok := logId.(string); ok {
		return realLogId, true
	} else {
		return "", false
	}
}

type Collector struct {
	q      goqueue.Queue
	writer io.Writer
}

func NewCollector(output io.Writer) *Collector {
	q := goqueue.NewInMemoryQueue(goqueue.Config{
		PushBlocking: true,
		PopBlocking:  true,
		MaxBuffer:    1024,
	})
	return &Collector{
		q:      q,
		writer: bufio.NewWriterSize(output, 10),
	}
}

type LogTask struct {
	LogContent string
	c          *Collector
}

func (l *LogTask) Execute(ctx context.Context) error {
	l.c.PersistLogs(ctx, l.LogContent)
	return nil
}

func (c *Collector) CtxLog(ctx context.Context, logLevel LogLevel, format string, args ...interface{}) {
	logId, ok := getLogIdFromContext(ctx)
	if !ok {
		logId = "NoLogId"
	}
	format = "LogId:" + logId + " " + format
	logContent := fmt.Sprintf(format, args...)
	err := c.q.Push(context.Background(), &LogTask{
		LogContent: logContent,
		c:          c,
	})
	if err == nil {
		return
	} else {
		fmt.Printf("output log failed:%s, log:%s\n", err, logContent)
	}
}

func (c *Collector) PersistLogs(ctx context.Context, logContent string) {
	write, err := c.writer.Write([]byte(logContent))
	if err != nil {
		fmt.Printf("persist log error: %s, write size:%d, log content: %s\n", err, write, logContent)
		return
	}
}

func (c *Collector) Start() {
	go func() {
		for {
			ctx := context.Background()
			task, err := c.q.Pop(ctx)
			if err != nil {
				continue
			}
			err = task.Execute(ctx)
			if err != nil {
				continue
			}
		}
	}()
}

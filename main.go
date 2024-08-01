package main

import (
	"context"
	"github.com/AyakuraYuki/bilibili-wallpaper/pkg/commands"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := notifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	commands.Cmd.Version = "2.0.0"
	if err := commands.Cmd.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}

// notifyContext 将信号绑定到上下文
func notifyContext(ctx context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan os.Signal, 5)
	signal.Notify(ch, signals...)
	if ctx.Err() == nil {
		go func() {
			// 第一次取消上下文
			select {
			case <-ctx.Done():
			case <-ch:
				cancel()
			}
			// 第二次结束程序
			select {
			case <-ctx.Done():
			case <-ch:
				os.Exit(1)
			}
		}()
	}
	return ctx, cancel
}

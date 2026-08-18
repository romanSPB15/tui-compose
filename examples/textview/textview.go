package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/extra"
)

func generateLog() string {
	now := time.Now().Format("15:04:05.000")
	services := []string{
		"api-gateway", "auth-service", "user-db", "order-processor",
		"payment-gateway", "notification-service", "logging-agent",
		"cache-redis", "message-queue", "frontend-server",
	}
	events := []string{
		"Started", "Stopped", "Initialized", "Connected", "Disconnected",
		"GC Triggered", "Health check passed", "Health check failed",
		"Configuration reloaded", "Connection pool exhausted",
		"Request timeout", "Retry attempt", "Task completed",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	svc := services[r.Intn(len(services))]
	evt := events[r.Intn(len(events))]

	var allowedLevels []string
	switch evt {
	case "Started", "Initialized", "Connected", "Health check passed", "Task completed":
		allowedLevels = []string{"INFO"}
	case "Stopped", "Disconnected", "GC Triggered", "Configuration reloaded":
		allowedLevels = []string{"INFO", "WARN"}
	case "Connection pool exhausted", "Request timeout", "Retry attempt":
		allowedLevels = []string{"WARN", "ERROR"}
	case "Health check failed":
		allowedLevels = []string{"ERROR"}
	default:
		allowedLevels = []string{"INFO", "WARN", "ERROR"}
	}

	level := allowedLevels[r.Intn(len(allowedLevels))]

	var detail string
	switch evt {
	case "Started", "Initialized", "Connected":
		detail = fmt.Sprintf(" (pid=%d)", r.Intn(1000)+1000)
	case "Stopped", "Disconnected":
		detail = fmt.Sprintf(" (elapsed=%dms)", r.Intn(500)+10)
	case "GC Triggered":
		detail = fmt.Sprintf(" (freed=%dMB)", r.Intn(100)+1)
	case "Health check failed":
		detail = fmt.Sprintf(" (status=%d)", r.Intn(500)+400)
	case "Request timeout":
		detail = fmt.Sprintf(" (timeout=%ds)", r.Intn(5)+1)
	case "Task completed":
		detail = fmt.Sprintf(" (duration=%dms)", r.Intn(200)+10)
	default:
		detail = ""
	}

	var color string
	switch level {
	case "INFO":
		color = "fg-green"
	case "WARN":
		color = "fg-yellow"
	case "ERROR":
		color = "fg-red"
	default:
		color = "fg-white"
	}

	raw := fmt.Sprintf("[%s] %-5s [%s] %s%s", now, level, svc, evt, detail)

	return fmt.Sprintf("[%s]%s[-]", color, raw)
}

func main() {
	wnd := tui.NewWindow()

	tv := extra.NewTextView(2).WithLines([]string{
		"Text [fg-red]Red [bold]Red Bold[-] [fg-blue]Blue[-]",
		"[bg-yellow]Yellow Background[reset]",
		"[italic]Italic[-]",
		"[cursive]Cursive[-]",
		"[reverse]Reversed[-]",
		"[underline]Underline[-]",
	}).WithFixedWidth(25)

	go func() {
		for {
			for range 5 {
				select {
				case <-wnd.OnQuit():
					return
				default:
					wnd.Commit(tv.ScrollDown)
					time.Sleep(time.Second / 3)
				}

			}
			for range 5 {
				select {
				case <-wnd.OnQuit():
					return
				default:
					wnd.Commit(tv.ScrollUp)
					time.Sleep(time.Second / 3)
				}

			}
		}
	}()

	logs := extra.NewTextView(10).WithLines([]string{"[bold]You can look at logs here...[-]"}).WithFixedWidth(75)

	go func() {
		ticker := time.NewTicker(time.Second / 2)
		defer ticker.Stop()
		for {
			select {
			case t := <-ticker.C:
				wnd.Commit(func() {
					_ = t
					logs.Append(generateLog())
					logs.ScrollDown()
				})
			case <-wnd.OnQuit():
				return
			}
		}
	}()

	wnd.SetContent(tui.NewHBox(tui.NewFrame(tv).WithTitle(tui.Title{Text: "TextView", Pos: tui.TitleTopCenter, Style: tui.FrYellow}).Rounded(), tui.NewFrame(logs).WithTitle(tui.Title{Text: "TextView - Logs", Pos: tui.TitleTopCenter, Style: tui.FrBrightMagenta}).Rounded()))

	wnd.Run()
}

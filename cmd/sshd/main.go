package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"terminal-echoware/internal/api"
	"terminal-echoware/internal/tui"
	"terminal-echoware/pkg/config"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg := config.GetConfig()

	apiBaseURL := os.Getenv("API_BASE_URL")
	if apiBaseURL == "" {
		apiBaseURL = cfg.APIBaseURL
	} else {
		cfg.APIBaseURL = apiBaseURL
	}

	sshPort := os.Getenv("SSH_PORT")
	if sshPort != "" {
		cfg.SSHPort = sshPort
	}

	showControls := os.Getenv("SHOW_CONTROLS")
	if showControls == "false" {
		cfg.ShowControls = false
	}

	apiClient := api.NewClient(apiBaseURL)

	s, err := wish.NewServer(
		wish.WithAddress(":"+cfg.SSHPort),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
				return tui.NewModel(apiClient), []tea.ProgramOption{
					tea.WithAltScreen(),
					tea.WithMouseCellMotion(),
				}
			}),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("SSH server starting on :%s", cfg.SSHPort)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-done
	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

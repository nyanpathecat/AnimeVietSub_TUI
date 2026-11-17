package main

import (
	"fmt"
	"os"

	"avs-tui/internal/api"
	"avs-tui/internal/config"
	"avs-tui/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	client, err := api.NewClient(api.DefaultBaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi khởi tạo API Client: %v\n", err)
		os.Exit(1)
	}

	// Cấu hình FlareSolverr: biến môi trường AVS_SOLVER_URL được ưu tiên hơn file cấu hình
	// để có thể chạy tạm trên máy khác mà không phải sửa .avs_config.json.
	if _, fromEnv := os.LookupEnv("AVS_SOLVER_URL"); !fromEnv {
		client.SetSolverURL(config.LoadConfig().SolverURL)
	}

	app := tui.NewAppModel(client)
	p := tea.NewProgram(app, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi chạy ứng dụng: %v\n", err)
		os.Exit(1)
	}
}

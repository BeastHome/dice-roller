package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/showr/dice-roller/internal/dice"
	"github.com/showr/dice-roller/internal/history"
	"github.com/showr/dice-roller/internal/presentation"
)

type app struct {
	screen tcell.Screen
	engine *dice.Engine

	// Presentation
	formatter presentation.Formatter

	// Input state
	input      string
	cursorPos  int
	helpOffset int

	// Input history
	inputHistory      []string
	inputHistoryIndex int // -1 = not browsing history

	// Output pane
	outputLines        []string
	outputOffset       int
	outputHitZones     []outputHitZone
	selectedOutputRoll int

	// History pane
	historyLines   []string
	historyResults []interface{}
	historyOffset  int
	historyStore   history.Store // abstracted history persistence

	// UI state
	activePane int // 0=input, 1=output, 2=history
	quitting   bool
}

func RunTUI(engine *dice.Engine) error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := s.Init(); err != nil {
		return err
	}
	defer s.Fini()

	s.EnableMouse()
	s.Clear()

	a := &app{
		screen:             s,
		engine:             engine,
		formatter:          presentation.NewDefaultFormatter(),
		outputLines:        []string{"Output will appear here."},
		selectedOutputRoll: -1,
		historyLines:       []string{},
		activePane:         1,
		inputHistoryIndex:  -1,
		historyStore:       history.NewFileStore(),
	}

	// Load all existing history files at startup
	a.loadAllHistory()
	if len(a.historyLines) > 0 {
		a.historyOffset = len(a.historyLines) - 1
		a.syncOutputFromHistory()
	}
	a.redraw()

	// Main event loop
	for {
		ev := s.PollEvent()
		switch e := ev.(type) {
		case *tcell.EventKey:
			if a.handleKeyEvent(e) {
				return nil // clean exit
			}
		case *tcell.EventMouse:
			a.handleMouseEvent(e)
		case *tcell.EventResize:
			s.Sync()
			a.redraw()
		}
	}
}

// ------------------------------------------------------------
// Load all history files at startup
// ------------------------------------------------------------
func (a *app) loadAllHistory() {
	dir := dice.HistoryDir()
	files, err := os.ReadDir(dir)
	if err != nil {
		// No history directory or unreadable; just start empty
		return
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name())
		if ext != ".json" && ext != ".log" {
			continue
		}

		path := filepath.Join(dir, f.Name())
		results, err := a.historyStore.Load(path)
		if err != nil {
			continue
		}

		for _, r := range results {
			switch v := r.(type) {
			case dice.Result:
				expr := strings.ToValidUTF8(v.Expression, "")
				summary := fmt.Sprintf("%s -> %d", expr, v.Total)

				a.historyLines = append(a.historyLines, summary)
				a.historyResults = append(a.historyResults, v)

			case dice.MultiRollResult:
				if v.Summary != "" {
					// ⭐ Use the persisted summary
					a.historyLines = append(a.historyLines, v.Summary)
				} else {
					// Fallback for old entries
					expr := strings.ToValidUTF8(v.Expression, "")
					a.historyLines = append(a.historyLines, expr)
				}
				a.historyResults = append(a.historyResults, v)

			default:
				summary := "(unknown history entry)"
				a.historyLines = append(a.historyLines, summary)
				a.historyResults = append(a.historyResults, v)
			}
		}
	}
}

// ------------------------------------------------------------
// redraw() — kept in app.go by your choice
// ------------------------------------------------------------
func (a *app) redraw() {
	if a.screen == nil {
		return // test mode: skip rendering
	}
	w, h := a.screen.Size()
	if w <= 10 || h <= 6 {
		a.drawCenteredText("Terminal too small", w, h)
		a.screen.Show()
		return
	}

	inputPane, outputPane, historyPane := computePaneLayout(w, h)

	a.drawBox(inputPane.x1, inputPane.y1, inputPane.x2, inputPane.y2, " Input ")
	a.drawBox(outputPane.x1, outputPane.y1, outputPane.x2, outputPane.y2, " Output ")
	a.drawBox(historyPane.x1, historyPane.y1, historyPane.x2, historyPane.y2, " History ")

	a.drawInput(inputPane.x1+1, inputPane.y1+1, inputPane.x2-1, inputPane.y2-1)
	a.drawOutputWithScroll(outputPane.x1+1, outputPane.y1+1, outputPane.x2-1, outputPane.y2-1)
	a.drawHistoryWithHighlight(historyPane.x1+1, historyPane.y1+1, historyPane.x2-1, historyPane.y2-1)

	if a.quitting {
		a.drawQuitDialog()
	}

	a.screen.Show()
}

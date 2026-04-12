package tui

import (
	"strings"

	"github.com/showr/dice-roller/internal/dice"
	"github.com/showr/dice-roller/internal/parse"
)

// ------------------------------------------------------------
// Enter Key — Roll Expression
// ------------------------------------------------------------
func (a *app) handleEnter() {
	raw := strings.TrimSpace(a.input)
	if raw == "" {
		return
	}

	// Store original input in history (what the user actually typed)
	a.inputHistory = append(a.inputHistory, raw)
	a.inputHistoryIndex = -1

	if a.handleMetaCommand(raw) {
		return
	}

	a.ensureSessionFile(raw)

	parsed, err := parse.ParseLine(raw)
	if err != nil {
		a.setOutputError(err)
		return
	}

	if len(parsed.Expressions) == 0 {
		a.clearInput()
		return
	}

	expr := parsed.Expressions[0]
	count := parsed.Multi
	if count <= 0 {
		count = 1
	}

	result, err := a.engine.Evaluate(expr, count)
	if err != nil {
		a.setOutputError(err)
		return
	}

	switch v := result.(type) {
	case dice.Result:
		a.handleSingleRollResult(v)
	case dice.MultiRollResult:
		v.Expression = expr
		a.handleConsolidatedMultiRoll(&v)
	}

	a.clearInput()
}

func (a *app) setOutputError(err error) {
	a.outputLines = []string{"Error: " + err.Error()}
	a.clearInput()
	a.outputOffset = 0
}

func (a *app) handleMetaCommand(raw string) bool {
	switch raw {
	case "--help", "-h":
		a.outputLines = append([]string{"Dice Roller Help", ""}, dice.HelpLines()...)
		a.clearInput()
		a.outputOffset = 0
		return true
	case "--version":
		a.outputLines = []string{"dice-roller version " + dice.Version}
		a.clearInput()
		a.outputOffset = 0
		return true
	default:
		return false
	}
}

// ------------------------------------------------------------
// Roll Result Handlers
// ------------------------------------------------------------
func (a *app) handleSingleRollResult(res dice.Result) {
	dice.AttachVerbose(&res)
	a.outputLines = strings.Split(res.Verbose, "\n")
	a.selectedOutputRoll = -1

	summary := a.formatter.FormatSingleSummary(res)
	a.ensureSessionFile(res.Expression)

	a.historyLines = append(a.historyLines, summary)
	a.historyResults = append(a.historyResults, res)
	a.historyOffset = len(a.historyResults) - 1

	_ = a.historyStore.Append(res)
	a.outputOffset = 0
}

func (a *app) handleConsolidatedMultiRoll(mr *dice.MultiRollResult) {
	if len(mr.Rolls) == 0 {
		return
	}
	a.selectedOutputRoll = defaultSelectedRollIndex(*mr)

	last := mr.Rolls[len(mr.Rolls)-1]
	dice.AttachVerbose(&last)
	a.outputLines = strings.Split(last.Verbose, "\n")

	mr.Expression = dice.FormatMultiExpression(mr.Expression, len(mr.Rolls))
	summary := a.formatter.FormatMultiSummary(*mr)
	mr.Summary = summary

	a.ensureSessionFile(mr.Expression)

	a.historyLines = append(a.historyLines, summary)
	a.historyResults = append(a.historyResults, *mr)
	a.historyOffset = len(a.historyResults) - 1

	_ = a.historyStore.Append(*mr)
	a.outputOffset = 0
}

// ------------------------------------------------------------
// Session File Creation
// ------------------------------------------------------------
func (a *app) ensureSessionFile(expr string) {
	if a.historyStore.CurrentSession() != "" {
		return
	}

	_, f, err := a.historyStore.NewSession(expr)
	if err != nil {
		return
	}
	f.Close()
}

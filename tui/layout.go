package tui

type paneRect struct {
	x1 int
	y1 int
	x2 int
	y2 int
}

func computePaneLayout(width, height int) (paneRect, paneRect, paneRect) {
	topHeight := height * 2 / 3
	if topHeight < 5 {
		topHeight = 5
	}

	inputWidth := width / 2
	input := paneRect{x1: 0, y1: 0, x2: inputWidth - 1, y2: topHeight - 1}
	output := paneRect{x1: inputWidth, y1: 0, x2: width - 1, y2: topHeight - 1}
	history := paneRect{x1: 0, y1: topHeight, x2: width - 1, y2: height - 1}
	return input, output, history
}

package api

import (
	"bytes"
	"fmt"
)

func (g Game) String() string {
	return fmt.Sprintf("Game %s, created on %s by %s (%s)",
		g.Identifier, g.CreationDate, g.Creator, g.Teaser)
}

func (g GameStatus) String() string {
	return fmt.Sprintf("Game %s, created on %s by %s (%s), turn %d (%s)",
		g.Identifier, g.CreationDate, g.Creator, g.Teaser, g.Turn, g.Status)
}

func (cmds Commands) String() string {
	return string(cmds)
}

func (resp baseResponse) String() string {
	return fmt.Sprintf(`{"status": "%s", "response": %s}`,
		resp.Status, resp.Response)
}

func (t Turn) String() string {
	return fmt.Sprintf("turn %d", t.Number)
}

func (t Turn) PrettyString() string {
	var buf bytes.Buffer

	pmap := NewPartialMap()

	for _, ant := range t.AntsStatuses {
		pmap.Combine(ant.Vision)
	}

	buf.WriteString(fmt.Sprintf("Turn %d\n\nMap:\n%s",
		t.Number, PrettyMap(pmap)))

	return buf.String()
}

func PrettyMap(m MapInterface) string {
	var buf bytes.Buffer

	w, h := m.Width(), m.Height()

	// "draw" a line at the beginning to see the map width
	for x := 0; x < w; x++ {
		buf.WriteString("-")
	}

	buf.WriteString("\n")

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := m.Cell(x, y)
			if c == nil {
				buf.WriteString(" ")
				continue
			}

			switch c.Content {
			case "grass":
				buf.WriteString("_")
			case "rock":
				buf.WriteString("#")
			case "sugar":
				buf.WriteString("s")
			case "mill":
				buf.WriteString("m")
			case "meat":
				buf.WriteString("M")
			case "water":
				buf.WriteString("~")
			default:
				buf.WriteString("?")
			}

		}

		buf.WriteString("\n")
	}

	return buf.String()
}

func (p Position) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func (d Direction) String() string {
	return fmt.Sprintf("(%d, %d)", d.X, d.Y)
}

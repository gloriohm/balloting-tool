package table

import (
	"ballot-tool/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ExportExcel(table ParsedTable) error {
	f := excelize.NewFile()
	sheet := utils.StripLabel(table.Label, "Table")
	sheet = utils.SanitizeFilename(sheet)
	f.SetSheetName("Sheet1", sheet)

	for i, row := range table.Rows {
		cIdx := 1
		rIdx := i + 1
		for _, cell := range row {
			if err := parseCell(f, sheet, rIdx, cIdx, cell); err != nil {
				return err
			}
			if cell.Colspan > 0 {
				cIdx = cIdx + cell.Colspan
			} else {
				cIdx++
			}
		}
	}

	path := sheet + ".xlsx"
	if err := writeXLSX(f, path); err != nil {
		return err
	}

	return nil
}

func parseCell(f *excelize.File, sheet string, row, col int, c Cell) error {
	if c.Colspan <= 0 {
		c.Colspan = 1
	}

	startCell, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return err
	}
	endCell, err := excelize.CoordinatesToCellName(col+c.Colspan-1, row)
	if err != nil {
		return err
	}

	if err := f.SetCellValue(sheet, startCell, c.Text); err != nil {
		return err
	}

	if c.Colspan > 1 {
		if err := f.MergeCell(sheet, startCell, endCell); err != nil {
			return err
		}
	}

	styleDef := &excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: mapAlign(c.Align),
			Vertical:   "center",
		},
		Border: parseCSSBorders(c.Style),
	}

	styleID, err := f.NewStyle(styleDef)
	if err != nil {
		return err
	}

	return f.SetCellStyle(sheet, startCell, endCell, styleID)
}

func writeXLSX(f *excelize.File, fName string) error {
	rootPath, err := genOutPath()
	if err != nil {
		return err
	}
	path := filepath.Join(rootPath, fName)
	if err := f.SaveAs(path); err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	return nil
}

func genOutPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(home, "downloads", "excel_out")

	if err = ensureDir(path); err != nil {
		return path, err
	}

	return path, nil
}

func ensureDir(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}

	return err
}

// everything below is vibed

func mapAlign(s string) string {
	switch s {
	case "left", "start":
		return "left"
	case "center", "centre", "middle":
		return "center"
	case "right", "end":
		return "right"
	default:
		return "left"
	}
}

func parseCSSBorders(css string) []excelize.Border {
	props := parseCSSDecls(css)

	sides := []string{"left", "right", "top", "bottom"}
	var borders []excelize.Border

	for _, side := range sides {
		styleKey := "border-" + side + "-style"
		widthKey := "border-" + side + "-width"
		colorKey := "border-" + side + "-color"

		cssStyle := strings.ToLower(strings.TrimSpace(props[styleKey]))
		cssWidth := strings.ToLower(strings.TrimSpace(props[widthKey]))
		cssColor := strings.TrimPrefix(strings.TrimSpace(props[colorKey]), "#")

		// No border unless style implies one.
		if cssStyle == "" || cssStyle == "none" {
			continue
		}

		borderStyle := mapBorderStyle(cssStyle, cssWidth)
		if borderStyle == 0 {
			continue
		}

		if cssColor == "" {
			cssColor = "000000"
		}

		borders = append(borders, excelize.Border{
			Type:  side,
			Color: cssColor,
			Style: borderStyle,
		})
	}

	return borders
}

func parseCSSDecls(css string) map[string]string {
	out := make(map[string]string)
	for _, part := range strings.Split(css, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			continue
		}
		k := strings.ToLower(strings.TrimSpace(kv[0]))
		v := strings.TrimSpace(kv[1])
		out[k] = v
	}
	return out
}

func mapBorderStyle(cssStyle, cssWidth string) int {
	// Excelize border styles are Excel indexes, not CSS widths.
	// This mapping is approximate.

	switch cssStyle {
	case "solid":
		// Approximate width handling.
		// 1px -> thin continuous
		// 2px -> medium continuous
		// >=3px -> thick continuous
		px := parsePx(cssWidth)
		switch {
		case px >= 3:
			return 5 // continuous thick
		case px >= 2:
			return 2 // continuous medium
		default:
			return 1 // continuous thin
		}
	case "dashed":
		return 3
	case "dotted":
		return 4
	case "double":
		return 6
	default:
		// fallback: treat unknown visible styles as thin solid
		return 1
	}
}

func parsePx(s string) int {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.TrimSuffix(s, "px")
	n, err := strconv.Atoi(s)
	if err != nil {
		return 1
	}
	return n
}

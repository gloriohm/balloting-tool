package table

import "ballot-tool/internal/utils"

type TableWrap struct {
	Title string `xml:"caption>title"`
	Label string `xml:"label"`
	Table Table  `xml:"table"`
}

type Table struct {
	TableHead TableHead `xml:"thead"`
	TableBody TableBody `xml:"tbody"`
	Cols      []Col     `xml:"col"`
}

type TableHead struct {
	Rows []HeaderRow `xml:"tr"`
}

type TableBody struct {
	Rows []BodyRow `xml:"tr"`
}

type Col struct {
	Width string `xml:"width,attr"`
}

type HeaderRow struct {
	TH []Cell `xml:"th"`
}

type BodyRow struct {
	TD []Cell `xml:"td"`
}

type Cell struct {
	Text    string `xml:",chardata"`
	Style   string `xml:"style,attr"`
	Align   string `xml:"align,attr"`
	Colspan int    `xml:"colspan,attr"`
}

type ParsedTable struct {
	Title      string
	Label      string
	Cols       []string
	HeaderRows int
	Rows       [][]Cell
}

func (tw TableWrap) toParsedTable() ParsedTable {
	out := ParsedTable{
		Title:      utils.NormalizeSpace(tw.Title),
		HeaderRows: len(tw.Table.TableHead.Rows),
		Label:      tw.Label,
	}

	if len(tw.Table.Cols) > 0 {
		for _, col := range tw.Table.Cols {
			out.Cols = append(out.Cols, col.Width)
		}
	}

	if len(tw.Table.TableHead.Rows) > 0 {
		for _, row := range tw.Table.TableHead.Rows {
			var vals []Cell

			for _, cell := range row.TH {
				vals = append(vals, cell)
			}

			out.Rows = append(out.Rows, vals)
		}
	}

	if len(tw.Table.TableBody.Rows) > 0 {
		for _, row := range tw.Table.TableBody.Rows {
			var vals []Cell

			for _, cell := range row.TD {
				vals = append(vals, cell)
			}

			out.Rows = append(out.Rows, vals)
		}
	}

	return out
}

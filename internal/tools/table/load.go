package table

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

func LoadFile(path string, secID string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := xml.NewDecoder(f)

	_, err = findElementByID(dec, "app", secID)
	if err != nil {
		return err
	}

	tables, err := parseTableWraps(dec)
	if err != nil {
		return err
	}

	log.Printf("%d tables successfully parsed\n", len(tables))

	for i, table := range tables {
		log.Printf("processing table %d\n", i+1)
		if err = ExportExcel(table); err != nil {
			return err
		}
	}

	return nil
}

func findElementByID(dec *xml.Decoder, tagName, wantID string) (xml.StartElement, error) {
	for {
		tok, err := dec.Token()
		if err != nil {
			return xml.StartElement{}, err
		}

		se, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}

		if se.Name.Local == tagName && attr(se, "id") == wantID {
			return se, nil
		}
	}
}

func attr(se xml.StartElement, name string) string {
	for _, a := range se.Attr {
		if a.Name.Local == name {
			return a.Value
		}
	}
	return ""
}

func parseTableWraps(dec *xml.Decoder) ([]ParsedTable, error) {
	var out []ParsedTable

	depth := 1

	for depth > 0 {
		tok, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return out, fmt.Errorf("unexpected EOF before section closed")
			}
			return nil, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == "table-wrap" {
				var tw TableWrap
				if err := dec.DecodeElement(&tw, &t); err != nil {
					return nil, err
				}
				out = append(out, tw.toParsedTable())
				continue
			}

			depth++

		case xml.EndElement:
			depth--
		}
	}

	return out, nil
}

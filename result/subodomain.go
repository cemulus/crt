package result

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

type Subdomain struct {
	Name string `json:"subdomain"`
}

type SubdomainResult []Subdomain

func (s SubdomainResult) Table() []byte {
	res := new(bytes.Buffer)
	table := tablewriter.NewWriter(res)

	table.SetHeader([]string{"Subdomains"})

	table.SetHeaderColor(tablewriter.Color(tablewriter.FgHiBlueColor))
	table.SetColumnColor(tablewriter.Color(tablewriter.FgHiYellowColor))

	for _, sub := range s {
		table.Append([]string{sub.Name})
	}

	table.SetRowLine(true)
	table.SetRowSeparator("—")
	table.Render()

	return res.Bytes()
}

func (s SubdomainResult) JSON() ([]byte, error) {
	res, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal results: %s", err)
	}

	return res, nil
}

func (s SubdomainResult) CSV() ([]byte, error) {
	res := new(bytes.Buffer)
	w := csv.NewWriter(res)

	if err := w.Write([]string{"subdomain"}); err != nil {
		return nil, fmt.Errorf("failed to write CSV headers: %s", err)
	}

	for _, sub := range s {
		if err := w.Write([]string{sub.Name}); err != nil {
			return nil, fmt.Errorf("failed to write CSV content: %s", err)
		}
	}
	w.Flush()

	return res.Bytes(), nil
}

func (s SubdomainResult) Size() int { return len(s) }

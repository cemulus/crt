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

func (s SubdomainResult) Table() string {
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

	return res.String()
}

func (s SubdomainResult) JSON() (string, error) {
	res, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %s", err)
	}

	return string(res), nil
}

func (s SubdomainResult) CSV() (string, error) {
	res := new(bytes.Buffer)
	w := csv.NewWriter(res)

	if err := w.Write([]string{"subdomain"}); err != nil {
		return "", fmt.Errorf("failed to write CSV headers: %s", err)
	}

	for _, sub := range s {
		if err := w.Write([]string{sub.Name}); err != nil {
			return "", fmt.Errorf("failed to write CSV content: %s", err)
		}
	}
	w.Flush()

	return res.String(), nil
}

func (s SubdomainResult) Size() int { return len(s) }

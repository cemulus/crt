package result

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

type Certificate struct {
	IssuerCaID     int       `json:"issuer_ca_id"`
	IssuerName     string    `json:"issuer_name"`
	CommonName     string    `json:"common_name"`
	NameValue      string    `json:"name_value"`
	ID             int       `json:"id"`
	EntryTimestamp time.Time `json:"entry_timestamp"`
	NotBefore      time.Time `json:"not_before"`
	NotAfter       time.Time `json:"not_after"`
	SerialNumber   string    `json:"serial_number"`
}

type CertResult []Certificate

func (r CertResult) Table() string {
	res := new(bytes.Buffer)
	table := tablewriter.NewWriter(res)

	info := []string{"Matching", "Logged At", "Not Before", "Not After", "Issuer"}
	table.SetHeader(info)
	table.SetFooter(info)

	blue := tablewriter.Color(tablewriter.FgHiBlueColor)
	yellow := tablewriter.Color(tablewriter.FgHiYellowColor)
	white := tablewriter.Color(tablewriter.FgWhiteColor)

	table.SetHeaderColor(blue, blue, blue, blue, blue)
	table.SetFooterColor(blue, blue, blue, blue, blue)
	table.SetColumnColor(yellow, white, white, white, white)

	for _, cert := range r {
		table.Append([]string{
			cert.NameValue,
			cert.EntryTimestamp.String()[0:10],
			cert.NotBefore.String()[0:10],
			cert.NotAfter.String()[0:10],
			strings.Trim(strings.Split(strings.Split(cert.IssuerName, "O=")[1], ",")[0], "\""),
		})
	}

	table.SetRowLine(true)
	table.SetRowSeparator("—")
	table.Render()

	return res.String()
}

func (r CertResult) JSON() (string, error) {
	res, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %s", err)
	}

	return string(res), nil
}

func (r CertResult) CSV() (string, error) {
	res := new(bytes.Buffer)
	w := csv.NewWriter(res)

	err := w.Write([]string{
		"issuer_ca_id", "issuer_name", "common_name", "name_value", "id",
		"entry_timestamp", "not_before", "not_after", "serial_number",
	})
	if err != nil {
		return "", fmt.Errorf("failed to write CSV headers: %s", err)
	}

	for _, v := range r {
		err = w.Write([]string{
			strconv.Itoa(v.IssuerCaID),
			v.IssuerName,
			v.CommonName,
			v.NameValue,
			strconv.Itoa(v.ID),
			v.EntryTimestamp.String(),
			v.NotBefore.String(),
			v.NotAfter.String(),
			v.SerialNumber})
		if err != nil {
			return "", fmt.Errorf("failed to write CSV content: %s", err)
		}
	}
	w.Flush()

	return res.String(), nil
}

func (r CertResult) Size() int { return len(r) }

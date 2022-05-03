package repository

import (
	"database/sql"
	"fmt"

	"github.com/cemulus/crt/result"

	_ "github.com/lib/pq"
)

var (
	driver = "postgres"
	host   = "crt.sh"
	port   = 5432
	user   = "guest"
	dbname = "certwatch"
	login  = fmt.Sprintf("host=%s port=%d user=%s dbname=%s", host, port, user, dbname)
)

type Repository struct {
	db *sql.DB
}

func New() (*Repository, error) {
	db, err := sql.Open(driver, login)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %s", err)
	}

	return &Repository{db}, nil
}

func (r *Repository) GetCertLogs(domain string, expired bool, limit int) (result.Certificates, error) {
	filter := ""

	if expired {
		filter = excludeExpiredFilter
	}

	stmt := fmt.Sprintf(certLogScript, domain, domain, filter, limit)

	rows, err := r.db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to query db: %s", err)
	}
	defer rows.Close()

	var res result.Certificates

	var (
		issuerCaID                                      sql.NullInt32
		id                                              sql.NullInt64
		issuerName, commonName, nameValue, serialNumber sql.NullString
		entryTimestamp, notBefore, notAfter             sql.NullTime
	)

	for rows.Next() {
		err = rows.Scan(
			&issuerCaID,
			&issuerName,
			&commonName,
			&nameValue,
			&id,
			&entryTimestamp,
			&notBefore,
			&notAfter,
			&serialNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %s", err)
		}

		certificate := result.Certificate{
			IssuerCaID:     int((issuerCaID).Int32),
			IssuerName:     issuerName.String,
			CommonName:     commonName.String,
			NameValue:      nameValue.String,
			ID:             int((id).Int64),
			EntryTimestamp: entryTimestamp.Time,
			NotBefore:      notBefore.Time,
			NotAfter:       notAfter.Time,
			SerialNumber:   serialNumber.String}

		res = append(res, certificate)
	}

	return res, nil
}

func (r *Repository) GetSubdomains(domain string, expired bool, limit int) (result.Subdomains, error) {
	filter := ""

	if expired {
		filter = excludeExpiredFilter
	}

	stmt := fmt.Sprintf(subdomainScript, domain, domain, filter, limit)

	rows, err := r.db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to query row: %s", err)
	}
	defer rows.Close()

	var res result.Subdomains
	var subdmn sql.NullString

	for rows.Next() {
		if err = rows.Scan(&subdmn); err != nil {
			return nil, fmt.Errorf("failed to scan row: %s", err)
		}

		res = append(res, result.Subdomain{Name: subdmn.String})
	}

	return res, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/cemulus/crt/result"
)

var (
	driver = "postgres"
	host   = "crt.sh"
	port   = 5432
	user   = "guest"
	dbname = "certwatch"
	login  = fmt.Sprintf("host=%s port=%d user=%s dbname=%s", host, port, user, dbname)
)

const (
	statement = `WITH ci AS (
	SELECT min(sub.CERTIFICATE_ID) ID,
		min(sub.ISSUER_CA_ID) ISSUER_CA_ID,
		array_agg(DISTINCT sub.NAME_VALUE) NAME_VALUES,
		x509_commonName(sub.CERTIFICATE) COMMON_NAME,
		x509_notBefore(sub.CERTIFICATE) NOT_BEFORE,
		x509_notAfter(sub.CERTIFICATE) NOT_AFTER,
		encode(x509_serialNumber(sub.CERTIFICATE), 'hex') SERIAL_NUMBER
	FROM (SELECT *
			FROM certificate_and_identities cai
			WHERE plainto_tsquery('certwatch', '%s') @@ identities(cai.CERTIFICATE)
				AND cai.NAME_VALUE ILIKE ('%%' || '%s' || '%%')
			LIMIT 10000
		) sub
	GROUP BY sub.CERTIFICATE
)
SELECT ci.ISSUER_CA_ID,
	ca.NAME ISSUER_NAME,
	ci.COMMON_NAME,
	array_to_string(ci.NAME_VALUES, chr(10)) NAME_VALUE,
	ci.ID ID,
	le.ENTRY_TIMESTAMP,
	ci.NOT_BEFORE,
	ci.NOT_AFTER,
	ci.SERIAL_NUMBER
FROM ci
	LEFT JOIN LATERAL (
		SELECT min(ctle.ENTRY_TIMESTAMP) ENTRY_TIMESTAMP
		FROM ct_log_entry ctle
		WHERE ctle.CERTIFICATE_ID = ci.ID
	) le ON TRUE,
	ca
WHERE ci.ISSUER_CA_ID = ca.ID
ORDER BY le.ENTRY_TIMESTAMP DESC NULLS LAST
LIMIT %d`
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

func (r *Repository) GetCertLogs(domain string, limit int) (result.CertResult, error) {
	stmt := fmt.Sprintf(statement, domain, domain, limit)

	rows, err := r.db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to query db: %s", err)
	}
	defer rows.Close()

	var res result.CertResult

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

func (r *Repository) Close() error {
	return r.db.Close()
}

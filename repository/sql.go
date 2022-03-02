package repository

const (
	certLogScript = `WITH ci AS (
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
				%s --filter
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

	subdomainScript = `SELECT DISTINCT cai.NAME_VALUE
FROM certificate_and_identities cai
WHERE plainto_tsquery('certwatch', '%s') @@ identities(cai.CERTIFICATE)
	AND cai.NAME_VALUE ILIKE ('%%' || '%s' || '%%')
	%s --filter
LIMIT %d`

	excludeExpiredFilter = `AND coalesce(x509_notAfter(cai.CERTIFICATE), 'infinity'::timestamp) >= date_trunc('year', now() AT TIME ZONE 'UTC')
	AND x509_notAfter(cai.CERTIFICATE) >= now() AT TIME ZONE 'UTC'`
)

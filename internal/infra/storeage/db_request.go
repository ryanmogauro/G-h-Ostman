package storeage

import (
	"database/sql"
	"strings"
	"time"

	"github.com/ryanmogauro/ghostman/internal/domain"
	_ "modernc.org/sqlite"
)

// Adds request to user's request history
func InsertRequest(db *sql.DB, req domain.Request) error {
	url := req.URL
	method := req.Method
	body := req.Body
	headers := req.Headers
	createdAt := time.Now().Format(time.RFC3339)

	headerString := ""

	for k, v := range headers {
		header := k + v + "\n"
		headerString += header
	}

	query := `
	INSERT INTO requests (url, method, body, headers, created_at)
	VALUES (?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query, url, method, body, headerString, createdAt)
	if err != nil {
		return err
	}

	return nil
}

// Retreives and returns user request history
func GetHistory(db *sql.DB) ([]domain.Request, error) {
	query := `
	SELECT id, url, method, body, headers, created_at
	FROM requests
	ORDER BY created_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []domain.Request
	for rows.Next() {
		var req domain.Request
		var headerStr string
		var createdAtStr string
		err := rows.Scan(&req.ID, &req.URL, &req.Method, &req.Body, &headerStr, &createdAtStr)
		if err != nil {
			return nil, err
		}

		//Parse headers, reconstruct header map
		headers := make(map[string]string)
		headerLines := strings.Split(headerStr, "\n")
		for _, line := range headerLines {
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
		req.Headers = headers

		// Parse created_at
		req.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, err
		}

		requests = append(requests, req)
	}

	return requests, nil
}

func GetRequest(db *sql.DB, id int) (domain.Request, error) {
	query := `
		SELECT id, url, method, body, headers, created_at
		FROM requests
		WHERE id = ?
	`

	row := db.QueryRow(query, id)

	var req domain.Request
	var headerStr string
	var createdAtStr string

	err := row.Scan(&req.ID, &req.URL, &req.Method, &req.Body, &headerStr, &createdAtStr)
	if err != nil {
		return domain.Request{}, err
	}

	// Parse headers, reconstruct header map
	headers := make(map[string]string)
	headerLines := strings.Split(headerStr, "\n")
	for _, line := range headerLines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	req.Headers = headers
	req.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return domain.Request{}, err
	}
	return req, nil
}

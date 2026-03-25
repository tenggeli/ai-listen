package model

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"
)

func splitSQLStatements(sqlText string) []string {
	parts := strings.Split(sqlText, ";")
	stmts := make([]string, 0, len(parts))
	for _, part := range parts {
		stmt := strings.TrimSpace(part)
		if stmt != "" {
			stmts = append(stmts, stmt)
		}
	}
	return stmts
}

func parseJSONStringArray(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var items []string
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil
	}
	return items
}

func toJSON(v any) string {
	if v == nil {
		return "[]"
	}
	data, err := json.Marshal(v)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func formatNullDate(v sql.NullTime) string {
	if !v.Valid {
		return ""
	}
	return v.Time.Format("2006-01-02")
}

func formatNullDateTime(v sql.NullTime) string {
	if !v.Valid {
		return ""
	}
	return v.Time.Format("2006-01-02 15:04:05")
}

func ptrTime(v sql.NullTime) *time.Time {
	if !v.Valid {
		return nil
	}
	t := v.Time
	return &t
}

func operatorRoleName(roleCode int) string {
	switch roleCode {
	case operatorRoleProvider:
		return "provider"
	case operatorRoleAdmin:
		return "admin"
	default:
		return "user"
	}
}

func containsInt(list []int, target int) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

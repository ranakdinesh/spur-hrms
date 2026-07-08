package postgres

import (
	"os"
	"strings"
	"testing"
)

func TestPolicyResolutionQueryPrecedence(t *testing.T) {
	content, err := os.ReadFile("../../sql/queries/policy_engine.sql")
	if err != nil {
		t.Fatalf("read policy engine query: %v", err)
	}
	query := string(content)
	expectedOrder := []string{
		"WHEN 'employee' THEN 60",
		"WHEN 'designation' THEN 50",
		"WHEN 'workforce_type' THEN 50",
		"WHEN 'role_group' THEN 45",
		"WHEN 'department' THEN 40",
		"WHEN 'branch' THEN 30",
		"WHEN 'tenant' THEN 20",
		"10 AS precedence",
	}

	lastIndex := -1
	for _, marker := range expectedOrder {
		index := strings.Index(query, marker)
		if index == -1 {
			t.Fatalf("policy resolution query is missing precedence marker %q", marker)
		}
		if index < lastIndex {
			t.Fatalf("policy resolution precedence marker %q appears before the previous marker", marker)
		}
		lastIndex = index
	}
	if !strings.Contains(query, "ORDER BY precedence DESC, assignment_created_at DESC") {
		t.Fatalf("policy resolution query must order highest precedence and newest assignment first")
	}
	if !strings.Contains(query, "ORDER BY precedence DESC, name ASC") {
		t.Fatalf("policy resolution candidates query must order highest precedence first")
	}
}

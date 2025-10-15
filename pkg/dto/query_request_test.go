package dto

import (
	"reflect"
	"testing"
)

func TestNewQueryRequestFromURL(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected *QueryRequest
	}{
		{
			name:  "default values",
			query: "",
			expected: &QueryRequest{
				Page:      1,
				PageSize:  10,
				Condition: []Condition{},
				Order:     []Order{},
			},
		},
		{
			name:  "custom pagination",
			query: "page=3&page_size=20",
			expected: &QueryRequest{
				Page:      3,
				PageSize:  20,
				Condition: []Condition{},
				Order:     []Order{},
			},
		},
		{
			name:  "equal condition",
			query: "name=test",
			expected: &QueryRequest{
				Page:     1,
				PageSize: 10,
				Condition: []Condition{
					{Key: "name", Value: "test", Op: Equal},
				},
				Order: []Order{},
			},
		},
		{
			name:  "greater than condition",
			query: "age_gt=18",
			expected: &QueryRequest{
				Page:     1,
				PageSize: 10,
				Condition: []Condition{
					{Key: "age", Value: "18", Op: Greater},
				},
				Order: []Order{},
			},
		},
		{
			name:  "less than condition",
			query: "price_lt=100",
			expected: &QueryRequest{
				Page:     1,
				PageSize: 10,
				Condition: []Condition{
					{Key: "price", Value: "100", Op: Less},
				},
				Order: []Order{},
			},
		},
		{
			name:  "not equal condition",
			query: "status_ne=inactive",
			expected: &QueryRequest{
				Page:     1,
				PageSize: 10,
				Condition: []Condition{
					{Key: "status", Value: "inactive", Op: NotEqual},
				},
				Order: []Order{},
			},
		},
		{
			name:  "greater than or equal condition",
			query: "score_gte=80",
			expected: &QueryRequest{
				Page:     1,
				PageSize: 10,
				Condition: []Condition{
					{Key: "score", Value: "80", Op: GreaterOrEqual},
				},
				Order: []Order{},
			},
		},
		{
			name:  "less than or equal condition",
			query: "quantity_lte=10",
			expected: &QueryRequest{
				Page:     1,
				PageSize: 10,
				Condition: []Condition{
					{Key: "quantity", Value: "10", Op: LessOrEqual},
				},
				Order: []Order{},
			},
		},
		{
			name:  "multiple conditions",
			query: "age_gte=18&age_lt=60&name=admin",
			expected: &QueryRequest{
				Page:     1,
				PageSize: 10,
				Condition: []Condition{
					{Key: "age", Value: "18", Op: GreaterOrEqual},
					{Key: "age", Value: "60", Op: Less},
					{Key: "name", Value: "admin", Op: Equal},
				},
				Order: []Order{},
			},
		},
		{
			name:  "single ascending order",
			query: "sort=name.asc",
			expected: &QueryRequest{
				Page:      1,
				PageSize:  10,
				Condition: []Condition{},
				Order: []Order{
					{Key: "name", Direction: Ascending},
				},
			},
		},
		{
			name:  "single descending order",
			query: "sort=created_at.desc",
			expected: &QueryRequest{
				Page:      1,
				PageSize:  10,
				Condition: []Condition{},
				Order: []Order{
					{Key: "created_at", Direction: Descending},
				},
			},
		},
		{
			name:  "multiple order fields",
			query: "sort=name.asc,age.desc",
			expected: &QueryRequest{
				Page:      1,
				PageSize:  10,
				Condition: []Condition{},
				Order: []Order{
					{Key: "name", Direction: Ascending},
					{Key: "age", Direction: Descending},
				},
			},
		},
		{
			name:  "complex query with conditions and ordering",
			query: "page=2&page_size=15&age_gte=21&status=active&sort=age.desc,name.asc",
			expected: &QueryRequest{
				Page:     2,
				PageSize: 15,
				Condition: []Condition{
					{Key: "age", Value: "21", Op: GreaterOrEqual},
					{Key: "status", Value: "active", Op: Equal},
				},
				Order: []Order{
					{Key: "age", Direction: Descending},
					{Key: "name", Direction: Ascending},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewQueryRequestFromURL(tt.query)

			if !reflect.DeepEqual(got.Page, tt.expected.Page) {
				t.Errorf("Page = %v, want %v", got.Page, tt.expected.Page)
			}

			if !reflect.DeepEqual(got.PageSize, tt.expected.PageSize) {
				t.Errorf("PageSize = %v, want %v", got.PageSize, tt.expected.PageSize)
			}

			if !reflect.DeepEqual(got.Condition, tt.expected.Condition) {
				t.Errorf("Condition = %v, want %v", got.Condition, tt.expected.Condition)
			}

			if !reflect.DeepEqual(got.Order, tt.expected.Order) {
				t.Errorf("Order = %v, want %v", got.Order, tt.expected.Order)
			}
		})
	}
}

func TestParseCondition(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected *Condition
	}{
		{
			name: "parse equal condition",
			key:  "name",
			expected: &Condition{
				Key: "name",
				Op:  Equal,
			},
		},
		{
			name: "parse greater than condition",
			key:  "age_gt",
			expected: &Condition{
				Key: "age",
				Op:  Greater,
			},
		},
		{
			name: "parse less than condition",
			key:  "price_lt",
			expected: &Condition{
				Key: "price",
				Op:  Less,
			},
		},
		{
			name: "parse not equal condition",
			key:  "status_ne",
			expected: &Condition{
				Key: "status",
				Op:  NotEqual,
			},
		},
		{
			name: "parse greater than or equal condition",
			key:  "score_gte",
			expected: &Condition{
				Key: "score",
				Op:  GreaterOrEqual,
			},
		},
		{
			name: "parse less than or equal condition",
			key:  "quantity_lte",
			expected: &Condition{
				Key: "quantity",
				Op:  LessOrEqual,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCondition(tt.key)
			if !reflect.DeepEqual(got.Key, tt.expected.Key) {
				t.Errorf("Key = %v, want %v", got.Key, tt.expected.Key)
			}
			if !reflect.DeepEqual(got.Op, tt.expected.Op) {
				t.Errorf("Op = %v, want %v", got.Op, tt.expected.Op)
			}
		})
	}
}

func TestParseOrder(t *testing.T) {
	tests := []struct {
		name     string
		sortStr  string
		expected []Order
	}{
		{
			name:    "empty sort string",
			sortStr: "",
			expected: []Order{},
		},
		{
			name:    "single ascending order",
			sortStr: "name.asc",
			expected: []Order{
				{Key: "name", Direction: Ascending},
			},
		},
		{
			name:    "single descending order",
			sortStr: "created_at.desc",
			expected: []Order{
				{Key: "created_at", Direction: Descending},
			},
		},
		{
			name:    "multiple order fields",
			sortStr: "name.asc,age.desc,score.asc",
			expected: []Order{
				{Key: "name", Direction: Ascending},
				{Key: "age", Direction: Descending},
				{Key: "score", Direction: Ascending},
			},
		},
		{
			name:    "invalid order direction",
			sortStr: "name.invalid",
			expected: []Order{},
		},
		{
			name:    "mixed valid and invalid orders",
			sortStr: "name.asc,invalid.direction,age.desc",
			expected: []Order{
				{Key: "name", Direction: Ascending},
				{Key: "age", Direction: Descending},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseOrder(tt.sortStr)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("parseOrder() = %v, want %v", got, tt.expected)
			}
		})
	}
}
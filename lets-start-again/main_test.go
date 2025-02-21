package main

import (
	"testing"
)

func TestGetNextID(t *testing.T) {
	tests := []struct {
		name     string
		tasks    []Task
		expected int
	}{
		{"Empty Slice", []Task{}, 0},
		{"Single Task", []Task{{ID: 1}}, 1},
		{"Two tasks", []Task{{ID: 1}, {ID: 3}, {ID: 2}}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getNextID(tt.tasks)
			if result != tt.expected {
				t.Errorf("getNextID(%v) = %d; want %d", tt.tasks, result, tt.expected)
			}
		})
	}

}

package reader

import (
	"reflect"
	"testing"
)

func TestTabularData_Keep(t *testing.T) {
	tests := []struct {
		name          string
		initialData   *TabularData
		headersToKeep []string
		expectedData  *TabularData
	}{
		{
			name: "keep all headers",
			initialData: &TabularData{
				Headers: []string{"ID", "Comment", "Response"},
				Records: [][]string{
					{"R1.1", "Comment A", "Response A"},
					{"R1.2", "Comment B", "Response B"},
				},
			},
			headersToKeep: []string{"ID", "Comment", "Response"},
			expectedData: &TabularData{
				Headers: []string{"ID", "Comment", "Response"},
				Records: [][]string{
					{"R1.1", "Comment A", "Response A"},
					{"R1.2", "Comment B", "Response B"},
				},
			},
		},
		{
			name: "keep subset of headers",
			initialData: &TabularData{
				Headers: []string{"ID", "Comment", "Response", "Action"},
				Records: [][]string{
					{"R1.1", "Comment A", "Response A", "Action A"},
					{"R1.2", "Comment B", "Response B", "Action B"},
				},
			},
			headersToKeep: []string{"ID", "Action"},
			expectedData: &TabularData{
				Headers: []string{"ID", "Action"},
				Records: [][]string{
					{"R1.1", "Action A"},
					{"R1.2", "Action B"},
				},
			},
		},
		{
			name: "keep handle non-existent headers",
			initialData: &TabularData{
				Headers: []string{"ID", "Comment", "Response"},
				Records: [][]string{
					{"R1.1", "Comment A", "Response A"},
					{"R1.2", "Comment B", "Response B"},
				},
			},
			headersToKeep: []string{"Address", "Phone"},
			expectedData: &TabularData{
				Headers: []string{},
				Records: [][]string{
					{},
					{},
				},
			},
		},
		{
			name: "keep with empty records",
			initialData: &TabularData{
				Headers: []string{"ID", "Comment", "Response"},
				Records: [][]string{},
			},
			headersToKeep: []string{"Comment"},
			expectedData: &TabularData{
				Headers: []string{"Comment"},
				Records: [][]string{},
			},
		},
		{
			name: "keep with empty selection",
			initialData: &TabularData{
				Headers: []string{"ID", "Comment", "Response"},
				Records: [][]string{
					{"R1.1", "Comment A", "Response A"},
					{"R1.2", "Comment B", "Response B"},
				},
			},
			headersToKeep: []string{},
			expectedData: &TabularData{
				Headers: []string{},
				Records: [][]string{
					{},
					{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			data := &TabularData{
				Headers: make([]string, len(tt.initialData.Headers)),
				Records: make([][]string, len(tt.initialData.Records)),
			}
			copy(data.Headers, tt.initialData.Headers)
			for i, record := range tt.initialData.Records {
				data.Records[i] = make([]string, len(record))
				copy(data.Records[i], record)
			}

			// Act
			data.Keep(tt.headersToKeep)

			// Assert
			if !reflect.DeepEqual(data.Headers, tt.expectedData.Headers) {
				t.Errorf("Headers = %v, want %v", data.Headers, tt.expectedData.Headers)
			}

			if !reflect.DeepEqual(data.Records, tt.expectedData.Records) {
				t.Errorf("Records = %v, want %v", data.Records, tt.expectedData.Records)
			}
		})
	}
}

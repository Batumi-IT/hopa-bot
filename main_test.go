package main

import "testing"

func Test_containsStupidQuestion(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    bool
	}{
		{
			name:    "1",
			message: "Где купить?",
			want:    true,
		},
		{
			name:    "2",
			message: "Где найти?",
			want:    true,
		},
		{
			name:    "3",
			message: "Где продаётся?",
			want:    true,
		},
		{
			name:    "4",
			message: "Где продается?",
			want:    true,
		},
		{
			name:    "5",
			message: "Где починить?",
			want:    true,
		},
		{
			name:    "6",
			message: "Где посмотреть?",
			want:    true,
		},
		{
			name:    "1f",
			message: "что такое залупа иваныча?",
			want:    false,
		},
		{
			name:    "2f",
			message: "где деньги Лебовски?",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsStupidQuestion(tt.message); got != tt.want {
				t.Errorf("containsStupidQuestion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_containsSmartQuestion(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		message string
		want    bool
	}{
		{
			name:    "1",
			message: "Где найти рынок Хопа?",
			want:    true,
		},
		{
			name:    "2",
			message: "Как попасть на Хопу?",
			want:    true,
		},
		{
			name:    "3",
			message: "Как добраться до Хопа?",
			want:    true,
		},
		{
			name:    "4",
			message: "Где найти рынок Хопу?",
			want:    true,
		}, {
			name:    "5",
			message: "Как добраться до Хопы?",
			want:    true,
		},
		{
			name:    "6f",
			message: "Ыхыхы ахаха?",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsSmartQuestion(tt.message); got != tt.want {
				t.Errorf("containsSmartQuestion() = %v, want %v", got, tt.want)
			}
		})
	}
}

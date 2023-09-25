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
			message: "где купить?",
			want:    true,
		},
		{
			name:    "2",
			message: "где найти?",
			want:    true,
		},
		{
			name:    "3",
			message: "где продаётся?",
			want:    true,
		},
		{
			name:    "4",
			message: "где продается?",
			want:    true,
		},
		{
			name:    "5",
			message: "где починить?",
			want:    true,
		},
		{
			name:    "6",
			message: "   где   посмотреть?",
			want:    true,
		},
		{
			name:    "1f",
			message: "что такое залупа иваныча?",
			want:    false,
		},
		{
			name:    "2f",
			message: "где деньги лебовски?",
			want:    false,
		},
		{
			name:    "7",
			message: "а без покрытия чугунные сковородки кто-то видел в продаже?",
			want:    true,
		},
		{
			name:    "8",
			message: "а где посмотреть что купить?",
			want:    true,
		},
		{
			name:    "9",
			message: "в где посмотреть что купить в продаже?",
			want:    true,
		},
		{
			name:    "3f",
			message: "как пройти в библиотеку?",
			want:    false,
		},
		{
			name:    "4f",
			message: "амвппв найти?",
			want:    false,
		},
		{
			name:    "5f",
			message: "гдерь найти?",
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
	tests := []struct {
		name    string
		message string
		want    bool
	}{
		{
			name:    "1",
			message: "где найти рынок хопа?",
			want:    true,
		},
		{
			name:    "2",
			message: "как попасть на хопу?",
			want:    true,
		},
		{
			name:    "3",
			message: "как добраться до хопа?",
			want:    true,
		},
		{
			name:    "4",
			message: "где найти рынок хопу?",
			want:    true,
		}, {
			name:    "5",
			message: "как добраться до хопы?",
			want:    true,
		},
		{
			name:    "6f",
			message: "ыхыхы ахаха?",
			want:    false,
		},
		{
			name:    "7f",
			message: "когде хопы?",
			want:    false,
		},
		{
			name:    "8f",
			message: "акак хопы?",
			want:    false,
		},
		{
			name:    "9f",
			message: "акака хопы?",
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

func Test_generateReply(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "1",
			message: "Где купить?",
			want:    "На рынке Хопа!",
		},
		{
			name:    "2",
			message: "Где найти?",
			want:    "На рынке Хопа!",
		},
		{
			name:    "3",
			message: "Где рынок хопа?",
			want:    "Держи ссылку с адресом рынка Хопа, раз в гугле забанили:\nhttps://goo.gl/maps/aqN4rzapdDXvRJNW9",
		},
		{
			name:    "4",
			message: "где найти хопу?",
			want:    "Хопа на рынке Хопа! Вот, ну:\nhttps://goo.gl/maps/aqN4rzapdDXvRJNW9",
		},
		{
			name:    "f5",
			message: "Ыхыхы ахаха?",
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateReply(tt.message); got != tt.want {
				t.Errorf("generateReply() = %v, want %v", got, tt.want)
			}
		})
	}
}

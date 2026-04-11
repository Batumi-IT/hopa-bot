package main

import "testing"

type messageBoolTest struct {
	name    string
	message string
	want    bool
}

func runMessageBoolTests(t *testing.T, fn func(string) bool, tests []messageBoolTest) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fn(tt.message); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_containsStupidQuestion(t *testing.T) {
	runMessageBoolTests(t, containsStupidQuestion, []messageBoolTest{
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
		{
			name:    "10",
			message: "где взять б/у машину?",
			want:    true,
		},
		{
			name:    "11",
			message: "где ты взяла б/у Кирилла?",
			want:    true,
		},
		{
			name:    "12",
			message: "где ты нашёл рыбью голову?",
			want:    true,
		},
		{
			name:    "13",
			message: "а где ты покупала морскую слизь?",
			want:    true,
		},
		{
			name:    "14",
			message: "а где ты купила его?",
			want:    true,
		},
		{
			name:    "15",
			message: "а где вы купили его?",
			want:    true,
		},
		{
			name:    "16",
			message: "а где достали это?",
			want:    true,
		},
		{
			name:    "17",
			message: "а где достать пиво?",
			want:    true,
		},
		{
			name:    "18",
			message: "а где ты достал это?",
			want:    true,
		},
	})
}

func Test_containsSmartQuestion(t *testing.T) {
	runMessageBoolTests(t, containsSmartQuestion, []messageBoolTest{
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
	})
}

func Test_generateReply(t *testing.T) {
	tests := []struct {
		name  string
		check Check
		want  string
	}{
		{
			name:  "1",
			check: Check{Stupid: true, Smart: false},
			want:  "На рынке Хопа!",
		},
		{
			name:  "2",
			check: Check{Stupid: false, Smart: true},
			want:  "Держи ссылку с адресом рынка Хопа, раз в гугле забанили:\nhttps://goo.gl/maps/aqN4rzapdDXvRJNW9",
		},
		{
			name:  "3",
			check: Check{Stupid: true, Smart: true},
			want:  "Хопа на рынке Хопа! Вот, ну:\nhttps://goo.gl/maps/aqN4rzapdDXvRJNW9",
		},
		{
			name:  "f4",
			check: Check{Stupid: false, Smart: false},
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateReply(tt.check); got != tt.want {
				t.Errorf("generateReply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateCheck(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    Check
	}{
		{
			name:    "1",
			message: "где купить?",
			want:    Check{Stupid: true, Smart: false},
		},
		{
			name:    "2",
			message: "где найти?",
			want:    Check{Stupid: true, Smart: false},
		},
		{
			name:    "3",
			message: "где найти рынок хопа?",
			want:    Check{Stupid: true, Smart: true},
		},
		{
			name:    "f4",
			message: "ыхыхы ахаха?",
			want:    Check{Stupid: false, Smart: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateCheck(tt.message); got != tt.want {
				t.Errorf("generateCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

type botUsernameTest struct {
	name        string
	message     string
	botUsername string
}

func Test_containsMention(t *testing.T) {
	tests := []struct {
		botUsernameTest
		want bool
	}{
		{botUsernameTest{"exact mention", "@hopabot где купить?", "hopabot"}, true},
		{botUsernameTest{"mention in middle", "привет @hopabot помоги", "hopabot"}, true},
		{botUsernameTest{"mention at end", "помоги @hopabot", "hopabot"}, true},
		{botUsernameTest{"case insensitive message", "@HopaBot привет", "hopabot"}, true},
		{botUsernameTest{"case insensitive username", "@hopabot привет", "HopaBot"}, true},
		{botUsernameTest{"no mention", "привет всем", "hopabot"}, false},
		{botUsernameTest{"empty username", "@hopabot привет", ""}, false},
		{botUsernameTest{"different bot", "@otherbot привет", "hopabot"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsMention(tt.message, tt.botUsername); got != tt.want {
				t.Errorf("containsMention() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stripMention(t *testing.T) {
	tests := []struct {
		botUsernameTest
		want string
	}{
		{botUsernameTest{"strip from start", "@hopabot где купить?", "hopabot"}, "где купить?"},
		{botUsernameTest{"strip from middle", "привет @hopabot помоги", "hopabot"}, "привет  помоги"},
		{botUsernameTest{"strip case insensitive", "привет @HopaBot помоги", "hopabot"}, "привет  помоги"},
		{botUsernameTest{"no mention present", "привет всем", "hopabot"}, "привет всем"},
		{botUsernameTest{"only mention", "@hopabot", "hopabot"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripMention(tt.message, tt.botUsername); got != tt.want {
				t.Errorf("stripMention() = %q, want %q", got, tt.want)
			}
		})
	}
}

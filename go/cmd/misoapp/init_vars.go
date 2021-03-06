package main

var (
	target_sounds []string
	target_text   []string
	target_regexp []string
	missile  int
)

const ReportInterval = 15

func init() {
	missile = 0
	target_sounds = []string{
		"1-misomiso",
		"2-siosio",
		"3-nyaan",
		"4-koyaan",
		"5-kitsune",
		"6-nnaaa",
		"7-yysk",
		"8-killmebaby",
	}

	target_text = []string{
		"みそみそ〜",
		"しおしお〜",
		"にゃーん",
		"こゃーん",
		"きつね",
		"んなぁ",
		"ゆゆ式",
		"キルミーベイベー!!!",
	}

	target_regexp = []string{
		"みそ|(みそ)",
		"しお|(しお)",
		"にゃ?ーん",
		"こゃ?ーん",
		"きつね",
		"んなぁ",
		"(ゆゆ(式|しき)|yysk|yuyush?iki)",
		"(キルミー)|(ベイベ?ー)",
	}
}


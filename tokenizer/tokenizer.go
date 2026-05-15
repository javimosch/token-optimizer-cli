package tokenizer

import "math"

type Model string

const (
	GPT4    Model = "gpt-4"
	GPT4o   Model = "gpt-4o"
	Claude3 Model = "claude-3"
	Generic Model = "generic"
)

type Config struct {
	Model Model
}

var defaultRates = map[Model]float64{
	GPT4:    4.0,
	GPT4o:   4.5,
	Claude3: 3.5,
	Generic: 4.0,
}

func Count(text string, model Model) int {
	rate, ok := defaultRates[model]
	if !ok {
		rate = defaultRates[Generic]
	}
	return int(math.Ceil(float64(len([]rune(text))) / rate))
}

func CountFile(text string) FileCount {
	return FileCount{
		Chars:      len([]rune(text)),
		Lines:      countLines(text),
		TokenEstimate: Count(text, Generic),
		TokensGPT4: Count(text, GPT4),
		TokensGPT4o: Count(text, GPT4o),
		TokensClaude: Count(text, Claude3),
	}
}

type FileCount struct {
	Chars         int `json:"chars"`
	Lines         int `json:"lines"`
	TokenEstimate int `json:"token_estimate"`
	TokensGPT4    int `json:"tokens_gpt4"`
	TokensGPT4o   int `json:"tokens_gpt4o"`
	TokensClaude  int `json:"tokens_claude"`
}

func countLines(s string) int {
	n := 0
	for _, c := range s {
		if c == '\n' {
			n++
		}
	}
	if len(s) > 0 && s[len(s)-1] != '\n' {
		n++
	}
	return n
}

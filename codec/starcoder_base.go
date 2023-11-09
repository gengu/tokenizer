package codec

import "github.com/dlclark/regexp2"

func NewStarCoder() *Codec {
	return &Codec{
		name:        "starcoder",
		vocabulary:  starcoderVocab,
		splitRegexp: regexp2.MustCompile(`(?i:'s|'t|'re|'ve|'m|'ll|'d)|[^\r\n\p{L}\p{N}]?\p{L}+|\p{N}{1,3}| ?[^\s\p{L}\p{N}]+[\r\n]*|\s*[\r\n]+|\s+(?!\S)|\s+`, regexp2.None),
		specialTokens: map[string]uint{
			"<|endoftext|>":    0,
			"<fim_prefix>":     1,
			"<fim_middle>":     2,
			"<fim_suffix>":     3,
			"<fim_pad>":        4,
			"<filename>":       5,
			"<gh_stars>":       6,
			"<issue_start>":    7,
			"<issue_comment>":  8,
			"<issue_closed>":   9,
			"<jupyter_start>":  10,
			"<jupyter_text>":   11,
			"<jupyter_code>":   12,
			"<jupyter_output>": 13,
			"<empty_output>":   14,
			"<commit_before>":  15,
			"<commit_msg>":     16,
			"<commit_after>":   17,
			"<reponame>":       18,
		},
	}
}

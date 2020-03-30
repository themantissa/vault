package random

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// CharsetRestriction requires a certain number of characters from the specified charset.
type CharsetRestriction struct {
	// Charset is the list of rules that candidate strings must contain a minimum number of.
	Charset []rune `mapstructure:"charset"`

	// MinChars indicates the minimum (inclusive) number of characters from the charset that should appear in the string.
	MinChars int `mapstructure:"min-chars"`
}

// ParseCharsetRestriction from the provided data map. The data map is expected to be parsed from HCL.
func ParseCharsetRestriction(data map[string]interface{}) (rule Rule, err error) {
	cr := &CharsetRestriction{}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:   nil,
		Result:     cr,
		DecodeHook: stringToRunesFunc,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to decode charset restriction: %w", err)
	}

	err = decoder.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse charset restriction: %w", err)
	}

	return cr, nil
}

// Chars returns the charset that this rule is looking for.
func (c CharsetRestriction) Chars() []rune {
	return c.Charset
}

// Pass returns true if the provided candidate string has a minimum number of chars in it.
// This adheres to the Rule interface
func (c CharsetRestriction) Pass(value []rune) bool {
	if c.MinChars <= 0 {
		return true
	}

	count := 0
	for _, r := range value {
		// charIn is sometimes faster than a map lookup because the data is so small
		// This is being kept rather than converted to a map to keep the code cleaner,
		// otherwise there would need to be additional parsing logic.
		if charIn(r, c.Charset) {
			count++
			if count >= c.MinChars {
				return true
			}
		}
	}

	return false
}

func charIn(search rune, charset []rune) bool {
	for _, r := range charset {
		if search == r {
			return true
		}
	}
	return false
}
package msi

import (
	"strings"
)

const (
	DIGITAL_SIGNATURE_STREAM_NAME        = "\x05DigitalSignature"
	MSI_DIGITAL_SIGNATURE_EX_STREAM_NAME = "\x05MsiDigitalSignatureEx"
	SUMMARY_INFO_STREAM_NAME             = "\x05SummaryInformation"

	TABLE_PREFIX = "\xE4\xA1\x80"
)

func NameEncode(name string, isTable bool) string {
	var sb strings.Builder
	if isTable {
		sb.WriteString(TABLE_PREFIX)
	}

	len := len(name)
	for i, j := 0, 1; i < len; i, j = i+1, j+1 {
		char := rune(name[i])
		if val1, match := toB64(char); match {
			//peek
			if j < len {
				charNext := rune(name[j])
				if val2, match := toB64(charNext); match {
					encoded := 0x3800 + (val2 << 6) + val1
					i++
					j++
					sb.WriteRune(rune(encoded))
					continue
				}
			}
			encoded := 0x4800 + val1
			sb.WriteRune(rune(encoded))
		} else {
			sb.WriteRune(char)
		}
	}

	return sb.String()
}

func NameDecode(name string) (string, bool) {
	var sb strings.Builder
	var isTable bool

	if hasTablePrefix(name) {
		isTable = true
	}

	for pos, char := range name {
		//skip table prefix
		if isTable && pos == 0 {
			continue
		}

		if contains := containsRune(0x3800, 0x4800, char); contains {
			val := char - 0x3800
			if vf, match := fromB64(val & 0x3f); match {
				sb.WriteRune(vf)
				vf, _ = fromB64(val >> 6)
				sb.WriteRune(vf)
			}
		} else if contains := containsRune(0x4800, 0x4840, char); contains {
			if vf, match := fromB64(char - 0x4800); match {
				sb.WriteRune(vf)
			}
		} else {
			sb.WriteRune(char)
		}
	}

	return sb.String(), isTable
}

func NameIsValid(name string, isTable bool) bool {
	if name == "" || (!isTable && hasTablePrefix(name)) {
		return false
	}

	enc := NameEncode(name, isTable)
	return len(enc) <= 31

}

func hasTablePrefix(name string) bool {
	return name[:len(TABLE_PREFIX)] == TABLE_PREFIX
}

func toB64(ch rune) (uint32, bool) {
	for i := '0'; i <= '9'; i++ {
		if i == ch {
			return uint32(ch - '0'), true
		}
	}

	for i := 'A'; i <= 'Z'; i++ {
		if i == ch {
			return uint32(10 + ch - 'A'), true
		}
	}

	for i := 'a'; i <= 'z'; i++ {
		if i == ch {
			return uint32(36 + ch - 'a'), true
		}
	}

	if ch == '.' {
		return 62, true
	}

	if ch == '_' {
		return 63, true
	}

	return 0, false
}

func fromB64(val rune) (rune, bool) {
	if val < 10 {
		return rune('0' + val), true
	}

	if val < 36 {
		return rune('A' + val - 10), true
	}

	if val < 62 {
		return rune('a' + val - 36), true
	}

	if val == 62 {
		return '.', true
	}

	if val == 63 {
		return '_', true
	}

	return 0, false
}

func containsRune(min, max int, search rune) bool {
	for i := min; i < max; i++ {
		if rune(i) == search {
			return true
		}
	}
	return false
}

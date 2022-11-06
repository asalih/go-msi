package msi

import (
	"fmt"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

type CodePage int

const (
	/// [Windows-932 (Japanese Shift JIS)](https://en.wikipedia.org/wiki/Code_page_932_(Microsoft_Windows))
	Windows932 CodePage = iota
	/// [Windows-936 (Chinese (simplified) GBK)](https://en.wikipedia.org/wiki/Code_page_936_(Microsoft_Windows))
	Windows936
	/// [Windows-949 (Korean Unified Hangul Code)](https://en.wikipedia.org/wiki/Unified_Hangul_Code)
	Windows949
	/// [Windows-950 (Chinese (traditional) Big5)](https://en.wikipedia.org/wiki/Code_page_950)
	Windows950
	/// [Windows-951 (Chinese (traditional) Big5-HKSCS)](https://en.wikipedia.org/wiki/Hong_Kong_Supplementary_Character_Set#Microsoft_Windows)
	Windows951
	/// [Windows-1250 (Latin 2)](https://en.wikipedia.org/wiki/Windows-1250)
	Windows1250
	/// [Windows-1251 (Cyrillic)](https://en.wikipedia.org/wiki/Windows-1251)
	Windows1251
	/// [Windows-1252 (Latin 1)](https://en.wikipedia.org/wiki/Windows-1252)
	Windows1252
	/// [Windows-1253 (Greek)](https://en.wikipedia.org/wiki/Windows-1253)
	Windows1253
	/// [Windows-1254 (Turkish)](https://en.wikipedia.org/wiki/Windows-1254)
	Windows1254
	/// [Windows-1255 (Hebrew)](https://en.wikipedia.org/wiki/Windows-1255)
	Windows1255
	/// [Windows-1256 (Arabic)](https://en.wikipedia.org/wiki/Windows-1256)
	Windows1256
	/// [Windows-1257 (Baltic)](https://en.wikipedia.org/wiki/Windows-1257)
	Windows1257
	/// [Windows-1258 (Vietnamese)](https://en.wikipedia.org/wiki/Windows-1258)
	Windows1258
	/// [Mac OS Roman](https://en.wikipedia.org/wiki/Mac_OS_Roman)
	MacintoshRoman
	/// [Macintosh
	/// Cyrillic](https://en.wikipedia.org/wiki/Macintosh_Cyrillic_encoding)
	MacintoshCyrillic
	/// [US-ASCII](https://en.wikipedia.org/wiki/ASCII)
	UsAscii
	/// [ISO-8859-1 (Latin 1)](https://en.wikipedia.org/wiki/ISO-8859-1)
	Iso88591
	/// [ISO-8859-2 (Latin 2)](https://en.wikipedia.org/wiki/ISO-8859-2)
	Iso88592
	/// [ISO-8859-3 (South European)](https://en.wikipedia.org/wiki/ISO-8859-3)
	Iso88593
	/// [ISO-8859-4 (North European)](https://en.wikipedia.org/wiki/ISO-8859-4)
	Iso88594
	/// [ISO-8859-5 (Cyrillic)](https://en.wikipedia.org/wiki/ISO-8859-5)
	Iso88595
	/// [ISO-8859-6 (Arabic)](https://en.wikipedia.org/wiki/ISO-8859-6)
	Iso88596
	/// [ISO-8859-7 (Greek)](https://en.wikipedia.org/wiki/ISO-8859-7)
	Iso88597
	/// [ISO-8859-8 (Hebrew)](https://en.wikipedia.org/wiki/ISO-8859-8)
	Iso88598
	/// [UTF-8](https://en.wikipedia.org/wiki/UTF-8)
	Utf8
)

// Returns the code page (if any) with the given ID number.
func CodePageFromID(id int) CodePage {
	switch id {
	case 0:
		return CodePageDefault()
	case 932:
		return Windows932
	case 936:
		return Windows936
	case 949:
		return Windows949
	case 950:
		return Windows950
	case 951:
		return Windows951
	case 1250:
		return Windows1250
	case 1251:
		return Windows1251
	case 1252:
		return Windows1252
	case 1253:
		return Windows1253
	case 1254:
		return Windows1254
	case 1255:
		return Windows1255
	case 1256:
		return Windows1256
	case 1257:
		return Windows1257
	case 1258:
		return Windows1258
	case 10000:
		return MacintoshRoman
	case 10007:
		return MacintoshCyrillic
	case 20127:
		return UsAscii
	case 28591:
		return Iso88591
	case 28592:
		return Iso88592
	case 28593:
		return Iso88593
	case 28594:
		return Iso88594
	case 28595:
		return Iso88595
	case 28596:
		return Iso88596
	case 28597:
		return Iso88597
	case 28598:
		return Iso88598
	case 65001:
		return Utf8
	default:
		return -1
	}
}

func CodePageDefault() CodePage {
	return Utf8
}

// Returns the ID number used within Windows to represent this code page.
func (c CodePage) ID() int {
	switch c {
	case Windows932:
		return 932
	case Windows936:
		return 936
	case Windows949:
		return 949
	case Windows950:
		return 950
	case Windows951:
		return 951
	case Windows1250:
		return 1250
	case Windows1251:
		return 1251
	case Windows1252:
		return 1252
	case Windows1253:
		return 1253
	case Windows1254:
		return 1254
	case Windows1255:
		return 1255
	case Windows1256:
		return 1256
	case Windows1257:
		return 1257
	case Windows1258:
		return 1258
	case MacintoshRoman:
		return 10000
	case MacintoshCyrillic:
		return 10007
	case UsAscii:
		return 20127
	case Iso88591:
		return 28591
	case Iso88592:
		return 28592
	case Iso88593:
		return 28593
	case Iso88594:
		return 28594
	case Iso88595:
		return 28595
	case Iso88596:
		return 28596
	case Iso88597:
		return 28597
	case Iso88598:
		return 28598
	case Utf8:
		return 65001
	default:
		return -1
	}
}

func (c CodePage) Decode(data []byte) (string, error) {
	enc := c.Encoding()
	if enc == nil {
		return "", fmt.Errorf("unsupported code page: %d", c)
	}

	dec := enc.NewDecoder()
	result, err := dec.Bytes(data)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (c CodePage) Encoding() encoding.Encoding {

	switch c {
	case Windows932:
		return japanese.EUCJP
	case Windows936:
		return simplifiedchinese.GBK
	case Windows949:
		return korean.EUCKR
	case Windows950:
		return traditionalchinese.Big5
	case Windows951:
		return traditionalchinese.Big5
	case Windows1250:
		return charmap.Windows1250
	case Windows1251:
		return charmap.Windows1251
	case Windows1252:
		return charmap.Windows1252
	case Windows1253:
		return charmap.Windows1253
	case Windows1254:
		return charmap.Windows1254
	case Windows1255:
		return charmap.Windows1255
	case Windows1256:
		return charmap.Windows1256
	case Windows1257:
		return charmap.Windows1257
	case Windows1258:
		return charmap.Windows1258
	case MacintoshRoman:
		return charmap.Macintosh
	case MacintoshCyrillic:
		return charmap.MacintoshCyrillic
	case UsAscii:
		return nil
	case Iso88591:
		return charmap.ISO8859_1
	case Iso88592:
		return charmap.ISO8859_2
	case Iso88593:
		return charmap.ISO8859_3
	case Iso88594:
		return charmap.ISO8859_4
	case Iso88595:
		return charmap.ISO8859_5
	case Iso88596:
		return charmap.ISO8859_6
	case Iso88597:
		return charmap.ISO8859_7
	case Iso88598:
		return charmap.ISO8859_8
	case Utf8:
		return nil
	default:
		return nil
	}
}

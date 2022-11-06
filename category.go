package msi

type Category int

const (
	/// An unrestricted text string.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/text) for this
	/// data type.
	CategoryText = iota
	/// A text string containing no lowercase letters.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/uppercase) for
	/// this data type.
	CategoryUpperCase
	/// A text string containing no uppercase letters.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/lowercase) for
	/// this data type.
	CategoryLowerCase
	/// A signed 16-bit integer.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/integer) for
	/// this data type.
	CategoryInteger
	/// A signed 32-bit integer.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/doubleinteger)
	/// for this data type.
	CategoryDoubleInteger
	/// Stores a civil datetime, with a 2-second resolution.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/time-date) for
	/// this data type.
	CategoryTimeDate
	/// A string identifier (such as a table or column name).  May only contain
	/// alphanumerics, underscores, and periods, and must start with a letter
	/// or underscore.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/identifier)
	/// for this data type.
	CategoryIdentifier
	/// A string that is either an identifier (see above), or a reference to an
	/// environment variable (which consists of a `%` character followed by an
	/// identifier).
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/property) for
	/// this data type.
	CategoryProperty
	/// The name of a file or directory.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/filename) for
	/// this data type.
	CategoryFilename
	/// A filename that can contain shell glob wildcards.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/wildcardfilename)
	/// for this data type.
	CategoryWildCardFilename
	/// A string containing an absolute filepath.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/path) for this
	/// data type.
	CategoryPath
	/// A string containing a semicolon-separated list of absolute filepaths.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/paths) for
	/// this data type.
	CategoryPaths
	/// A string containing an absolute or relative filepath.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/anypath) for
	/// this data type.
	CategoryAnyPath
	/// A string containing either a filename or an identifier.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/defaultdir)
	/// for this data type.
	CategoryDefaultDir
	/// A string containing a registry path.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/regpath) for
	/// this data type.
	CategoryRegPath
	/// A string containing special formatting escapes, such as environment
	/// variables.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/formatted) for
	/// this data type.
	CategoryFormatted
	/// A security descriptor definition language (SDDL) text string written in
	/// valid [Security Descriptor String
	/// Format](https://docs.microsoft.com/en-us/windows/win32/secauthz/security-descriptor-string-format).
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/formattedsddltext)
	/// for this data type.
	CategoryFormattedSddlText
	/// Like `Formatted`, but allows additional escapes.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/template) for
	/// this data type.
	CategoryTemplate
	/// A string represeting a boolean predicate.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/condition) for
	/// this data type.
	CategoryCondition
	/// A hyphenated, uppercase GUID string, enclosed in curly braces.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/guid) for
	/// this data type.
	CategoryGuid
	/// A string containing a version number.  The string must consist of at
	/// most four period-separated numbers, with the value of each number being
	/// at most 65535.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/version) for
	/// this data type.
	CategoryVersion
	/// A string containing a comma-separated list of decimal language ID
	/// numbers.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/language) for
	/// this data type.
	CategoryLanguage
	/// A string that refers to a binary data stream.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/binary) for
	/// this data type.
	CategoryBinary
	/// A string that refers to a custom source.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/customsource)
	/// for this data type.
	CategoryCustomSource
	/// A string that refers to a cabinet.  If it starts with a `#` character,
	/// then the rest of the string is an identifier (see above) indicating a
	/// data stream in the package where the cabinet is stored.  Otherwise, the
	/// string is a short filename (at most eight characters, a period, and a
	/// three-character extension).
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/cabinet) for
	/// this data type.
	///
	/// # Examples
	///
	/// ```
	/// // Valid:
	/// assert!(msi::Category::Cabinet.validate("hello.txt"));
	/// assert!(msi::Category::Cabinet.validate("#HelloWorld"));
	/// // Invalid:
	/// assert!(!msi::Category::Cabinet.validate("longfilename.long"));
	/// assert!(!msi::Category::Cabinet.validate("#123.456"));
	/// ```
	CategoryCabinet
	/// A string that refers to a shortcut.
	///
	/// For more details, see the [MSI
	/// docs](https://docs.microsoft.com/en-us/windows/win32/msi/shortcut) for
	/// this data type.
	CategoryShortcut
)

var AllCategories = []Category{
	CategoryText, CategoryUpperCase, CategoryLowerCase, CategoryInteger,
	CategoryDoubleInteger, CategoryTimeDate, CategoryIdentifier, CategoryProperty,
	CategoryFilename, CategoryWildCardFilename, CategoryPath, CategoryPaths,
	CategoryAnyPath, CategoryDefaultDir, CategoryRegPath, CategoryFormatted,
	CategoryFormattedSddlText, CategoryTemplate, CategoryCondition, CategoryGuid,
	CategoryVersion, CategoryLanguage, CategoryBinary, CategoryCustomSource,
	CategoryCabinet, CategoryShortcut,
}

func (c Category) String() string {
	switch c {
	case CategoryAnyPath:
		return "AnyPath"
	case CategoryBinary:
		return "Binary"
	case CategoryCabinet:
		return "Cabinet"
	case CategoryCondition:
		return "Condition"
	case CategoryCustomSource:
		return "CustomSource"
	case CategoryDefaultDir:
		return "DefaultDir"
	case CategoryDoubleInteger:
		return "DoubleInteger"
	case CategoryFilename:
		return "Filename"
	case CategoryFormatted:
		return "Formatted"
	case CategoryFormattedSddlText:
		return "FormattedSddlText"
	case CategoryGuid:
		return "GUID"
	case CategoryIdentifier:
		return "Identifier"
	case CategoryInteger:
		return "Integer"
	case CategoryLanguage:
		return "Language"
	case CategoryLowerCase:
		return "LowerCase"
	case CategoryPath:
		return "Path"
	case CategoryPaths:
		return "Paths"
	case CategoryProperty:
		return "Property"
	case CategoryRegPath:
		return "RegPath"
	case CategoryShortcut:
		return "Shortcut"
	case CategoryTemplate:
		return "Template"
	case CategoryText:
		return "Text"
	case CategoryTimeDate:
		return "TimeDate"
	case CategoryUpperCase:
		return "UpperCase"
	case CategoryVersion:
		return "Version"
	case CategoryWildCardFilename:
		return "WildCardFilename"
	default:
		return ""
	}
}

func CategoryFromString(str string) Category {
	switch str {
	case "AnyPath":
		return CategoryAnyPath
	case "Binary":
		return CategoryBinary
	case "Cabinet":
		return CategoryCabinet
	case "Condition":
		return CategoryCondition
	case "CustomSource":
		return CategoryCustomSource
	case "DefaultDir":
		return CategoryDefaultDir
	case "DoubleInteger":
		return CategoryDoubleInteger
	case "Filename":
		return CategoryFilename
	case "Formatted":
		return CategoryFormatted
	case "FormattedSddlText":
		return CategoryFormattedSddlText
	case "GUID", "Guid":
		return CategoryGuid
	case "Identifier":
		return CategoryIdentifier
	case "Integer":
		return CategoryInteger
	case "Language":
		return CategoryLanguage
	case "LowerCase":
		return CategoryLowerCase
	case "Path":
		return CategoryPath
	case "Paths":
		return CategoryPaths
	case "Property":
		return CategoryProperty
	case "RegPath":
		return CategoryRegPath
	case "Shortcut":
		return CategoryShortcut
	case "Template":
		return CategoryTemplate
	case "Text":
		return CategoryText
	case "TimeDate":
		return CategoryTimeDate
	case "UpperCase":
		return CategoryUpperCase
	case "Version":
		return CategoryVersion
	case "WildCardFilename":
		return CategoryWildCardFilename
	default:
		return -1
	}
}

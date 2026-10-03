// Package validate provides simple syntactic validators for common
// user-input formats.
//
// Homework — Task 2 (Lesson 6: File I/O, JSON and Testing):
// Implement ValidateEmail and/or ValidatePhone below (your mentor may
// ask for just one) and extend the test tables in validate_test.go to
// at least 8 cases each, including edge cases.
package validate

import "strings"

// ValidateEmail reports whether s is a syntactically valid email address.
//
// This validator intentionally accepts a conservative, common subset of
// email addresses: the local part contains ASCII dot-atom characters, and
// the domain contains ASCII letters, digits, dots, and hyphens. It rejects
// Unicode and quoted local parts, as well as leading, trailing, or
// consecutive dots. The local part is limited to 64 bytes and the domain to
// 253 bytes. Domain labels are not checked individually; a single-label
// domain (for example, "student@localhost") is allowed.
func ValidateEmail(s string) bool {
	if s == "" || strings.Count(s, "@") != 1 {
		return false
	}

	at := strings.IndexByte(s, '@')
	local, domain := s[:at], s[at+1:]
	if len(local) == 0 || len(local) > 64 || len(domain) == 0 || len(domain) > 253 {
		return false
	}

	// Checking the bytes also rejects whitespace and all other non-ASCII
	// characters. unicode.IsSpace is not needed because neither side allows
	// any character outside the explicit ASCII sets below.
	if local[0] == '.' || local[len(local)-1] == '.' || strings.Contains(local, "..") {
		return false
	}
	for i := 0; i < len(local); i++ {
		c := local[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || strings.ContainsRune(".!#$%&'*+-/=?^_`{|}~", rune(c)) {
			continue
		}
		return false
	}

	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") || strings.Contains(domain, "..") {
		return false
	}
	for i := 0; i < len(domain); i++ {
		c := domain[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '.' || c == '-' {
			continue
		}
		return false
	}

	return true
}

// ValidatePhone reports whether s is a syntactically valid phone number.
//
// The accepted formats are an international number consisting of '+' and
// 10–15 digits (with no separators), a ten-digit national number beginning
// with 0, or that same national number formatted as XXX-XXX-XXXX. Spaces,
// parentheses, dots, and other separator styles are deliberately rejected.
func ValidatePhone(s string) bool {
	if s == "" {
		return false
	}

	if strings.HasPrefix(s, "+") {
		digits := s[1:]
		if len(digits) < 10 || len(digits) > 15 || digits[0] == '0' {
			return false
		}
		return allDigits(digits)
	}

	if len(s) == 10 && s[0] == '0' {
		return allDigits(s)
	}

	if len(s) == 12 && s[3] == '-' && s[7] == '-' && s[0] == '0' {
		return allDigits(s[:3]) && allDigits(s[4:7]) && allDigits(s[8:])
	}

	return false
}

func allDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

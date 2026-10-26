// Package avatar creates deterministic initials-based avatars for Mystic Square users.
package avatar

import (
	"crypto/sha256"
	"fmt"
	"html"
	"strings"
	"unicode/utf8"
)

const maxNameRunes = 200

var palette = [...]string{
	"#1D4ED8", // blue
	"#6D28D9", // violet
	"#047857", // emerald
	"#B45309", // amber
	"#BE123C", // rose
	"#0F766E", // teal
}

// Result is the generated representation returned by the avatar Lambda.
type Result struct {
	Initials   string `json:"initials"`
	Background string `json:"background"`
	Foreground string `json:"foreground"`
	SVG        string `json:"svg"`
}

// Generate creates a stable avatar from userID and fullName. The same user and
// name always produce the same initials, colors, and SVG.
func Generate(userID uint64, fullName string) (Result, error) {
	name := strings.Join(strings.Fields(fullName), " ")
	if userID == 0 {
		return Result{}, fmt.Errorf("user_id must be positive")
	}
	if name == "" {
		return Result{}, fmt.Errorf("full_name is required")
	}
	if utf8.RuneCountInString(name) > maxNameRunes {
		return Result{}, fmt.Errorf("full_name must not exceed %d characters", maxNameRunes)
	}

	initials := makeInitials(name)
	digest := sha256.Sum256([]byte(fmt.Sprintf("%d:%s", userID, name)))
	background := palette[int(digest[0])%len(palette)]
	foreground := "#FFFFFF"

	return Result{
		Initials:   initials,
		Background: background,
		Foreground: foreground,
		SVG:        makeSVG(name, initials, background, foreground),
	}, nil
}

func makeInitials(name string) string {
	words := strings.Fields(name)
	if len(words) == 0 {
		return ""
	}
	if len(words) == 1 {
		runes := []rune(words[0])
		if len(runes) == 1 {
			return strings.ToUpper(string(runes[0]))
		}
		return strings.ToUpper(string(runes[0:2]))
	}

	first := []rune(words[0])
	last := []rune(words[len(words)-1])
	return strings.ToUpper(string([]rune{first[0], last[0]}))
}

func makeSVG(name, initials, background, foreground string) string {
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 128 128" role="img" aria-label="%s avatar"><rect width="128" height="128" rx="64" fill="%s"/><text x="64" y="64" dy=".35em" text-anchor="middle" font-family="Arial, sans-serif" font-size="48" font-weight="700" fill="%s">%s</text></svg>`,
		html.EscapeString(name), background, foreground, html.EscapeString(initials))
}

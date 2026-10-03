# Task 3 — AI blind test generation

## Function under test

```go
func ValidatePhone(s string) bool
```

The implementation's accepted formats are:

- `+` followed by 10–15 digits, with no separators;
- ten digits beginning with `0`; or
- the same national number formatted as `XXX-XXX-XXXX`.

Other separators, including spaces, dots, and parentheses, are rejected.

## Prompt you gave the AI

```text
Generate a complete table-driven Go test suite for this function. Include
valid inputs, invalid inputs, empty input, boundary lengths, letters,
whitespace, punctuation, and Unicode digits. Return only the test code.

func ValidatePhone(s string) bool
```

The function signature was provided without its implementation, existing test
cases, or format description. Therefore, the AI had to make assumptions about
which phone-number formats should be considered valid.

## Edge cases the AI found that you had missed

The generated suite included these cases that were not in my original table:

- A string containing only `+`.
- An international number beginning with `0` after the plus sign, such as
  `+0123456789`.
- A string containing two plus signs, such as `++380501234567`.
- Unicode numerals, such as `+٣٨٠٥٠١٢٣٤٥٦٧`.
- A plus sign followed by digits but fewer than the minimum number of digits.

These did not expose bugs in the implementation. `ValidatePhone` rejects the
empty or incomplete international form, requires the first international digit
to be nonzero, and checks every remaining character with `allDigits`, which
rejects both additional plus signs and non-ASCII numerals.

## Edge cases you had that the AI missed

The AI did not test the two accepted national formats separately:

- `0501234567` (ten digits without separators); and
- `050-123-4567` (the exact `XXX-XXX-XXXX` format).

It also did not test the explicit policy that spaces, dots, and parentheses are
not accepted. Those cases matter because a phone number can be written with
many different conventional separators, but this function deliberately accepts
only one hyphenated national format and an unseparated international format.

The AI also omitted the upper international boundary of 15 digits and the
sixteen-digit rejection case. I added both to make the length contract clear.

## Cases where the AI's expected output was wrong

Yes. The AI classified values such as `+380 50 123 4567` and
`+1-202-555-0100` as valid because they are recognizable real-world phone
formats. The implementation correctly rejects them according to the chosen
contract: international numbers may not contain separators, and hyphens are
allowed only in the exact national `XXX-XXX-XXXX` layout.

This was not an implementation bug. It showed that a function signature alone
cannot communicate a format policy. The expected result depends on the
requirements chosen by the developer.

## What you'd change about your own test-writing process after this

I would write the accepted input grammar down before writing the test table.
For a format validator, I would then create a small matrix for every accepted
format: a normal value, the minimum and maximum boundaries, and a value just
outside each boundary. I would also vary one character at a time to test
letters, Unicode digits, whitespace, punctuation, and misplaced separators.

Finally, I would review AI-generated expectations rather than treating them as
requirements. In particular, common real-world formatting conventions are not
necessarily valid for a deliberately strict validator such as this one.

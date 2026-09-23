# Homework Report · Lesson 3

## 1. Link to the code

[Lesson 3 source code in this repository](https://github.com/seeman512/golang-genai/tree/main/lesson-03)

## 2. Prompt from Section 2

[Interface extraction prompt](https://github.com/seeman512/golang-genai/blob/main/lesson-03/prompts/interface-extraction-prompt.md)

## 3. Mixed receivers and consumer-side interfaces

Mixing value and pointer receivers on the same type can be dangerous because
it makes method behavior and interface satisfaction less predictable. A value
receiver operates on a copy, while a pointer receiver can modify the original
value. In addition, a type's method set differs from that of its pointer, so a
value may no longer satisfy an interface when one of its methods uses a pointer
receiver. Using pointer receivers consistently for stateful types makes it
clear that methods operate on the original value and avoids accidental copies.

AI can help design consumer-side interfaces by inspecting how a service uses a
dependency and extracting only the methods required by that service. This
keeps the abstraction small, makes dependency injection straightforward, and
allows lightweight fakes or mocks to be used in tests without a real database.
The generated design should still be reviewed and verified against the
project's existing API, tests, and behavior.

## 4. Screenshots / video of the completed work

[Screen recording and screenshots](https://drive.google.com/file/d/1wU9mn5YuG-uukj0E-ZKW9DNWwx7OuXQL/view?usp=drive_link)

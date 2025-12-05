# gogrep

A lightweight regular expression engine implementation in Go that uses Non-deterministic Finite Automaton (NFA) for pattern matching.

## Overview

gogrep is a command-line tool that evaluates regular expression patterns against input strings. It implements a classic NFA-based regex engine using Thompson's construction algorithm to convert regular expressions into state machines.

## Features

- **Core regex operators:**
  - `.` - Concatenation (implicit)
  - `|` - Alternation (OR)
  - `*` - Kleene star (zero or more)
  - `+` - Plus (one or more)
  - `?` - Question (zero or one)
  - `()` - Grouping

- **Efficient NFA evaluation:** Uses state buffering to simulate NFA execution without backtracking
- **Infix to postfix conversion:** Automatically handles operator precedence and parentheses

## Installation

```bash
go build -o gogrep
```

## Usage

```bash
./gogrep <pattern> <string>
```

**Examples:**

```bash
# Match single character
./gogrep "a" "a"
# Output: is match !

# Alternation
./gogrep "a|b" "b"
# Output: is match !

# Kleene star
./gogrep "a*b" "aaab"
# Output: is match !

# Grouping with alternation
./gogrep "(a|b)c" "bc"
# Output: is match !

# One or more
./gogrep "a+b" "ab"
# Output: is match !

# Optional
./gogrep "a?b" "b"
# Output: is match !
```

## Architecture

The project consists of several key components:

### 1. State Machine (`statemachine/`)

- **`nfa.go`**: Implements Thompson's NFA construction
  - Converts infix expressions to postfix notation
  - Builds NFA fragments for each operator
  - Patches state transitions to create final automaton

- **`evaluate.go`**: NFA evaluation engine
  - Maintains active state buffers during evaluation
  - Steps through input characters
  - Determines match status

### 2. Primitives (`primitives/`)

- **`stack.go`**: Generic stack implementation used for expression parsing and NFA construction

### 3. Main (`main.go`)

- Command-line interface and entry point

## How It Works

1. **Parsing**: The input pattern is preprocessed to insert explicit concatenation operators
2. **Conversion**: Infix notation is converted to postfix using the Shunting Yard algorithm
3. **NFA Construction**: Postfix expression is converted to an NFA using Thompson's construction
4. **Evaluation**: The NFA is simulated against the input string using state buffering

### State Types

- **MATCH_STATE (257)**: Accepting state indicating a successful match
- **SPLIT_STATE (256)**: ε-transition state for alternation and repetition
- **Character states**: Transitions that consume specific input characters

## Limitations

- Currently supports only literal character matching (no character classes like `[a-z]`)
- No support for anchors (`^`, `$`)
- No support for escape sequences
- Patterns and strings must consist of ASCII characters

## Testing

Run the test suite:

```bash
go test ./...
```

## Project Structure

```
gogrep/
├── main.go                      # CLI entry point
├── primitives/
│   └── stack.go                 # Generic stack data structure
└── statemachine/
    ├── nfa.go                   # NFA construction
    ├── nfa_test.go             # NFA construction tests
    ├── evaluate.go              # NFA evaluation engine
    └── evaluate_test.go         # Evaluation tests
```

## License

This is an educational project demonstrating regex engine implementation.

## References

This implementation is based on Russ Cox's article series on regular expression matching:
- [Regular Expression Matching Can Be Simple And Fast](https://swtch.com/~rsc/regexp/regexp1.html)

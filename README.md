# Regerize

Regerize is a language that compiles to Golang regexp2 regular expressions.
It is heavily based on the work of (Melody)[https://github.com/yoav-lavi/melody], which utilizes Rust to build ECMAScript regular expressions. You should check out their work!

## Why?

Years ago, I worked for a company that had a large need to normalize HTML information to identify unique and repeating structures in DOMs with the hopes of building fully automated test cases across a domain. A heavy burden existed in creating regular expressions to normalize the documents reliably, and an even heaver burden of educating the masses about how to properly write regular expressions.

After seeing (Melody)[https://github.com/yoav-lavi/melody], I wanted to better my Golang writing while performing something similar. Oddly enough, I use regular expressions to generate regular expressions.

This project was fun and is in no way comprehensive. If you would like to use its features, please use (Melody)[https://github.com/yoav-lavi/melody] as it is fully supported!

## Examples

### Batman Theme

```rust
16 of "na";

2 of match {
  <space>;
  "batman";
}
```

Turns into

```regex
(na){16}\sbatman
```

### Web URLs

```rust
before {
  "www.";
};
some of <word>;
".com";
```

Turns into

```regex
(?<=www\.)\w+\.com
```

### Variables

```rust
let .banana {
    "ba";
    2 of "na";
};
let .ringing {
    8 of "ring ";
};
match {
    .ringing;
    ",";
    <space>;
    .banana;
    <space>;
    `phone`
};
```

Turns into

```regex
(?:(ring ){8},\sba(na){2}\sphone)
```

## Install

### Go Install

As of Go 1.18+, use `go install`

```sh
go install github.com/lawrencemq/regerize@latest
```

### From Source

```sh
git clone git@github.com:lawrencemq/regerize.git
cd regerize
go build
```

## Usage

```go
regex, err := parser.ParseFile(filename)
if err != nil {
  fmt.Println("Unable to parse file: ", err)
  return
}

```

## Syntax

### Quantifiers

- `... of` - used to express a specific amount of a pattern. equivalent to regex `{5}` (assuming `5 of ...`)
- `... to ... of` - used to express an amount within a range of a pattern. equivalent to regex `{5,9}` (assuming `5 to 9 of ...`)
- `over ... of` - used to express more than an amount of a pattern. equivalent to regex `{6,}` (assuming `over 5 of ...`)
- `some of` - used to express 1 or more of a pattern. equivalent to regex `+`
- `any of` - used to express 0 or more of a pattern. equivalent to regex `*`
- `option of` - used to express 0 or 1 of a pattern. equivalent to regex `?`

### Symbols

- `<char>` - matches any single character. equivalent to regex `.`
- `<space>` - matches a space character. equivalent to regex ` `
- `<whitespace>` - matches any kind of whitespace character. equivalent to regex `\s` or `[ \t\n\v\f\r]`
- `<newline>` - matches a newline character. equivalent to regex `\n`
- `<tab>` - matches a tab character. equivalent to regex `\t`
- `<return>` - matches a carriage return character. equivalent to regex `\r`
- `<feed>` - matches a form feed character. equivalent to regex `\f`
- `<null>` - matches a null characther. equivalent to regex `\0`
- `<num>` - matches any single digit. equivalent to regex `\d` or `[0-9]`
- `<!num>` - matches any single non-digit. equivalent to regex `[!\d]` or `[!0-9]`
- `<vertical>` - matches a vertical tab character. equivalent to regex `\v`
- `<alphanum>` - matches a word character (any latin letter, any digit or an underscore). equivalent to regex `\w` or `[a-zA-Z0-9_]`
- `<!alphanum>` - matches a non-word character (any latin letter, any digit or an underscore). equivalent to regex `[!\w]` or `[!a-zA-Z0-9_]`
- `<alpha>` - matches any single latin letter. equivalent to regex `[a-zA-Z]`
- `<!alpha>` - matches any single non-latin letter. equivalent to regex `[!a-zA-Z]`
- `<hex>` - matches any hex value regardless of case. equivalent to regex `[0-9a-fA-F]`
- `<start>` - matches the beginning of a line. equivalent to regex `^`
- `<end>` - matches the en dof a line. equivalent to regex `$`

### Character Ranges

- `... to ...` - used with digits or alphabetic characters to express a character range. equivalent to regex `[5-9]` (assuming `5 to 9`) or `[a-z]` (assuming `a to z`)

### Literals

- `"..."` or `'...'` - used to mark a literal part of the match. Melody will automatically escape characters as needed. Quotes (of the same kind surrounding the literal) should be escaped

### Raw

- <code>\`...\`</code> - added directly to the output without any escaping

### Groups

- `capture` - used to open a `capture` or named `capture` block. capture patterns are later available in the list of matches (either positional or named). equivalent to regex `(...)`
- `match` - used to open a `match` block, matches the contents without capturing. equivalent to regex `(?:...)`
- `either` - used to open an `either` block, matches one of the statements within the block. equivalent to regex `(?:...|...)`

### Assertions

- `ahead` - used to open an `ahead` block. equivalent to regex `(?=...)`. use after an expression
- `behind` - used to open an `behind` block. equivalent to regex `(?<=...)`. use before an expression

Assertions can be preceeded by `not` to create a negative assertion (equivalent to regex `(?!...)`, `(?<!...)`)

### Variables

- `let .variable_name = { ... }` - defines a variable from a block of statements. can later be used with `.variable_name`. Variables must be declared before being used. Variable invocations cannot be quantified directly, use a group if you want to quantify a variable invocation

  example:

  ```rs
  let .a_and_b = {
    "a";
    "b";
  }

  .a_and_b;
  "c";

  // abc
  ```

### Comments

- `/* ... */`, `// ...` - used to mark comments (note: `// ...` comments must be on separate line)

## File Extension

The Regerize file extension is `.rgr`

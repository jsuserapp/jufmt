// Package jufmt provides colored fmt/log-style printing with optional timestamps
// and call-site locations (file:line) for quick navigation in the IDE.
//
// It is a partial replacement for the standard fmt and log packages. Most call
// patterns are compatible. Because each print adds a prefix, output is usually
// one line at a time; chaining partial-line writes will look messy.
//
// Use SetLogOutputHook for structured logging (full file paths, plain messages)
// Return nil from the hook to keep default console formatting; see LogHookResult.
package jufmt

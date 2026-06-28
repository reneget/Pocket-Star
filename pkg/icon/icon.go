package icon

import (
	_ "embed"
)

//go:embed logo.svg
var Logo string

func String() string {
	return `
     ╭────╮
     │ ╱╲ │
     │ ╲╱ │
     ╰──╥─╯
    ╭───╨───╮
    │       │
    ╰───────╯
`
}

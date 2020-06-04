
## File: ./regexpscommands/regexpscommands.go
- Package: github.com/go-worker/regexpscommands
- Import
	- fmt: fmt
	- commands: github.com/go-worker/commands
	- regexps: github.com/go-worker/regexps
	- ui: github.com/go-worker/ui
- Function
	- ParseCommandByRegexp
		- Params
			- inputString: string
		- Return
			- 0: commands.CommandInterface
			- 1: error
		- Call
			- commands: CreateCommand
			- fmt: Errorf

## File: ./regexpscommands/regexpscommands.go
- Package: github.com/go-worker/regexpscommands
- Import
	- ui: github.com/go-worker/ui
	- fmt: fmt
	- commands: github.com/go-worker/commands
	- regexps: github.com/go-worker/regexps
- Function
	- ParseCommandByRegexp
		- Params
			- inputString: string
		- Return
			- 0: CommandInterface
			- 1: error
		- Call
			- commands: CreateCommand
			- fmt: Errorf

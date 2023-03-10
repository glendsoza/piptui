# piptui -> Terminal UI to manage pip packages


## How to run
- Download the binary from [releases page](https://github.com/glendsoza/piptui/releases)
- Activate the virtual environemnt and run the binary

## How to build 

- Install [go](https://go.dev/doc/install)
- Go to the project root and run `go mod tidy && go run main.go` 

## Supported platforms

- Currently this tool is tested only on `linux` with `bash` (yes, it requires bash!)

## How to use
![Demo](./demo/demo.gif)

- go to the project root
- activate your python virtual environment
- on the same terminal session run `go run main.go`
- once the ui loads following short cuts can be used to go into different screens from the main menu (short-cuts to different screens will only work from the main menu )
  - `ctrl+I` -> Install screen
  - `ctrl+T` -> Dependency Tree
  - `ctrl+U` -> Uninstall screen
  - `ctrl+D` -> Switch tabs (if there are multiple tabs)
  - `Esc` -> go back to the main screen

  
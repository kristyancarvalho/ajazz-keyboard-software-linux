module github.com/kristyancarvalho/ak820pro

go 1.23.8

require (
	gioui.org v0.9.0
	gioui.org/x v0.9.0
	github.com/sstallion/go-hid v0.15.0
	github.com/stretchr/testify v1.11.1
	golang.org/x/image v0.26.0
)

require (
	gioui.org/shader v1.0.8 // indirect
	git.wow.st/gmp/jni v0.0.0-20210610011705-34026c7e22d0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-text/typesetting v0.3.0 // indirect
	github.com/godbus/dbus/v5 v5.0.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/exp/shiny v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/sstallion/go-hid => ./third_party/go-hid

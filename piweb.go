// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piweb

var (
	GoBuild        = "go build -buildvcs=false -o {{.Output}} ." // executed on each webpage refresh.
	ReleaseGoBuild = "go build -o {{.Output}} ."                 // executed when creating release build.
)

func Run() {
	run()
}

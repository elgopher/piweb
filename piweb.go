// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piweb

var (
	GoBuild        = "go build -buildvcs=false -o {{.Output}} ." // executed on each webpage refresh.
	ReleaseGoBuild = "go build -o {{.Output}} ."                 // executed when creating release build.
)

var HtmlDir = "." // change if you want to store html/javascript/css files in different directory than working one.

func Run() {
	run()
}

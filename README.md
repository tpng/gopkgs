Command gopkgs list your installed Go packages for import.

	$ go get github.com/tpng/gopkgs

It aims to provide a faster alternative to `go list all`
to list available packages for import.

It is extracted from the margo import paths function
(https://github.com/DisposaBoy/GoSublime/blob/master/src/gosubli.me/margo/m_import_paths.go)
bundled in GoSublime.

The difference of `go list all` and `gopkgs` is that `go list all` looks for go packages in your `$GOPATH/src` while `gopkgs` looks in your `$GOPATH/pkg`.
As a result of this, only importable packages that have been installed (either by go get or go install) are listed by `gopkgs`.

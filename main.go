// The MIT License (MIT)
// Copyright © 2015 DisposaBoy <disposaboy@dby.me>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the “Software”), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	// output GOROOT pkgs
	for _, p := range importPaths(runtime.GOROOT()) {
		fmt.Println(p)
	}
	fmt.Println("unsafe")

	// output GOPATH pkgs
	for _, p := range importPaths(os.Getenv("GOPATH")) {
		fmt.Println(p)
	}
}

func importPaths(envs ...string) (paths []string) {
	roots := map[string]bool{}
	for _, env := range envs {
		for _, root := range filepath.SplitList(env) {
			if root != "" {
				roots[root] = true
			}
		}
	}

	seen := map[string]bool{}
	pfx := strings.HasPrefix
	sfx := strings.HasSuffix
	osArchSfx := runtime.GOOS + "_" + runtime.GOARCH
	for root, _ := range roots {
		root = resolveSymlink(filepath.Join(root, "pkg", osArchSfx))
		walkF := func(p string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}

			p, e := filepath.Rel(root, p)
			if e != nil || !sfx(p, ".a") {
				return nil
			}

			p = p[:len(p)-2]
			if pfx(p, ".") || pfx(p, "_") || sfx(p, "_test") {
				return nil
			}

			p = path.Clean(filepath.ToSlash(p))
			if !seen[p] {
				seen[p] = true
				paths = append(paths, p)
			}

			return nil
		}
		filepath.Walk(root, walkF)
	}
	return paths
}

func resolveSymlink(path string) string {
	resolvePath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path
	}
	return resolvePath
}

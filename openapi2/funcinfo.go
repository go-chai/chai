package openapi2

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/go-chai/chai/internal/log"
)

type funcInfo struct {
	Pkg          string         `json:"pkg"`
	PkgPath      string         `json:"pkgPath"`
	Name         string         `json:"name"`
	Func         string         `json:"func"`
	Comment      string         `json:"comment"`
	File         string         `json:"file,omitempty"`
	ASTFile      *ast.File      `json:"ast_file,omitempty"`
	FSet         *token.FileSet `json:"fset,omitempty"`
	Line         int            `json:"line,omitempty"`
	Anonymous    bool           `json:"anonymous,omitempty"`
	Unresolvable bool           `json:"unresolvable,omitempty"`

	fn any
}

func (fi *funcInfo) Dump() {
	log.Printf("-----------\n")
	log.Printf("PkgPath %q\n", fi.PkgPath)
	log.Printf("Name %q\n", fi.Name)
	log.Printf("Pkg: %q\n", fi.Pkg)
	log.Printf("Func %q\n", fi.Func)
	log.Printf("Line %d\n", fi.Line)
	log.Printf("-----------\n")
}

func getFuncInfo(fn any) funcInfo {
	return getFuncInfoWithSrc(fn, nil)
}

func getFuncInfoWithSrc(fn any, src any) funcInfo {
	fi := funcInfo{
		fn: fn,
	}
	frame := getCallerFrame(fn)
	goPathSrc := filepath.Join(os.Getenv("GOPATH"), "src")

	if frame == nil {
		fi.Unresolvable = true
		return fi
	}

	pkgName := getPkgName(frame.File, src)
	fi.PkgPath = reflect.TypeOf(fn).PkgPath()
	fi.Name = reflect.TypeOf(fn).Name()

	funcPath := frame.Func.Name()

	idx := strings.Index(funcPath, "/"+pkgName)
	if idx > 0 {
		fi.Pkg = funcPath[:idx+1+len(pkgName)]
		fi.Func = funcPath[idx+2+len(pkgName):]
	} else {
		fi.Func = funcPath
	}

	if strings.Index(fi.Func, ".func") > 0 {
		fi.Anonymous = true
	}

	fi.File = frame.File
	fi.Line = frame.Line
	if filepath.HasPrefix(fi.File, goPathSrc) {
		fi.File = fi.File[len(goPathSrc)+1:]
	}

	if !fi.Unresolvable {
		fi.Comment, fi.ASTFile, fi.FSet = getFuncComment(frame.File, frame.Line, src)
	}
	return fi
}

func getCallerFrame(fn any) *runtime.Frame {
	pc := reflect.ValueOf(fn).Pointer()
	frames := runtime.CallersFrames([]uintptr{pc})
	if frames == nil {
		return nil
	}

	frame, _ := frames.Next()
	if frame.Entry == 0 {
		return nil
	}
	return &frame
}

func getPkgName(file string, src any) string {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, file, src, parser.PackageClauseOnly)
	if err != nil {
		return ""
	}
	if astFile.Name == nil {
		return ""
	}
	return astFile.Name.Name
}

func getFuncComment(file string, line int, src any) (string, *ast.File, *token.FileSet) {
	fset := token.NewFileSet()

	astFile, err := parser.ParseFile(fset, file, src, parser.ParseComments)
	if err != nil {
		return "", nil, nil
	}

	if len(astFile.Comments) == 0 {
		return "", astFile, fset
	}

	line = fixFuncLine(line, fset, astFile)
	for _, cmt := range astFile.Comments {
		if fset.Position(cmt.End()).Line+1 == line {
			return cmt.Text(), astFile, fset
		}
	}

	return "", astFile, fset
}

func pos(fset *token.FileSet, n ast.Node) int {
	return fset.Position(n.Pos()).Line
}

// If the compiler inlined the function, we get the line of the return statement rather than the line of the function definition.
// fixFuncLine checks if the specified line contains any return statements
// and if so, returns the line of the function definition that the first return belongs to.
func fixFuncLine(line int, fset *token.FileSet, astFile *ast.File) int {
	fixedFuncLine := line

	var stack []ast.Node
	ast.Inspect(astFile, func(n ast.Node) bool {
		if n != nil {
			stack = append(stack, n)
		} else {
			stack = stack[:len(stack)-1]
		}

		// Check if the current node is on the specified line.
		if n == nil || fset.Position(n.Pos()).Line != line {
			return true
		}
		// Check if the current node is a return statement.
		_, ok := n.(*ast.ReturnStmt)
		if !ok {
			return true
		}
		// Starting at the return statement, go up the node stack until we find the first function definition
		for i := len(stack) - 1; i >= 0; i-- {
			switch parent := stack[i].(type) {
			case *ast.FuncDecl, *ast.FuncLit:
				fixedFuncLine = pos(fset, parent)

				// Stop looking after the function definition of the first return statement
				return false
			}
		}

		return true
	})

	return fixedFuncLine
}

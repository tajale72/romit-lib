package sequence

import (
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"
)

func SequenceGenarotor() {
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, ".", nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	funcDecls := make(map[string]*goast.FuncDecl)
	var allCalls []Call
	seenEdges := make(map[edgeKey]bool)

	// Collect all function declarations from all files in the current directory
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				if fn, ok := decl.(*goast.FuncDecl); ok {
					funcDecls[fn.Name.Name] = fn
				}
			}
		}
	}

	var visitFunc func(caller string, fn *goast.FuncDecl)
	visitFunc = func(caller string, fn *goast.FuncDecl) {
		if fn == nil || fn.Body == nil {
			return
		}

		goast.Inspect(fn.Body, func(n goast.Node) bool {
			callExpr, ok := n.(*goast.CallExpr)
			if !ok {
				return true
			}

			switch fun := callExpr.Fun.(type) {
			case *goast.SelectorExpr:
				if ident, ok := fun.X.(*goast.Ident); ok && ident.Name == "http" {
					url := getArgValue(callExpr.Args)
					key := edgeKey{Caller: caller, Callee: extractServiceFromURL(url), Note: fmt.Sprintf("HTTP %s %s", fun.Sel.Name, url)}
					if !seenEdges[key] {
						allCalls = append(allCalls, Call{Caller: caller, Callee: key.Callee, Note: key.Note})
						seenEdges[key] = true
					}
				}

			case *goast.Ident:
				if _, exists := funcDecls[fun.Name]; exists {
					key := edgeKey{Caller: caller, Callee: fun.Name, Note: "internal call"}
					if !seenEdges[key] {
						allCalls = append(allCalls, Call{Caller: caller, Callee: fun.Name, Note: key.Note})
						seenEdges[key] = true
					}
					visitFunc(fun.Name, funcDecls[fun.Name])
				}
			}
			return true
		})
	}

	// Start with callService (entry point)
	if fn, ok := funcDecls["main"]; ok {
		visitFunc("callService", fn)
	}

	mermaidFile := "diagram.mmd"
	os.Create(mermaidFile)
	if err != nil {
		fmt.Println("Error creating diagram file:", err)
		return
	}
	f, err := os.OpenFile("diagram.mmd", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening diagram file:", err)
	}

	f.WriteString("sequenceDiagram\n")

	// Output diagram
	fmt.Println("```mermaid")
	fmt.Println("sequenceDiagram")
	for _, c := range allCalls {
		str := fmt.Sprintf("    %s->>%s: %s\n", c.Caller, c.Callee, c.Note)
		f.WriteString(str)
	}
	generateImageFromMermaid(mermaidFile)
	fmt.Println("```")
}

func generateImageFromMermaid(mermaidFile string) {
	cmd := exec.Command("mmdc", "-i", mermaidFile, "-o", "diagram.png")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error running mmdc:", err)
	} else {
		fmt.Println("✅ diagram.png generated")
	}
}

func getArgValue(args []goast.Expr) string {
	if len(args) == 0 {
		return ""
	}
	if bl, ok := args[0].(*goast.BasicLit); ok {
		return strings.Trim(bl.Value, `"`)
	}
	return "<dynamic>"
}

func extractServiceFromURL(url string) string {
	switch {
	case strings.Contains(url, "service-b"):
		return "ServiceB"
	case strings.Contains(url, "service-c"):
		return "ServiceC"
	}
	return "ExternalService"
}

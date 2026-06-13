package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

func TestDocHasResetTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cg   *ast.CommentGroup
		want bool
	}{
		{name: "nil", cg: nil, want: false},
		{name: "empty", cg: &ast.CommentGroup{}, want: false},
		{
			name: "wrong text",
			cg:   &ast.CommentGroup{List: []*ast.Comment{{Text: "// some comment"}}},
			want: false,
		},
		{
			name: "has tag",
			cg:   &ast.CommentGroup{List: []*ast.Comment{{Text: "// generate:reset"}}},
			want: true,
		},
		{
			name: "tag among others",
			cg: &ast.CommentGroup{List: []*ast.Comment{
				{Text: "// some comment"},
				{Text: "// generate:reset"},
				{Text: "// another comment"},
			}},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, docHasResetTag(tt.cg))
		})
	}
}

func TestResetCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expr     string
		typ      ast.Expr
		contains string
		empty    bool
	}{
		{
			name:     "int",
			expr:     "s.X",
			typ:      &ast.Ident{Name: "int"},
			contains: "s.X = 0",
		},
		{
			name:     "string",
			expr:     "s.Name",
			typ:      &ast.Ident{Name: "string"},
			contains: `s.Name = ""`,
		},
		{
			name:     "bool",
			expr:     "s.Flag",
			typ:      &ast.Ident{Name: "bool"},
			contains: "s.Flag = false",
		},
		{
			name:     "float64",
			expr:     "s.Val",
			typ:      &ast.Ident{Name: "float64"},
			contains: "s.Val = 0",
		},
		{
			name:     "named struct",
			expr:     "s.Sub",
			typ:      &ast.Ident{Name: "SubStruct"},
			contains: "resetter",
		},
		{
			name:     "slice",
			expr:     "s.Items",
			typ:      &ast.ArrayType{Len: nil, Elt: &ast.Ident{Name: "int"}},
			contains: "s.Items = s.Items[:0]",
		},
		{
			name:  "fixed array",
			expr:  "s.Arr",
			typ:   &ast.ArrayType{Len: &ast.BasicLit{}, Elt: &ast.Ident{Name: "int"}},
			empty: true,
		},
		{
			name:     "map",
			expr:     "s.M",
			typ:      &ast.MapType{Key: &ast.Ident{Name: "string"}, Value: &ast.Ident{Name: "int"}},
			contains: "clear(s.M)",
		},
		{
			name:     "selector expr",
			expr:     "s.T",
			typ:      &ast.SelectorExpr{X: &ast.Ident{Name: "time"}, Sel: &ast.Ident{Name: "Time"}},
			contains: "resetter",
		},
		{
			name:     "pointer",
			expr:     "s.Ptr",
			typ:      &ast.StarExpr{X: &ast.Ident{Name: "int"}},
			contains: "s.Ptr != nil",
		},
		{
			name:  "channel (unknown)",
			expr:  "s.Ch",
			typ:   &ast.ChanType{},
			empty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := resetCode(tt.expr, tt.typ)
			if tt.empty {
				require.Empty(t, got)
			} else {
				require.Contains(t, got, tt.contains)
			}
		})
	}
}

func TestPtrResetCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expr     string
		inner    ast.Expr
		contains string
		empty    bool
	}{
		{
			name:     "ptr to int",
			expr:     "s.X",
			inner:    &ast.Ident{Name: "int"},
			contains: "*s.X = 0",
		},
		{
			name:     "ptr to named struct",
			expr:     "s.Sub",
			inner:    &ast.Ident{Name: "SubStruct"},
			contains: "resetter",
		},
		{
			name:     "ptr to slice",
			expr:     "s.Items",
			inner:    &ast.ArrayType{Len: nil, Elt: &ast.Ident{Name: "int"}},
			contains: "*s.Items",
		},
		{
			name:  "ptr to fixed array",
			expr:  "s.Arr",
			inner: &ast.ArrayType{Len: &ast.BasicLit{}, Elt: &ast.Ident{Name: "int"}},
			empty: true,
		},
		{
			name:     "ptr to map",
			expr:     "s.M",
			inner:    &ast.MapType{Key: &ast.Ident{Name: "string"}, Value: &ast.Ident{Name: "int"}},
			contains: "clear(*s.M)",
		},
		{
			name:     "ptr to selector",
			expr:     "s.T",
			inner:    &ast.SelectorExpr{X: &ast.Ident{Name: "time"}, Sel: &ast.Ident{Name: "Time"}},
			contains: "resetter",
		},
		{
			name:  "ptr to channel (unknown)",
			expr:  "s.Ch",
			inner: &ast.ChanType{},
			empty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ptrResetCode(tt.expr, tt.inner)
			if tt.empty {
				require.Empty(t, got)
			} else {
				require.Contains(t, got, tt.contains)
			}
		})
	}
}

func parseSource(t *testing.T, src string) *ast.File {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	require.NoError(t, err)
	return f
}

func TestCollectStructs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		src       string
		wantLen   int
		wantNames []string
	}{
		{
			name:    "no declarations",
			src:     `package test`,
			wantLen: 0,
		},
		{
			name: "struct without tag",
			src: `package test
type Foo struct { X int }`,
			wantLen: 0,
		},
		{
			name: "struct with tag",
			src: `package test
// generate:reset
type Foo struct { X int }`,
			wantLen:   1,
			wantNames: []string{"Foo"},
		},
		{
			name: "multiple structs mixed",
			src: `package test

// generate:reset
type Foo struct { X int }

type Bar struct { Y string }

// generate:reset
type Baz struct { Z bool }
`,
			wantLen:   2,
			wantNames: []string{"Foo", "Baz"},
		},
		{
			name: "non-struct type with tag",
			src: `package test
// generate:reset
type MyInt int`,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := parseSource(t, tt.src)
			got := collectStructs(f)
			require.Len(t, got, tt.wantLen)
			for i, name := range tt.wantNames {
				require.Equal(t, name, got[i].name)
			}
		})
	}
}

func TestBuildFileData(t *testing.T) {
	t.Parallel()

	t.Run("empty structs", func(t *testing.T) {
		t.Parallel()
		data := buildFileData("mypkg", nil)
		require.Equal(t, "mypkg", data.Package)
		require.Empty(t, data.Structs)
	})

	t.Run("struct with various fields", func(t *testing.T) {
		t.Parallel()
		src := `package test
// generate:reset
type Sample struct {
	Count   int
	Name    string
	Flag    bool
	Items   []int
	Mapping map[string]int
	Sub     SubType
	Ptr     *int
}
`
		f := parseSource(t, src)
		structs := collectStructs(f)
		require.Len(t, structs, 1)

		data := buildFileData("test", structs)

		require.Equal(t, "test", data.Package)
		require.Len(t, data.Structs, 1)

		sd := data.Structs[0]
		require.Equal(t, "Sample", sd.Name)
		require.Equal(t, "s", sd.Recv)
		require.Contains(t, sd.Body, "s.Count = 0")
		require.Contains(t, sd.Body, `s.Name = ""`)
		require.Contains(t, sd.Body, "s.Flag = false")
		require.Contains(t, sd.Body, "s.Items = s.Items[:0]")
		require.Contains(t, sd.Body, "clear(s.Mapping)")
	})

	t.Run("receiver is first letter lowercase", func(t *testing.T) {
		t.Parallel()
		src := `package test
// generate:reset
type ResetableStruct struct { X int }
`
		f := parseSource(t, src)
		structs := collectStructs(f)
		data := buildFileData("test", structs)

		require.Equal(t, "r", data.Structs[0].Recv)
	})
}

func TestWriteGenFile(t *testing.T) {
	t.Parallel()

	tmpl := template.Must(template.New("reset").Parse(genFileTemplate))
	dir := t.TempDir()

	data := fileData{
		Package: "testpkg",
		Structs: []structData{
			{
				Name: "Foo",
				Recv: "f",
				Body: "f.X = 0\n",
			},
		},
	}

	err := writeGenFile(tmpl, dir, data)
	require.NoError(t, err)

	outPath := filepath.Join(dir, "reset.gen.go")
	content, err := os.ReadFile(outPath)
	require.NoError(t, err)

	src := string(content)
	require.Contains(t, src, "package testpkg")
	require.Contains(t, src, "func (f *Foo) Reset()")
	require.Contains(t, src, "f.X = 0")
	require.Contains(t, src, "DO NOT EDIT")
}

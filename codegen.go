// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"net/textproto"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
)

type Generator struct {
	buf bytes.Buffer
}

func (g *Generator) Printf(s string, args ...interface{}) {
	fmt.Fprintf(&g.buf, s, args...)
}

func (g *Generator) WriteFile(name string) error {
	src, err := g.Format()
	if err != nil {
		return fmt.Errorf("format: %s: %s:\n\n%s\n", name, err, g.Bytes())
	} else if err := ioutil.WriteFile(name, src, 0644); err != nil {
		return err
	}
	return nil
}

func (g *Generator) Bytes() []byte {
	return g.buf.Bytes()
}

func (g *Generator) Format() ([]byte, error) {
	return format.Source(g.Bytes())
}

func (g *Generator) generate() error {
	g.Printf(`// Code generated by codegen.go; DO NOT EDIT
	
	package s3protocol
	
	import (
		"fmt"
		"net/http"
		"net/url"
		"strconv"
	
		"github.com/aws/aws-sdk-go/aws"
		"github.com/aws/aws-sdk-go/service/s3"
	)
	`)
	if err := g.generateInput(s3.GetObjectInput{}); err != nil {
		return err
	}
	if err := g.generateInput(s3.HeadObjectInput{}); err != nil {
		return err
	}
	if err := g.generateOutput(s3.GetObjectOutput{}); err != nil {
		return err
	}
	if err := g.generateOutput(s3.HeadObjectOutput{}); err != nil {
		return err
	}
	return nil
}

var typeTime = reflect.TypeOf(time.Time{})

func (g *Generator) generateInput(target interface{}) error {
	typ := reflect.TypeOf(target)
	name := typ.Name()

	g.Printf("func new%s(req *http.Request) (*s3.%s, error) {\n", name, name)
	g.Printf("var in s3.%s\n", name)
	g.Printf(`header := req.Header
	if header == nil {
		header = make(http.Header)
	}
	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return nil, err
	}
	`)
	num := typ.NumField()
	for i := 0; i < num; i++ {
		f := typ.Field(i)
		tag := f.Tag
		name := tag.Get("locationName")
		switch tag.Get("location") {
		case "header":
			name = textproto.CanonicalMIMEHeaderKey(name)
			g.Printf("if v, ok := header[%q]; ok && len(v) > 0 {\n", name)
		case "querystring":
			g.Printf("if v, ok := query[%q]; ok && len(v) > 0 {\n", name)
		default:
			continue
		}

		switch f.Type.Elem().Kind() {
		case reflect.String:
			g.Printf("in.%s = aws.String(v[0])\n", f.Name)
		case reflect.Int64:
			g.Printf(`i, err := strconv.ParseInt(v[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("s3protocol: failed to parse %s: %%v", err)
			}
			in.%s = aws.Int64(i)
			`, name, f.Name)
		case reflect.Struct:
			if f.Type.Elem() == typeTime {
				g.Printf(`t, err := http.ParseTime(v[0])
				if err != nil {
					return nil, fmt.Errorf("s3protocol: failed to parse %s: %%v", err)
				}
				in.%s = aws.Time(t)
				`, name, f.Name)
			} else {
				return fmt.Errorf("unknown type: %v", f.Type.Elem())
			}
		default:
			return fmt.Errorf("unknown type: %v", f.Type.Elem())
		}
		g.Printf("}\n")
	}
	g.Printf("return &in, nil\n}\n\n")
	return nil
}

func (g *Generator) generateOutput(target interface{}) error {
	typ := reflect.TypeOf(target)
	name := typ.Name()

	g.Printf("func makeHeaderFrom%[1]s(out *s3.%[1]s) (http.Header, error) {\n", name)
	g.Printf("header := make(http.Header)\n")
	num := typ.NumField()
	for i := 0; i < num; i++ {
		f := typ.Field(i)
		tag := f.Tag
		if tag.Get("location") != "header" {
			continue
		}
		name := textproto.CanonicalMIMEHeaderKey(tag.Get("locationName"))
		g.Printf("if out.%s != nil {\n", f.Name)
		switch f.Type.Elem().Kind() {
		case reflect.Bool:
			g.Printf("header.Set(%q, strconv.FormatBool(aws.BoolValue(out.%s)))\n", name, f.Name)
		case reflect.String:
			g.Printf("header.Set(%q, aws.StringValue(out.%s))\n", name, f.Name)
		case reflect.Int64:
			g.Printf("header.Set(%q, strconv.FormatInt(aws.Int64Value(out.%s), 10))\n", name, f.Name)
		case reflect.Struct:
			if f.Type.Elem() == typeTime {
				g.Printf("header.Set(%q, out.%s.Format(http.TimeFormat))\n", name, f.Name)
			} else {
				return fmt.Errorf("unknown type: %v", f.Type.Elem())
			}
		default:
			return fmt.Errorf("unknown type: %v", f.Type.Elem())
		}
		g.Printf("}\n")
	}
	g.Printf("return header, nil\n}\n\n")
	return nil
}

func main() {
	var g Generator
	if err := g.generate(); err != nil {
		log.Fatal(err)
	}
	if err := g.WriteFile("generated.go"); err != nil {
		log.Fatal(err)
	}
}

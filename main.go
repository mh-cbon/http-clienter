// Package http-clienter generates http client of a type
package main

//go:generate lister string:utils/StringSlice

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mh-cbon/astutil"
	"github.com/mh-cbon/http-clienter/utils"
	httper "github.com/mh-cbon/httper/lib"
)

var name = "http-clienter"
var version = "0.0.0"

func main() {

	var help bool
	var h bool
	var ver bool
	var v bool
	var outPkg string
	var mode string
	flag.BoolVar(&help, "help", false, "Show help.")
	flag.BoolVar(&h, "h", false, "Show help.")
	flag.BoolVar(&ver, "version", false, "Show version.")
	flag.BoolVar(&v, "v", false, "Show version.")
	flag.StringVar(&outPkg, "p", "", "Package name of the new code.")
	flag.StringVar(&mode, "mode", "gorilla", "The genration mode gorilla|std.")

	flag.Parse()

	if ver || v {
		showVer()
		return
	}
	if help || h {
		showHelp()
		return
	}

	if mode != "gorilla" && mode != "std" {
		wrongInput("invalid mode value %q", mode)
		return
	}

	if flag.NArg() < 1 {
		wrongInput("not enough type to trasnform")
		return
	}
	args := flag.Args()

	out := ""
	if args[0] == "-" {
		args = args[1:]
		out = "-"
	}

	todos, err := utils.NewTransformsArgs(utils.GetPkgToLoad()).Parse(args)
	if err != nil {
		panic(err)
	}

	filesOut := utils.NewFilesOut("github.com/mh-cbon/" + name)

	for _, todo := range todos.Args {
		if todo.FromPkgPath == "" {
			log.Println("Skipped ", todo.FromTypeName)
			continue
		}

		fileOut := filesOut.Get(todo.ToPath)

		fileOut.PkgName = outPkg
		if fileOut.PkgName == "" {
			fileOut.PkgName = findOutPkg(todo)
		}

		if err := processType(mode, todo, fileOut); err != nil {
			log.Println(err)
		}
	}
	filesOut.Write(out)
}

func wrongInput(r string, a ...interface{}) {
	showHelp()
	fmt.Printf(`

Please check your input parameters: %v
`, fmt.Sprintf(r, a...))
}

func showVer() {
	fmt.Printf("%v %v\n", name, version)
}

func showHelp() {
	showVer()
	fmt.Println()
	fmt.Println("Usage")
	fmt.Println()
	fmt.Printf("	%v [-p name] [-mode name] [...types]\n\n", name)
	fmt.Printf("  types:  A list of types such as src:dst.\n")
	fmt.Printf("          A type is defined by its package path and its type name,\n")
	fmt.Printf("          [pkgpath/]name\n")
	fmt.Printf("          If the Package path is empty, it is set to the package name being generated.\n")
	// fmt.Printf("          If the Package path is a directory relative to the cwd, and the Package name is not provided\n")
	// fmt.Printf("          the package path is set to this relative directory,\n")
	// fmt.Printf("          the package name is set to the name of this directory.\n")
	fmt.Printf("          Name can be a valid type identifier such as TypeName, *TypeName, []TypeName \n")
	fmt.Printf("  -p:     The name of the package output.\n")
	fmt.Printf("  -mode:  The generation mode gorilla|std.\n")
	fmt.Println()
}

func findOutPkg(todo utils.TransformArg) string {
	if todo.ToPkgPath != "" {
		prog := astutil.GetProgramFast(todo.ToPkgPath)
		if prog != nil {
			pkg := prog.Package(todo.ToPkgPath)
			return pkg.Pkg.Name()
		}
	}
	if todo.ToPkgPath == "" {
		prog := astutil.GetProgramFast(utils.GetPkgToLoad())
		if len(prog.Imported) < 1 {
			panic("impossible, add [-p name] option")
		}
		for _, p := range prog.Imported {
			return p.Pkg.Name()
		}
	}
	if strings.Index(todo.ToPkgPath, "/") > -1 {
		return filepath.Base(todo.ToPkgPath)
	}
	return todo.ToPkgPath
}

func processType(mode string, todo utils.TransformArg, fileOut *utils.FileOut) error {
	dest := &fileOut.Body
	srcName := todo.FromTypeName
	destName := todo.ToTypeName

	prog := astutil.GetProgramFast(todo.FromPkgPath)
	pkg := prog.Package(todo.FromPkgPath)
	foundMethods := astutil.FindMethods(pkg)

	srcConcrete := astutil.GetUnpointedType(srcName)
	// the json input must provide a key/value for each params.
	structType := astutil.FindStruct(pkg, srcConcrete)
	structComment := astutil.GetComment(prog, structType.Pos())
	// todo: might do better to send only annotations or do other improvemenets.
	structComment = makeCommentLines(structComment)
	structAnnotations := astutil.GetAnnotations(structComment, "@")

	dstConcrete := astutil.GetUnpointedType(destName)

	// Declare the new type

	fmt.Fprintf(dest, `
// %v is an http-clienter of %v.
%v`, dstConcrete, srcName, structComment)

	if mode == gorillaMode {
		fmt.Fprintf(dest, `
	type %v struct{
		router *mux.Router
	  Base string
	}
			`, dstConcrete)

	} else {
		fmt.Fprintf(dest, `
	type %v struct{
	  Base string
	}
			`, dstConcrete)
	}

	// Make the constructor
	// should param *http.Client be an interface ?
	fmt.Fprintf(dest, `// New%v constructs an http-clienter of %v
`, dstConcrete, srcName)

	if mode == gorillaMode {
		fmt.Fprintf(dest, `func New%v(router *mux.Router) *%v {
	ret := &%v{
		router: router,
	}
  return ret
}
`, dstConcrete, dstConcrete, dstConcrete)

	} else {
		fmt.Fprintf(dest, `func New%v() *%v {
	ret := &%v{
	}
  return ret
}
`, dstConcrete, dstConcrete, dstConcrete)
	}

	if mode == gorillaMode {
		fileOut.AddImport("fmt", "")
		fileOut.AddImport("io", "")
		fileOut.AddImport("net/http", "")
		fileOut.AddImport("net/url", "")
		fileOut.AddImport("strings", "")
		fileOut.AddImport("github.com/gorilla/mux", "")
	} else {
		fileOut.AddImport("net/http", "")
	}

	for _, m := range foundMethods[srcConcrete] {
		methodName := astutil.MethodName(m)

		comment := astutil.GetComment(prog, m.Pos())
		annotations := astutil.GetAnnotations(comment, "@")
		annotations = mergeAnnotations(structAnnotations, annotations)
		params := astutil.MethodParams(m)
		lParams := commaArgsToSlice(params)
		paramNames := astutil.MethodParamNames(m)
		lParamNames := commaArgsToSlice(paramNames)

		if mode == "std" {

			importIDs := astutil.GetSignatureImportIdentifiers(m)
			for _, i := range importIDs {
				fileOut.AddImport(astutil.GetImportPath(pkg, i), i)
			}

			fileOut.AddImport("errors", "")

			fmt.Fprintf(dest, `// %v constructs a request to %v
		`, methodName, methodName)

			fmt.Fprintf(dest, `func(t %v) %v(%v) (*http.Request, error) {
			return nil, errors.New("todo")
		}
		`, destName, methodName, params)

		} else if route, ok := annotations["route"]; ok {

			importIDs := astutil.GetSignatureImportIdentifiers(m)
			for _, i := range importIDs {
				fileOut.AddImport(astutil.GetImportPath(pkg, i), i)
			}

			getParams := ""
			postParams := ""
			routeName, _ := annotations["name"]

			// - look for every route params
			managedParamNames := utils.NewStringSlice()
			routeParamsExpr := []string{}
			routeParamNames := getRouteParamsFromRoute(mode, route)
			for _, p := range routeParamNames {
				routeParamsExpr = append(routeParamsExpr, fmt.Sprintf("%q", p))
				routeParamsExpr = append(routeParamsExpr, p)
				methodParam := getMethodParamForRouteParam(mode, lParamNames, p, managedParamNames)
				if methodParam == "" {
					log.Println("route param not identified into the method parameters " + p)
					continue
				}
				managedParamNames.Push(methodParam)
			}

			// - look for url/req/post params, not already managed by the route params
			for _, p := range lParamNames {
				if p == reqBodyVarName {
					continue
				}
				if !managedParamNames.Contains(p) {
					prefix := getVarPrefix(mode, p)
					rParamName := getVarValueName(mode, p)
					if prefix == "get" || prefix == "url" || prefix == "req" {
						getParams += fmt.Sprintf("url.Query().Add(%q, %v)", rParamName, p)
						managedParamNames.Push(p)

					} else if prefix == "post" {
						getParams += fmt.Sprintf("form.Add(%q, %v)", rParamName, p)
						managedParamNames.Push(p)
					}
				}
			}

			// - forge url from the router using the route name
			url := ""
			if routeName != "" {
				k := ""
				if len(routeParamsExpr) > 0 {
					k = strings.Join(routeParamsExpr, ", ")
					k = k[:len(k)-2]
				}
				url = fmt.Sprintf(`url, URLerr := t.router.Get(%q).URL(%v)
									`, routeName, k)
			} else {
				// - a route without name neeeds a jit update.
				url += fmt.Sprintf(`surl := %q
										`, route)

				managedParamNames = utils.NewStringSlice()
				for _, p := range routeParamNames {
					methodParam := getMethodParamForRouteParam(mode, lParamNames, p, managedParamNames)
					if methodParam == "" {
						log.Println("route param not identified into the method parameters " + p)
						continue
					}
					url += fmt.Sprintf(`surl = strings.Replace(surl, "{%v}", fmt.Sprintf("%%v", %v), 1)
													`, p, methodParam)
					managedParamNames.Push(methodParam)
				}

				url += fmt.Sprintf(`url, URLerr := url.ParseRequestURI(surl)
									`)
			}
			url += handleErr("URLerr")

			// - if any GET params, handle them
			if getParams != "" {
				url += fmt.Sprintf("%v\n", getParams)
			}
			// - if any GET params, handle them
			if postParams != "" {
				url += fmt.Sprintf("form := url.Values{}\n%v\n", postParams)
			}

			url += fmt.Sprint("finalURL := url.String()\n")

			// - build the final url
			if base, ok := annotations["base"]; ok {
				url += fmt.Sprintf("finalURL = fmt.Sprint(%q, %q, finalURL)\n", "%v%v", base)
			}
			url += fmt.Sprintf("finalURL = fmt.Sprintf(%q, t.Base, finalURL)\n", "%v%v")

			// modify method params to transform a reqBody ? to reqBody io.Reader
			methodParams := changeParamType(lParams, "reqBody", "io.Reader")

			body := ""
			// - handle the request body
			if postParams != "" {
				body += fmt.Sprintf("body = strings.NewReader(form.Encode())\n")

			} else if hasReqBody(lParamNames) {
				body += fmt.Sprintf("body = %v\n", "reqBody")
			}

			// - create the request object
			preferedMethod := getPreferredMethod(annotations)
			body += fmt.Sprintf("%v\n", url)
			body += fmt.Sprintf(" req, reqErr := http.NewRequest(%q, finalURL, body)\n", preferedMethod)
			body += handleErr("reqErr")
			body += fmt.Sprintf("ret = req\n")

			// - print the method
			fmt.Fprintf(dest, "// %v constructs a request to %v\n", methodName, route)
			fmt.Fprintf(dest, `func(t %v) %v(%v) (*http.Request, error) {
					        var ret *http.Request
					        var body io.Reader
					        // var err error
					        %v
					        return ret , nil
					      }
					      `, destName, methodName, strings.Join(methodParams, ","), body)
		}

	}

	return nil
}

func getMethodParamForRouteParam(mode string,
	methodParamNames []string,
	routeParamName string,
	managed *utils.StringSlice) string {
	for _, methodParamName := range methodParamNames {
		prefix := getVarPrefix(mode, methodParamName)
		if prefix == "route" || prefix == "get" || prefix == "req" || prefix == "url" {
			valueName := getVarValueName(mode, methodParamName)
			if strings.ToLower(valueName) == strings.ToLower(routeParamName) {
				return methodParamName
			}
		}
	}
	return ""
}

func handleErr(errVarName string) string {
	return fmt.Sprintf(`if %v!= nil {
return nil, %v
}
`, errVarName, errVarName)
}

func hasReqBody(paramNames []string) bool {
	for _, p := range paramNames {
		if p == "reqBody" {
			return true
		}
	}
	return false
}

func getReqBodyType(params []string) string {
	for _, p := range params {
		k := strings.Split(p, " ")
		k[0] = strings.TrimSpace(k[0])
		if k[0] == "reqBody" && len(k) > 1 {
			return strings.TrimSpace(k[1])
		}
	}
	return ""
}

func changeParamType(lParams []string, name, t string) []string {
	ret := []string{}
	for _, p := range lParams {
		p = strings.TrimSpace(p)
		if strings.Index(p, name) == 0 {
			p = name + " " + t
		}
		ret = append(ret, p)
	}
	return ret
}

var re = regexp.MustCompile(`({[^}]+(:[^}]+|)})`)

// func routeHasParam(mode, route, paramName string) bool {
// 	if mode == gorillaMode {
// 		for _, p := range getRouteParamNamesFromRoute(mode, route) {
// 			if p == paramName {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

func getRouteParamsFromRoute(mode, route string) []string {
	ret := []string{}
	if mode == gorillaMode {
		//todo: find a better way.
		res := re.FindAllStringSubmatch(route, -1)
		for _, r := range res {
			if len(r) > 0 {
				k := strings.TrimSpace(r[0])
				if len(k) > 2 { // there is braces inside
					k = k[1 : len(k)-1]
					ret = append(ret, k)
				}
			}
		}
	}
	return ret
}

func getRouteParamNamesFromRoute(mode, route string) []string {
	ret := []string{}
	if mode == gorillaMode {
		//todo: find a better way.
		res := re.FindAllStringSubmatch(route, -1)
		for _, r := range res {
			if len(r) > 0 {
				k := strings.TrimSpace(r[0])
				if len(k) > 2 { // there is braces inside
					k = k[1 : len(k)-1]
					if strings.Index(k, ":") > -1 {
						j := strings.Split(k, ":")
						ret = append(ret, j[0])
					} else {
						ret = append(ret, k)
					}
				}
			}
		}
	}
	return ret
}

var gorillaMode = "gorilla"
var stdMode = "std"
var reqBodyVarName = "reqBody"

func isConvetionnedParam(mode, varName string) bool {
	if varName == reqBodyVarName {
		return true
	}
	return getVarPrefix(mode, varName) != ""
}

func getDataProviderFactory(mode string) httper.DataerProvider {
	var factory httper.DataerProvider
	if mode == stdMode {
		factory = &httper.StdHTTPDataProvider{}
	} else if mode == gorillaMode {
		factory = &httper.GorillaHTTPDataProvider{}
	}
	return factory
}

func getDataProvider(mode string) *httper.DataProviderFacade {
	return getDataProviderFactory(mode).MakeEmpty().(*httper.DataProviderFacade)
}

func getVarPrefix(mode, varName string) string {
	ret := ""
	provider := getDataProvider(mode)
	for _, p := range provider.Providers {
		prefix := p.GetName()
		if strings.HasPrefix(varName, strings.ToLower(prefix)) {
			f := string(varName[len(prefix):][0])
			if f == strings.ToUpper(f) {
				ret = prefix
				break
			}
		} else if strings.HasPrefix(varName, strings.ToUpper(prefix)) {
			f := string(varName[len(prefix):][0])
			if f == strings.ToLower(f) {
				ret = prefix
				break
			}
		}
	}
	return ret
}

func getVarValueName(mode, varName string) string {
	ret := ""
	provider := getDataProvider(mode)
	for _, p := range provider.Providers {
		prefix := p.GetName()
		if strings.HasPrefix(varName, strings.ToLower(prefix)) {
			f := string(varName[len(prefix):][0])
			if f == strings.ToUpper(f) {
				ret = varName[len(prefix):]
				break
			}
		} else if strings.HasPrefix(varName, strings.ToUpper(prefix)) {
			f := string(varName[len(prefix):][0])
			if f == strings.ToLower(f) {
				ret = varName[len(prefix):]
				break
			}
		}
	}
	return ret
}

func getPreferredMethod(annotations map[string]string) string {
	preferedMethod := "GET"
	if m, ok := annotations["metods"]; ok {
		methods := commaArgsToSlice(m)
		if len(methods) > 0 {
			preferedMethod = strings.ToUpper(methods[0])
		}
	}
	return preferedMethod
}

func mergeAnnotations(structAnnot, methodAnnot map[string]string) map[string]string {
	ret := map[string]string{}
	for k, v := range methodAnnot {
		ret[k] = v
	}
	for k, v := range structAnnot {
		if _, ok := ret[k]; !ok {
			ret[k] = v
		}
	}
	return ret
}

func makeCommentLines(s string) string {
	s = strings.TrimSpace(s)
	comment := ""
	for _, k := range strings.Split(s, "\n") {
		comment += "// " + k + "\n"
	}
	comment = strings.TrimSpace(comment)
	if comment == "" {
		comment = "//"
	}
	return comment
}

func commaArgsToSlice(s string) []string {
	ret := []string{}
	for _, l := range strings.Split(s, ",") {
		l = strings.TrimSpace(l)
		if l != "" {
			ret = append(ret, l)
		}
	}
	return ret
}

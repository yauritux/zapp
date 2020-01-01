package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode"

	"github.com/alecthomas/template"
)

// TheField is
type TheField struct {
	Name     string
	DataType string
}

// TheClass is
type TheClass struct {
	Name   string
	Fields []TheField
}

// ThePackage is
type ThePackage struct {
	appName     string
	packageName string
	// datatypeid  string
	Classes []TheClass
}

// CamelCase is
func CamelCase(name string) string {

	// force it!
	if name == "IPAddress" {
		return "ipAddress"
	}

	out := []rune(name)
	out[0] = unicode.ToLower([]rune(name)[0])
	return string(out)
}

// UpperCase is
func UpperCase(name string) string {
	return strings.ToUpper(name)
}

// LowerCase is
func LowerCase(name string) string {
	return strings.ToLower(name)
}

// PascalCase is
func PascalCase(name string) string {
	return name
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// SnakeCase is
func SnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// HasTime is
func HasTime(dataTypes []TheField) bool {
	for _, tm := range dataTypes {
		if tm.DataType == "time.Time" {
			return true
		}
	}
	return false
}

func main() {

	filename := "skrip.txt"

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	tp := ThePackage{}

	var theClass *TheClass

	for scanner.Scan() {
		row := scanner.Text()

		trimmedStr := strings.TrimSpace(row)
		if trimmedStr == "" {
			continue
		}

		if strings.HasPrefix(trimmedStr, "//") {
			continue
		}

		if strings.HasPrefix(row, "package") {
			pkgs := strings.Split(row, ",")
			packageName := pkgs[1]
			tp.packageName = strings.TrimSpace(packageName)
			tp.Classes = []TheClass{}
			s := strings.Split(packageName, "/")
			tp.appName = s[len(s)-1]
			continue
		}

		// if strings.HasPrefix(row, "datatypeid") {
		// 	dtt := strings.Split(row, ",")
		// 	tp.datatypeid = strings.TrimSpace(dtt[1])
		// 	continue
		// }

		if strings.HasPrefix(row, "table") {
			clName := strings.Split(row, " ")
			theClass = &TheClass{
				Name:   strings.TrimSpace(clName[1]),
				Fields: []TheField{},
			}
			continue
		}

		if strings.HasPrefix(row, "endtable") {
			tp.Classes = append(tp.Classes, *theClass)
			theClass = nil
			continue
		}

		if strings.HasPrefix(row, "field") {
			flField := strings.Split(row, ",")
			theClass.Fields = append(theClass.Fields, TheField{
				Name:     strings.TrimSpace(flField[1]),
				DataType: strings.TrimSpace(flField[2]),
			})
			continue
		}

	}

	tp.Run()

	cmd := exec.Command("/usr/local/bin/go", "fmt", fmt.Sprintf("%s/...", tp.packageName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

// Run is
func (tp *ThePackage) Run() {

	// create app folder
	{
		dir := fmt.Sprintf("../../../../%s/app", tp.packageName)
		os.MkdirAll(dir, 0777)
	}

	// create controller folder
	{
		dir := fmt.Sprintf("../../../../%s/controller", tp.packageName)
		os.MkdirAll(dir, 0777)
	}

	// create model folder
	{
		dir := fmt.Sprintf("../../../../%s/model/basic", tp.packageName)
		os.MkdirAll(dir, 0777)
	}

	// create service folder
	{
		dir := fmt.Sprintf("../../../../%s/service", tp.packageName)
		os.MkdirAll(dir, 0777)
	}

	// create dao folder
	{
		dir := fmt.Sprintf("../../../../%s/dao", tp.packageName)
		os.MkdirAll(dir, 0777)
	}

	// create utils folder
	{
		dir := ""

		dir = fmt.Sprintf("../../../../%s/utils/common", tp.packageName)
		os.MkdirAll(dir, 0777)

		dir = fmt.Sprintf("../../../../%s/utils/config", tp.packageName)
		os.MkdirAll(dir, 0777)

		dir = fmt.Sprintf("../../../../%s/utils/log", tp.packageName)
		os.MkdirAll(dir, 0777)

		dir = fmt.Sprintf("../../../../%s/utils/transaction", tp.packageName)
		os.MkdirAll(dir, 0777)
	}

	{
		dir := fmt.Sprintf("../../../../%s/webapp/public", tp.packageName)
		os.MkdirAll(dir, 0777)
	}

	// {
	// 	dir := fmt.Sprintf("../../../../%s/webapp/src/api/modules", tp.packageName)
	// 	os.MkdirAll(dir, 0777)
	// }

	{
		dir := fmt.Sprintf("../../../../%s/webapp/src/store/modules", tp.packageName)
		os.MkdirAll(dir, 0777)
	}

	{
		dir := fmt.Sprintf("../../../../%s/webapp/src/router/modules", tp.packageName)
		os.MkdirAll(dir, 0777)
	}

	{
		f := fmt.Sprintf("../../../../%s/webapp/src/pages", tp.packageName)
		os.Mkdir(f, 0777)
	}

	// {
	// 	f := fmt.Sprintf("../../../../%s/webapp/src/components", tp.packageName)
	// 	os.Mkdir(f, 0777)
	// }

	{
		f := fmt.Sprintf("../../../../%s/webapp/src/utils", tp.packageName)
		os.Mkdir(f, 0777)
	}

	for _, et := range tp.Classes {

		// create dao
		{
			templateFile := fmt.Sprintf("../templates/backend/dao._go")
			outputFile := fmt.Sprintf("../../../../%s/dao/%s.go", tp.packageName, LowerCase(et.Name))
			basic(tp, templateFile, outputFile, et)
		}

		// create controller
		{
			templateFile := fmt.Sprintf("../templates/backend/controller._go")
			outputFile := fmt.Sprintf("../../../../%s/controller/%s.go", tp.packageName, LowerCase(et.Name))
			basic(tp, templateFile, outputFile, et)
		}

		// create model
		{
			templateFile := fmt.Sprintf("../templates/backend/model._go")
			outputFile := fmt.Sprintf("../../../../%s/model/%s.go", tp.packageName, LowerCase(et.Name))
			basic(tp, templateFile, outputFile, et)
		}

		// create model
		{
			templateFile := fmt.Sprintf("../templates/backend/model_basic._go")
			outputFile := fmt.Sprintf("../../../../%s/model/basic/model.go", tp.packageName)
			basic(tp, templateFile, outputFile, et)
		}

		// create service
		{
			templateFile := fmt.Sprintf("../templates/backend/service._go")
			outputFile := fmt.Sprintf("../../../../%s/service/%s.go", tp.packageName, LowerCase(et.Name))
			basic(tp, templateFile, outputFile, et)
		}

		// create api modules
		// {
		// 	templateFile := fmt.Sprintf("../templates/frontend/src/api/modules/file._js")
		// 	outputFile := fmt.Sprintf("../../../../%s/webapp/src/api/modules/%s.js", tp.packageName, LowerCase(et.Name))
		// 	basic(tp, templateFile, outputFile, et)
		// }

		// create store modules
		// {
		// 	templateFile := fmt.Sprintf("../templates/frontend/src/store/modules/file._js")
		// 	outputFile := fmt.Sprintf("../../../../%s/webapp/src/store/modules/%s.js", tp.packageName, LowerCase(et.Name))
		// 	basic(tp, templateFile, outputFile, et)
		// }

		// create router modules
		{
			templateFile := fmt.Sprintf("../templates/frontend/src/router/modules/file._js")
			outputFile := fmt.Sprintf("../../../../%s/webapp/src/router/modules/%s.js", tp.packageName, LowerCase(et.Name))
			basic(tp, templateFile, outputFile, et)
		}

		{
			// create folder pages
			f := fmt.Sprintf("../../../../%s/webapp/src/pages/%s", tp.packageName, LowerCase(et.Name))
			os.Mkdir(f, 0777)

			// create file table under folder
			{
				templateFile := fmt.Sprintf("../templates/frontend/src/pages/folder/list._vue")
				outputFile := fmt.Sprintf("../../../../%s/webapp/src/pages/%s/list.vue", tp.packageName, LowerCase(et.Name))
				basic(tp, templateFile, outputFile, et)
			}

			// create file input under folder
			{
				templateFile := fmt.Sprintf("../templates/frontend/src/pages/folder/input._vue")
				outputFile := fmt.Sprintf("../../../../%s/webapp/src/pages/%s/input.vue", tp.packageName, LowerCase(et.Name))
				basic(tp, templateFile, outputFile, et)
			}

		}

	}

	{
		templateFile := fmt.Sprintf("../templates/backend/main._go")
		outputFile := fmt.Sprintf("../../../../%s/app/main.go", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/backend/router._go")
		outputFile := fmt.Sprintf("../../../../%s/app/router.go", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/backend/utils/common/identifier._go")
		outputFile := fmt.Sprintf("../../../../%s/utils/common/identifier.go", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/backend/utils/common/json._go")
		outputFile := fmt.Sprintf("../../../../%s/utils/common/json.go", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/backend/utils/common/strings._go")
		outputFile := fmt.Sprintf("../../../../%s/utils/common/strings.go", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/backend/utils/config/config._go")
		outputFile := fmt.Sprintf("../../../../%s/utils/config/config.go", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/backend/utils/log/log._go")
		outputFile := fmt.Sprintf("../../../../%s/utils/log/log.go", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/backend/utils/transaction/transaction._go")
		outputFile := fmt.Sprintf("../../../../%s/utils/transaction/transaction.go", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/backend/._gitignore")
		outputFile := fmt.Sprintf("../../../../%s/.gitignore", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/config._toml")
		outputFile := fmt.Sprintf("../../../../%s/config.toml", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/README._md")
		outputFile := fmt.Sprintf("../../../../%s/README.md", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/public/favicon._ico")
		outputFile := fmt.Sprintf("../../../../%s/webapp/public/favicon.ico", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/public/index._html")
		outputFile := fmt.Sprintf("../../../../%s/webapp/public/index.html", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/babel.config._js")
		outputFile := fmt.Sprintf("../../../../%s/webapp/babel.config.js", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/vue.config._js")
		outputFile := fmt.Sprintf("../../../../%s/webapp/vue.config.js", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/package._json")
		outputFile := fmt.Sprintf("../../../../%s/webapp/package.json", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/src/App._vue")
		outputFile := fmt.Sprintf("../../../../%s/webapp/src/App.vue", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/src/pages/forgotpassword._vue")
		outputFile := fmt.Sprintf("../../../../%s/webapp/src/pages/forgotpassword.vue", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/src/pages/home._vue")
		outputFile := fmt.Sprintf("../../../../%s/webapp/src/pages/home.vue", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/src/pages/login._vue")
		outputFile := fmt.Sprintf("../../../../%s/webapp/src/pages/login.vue", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/src/pages/notfound._vue")
		outputFile := fmt.Sprintf("../../../../%s/webapp/src/pages/notfound.vue", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/src/pages/register._vue")
		outputFile := fmt.Sprintf("../../../../%s/webapp/src/pages/register.vue", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/src/pages/successregister._vue")
		outputFile := fmt.Sprintf("../../../../%s/webapp/src/pages/successregister.vue", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/src/main._js")
		outputFile := fmt.Sprintf("../../../../%s/webapp/src/main.js", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	// {
	// 	templateFile := fmt.Sprintf("../templates/frontend/src/api/index._js")
	// 	outputFile := fmt.Sprintf("../../../../%s/webapp/src/api/index.js", tp.packageName)
	// 	basic(tp, templateFile, outputFile, tp)
	// }

	{
		templateFile := fmt.Sprintf("../templates/frontend/src/store/index._js")
		outputFile := fmt.Sprintf("../../../../%s/webapp/src/store/index.js", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/src/store/basiccrud._js")
		outputFile := fmt.Sprintf("../../../../%s/webapp/src/store/basiccrud.js", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/src/router/index._js")
		outputFile := fmt.Sprintf("../../../../%s/webapp/src/router/index.js", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/src/utils/httprequest._js")
		outputFile := fmt.Sprintf("../../../../%s/webapp/src/utils/httprequest.js", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

	{
		templateFile := fmt.Sprintf("../templates/frontend/src/utils/auth._js")
		outputFile := fmt.Sprintf("../../../../%s/webapp/src/utils/auth.js", tp.packageName)
		basic(tp, templateFile, outputFile, tp)
	}

}

func basic(pkg *ThePackage, templateFile, outputFile string, object interface{}) {

	fmt.Println(templateFile)
	file, err := os.Open(templateFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var buffer bytes.Buffer

	for scanner.Scan() {
		row := scanner.Text()
		buffer.WriteString(row)
		buffer.WriteString("\n")
	}

	FuncMap := template.FuncMap{
		"HasTime":     HasTime,
		"CamelCase":   CamelCase,
		"PascalCase":  PascalCase,
		"SnakeCase":   SnakeCase,
		"UpperCase":   UpperCase,
		"LowerCase":   LowerCase,
		"AppName":     func() string { return pkg.appName },
		"PackageName": func() string { return pkg.packageName },
		// "DataTypeId":  func() string { return pkg.datatypeid },
	}

	t, err := template.
		New("todos").
		Funcs(FuncMap).
		Parse(buffer.String())

	if err != nil {
		panic(err)
	}

	var bf bytes.Buffer
	err = t.Execute(&bf, object)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(outputFile, bf.Bytes(), 0664)
	if err != nil {
		panic(err)
	}

}

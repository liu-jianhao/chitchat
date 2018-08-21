package main

import (
	"chitchat/data"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

var config Configuration
var logger *log.Logger

// 为了方便而实现的一个打印函数
// Convenience function for printing to stdout
func p(a ...interface{}) {
	fmt.Println(a)
}

// 在这个包中（main包），会首先执行init函数
func init() {
	loadConfig()
	file, err := os.OpenFile("chitchat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

// 将配置文件读取到全局变量config中
func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	// 根据给定的JSON文件，创建出相应的解码器
	decoder := json.NewDecoder(file)
	config = Configuration{}
	// 将JSON数据解码至config结构
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}

// Convenience function to redirect to the error message page
func error_message(writer http.ResponseWriter, request *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(writer, request, strings.Join(url, ""), 302)
}

// 检查用户是否登录
// Checks if the user is logged in and has a session, if not err is not nil
func session(writer http.ResponseWriter, request *http.Request) (sess data.Session, err error) {
	// 从请求中取出cookie
	cookie, err := request.Cookie("_cookie")
	if err == nil {
		sess = data.Session{Uuid: cookie.Value}
		// 还要进行第二项检查——访问数据库并核实会话的唯一ID是否存在
		if ok, _ := sess.Check(); !ok {
			err = errors.New("Invalid session")
		}
	}
	return
}

// parse HTML templates
// pass in a list of file names, and get a template
func parseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))
	return
}

// 最后一个参数以3个点开头，表示该函数是一个可变参数函数，可在最后的可变参数中接受零个或任意多个值作为参数
// 注意：可变参数必须是最后一个参数
func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	// ParseFiles函数对这些模板文件进行语法分析，并创建出相应的模板
	// 为了捕捉语法分析过程中过程中可能会产生的错误，程序使用了Must函数去包围ParseFiles函数的执行结果
	// 这样当ParseFiles返回错误时，Must函数会行用户返回相应的错误报告
	// 以layout.html模板文件为例（见template/layout.html）
	// 源代码中使用了define动作。这个动作通过文件开头的{{ define "layout" }}和文件结尾的{{ end }}，
	// 把被包围的文本快定义成layout模板的一部分
	// 跟在引用模板名字之后的点(.)代表了传递给被引用模板的数据
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(writer, "layout", data)
}

// for logging
func info(args ...interface{}) {
	logger.SetPrefix("INFO ")
	logger.Println(args...)
}

func danger(args ...interface{}) {
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}

func warning(args ...interface{}) {
	logger.SetPrefix("WARNING ")
	logger.Println(args...)
}

// version
func version() string {
	return "0.1"
}

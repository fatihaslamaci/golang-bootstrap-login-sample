package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var templatesNavbar = template.Must(template.ParseFiles(
	"./templates/basenavbar.html",
	"./templates/head.html",
	"./templates/basefooter.html"))

var templatesBlank = template.Must(template.ParseFiles(
	"./templates/baseblank.html",
	"./templates/head.html",
	"./templates/basefooter.html"))

type Context struct {
	Title string
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func render(w http.ResponseWriter, tmpl string, js_tmpl string, context Context, template *template.Template) {
	tmpl_list := []string{fmt.Sprintf("templates/%s.html", tmpl),
		fmt.Sprintf("templates/%s.html", js_tmpl)}

	t, err := template.Clone()
	if err != nil {
		log.Print("template clone error: ", err)
	}

	t, err = t.ParseFiles(tmpl_list...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, context)
	if err != nil {
		log.Print("template executing error: ", err)

	}
}

func renderBlank(w http.ResponseWriter, tmpl string, js_tmpl string, context Context) {
	render(w, tmpl, js_tmpl, context, templatesBlank)
}

func renderNavbar(w http.ResponseWriter, tmpl string, js_tmpl string, context Context) {
	render(w, tmpl, js_tmpl, context, templatesNavbar)
}

func blankHandler(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "Blak Page"}
	renderNavbar(w, "blank", "blank_js", context)
}

func loginPageHandler(response http.ResponseWriter, request *http.Request) {
	context := Context{Title: "Login Page"}
	renderBlank(response, "login", "blank_js", context)
}

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	context := Context{Title: "Login Page"}
	renderBlank(response, "login", "blank_js", context)
}

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	context := Context{Title: "Home Page"}
	if userName != "" {
		renderNavbar(response, "main", "blank_js", context)
	} else {
		http.Redirect(response, request, "/", 302)
	}

}

func morrisPageHandler(response http.ResponseWriter, request *http.Request) {
	context := Context{Title: "Morris Page"}
	renderNavbar(response, "morris", "morris_js", context)
}

func flotPageHandler(response http.ResponseWriter, request *http.Request) {
	context := Context{Title: "Flot Page Page"}
	renderNavbar(response, "flot", "flot_js", context)
}

func tablesHandler(response http.ResponseWriter, request *http.Request) {
	context := Context{Title: "tables Page"}
	renderNavbar(response, "tables", "tables_js", context)
}

func formsHandler(response http.ResponseWriter, request *http.Request) {
	context := Context{Title: "forms Page"}
	renderNavbar(response, "forms", "blank_js", context)
}

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["email"]
		}
	}
	return userName
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"email": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

// login handler

func loginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("email")
	pass := request.FormValue("password")
	redirectTarget := "/"
	if name != "" && pass != "" {
		// .. check credentials ..
		setSession(name, response)
		redirectTarget = "/main"
	}
	http.Redirect(response, request, redirectTarget, 302)
}

// logout handler

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

// server main method

var router = mux.NewRouter()

func panels_wellsHandler(response http.ResponseWriter, request *http.Request) {

	context := Context{Title: "panels_wells Page"}
	renderNavbar(response, "panels-wells", "blank_js", context)

}

func buttonsHandler(response http.ResponseWriter, request *http.Request) {

	context := Context{Title: "buttons Page"}
	renderNavbar(response, "buttons", "blank_js", context)

}

func notificationsHandler(response http.ResponseWriter, request *http.Request) {

	context := Context{Title: "notifications Page"}
	renderNavbar(response, "notifications", "blank_js", context)

}

func typographyHandler(response http.ResponseWriter, request *http.Request) {

	context := Context{Title: "typography Page"}
	renderNavbar(response, "typography", "blank_js", context)

}

func iconsHandler(response http.ResponseWriter, request *http.Request) {

	context := Context{Title: "icons Page"}
	renderNavbar(response, "icons", "blank_js", context)

}

func gridHandler(response http.ResponseWriter, request *http.Request) {

	context := Context{Title: "grid Page"}
	renderNavbar(response, "grid", "blank_js", context)

}

func addStaticDir(s string) {
	http.Handle("/"+s+"/", http.StripPrefix("/"+s, http.FileServer(http.Dir("./statics/"+s))))
}

func addStaticDirAll() {
	addStaticDir("bower_components")
	addStaticDir("dist")
	addStaticDir("js")
	addStaticDir("less")
}

func main() {

	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/loginpage", loginPageHandler)
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler)
	router.HandleFunc("/main", internalPageHandler)

	router.HandleFunc("/flot", flotPageHandler)
	router.HandleFunc("/morris", morrisPageHandler)
	router.HandleFunc("/tables", tablesHandler)
	router.HandleFunc("/forms", formsHandler)
	router.HandleFunc("/panels-wells", panels_wellsHandler)
	router.HandleFunc("/buttons", buttonsHandler)
	router.HandleFunc("/notifications", notificationsHandler)
	router.HandleFunc("/typography", typographyHandler)
	router.HandleFunc("/icons", iconsHandler)
	router.HandleFunc("/grid", gridHandler)
	router.HandleFunc("/blank", blankHandler)

	http.Handle("/", router)

	addStaticDirAll()

	http.ListenAndServe(":8000", nil)
}

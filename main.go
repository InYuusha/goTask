package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)
type Note struct{
	Title string
	Desc string 
	CreatedOn time.Time

}
var noteStore = make(map[string]Note)

var id int =0

var templates = make(map[string]*template.Template)

func init(){
templates["index"]=template.Must(template.ParseFiles("public/templates/index.html","public/templates/base.html"))
templates["add"]=template.Must(template.ParseFiles("public/templates/add.html","public/templates/base.html"))
//templates["edit"]=template.Must(template.ParseFiles("templates/edit.html","templates/base.html"))

}
func renderTemplate(res http.ResponseWriter, name string , template string , viewModel interface{}){
	temp,exists:=templates[name]
	if !exists{
		http.Error(res,"The Template Doesnt exists",http.StatusInternalServerError)
	}
	err:=temp.ExecuteTemplate(res,template,viewModel)
	if err!=nil{
		http.Error(res,err.Error(),http.StatusInternalServerError)
	}
}
func getNotes(res http.ResponseWriter ,req*http.Request){
	renderTemplate(res,"index","base",noteStore)

}
func addNote(res http.ResponseWriter,req*http.Request){
	renderTemplate(res,"add","base",nil)
}
func saveNote(res http.ResponseWriter,req *http.Request){
	
	req.ParseForm()
	title:=req.PostFormValue("title")
	desc:=req.PostFormValue("desc")

	note:=Note{title,desc,time.Now()}
	id++
	k:=strconv.Itoa(id)
	noteStore[k]=note
	http.Redirect(res,req,"/",302)
}
func delNote(res http.ResponseWriter,req *http.Request){
	vars:=mux.Vars(req)
	id:=vars["id"]
	if _,ok:=noteStore[id];ok{
		delete(noteStore,id)
	}else{
		http.Error(res,"Doesnt exists",404)
	}
	http.Redirect(res,req,"/",302)
}
func main(){
	r:=mux.NewRouter().StrictSlash(false)
	
	r.HandleFunc("/",getNotes)
	r.HandleFunc("/add",addNote)
	r.HandleFunc("/save",saveNote)
	r.HandleFunc("/del/{id}",delNote)

	server:=&http.Server{
		Addr:":3000",
		Handler: r,
	}
	log.Println("Listening")
	server.ListenAndServe()
}
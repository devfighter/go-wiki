package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/renstrom/go-wiki/vendor/_nuts/github.com/gorilla/mux"
)

const imageTypes = ".jpg .jpeg .png .gif"

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	md, err := ioutil.ReadFile(options.Dir + "/index.md")
	if err != nil {
		log.Fatalln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	wiki := Wiki{Markdown: md, template: options.template}
	wiki.Write(w)
}

func WikiHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Deny requests trying to traverse up the directory structure using
	// relative paths
	if strings.Contains(vars["filepath"], "..") {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Path to the file as it is on the the local file system
	fsPath := fmt.Sprintf("%s/%s", options.Dir, vars["filepath"])

	// Serve (accepted) images
	for _, filext := range strings.Split(imageTypes, " ") {
		if path.Ext(r.URL.Path) == filext {
			http.ServeFile(w, r, fsPath)
			return
		}
	}

	md, err := ioutil.ReadFile(fsPath + ".md")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	wiki := Wiki{
		Markdown: md,
		filepath: fsPath,
		template: options.template,
	}

	wiki.Commits, err = Commits(vars["filepath"]+".md", 5)
	if err != nil {
		log.Println("ERROR", "Failed to get commits")
	}

	wiki.Write(w)
}

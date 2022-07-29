/*  Star Wars Role-Playing Mud
 *  Copyright (C) 2022 @{See Authors}
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */
package swr

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type Editor struct {
}

func EditorStart() {
	go func() {
		fs := http.FileServer(http.Dir("./web/public"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))
		http.HandleFunc("/", serveTemplate)
		http.HandleFunc("/data", dataHandler)
		http.ListenAndServe(":8080", nil)
		log.Println("Editor now accepting connections on 0.0.0.0:8080")
	}()
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	fn := filepath.Clean(r.URL.Path)
	if fn == "/" {
		fn = "/index.html"
	}
	if !strings.HasSuffix(fn, ".html") {
		fn = fmt.Sprintf("%s.html", fn)
	}
	lp := filepath.Join("web", "templates", "layout.html")
	fp := filepath.Join("web", "templates", fn)
	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Print(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		dataGet(w, r)
	case "PUT":
		dataPut(w, r)
	case "POST":
		dataPost(w, r)
	case "DELETE":
		dataDelete(w, r)
	}
}

func writeError(w http.ResponseWriter, err error) {
	if err != nil {
		w.Header().Write(w)
		w.WriteHeader(500)
		w.Write([]byte("---\nstatus: 500\n"))
	}
}

func writeData(w http.ResponseWriter, obj interface{}) {
	buf, err := yaml.Marshal(obj)
	ErrorCheck(err)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Write(w)
	w.WriteHeader(200)
	w.Write(buf)
}
func dataGet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query.Has("type") {
		t := query.Get("type")
		switch t {
		case "area":
			area := DB().areas[query.Get("name")]
			if area != nil {
				writeData(w, area)
				return
			}
		case "room":
			sid := query.Get("id")
			id, err := strconv.Atoi(sid)
			if err != nil {
				writeError(w, err)
				return
			}
			room := DB().rooms[uint(id)]
			if room != nil {
				writeData(w, room)
				return
			}
		case "item":
			sid := query.Get("id")
			id, err := strconv.Atoi(sid)
			if err != nil {
				writeError(w, err)
				return
			}
			item := DB().items[uint(id)]
			if item != nil {
				writeData(w, item)
				return
			}
		case "entity":
			sid := query.Get("id")
			id, err := strconv.Atoi(sid)
			if err != nil {
				writeError(w, err)
				return
			}
			mob := DB().mobs[uint(id)]
			if mob != nil {
				writeData(w, mob)
				return
			}
		case "ship":
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{ \"status\": 404 }"))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{ \"status\": 404 }"))
	}
}

func dataPut(w http.ResponseWriter, r *http.Request) {

}

func dataPost(w http.ResponseWriter, r *http.Request) {

}

func dataDelete(w http.ResponseWriter, r *http.Request) {

}

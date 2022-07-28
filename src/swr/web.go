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
	"strings"
	"text/template"
)

type Editor struct {
}

func EditorStart() {
	go func() {
		fs := http.FileServer(http.Dir("./web/public"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))
		http.HandleFunc("/", serveTemplate)
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
	fmt.Println(fn)
	lp := filepath.Join("web", "templates", "layout.html")
	fp := filepath.Join("web", "templates", fn)
	fmt.Println(fp)
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

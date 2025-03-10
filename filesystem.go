package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

/*
<!doctype html>
<meta name="viewport" content="width=device-width">
<pre>
<a href="ai.txt">ai.txt</a>
<a href="auth.go">auth.go</a>
<a href="cert.pem">cert.pem</a>
<a href="go.mod">go.mod</a>
<a href="go.sum">go.sum</a>
<a href="key.pem">key.pem</a>
<a href="log-share">log-share</a>
<a href="main.go">main.go</a>
<a href="utils.go">utils.go</a>
</pre>



<!doctype html>
<meta name="viewport" content="width=device-width">
<pre>
<a href="drwxr-xr-x  2 danl  staff       64 Mar  7 18:29 folder1">/files/folder1</a>
<a href="-rw-r--r--  1 danl  staff        0 Mar  7 15:00 ai.txt">/files/ai.txt</a>
</pre>

*/

var htmlHeader = "<!doctype html>\n<meta name=\"viewport\" content=\"width=device-width\">\n"
var notFoundError = htmlHeader + "<pre>File Not Found</pre>"

func scanFolder(startdir string) string {
	fileSystem := os.DirFS(startdir)
	html := htmlHeader + "<pre>"

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
		}
		if path == "." {
			return nil
		}
		fi, err := d.Info()
		if err != nil {
			fmt.Println(err)
			return nil
		}

		//fmt.Println(path)
		// spath, ok := strings.CutPrefix(path, directory)
		// if !ok {
		//      log.Println("could not match prefix", path, directory)
		//      return nil
		// }

		fpath := filepath.Join(startdir, path)
		if d.IsDir() && path != "." {
			html += fmt.Sprintf("d%s\t%d\t%s\t<a href=\"/path?file=%s\">%s</a><br>",
				fi.Mode().String(),
				fi.Size(),
				fi.ModTime(),
				fpath,
				d.Name())

			return filepath.SkipDir // don't scan recursively
		} else {

			sizeType := "b"
			sizeMod := fi.Size()

			if sizeMod > 1024*1024*1024 { // GB
				sizeMod = sizeMod / 1024 / 1024 / 1204
				sizeType = "G"
			} else if sizeMod > 1024*1024 { // MB
				sizeMod = sizeMod / 1024 / 1024
				sizeType = "M"
			} else if sizeMod > 1024 { // KB
				sizeMod = sizeMod / 1024
				sizeType = "K"
			}

			html += fmt.Sprintf("-%s\t%d%s\t%s\t<a href=\"/path?file=%s\" download=\"%s\">%s</a><br>",
				fi.Mode().String(),
				sizeMod,
				sizeType,
				fi.ModTime(),
				fpath,
				d.Name(),
				d.Name())

		}
		return nil
	})
	html += "</pre>"
	return html
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(scanFolder(directory)))
}

func servePath(w http.ResponseWriter, r *http.Request) {
	file := r.FormValue("file")
	if file == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(notFoundError))
		return
	}

	f, err := os.Stat(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%s<pre>%s</pre>", htmlHeader, err)))
		return
	}
	//fmt.Println(file, f.Name())

	if f.IsDir() {
		w.Write([]byte(scanFolder(file)))
		return
	}

	fh, err := os.Open(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%s<pre>%s</pre>", htmlHeader, err)))
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, fh)
}

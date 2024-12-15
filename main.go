package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/phillip-england/vbf"
)

const KeyTemplates = "KEYTEMPLATES"

func main() {

	mux, gCtx := vbf.VeryBestFramework()

	strEquals := func(input string, value string) bool {
		return input == value
	}

	funcMap := template.FuncMap{
		"strEquals": strEquals,
	}

	templates, err := vbf.ParseTemplates("./templates", funcMap)
	if err != nil {
		panic(err)
	}

	vbf.SetGlobalContext(gCtx, KeyTemplates, templates)
	vbf.HandleStaticFiles(mux)
	vbf.HandleFavicon(mux)

	vbf.AddRoute("GET /", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			templates, _ := vbf.GetContext(KeyTemplates, r).(*template.Template)
			mdContent, err := vbf.LoadMarkdown("./content/index.md", "dracula")
			if err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(500)
				return
			}
			vbf.ExecuteTemplate(w, templates, "root.html", map[string]interface{}{
				"Title":   "phillip england",
				"Content": template.HTML(mdContent),
				"ReqPath": r.URL.Path,
        "PostName": "[phillip england]",
        "SubText": "learning in public",
        "ImagePath": "/static/img/me2.jpg",
        "DateWritten": "🎂 born on 12/9/2024",
			})
		} else {
			vbf.WriteString(w, "404 not found")
		}
	}, vbf.MwLogger)

	vbf.AddRoute("GET /posts", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		templates, _ := vbf.GetContext(KeyTemplates, r).(*template.Template)
		mdContent, err := vbf.LoadMarkdown("./content/posts.md", "dracula")
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(500)
			return
		}
		vbf.ExecuteTemplate(w, templates, "root.html", map[string]interface{}{
			"Title":   "phillip england",
			"Content": template.HTML(mdContent),
			"ReqPath": r.URL.Path,
      "PostName": "[posts]",
      "SubText": "see what i've been learning about",
      "ImagePath": "/static/img/posts.webp",
      "DateWritten": "🎂 born on 12/9/2024",
		})
	}, vbf.MwLogger)

  filepath.Walk("./content", func(path string, info fs.FileInfo, err error) error {
    if !strings.Contains(path, "/post/") {
      return nil
    }
    parts := strings.Split(path, "/")
    if len(parts) == 0 {
      return nil
    }
    lastPart := parts[len(parts)-1]
    moreParts := strings.Split(lastPart, ".")
    if len(moreParts) < 3 {
      return nil
    }
    index := moreParts[0]
    title := moreParts[1]
    extension := moreParts[2]
    if extension != "md" {
      return nil
    }
    imagePath := "/static/img/post/"+index+".webp"
    appPath := "/post/"+title
    vbf.AddRoute("GET "+appPath, mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
      templates, _ := vbf.GetContext(KeyTemplates, r).(*template.Template)
      mdContent, err := vbf.LoadMarkdown("/"+path, "dracula")
      if err != nil {
        fmt.Println(err.Error())
        w.WriteHeader(500)
        return
      }
      doc, err := goquery.NewDocumentFromReader(strings.NewReader(mdContent))
      if err != nil {
        fmt.Println(err.Error())
        w.WriteHeader(500)
        return
      }
      metaDataElm := doc.Find("#meta-data")
      var subText string
      var dob string
      metaDataElm.Find("*").Each(func(i int, sel *goquery.Selection) {
        key, _ := sel.Attr("key")
        if key == "subtext" {
          val, _ := sel.Attr("value")          
          subText = val
        }
        if key == "dob" {
          val, _ := sel.Attr("value")
          dob = val
        }
      })
      vbf.ExecuteTemplate(w, templates, "root.html", map[string]interface{}{
        "Title":   "phillip england",
        "Content": template.HTML(mdContent),
        "ReqPath": r.URL.Path,
        "PostName": "["+ strings.ReplaceAll(title, "-", " ") +"]",
        "SubText": subText,
        "ImagePath": imagePath,
        "DateWritten": "written " +dob,
      })
    }, vbf.MwLogger)
    return nil
  })

	err = vbf.Serve(mux, "8080")
	if err != nil {
		panic(err)
	}

}

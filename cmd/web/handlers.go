package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ga676005/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, 200, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// 前端傳的 formData 要先用 r.ParseForm() 讀到 r.PostForm 裡
	// 建議不要用直接用 r.PostFormValue()，因為它會忽略錯誤
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// 用 r.PostForm.Get() 去抓欄位
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// r.PostForm.Get() 回傳的永遠都是 string，而且只能抓第一個值
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// 1. 如果是 checkbox 有多個值的時候
	// <input type="checkbox" name="items" value="foo"> Foo
	// <input type="checkbox" name="items" value="bar"> Bar
	// <input type="checkbox" name="items" value="baz"> Baz
	// 要用這種方式寫
	// for i, item := range r.PostForm["items"] {
	// 	fmt.Fprintf(w, "%d: Item %s\n", i, item)
	// }

	// 2. 限制 form size
	// POST 預設是 10 MB，除非有在 <form> 加上 enctype="multipart/form-data"
	// 然後傳的東西也是 multipart data，才能超過上限
	// 另一個設上限的方式 r.Body = http.MaxBytesReader(w, r.Body, 4096)，要寫在 r.ParseForm() 前面
	// 如果超過上限 MaxBytesReader 會在 http.ResponseWriter 設個 flag 來關閉 TCP 連線

	// 3. query string parameters
	// 如果 form 的 method 是 get，它會這樣送 /foo/bar?title=value&content=value
	// 那就用 r.URL.Query().Get("title") 去抓，這一樣永遠是字串

	// 4. r.Form
	// r.Form.Get() 會抓 r.body 和 query string 的欄位
	// 如果兩個都存在的話會優先用 r.body 的

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

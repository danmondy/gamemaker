package site

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/danmondy/gamemaker/data"
)

var (
	templates *template.Template
	cards     []data.Card
	designs   []data.Design
	deckFile  string
)

func InitSite(cfg data.Config) {
	FuncMap := BuildFuncMap()
	deckFile = cfg.FileLocation
	LoadDeck(cfg.FileLocation)
	fmt.Println("Docroot:", cfg.Docroot)
	templates = template.Must(template.New("mass").Funcs(FuncMap).ParseGlob(fmt.Sprintf("%s/templates/*", cfg.Docroot)))
}

func BuildFuncMap() template.FuncMap {
	return template.FuncMap{
		"Loop": func(n int) []int {
			r := make([]int, n)
			for i := 0; i < n; i++ {
				r[i] = i
			}
			return r
		},
		"Add":         func(a int, b int) int { return a + b },
		"PrettyYear":  func(t time.Time) string { return t.Format("2006") },
		"PrettyMonth": func(m time.Time) string { return m.Month().String()[0:3] + "." },
		"Elipses":     func(s string) string { return fmt.Sprintf("%s...", []byte(s)[0:3]) },
	}
}

func SetRoutes(r *mux.Router) {

	fs := http.FileServer(http.Dir(deckFile))
	r.PathPrefix("/deckImg/").Handler(http.StripPrefix("/deckImg/", fs))
	r.HandleFunc("/", IndexHandler)

	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/logout", LogOutHandler)
	r.HandleFunc("/signup", NewUserHandler)

	//r.HandleFunc("/places", Auth(PlacesHandler))
	//r.HandleFunc("/places/new", Auth(NewPlaceHandler))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	m := ViewModel{Page: "index"}

	m.Model = struct {
		Designs []data.Design
		Cards   []data.Card
	}{Designs: designs, Cards: cards}

	err := renderTemplate(w, "index", m)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func MaintenanceHandler(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, fmt.Sprintf("%s/templates/maintenance.html", cfg.Docroot))
	err := renderTemplate(w, "maintenance", nil)
	if err != nil {
		fmt.Println(err)
	}
}

type UserHandlerFunc func(http.ResponseWriter, *http.Request, *data.User)

func Admin(hf UserHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func Auth(hf UserHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, ":", r.URL.Host, r.URL.Path)
		cookie, err := r.Cookie("authtoken")
		if err != nil {
			fmt.Println("Cookie not found")
			http.Redirect(w, r, "/login?msg=You are not logged in.", http.StatusTemporaryRedirect)
			return
		}
		sessId := fmt.Sprintf("%s-%s", cookie.Value, r.RemoteAddr)
		session, ok := data.SessionCache.Get(sessId)
		if ok {
			fmt.Println("User is logged in and in an active session.")
			hf(w, r, session.User)
			return
		}
		fmt.Println("session not found.")
		http.Redirect(w, r, fmt.Sprintf("/login?msg=You are not logged in.&url=%s", r.URL), http.StatusTemporaryRedirect)
		return
	}
}

/*func PlacesHandler(w http.ResponseWriter, r *http.Request, u *data.User) {
	places, err := data.GetUserPlaces(u.Id)
	msgs := []string{}
	if err != nil {
		msgs = append(msgs, err.Error())
	}

	model := struct {
		Places []data.Place
		User   data.User
		Msg    []string
	}{places, *u, msgs}

	err = renderTemplate(w, "places", model)
	if err != nil {
		fmt.Println(err)
	}
}
func NewPlaceHandler(w http.ResponseWriter, r *http.Request, u *data.User) {
	switch r.Method {
	case "GET":
		page := struct {
			User data.Place
			Msg  []string
		}{}
		renderTemplate(w, "newplace", page)
	case "POST":
		if err := r.ParseForm(); err != nil {
			renderTemplate(w, "/places/new", err.Error())
			return
		}
		title := r.FormValue("title")
		subtitle := r.FormValue("subtitle")
		description := r.FormValue("description")

		place := data.NewPlace(u)
		place.Title = title
		place.Subtitle = subtitle
		place.Description = description

		err := place.Insert()
		if err != nil {
			renderTemplate(w, "newplace", user)
			fmt.Println(err)
			return
		}
		http.Redirect(w, r, "/places", http.StatusFound)
	}
}*/
func NewUserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		page := struct {
			User data.User
			Msg  []string
		}{}
		renderTemplate(w, "newuser", page)
	case "POST":
		if err := r.ParseForm(); err != nil {
			renderTemplate(w, "newuser", err.Error())
			return
		}
		email := r.FormValue("email")
		pword := r.FormValue("password")

		user := data.NewUser(email, pword)
		err := user.Insert()
		if err != nil {
			renderTemplate(w, "newuser", user)
			fmt.Println(err)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authtoken")
	if err != nil {
		fmt.Println("Cookie not found")
		http.Redirect(w, r, "/login?msg=You are not logged in.", http.StatusTemporaryRedirect)
		return
	}
	sessId := fmt.Sprintf("%s-%s", cookie.Value, r.RemoteAddr)
	cookie.Expires = time.Now()
	data.SessionCache.Remove(sessId)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		values := r.URL.Query()
		msg := values.Get("msg")
		url := values.Get("url")
		model := struct{ Message, Url string }{msg, url}
		renderTemplate(w, "login", model)
		return
	case "POST":
		//READ FORM
		if err := r.ParseForm(); err != nil {
			renderTemplate(w, "login", err.Error())
			return
		}
		email := r.FormValue("email")
		pword := r.FormValue("password")
		url := r.FormValue("url")
		user, ok := data.Authenticate(email, pword)
		if ok { //IF AUTH OK, CREATE SESSION
			fmt.Println("User auth ok")
			cookie := &http.Cookie{}
			cookie.Name = "authtoken"
			authToken := time.Now().Format(time.RFC822) //TODO: add random characters to end
			cookie.Value = authToken                    // needs mutex
			cookie.Expires = time.Now().Add(time.Hour)
			http.SetCookie(w, cookie)
			session := data.Session{Id: fmt.Sprintf("%s-%s", authToken, r.RemoteAddr), User: user, LastTouched: time.Now()}
			data.SessionCache.Add(session)
			if url != "" {
				http.Redirect(w, r, url, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "/maps", http.StatusSeeOther)
			return
		} else {
			fmt.Println("User auth not so ok")
			renderTemplate(w, "login", LoginView{"No Username/Password combination was found.", url})
			return
		}
	}
	renderErr(http.StatusMethodNotAllowed)
}

func renderTemplate(w http.ResponseWriter, tmpl string, model interface{}) error {
	err := templates.ExecuteTemplate(w, tmpl+".html", model)
	return err
}

func renderErr(status int) {

}

type LoginView struct {
	Message, Url string
}

type ViewModel struct {
	Page  string
	User  data.User
	Msg   []string
	Model interface{}
	Alert string
}

func LoadDeck(file string) {
	dat, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", file, "cards.json"))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = json.Unmarshal(dat, &cards)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	dat, err = ioutil.ReadFile(fmt.Sprintf("%s/%s", file, "designs.json"))
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(dat, &designs)
	if err != nil {
		panic(err)
	}
}

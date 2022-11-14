package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	// "golang.org/x/text/message"
)

func main() {
	connection.DatabaseConnect()
	route := mux.NewRouter()

	route.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")

	route.HandleFunc("/form-add-project", formAddProject).Methods("GET")
	route.HandleFunc("/send-data-add-project", sendDataAddProject).Methods("POST")
	route.HandleFunc("/form-edit-project/{id}", formEditProject).Methods("GET")
	route.HandleFunc("/send-data-edit-project/{id}", sendDataEditProject).Methods("POST")
	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")

	route.HandleFunc("/contact", contact).Methods("GET")

	route.HandleFunc("/form-register", formRegister).Methods("GET")
	route.HandleFunc("/register", register).Methods("POST")

	route.HandleFunc("/form-login", formLogin).Methods("GET")
	route.HandleFunc("/login", login).Methods("POST")

	route.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println("Server running on localhost:8000")
	http.ListenAndServe("localhost:8000", route)
}

type SessionData struct {
	UserId    int
	IsLogin   bool
	UserName  string
	FlashData string
}

type projectDataStruc struct {
	Id              int
	ProjectName     string
	StartDate       time.Time
	EndDate         time.Time
	StartDateFormat string
	EndDateFormat   string
	Duration        string
	Description     string
	Technologies    []string
	Image           string
	IsLogin         bool
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

var Data = SessionData{}

var projectData = []projectDataStruc{
	// {
	// 	Id:           0,
	// 	ProjectName:  "Dummy Project 1",
	// 	StartDate:    "2022-09-12",
	// 	EndDate:      "2022-09-19",
	// 	Duration:     "1 Weeks",
	// 	Description:  "Description Dummy Project 1",
	// 	Technologies: []string{"NodeJs", "ReactJs"},
	// 	Image:        "barelang.jpeg",
	// },
	// {
	// 	Id:           1,
	// 	ProjectName:  "Dummy Project 2",
	// 	StartDate:    "2022-09-20",
	// 	EndDate:      "2022-09-25",
	// 	Duration:     "5 Days",
	// 	Description:  "Description Dummy Project 2",
	// 	Technologies: []string{"NodeJs", "ReactJs", "Golang"},
	// 	Image:        "cendrawasih.jpeg",
	// },
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		getFlashMessage := session.Flashes("Message")
		session.Save(r, w)
		var buildFlashMessage []string
		if len(getFlashMessage) > 0 {
			for _, fMLetter := range getFlashMessage {
				buildFlashMessage = append(buildFlashMessage, fMLetter.(string))
			}
		}

		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserId = session.Values["Id"].(int)
		Data.UserName = session.Values["Name"].(string)
		Data.FlashData = strings.Join(buildFlashMessage, "")
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	} else {

		var result []projectDataStruc
		data, _ := connection.Conn.Query(context.Background(), "SELECT id, project_name, start_date, end_date, duration, description, technologies, image FROM db_project")

		for data.Next() {
			var each = projectDataStruc{}
			err := data.Scan(&each.Id, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Duration, &each.Description, &each.Technologies, &each.Image)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			result = append(result, each)
		}

		response := map[string]interface{}{
			"DataSession": Data,
			"ProjectData": result,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/detail-project.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	} else {
		where_id, _ := strconv.Atoi(mux.Vars(r)["id"])

		selectedProjectData := projectDataStruc{}

		err = connection.Conn.QueryRow(context.Background(), "SELECT id, project_name, start_date, end_date, duration, description, technologies, image FROM db_project WHERE id=$1", where_id).
			Scan(&selectedProjectData.Id, &selectedProjectData.ProjectName, &selectedProjectData.StartDate, &selectedProjectData.EndDate, &selectedProjectData.Duration, &selectedProjectData.Description, &selectedProjectData.Technologies, &selectedProjectData.Image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		selectedProjectData.StartDateFormat = selectedProjectData.StartDate.Format("2006-01-02")
		selectedProjectData.EndDateFormat = selectedProjectData.EndDate.Format("2006-01-02")

		response := map[string]interface{}{
			"selectedProjectData": selectedProjectData,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func formAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/add-project.html")

	if tmpl == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, nil)
	}
}

func sendDataAddProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		projectName := r.PostForm.Get("project-name")
		startDate := r.PostForm.Get("start-date")
		endDate := r.PostForm.Get("end-date")
		var duration string
		description := r.PostForm.Get("description")
		var technologies []string
		technologies = r.Form["technologies"]
		image := r.PostForm.Get("project-image")

		layoutDate := "2006-01-02"
		startDateParse, _ := time.Parse(layoutDate, startDate)
		endDateParse, _ := time.Parse(layoutDate, endDate)

		hour := 1
		day := hour * 24
		week := hour * 24 * 7
		month := hour * 24 * 30
		year := hour * 24 * 365

		differHour := endDateParse.Sub(startDateParse).Hours()
		var differHours int = int(differHour)
		// fmt.Println(differHours)
		days := differHours / day
		weeks := differHours / week
		months := differHours / month
		years := differHours / year

		if differHours < week {
			duration = strconv.Itoa(int(days)) + " Days"
		} else if differHours < month {
			duration = strconv.Itoa(int(weeks)) + " Weeks"
		} else if differHours < year {
			duration = strconv.Itoa(int(months)) + " Months"
		} else if differHours > year {
			duration = strconv.Itoa(int(years)) + " Years"
		}

		_, err = connection.Conn.Exec(context.Background(), "INSERT INTO db_project(project_name, start_date, end_date, duration, description, technologies, image) VALUES ($1, $2, $3, $4, $5, $6, $7)", projectName, startDate, endDate, duration, description, technologies, image)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
			return
		}

		// projectData = append(projectData)

		// fmt.Println(projectData)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func formEditProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/edit-project.html")

	if tmpl == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
	} else {
		where_id, _ := strconv.Atoi(mux.Vars(r)["id"])

		selectedProjectData := projectDataStruc{}

		err = connection.Conn.QueryRow(context.Background(), "SELECT id, project_name, start_date, end_date, description, technologies, image FROM db_project WHERE id=$1", where_id).
			Scan(&selectedProjectData.Id, &selectedProjectData.ProjectName, &selectedProjectData.StartDate, &selectedProjectData.EndDate, &selectedProjectData.Description, &selectedProjectData.Technologies, &selectedProjectData.Image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		selectedProjectData.StartDateFormat = selectedProjectData.StartDate.Format("2006-01-02")
		selectedProjectData.EndDateFormat = selectedProjectData.EndDate.Format("2006-01-02")

		response := map[string]interface{}{
			"selectedProjectData": selectedProjectData,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
}

func sendDataEditProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		where_id, _ := strconv.Atoi(mux.Vars(r)["id"])

		projectName := r.PostForm.Get("project-name")
		startDate := r.PostForm.Get("start-date")
		endDate := r.PostForm.Get("end-date")
		var duration string
		description := r.PostForm.Get("description")
		var technologies []string
		technologies = r.Form["technologies"]
		image := r.PostForm.Get("project-image")

		layoutDate := "2006-01-02"
		startDateParse, _ := time.Parse(layoutDate, startDate)
		endDateParse, _ := time.Parse(layoutDate, endDate)

		hour := 1
		day := hour * 24
		week := hour * 24 * 7
		month := hour * 24 * 30
		year := hour * 24 * 365

		differHour := endDateParse.Sub(startDateParse).Hours()
		var differHours int = int(differHour)
		// fmt.Println(differHours)
		days := differHours / day
		weeks := differHours / week
		months := differHours / month
		years := differHours / year

		if differHours < week {
			duration = strconv.Itoa(int(days)) + " Days"
		} else if differHours < month {
			duration = strconv.Itoa(int(weeks)) + " Weeks"
		} else if differHours < year {
			duration = strconv.Itoa(int(months)) + " Months"
		} else if differHours > year {
			duration = strconv.Itoa(int(years)) + " Years"
		}

		_, err = connection.Conn.Exec(context.Background(), "UPDATE db_project SET project_name=$1, start_date=$2, end_date=$3, duration=$4, description=$5, technologies=$6, image=$7 WHERE id=$8",
			projectName, startDate, endDate, duration, description, technologies, image, where_id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM db_project WHERE id=$1", id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func formRegister(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/form-register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("project-name")
	var email = r.PostForm.Get("project-email")
	var password = r.PostForm.Get("project-password")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	// fmt.Println(passwordHash)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
}

func formLogin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/form-login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	getFlashMessage := session.Flashes("Message")
	session.Save(r, w) //untuk mereset flash message dari session browser
	var buildFlashMessage []string
	if len(getFlashMessage) > 0 {
		for _, fMLetter := range getFlashMessage {
			buildFlashMessage = append(buildFlashMessage, fMLetter.(string))
		}
	}
	Data.FlashData = strings.Join(buildFlashMessage, "")

	response := map[string]interface{}{
		"DataSession": Data,
	}

	tmpl.Execute(w, response)
}

func login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	var email = r.PostForm.Get("project-email")
	var password = r.PostForm.Get("project-password")

	user := User{}
	//mengambil data email dan melakukan pengecekan email
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM  tb_user WHERE email=$1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		session.AddFlash("Email belum terdaftar!", "Message")
		session.Save(r, w)

		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
		return
	}

	fmt.Println(user)
	//melakukan pengecekan password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		session.AddFlash("Password salah!", "Message")
		session.Save(r, w)

		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	//menyimpan data ke dalam session
	session.Values["Name"] = user.Name
	session.Values["Email"] = user.Email
	session.Values["Id"] = user.ID
	session.Values["IsLogin"] = true
	session.Options.MaxAge = 10800 //3jam

	session.AddFlash("succesfull login", "Message")

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

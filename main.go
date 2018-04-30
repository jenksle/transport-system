package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/julienschmidt/httprouter"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

//Config file structure
type Config struct {
	Server         string `json:"server"`
	User           string `json:"user"`
	Pwd            string `json:"pwd"`
	Db             string `json:"db"`
	Port           int    `json:"port"`
	SendgridAPIkey string `json:"sendgridAPIkey"`
	EmailName      string `json:"emailName"`
	EmailAddress   string `json:"emailAddress"`
}

type report struct {
	TransportID    int
	ActionFlag     string
	JobNo          int
	Customer       string
	Department     string
	JobDescription string
	TransportDate  string
	DateDiff       int
}

type data struct {
	ReportsCollect       []report
	ReportsDeliver       []report
	ReportsCollectFuture []report
	ReportsDeliverFuture []report
	ReportsAll           []report
	ServerTime           string
}

type jobdetails struct {
	JobNo          string `json:"jobNo"`
	Department     string `json:"department"`
	CustomerName   string `json:"customerName"`
	JobDescription string `json:"jobDescription"`
}

type validation struct {
	DisplayMessage string
	SubmitMsg      string
	JobDescription string
	JobNo          string
	Department     string
	Customer       string
	TransportDate  string
	ActionFlag     string
}

type jobEdit struct {
	DisplayMessage string
	TransportID    string
	ActionFlag     string
	JobNo          string
	Customer       string
	Department     string
	JobDescription string
	TransportDate  string
}

var config Config
var db *sql.DB

func init() {

	// Load application configuration from settings file
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the database and test connection

	connection := fmt.Sprintf("Server=%s;User ID=%s;Password=%s;database=%s;",
		config.Server,
		config.User,
		config.Pwd,
		config.Db)

	db, err = sql.Open("mssql", connection)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

}

func main() {

	router := httprouter.New()

	router.ServeFiles("/assets/*filepath", http.Dir("./assets"))

	router.GET("/transport/create", transportCreate)
	router.POST("/transport/create", transportCreate)
	router.GET("/transport/create/:jobno", createGetJobDetails)
	router.GET("/transport/complete", completeList)
	router.GET("/transport/complete/:transportid", completeJob)
	router.GET("/transport/edit/:transportid", editJob)
	router.POST("/transport/edit", editJob)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), router))

}

func transportCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var output validation

	t, _ := template.ParseFiles("assets/templates/create.tpl")

	output.ActionFlag = "C"

	if r.Method == "POST" {

		output.ActionFlag = r.FormValue("action_flag")
		output.Department = r.FormValue("department") //electrical, mechanical, electronics
		output.JobNo = r.FormValue("job_no")
		output.Customer = r.FormValue("customer")
		output.JobDescription = r.FormValue("job_description")
		output.TransportDate = r.FormValue("date")
		const shortForm = "2006-01-02"

		if strings.TrimSpace(output.Customer) == "" {
			output.DisplayMessage = "No customer entered"
		} else if strings.TrimSpace(output.Department) == "" {
			output.DisplayMessage = "No Department entered"
		} else if strings.TrimSpace(output.TransportDate) == "" {
			output.DisplayMessage = "No date entered"
		} else {

			sql := `INSERT INTO transport (action_flag, department, job_no, customer, job_description, transport_date, is_active)
					VALUES (?, ?, ?, ?, ?, ?, 'Y')`

			d, _ := time.Parse(shortForm, output.TransportDate)

			_, err := db.Exec(sql, output.ActionFlag, output.Department, output.JobNo, output.Customer, output.JobDescription, d)
			if err != nil {
				log.Fatal(err)
			}

			output.DisplayMessage = "Record submitted"

			var email bytes.Buffer

			t, _ := template.ParseFiles("assets/templates/email.tpl")

			err = t.Execute(&email, output)
			if err != nil {
				log.Fatal(err)
			}

			m := mail.NewV3Mail()
			m.SetFrom(mail.NewEmail("Rewinds & J Windsor Ltd", "donotreply@rjweng.com"))
			m.Subject = fmt.Sprintf("New transportation job")

			p := mail.NewPersonalization()
			tos := []*mail.Email{
				mail.NewEmail(config.EmailName, config.EmailAddress),
			}
			p.AddTos(tos...)

			m.AddPersonalizations(p)
			m.AddContent(mail.NewContent("text/html", email.String()))

			request := sendgrid.GetRequest(config.SendgridAPIkey, "/v3/mail/send", "https://api.sendgrid.com")
			request.Method = "POST"
			request.Body = mail.GetRequestBody(m)

			response, err := sendgrid.API(request)
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println(response.StatusCode)
				fmt.Println(response.Body)
				fmt.Println(response.Headers)
			}

		}
	}

	err := t.Execute(w, output)
	if err != nil {
		log.Fatal(err)
	}

}

func createGetJobDetails(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	JobNo := p.ByName("jobno")

	if len(JobNo) > 10 {
		fmt.Println("Job number greater than 10 characters")
	} else if len(JobNo) < 5 {
		fmt.Println("Job number less than 5 characters")
	} else {

		var output jobdetails

		sql1 := `SELECT sc.[name], LEFT(j.contract, 2)
					AS department, LTRIM(RTRIM(j.description1)) AS job_description
				 FROM rjw.scheme.jcmastm j
					JOIN rjw.scheme.slcustm sc ON j.customer = sc.customer
				 WHERE j.job_code = ?`

		err := db.QueryRow(sql1, JobNo).Scan(&output.CustomerName, &output.Department, &output.JobDescription)
		if err != nil && err != sql.ErrNoRows {
			log.Fatal(err)
		}

		b, _ := json.Marshal(output)
		fmt.Fprintf(w, string(b))
	}

}

func completeList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var output data

	sql1 := `SELECT transport_id, action_flag, job_no, customer, department, job_description, COALESCE(CONVERT(NVARCHAR(11), transport_date, 106), '-')
			 FROM transport
			 WHERE is_active = 'Y'
			 ORDER BY transport_date DESC`

	rows, err := db.Query(sql1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	output.ReportsAll = make([]report, 0)

	for rows.Next() {

		var r report

		err := rows.Scan(&r.TransportID, &r.ActionFlag, &r.JobNo, &r.Customer, &r.Department, &r.JobDescription, &r.TransportDate)
		if err != nil {
			log.Fatal(err)
		}

		if r.ActionFlag == "C" {
			r.ActionFlag = "Collect"
		} else if r.ActionFlag == "D" {
			r.ActionFlag = "Deliver"
		} else {
			r.ActionFlag = "Future"
		}

		output.ReportsAll = append(output.ReportsAll, r)

	}

	t, err := template.ParseFiles("assets/templates/complete.tpl")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, output)
	if err != nil {
		log.Fatal(err)
	}
}

func completeJob(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	TransportID := p.ByName("transportid")

	sql2 := `UPDATE transport
				SET is_active = 'N'
			 WHERE transport_id = ?`

	_, err := db.Exec(sql2, TransportID)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/transport/complete", 303)

}

func editJob(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	TransportID := p.ByName("transportid")

	var output jobEdit

	sql1 := `SELECT transport_id, action_flag, job_no, customer, department, job_description, CONVERT(NVARCHAR(10), transport_date, 120)
					FROM transport
				WHERE transport_id = ?`

	err := db.QueryRow(sql1, TransportID).Scan(&output.TransportID, &output.ActionFlag, &output.JobNo, &output.Customer, &output.Department, &output.JobDescription, &output.TransportDate)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	if r.Method == "POST" {

		output.TransportID = r.FormValue("transport_id")
		output.ActionFlag = r.FormValue("action_flag")
		output.Department = r.FormValue("department")
		output.JobNo = r.FormValue("job_no")
		output.Customer = r.FormValue("customer")
		output.JobDescription = r.FormValue("job_description")
		output.TransportDate = r.FormValue("date")

		const shortForm = "2006-01-02"

		if strings.TrimSpace(output.Customer) == "" {
			output.DisplayMessage = "No customer entered"
		} else if strings.TrimSpace(output.Department) == "" {
			output.DisplayMessage = "No Department entered"
		} else if strings.TrimSpace(output.TransportDate) == "" {
			output.DisplayMessage = "No date entered"
		} else {

			sql2 := `UPDATE transport
				 	SET action_flag = ?, department = ?, job_no = ?, customer = ?, job_description = ?, transport_date = ?
				 WHERE transport_id = ?`

			d, _ := time.Parse(shortForm, output.TransportDate)

			_, err := db.Exec(sql2, output.ActionFlag, output.Department, output.JobNo, output.Customer, output.JobDescription, d, output.TransportID)
			if err != nil && err != sql.ErrNoRows {
				log.Fatal(err)
			}

			http.Redirect(w, r, "/transport/complete", 303)
		}

	}

	t, err := template.ParseFiles("assets/templates/edit.tpl")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, output)
	if err != nil {
		log.Fatal(err)
	}

}

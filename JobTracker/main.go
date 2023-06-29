package main

/*------------------------------------------------------------
future:
- support for multi file upload
- (apply) is mainly and only seen for user
- picture of company showed in read list from creation input
- once in the list, pop up upload document for user to apply
- document will be sent to the creator of that specific job id
- server side calculation biar no tipuan, jgn js trs utk validasi
- google map like gojek destinator intregation
- price calculation, jwt auth, pagination, customer support gpt, task
- sediakan template kualifikasi seperti jenis kelamin dan umur
- status of all aplication list
- placeholder text di tiap page seperti "tempat job list",
lalu hilangkan placeholder bila terdapat minim 1 data
- search profile of job seeker and creator by name & guuid
-------------------------------------------------------------*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/jobtracker_db"))
	if err != nil {
		panic(err)
	}

	database.AutoMigrate(&Job{})
	database.AutoMigrate(&Application{})

	DB = database
}

// ------------------ Job REST API -----------------------------------------------
type Job struct {
	Id          int64  `gorm:"primaryKey" json:"id"`
	Nama        string `gorm:"type:varchar(300)" json:"nama"`
	Guid        string `gorm:"type:varchar(300)" json:"guid"`
	Organisasi  string `gorm:"type:varchar(300)" json:"organisasi"`
	Lokasi      string `gorm:"type:varchar(300)" json:"lokasi"`
	Deskripsi   string `gorm:"type:text" json:"deskripsi"`
	Kualifikasi string `gorm:"type:text" json:"kualifikasi"`
}
type Application struct {
	Id            int64  `gorm:"primaryKey" json:"id"`
	Nama          string `gorm:"type:varchar(300)" json:"nama"`
	Kontak        string `gorm:"type:varchar(300)" json:"kontak"`
	Umur          string `gorm:"type:varchar(300)" json:"umur"`
	Jekel         string `gorm:"type:varchar(300)" json:"jekel"`
	Deskripsi     string `gorm:"type:text" json:"deskripsi"`
	JobId         int64  `gorm:"type:bigint(20)" json:"jobid"`
	JobOrganisasi string `gorm:"type:varchar(300)" json:"joborganisasi"`
	Filepath      string `gorm:"type:varchar(300)" json:"filepath"`
	Approval      string `gorm:"type:varchar(300)" json:"approval"`
}

// Progress is used to track the progress of a file upload.
// It implements the io.Writer interface so it can be passed
// to an io.TeeReader()
type Progress struct {
	TotalSize int64
	BytesRead int64
}

func API_IndexJob(c *gin.Context) {

	//var products []Product
	var jobs []Job

	res := DB.Find(&jobs)
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": res.Error})
		return
	}
	log.Print(res.RowsAffected)
	c.JSON(http.StatusOK, jobs)

}

func API_ShowJob(c *gin.Context) {
	var job Job
	id := c.Param("id")

	if err := DB.First(&job, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, job)
}

func API_CreateJob(c *gin.Context) {

	var job Job

	if err := c.ShouldBindJSON(&job); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	DB.Create(&job)
	c.JSON(http.StatusOK, job)
}

func API_UpdateJob(c *gin.Context) {
	var job Job
	id := c.Param("id")

	if err := c.ShouldBindJSON(&job); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if DB.Model(&job).Where("id = ?", id).Updates(&job).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "tidak dapat mengupdate product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui"})

}

func API_DeleteJob(c *gin.Context) {

	var job Job
	id := c.Param("id")

	if DB.Delete(&job, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat menghapus product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}

// ------------------ Application REST API -----------------------------------------------

func API_IndexApplication(c *gin.Context) {
	var applications []Application

	res := DB.Find(&applications)
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": res.Error})
		return
	}
	log.Print(res.RowsAffected)
	c.JSON(http.StatusOK, applications)
}

func API_ShowApplication(c *gin.Context) {
	var application Application
	id := c.Param("id")

	if err := DB.First(&application, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, application)
}

func API_CreateApplication(c *gin.Context) {

	var application Application

	if err := c.ShouldBindJSON(&application); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	DB.Create(&application)
	c.JSON(http.StatusOK, application)
}

func API_UpdateApplication(c *gin.Context) {
	var application Application
	id := c.Param("id")

	if err := c.ShouldBindJSON(&application); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if DB.Model(&application).Where("id = ?", id).Updates(&application).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "tidak dapat mengupdate product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui"})

}

func API_DeleteApplication(c *gin.Context) {

	var application Application
	var app Application
	id := c.Param("id")

	if DB.Model(&application).Where("id = ?", id).Find(&app).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "tidak dapat ketemu application"})
		return
	}

	//Removing file from the directory
	//Using Remove() function not removeAll(for folder&files)
	path := fmt.Sprintf("./static%s", app.Filepath)
	e := os.Remove(path)
	if e != nil {
		log.Fatal(e)
	}

	if DB.Delete(&application, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat menghapus product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}

// ------------------ Job Consume API -----------------------------------------------

// Reference the REST API server
var BASE_URL = "http://localhost:3333"

func redirectToPath(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusSeeOther)
		fmt.Println("Hello: " + r.Host)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/home.html")
}

func ApplicationIndex(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, "views/application_job.html")
	var application Application
	var data map[string]interface{}
	//getting the id from the url when pressing one of the api data list
	id := r.URL.Query().Get("id")
	// this if is for personal contact to application i think idk, not fixed
	if id != "" {
		res, err := http.Get(BASE_URL + "/application/" + id)
		if err != nil {
			log.Print(err)
		}
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&application); err != nil {
			log.Print(err)
		}
		//fill the post textbox with the previous data when editing include id
		data = map[string]interface{}{
			"application": application,
		}
		temp, err := template.ParseFiles("views/detail_application.html")
		if err != nil {
			// Handle the error appropriately, e.g., log the error or return an HTTP error response.
			log.Println("Error parsing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = temp.Execute(w, data)
		if err != nil {
			// Handle the error appropriately, e.g., log the error or return an HTTP error response.
			log.Println("Error executing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		//getting the api to send us the datas
		response, err := http.Get(BASE_URL + "/applications")
		if err != nil {
			log.Print(err)
		}
		defer response.Body.Close()
		//decode the apis data into json
		//PROBLEM: is in the Decode(&posts)
		var applications []Application
		decoder := json.NewDecoder(response.Body)
		if err := decoder.Decode(&applications); err != nil { // <-- here
			log.Print(err)
		}
		fmt.Println("-------------------")
		fmt.Println(response.StatusCode)
		fmt.Println(response.Status)
		fmt.Println(applications)
		fmt.Println("-------------------")
		//send it into our html tags with the "posts"
		datas := map[string]interface{}{
			"applications": applications,
		}
		temp, err := template.ParseFiles("views/application_job.html")
		if err != nil {
			// Handle the error appropriately, e.g., log the error or return an HTTP error response.
			log.Println("Error parsing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = temp.Execute(w, datas)
		if err != nil {
			// Handle the error appropriately, e.g., log the error or return an HTTP error response.
			log.Println("Error executing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	}
}

func JobIndex(w http.ResponseWriter, r *http.Request) {
	var job Job
	var data map[string]interface{}
	//getting the id from the url when pressing one of the api data list
	id := r.URL.Query().Get("id")
	if id != "" {
		res, err := http.Get(BASE_URL + "/product/" + id)
		if err != nil {
			log.Print(err)
		}
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&job); err != nil {
			log.Print(err)
		}
		//fill the post textbox with the previous data when editing include id
		data = map[string]interface{}{
			"job": job,
		}
		temp, err := template.ParseFiles("views/detail_job.html")
		if err != nil {
			// Handle the error appropriately, e.g., log the error or return an HTTP error response.
			log.Println("Error parsing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = temp.Execute(w, data)
		if err != nil {
			// Handle the error appropriately, e.g., log the error or return an HTTP error response.
			log.Println("Error executing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		//getting the api to send us the datas
		response, err := http.Get(BASE_URL + "/products")
		if err != nil {
			log.Print(err)
		}
		defer response.Body.Close()
		//decode the apis data into json
		//PROBLEM: is in the Decode(&posts)
		var jobs []Job
		decoder := json.NewDecoder(response.Body)
		if err := decoder.Decode(&jobs); err != nil { // <-- here
			log.Print(err)
		}
		fmt.Println("-------------------")
		fmt.Println(response.StatusCode)
		fmt.Println(response.Status)
		fmt.Println(jobs)
		fmt.Println("-------------------")
		//send it into our html tags with the "posts"
		datas := map[string]interface{}{
			"jobs": jobs,
		}
		temp, err := template.ParseFiles("views/joblist.html")
		if err != nil {
			// Handle the error appropriately, e.g., log the error or return an HTTP error response.
			log.Println("Error parsing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = temp.Execute(w, datas)
		if err != nil {
			// Handle the error appropriately, e.g., log the error or return an HTTP error response.
			log.Println("Error executing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	}

}

func Create(w http.ResponseWriter, r *http.Request) {

	var job Job
	var data map[string]interface{}
	//getting the id from the url when pressing one of the api data list as PUT
	id := r.URL.Query().Get("id")
	if id != "" {
		res, err := http.Get(BASE_URL + "/product/" + id)
		if err != nil {
			log.Print(err)
		}
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&job); err != nil {
			log.Print(err)
		}
		//fill the post textbox with the previous data when editing include id
		data = map[string]interface{}{
			"job": job,
		}
	}
	//the part of the code that creates POST if there's no id present meaning null
	temp, _ := template.ParseFiles("views/create_job.html")
	temp.Execute(w, data)

}

func Store(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	//getting the previous filled id from form into a var
	id := r.Form.Get("job_id")
	println(" id => " + id)
	//find last index record from our table struct
	var lastVersion struct {
		ID int
	}
	DB.Table("jobs").Last(&lastVersion)
	// set something if no data in table
	var nextInt = int64(lastVersion.ID)
	if int64(lastVersion.ID) < 1 {
		//if first data is 1 or nothing set default
		nextInt = 1
	} else if id == "" {
		//if creating job that is not 1st or nothing data then plus one
		nextInt = nextInt + 1
	} // else if updating then .. no need for config

	newPost := Job{
		Id:          nextInt,
		Nama:        r.Form.Get("job_nama"),
		Organisasi:  r.Form.Get("job_organisasi"),
		Lokasi:      r.Form.Get("job_lokasi"),
		Deskripsi:   r.Form.Get("job_deskripsi"),
		Kualifikasi: r.Form.Get("job_kualifikasi"),
	}

	jsonValue, _ := json.Marshal(newPost)
	buff := bytes.NewBuffer(jsonValue)

	var req *http.Request
	var err error
	if id != "" {
		//update
		fmt.Println("Proses update")
		req, err = http.NewRequest(http.MethodPut, BASE_URL+"/product/"+id, buff)
	} else {
		// create
		fmt.Println("Proses create")
		req, err = http.NewRequest(http.MethodPost, BASE_URL+"/product", buff)
	}

	if err != nil {
		log.Print(err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Print(err)
	}
	defer res.Body.Close()

	var postResponse Job

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&postResponse); err != nil {
		log.Print(err)
	}

	fmt.Println(res.StatusCode)
	fmt.Println(res.Status)
	fmt.Println(postResponse)

	if res.StatusCode == 201 || res.StatusCode == 200 || res.StatusCode == 400 {
		http.Redirect(w, r, "/posts/job", http.StatusSeeOther)
	}

}

func Delete(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	req, err := http.NewRequest(http.MethodDelete, BASE_URL+"/product/"+id, nil)
	if err != nil {
		log.Print(err)
	}

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Print(err)
	}

	defer res.Body.Close()

	fmt.Println(res.StatusCode)
	fmt.Println(res.Status)

	if res.StatusCode == 200 {
		http.Redirect(w, r, "/posts/job", http.StatusSeeOther)
	}

}

// ------------------ Application Consume API -----------------------------------------------

func CreateApplication(w http.ResponseWriter, r *http.Request) {

	var job Job
	var application Application
	var data map[string]interface{}
	//getting the id from the url when pressing one of the api data list as PUT
	id := r.URL.Query().Get("id")
	appid := r.URL.Query().Get("appid")
	if id != "" {
		res, err := http.Get(BASE_URL + "/product/" + id)
		if err != nil {
			log.Print(err)
		}
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&job); err != nil {
			log.Print(err)
		}

		//fill the post textbox with the previous data when editing include id
		data = map[string]interface{}{
			"job": job,
		}
	} else if appid != "" {
		resa, err := http.Get(BASE_URL + "/application/" + appid)
		if err != nil {
			log.Print(err)
		}
		defer resa.Body.Close()

		decodera := json.NewDecoder(resa.Body)
		if err := decodera.Decode(&application); err != nil {
			log.Print(err)
		}

		res, err := http.Get(BASE_URL + "/product/" + fmt.Sprint(application.JobId))
		if err != nil {
			log.Print(err)
		}
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&job); err != nil {
			log.Print(err)
		}

		//fill the post textbox with the previous data when editing include id
		data = map[string]interface{}{
			"application": application,
			"job":         job,
		}

	}
	//the part of the code that creates POST if there's no id present meaning null
	temp, _ := template.ParseFiles("views/create_application.html")
	temp.Execute(w, data)

}

const MAX_UPLOAD_SIZES = (1024 * 1024) // 1MB
func StoreApplication(w http.ResponseWriter, r *http.Request) {
	t := r.FormValue("application_id")
	app_name := r.FormValue("application_nama")
	app_kontak := r.FormValue("application_kontak")
	app_umur := r.FormValue("application_umur")
	app_jekel := r.FormValue("application_jekel")
	app_deskripsi := r.FormValue("application_deskripsi")
	job_organisasi := r.FormValue("job_organisasi")
	id := r.FormValue("job_id")
	idInt, _ := strconv.ParseInt(id, 10, 64)
	var application_path = ""
	// #UPLOAD# get a reference to the fileHeaders
	files := r.MultipartForm.File["file"]
	// #UPLOAD# looping through each of them all multi files uploaded
	for _, fileHeader := range files {
		if fileHeader.Size > MAX_UPLOAD_SIZES {
			http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 1MB in size", fileHeader.Filename), http.StatusBadRequest)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" {
			http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = os.MkdirAll("./static/uploads", os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		timeN := time.Now().UnixNano()
		application_path = fmt.Sprintf("/uploads/%s_%d%s", fileHeader.Filename, timeN, filepath.Ext(fileHeader.Filename))
		full_path := fmt.Sprintf("./static/uploads/%s_%d%s", fileHeader.Filename, timeN, filepath.Ext(fileHeader.Filename))
		f, err := os.Create(full_path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer f.Close()

		pr := &Progress{
			TotalSize: fileHeader.Size,
		}

		_, err = io.Copy(f, io.TeeReader(file, pr))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	// #FORM# find last index record from our table struct
	var lastVersion struct {
		ID int
	}
	DB.Table("applications").Last(&lastVersion)
	// set something if no data in table
	var nextInt = int64(lastVersion.ID)
	if int64(lastVersion.ID) < 1 {
		//if first data is 1 or nothing set default
		nextInt = 1
	} else if t == "" {
		//if creating job that is not 1st or nothing data then plus one
		nextInt = nextInt + 1
	} // else if updating then .. no need for config

	newPost := Application{
		Id:            nextInt,
		Nama:          app_name,
		Kontak:        app_kontak,
		Umur:          app_umur,
		Jekel:         app_jekel,
		Deskripsi:     app_deskripsi,
		JobId:         idInt,
		JobOrganisasi: job_organisasi,
		Filepath:      application_path,
		Approval:      "false",
	}

	jsonValue, _ := json.Marshal(newPost)
	buff := bytes.NewBuffer(jsonValue)

	var req *http.Request
	var err error

	if t != "" {
		//update
		fmt.Println("Proses update application")
		req, err = http.NewRequest(http.MethodPut, BASE_URL+"/application/"+fmt.Sprint(nextInt), buff)
	} else {
		// create
		fmt.Println("Proses create application")
		req, err = http.NewRequest(http.MethodPost, BASE_URL+"/application", buff)
	}

	if err != nil {
		log.Print(err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Print(err)
	}
	defer res.Body.Close()

	var postResponse Application

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&postResponse); err != nil {
		log.Print(err)
	}

	fmt.Println(res.StatusCode)
	fmt.Println(res.Status)
	fmt.Println(postResponse)
	// After succesful apply documents application then go to..
	if res.StatusCode == 201 || res.StatusCode == 200 || res.StatusCode == 400 {
		http.Redirect(w, r, "/status/application", http.StatusSeeOther)
	}

}

func DeleteApplication(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	req, err := http.NewRequest(http.MethodDelete, BASE_URL+"/application/"+id, nil)
	if err != nil {
		log.Print(err)
	}

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Print(err)
	}

	defer res.Body.Close()

	fmt.Println(res.StatusCode)
	fmt.Println(res.Status)

	if res.StatusCode == 200 {
		http.Redirect(w, r, "/status/application", http.StatusSeeOther)
	}

}

// Write is used to satisfy the io.Writer interface.
// Instead of writing somewhere, it simply aggregates
// the total bytes on each read
func (pr *Progress) Write(p []byte) (n int, err error) {
	n, err = len(p), nil
	pr.BytesRead += int64(n)
	pr.Print()
	return
}

// Print displays the current progress of the file upload
func (pr *Progress) Print() {
	if pr.BytesRead == pr.TotalSize {
		fmt.Println("DONE!")
		return
	}

	fmt.Printf("File upload in progress: %d/%d\n", pr.BytesRead, pr.TotalSize)
}

func ApprovalApplication(w http.ResponseWriter, r *http.Request) {
	println("approving,,")
	id := r.URL.Query().Get("id")
	println(id)
	//var application Application
	app_id := r.FormValue("app_id")
	app_nama := r.FormValue("app_nama")
	app_kontak := r.FormValue("app_kontak")
	app_umur := r.FormValue("app_umur")
	app_jekel := r.FormValue("app_jekel")
	app_deskripsi := r.FormValue("app_deskripsi")
	app_jobid := r.FormValue("app_jobid")
	app_joborganisasi := r.FormValue("app_joborganisasi")
	app_filepath := r.FormValue("app_filepath")
	app_approval := r.FormValue("app_approval")
	int_app_id, _ := strconv.ParseInt(app_id, 10, 64)
	int_app_jobid, _ := strconv.ParseInt(app_jobid, 10, 64)

	//toggle
	println("before toggle: " + app_approval)
	if app_approval == "false" {
		app_approval = "true"
	} else if app_approval == "true" {
		app_approval = "false"
	}
	println("after toggle: " + app_approval)

	newPost := Application{
		Id:            int_app_id,
		Nama:          app_nama,
		Kontak:        app_kontak,
		Umur:          app_umur,
		Jekel:         app_jekel,
		Deskripsi:     app_deskripsi,
		JobId:         int_app_jobid,
		JobOrganisasi: app_joborganisasi,
		Filepath:      app_filepath,
		Approval:      app_approval,
	}

	jsonValue, _ := json.Marshal(newPost)
	buff := bytes.NewBuffer(jsonValue)

	var req *http.Request
	var err error

	//update
	fmt.Println("Proses update approval")
	req, err = http.NewRequest(http.MethodPut, BASE_URL+"/application/"+app_id, buff)

	if err != nil {
		log.Print(err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Print(err)
	}
	defer res.Body.Close()

	var postResponse Application

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&postResponse); err != nil {
		log.Print(err)
	}

	fmt.Println(res.StatusCode)
	fmt.Println(res.Status)
	fmt.Println(postResponse)

	// After succesful apply documents application then go to..
	if res.StatusCode == 200 || res.StatusCode == 400 {
		http.Redirect(w, r, "/status/application", http.StatusSeeOther)
	}

}

func main() {
	runtime.GOMAXPROCS(3) // use only 3 processor core
	fmt.Println("Version", runtime.Version())
	fmt.Println("NumCPU", runtime.NumCPU())
	fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))

	r := gin.Default()
	ConnectDatabase()

	r.GET("/products", API_IndexJob)
	r.GET("/product/:id", API_ShowJob)
	r.POST("/product", API_CreateJob)
	r.PUT("/product/:id", API_UpdateJob)
	r.DELETE("/product/:id", API_DeleteJob)

	r.GET("/applications", API_IndexApplication)
	r.GET("/application/:id", API_ShowApplication)
	r.POST("/application", API_CreateApplication)
	r.PUT("/application/:id", API_UpdateApplication)
	r.DELETE("/application/:id", API_DeleteApplication)

	r0 := http.NewServeMux()
	r1 := http.NewServeMux()

	r0.HandleFunc("/wiki", redirectToPath("https://www.wikihow.com/Main-Page"))
	r0.HandleFunc("/", Home)
	r0.HandleFunc("/posts/job", JobIndex)
	r0.HandleFunc("/post/create", Create)
	r0.HandleFunc("/post/store", Store)
	r0.HandleFunc("/post/delete", Delete)
	r0.HandleFunc("/status/application", ApplicationIndex)
	r0.HandleFunc("/post/createapplication", CreateApplication)
	r0.HandleFunc("/post/storeapplication", StoreApplication)
	r0.HandleFunc("/post/deleteapplication", DeleteApplication)
	r0.HandleFunc("/post/application/approval", ApprovalApplication)

	r1.HandleFunc("/wiki", redirectToPath("https://www.wikihow.com/Main-Page"))
	r1.HandleFunc("/", Home)
	r1.HandleFunc("/posts/job", JobIndex)
	r1.HandleFunc("/post/create", Create)
	r1.HandleFunc("/post/store", Store)
	r1.HandleFunc("/post/delete", Delete)
	r1.HandleFunc("/status/application", ApplicationIndex)
	r1.HandleFunc("/post/createapplication", CreateApplication)
	r1.HandleFunc("/post/storeapplication", StoreApplication)
	r1.HandleFunc("/post/deleteapplication", DeleteApplication)
	r1.HandleFunc("/post/application/approval", ApprovalApplication)

	// Create a file server.
	fileServer := http.FileServer(http.Dir("static"))

	// Serve the file server on port 2222.
	go func() { log.Fatal(http.ListenAndServe(":4444", r1)) }()
	go func() { log.Fatal(r.Run(":3333")) }()
	go func() { log.Fatal(http.ListenAndServe(":2222", fileServer)) }()
	go func() { log.Fatal(http.ListenAndServe(":1111", r0)) }()

	log.Print("Backup Server started on: http://localhost:4444")
	log.Print("API Server started on: http://localhost:3333")
	log.Print("File Server started on: http://localhost:2222")
	log.Print("Main Frontend Server started on: http://localhost:1111")
	select {}
}

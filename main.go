package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type App struct {
	db        *sql.DB
	templates *template.Template
	adminUser string
	adminPass string
}

type Service struct {
	ID           int
	Slug         string
	Name         string
	City         string
	Address      string
	Phone        string
	Email        string
	Website      string
	Description  string
	Specialties  string
	WorkingHours string
	Lat          float64
	Lng          float64
	Rating       float64
	ReviewCount  int
	Featured     bool
	ImageURL     string
	Offerings    []Offering
	Brands       []string
	Amenities    []string
	Gallery      []string
	Reviews      []Review
}

type Offering struct {
	ID        int
	ServiceID int
	Name      string
	PriceFrom string
}

type Review struct {
	ID             int
	ServiceID      int
	ServiceName    string
	AuthorName     string
	Title          string
	Comment        string
	Rating         int
	SpeedRating    int
	PriceRating    int
	QualityRating  int
	KindnessRating int
	Approved       bool
	CreatedAt      time.Time
}

type Article struct {
	ID          int
	Slug        string
	Title       string
	Excerpt     string
	Category    string
	ReadTime    string
	ImageURL    string
	RawBody     string
	Body        template.HTML
	Featured    bool
	PublishedAt time.Time
}

type CarModel struct {
	ID              int
	Slug            string
	Brand           string
	ModelName       string
	Years           string
	Engine          string
	ServiceInterval string
	KnownIssues     string
	TypicalCosts    string
	Summary         string
	ImageURL        string
	Featured        bool
}

type HomePageData struct {
	Title            string
	FeaturedServices []Service
	LatestArticles   []Article
	FeaturedModels   []CarModel
	Cities           []string
	Brands           []string
	ServiceTypes     []string
}

type ServicesPageData struct {
	Title         string
	Services      []Service
	Cities        []string
	Brands        []string
	ServiceTypes  []string
	Query         string
	City          string
	Brand         string
	ServiceType   string
	MinRating     string
	TotalServices int
}

type ServicePageData struct {
	Title           string
	Service         Service
	ReviewSubmitted bool
}

type ArticlesPageData struct {
	Title    string
	Articles []Article
}

type ArticlePageData struct {
	Title   string
	Article Article
}

type ModelsPageData struct {
	Title  string
	Models []CarModel
	Brands []string
	Brand  string
	Query  string
}

type ModelPageData struct {
	Title string
	Model CarModel
}

type AdminPageData struct {
	Title          string
	Section        string
	Services       []Service
	Articles       []Article
	Models         []CarModel
	Reviews        []Review
	EditService    Service
	EditArticle    Article
	EditModel      CarModel
	Message        string
	Today          string
	ServiceCount   int
	ArticleCount   int
	ModelCount     int
	PendingReviews int
}

func main() {
	dbPath := envOrDefault("DB_PATH", filepath.Join("data", "autoslot.db"))
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		log.Fatalf("create data dir: %v", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := initDB(db); err != nil {
		log.Fatalf("init db: %v", err)
	}

	tpl, err := template.New("").Funcs(template.FuncMap{
		"join":       strings.Join,
		"nl2br":      nl2br,
		"stars":      stars,
		"formatDate": func(t time.Time) string { return t.Format("02.01.2006") },
	}).ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		log.Fatalf("parse templates: %v", err)
	}

	app := &App{
		db:        db,
		templates: tpl,
		adminUser: envOrDefault("ADMIN_USER", "admin"),
		adminPass: envOrDefault("ADMIN_PASS", "admin123"),
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", app.handleHome)
	mux.HandleFunc("/services", app.handleServices)
	mux.HandleFunc("/services/", app.handleServiceRoutes)
	mux.HandleFunc("/articles", app.handleArticles)
	mux.HandleFunc("/articles/", app.handleArticleDetail)
	mux.HandleFunc("/models", app.handleModels)
	mux.HandleFunc("/models/", app.handleModelDetail)
	mux.Handle("/admin", app.basicAuth(http.HandlerFunc(app.handleAdminDashboard)))
	mux.Handle("/admin/", app.basicAuth(http.HandlerFunc(app.handleAdminRoutes)))

	addr := envOrDefault("APP_ADDR", ":8080")
	log.Printf("autoslot listening on %s", addr)
	if err := http.ListenAndServe(addr, logRequest(mux)); err != nil {
		log.Fatal(err)
	}
}

func initDB(db *sql.DB) error {
	schema, err := os.ReadFile(filepath.Join("data", "schema.sql"))
	if err != nil {
		return err
	}
	if _, err := db.Exec(string(schema)); err != nil {
		return err
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM services`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	seed, err := os.ReadFile(filepath.Join("data", "seed.sql"))
	if err != nil {
		return err
	}
	_, err = db.Exec(string(seed))
	return err
}

func (a *App) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	services, _ := a.listFeaturedServices()
	articles, _ := a.listLatestArticles(3)
	models, _ := a.listFeaturedModels(3)
	cities, brands, serviceTypes, _ := a.loadDirectoryFilters()

	a.render(w, "home.html", HomePageData{
		Title:            "AutoSlot | Pronadji servis bez lutanja",
		FeaturedServices: services,
		LatestArticles:   articles,
		FeaturedModels:   models,
		Cities:           cities,
		Brands:           brands,
		ServiceTypes:     serviceTypes,
	})
}

func (a *App) handleServices(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/services" {
		http.NotFound(w, r)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	city := strings.TrimSpace(r.URL.Query().Get("city"))
	brand := strings.TrimSpace(r.URL.Query().Get("brand"))
	serviceType := strings.TrimSpace(r.URL.Query().Get("service"))
	minRating := strings.TrimSpace(r.URL.Query().Get("min_rating"))

	services, err := a.searchServices(query, city, brand, serviceType, minRating)
	if err != nil {
		http.Error(w, "greska pri ucitavanju servisa", http.StatusInternalServerError)
		return
	}

	cities, brands, serviceTypes, _ := a.loadDirectoryFilters()
	a.render(w, "services.html", ServicesPageData{
		Title:         "Auto servisi",
		Services:      services,
		Cities:        cities,
		Brands:        brands,
		ServiceTypes:  serviceTypes,
		Query:         query,
		City:          city,
		Brand:         brand,
		ServiceType:   serviceType,
		MinRating:     minRating,
		TotalServices: len(services),
	})
}

func (a *App) handleServiceRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/services/")
	path = strings.Trim(path, "/")
	if path == "" {
		http.Redirect(w, r, "/services", http.StatusSeeOther)
		return
	}

	parts := strings.Split(path, "/")
	slug := parts[0]
	if len(parts) == 1 && r.Method == http.MethodGet {
		a.handleServiceDetail(w, r, slug)
		return
	}
	if len(parts) == 2 && parts[1] == "reviews" && r.Method == http.MethodPost {
		a.handleReviewSubmit(w, r, slug)
		return
	}

	http.NotFound(w, r)
}

func (a *App) handleServiceDetail(w http.ResponseWriter, r *http.Request, slug string) {
	service, err := a.getServiceBySlug(slug, true)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	a.render(w, "service_detail.html", ServicePageData{
		Title:           service.Name,
		Service:         service,
		ReviewSubmitted: r.URL.Query().Get("review") == "submitted",
	})
}

func (a *App) handleReviewSubmit(w http.ResponseWriter, r *http.Request, slug string) {
	service, err := a.getServiceBySlug(slug, false)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "neispravan formular", http.StatusBadRequest)
		return
	}

	review := Review{
		ServiceID:      service.ID,
		AuthorName:     strings.TrimSpace(r.FormValue("author_name")),
		Title:          strings.TrimSpace(r.FormValue("title")),
		Comment:        strings.TrimSpace(r.FormValue("comment")),
		Rating:         clampRating(parseInt(r.FormValue("rating"))),
		SpeedRating:    clampRating(parseInt(r.FormValue("speed_rating"))),
		PriceRating:    clampRating(parseInt(r.FormValue("price_rating"))),
		QualityRating:  clampRating(parseInt(r.FormValue("quality_rating"))),
		KindnessRating: clampRating(parseInt(r.FormValue("kindness_rating"))),
	}

	if review.AuthorName == "" || review.Comment == "" || review.Rating == 0 {
		http.Redirect(w, r, "/services/"+slug, http.StatusSeeOther)
		return
	}
	if review.Title == "" {
		review.Title = "Utisak korisnika"
	}

	_, err = a.db.Exec(`
		INSERT INTO reviews (
			service_id, author_name, title, comment, rating,
			speed_rating, price_rating, quality_rating, kindness_rating,
			approved, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0, CURRENT_TIMESTAMP)
	`,
		review.ServiceID, review.AuthorName, review.Title, review.Comment, review.Rating,
		review.SpeedRating, review.PriceRating, review.QualityRating, review.KindnessRating,
	)
	if err != nil {
		http.Error(w, "greska pri cuvanju recenzije", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/services/"+slug+"?review=submitted", http.StatusSeeOther)
}

func (a *App) handleArticles(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/articles" {
		http.NotFound(w, r)
		return
	}
	articles, err := a.listAllArticles()
	if err != nil {
		http.Error(w, "greska pri ucitavanju tekstova", http.StatusInternalServerError)
		return
	}
	a.render(w, "articles.html", ArticlesPageData{Title: "Saveti za vozace", Articles: articles})
}

func (a *App) handleArticleDetail(w http.ResponseWriter, r *http.Request) {
	slug := strings.Trim(strings.TrimPrefix(r.URL.Path, "/articles/"), "/")
	if slug == "" {
		http.Redirect(w, r, "/articles", http.StatusSeeOther)
		return
	}
	article, err := a.getArticleBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	a.render(w, "article_detail.html", ArticlePageData{Title: article.Title, Article: article})
}

func (a *App) handleModels(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/models" {
		http.NotFound(w, r)
		return
	}
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	brand := strings.TrimSpace(r.URL.Query().Get("brand"))
	models, brands, err := a.listModels(query, brand)
	if err != nil {
		http.Error(w, "greska pri ucitavanju modela", http.StatusInternalServerError)
		return
	}
	a.render(w, "models.html", ModelsPageData{
		Title:  "Modeli vozila i poznati problemi",
		Models: models,
		Brands: brands,
		Brand:  brand,
		Query:  query,
	})
}

func (a *App) handleModelDetail(w http.ResponseWriter, r *http.Request) {
	slug := strings.Trim(strings.TrimPrefix(r.URL.Path, "/models/"), "/")
	if slug == "" {
		http.Redirect(w, r, "/models", http.StatusSeeOther)
		return
	}
	model, err := a.getModelBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	a.render(w, "model_detail.html", ModelPageData{Title: model.Brand + " " + model.ModelName, Model: model})
}

func (a *App) handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/admin/services", http.StatusSeeOther)
}

func (a *App) handleAdminRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/admin/"), "/")
	switch {
	case path == "services" && r.Method == http.MethodGet:
		a.handleAdminServices(w, r)
	case path == "services" && r.Method == http.MethodPost:
		a.handleAdminServiceSave(w, r)
	case path == "articles" && r.Method == http.MethodGet:
		a.handleAdminArticles(w, r)
	case path == "articles" && r.Method == http.MethodPost:
		a.handleAdminArticleSave(w, r)
	case path == "models" && r.Method == http.MethodGet:
		a.handleAdminModels(w, r)
	case path == "models" && r.Method == http.MethodPost:
		a.handleAdminModelSave(w, r)
	case path == "reviews" && r.Method == http.MethodGet:
		a.handleAdminReviews(w, r)
	case strings.HasPrefix(path, "reviews/") && r.Method == http.MethodPost:
		a.handleAdminReviewAction(w, r, strings.TrimPrefix(path, "reviews/"))
	default:
		http.NotFound(w, r)
	}
}

func (a *App) handleAdminServices(w http.ResponseWriter, r *http.Request) {
	services, _ := a.listAllServices()
	editID := parseInt(r.URL.Query().Get("edit"))
	var edit Service
	if editID > 0 {
		edit, _ = a.getServiceByID(editID)
	}
	data := a.adminData("services", r.URL.Query().Get("msg"))
	data.Services = services
	data.EditService = edit
	a.render(w, "admin_services.html", data)
}

func (a *App) handleAdminServiceSave(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "neispravan formular", http.StatusBadRequest)
		return
	}

	id := parseInt(r.FormValue("id"))
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Redirect(w, r, "/admin/services?msg=Ime+servisa+je+obavezno", http.StatusSeeOther)
		return
	}

	slug := slugify(strings.TrimSpace(r.FormValue("slug")))
	if slug == "" {
		slug = slugify(name)
	}

	service := Service{
		ID:           id,
		Slug:         slug,
		Name:         name,
		City:         strings.TrimSpace(r.FormValue("city")),
		Address:      strings.TrimSpace(r.FormValue("address")),
		Phone:        strings.TrimSpace(r.FormValue("phone")),
		Email:        strings.TrimSpace(r.FormValue("email")),
		Website:      strings.TrimSpace(r.FormValue("website")),
		Description:  strings.TrimSpace(r.FormValue("description")),
		Specialties:  strings.TrimSpace(r.FormValue("specialties")),
		WorkingHours: strings.TrimSpace(r.FormValue("working_hours")),
		Lat:          parseFloat(r.FormValue("lat")),
		Lng:          parseFloat(r.FormValue("lng")),
		ImageURL:     strings.TrimSpace(r.FormValue("image_url")),
		Featured:     r.FormValue("featured") == "on",
	}

	tx, err := a.db.Begin()
	if err != nil {
		http.Error(w, "greska baze", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	if service.ID > 0 {
		_, err = tx.Exec(`
			UPDATE services SET slug=?, name=?, city=?, address=?, phone=?, email=?, website=?,
				description=?, specialties=?, working_hours=?, lat=?, lng=?, featured=?, image_url=?
			WHERE id=?
		`, service.Slug, service.Name, service.City, service.Address, service.Phone, service.Email, service.Website,
			service.Description, service.Specialties, service.WorkingHours, service.Lat, service.Lng, boolToInt(service.Featured), service.ImageURL, service.ID)
	} else {
		res, execErr := tx.Exec(`
			INSERT INTO services (
				slug, name, city, address, phone, email, website, description, specialties,
				working_hours, lat, lng, rating, review_count, featured, image_url
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, 0, ?, ?)
		`, service.Slug, service.Name, service.City, service.Address, service.Phone, service.Email, service.Website,
			service.Description, service.Specialties, service.WorkingHours, service.Lat, service.Lng, boolToInt(service.Featured), service.ImageURL)
		err = execErr
		if err == nil {
			insertID, _ := res.LastInsertId()
			service.ID = int(insertID)
		}
	}
	if err != nil {
		http.Redirect(w, r, "/admin/services?msg=Greska+pri+snimanju+servisa", http.StatusSeeOther)
		return
	}

	if err := syncStringsTx(tx, "service_offerings", "service_id", "name", "price_from", service.ID, parsePairs(r.FormValue("offerings"))); err != nil {
		http.Error(w, "greska pri snimanju usluga", http.StatusInternalServerError)
		return
	}
	if err := syncSingleColumnTx(tx, "service_brands", "service_id", "brand", service.ID, parseCSV(r.FormValue("brands"))); err != nil {
		http.Error(w, "greska pri snimanju marki", http.StatusInternalServerError)
		return
	}
	if err := syncSingleColumnTx(tx, "service_amenities", "service_id", "amenity", service.ID, parseCSV(r.FormValue("amenities"))); err != nil {
		http.Error(w, "greska pri snimanju pogodnosti", http.StatusInternalServerError)
		return
	}
	if err := syncSingleColumnTx(tx, "service_gallery", "service_id", "image_url", service.ID, parseLines(r.FormValue("gallery"))); err != nil {
		http.Error(w, "greska pri snimanju galerije", http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, "greska pri potvrdi izmene", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/services?msg=Servis+je+sacuvan", http.StatusSeeOther)
}

func (a *App) handleAdminArticles(w http.ResponseWriter, r *http.Request) {
	articles, _ := a.listAllArticles()
	editID := parseInt(r.URL.Query().Get("edit"))
	var edit Article
	if editID > 0 {
		edit, _ = a.getArticleByID(editID)
	}
	data := a.adminData("articles", r.URL.Query().Get("msg"))
	data.Articles = articles
	data.EditArticle = edit
	a.render(w, "admin_articles.html", data)
}

func (a *App) handleAdminArticleSave(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "neispravan formular", http.StatusBadRequest)
		return
	}

	id := parseInt(r.FormValue("id"))
	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		http.Redirect(w, r, "/admin/articles?msg=Naslov+je+obavezan", http.StatusSeeOther)
		return
	}
	slug := slugify(strings.TrimSpace(r.FormValue("slug")))
	if slug == "" {
		slug = slugify(title)
	}
	publishedAt := strings.TrimSpace(r.FormValue("published_at"))
	if publishedAt == "" {
		publishedAt = time.Now().Format("2006-01-02")
	}

	var err error
	if id > 0 {
		_, err = a.db.Exec(`
			UPDATE articles SET slug=?, title=?, excerpt=?, category=?, read_time=?, image_url=?, body=?, featured=?, published_at=?
			WHERE id=?
		`, slug, title, strings.TrimSpace(r.FormValue("excerpt")), strings.TrimSpace(r.FormValue("category")),
			strings.TrimSpace(r.FormValue("read_time")), strings.TrimSpace(r.FormValue("image_url")),
			strings.TrimSpace(r.FormValue("body")), boolToInt(r.FormValue("featured") == "on"), publishedAt, id)
	} else {
		_, err = a.db.Exec(`
			INSERT INTO articles (slug, title, excerpt, category, read_time, image_url, body, featured, published_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, slug, title, strings.TrimSpace(r.FormValue("excerpt")), strings.TrimSpace(r.FormValue("category")),
			strings.TrimSpace(r.FormValue("read_time")), strings.TrimSpace(r.FormValue("image_url")),
			strings.TrimSpace(r.FormValue("body")), boolToInt(r.FormValue("featured") == "on"), publishedAt)
	}
	if err != nil {
		http.Redirect(w, r, "/admin/articles?msg=Greska+pri+snimanju+teksta", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/articles?msg=Tekst+je+sacuvan", http.StatusSeeOther)
}

func (a *App) handleAdminModels(w http.ResponseWriter, r *http.Request) {
	models, _, _ := a.listModels("", "")
	editID := parseInt(r.URL.Query().Get("edit"))
	var edit CarModel
	if editID > 0 {
		edit, _ = a.getModelByID(editID)
	}
	data := a.adminData("models", r.URL.Query().Get("msg"))
	data.Models = models
	data.EditModel = edit
	a.render(w, "admin_models.html", data)
}

func (a *App) handleAdminModelSave(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "neispravan formular", http.StatusBadRequest)
		return
	}

	id := parseInt(r.FormValue("id"))
	brand := strings.TrimSpace(r.FormValue("brand"))
	modelName := strings.TrimSpace(r.FormValue("model_name"))
	if brand == "" || modelName == "" {
		http.Redirect(w, r, "/admin/models?msg=Marka+i+model+su+obavezni", http.StatusSeeOther)
		return
	}

	slug := slugify(strings.TrimSpace(r.FormValue("slug")))
	if slug == "" {
		slug = slugify(brand + "-" + modelName)
	}

	var err error
	if id > 0 {
		_, err = a.db.Exec(`
			UPDATE models SET slug=?, brand=?, model_name=?, years=?, engine=?, service_interval=?, known_issues=?, typical_costs=?, summary=?, image_url=?, featured=?
			WHERE id=?
		`, slug, brand, modelName, strings.TrimSpace(r.FormValue("years")), strings.TrimSpace(r.FormValue("engine")),
			strings.TrimSpace(r.FormValue("service_interval")), strings.TrimSpace(r.FormValue("known_issues")),
			strings.TrimSpace(r.FormValue("typical_costs")), strings.TrimSpace(r.FormValue("summary")),
			strings.TrimSpace(r.FormValue("image_url")), boolToInt(r.FormValue("featured") == "on"), id)
	} else {
		_, err = a.db.Exec(`
			INSERT INTO models (slug, brand, model_name, years, engine, service_interval, known_issues, typical_costs, summary, image_url, featured)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, slug, brand, modelName, strings.TrimSpace(r.FormValue("years")), strings.TrimSpace(r.FormValue("engine")),
			strings.TrimSpace(r.FormValue("service_interval")), strings.TrimSpace(r.FormValue("known_issues")),
			strings.TrimSpace(r.FormValue("typical_costs")), strings.TrimSpace(r.FormValue("summary")),
			strings.TrimSpace(r.FormValue("image_url")), boolToInt(r.FormValue("featured") == "on"))
	}
	if err != nil {
		http.Redirect(w, r, "/admin/models?msg=Greska+pri+snimanju+modela", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/models?msg=Model+je+sacuvan", http.StatusSeeOther)
}

func (a *App) handleAdminReviews(w http.ResponseWriter, r *http.Request) {
	rows, err := a.db.Query(`
		SELECT r.id, r.service_id, s.name, r.author_name, r.title, r.comment, r.rating,
			r.speed_rating, r.price_rating, r.quality_rating, r.kindness_rating, r.approved, r.created_at
		FROM reviews r
		JOIN services s ON s.id = r.service_id
		ORDER BY r.created_at DESC
	`)
	if err != nil {
		http.Error(w, "greska pri ucitavanju recenzija", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		if err := rows.Scan(&review.ID, &review.ServiceID, &review.ServiceName, &review.AuthorName, &review.Title,
			&review.Comment, &review.Rating, &review.SpeedRating, &review.PriceRating, &review.QualityRating,
			&review.KindnessRating, &review.Approved, &review.CreatedAt); err == nil {
			reviews = append(reviews, review)
		}
	}

	data := a.adminData("reviews", r.URL.Query().Get("msg"))
	data.Reviews = reviews
	a.render(w, "admin_reviews.html", data)
}

func (a *App) handleAdminReviewAction(w http.ResponseWriter, r *http.Request, path string) {
	id := parseInt(strings.Trim(path, "/"))
	if id == 0 {
		http.NotFound(w, r)
		return
	}
	action := r.FormValue("action")
	var approved int
	if action == "approve" {
		approved = 1
	}
	_, err := a.db.Exec(`UPDATE reviews SET approved=? WHERE id=?`, approved, id)
	if err != nil {
		http.Redirect(w, r, "/admin/reviews?msg=Greska+pri+izmeni+recenzije", http.StatusSeeOther)
		return
	}
	var serviceID int
	if err := a.db.QueryRow(`SELECT service_id FROM reviews WHERE id=?`, id).Scan(&serviceID); err == nil {
		_ = a.refreshServiceRating(serviceID)
	}
	http.Redirect(w, r, "/admin/reviews?msg=Recenzija+je+azurirana", http.StatusSeeOther)
}

func (a *App) adminData(section, message string) AdminPageData {
	data := AdminPageData{
		Title:   "Admin panel",
		Section: section,
		Message: strings.ReplaceAll(message, "+", " "),
		Today:   time.Now().Format("2006-01-02"),
	}
	_ = a.db.QueryRow(`SELECT COUNT(*) FROM services`).Scan(&data.ServiceCount)
	_ = a.db.QueryRow(`SELECT COUNT(*) FROM articles`).Scan(&data.ArticleCount)
	_ = a.db.QueryRow(`SELECT COUNT(*) FROM models`).Scan(&data.ModelCount)
	_ = a.db.QueryRow(`SELECT COUNT(*) FROM reviews WHERE approved = 0`).Scan(&data.PendingReviews)
	return data
}

func (a *App) render(w http.ResponseWriter, name string, data any) {
	if err := a.templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *App) basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != a.adminUser || pass != a.adminPass {
			w.Header().Set("WWW-Authenticate", `Basic realm="autoslot-admin"`)
			http.Error(w, "autorizacija je potrebna", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *App) searchServices(query, city, brand, serviceType, minRating string) ([]Service, error) {
	sqlBuilder := `
		SELECT DISTINCT s.id, s.slug, s.name, s.city, s.address, s.phone, s.email, s.website,
			s.description, s.specialties, s.working_hours, s.lat, s.lng, s.rating, s.review_count,
			s.featured, s.image_url
		FROM services s
		LEFT JOIN service_brands sb ON sb.service_id = s.id
		LEFT JOIN service_offerings so ON so.service_id = s.id
		WHERE 1=1
	`
	var args []any
	if query != "" {
		sqlBuilder += ` AND (LOWER(s.name) LIKE ? OR LOWER(s.description) LIKE ? OR LOWER(s.specialties) LIKE ?)`
		q := "%" + strings.ToLower(query) + "%"
		args = append(args, q, q, q)
	}
	if city != "" {
		sqlBuilder += ` AND s.city = ?`
		args = append(args, city)
	}
	if brand != "" {
		sqlBuilder += ` AND sb.brand = ?`
		args = append(args, brand)
	}
	if serviceType != "" {
		sqlBuilder += ` AND so.name = ?`
		args = append(args, serviceType)
	}
	if minRating != "" {
		sqlBuilder += ` AND s.rating >= ?`
		args = append(args, parseFloat(minRating))
	}
	sqlBuilder += ` ORDER BY s.featured DESC, s.rating DESC, s.review_count DESC, s.name ASC`

	rows, err := a.db.Query(sqlBuilder, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		service, err := scanService(rows)
		if err != nil {
			return nil, err
		}
		service.Offerings, _ = a.listOfferings(service.ID)
		service.Brands, _ = a.listStrings(`SELECT brand FROM service_brands WHERE service_id=? ORDER BY brand`, service.ID)
		service.Amenities, _ = a.listStrings(`SELECT amenity FROM service_amenities WHERE service_id=? ORDER BY amenity`, service.ID)
		services = append(services, service)
	}
	return services, nil
}

func (a *App) listFeaturedServices() ([]Service, error) {
	rows, err := a.db.Query(`
		SELECT id, slug, name, city, address, phone, email, website, description, specialties,
			working_hours, lat, lng, rating, review_count, featured, image_url
		FROM services
		ORDER BY featured DESC, rating DESC, review_count DESC
		LIMIT 6
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		service, err := scanService(rows)
		if err != nil {
			return nil, err
		}
		service.Offerings, _ = a.listOfferings(service.ID)
		service.Brands, _ = a.listStrings(`SELECT brand FROM service_brands WHERE service_id=? ORDER BY brand`, service.ID)
		services = append(services, service)
	}
	return services, nil
}

func (a *App) listAllServices() ([]Service, error) {
	rows, err := a.db.Query(`
		SELECT id, slug, name, city, address, phone, email, website, description, specialties,
			working_hours, lat, lng, rating, review_count, featured, image_url
		FROM services
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		service, err := scanService(rows)
		if err != nil {
			return nil, err
		}
		service.Offerings, _ = a.listOfferings(service.ID)
		service.Brands, _ = a.listStrings(`SELECT brand FROM service_brands WHERE service_id=? ORDER BY brand`, service.ID)
		service.Amenities, _ = a.listStrings(`SELECT amenity FROM service_amenities WHERE service_id=? ORDER BY amenity`, service.ID)
		service.Gallery, _ = a.listStrings(`SELECT image_url FROM service_gallery WHERE service_id=? ORDER BY id`, service.ID)
		services = append(services, service)
	}
	return services, nil
}

func (a *App) getServiceBySlug(slug string, withReviews bool) (Service, error) {
	var service Service
	err := a.db.QueryRow(`
		SELECT id, slug, name, city, address, phone, email, website, description, specialties,
			working_hours, lat, lng, rating, review_count, featured, image_url
		FROM services WHERE slug = ?
	`, slug).Scan(
		&service.ID, &service.Slug, &service.Name, &service.City, &service.Address, &service.Phone, &service.Email,
		&service.Website, &service.Description, &service.Specialties, &service.WorkingHours, &service.Lat, &service.Lng,
		&service.Rating, &service.ReviewCount, &service.Featured, &service.ImageURL,
	)
	if err != nil {
		return service, err
	}
	service.Offerings, _ = a.listOfferings(service.ID)
	service.Brands, _ = a.listStrings(`SELECT brand FROM service_brands WHERE service_id=? ORDER BY brand`, service.ID)
	service.Amenities, _ = a.listStrings(`SELECT amenity FROM service_amenities WHERE service_id=? ORDER BY amenity`, service.ID)
	service.Gallery, _ = a.listStrings(`SELECT image_url FROM service_gallery WHERE service_id=? ORDER BY id`, service.ID)

	if withReviews {
		rows, err := a.db.Query(`
			SELECT id, service_id, author_name, title, comment, rating, speed_rating, price_rating,
				quality_rating, kindness_rating, approved, created_at
			FROM reviews
			WHERE service_id = ? AND approved = 1
			ORDER BY created_at DESC
		`, service.ID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var review Review
				if scanErr := rows.Scan(&review.ID, &review.ServiceID, &review.AuthorName, &review.Title, &review.Comment,
					&review.Rating, &review.SpeedRating, &review.PriceRating, &review.QualityRating, &review.KindnessRating,
					&review.Approved, &review.CreatedAt); scanErr == nil {
					service.Reviews = append(service.Reviews, review)
				}
			}
		}
	}
	return service, nil
}

func (a *App) getServiceByID(id int) (Service, error) {
	var service Service
	err := a.db.QueryRow(`
		SELECT id, slug, name, city, address, phone, email, website, description, specialties,
			working_hours, lat, lng, rating, review_count, featured, image_url
		FROM services WHERE id = ?
	`, id).Scan(
		&service.ID, &service.Slug, &service.Name, &service.City, &service.Address, &service.Phone, &service.Email,
		&service.Website, &service.Description, &service.Specialties, &service.WorkingHours, &service.Lat, &service.Lng,
		&service.Rating, &service.ReviewCount, &service.Featured, &service.ImageURL,
	)
	if err != nil {
		return service, err
	}
	service.Offerings, _ = a.listOfferings(service.ID)
	service.Brands, _ = a.listStrings(`SELECT brand FROM service_brands WHERE service_id=? ORDER BY brand`, service.ID)
	service.Amenities, _ = a.listStrings(`SELECT amenity FROM service_amenities WHERE service_id=? ORDER BY amenity`, service.ID)
	service.Gallery, _ = a.listStrings(`SELECT image_url FROM service_gallery WHERE service_id=? ORDER BY id`, service.ID)
	return service, nil
}

func (a *App) listLatestArticles(limit int) ([]Article, error) {
	rows, err := a.db.Query(`
		SELECT id, slug, title, excerpt, category, read_time, image_url, body, featured, published_at
		FROM articles
		ORDER BY date(published_at) DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		article, scanErr := scanArticle(rows)
		if scanErr == nil {
			articles = append(articles, article)
		}
	}
	return articles, nil
}

func (a *App) listAllArticles() ([]Article, error) {
	rows, err := a.db.Query(`
		SELECT id, slug, title, excerpt, category, read_time, image_url, body, featured, published_at
		FROM articles
		ORDER BY date(published_at) DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		article, scanErr := scanArticle(rows)
		if scanErr == nil {
			articles = append(articles, article)
		}
	}
	return articles, nil
}

func (a *App) getArticleBySlug(slug string) (Article, error) {
	row := a.db.QueryRow(`
		SELECT id, slug, title, excerpt, category, read_time, image_url, body, featured, published_at
		FROM articles WHERE slug = ?
	`, slug)
	return scanArticleRow(row)
}

func (a *App) getArticleByID(id int) (Article, error) {
	row := a.db.QueryRow(`
		SELECT id, slug, title, excerpt, category, read_time, image_url, body, featured, published_at
		FROM articles WHERE id = ?
	`, id)
	return scanArticleRow(row)
}

func (a *App) listFeaturedModels(limit int) ([]CarModel, error) {
	rows, err := a.db.Query(`
		SELECT id, slug, brand, model_name, years, engine, service_interval, known_issues, typical_costs, summary, image_url, featured
		FROM models
		ORDER BY featured DESC, brand ASC, model_name ASC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []CarModel
	for rows.Next() {
		model, scanErr := scanModel(rows)
		if scanErr == nil {
			models = append(models, model)
		}
	}
	return models, nil
}

func (a *App) listModels(query, brand string) ([]CarModel, []string, error) {
	sqlBuilder := `
		SELECT id, slug, brand, model_name, years, engine, service_interval, known_issues, typical_costs, summary, image_url, featured
		FROM models
		WHERE 1=1
	`
	var args []any
	if query != "" {
		sqlBuilder += ` AND (LOWER(brand) LIKE ? OR LOWER(model_name) LIKE ? OR LOWER(summary) LIKE ? OR LOWER(known_issues) LIKE ?)`
		q := "%" + strings.ToLower(query) + "%"
		args = append(args, q, q, q, q)
	}
	if brand != "" {
		sqlBuilder += ` AND brand = ?`
		args = append(args, brand)
	}
	sqlBuilder += ` ORDER BY featured DESC, brand ASC, model_name ASC`

	rows, err := a.db.Query(sqlBuilder, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var models []CarModel
	for rows.Next() {
		model, scanErr := scanModel(rows)
		if scanErr == nil {
			models = append(models, model)
		}
	}
	brands, _ := a.listStrings(`SELECT DISTINCT brand FROM models ORDER BY brand`)
	return models, brands, nil
}

func (a *App) getModelBySlug(slug string) (CarModel, error) {
	row := a.db.QueryRow(`
		SELECT id, slug, brand, model_name, years, engine, service_interval, known_issues, typical_costs, summary, image_url, featured
		FROM models WHERE slug = ?
	`, slug)
	return scanModelRow(row)
}

func (a *App) getModelByID(id int) (CarModel, error) {
	row := a.db.QueryRow(`
		SELECT id, slug, brand, model_name, years, engine, service_interval, known_issues, typical_costs, summary, image_url, featured
		FROM models WHERE id = ?
	`, id)
	return scanModelRow(row)
}

func (a *App) listOfferings(serviceID int) ([]Offering, error) {
	rows, err := a.db.Query(`SELECT id, service_id, name, price_from FROM service_offerings WHERE service_id=? ORDER BY id`, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var offerings []Offering
	for rows.Next() {
		var offering Offering
		if err := rows.Scan(&offering.ID, &offering.ServiceID, &offering.Name, &offering.PriceFrom); err == nil {
			offerings = append(offerings, offering)
		}
	}
	return offerings, nil
}

func (a *App) loadDirectoryFilters() ([]string, []string, []string, error) {
	cities, err := a.listStrings(`SELECT DISTINCT city FROM services WHERE city <> '' ORDER BY city`)
	if err != nil {
		return nil, nil, nil, err
	}
	brands, err := a.listStrings(`SELECT DISTINCT brand FROM service_brands ORDER BY brand`)
	if err != nil {
		return nil, nil, nil, err
	}
	serviceTypes, err := a.listStrings(`SELECT DISTINCT name FROM service_offerings ORDER BY name`)
	return cities, brands, serviceTypes, err
}

func (a *App) listStrings(query string, args ...any) ([]string, error) {
	rows, err := a.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err == nil {
			values = append(values, value)
		}
	}
	return values, nil
}

func (a *App) refreshServiceRating(serviceID int) error {
	var avg sql.NullFloat64
	var count int
	if err := a.db.QueryRow(`SELECT AVG(rating), COUNT(*) FROM reviews WHERE service_id=? AND approved=1`, serviceID).Scan(&avg, &count); err != nil {
		return err
	}
	rating := 0.0
	if avg.Valid {
		rating = avg.Float64
	}
	_, err := a.db.Exec(`UPDATE services SET rating=?, review_count=? WHERE id=?`, rating, count, serviceID)
	return err
}

func syncSingleColumnTx(tx *sql.Tx, table, fkColumn, valueColumn string, id int, values []string) error {
	if _, err := tx.Exec(fmt.Sprintf(`DELETE FROM %s WHERE %s=?`, table, fkColumn), id); err != nil {
		return err
	}
	for _, value := range values {
		if _, err := tx.Exec(fmt.Sprintf(`INSERT INTO %s (%s, %s) VALUES (?, ?)`, table, fkColumn, valueColumn), id, value); err != nil {
			return err
		}
	}
	return nil
}

func syncStringsTx(tx *sql.Tx, table, fkColumn, leftColumn, rightColumn string, id int, pairs [][2]string) error {
	if _, err := tx.Exec(fmt.Sprintf(`DELETE FROM %s WHERE %s=?`, table, fkColumn), id); err != nil {
		return err
	}
	for _, pair := range pairs {
		if _, err := tx.Exec(fmt.Sprintf(`INSERT INTO %s (%s, %s, %s) VALUES (?, ?, ?)`, table, fkColumn, leftColumn, rightColumn), id, pair[0], pair[1]); err != nil {
			return err
		}
	}
	return nil
}

func parseCSV(value string) []string {
	parts := strings.Split(value, ",")
	var out []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	sort.Strings(out)
	return out
}

func parseLines(value string) []string {
	lines := strings.Split(value, "\n")
	var out []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func parsePairs(value string) [][2]string {
	lines := parseLines(value)
	var pairs [][2]string
	for _, line := range lines {
		parts := strings.SplitN(line, "|", 2)
		left := strings.TrimSpace(parts[0])
		right := ""
		if len(parts) > 1 {
			right = strings.TrimSpace(parts[1])
		}
		if left != "" {
			pairs = append(pairs, [2]string{left, right})
		}
	}
	return pairs
}

func parseInt(value string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(value))
	return n
}

func parseFloat(value string) float64 {
	n, _ := strconv.ParseFloat(strings.TrimSpace(value), 64)
	return n
}

func clampRating(value int) int {
	if value < 1 {
		return 0
	}
	if value > 5 {
		return 5
	}
	return value
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func scanService(scanner interface{ Scan(dest ...any) error }) (Service, error) {
	var service Service
	err := scanner.Scan(
		&service.ID, &service.Slug, &service.Name, &service.City, &service.Address, &service.Phone, &service.Email,
		&service.Website, &service.Description, &service.Specialties, &service.WorkingHours, &service.Lat, &service.Lng,
		&service.Rating, &service.ReviewCount, &service.Featured, &service.ImageURL,
	)
	return service, err
}

func scanArticle(scanner interface{ Scan(dest ...any) error }) (Article, error) {
	var article Article
	var body string
	err := scanner.Scan(&article.ID, &article.Slug, &article.Title, &article.Excerpt, &article.Category,
		&article.ReadTime, &article.ImageURL, &body, &article.Featured, &article.PublishedAt)
	article.RawBody = body
	article.Body = nl2br(body)
	return article, err
}

func scanArticleRow(row *sql.Row) (Article, error) {
	article, err := scanArticle(row)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return article, err
	}
	return article, err
}

func scanModel(scanner interface{ Scan(dest ...any) error }) (CarModel, error) {
	var model CarModel
	err := scanner.Scan(&model.ID, &model.Slug, &model.Brand, &model.ModelName, &model.Years, &model.Engine,
		&model.ServiceInterval, &model.KnownIssues, &model.TypicalCosts, &model.Summary, &model.ImageURL, &model.Featured)
	return model, err
}

func scanModelRow(row *sql.Row) (CarModel, error) {
	model, err := scanModel(row)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return model, err
	}
	return model, err
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacements := map[string]string{
		"č": "c",
		"ć": "c",
		"š": "s",
		"ž": "z",
		"đ": "dj",
	}
	for old, newValue := range replacements {
		value = strings.ReplaceAll(value, old, newValue)
	}
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteRune('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func nl2br(value string) template.HTML {
	safe := template.HTMLEscapeString(value)
	safe = strings.ReplaceAll(safe, "\n", "<br>")
	return template.HTML(safe)
}

func stars(rating int) string {
	if rating < 0 {
		rating = 0
	}
	if rating > 5 {
		rating = 5
	}
	return strings.Repeat("*", rating) + strings.Repeat("-", 5-rating)
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start).Round(time.Millisecond))
	})
}

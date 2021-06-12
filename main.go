package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type Team struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

type ScoreLog struct {
	Id        int       `db:"id"`
	TeamId    int       `db:"team_id"`
	Score     int       `db:"score"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

type Log struct {
	TeamName string    `json:"team_name"`
	Score    int       `json:"score"`
	Message  string    `json:"message"`
	LogTime  time.Time `json:"log_time"`
}
type Report struct {
	TeamName string `json:"team_name"`
	Score    int    `json:"score"`
	Message  string `json:"message"`
}
type ScoreBoard struct {
	Len  int   `json:"len"`
	Logs []Log `json:"logs"`
}

type Template struct {
	templates *template.Template
}

type State struct {
	DB *sqlx.DB
}

type CustomValidator struct {
	validator *validator.Validate
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (s *State) getBoard(c echo.Context) error {
	rows, err := s.DB.Queryx("SELECT * FROM score_log")
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch data from DB.")
	}
	var logs []Log
	for rows.Next() {
		var scoreLog ScoreLog
		err = rows.StructScan(&scoreLog)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to convert data to struct.")
		}
		var TeamName string
		err = s.DB.Get(&TeamName, "SELECT name FROM team WHERE id = ?", scoreLog.TeamId)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to Get TeamName.")
		}
		res := Log{
			TeamName: TeamName,
			Score:    scoreLog.Score,
			Message:  scoreLog.Message,
			LogTime:  scoreLog.CreatedAt,
		}
		logs = append(logs, res)
	}
	sort.SliceStable(logs, func(i, j int) bool { return logs[i].LogTime.After(logs[j].LogTime) })
	return c.Render(http.StatusOK, "index.html", logs)
}

func (s *State) reportScore(c echo.Context) (err error) {
	report := new(Report)
	if err = c.Validate(report); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to validate report. Check fields.")
	}
	if err = c.Bind(report); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to convert report to Struct")
	}
	// Validate Struct
	if _, err = s.DB.Exec("INSERT INTO team (name) SELECT tmp.name FROM team AS tmp WHERE NOT EXISTS (SELECT name FROM team WHERE name = ?)", report.TeamName); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add team_name to DB")
	}
	var teamId int
	if err = s.DB.Get(&teamId, "SELECT id FROM team WHERE name = ?", report.TeamName); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to select from DB")
	}
	if _, err = s.DB.Exec("INSERT INTO score_log (team_id, score, message) VALUES (?, ?, ?)", teamId, report.Score, report.Message); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add report to DB")
	}
	return nil
}

func dbconfig() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
}

func main() {
	//db init
	db, err := sqlx.Connect("mysql", dbconfig())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Echo
	s := State{DB: db}
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	// report token
	// work:  curl -v localhost:8080 -H "ReportToken: AAA"
	reportToken := os.Getenv("REPORT_TOKEN")
	g := e.Group("/report")
	g.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:ReportToken",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == reportToken, nil
		},
	}))

	// Template
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t

	// Register Routes
	e.GET("/", s.getBoard)
	e.POST("/report", s.reportScore)
	e.POST("/bench", nil)
	e.Logger.Fatal(e.Start(":8080"))
}

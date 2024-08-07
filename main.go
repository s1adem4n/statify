package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"statify/frontend"
	_ "statify/migrations"
	"statify/pkg/analyze"
	"statify/pkg/scripts"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func SPAMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			if he, ok := err.(*echo.HTTPError); ok {
				if he.Code == http.StatusNotFound {
					c.Request().URL.Path = "/"
					return next(c)
				}
			}
		}
		return err
	}
}

func GetAnalyzeViews(dao *daos.Dao, c echo.Context) ([]analyze.View, error) {
	domain := c.QueryParam("domain")
	start, err := time.Parse(time.RFC3339, c.QueryParam("start"))
	if err != nil {
		return nil, c.String(400, "invalid start time: "+err.Error())
	}
	end, err := time.Parse(time.RFC3339, c.QueryParam("end"))
	if err != nil {
		end = time.Now()
	}
	dedupe := c.QueryParam("dedupe") == "true"

	records, err := dao.FindRecordsByFilter("views", "domain = {:domain} && created >= {:start} && created <= {:end}", "created", 0, 0, dbx.Params{
		"domain": domain,
		"start":  start,
		"end":    end,
	})
	if err != nil {
		return nil, c.String(500, "error fetching views")
	}

	views := make([]analyze.View, len(records))
	for i, record := range records {
		views[i].FromRecord(record)
	}
	if dedupe {
		views = analyze.RemoveDuplicates(views)
	}

	return views, nil
}

func main() {
	app := pocketbase.New()

	var serverAddress string
	app.RootCmd.PersistentFlags().StringVar(&serverAddress, "address", "http://localhost:8080", "The address of the server")

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		router := e.Router
		dao := app.Dao()

		// gzip middleware
		router.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: func(c echo.Context) bool {
				path := c.Request().URL.Path
				return strings.HasPrefix(path, "/_")
			},
		}))
		router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
		}))
		router.Pre(SPAMiddleware)

		subFS := echo.MustSubFS(frontend.Assets, "build")
		router.StaticFS("/", subFS)

		router.GET("/ip", func(c echo.Context) error {
			return c.String(200, c.RealIP())
		})

		router.GET("/tracker.js", func(c echo.Context) error {
			if c.QueryParam("hash") != scripts.TrackerHash {
				return c.Redirect(302, "/tracker.js?hash="+scripts.TrackerHash)
			}

			script := scripts.RenderTracker(serverAddress)
			return c.Blob(200, "application/javascript", []byte(script))
		})

		router.GET("/stats/views/count", func(c echo.Context) error {
			views, err := GetAnalyzeViews(dao, c)
			if err != nil {
				return err
			}

			return c.JSON(200, len(views))
		})
		router.GET("/stats/views/paths", func(c echo.Context) error {
			views, err := GetAnalyzeViews(dao, c)
			if err != nil || views == nil {
				return err
			}

			paths := analyze.CountViewsByPath(views)
			return c.JSON(200, paths)
		})
		router.GET("/stats/views/devices", func(c echo.Context) error {
			views, err := GetAnalyzeViews(dao, c)
			if err != nil || views == nil {
				return err
			}

			devices := analyze.CountViewsByDevice(views)
			stringMap := make(map[string]int)
			stringMap["desktop"] = devices[analyze.DeviceDesktop]
			stringMap["mobile"] = devices[analyze.DeviceMobile]

			return c.JSON(200, stringMap)
		})
		router.GET("/stats/views/sessions", func(c echo.Context) error {
			views, err := GetAnalyzeViews(dao, c)
			if err != nil || views == nil {
				return err
			}

			sessions := analyze.CountViewsBySession(views)
			return c.JSON(200, sessions)
		})
		router.GET("/stats/views/time", func(c echo.Context) error {
			views, err := GetAnalyzeViews(dao, c)
			if err != nil || views == nil {
				return err
			}

			interval, err := time.ParseDuration(c.QueryParam("interval"))
			if err != nil {
				return c.String(400, "invalid interval: "+err.Error())
			}

			counts := analyze.CountViewsOverTime(views, interval)
			return c.JSON(200, counts)
		})
		router.GET("/stats/views/domains", func(c echo.Context) error {
			records, err := dao.FindRecordsByFilter("views", "domain != \"\"", "", 0, 0, nil)
			if err != nil {
				return c.String(500, "error fetching views")
			}

			views := make([]analyze.View, len(records))
			for i, record := range records {
				views[i].FromRecord(record)
			}

			domains := analyze.CountViewsByDomain(views)
			return c.JSON(200, domains)
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

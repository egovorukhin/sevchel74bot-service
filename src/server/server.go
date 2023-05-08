package server

import (
	"bytes"
	"fmt"
	info "github.com/egovorukhin/egoappinfo"
	"github.com/egovorukhin/egolog"
	"github.com/ewa-go/ewa"
	ef "github.com/ewa-go/ewa-fiber"
	"github.com/ewa-go/ewa/security"
	"github.com/ewa-go/ewa/session"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"regexp"
	mApp "sevchel74botService/models/app"
	"time"
)

type Config struct {
	Root        string       `yaml:"root,omitempty"`
	Addr        string       `yaml:"addr"`
	Timeout     Timeout      `yaml:"timeout"`
	Secure      bool         `yaml:"secure"`
	Certificate *Certificate `yaml:"certificate,omitempty"`
	Logger      *Logger      `yaml:"logger,omitempty"`
}

type Timeout struct {
	Read  int `yaml:"read"`
	Write int `yaml:"write"`
	Idle  int `yaml:"idle"`
}

type Certificate struct {
	Cert       string `yaml:"cert"`
	Key        string `yaml:"key"`
	ClientCert string `yaml:"clientCert"`
}

type Logger struct {
	Format   string `yaml:"format"`
	Filename string `yaml:"filename"`
	Time     Time   `yaml:"time"`
}

type Time struct {
	Format   string `yaml:"time_format"`
	Zone     string `yaml:"time_zone"`
	Interval int    `yaml:"time_interval"`
}

func Init(cfg Config) error {

	wsCfg := ewa.Config{
		Addr: cfg.Addr,
		Views: &ewa.Views{
			Root:   cfg.Root,
			Engine: ef.Html,
		},
		Static: &ewa.Static{
			Prefix: "/",
			Root:   cfg.Root,
		},
		Session: &session.Config{
			RedirectPath: "/login",
			//Expires:              time.Second * 30,
			SessionHandler:       sessionHandler,
			DeleteSessionHandler: deleteSessionHandler,
		},
		Authorization: security.Authorization{
			ApiKey: &security.ApiKey{
				KeyName: "Token",
				Param:   security.ParamHeader,
				Handler: apiKeyHandler,
			},
		},
		ContextHandler: contextHandler,
	}
	if cfg.Secure && cfg.Certificate != nil {
		wsCfg.Secure = &ewa.Secure{
			Key:        cfg.Certificate.Key,
			Cert:       cfg.Certificate.Cert,
			ClientCert: cfg.Certificate.ClientCert,
		}
	}
	app := fiber.New(fiber.Config{
		Views: ef.NewViews(cfg.Root, ef.Html, &ef.Engine{
			Reload: true,
		}),
		ReadTimeout:  time.Duration(cfg.Timeout.Read) * time.Second,
		WriteTimeout: time.Duration(cfg.Timeout.Write) * time.Second,
		IdleTimeout:  time.Duration(cfg.Timeout.Idle) * time.Second,
	})
	// CORS
	app.Use(cors.New())
	// Logger
	if cfg.Logger != nil {
		app.Use(logger.New(logger.Config{
			Format:       cfg.Logger.Format,
			TimeFormat:   cfg.Logger.Time.Format,
			TimeZone:     cfg.Logger.Time.Zone,
			TimeInterval: time.Duration(cfg.Logger.Time.Interval) * time.Millisecond,
			Output:       cfg.Logger,
		}))
	}
	// Новый сервер
	server := ewa.New(&ef.Server{App: app}, wsCfg)
	//api
	/*server.Register(new(controllers.Api)).NotShow()
	//web
	server.Register(new(web.Home)).SetPath("/").NotShow()
	server.Register(new(web.About)).SetPath("/about").NotShow()
	server.Register(new(web.Servers)).SetPath("/servers").NotShow()
	server.Register(new(web.Users)).SetPath("/users").NotShow()
	server.Register(new(web.Roles)).SetPath("/roles").NotShow()
	server.Register(new(web.Commands)).SetPath("/commands").NotShow()
	server.Register(new(web.Params)).SetPath("/params").NotShow()
	server.Register(new(web.Settings)).SetPath("/settings").NotShow()
	server.Register(new(web.Login)).SetPath("/login").NotShow()
	server.Register(new(web.Logout)).SetPath("/logout").NotShow()

	//api
	server.Register(new(controllers.Auth))
	//api.adbridge
	server.Register(new(bridge.Server))
	server.Register(new(bridge.AdMember))
	server.Register(new(bridge.Command))
	server.Register(new(bridge.Param))
	server.Register(new(bridge.Role))
	server.Register(new(bridge.Config))
	server.Register(new(bridge.Client))
	server.Register(new(bridge.ServerMemberRole))
	server.Register(new(bridge.ParamCommand))
	server.Register(new(bridge.RoleCommand))
	//api.ad
	server.Register(new(ad.User))
	//api.settings
	server.Register(new(settings.Audit))

	// Модели для свагера
	server.Swagger.SetModels(ewa.Models{
		mAdbridge.ModelAdMember:    mAdbridge.AdMember{},
		mSettings.ModelAudit:       mSettings.Audit{},
		model.ModelCrudResponse:    model.CrudResponse{},
		bridge.ModelServerResponse: bridge.ServerResponse{},
		bridge.ModelServer:         bridge.Server{},
	})*/

	// Описываем swagger
	server.Swagger.SetBasePath("/api").SetInfo(info.GetApplication().Hostname, &ewa.Info{
		Description: "Bridge - very nice application!",
		Version:     info.GetVersion().String(),
		Title:       info.GetApplicationName(),
		Contact: &ewa.Contact{
			Email: "yegor.govorukhin@kaspi.kz",
		},
		License: &ewa.License{
			Name: "Freeware license to Kaspi Bank",
		},
	}, nil)

	//Запускаем web сервер
	return server.Start()
}

// Обработчик для fiber контекста
func contextHandler(handler ewa.Handler) interface{} {
	return func(ctx *fiber.Ctx) error {
		return handler(ewa.NewContext(&ef.Context{Ctx: ctx}))
	}
}

func apiKeyHandler(token string) (username string, err error) {
	return sessionHandler(token)
}

func sessionHandler(id string) (user string, err error) {
	s, err := mApp.GetSession(fmt.Sprintf("id='%s'", id))
	if err != nil {
		return "", err
	}

	u, err := mApp.GetUser(fmt.Sprintf("id=%d", *s.UserId))
	if err != nil {
		return "", err
	}

	now := time.Now()
	s.LastTime = &now
	err = s.Update()
	if err != nil {
		return "", err
	}

	return u.Login, nil
}

func deleteSessionHandler(id string) bool {
	err := mApp.DeleteSession(id)
	if err != nil {
		return false
	}
	return true
}

func (l *Logger) Write(data []byte) (n int, err error) {
	egolog.Infofn(l.Filename, string(maskPwd(data, []byte("password"), []byte("******"))))
	return len(data), nil
}

func maskPwd(data, wordPwd, mask []byte) []byte {
	if bytes.Contains(data, wordPwd) {
		index := -1
		data = regexp.MustCompile(`[^"\\]+(?:\\.[^"\\]*)*`).ReplaceAllFunc(data, func(b []byte) []byte {
			if bytes.Contains(b, wordPwd) {
				index = 0
			}
			if index > -1 {
				index++
			}
			if index == 3 {
				return mask
			}
			return b
		})
	}
	return data
}

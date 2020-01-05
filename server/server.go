package server

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
)

type Server struct {
	e    *echo.Echo
	port string
}

func New(port string) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Logger.SetOutput(ioutil.Discard)
	return &Server{e: e, port: port}
}

func (s *Server) AddFileHandlerFromPath(rootPath string) error {
	return s.AddFileHandlerFromFs(http.Dir(rootPath))
}

func (s *Server) AddFileHandlerFromFs(fs http.FileSystem) error {
	assetHandler := http.FileServer(fs)
	s.e.GET("/*", echo.WrapHandler(assetHandler))
	return nil
}

func (s *Server) Start() error {
	return s.e.Start(":" + s.port)
}

func Api() error {
	return nil
}

package lib

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/zserge/lorca"

	"github.com/labstack/echo/v4"
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

func (s *Server) StartWithApp() error {
	go func() {
		s.Start()
	}()

	// FIXME convert widht/height to variable
	ui, _ := lorca.New("http://localhost:"+s.port, "", 720, 480)
	defer ui.Close()

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return s.Shutdown(ctx)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func (s *Server) AddApiHandler() error {
	handler := func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer func() {
			err := src.Close()
			if err != nil {
				log.Println("failed to close uploaded file")
			}
		}()

		// Destination
		dst, err := os.Create(file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully.</p>", file.Filename))
	}
	s.e.POST("/api/upload", handler)
	return nil
}

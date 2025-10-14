package server

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/mbilarusdev/durak_auth_bot/internal/endpoint"
)

type HttpServer struct {
	SendCodeEndpoint    *endpoint.SendCodeEndpoint
	ConfirmCodeEndpoint *endpoint.ConfirmCodeEndpoint
	CheckAuthEndpoint   *endpoint.CheckAuthEndpoint
	LogoutEndpoint      *endpoint.LogoutEndpoint
}

func NewHttpServer(sendCodeEndpoint *endpoint.SendCodeEndpoint,
	confirmCodeEndpoint *endpoint.ConfirmCodeEndpoint,
	checkAuthEndpoint *endpoint.CheckAuthEndpoint,
	logoutEndpoint *endpoint.LogoutEndpoint) *HttpServer {
	server := new(HttpServer)
	server.SendCodeEndpoint = sendCodeEndpoint
	server.ConfirmCodeEndpoint = confirmCodeEndpoint
	server.CheckAuthEndpoint = checkAuthEndpoint
	server.LogoutEndpoint = logoutEndpoint
	return server
}

func (httpServer *HttpServer) ListenAndServe() {
	addr := ":8080"
	r := mux.NewRouter()
	swaggerDir := filepath.Join("docs", "swagger-ui")
	r.PathPrefix("/swagger/").
		Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir(swaggerDir))))
	r.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})
	r.HandleFunc("/code/send", httpServer.SendCodeEndpoint.Call).Methods(http.MethodPost)
	r.HandleFunc("/code/confirm", httpServer.ConfirmCodeEndpoint.Call).Methods(http.MethodPost)
	r.HandleFunc("/login/check", httpServer.CheckAuthEndpoint.Call).Methods(http.MethodPost)
	r.HandleFunc("/logout", httpServer.LogoutEndpoint.Call).Methods(http.MethodPost)
	server := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Starting HTTP server at %s", addr)
	log.Fatal(server.ListenAndServe())
}

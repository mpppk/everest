package serve

import "net/http"

func Files(rootPath, port string) error {
	http.Handle("/", http.FileServer(http.Dir(rootPath)))
	return http.ListenAndServe(":"+port, nil)
}

func FileSystem(fs http.FileSystem, port string) error {
	http.Handle("/", http.FileServer(fs))
	return http.ListenAndServe(":"+port, nil)
}

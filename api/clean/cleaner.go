package clean

import "net/http"

func CleanDatabase(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Clean Db"))
}
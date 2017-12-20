package server

import (
  	"net/http"
)

func Start() {
	router := newRouter()
	go router.run()
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		handleUpgradeRequest(router, w, r)
	})
	http.ListenAndServe(":3000", nil)
}

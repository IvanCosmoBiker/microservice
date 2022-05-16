package listener

import(
    "net/http"
    "log"
)

func StartListener(url string ,listenPort string, f func(w http.ResponseWriter, req *http.Request)){
    http.HandleFunc(url, f)
    log.Println(listenPort)
	if err := http.ListenAndServe(listenPort, nil); err != nil {
		log.Println(err)
	}
}
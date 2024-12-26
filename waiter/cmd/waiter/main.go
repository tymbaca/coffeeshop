package waiter

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OKay"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

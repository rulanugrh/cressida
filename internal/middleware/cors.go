package middleware

import "net/http"

func CORSHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "POST, PUT, GET, OPTIONS, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Conteng-Lenght")
		w.Header().Add("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CORSStreaming(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "POST, PUT, GET, OPTIONS, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Conteng-Lenght")
		w.Header().Add("Content-Type", "text/event-stream")

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		next.ServeHTTP(w, r)
	})
}

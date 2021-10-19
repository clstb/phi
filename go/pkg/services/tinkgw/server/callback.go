package server

import "net/http"

func (s *Server) Callback(callbacks <-chan PendingCallback) http.HandlerFunc {
	m := make(map[string]chan<- string)
	go func() {
		for {
			cb := <-callbacks
			if cb.del {
				delete(m, cb.state)
			} else {
				m[cb.state] = cb.ch
			}
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		state := q.Get("state")
		if state == "" {
			http.Error(w, "bad request: missing state", http.StatusBadRequest)
			return
		}

		ch, ok := m[state]
		if !ok {
			http.Error(w, "precondition failed: state not found", http.StatusPreconditionFailed)
			return
		}
		delete(m, state)

		ch <- ""

		_, err := w.Write([]byte("ok"))
		if err != nil {
			http.Error(w, "internal server error: writing reponse", http.StatusInternalServerError)
			return
		}
	}
}

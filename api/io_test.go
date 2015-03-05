package api

import (
	"fmt"
	"github.com/franela/goblin"
	o "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewFakeAPIServer() *httptest.Server {
	return httptest.NewTLSServer(http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {

		method := r.Method
		path := r.URL.Path

		r.ParseForm()

		switch path {
		case "/0/method":
			w.WriteHeader(200)
			fmt.Fprint(w, fmt.Sprintf("%v %v", method, path))

		case "/0/idontexist":
			w.WriteHeader(404)
			fmt.Fprint(w, "nope")

		case "/0/geturlparams":
			if method == "GET" {
				w.WriteHeader(200)
				fmt.Fprint(w, fmt.Sprintf("%v %s", method, r.Form.Encode()))
			}

		case "/0/getpostparams":
			if method == "POST" {
				w.WriteHeader(200)
				fmt.Fprint(w, fmt.Sprintf("%v %s", method, r.PostForm.Encode()))
			}

			// API URLs. We don't mock the API here since we're only testing
			// the low-level stuff. The parsing tests will be done in
			// client_test.

		case // GET
			"/0/api",
			"/0/create",
			"/0/destroy",
			"/0/games",
			"/0/join",
			"/0/log",
			"/0/logout",
			"/0/play",
			"/0/shutdown",
			"/0/status",
			"/0/whoami",
			// POST
			"/0/register",
			"/0/auth":

			w.WriteHeader(200)
			fmt.Fprint(w, fmt.Sprintf("%v %s %s", method, path, r.Form.Encode()))
		}
	}))
}

func TestIO(t *testing.T) {

	g := goblin.Goblin(t)

	o.RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("NewHTTClient", func() {
		g.It("Should not return nil", func() {
			o.Expect(NewHTTClient()).NotTo(o.BeNil())
		})
	})

	g.Describe("getError", func() {
		g.It("Should return an Err4XX if the code is 4XX", func() {
			o.Expect(getError(400)).To(o.Equal(Err4XX))
			o.Expect(getError(403)).To(o.Equal(Err4XX))
			o.Expect(getError(404)).To(o.Equal(Err4XX))
		})

		g.It("Should return an Err5XX if the code is 5XX", func() {
			o.Expect(getError(500)).To(o.Equal(Err5XX))
		})
	})

	g.Describe("Httclient", func() {
		g.Describe(".call", func() {
			var ts *httptest.Server
			var h *Httclient

			g.Before(func() { ts = NewFakeAPIServer() })
			g.After(func() { ts.Close() })

			g.BeforeEach(func() {
				h = NewHTTClient()
				h.baseURL = ts.URL
			})

			g.It("Should call the remote server", func() {
				b := h.call(get, "/method", nil)

				o.Expect(b).NotTo(o.BeNil())
				o.Expect(b.Error()).To(o.BeNil())
				o.Expect(b.IsEmpty()).To(o.BeFalse())
				o.Expect(b.StatusCode).To(o.Equal(200))
				o.Expect(b.Content.ToString()).To(o.Equal("GET /0/method"))
			})

			g.It("Should use POST if it was given as the method", func() {
				b := h.call(post, "/method", nil)

				o.Expect(b).NotTo(o.BeNil())
				o.Expect(b.Error()).To(o.BeNil())
				o.Expect(b.IsEmpty()).To(o.BeFalse())
				o.Expect(b.StatusCode).To(o.Equal(200))
				o.Expect(b.Content.ToString()).To(o.Equal("POST /0/method"))
			})

			g.It("Should set body.err to Err4XX if the status code is 4XX", func() {
				b := h.call(get, "/idontexist", nil)

				o.Expect(b).NotTo(o.BeNil())
				o.Expect(b.Error()).To(o.Equal(Err4XX))
				o.Expect(b.StatusCode).To(o.Equal(404))
			})

			g.It("Should send parameters in the URL for GET requests", func() {
				b := h.call(get, "/geturlparams", struct {
					Param string
				}{"foo"})

				o.Expect(b).NotTo(o.BeNil())
				o.Expect(b.Error()).To(o.BeNil())
				o.Expect(b.IsEmpty()).To(o.BeFalse())
				o.Expect(b.StatusCode).To(o.Equal(200))
				o.Expect(b.Content.ToString()).To(o.Equal("GET param=foo"))
			})

			g.It("Should send parameters in the body for POST requests", func() {
				b := h.call(post, "/getpostparams", struct {
					Param string
				}{"foo"})

				o.Expect(b).NotTo(o.BeNil())
				o.Expect(b.Error()).To(o.BeNil())
				o.Expect(b.IsEmpty()).To(o.BeFalse())
				o.Expect(b.StatusCode).To(o.Equal(200))
				o.Expect(b.Content.ToString()).To(o.Equal("POST param=foo"))
			})
		})

		g.Describe("API calls", func() {
			var ts *httptest.Server
			var h *Httclient

			g.Before(func() { ts = NewFakeAPIServer() })
			g.After(func() { ts.Close() })

			g.BeforeEach(func() {
				h = NewHTTClient()
				h.baseURL = ts.URL
			})

			g.Describe("CallAPI", func() {
				g.It("Should GET /api without parameters", func() {
					b := h.CallAPI()

					o.Expect(b).NotTo(o.BeNil())
					o.Expect(b.Error()).To(o.BeNil())
					o.Expect(b.StatusCode).To(o.Equal(200))
					o.Expect(b.Content.ToString()).To(o.Equal("GET /0/api "))
				})
			})

			g.Describe("CallAuth", func() {
				g.It("Should POST /auth with a login and a password", func() {
					b := h.CallAuth(UserCredentialsParams{
						Login:    "foo",
						Password: "bar",
					})

					o.Expect(b).NotTo(o.BeNil())
					o.Expect(b.Error()).To(o.BeNil())
					o.Expect(b.StatusCode).To(o.Equal(200))
					o.Expect(b.Content.ToString()).To(o.Equal(
						"POST /0/auth login=foo&password=bar"))
				})
			})

			g.Describe("CallCreate", func() {
				g.It("Should GET /create with game parameters", func() {
					b := h.CallCreate(GameSpecParams{
						Users:             "foo",
						Teaser:            "desc",
						Pace:              12,
						Nb_turn:           34,
						Nb_ant_per_player: 2,
						Nb_player:         1,
						Minimal_nb_player: 1,
						Initial_energy:    40,
						Initial_acid:      35,
					})

					o.Expect(b).NotTo(o.BeNil())
					o.Expect(b.Error()).To(o.BeNil())
					o.Expect(b.StatusCode).To(o.Equal(200))
					o.Expect(b.Content.ToString()).To(o.Equal(
						"GET /0/create " +
							"initial_acid=35&initial_energy=40" +
							"&minimal_nb_player=1&nb_ant_per_player=2" +
							"&nb_player=1&nb_turn=34&pace=12&teaser=desc" +
							"&users=foo"))
				})
			})

			g.Describe("CallDestroy", func() {
				g.It("Should GET /destroy with a game ID", func() {
					b := h.CallDestroy(GameIDParams{Id: GameID("foobar")})

					o.Expect(b).NotTo(o.BeNil())
					o.Expect(b.Error()).To(o.BeNil())
					o.Expect(b.StatusCode).To(o.Equal(200))
					o.Expect(b.Content.ToString()).To(o.Equal(
						"GET /0/destroy id=foobar"))
				})
			})

			g.Describe("CallGames", func() {
				g.It("Should GET /games without parameters", func() {
					b := h.CallGames()

					o.Expect(b).NotTo(o.BeNil())
					o.Expect(b.Error()).To(o.BeNil())
					o.Expect(b.StatusCode).To(o.Equal(200))
					o.Expect(b.Content.ToString()).To(o.Equal("GET /0/games "))
				})
			})

			g.Describe("CallJoin", func() {
				g.It("Should GET /join with a game ID", func() {
					b := h.CallJoin(GameIDParams{Id: GameID("foobar1")})

					o.Expect(b).NotTo(o.BeNil())
					o.Expect(b.Error()).To(o.BeNil())
					o.Expect(b.StatusCode).To(o.Equal(200))
					o.Expect(b.Content.ToString()).To(o.Equal(
						"GET /0/join id=foobar1"))
				})
			})
			g.Describe("CallLog", func() {
				g.It("Should GET /log with a game ID", func() {
					b := h.CallLog(GameIDParams{Id: GameID("foobar42")})

					o.Expect(b).NotTo(o.BeNil())
					o.Expect(b.Error()).To(o.BeNil())
					o.Expect(b.StatusCode).To(o.Equal(200))
					o.Expect(b.Content.ToString()).To(o.Equal(
						"GET /0/log id=foobar42"))
				})
			})

			g.Describe("CallLogout", func() {
				g.It("Should GET /logout without parameters", func() {
					b := h.CallLogout()

					o.Expect(b).NotTo(o.BeNil())
					o.Expect(b.Error()).To(o.BeNil())
					o.Expect(b.StatusCode).To(o.Equal(200))
					o.Expect(b.Content.ToString()).To(o.Equal("GET /0/logout "))
				})
			})

			g.Describe("CallPlay", func() {
				g.It("Should GET /play with a game id and some commands", func() {
					b := h.CallPlay(PlayParams{Id: GameID("a"), Cmds: "xyz"})

					o.Expect(b).NotTo(o.BeNil())
					o.Expect(b.Error()).To(o.BeNil())
					o.Expect(b.StatusCode).To(o.Equal(200))
					o.Expect(b.Content.ToString()).To(o.Equal(
						"GET /0/play cmds=xyz&id=a"))
				})
			})

			g.Describe("CallRegister", func() {
				g.It("Should POST /register with a login and a password", func() {
					b := h.CallRegister(UserCredentialsParams{
						Login:    "foo1",
						Password: "bar",
					})

					o.Expect(b).NotTo(o.BeNil())
					o.Expect(b.Error()).To(o.BeNil())
					o.Expect(b.StatusCode).To(o.Equal(200))
					o.Expect(b.Content.ToString()).To(o.Equal(
						"POST /0/register login=foo1&password=bar"))
				})
			})

			g.Describe("CallShutdown", func() {
				g.It("Should GET /shutdown with an ID", func() {
					b := h.CallShutdown(GenericIDParams{Id: "abc"})

					o.Expect(b).NotTo(o.BeNil())
					o.Expect(b.Error()).To(o.BeNil())
					o.Expect(b.StatusCode).To(o.Equal(200))
					o.Expect(b.Content.ToString()).To(o.Equal(
						"GET /0/shutdown id=abc"))
				})
			})

			g.Describe("CallStatus", func() {
				g.It("Should GET /status with a game ID", func() {
					b := h.CallStatus(GameIDParams{Id: GameID("f42")})

					o.Expect(b).NotTo(o.BeNil())
					o.Expect(b.Error()).To(o.BeNil())
					o.Expect(b.StatusCode).To(o.Equal(200))
					o.Expect(b.Content.ToString()).To(o.Equal(
						"GET /0/status id=f42"))
				})
			})

			g.Describe("CallWhoAmI", func() {
				g.It("Should GET /whoami without parameters", func() {
					b := h.CallWhoAmI()

					o.Expect(b).NotTo(o.BeNil())
					o.Expect(b.Error()).To(o.BeNil())
					o.Expect(b.StatusCode).To(o.Equal(200))
					o.Expect(b.Content.ToString()).To(o.Equal("GET /0/whoami "))
				})
			})
		})
	})
}

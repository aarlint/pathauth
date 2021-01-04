// Package plugindemo a PathAuth plugin.
package pathauth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	AccessLogger  *log.Logger
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	file, err := os.OpenFile("/var/log/pathauth.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	AccessLogger = log.New(file, "ACCESS: ", log.Ldate|log.Ltime)
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
}

// Config the plugin configuration.
type Config struct {
	Groups map[string][]string `json:"groups,omitempty"`
	Paths  []Path              `json:"paths,omitempty"`
}

// Path holds one path vs user check
type Path struct {
	Regex  string   `json:"path,omitempty"`
	Users  []string `json:"users,omitempty"`
	Groups []string `json:"groups,omitempty"`
	Public bool     `json:"public,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// PathAuth a PathAuth plugin.
type PathAuth struct {
	next        http.Handler
	publicPaths []*regexp.Regexp
	users       map[string][]*regexp.Regexp
	name        string
}

// New created a new PathAuth plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	publicPaths := []*regexp.Regexp{}
	users := make(map[string][]*regexp.Regexp)
	// create map of user and regexs
	for _, path := range config.Paths {
		re := regexp.MustCompile(path.Regex)
		if path.Public {
			publicPaths = append(publicPaths, re)
		}
		u := make(map[string]bool)
		for _, group := range path.Groups {
			for _, user := range config.Groups[group] {
				u[user] = true
			}
		}
		for _, user := range path.Users {
			u[user] = true
		}
		fmt.Println(u)
		for user, _ := range u {
			users[user] = append(users[user], re)
		}
	}
	// fmt.Println(publicPaths, "\n", users)
	return &PathAuth{
		publicPaths: publicPaths,
		users:       users,
		next:        next,
		name:        name,
	}, nil
}

func (a *PathAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	user := ""
	// Read cookie
	cookie, cerr := req.Cookie("_forward_auth")
	if cerr != nil {
		fmt.Printf("No _forward_auth cookie\r\n")
	} else {
		user = strings.Split(cookie.Value, "|")[2]
	}
	basicAuthUser, _, ok := req.BasicAuth()
	if !ok {
		fmt.Printf("No BasicAuth\r\n")
	} else {
		user = basicAuthUser
	}
	if user == "" {
		WarningLogger.Println("No user found in header _forward_auth or basic auth")
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	// loop regex's allowed by user
	realIp := req.Header.Get("X-Real-Ip")
	requestedURL := req.Header.Get("X-Forwarded-Host")
	requestedPath := req.URL.EscapedPath()
	for _, regex := range a.users[user] {
		if regex.MatchString(requestedPath) {
			fmt.Printf(user)
			req.Header.Add("User", user)
			req.Header.Set("User", user)
			AccessLogger.Println(user, "accessed", requestedURL+requestedPath, "from", realIp)
			a.next.ServeHTTP(rw, req)
			return
		}
	}

	// if not authenticated check public paths

	for _, regex := range a.publicPaths {
		if regex.MatchString(requestedPath) {
			fmt.Printf(user)
			req.Header.Add("User", user)
			req.Header.Set("User", user)
			a.next.ServeHTTP(rw, req)
			AccessLogger.Println(user, "accessed", "public path", requestedURL+requestedPath, "from", realIp)
			return
		}
	}

	// if still not authed then return status forbidden
	AccessLogger.Println(user, "was denied access to", requestedURL+requestedPath)
	rw.WriteHeader(http.StatusForbidden)
	return

}

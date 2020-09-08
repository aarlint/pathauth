// Package plugindemo a PathAuth plugin.
package pathauth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"regexp"
	
)

// Config the plugin configuration.
type Config struct {
	Base string `json:"base,omitempty"`
	Paths []Path `json:"paths,omitempty"`
}

// Path holds one path vs user check
type Path struct {
	Regex string `json:"path,omitempty"`
	Users []string `json:"users,omitempty"`
}



// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// PathAuth a PathAuth plugin.
type PathAuth struct {
	next     http.Handler
	base	 string
	paths    []Path
	name     string
}

// New created a new PathAuth plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	

	return &PathAuth{
		base: config.Base,
		paths: config.Paths,
		next:     next,
		name:     name,
	}, nil
}

func (a *PathAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	authenticated := false
	// Read cookie
	cookie, err := req.Cookie("_forward_auth")
	if err != nil {
		fmt.Printf("Cant find cookie :/\r\n")
		return
	}
	user := strings.Split(cookie.Value,"|")[2]
	currentPath := req.URL.EscapedPath()
	fmt.Println(user)
	for _, path := range a.paths {
		fmt.Println("looking for ", path.Regex)
		re := regexp.MustCompile(path.Regex)
		if re.MatchString(currentPath) {
			fmt.Println("current path is ", path.Regex)
			for _, allowedUser := range path.Users {
				fmt.Println("looking for ", allowedUser)
				if user == allowedUser {
					fmt.Println(user, " You're in!")
					authenticated = true
					break
				}
			}
		}
	}
	if !authenticated {
		rw.WriteHeader(http.StatusForbidden)
	}
	a.next.ServeHTTP(rw, req)
}

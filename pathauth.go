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
	Groups map[string][]string `json:"groups,omitempty"`
	Paths []Path `json:"paths,omitempty"`
}


// Path holds one path vs user check
type Path struct {
	Regex string `json:"path,omitempty"`
	Users []string `json:"users,omitempty"`
	Groups []string `json:"groups,omitempty"`
	Public bool `json:"public,omitempty"`
}



// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// PathAuth a PathAuth plugin.
type PathAuth struct {
	next     http.Handler
	publicPaths []string
	users	 map[string][]string
	name     string
}

func find(source []string, value string) bool {
    for _, item := range source {
        if item == value {
			fmt.Println(item)
            return true
        }
    }
    return false
}

// New created a new PathAuth plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	fmt.Println(config.Groups)
	publicPaths := []string{}
	users := make(map[string][]string)
	// create map of user and regexs
	for _, path := range config.Paths {
		if path.Public {
			publicPaths = append(publicPaths, path.Regex)
		}
		u := []string{}
		for _, group := range path.Groups {
			fmt.Println(config.Groups[group])
			for _, user := range config.Groups[group] {
				if find(u,user) == false {
					u = append(u, user)
				}
			}
		}
		for _, user := range path.Users {
			if find(u,user) == false {
					u = append(u, user)
				}
			}
		for _, user := range u {
			users[user] = append(users[user], path.Regex)
		}	
	}
	fmt.Println(users)
	return &PathAuth{
		publicPaths: publicPaths,
		users: users,
		next:     next,
		name:     name,
	}, nil
}

func (a *PathAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	authenticated := false
	user := ""
	// Read cookie
	cookie, cerr := req.Cookie("_forward_auth")
	if cerr != nil {
		fmt.Printf("Can't find forward auth cookie :/\r\n")
	} else {
		user = strings.Split(cookie.Value, "|")[2]
	}
	basicAuthUser, _, ok := req.BasicAuth()
	if !ok {
		fmt.Printf("Can't find BasicAuth :/\r\n")
	} else {
		user = basicAuthUser
	}
	if user == "" {
		fmt.Printf("cannot determine current user :/\r\n")
		rw.WriteHeader(http.StatusForbidden)
			return
	}
	currentPath := req.URL.EscapedPath()
	for _, path := range a.users[user] {
		fmt.Println("looking for ", path)
		re := regexp.MustCompile(path)
		if re.MatchString(currentPath) {
			authenticated = true
			break
				}
			}
	// if not authenticated check public paths
	if !authenticated {
		for _, path := range a.publicPaths {
			re := regexp.MustCompile(path)
			if re.MatchString(currentPath) {
				authenticated = true
				break
					}
		}
		if !authenticated {
			rw.WriteHeader(http.StatusForbidden)
			return

		}
	}
	a.next.ServeHTTP(rw, req)
}

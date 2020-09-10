# Path Auth

Path auth is a middleware plugin for [Traefik](https://github.com/containous/traefik) which enables authorization on the path when chained behind [traefik-forward-auth](https://github.com/thomseddon/traefik-forward-auth) or basic auth, making it possible to place simple user authorization based on regex and userlists defined in middleware.

## Configuration

### Static

```toml
[experimental.pilot]
  token = "xxxx"

[experimental.plugins.pathauth]
  modulename = "github.com/aarlint/pathauth"
  version = "v0.1.0"
```

### Dynamic

To configure the `Path Auth` plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in 
your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/). The following example creates
and uses the `pathauth` middleware plugin to ensure the user's email is allowed to visit the regexed path.

In order to use this middleware you must chain it behind [traefik-forward-auth](https://github.com/thomseddon/traefik-forward-auth).


```yaml
http:
  # Add the router
  routers:
    my-router:
      entryPoints:
      - http
      middlewares:
      - pathauth
      service: service-foo
      rule: Path(`/foo`)

  # Add the middleware
  middlewares:
    pathauth:
      plugin:
        groups:
        #group key names can be anything you want
          admin:
            - professor.professerton@gmail.com
            - austin.arlint@gmail.com
            
        paths:
          - regex: ^/notls$
            groups:
              - admin
          # this path is accessible by anyone who makes it here because public key = true
          - regex: ^/public$
            public: true
          # this path is accessible by anyone listed in the admin group plus new.breath@gmail.com
          - regex: ^/other$
            users: 
              - new.breath@gmail.com
            groups:
              - admin
  # Add the service
  services:
    service-foo:
      loadBalancer:
        servers:
        - url: http://localhost:5000/
        passHostHeader: false
```

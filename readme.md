# Path Auth

Path auth is a middleware plugin for [Traefik](https://github.com/containous/traefik) which enables authorization on the path when chained behind [traefik-forward-auth](https://github.com/thomseddon/traefik-forward-auth), making it possible to place simple user authorization based on regex and userlists defined in middleware.

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
        paths:
          - regex: ^/notls
            users: 
              - austin.arlint@gmail.com
              - new.breath@gmail.com
          - regex: ^/yourmom
            users:
              - test.user@gmail.com
              - other.user@gmail.com

  # Add the service
  services:
    service-foo:
      loadBalancer:
        servers:
        - url: http://localhost:5000/
        passHostHeader: false
```

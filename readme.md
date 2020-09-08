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


```toml
[http.routers]
  [http.routers.my-router]
    rule = "Host(`localhost`)"
    middlewares = ["pathauth-foo"]
    service = "my-service"

[http.middlewares]
  [http.middlewares.pathauth-foo.plugin.pathauth]


    # allows other.user@gmail.com to access path ^/yourmom
    [[http.middlewares.rewrite-foo.plugin.pathauth]]
      paths:
        - regex: ^/notls
          users: 
            - austin.arlint@gmail.com
            - poop.breath@gmail.com
        - regex: ^/yourmom
          users:
            - test.user@gmail.com
            - other.user@gmail.com

[http.services]
  [http.services.my-service]
    [http.services.my-service.loadBalancer]
      [[http.services.my-service.loadBalancer.servers]]
        url = "http://127.0.0.1"
```

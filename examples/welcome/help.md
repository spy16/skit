Hello from Skit! :heart:
Looks like this skit instance is not configured yet.
You can configure skit using the following `skit.toml` file:

```
token = "your-token-here"

[[handlers]]
name = "match_all"
type = "command"
match = [
    ".*"
]
cmd = "./hello.sh"
timeout = "5s"
```

As simple as that! :kungfu:

Visit https://github.com/spy16/skit for more help.

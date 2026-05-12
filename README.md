# RoadRunner / Velox `uax29` Repro

This repo contributes the conflicting `github.com/clipperhouse/uax29/v2@v2.7.0`
selection that breaks a RoadRunner / Velox build when Velox also compiles its
own `github.com/olekukonko/tablewriter@v1.1.3` dependency.

The plugin package itself should still build:

```sh
go test ./...
```

The failure appears in the Velox/RoadRunner build, where Velox contributes:

- `github.com/olekukonko/tablewriter@v1.1.3`
- `github.com/clipperhouse/displaywidth@v0.6.2`

and this plugin contributes:

- `github.com/clipperhouse/uax29/v2@v2.7.0`

Expected Velox failure:

```text
# github.com/clipperhouse/displaywidth
.../github.com/clipperhouse/displaywidth@v0.6.2/graphemes.go:48:12: cannot use graphemes.FromString(s) (value of type *graphemes.Iterator[string]) as graphemes.Iterator[string] value in struct literal
.../github.com/clipperhouse/displaywidth@v0.6.2/graphemes.go:69:12: cannot use graphemes.FromBytes(s) (value of type *graphemes.Iterator[[]byte]) as graphemes.Iterator[[]byte] value in struct literal
```

The `uax29` selection is pulled in transitively, matching `../pcx-frontend-api-go`:

- `plugin.go` imports `gopkg.in/DataDog/dd-trace-go.v1/contrib/gofiber/fiber.v2`
- `gopkg.in/DataDog/dd-trace-go.v1@v1.74.8` selects DataDog v2 modules
- `github.com/DataDog/dd-trace-go/contrib/gofiber/fiber.v2/v2@v2.8.1` selects `github.com/clipperhouse/uax29/v2@v2.7.0`
- `github.com/gofiber/fiber/v2` compiles `github.com/mattn/go-runewidth`, which compiles `github.com/clipperhouse/uax29/v2/graphemes`

To inspect the relevant graph:

```sh
go list -m gopkg.in/DataDog/dd-trace-go.v1 github.com/DataDog/dd-trace-go/contrib/gofiber/fiber.v2/v2 github.com/clipperhouse/uax29/v2
go mod graph | grep 'github.com/clipperhouse/uax29/v2'
go mod graph | grep 'github.com/clipperhouse/displaywidth'
```

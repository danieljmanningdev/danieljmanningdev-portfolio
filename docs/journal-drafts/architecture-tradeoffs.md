# Server-rendered Go: a deliberate default, not a universal rule

My portfolio combines public service pages with a private workspace for clients, projects, contracts and publishing. Most interactions are a form submission followed by validation, a database operation and an updated document. Go renders that document; HTMX enhances selected interactions.

That shape of product is a good reason to start with server rendering. It is not a reason to declare every client-side framework unnecessary.

## Keep decisions close to the data

Middleware handles cross-cutting concerns, handlers interpret an operation, repositories perform parameterised database queries and templates present the result. The value is a clear boundary, not the maximum possible number of packages.

Public headings, canonical URLs, descriptions and JSON-LD arrive in the HTML response. A visitor can follow ordinary links without waiting for a client-side router to initialise. Enhancements should improve that baseline rather than make it fragile.

## The trade-offs still exist

A server-rendered form needs understandable errors and a clear result. A backend written in Go does not repair an unclear label or a page that overflows on a phone. Offline editing also does not appear automatically: it requires local storage, synchronisation and decisions about conflicts.

For a visual editor, canvas-heavy product or complex local state, I would reassess the frontend rather than insist that every interaction use the same approach.

## Make reusable work accurately identifiable

The Go web libraries and `djm-cli` are companion projects. The portfolio still has application-specific packages under `internal/`; showing those libraries on the homepage should not imply that all are runtime dependencies here.

The [Portfolio & Client Workspace case study](https://danieljmanningdev.com/work/portfolio) and [source repository](https://github.com/danieljmanningdev/danieljmanningdev-portfolio) make those choices inspectable. The useful comparison is whether an approach supports the product's requirements, maintenance and user tasks—not whether it wins a stack argument.

Further reading: [Go documentation](https://go.dev/doc/).

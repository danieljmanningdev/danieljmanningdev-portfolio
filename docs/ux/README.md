# UX JSON documentation

This directory is a developer-first map of the portfolio and private workspace.

The JSON files do not reproduce every visual detail from Figma or every implementation detail from Go, HTML and CSS. They record the information that is useful when planning, locating and changing an experience:

- what a screen, component or flow is for;
- where its route, handler, template and layout live;
- which reusable components and flows it belongs to;
- which meaningful states exist;
- which UX decisions affect it;
- whether JavaScript or HTMX is involved.

## Structure

```text
docs/ux/
├── components/   Reusable interface components and patterns
├── definitions/  Shared JSON Schema definitions
├── flows/        Multi-screen user and system journeys
├── schemas/      Schemas for screens, components, flows and optional tokens
└── screens/      Route-level interface documentation
```

## File conventions

| File | ID prefix | Purpose |
|---|---|---|
| `*.screen.json` | `SCR-` | A route-level user interface |
| `*.component.json` | `CMP-` | A reusable interface component or pattern |
| `*.flow.json` | `FLOW-` | A journey across screens and actions |
| `*.tokens.json` | — | Optional structured design tokens |

VS Code associates each filename pattern with its schema through `.vscode/settings.json`, so invalid properties and enum values are reported while editing.

## Sources of truth

```text
Figma       Visual composition and prototypes
UX JSON     Purpose, relationships, states, decisions and code locations
Go / HTML   Behaviour and implementation
CSS :root   Design-token values used by this application
```

The token schema remains available for projects that need a portable token file, but this application keeps its implemented visual values in CSS custom properties.

## Writing guidance

Keep each JSON file as small as the subject allows.

Add a property when it answers a real planning or maintenance question. Do not restate content that is already obvious from the template or encode coordinates and pixel values that Figma and CSS already own.

Use descriptive, stable IDs such as:

```text
SCR-ADMIN-CLIENTS-SHOW
CMP-ADMIN-HEADER
FLOW-CLIENT-MANAGEMENT
UXD-AUTH-002
```

Update a document when a meaningful route, handler, template, state, relationship or UX decision changes.

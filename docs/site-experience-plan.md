# Site Experience Redesign Plan

## Objective

Turn the portfolio and personal workspace into one coherent product system: polished and expressive in public, calm and useful in private, and consistently fast because the experience remains server-rendered.

The redesign should feel more distinctive than a generic agency template without becoming visually noisy, overly editorial or dependent on a large JavaScript runtime.

## Reference synthesis

The direction combines the strongest patterns from the supplied references rather than copying any one visual language:

- **Kukie:** product clarity, measurable proof, modular utility cards and visible process.
- **GPTAgency:** decisive problem-led messaging, high-contrast accents and narrative scroll rhythm.
- **UXCentury:** enterprise confidence, restrained copy, generous spacing and an inspectable delivery process.
- **Promo Tech Store:** compact status language and precise interface structure.
- **Atomic:** cinematic scale, layered composition and portfolio energy.

## Design concept: Clarity / Control

The public site communicates the value: what Daniel builds, why it matters and how he thinks.

The private workspace provides control: clients, projects, contracts and journal publishing without unnecessary explanatory copy.

Both surfaces share:

- the original deep navy identity rather than a neutral black or green-led palette;
- polished cobalt, polo-blue and restrained icy-cyan accents;
- modern sans-serif typography with precise mono details;
- thin borders, controlled translucency and subtle blue light rather than heavy glassmorphism;
- large, tightly tracked headings balanced by compact operational labels;
- asymmetry used only where it improves hierarchy;
- motion that is brief, progressive and disabled when reduced motion is requested;
- the existing DJM / DEV logo system used consistently in public and private headers.

## Public experience

### Navigation

- Compact sticky header with the full branded logo.
- Direct routes to work, capabilities, about and journal.
- Accessible native-details mobile menu; no custom navigation JavaScript.
- One clear project CTA.

### Homepage

1. **Problem-led hero:** “Complex products. Clear experiences.”
2. **Proof layer:** product focus, design-engineering scope and core stack.
3. **Product-system visual:** a lightweight abstract interface composition rather than a screenshot that becomes stale after each redesign.
4. **Capability ticker:** a restrained marquee as atmosphere, not navigation.
5. **Selected work:** a purpose-built workspace composition that communicates the product without presenting obsolete screenshots.
6. **Capabilities:** fewer, stronger service statements with clear outcomes.
7. **Process:** visible product delivery sequence.
8. **Journal bridge:** technical writing becomes part of the portfolio narrative.
9. **Contact:** strong closing statement and a single obvious action.

### Case study

- Modern opening with product facts and clear positioning.
- A custom interface composition representing the public/private relationship instead of screenshots of the previous visual design.
- Clear problem, model, architecture, security and outcome sections.
- Technical detail remains readable without making the page feel like internal documentation.

### Journal

- Archive-first editorial index without serif-heavy styling.
- Comfortable article line length and strong Markdown typography.
- Metadata separated from prose on larger screens.
- Clear continuation paths back to the archive and contact.

## Personal workspace

### Shell

- Dedicated workspace navigation instead of reusing the public header.
- The real DJM / DEV mark rather than a temporary “DM” monogram.
- Persistent routes for overview, clients, projects, contracts and journal.
- Compact footer with no irrelevant technical status copy.
- Existing admin pages inherit the new shell while retaining their proven CRUD flows.

### Dashboard

- Sign out is grouped with genuine workspace actions instead of floating above the page.
- Four equal product-area cards replace the accidental three-plus-one layout.
- Quick actions form a coherent four-item group.
- Copy is written for Daniel as the sole user; it does not sell or explain the technology back to him.
- Database, framework and infrastructure labels are omitted unless they directly help complete a task.

### Journal administration

- Archive view prioritises status, title, slug and actions.
- Editor separates writing from publication controls.
- Markdown guidance is visible only where it is useful.
- Destructive actions are clearly isolated.

### Authentication

- Login gets its own distraction-free shell.
- The branded split composition reinforces that the workspace belongs to the same product.
- Existing server-side session, throttling and CSRF behaviour remains unchanged.

## Interaction and HTMX principles

- Every core journey must work as standard links and forms first.
- HTMX is used for small, local enhancements such as deletion and request progress.
- No client-side router, duplicated JSON API or hydration layer.
- Native browser behaviour is preferred for menus, focus, navigation and form semantics.
- View transitions are progressive enhancement only.

## Tailwind CSS v4 approach

- CSS-first configuration through `@theme`.
- OKLCH design tokens for the ink, mist, blue-signal and icy-electric palettes.
- Modern utilities such as `bg-linear-to-*`, dynamic spacing and arbitrary values where necessary.
- Reusable semantic component classes for buttons, panels, article prose and shells.
- Existing admin templates remain compatible through scoped blue variables while they are progressively refined.

## Accessibility and quality bar

- WCAG 2.2 AA contrast targets.
- Visible focus states and skip links on public, admin and auth shells.
- Semantic landmarks and descriptive navigation labels.
- Large controls with usable touch targets.
- Reduced-motion support.
- No content hidden behind fixed elements.
- Responsive testing from 320px through wide desktop.
- Template parsing, Go tests and Tailwind compilation must pass before merge.

## Implementation phases

### Phase 1 — Experience foundation

- Design tokens and reusable CSS components.
- Separate public, admin and auth shells.
- New public and admin navigation/footer components.
- Rendering helper routes pages to the correct shell.

### Phase 2 — High-impact surfaces

- Homepage.
- Portfolio case study.
- Journal archive and article.
- Login, dashboard and journal administration.
- Custom public 404 experience.

### Phase 3 — Progressive refinement

- Bring clients, projects and contracts to the same component quality.
- Admin-specific error states and empty-state refinement.
- Active navigation states and contextual breadcrumbs.
- Optional HTMX previews and autosave only after the core forms remain dependable.

### Phase 4 — Validation

- Full responsive and keyboard pass.
- Lighthouse and accessibility audit.
- Production CSS verification.
- Screenshot review before merge.

# Palette and interaction follow-up — 11 September 2026

## Intent

The owner prefers the newer page layouts but does not find the colours harmonious.
Keep those layouts. Use less saturated charcoal-blue surfaces, match the header and
footer to the page background, and use one softened tint of the existing Signal 500
blue for prominent accents. The public theme is scoped to `.site-shell` in
`editorial.css`; the private workspace and the original brand foundations remain
unchanged in `tokens.css`. Semantic aliases are explicitly rebound at the theme
boundary rather than relying on inherited aliases to recompute.

The palette is a design judgement to review visually, not an objectively optimal
choice or a colour prescription attributed to Don Norman.

## Principle-based changes

This is an application of Norman's published principles, not his review or endorsement.

- **Signifiers:** a visible Menu/Close label accompanies the mobile control; linked
  project/article rows retain visual cues without requiring hover. Static examples
  no longer show misleading Live/Clear status badges. The Journal's existing plain
  topic labels remain non-interactive, rather than pretending to be filters.
- **Conceptual clarity and mapping:** navigation says Services, matching the offer.
  The server renders the current Journal/Work/Services section, including without
  JavaScript. `aria-current="page"` is used only for the Journal archive itself;
  section ancestors use `location`. Homepage hash navigation updates its marker.
- **Feedback:** Copy email address reports successful copying only after the
  clipboard promise resolves. Failure says how to copy manually. No fake sent
  message, invented response-time promise or background submission is introduced.
- **Recovery/control:** Escape closes the menu and restores focus. Selecting a
  same-page destination closes the menu and moves focus to the destination heading.
  Ordinary navigation, history and modified clicks are left to the browser.
- **Predictable next steps:** the contact block explicitly says Email Daniel opens
  the visitor's email app, displays the selectable address, and suggests what to
  include. The copy button only appears when JavaScript and the Clipboard API exist.
- **Text enlargement:** a compact two-item phone header keeps the logo and labelled
  menu separate at 200% text size. The centred desktop mark remains unchanged.

The tiny self-hosted enhancement script has no runtime library, analytics or network
calls. Emailing still requires the visitor to send their message in their own app.

## Primary source

Don Norman, “Signifiers, not affordances” (author's version, ACM Interactions, 2008):
https://jnd.org/signifiers-not-affordances/

The interpretation is specific to this portfolio. Real usability testing with
prospective clients is still needed; heuristic improvements do not prove that
visitors understand the offer or prefer this palette.

## Checks

Retain the existing application/browser regression suites. The palette assertion
now verifies the explicitly chosen shared page/header/footer colour rather than
the superseded Ink 850 requirement. Additional interaction tests cover server-rendered
current-section markers, native navigation without JavaScript, Escape, same-page
focus, copying success/failure and visible labels. Accessibility scans and text-size
stress checks must continue to pass. CI screenshots are the visual review evidence;
historical Lighthouse scores are not results for this revision.

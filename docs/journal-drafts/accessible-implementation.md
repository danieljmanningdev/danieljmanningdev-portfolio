# Design tokens become useful when they survive implementation

A design system is not complete when its colours have names and its components look consistent in Figma. It becomes useful when the browser implementation can handle real text, different devices, keyboard input and failure states.

## Use semantic decisions

A page gutter and the padding inside a card have different purposes. The portfolio's shared container uses a responsive gutter token, allowing edge clearance to improve across public pages without unrelated margins on each section.

That change also reduces the space available for text. Headings and grid children need to be checked alongside the spacing, not treated as separate screenshots. Setting `min-width: 0` on a grid child can be necessary for its content to shrink as intended; clipping the entire page is not a repair for overflowing content.

## Implement the state, not only its appearance

A form needs labels and understandable error recovery. Navigation needs a usable keyboard order. A focus ring must remain visible against the actual background, and secondary text still needs sufficient contrast.

Native elements provide a reliable starting point. A link navigates; a button performs an action; `details` and `summary` offer disclosure behaviour. A Figma layer named “button” cannot supply those semantics by itself.

## Test the uncomfortable sizes

WCAG reflow guidance includes a width equivalent to 320 CSS pixels for vertically scrolling content, with defined exceptions. A clean desktop screenshot is not sufficient evidence.

The review adds browser checks across narrow and wide viewports, accessibility scans, keyboard navigation and no-JavaScript content checks. Automated results still cannot certify every reading sequence, label, complex background contrast or screen-reader interaction. Manual keyboard, zoom, real-device and assistive-technology review remain necessary.

A useful design-development loop is intention, implementation, observation and revision. The goal is not to defend the first mockup; it is to make the next version clearer and more dependable.

See [UI/UX design](https://danieljmanningdev.com/ui-ux-design/) and the [portfolio implementation](https://github.com/danieljmanningdev/danieljmanningdev-portfolio).

References: [WCAG 2.2](https://www.w3.org/WAI/WCAG22/quickref/), [Understanding Reflow](https://www.w3.org/WAI/WCAG22/Understanding/reflow.html), [Playwright accessibility-testing limits](https://playwright.dev/docs/accessibility-testing).

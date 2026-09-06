# Accessibility Audit

**Status:** Not started
**Target:** WCAG 2.2 AA
**Site:** danieljmanningdev.com

## Tools

- axe DevTools
- WCAG 2.2 Quick Reference
- WebAIM Contrast Checker
- Browser DevTools

## Automated Audit

### axe DevTools

Test representative pages:

- [ ] Homepage
- [ ] Portfolio
- [ ] Case study
- [ ] Journal index
- [ ] Journal article
- [ ] Web Design
- [ ] Web Development
- [ ] Software Development
- [ ] UI/UX Design
- [ ] Web Design Leeds

For each page:

- [ ] Run axe
- [ ] Record violations
- [ ] Manually verify findings
- [ ] Create GitHub issues for genuine defects

## Keyboard Audit

Disconnect/ignore the mouse.

For each representative page:

- [ ] Tab through every interactive element
- [ ] Focus indicator is always visible
- [ ] Focus order follows the visual/logical order
- [ ] No keyboard traps
- [ ] Links can be activated
- [ ] Buttons can be activated
- [ ] Navigation can be operated
- [ ] Skip link works, if provided
- [ ] Focus does not disappear behind content

## Zoom

Test browser zoom:

- [ ] 200%
- [ ] 400%

Check:

- [ ] No important content disappears
- [ ] Text remains readable
- [ ] Controls remain usable
- [ ] Horizontal scrolling is avoided where WCAG reflow requires it
- [ ] Navigation remains usable

## Colour and Contrast

Check:

- [ ] Normal text contrast
- [ ] Large text contrast
- [ ] Interactive controls
- [ ] Focus indicators
- [ ] Error states
- [ ] Information does not depend on colour alone

## Forms

Where applicable:

- [ ] Every control has an accessible label
- [ ] Placeholder is not the only label
- [ ] Required fields are understandable
- [ ] Errors identify the affected field
- [ ] Errors explain how to fix the problem
- [ ] Autocomplete attributes are appropriate
- [ ] Keyboard interaction works

## Images and Media

- [ ] Informative images have useful alt text
- [ ] Decorative images are ignored appropriately
- [ ] Alt text does not unnecessarily repeat nearby text

## Structure

- [ ] One sensible page-level heading
- [ ] Heading hierarchy is logical
- [ ] Landmarks are appropriate
- [ ] Lists use list semantics
- [ ] Navigation uses appropriate semantics
- [ ] Links make sense from their accessible name/context

## Screen Reader Sanity Test

Test representative critical flows.

- [ ] Page title is useful
- [ ] Headings provide useful navigation
- [ ] Landmarks make sense
- [ ] Links/buttons have understandable names
- [ ] Dynamic changes are announced where required
- [ ] Forms can be understood and completed

## Motion

- [ ] prefers-reduced-motion is respected where appropriate
- [ ] No essential information depends on animation
- [ ] No problematic flashing content

## Findings

| ID | Page | Finding | WCAG | Priority | Issue | Status |
| --- | --- | --- | --- | --- | --- | --- |
| A11Y-001 | | | | | | Open |

## Retest

After fixes:

- [ ] Automated tests rerun
- [ ] Keyboard tests rerun
- [ ] Affected screen-reader behaviour retested
- [ ] Affected viewport/zoom behaviour retested
- [ ] Closed issues verified rather than assumed fixed
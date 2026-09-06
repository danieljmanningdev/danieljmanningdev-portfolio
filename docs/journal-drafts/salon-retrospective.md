# Revisiting a salon website as a design and development exercise

Salon Rebuild is a fictional portfolio project, not a commissioned redesign or an operating salon. Its illustrative service information and enquiry flow are not evidence of actual bookings or a conversion increase.

The exercise asks a practical question: what should a visitor understand, and what should be easy to do next?

## Give the visitor a clear route

Services, team, gallery and contact pages have different jobs. Photography can establish character, but service information and the next action must remain understandable. A demonstration enquiry flow must not imply that it creates a real appointment.

## Make consistency work with real content

Typography, spacing and repeated components should keep the pages related. Cards need to tolerate longer titles; navigation should not change meaning between desktop and mobile. Responsive design is more than shrinking the desktop composition until it technically fits.

For a commissioned site, I would validate those decisions through user tasks and real business requirements. This exercise can explain intended behaviour and show implementation, but it cannot supply invented research findings.

## Keep development proportional

Go and server-rendered templates suit the project's informational pages. Repeated sections can be templates while ordinary links and headings remain meaningful.

The portfolio screenshot is a separate performance concern. Sending a full-resolution PNG at every viewport is wasteful. The review supplies multiple image sizes in AVIF and WebP, retains a fallback and declares intrinsic dimensions. Those are implementation improvements, not a claim that every visitor now experiences a particular speed.

## Make the evidence easy to inspect

The [case study](https://danieljmanningdev.com/work/salon-rebuild/) connects the visual decisions to a live demonstration and [source](https://github.com/danieljmanningdev/salon-rebuild). A real client version would still need an approved content set, privacy review, production enquiry handling and measured outcomes.

This is a record of design judgement and development practice—not a substitute for genuine client evidence.

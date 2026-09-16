package casestudies

import "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"

var WireframeKit = models.CaseStudy{
	Slug:    "wireframe-kit",
	Year:    "2026",
	Title:   "Wireframe",
	Accent:  "Kit",
	Summary: "A reusable Figma wireframe kit with 40+ components, built to make early product exploration faster, clearer and more consistent.",

	Tags: []string{
		"Figma",
		"Wireframing",
		"UI / UX",
		"Prototyping",
	},

	Metrics: []models.CaseStudyMetric{
		{
			Label: "Project",
			Value: "Wireframe Kit",
		},
		{
			Label: "Components",
			Value: "40+",
		},
		{
			Label: "Focus",
			Value: "Rapid Prototyping",
		},
	},

	Hero: models.CaseStudyMedia{
		Src: "/static/images/wireframe-kit-cover.png",

		Alt:     "Cover of the Wireframe Kit created in Figma",
		Caption: "Wireframe Kit — a reusable collection of low-fidelity components for exploring product layouts and flows.",
		Width:   1920,
		Height:  1080,
	},

	Sections: []models.CaseStudySection{
		{
			Number: "01",
			Label:  "The problem",
			Title:  "Good interface design needs somewhere to start.",
			Type:   "editorial",

			Paragraphs: []string{
				"Jumping directly into polished visual design can make it too easy to commit to colours, typography and styling before the underlying layout and hierarchy have been properly explored.",
				"Wireframes provide a lower-fidelity space for testing structure, content and flow first. They make it easier to move things around, challenge assumptions and compare different approaches without treating every early decision as final.",
				"I built this kit to create a faster and more consistent starting point for that part of the product design process.",
			},
		},

		{
			Number: "02",
			Label:  "Design direction",
			Type:   "cards",
			Soft:   true,

			Cards: []models.CaseStudyCard{
				{
					Index: "01",
					Title: "Structure before styling",
					Copy:  "The components deliberately stay visually restrained so attention remains on layout, hierarchy and interaction rather than finished visual design.",
				},
				{
					Index: "02",
					Title: "Reusable by default",
					Copy:  "The kit is designed around reusable components so common interface patterns can be assembled quickly without rebuilding the same foundations for every screen.",
				},
				{
					Index: "03",
					Title: "Fast to explore",
					Copy:  "Low-fidelity components make it easier to experiment with multiple directions and change layouts before more time is invested in polished UI.",
				},
			},
		},

		{
			Number: "03",
			Label:  "What it enables",
			Title:  "A reusable foundation for early product thinking.",
			Type:   "list",

			Items: []string{
				"Faster exploration of page and screen layouts",
				"More consistent low-fidelity interface components",
				"Clearer focus on hierarchy before visual styling",
				"Reusable starting points across different product ideas",
				"Easier iteration when requirements or layouts change",
				"A clearer transition from wireframe to finished visual UI",
			},
		},

		{
			Number: "04",
			Label:  "Workflow",
			Title:  "Keep early decisions cheap to change.",
			Type:   "editorial",

			Paragraphs: []string{
				"The kit is intended for the stage where an idea needs enough structure to become tangible, but not so much visual polish that changing direction becomes expensive.",
				"Components can be combined into rough screens and flows, rearranged as the product develops and then gradually replaced or developed into the finished visual interface.",
				"Keeping those stages separate makes the progression from bare structure to polished UI much easier to see and evaluate.",
			},
		},

		{
			Number: "05",
			Label:  "Retrospective",
			Title:  "A small resource that improves the process around larger designs.",
			Type:   "editorial",
			Soft:   true,

			Paragraphs: []string{
				"Building the kit made me think more deliberately about which interface patterns are genuinely reusable and which decisions belong to an individual product.",
				"It also reinforced the value of separating structural decisions from visual ones. A wireframe does not need to look finished to be useful; its job is to make the next design decision easier.",
				"Publishing the kit on Figma Community also turns something I built for my own workflow into a resource that other designers can use and adapt.",
			},
		},
	},

	Explore: models.CaseStudyExploreSection{
		Title:       "Explore the wireframe kit",
		Accent:      "and use it in your own workflow.",
		Description: "The complete kit is published on Figma Community with 40+ reusable components for early-stage product design and rapid prototyping.",

		Links: []models.CaseStudyExploreLink{

			{
				Label: "View on Figma",
				URL:   "https://www.figma.com/community/file/1678446211261822616/wireframe-kit-danieljmanningdev-com?q_id=80b050b8-60e0-4b24-a71e-005f6cf48595",
				Style: "primary",
			},
		},

		Note: "Available as a Figma Community resource.",
	},

	Outcome: models.CaseStudyOutcomeSection{
		Title:    "A reusable starting point for turning rough product ideas into structured wireframes",
		Accent:   "before committing to polished visual UI.",
		CTALabel: "Discuss a project",
		CTAURL:   "mailto:daniel@danieljmanningdev.com",
	},
}

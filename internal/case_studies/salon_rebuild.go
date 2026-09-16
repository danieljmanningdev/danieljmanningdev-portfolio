package repository

import "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"

var SalonRebuild = models.CaseStudy{
	Slug:    "salon-rebuild",
	Year:    "2026",
	Title:   "Salon",
	Accent:  "Rebuild",
	Summary: "The project revisits the type of salon website work I produced early in my freelance career and explores how I would approach the same broad category today using a more considered UI/UX and engineering process.",

	Tags: []string{
		"Website Redesign",
		"Go",
		"Server-rendered HTML",
		"UI / UX",
	},

	Metrics: []models.CaseStudyMetric{
		{
			Label: "Project",
			Value: "Full Site Rebuild",
		},
		{
			Label: "Content",
			Value: "Fictional Concept",
		},
		{
			Label: "Experience",
			Value: "Sleek Luxury",
		},
	},

	Hero: models.CaseStudyMedia{
		Src:     "/static/images/salon-rebuild-home.png",
		Alt:     "Homepage of the fictional Salon Rebuild concept",
		Caption: "Salon Rebuild — a fictional retrospective concept, not a commissioned redesign for the original business.",
		Width:   2880,
		Height:  1561,
	},

	Sections: []models.CaseStudySection{
		{
			Number: "01",
			Label:  "The problem",
			Title:  "Revisiting the kind of salon website I originally built with WordPress.",
			Type:   "editorial",

			Paragraphs: []string{
				"One of my earliest freelance projects was a salon website built with WordPress. Several years later, I wanted to revisit the same broad type of brief and see how differently I would approach it with the design and engineering experience I have now.",
				"This rebuild is not a redesign commissioned by the original business and does not reuse its content, staff information, pricing, testimonials or branding. Instead, it is a fictional salon concept designed specifically as a retrospective portfolio project.",
				"The aim was to create a much more modern, responsive and deliberate experience while keeping the visual direction restrained: premium enough to feel appropriate for a contemporary salon without becoming decorative for the sake of it.",
			},
		},

		{
			Number: "02",
			Label:  "Design direction",
			Type:   "cards",

			Cards: []models.CaseStudyCard{
				{
					Index: "01",
					Title: "Luxury without excess",
					Copy:  "A restrained palette, large editorial typography and generous spacing create a premium feel without relying on excessive decoration or visual noise.",
				},
				{
					Index: "02",
					Title: "Mobile-first hierarchy",
					Copy:  "Content is structured to remain clear on smaller screens, with strong visual hierarchy, readable typography and responsive layouts that adapt rather than simply shrink.",
				},
				{
					Index: "03",
					Title: "Clear booking journey",
					Copy:  "Services, team information, gallery content and appointment calls-to-action are organised around the actions a potential salon customer is most likely to take.",
				},
			},
		},

		{
			Number: "03",
			Label:  "Technical approach",
			Title:  "Simple architecture for a simple product.",
			Type:   "architecture",

			Paragraphs: []string{
				"The rebuild does not need a database, authentication system or client-side application framework. Go handles routing and renders reusable HTML templates on the server, while Tailwind provides the styling layer.",
				"Keeping the implementation deliberately small makes the project easy to understand, maintain and deploy while still demonstrating a structured Go web application rather than a collection of disconnected static files.",
			},

			Steps: []models.CaseStudyStep{
				{
					Number: "01",
					Label:  "Request",
					Value:  "Go ServeMux",
				},
				{
					Number: "02",
					Label:  "Application",
					Value:  "Page handlers + PageData",
				},
				{
					Number: "03",
					Label:  "Presentation",
					Value:  "Shared Go templates + Tailwind",
				},
				{
					Number: "04",
					Label:  "Response",
					Value:  "Server-rendered HTML",
				},
			},
		},

		{
			Number: "04",
			Label:  "What changed",
			Title:  "Same category. Completely different process.",
			Type:   "list",

			Items: []string{
				"Stronger typography and visual hierarchy",
				"Responsive layouts designed intentionally for mobile",
				"Clearer information architecture and calls-to-action",
				"Reusable Go layouts and page templates",
				"Version-controlled development with signed Git history",
				"Defined scope instead of allowing the project to expand indefinitely",
			},
		},

		{
			Number: "05",
			Label:  "Retrospective",
			Title:  "The rebuild is as much about process as visual improvement.",
			Type:   "editorial",

			Paragraphs: []string{
				"The original salon project came from a point in my career where simply completing a real client website felt like the achievement. I had much less confidence around pricing, scope, design systems and the value of my own time.",
				"Rebuilding the idea now made the progression easier to see. The difference is not just cleaner typography or a more contemporary colour palette. I approach structure, responsiveness, maintainability and project boundaries differently too.",
				"That makes this project useful to me as more than a visual redesign. It is a record of how both my design thinking and engineering process have developed.",
			},
		},
	},

	Explore: models.CaseStudyExploreSection{
		Title:       "See the finished rebuild and the implementation behind it.",
		Description: "The complete concept is deployed as a live Go web service, while the public repository exposes the implementation and development history behind the finished design.",

		Links: []models.CaseStudyExploreLink{
			{
				Label: "View live site",
				URL:   "https://salon-rebuild.onrender.com/",
				Style: "primary",
			},
			{
				Label: "View source",
				URL:   "https://github.com/danieljmanningdev/salon-rebuild",
				Style: "secondary",
			},
		},

		Note: "Live demo hosted on Render's free tier. The first request may take a moment if the service has been idle.",
	},

	Outcome: models.CaseStudyOutcomeSection{
		Title:    "A modern salon experience that shows how much my approach has changed",
		Accent:   "since my earliest freelance work.",
		CTALabel: "Discuss a project",
		CTAURL:   "mailto:daniel@danieljmanningdev.com",
	},
}

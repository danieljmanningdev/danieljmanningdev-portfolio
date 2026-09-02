// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package http

type relatedLink struct {
	Label       string
	Title       string
	Description string
	URL         string
}

func salonCaseStudyRelatedLink() relatedLink {
	return relatedLink{
		Label:       "Case study",
		Title:       "Salon Rebuild",
		Description: "See a salon concept rebuilt with a more considered UI/UX, responsive design and Go-based delivery process.",
		URL:         "/work/salon-rebuild/",
	}
}

func portfolioCaseStudyRelatedLink() relatedLink {
	return relatedLink{
		Label:       "Case study",
		Title:       "Portfolio & Client Workspace",
		Description: "Explore the product architecture, security and server-rendered engineering behind this portfolio and private workspace.",
		URL:         "/work/portfolio",
	}
}

func webDesignRelatedLink() relatedLink {
	return relatedLink{
		Label:       "Service",
		Title:       "Web Design",
		Description: "Responsive, accessible website design shaped around clear journeys, useful content and maintainable systems.",
		URL:         "/web-design/",
	}
}

func webDesignLeedsRelatedLink() relatedLink {
	return relatedLink{
		Label:       "Leeds service",
		Title:       "Web Design in Leeds",
		Description: "Web design and development for Leeds businesses that need a clearer, faster and more dependable online presence.",
		URL:         "/web-design-leeds/",
	}
}

func webDevelopmentRelatedLink() relatedLink {
	return relatedLink{
		Label:       "Service",
		Title:       "Web Development",
		Description: "Server-rendered websites and web applications built for performance, accessibility and long-term maintainability.",
		URL:         "/web-development/",
	}
}

func softwareDevelopmentRelatedLink() relatedLink {
	return relatedLink{
		Label:       "Service",
		Title:       "Software Development",
		Description: "Focused software for real workflows, from internal tools and portals to custom web applications.",
		URL:         "/software-development/",
	}
}

func uiUXDesignRelatedLink() relatedLink {
	return relatedLink{
		Label:       "Service",
		Title:       "UI & UX Design",
		Description: "Interface and experience design that reduces friction, clarifies decisions and keeps complex products coherent.",
		URL:         "/ui-ux-design/",
	}
}

func journalRelatedLink() relatedLink {
	return relatedLink{
		Label:       "Journal",
		Title:       "Design and engineering notes",
		Description: "Read practical notes on Go, server-rendered architecture, product design, security and project decisions.",
		URL:         "/blog/",
	}
}

func pricingLessonsRelatedLink() relatedLink {
	return relatedLink{
		Label:       "Project retrospective",
		Title:       "I Built a Website for 50p an Hour. Never Again.",
		Description: "The pricing, scope and ownership lessons behind the early freelance project that prompted the Salon Rebuild.",
		URL:         "/blog/i-built-a-website-for-50p-an-hour-never-again",
	}
}

func relatedLinksForJournalPost(slug string) []relatedLink {
	switch slug {
	case "i-built-a-website-for-50p-an-hour-never-again":
		return []relatedLink{
			salonCaseStudyRelatedLink(),
			webDesignRelatedLink(),
			uiUXDesignRelatedLink(),
		}

	case "why-i-use-htmx-over-javascript-bloatware":
		return []relatedLink{
			webDevelopmentRelatedLink(),
			portfolioCaseStudyRelatedLink(),
			softwareDevelopmentRelatedLink(),
		}

	case "sqlite-can-go-to-production":
		return []relatedLink{
			softwareDevelopmentRelatedLink(),
			portfolioCaseStudyRelatedLink(),
			webDevelopmentRelatedLink(),
		}

	case "why-i-chose-go-htmx-sqlite-instead-of-react-for-my-portfolio":
		return []relatedLink{
			portfolioCaseStudyRelatedLink(),
			webDevelopmentRelatedLink(),
			softwareDevelopmentRelatedLink(),
		}

	default:
		return []relatedLink{
			portfolioCaseStudyRelatedLink(),
			webDevelopmentRelatedLink(),
			journalRelatedLink(),
		}
	}
}

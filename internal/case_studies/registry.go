package casestudies

import "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"

var studies = map[string]models.CaseStudy{
	SalonRebuild.Slug: SalonRebuild,
	WireframeKit.Slug: WireframeKit,
}

func Get(slug string) (models.CaseStudy, bool) {
	study, ok := studies[slug]
	return study, ok
}

package models

type CaseStudy struct {
	Slug    string
	Year    string
	Title   string
	Accent  string
	Summary string
	Tags    []string

	Metrics []CaseStudyMetric
	Hero    CaseStudyMedia

	Sections []CaseStudySection

	Explore CaseStudyExploreSection
	Outcome CaseStudyOutcomeSection
}

type CaseStudyMetric struct {
	Label string
	Value string
}

type CaseStudyMedia struct {
	Src     string
	Alt     string
	Caption string
	Width   int
	Height  int
}

type CaseStudySection struct {
	Number string
	Label  string
	Title  string
	Type   string

	Paragraphs []string
	Cards      []CaseStudyCard
	Items      []string
	Steps      []CaseStudyStep
	Media      *CaseStudyMedia
}

type CaseStudyCard struct {
	Index string
	Title string
	Copy  string
}

type CaseStudyStep struct {
	Number string
	Label  string
	Value  string
}

type CaseStudyExploreSection struct {
	Title       string
	Description string
	Links       []CaseStudyExploreLink
	Note        string
}

type CaseStudyExploreLink struct {
	Label string
	URL   string
	Style string
}

type CaseStudyOutcomeSection struct {
	Title    string
	Accent   string
	CTALabel string
	CTAURL   string
}

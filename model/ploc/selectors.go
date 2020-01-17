package ploc

// SelectByKeyword traverses a list of subjects and returns the first that matches the specified keyword.
func (subs *Subjects) SelectByKeyword(keyword string) Subject {
	for _, s := range *subs {
		if s.Keyword == keyword {
			return s
		}
	}

	return Subject{}
}

// SelectByName traverses a list of expert previews and returns the first that matches the specified expert name.
func (exps *ExpertPreviews) SelectByName(name string) ExpertPreview {
	for _, e := range *exps {
		if e.Name == name {
			return e
		}
	}

	return ExpertPreview{}
}

// SelectByTitle traverses a list of collection names and returns the first that matches the specified title.
func (collections *Collections) SelectByTitle(title string) Collection {
	for _, c := range *collections {
		if c.Title == title {
			return c
		}
	}

	return Collection{}
}

// SelectByTitle traverses a list of bookmarked records and returns the first that matches the specified title.
func (bookmarks *RecordBookmarks) SelectByTitle(title string) RecordBookmark {
	for _, bm := range *bookmarks {
		if bm.Title == title {
			return bm
		}
	}

	return RecordBookmark{}
}

// SelectByTitle traverses a list of record previews and returns the first that matches the specified title.
func (recs *RecordPreviews) SelectByTitle(title string) RecordPreview {
	for _, r := range *recs {
		if r.Title == title {
			return r
		}
	}

	return RecordPreview{}
}

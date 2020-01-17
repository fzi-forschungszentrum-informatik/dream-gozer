package ploc

// Collection is used to send the name of a bookmark collection in JSON format to the ploc client app.
type Collection struct {
	Id    int64  `json:"id"`
	Title string `json:"name"`
}

// Collections is used to send a list of bookmark collection names in JSON format to the ploc client app.
// Such a list is used for example to inform the ploc client app about all the user's existing collection names.
type Collections []Collection

// ExpertBookmark is used to send a preview of a bookmarked expert in JSON format to the ploc client app.
type ExpertBookmark ExpertPreview

// ExpertBookmarks is used to communicate a list of bookmarked experts in JSON format to the ploc client app.
// Such a list is used for example to present a user's expert bookmark list.
type ExpertBookmarks ExpertPreviews

// ExpertPreview is used to send a preview of an expert in JSON format to the ploc client app.
// ExpertPreview is used to preview an expert in the expert feed.
type ExpertPreview struct {
	Id                    int64    `json:"id"`
	Name                  string   `json:"name"`
	LastPublicationYear   int64    `json:"last_publication_year"`
	TotalPublicationCount int64    `json:"total_publication_count"`
	Subjects              []string `json:"subjects"`
}

// ExpertBookmarks is used to communicate a list of expert previews in JSON format to the ploc client app.
// Such a list is used for example to present a segement of a user's expert feed.
type ExpertPreviews []ExpertPreview

// Feedback is used to send a lightweight review in JSON format to the ploc client app.
// The review defines binary flags for high relevance, high quality of presentation, and sound methodology (0=no,1=yes).
type Feedback struct {
	RecordId     int64  `json:"record_id"`
	OrcId        string `json:"orcid"`
	Relevance    int64  `json:"relevance"`
	Presentation int64  `json:"presentation"`
	Methodology  int64  `json:"methodology"`
}

// Feedbacks is used to send a list of feedback in JSON format to the ploc client app.
// Such a list is used for example to summarize all the feedbacks of a specific publication.
type Feedbacks []Feedback

// Keywords is used to send a list of subjects or topics in JSON format to the ploc client app.
// Such a list is used for example to classify a publication or person in terms of topics.
type Keywords []string

// Names is used to send a list of names in JSON format to the ploc client app.
// Such a list used for example to name all authors of a publication.
// The name can be in its abreviated form (e.g. "J. Doe") or its full form (e.g. "John Doe")
type Names []string

// RecordBookmark is used to send a preview of a bookmarked record in JSON format to the ploc client app.
// The preview contains the bookmarks database ID, the title of the record, the publication year, the
// names of the creators, the related subjects, the kind of publication (e.g. thesis), a flag that signifies if
// the details of that record were watched before, and the collections of which the bookmarked record is part of.
type RecordBookmark struct {
	Id            int64    `json:"id"`
	Title         string   `json:"title"`
	Year          int64    `json:"year"`
	Creators      string   `json:"creators"`
	Subjects      []string `json:"subjects"`
	Type          int64    `json:"type"`
	visited       bool     `json:"visited"`
	CollectionIds []int64  `json:"collection_ids"`
}

// RecordBookmarks is used to send a list of bookmarked records in JSON format to the ploc client app.
type RecordBookmarks []RecordBookmark

// RecordPreview is used to send a preview of a publication record in JSON format to the ploc client app.
// A record preview is used for example to preview a record in the record and feedback feed.
type RecordPreview struct {
	Id       int64    `json:"id"`
	Title    string   `json:"title"`
	Year     int64    `json:"year"`
	Creators string   `json:"creators"`
	Subjects []string `json:"subjects"`
	Abstract string   `json:"abstract"`
	Type     int64    `json:"type"`
	Visited  bool     `json:"visited"`
}

// RecordPreviews is used to send a list of record previews.
// Such a list is used for example to present a segement of a user's record or feedback feed.
type RecordPreviews []RecordPreview

// RecordType is used to send a kind of publication in JSON format to the ploc client app.
// By kind we mean for example PhD thesis or conference paper.
type RecordType struct {
	Id      int64  `json:"id"`
	Keyword string `json:"keyword"`
}

// RecordTypes is used to send a list of all supported kind of publications in JSON format to the ploc client app.
type RecordTypes []RecordType

// Subject is used to send a single topic or keyword in JSON format to the ploc client app.
// A subject is used for example to classify a publication, creator or expert.
type Subject struct {
	Id      int64  `json:"id"`
	Keyword string `json:"keyword"`
}

// Subjects is used to send a list of topics or keywords in JSON format to the ploc client app.
// Such a list is used for example to name all the topics that a publications relates to.
type Subjects []Subject

// TinyRecord is used to send a condensed summary of a publication record in JSON format to the ploc client app.
// TinyRecord is shorter than a RecordPreview and it used for example in the publication history of an expert.
type TinyRecord struct {
	Id       int64    `json:"id"`
	Title    string   `json:"title"`
	Year     int64    `json:"year"`
	Creators []string `json:"creators"`
}

// TinyRecords is used to send a list of topics or keywords in JSON format to the ploc client app.
// Such a list is used for example to present the publication history of an expert.
type TinyRecords []TinyRecord

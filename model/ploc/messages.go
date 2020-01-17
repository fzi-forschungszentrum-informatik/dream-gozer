package ploc

import (
	"encoding/json"
)

// *** USER PROFILE ***************************************

// CreateUserProfileRequest defines a request of a user to create a new user profile.
// The provided secret is stored in its hashed form and used for authentication in the future.
type CreateUserProfileRequest struct {
	Secret string `json:"secret"`
}

// CreateUserProfileResponse defines a response to a user after creating a new user profile.
// The provided GUID is used to relate to a user's profile in future requests (as part of HTTP Basic Auth).
type CreateUserProfileResponse struct {
	GUID string `json:"guid"`
}

// *** EXPERT PROFILE *************************************

// CreateExpertProfileRequest defines a request of a user to register as an expert.
// The provided ORCiD identfifier is used to bind a user against the specified ORCiD profile.
type CreateExpertProfileRequest struct {
	OrcId string `json:"orcid"`
}

// ReadExpertProfileResponse defines a response to a user returning its ORCiD.
type ReadExpertProfileResponse struct {
	OrcId string `json:"orcid,omitempty"`
}

// *** SUBJECTS *******************************************

// ReadSubjectsResponse defines a response returning all supported subjects.
// The provided list of subjects can be used to specify a user's interests.
type ReadSubjectsResponse struct {
	Subjects Subjects `json:"subjects"`
}

// *** INTERESTS ******************************************

// CreateExpertProfileRequest defines a request of a user to add a specific subject his interests.
type CreateInterestRequest struct {
	SubjectId int64 `json:"subject_id"`
}

// ReadSubjectsResponse defines a response after adding a specific subject to the user's interests.
// The provided record count defines how many publication records still match the user's interests.
type CreateInterestResponse struct {
	RecordCount int64 `json:"record_count"`
}

// CreateExpertProfileRequest defines a request of a user to remove a specific subject from the list of his interests.
type DeleteInterestRequest struct {
	SubjectId int64 `json:"subject_id"`
}

// DeleteInterestResponse defines a response after a specific subject was removed from the list of the user's interests.
// The provided record count defines how many publication records still match the user's interests.
type DeleteInterestResponse struct {
	RecordCount int64 `json:"record_count"`
}

// ReadInterestsResponse defines a response returning all the interests of an user.
// The provided subjects define the subjects of user interests.
// The provided record count defines about how many publication records still match the user's interests.
type ReadInterestsResponse struct {
	Subjects    Subjects `json:"subjects"`
	RecordCount int64    `json:"record_count"`
}

// *** RECORD-TYPES ***************************************

// ReadRecordTypesResponse defines a response returning all supported record types.
// The provided list defines which kind of records are supported in general (e.g. thesis).
type ReadRecordTypesResponse struct {
	Types RecordTypes `json:"record_types"`
}

// *** RECORD-FEED ****************************************

// ReadRecordFeedRequest defines a request of a user for a segment a publication feed that respects his interests.
// The limit defines the maximum number of records that should be returned, while the offset defines the start index within
// the feed.
type ReadRecordFeedRequest struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

// ReadRecordFeedResponse defines a response that returns the user's personal publication feed.
// Offset and limit duplicate the requested position and number of records from the request.
// The list of records contains the specified segment with a preview for each record.
// The fields Records and RawRecords are used for either marshalling (RawRecords) or unmarshalling (Records).
// RawRecords directly map to precomputed JSON-data from the database for performance reasons.
type ReadRecordFeedResponse struct {
	Offset     int64             `json:"offset"`
	Limit      int64             `json:"limit"`
	Records    RecordPreviews    `json:"records"`
	RawRecords []json.RawMessage `json:"-"`
}

// SearchRecordFeedRequest defines a request of a user to show only the publications from his personalized feed that contain the
// specified search term. The limit defines the maximum number of records that should be returned, while the offset defines
// the start index within the feed.
type SearchRecordFeedRequest struct {
	SearchTerm string `json:"search_term"`
	Offset     int64  `json:"offset"`
	Limit      int64  `json:"limit"`
}

// SearchRecordFeedResponse defines a response returning a user's search results in his personal publication feed.
// Offset and limit duplicate the requested position and number of records from the request.
// The list of records contains the specified segment within the search results with a preview for each record.
// The fields Records and RawRecords are used for either marshalling (RawRecords) or unmarshalling (Records).
// RawRecords directly map to precomputed JSON-data from the database for performance reasons.
type SearchRecordFeedResponse struct {
	Offset     int64             `json:"offset"`
	Limit      int64             `json:"limit"`
	Records    RecordPreviews    `json:"records"`
	RawRecords []json.RawMessage `json:"-"`
}

// ReadRecordDetailsRequest defines a request for detailed information about a record.
type ReadRecordDetailsRequest struct {
	RecordId int64 `json:"record_id"`
}

// ReadRecordDetailsResponse defines a response returning detailed information about a record.
// The details include the record's title, related subjects, name of the authors, year of publication, abstract,
// type of publication, document object identifier, front page of the publisher and link to a PDF version of the document.
// These fields and RawDetails are used for either marshalling (RawDetails) or unmarshalling (Id,Title,...,PDFLink).
// RawDetails directly map to precomputed JSON-data from the database for performance reasons.
type ReadRecordDetailsResponse struct {
	Id             int64           `json:"id"`
	Title          string          `json:"title"`
	Creators       Names           `json:"creators"`
	Subjects       Keywords        `json:"subjects"`
	Year           int64           `json:"year"`
	Teaser         string          `json:"abstract"`
	Type           int64           `json:"type"`
	Doi            string          `json:"doi,omitempty"`
	RepositoryLink string          `json:"repository_link,omitempty"`
	PDFLink        string          `json:"pdf_link,omitempty"`
	RawDetails     json.RawMessage `json:"-"`
}

// *** EXPERT-FEED ****************************************

// ReadExpertFeedRequest defines a request of a user its expert feed. The limit defines the maximum
// number of experts that should be returned, while the offset defines the start index within the overall feed.
type ReadExpertFeedRequest struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

// ReadExpertFeedResponse defines a response returning a user's personal expert feed.
// Offset and limit duplicate the requested position and number of records from the request.
// The list of experts contains the specified segment with preview information for each expert.
// The fields Experts and RawExperts are used for either marshalling (RawExperts) or unmarshalling (Experts).
// RawExperts directly map to precomputed JSON-data from the database for performance reasons.
type ReadExpertFeedResponse struct {
	Offset     int64             `json:"offset"`
	Limit      int64             `json:"limit"`
	Experts    ExpertPreviews    `json:"experts"`
	RawExperts []json.RawMessage `json:"-"`
}

// SearchExpertFeedRequest defines a request of a user to show only the experts his expert feed that contain the
// specified search term. The limit defines the maximum number of experts that should be returned, while the offset defines
// the start index within the feed.
type SearchExpertFeedRequest struct {
	SearchTerm string `json:"search_term"`
	Offset     int64  `json:"offset"`
	Limit      int64  `json:"limit"`
}

// SearchExpertFeedResponse defines a response returning search results from the user's personal expert feed.
// Offset and limit duplicate the requested position and number of records from the request.
// The list of experts contains the specified segment within the search results with preview information for each expert.
// The fields Experts and RawExperts are used for either marshalling (RawExperts) or unmarshalling (Experts).
// RawExperts directly map to precomputed JSON-data from the database for performance reasons.
type SearchExpertFeedResponse struct {
	Offset     int64             `json:"offset"`
	Limit      int64             `json:"limit"`
	Experts    ExpertPreviews    `json:"experts"`
	RawExperts []json.RawMessage `json:"-"`
}

// ReadExpertDetailsRequest defines a request for detailed information about an expert.
type ReadExpertDetailsRequest struct {
	ExpertId int64 `json:"expert_id"`
}

// ReadExpertDetailsResponse defines a response returning detailed information about an expert.
// The details include the expert's name, subjects of expertise, ORCiD identifier, last known year of publication,
// and a list of all its publications.
// These fields and RawDetails are used for either marshalling (RawDetails) or unmarshalling (ExpertId,Name,...,Records).
// RawDetails directly map to precomputed JSON-data from the database for performance reasons.
type ReadExpertDetailsResponse struct {
	ExpertId            int64           `json:"expert_id"`
	Name                string          `json:"name"`
	Subjects            Keywords        `json:"subjects"`
	OrcId               string          `json:"orcid"`
	LastPublicationYear int64           `json:"last_publication_year"`
	Records             TinyRecords     `json:"records"`
	RawDetails          json.RawMessage `json:"-"`
}

// *** FEEDBACK-FEED **************************************

// ReadFeedbackFeedRequest defines a request of a user in its role as an domain expert to return the publications that may need a review.
// The result contains publication records that are in the domain of expertise and that may be of interest.
// A user must register as an expert (via ORCiD) before he can access its personal feedback feed.
// The limit defines the maximum number of records that should be returned, while the offset defines the start index within
// the feed.
type ReadFeedbackFeedRequest struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

// ReadFeedbackFeedResponse defines a response returning a user's feedback feed with publications that may require a review.
// The feedback feed contains publication records that are in the domain of expertise of that user and that may be of interest to him.
// Offset and limit duplicate the requested position and number of records from the request.
// The list of records contains the specified segment with a preview for each record.
// The fields Records and RawRecords are used for either marshalling (RawRecords) or unmarshalling (Records).
// RawRecords directly map to precomputed JSON-data from the database for performance reasons.
type ReadFeedbackFeedResponse struct {
	Offset     int64             `json:"offset"`
	Limit      int64             `json:"limit"`
	Records    RecordPreviews    `json:"records"`
	RawRecords []json.RawMessage `json:"-"`
}

// *** EXPERT-BOOKMARKS ***********************************

// CreateExpertBookmarkRequest defines a request of a user to bookmark an expert.
type CreateExpertBookmarkRequest struct {
	ExpertId int64 `json:"expert_id"`
}

// DeleteExpertBookmarkRequest defines a request of a user to delete an expert from the bookmarks.
type DeleteExpertBookmarkRequest struct {
	ExpertId int64 `json:"expert_id"`
}

// ReadExpertBookmarksResponse defines a response returning all the experts that were bookmarked by a user.
// The fields Bookmarks and RawBookmarks are used for either marshalling (RawBookmarks) or unmarshalling (Bookmarks).
// RawBookmarks directly map to precomputed JSON-data from the database for performance reasons.
type ReadExpertBookmarksResponse struct {
	Bookmarks    ExpertBookmarks   `json:"bookmarks"`
	RawBookmarks []json.RawMessage `json:"-"`
}

// *** RECORD-BOOKMARKS ***********************************

// UpdateRecordBookmarkCollectionsRequest defines a request of a user to first remove a record from all his collections
// and then to add this record to the specified collections only.
type UpdateRecordBookmarkCollectionsRequest struct {
	RecordId      int64   `json:"record_id"`
	CollectionIds []int64 `json:"collection_ids"`
}

// CreateRecordBookmarkRequest defines a request of a user to bookmark a record.
type CreateRecordBookmarkRequest struct {
	RecordId int64 `json:"record_id"`
}

// DeleteRecordBookmarkRequest defines a request of a user to delete a record from his bookmarks.
type DeleteRecordBookmarkRequest struct {
	RecordId int64 `json:"record_id"`
}

// ReadRecordBookmarksResponse defines a response returning all the publication records that were bookmarked by a user.
// The fields Bookmarks and RawBookmarks are used for either marshalling (RawBookmarks) or unmarshalling (Bookmarks).
// RawBookmarks directly map to precomputed JSON-data from the database for performance reasons.
type ReadRecordBookmarksResponse struct {
	Bookmarks    RecordBookmarks   `json:"bookmarks"`
	RawBookmarks []json.RawMessage `json:"-"`
}

// *** COLLECTIONS *********************************

// CreateRecordBookmarkRequest defines a request of a user to create a new, named collection for publication bookmarks.
type CreateCollectionRequest struct {
	Title string `json:"name"`
}

// CreateCollectionResponse defines a response that acknowledges that a new named collection was added to a user's profile.
// On success the response gives the ID for the newly created collection.
type CreateCollectionResponse struct {
	CollectionId int64 `json:"collection_id"`
}

// CreateRecordBookmarkRequest defines a request of a user to delete a named collection and all related publication bookmarks.
type DeleteCollectionRequest struct {
	CollectionId int64 `json:"collection_id"`
}

// UpdateCollectionRequest defines a request of a user to rename a collection of publication bookmarks.
type UpdateCollectionRequest struct {
	CollectionId int64  `json:"collection_id"`
	Title        string `json:"name"`
}

// ReadCollectionsResponse defines a response that returns all the named collections of a user.
type ReadCollectionsResponse struct {
	Collections Collections `json:"collections"`
}

// *** FEEDBACK *******************************************

// CreateFeedbackRequest defines a request of a user in the role of a domain expert to add feedback for a specific publication.
// A user in the role of an expert must register as an expert (via ORCiD) before he can provide feedback to a publication.
// The feedback consists of binary flags for high relevance, high quality of presentation, and sound methodology (0=no,1=yes).
type CreateFeedbackRequest struct {
	RecordId     int64 `json:"record_id"`
	Relevance    int64 `json:"relevance"`
	Presentation int64 `json:"presentation"`
	Methodology  int64 `json:"methodology"`
}

// ReadFeedbackRequest defines a request of a user to return all the feedback for a specific publication.
type ReadFeedbackRequest struct {
	RecordId int64 `json:"record_id"`
}

// ReadFeedbackResponse defines a response that returns all the expert feedbacks for a specific publication.
type ReadFeedbackResponse struct {
	Feedbacks Feedbacks `json:"feedbacks"`
}

// *** PERSONALIZATION ************************************

// CreateRecordDislikeRequest defines a request of a user to mark a specific publication as uninteresting.
// This information is used to improve the user's publication and expert feed composition, by removing similar content.
type CreateRecordDislikeRequest struct {
	RecordId int64 `json:"record_id"`
}

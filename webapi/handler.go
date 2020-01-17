package webapi

import (
	"log"
	"net/http"
	"os"
	"unicode/utf8"
)

import (
	"github.com/fzi-forschungszentrum-informatik/dream-gozer/model"
	"github.com/fzi-forschungszentrum-informatik/dream-gozer/model/ploc"
)

// createCollection is a Web request handler that creates a new collection.
// The user specifies the authorized user profile to which this operation is related.
func (c *Context) createCollection(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.CreateCollectionRequest
	var response ploc.CreateCollectionResponse

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Update database

	collectionId, err := c.db.CreateCollection(u.Id, request.Title)
	if err != nil {
		handleInternalError(w, "Database error. Could not create collection.", err)
		return
	}

	// Build response

	response.CollectionId = collectionId

	// Respond

	writeResponse(w, response)
}

// createExpertBookmark is a Web request handler that bookmarks an expert.
// The user specifies the authorized user profile to which this operation is related.
func (c *Context) createExpertBookmark(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.CreateExpertBookmarkRequest

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Update database

	err := c.db.CreateExpertBookmark(u.Id, request.ExpertId)
	if err != nil {
		handleInternalError(w, "Database error. Could not bookmark expert.", err)
		return
	}

	// Write response (no payload)

	w.WriteHeader(http.StatusOK)
}

// createExpertProfile is a Web request handler that registers a user as an expert.
func (c *Context) createExpertProfile(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.CreateExpertProfileRequest

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// TODO: Check ORCiD identity via OAuth2.

	// Update database

	err := c.db.CreateExpertProfile(u.Id, request.OrcId)
	if err != nil {
		handleInternalError(w, "Database error. Could create expert profile.", err)
		return
	}

	// Write response (no payload)

	w.WriteHeader(http.StatusOK)
}

// createFeedback is a Web request handler that adds a user's feedback about a record.
// The feedback is also made public by writing it to a public distributed ledger.
// A unique bibliographic hash is used to address a record in the ledger.
func (c *Context) createFeedback(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.CreateFeedbackRequest

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Update database

	err := c.db.CreateFeedback(u.Id, request.RecordId, request.Relevance, request.Presentation, request.Methodology)
	if err != nil {
		handleInternalError(w, "Database error. Could create feedback.", err)
		return
	}

	// Add feedback publically to ethereum blockchain

	bibHash, err := c.db.ReadBibHashByRecordId(request.RecordId)
	if err != nil {
		handleInternalError(w, "Database error. Could not determine bibliographic hash for specified record.", err)
		return
	}

	// Write response (no payload)

	w.WriteHeader(http.StatusOK)

	// Write feedback publically to the ledger.
	// TODO: Implement failover strategy if writing to the ledger fails (e.g. via producer-consumer).

	if c.ledger != nil {
		err = c.ledger.AddFeedback(u.OrcId, bibHash, uint8(request.Relevance), uint8(request.Presentation), uint8(request.Methodology))
		if err != nil {
			log.Printf("Could not add feedback to the ledger. %s", err)
		}
	}
}

// createInterest is a Web request handler that defines a user's interest in a specific subject.
func (c *Context) createInterest(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.CreateInterestRequest
	var response ploc.CreateInterestResponse

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Update database

	err := c.db.CreateInterest(u.Id, request.SubjectId)
	if err != nil {
		handleInternalError(w, "Could not add user's interest to database.", err)
		return
	}

	// Build response

	recordCount, err := c.db.ReadRecordFeedCount(u.Id)
	if err != nil {
		handleInternalError(w, "Could not read the number of records matching the user's interests from database.", err)
		return
	}

	response.RecordCount = recordCount

	// Respond

	writeResponse(w, response)
}

// createRecordBookmark is a Web request handler that bookmarks a specific record.
// The user specifies the authorized user profile to which this operation is related.
func (c *Context) createRecordBookmark(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.CreateRecordBookmarkRequest

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Update database

	err := c.db.CreateRecordBookmark(u.Id, request.RecordId)
	if err != nil {
		handleInternalError(w, "Database error. Could not bookmark record.", err)
		return
	}

	// Write response (no payload)

	w.WriteHeader(http.StatusOK)
}

// createRecordDislike is a Web request handler that defines a user's disinterest in a specific record.
// Such records are of no interest and should be avoided in future search results.
func (c *Context) createRecordDislike(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.CreateRecordDislikeRequest

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Update database

	err := c.db.CreateRecordDislike(u.Id, request.RecordId)
	if err != nil {
		handleInternalError(w, "Database error. Could not mark record as disliked.", err)
		return
	}

	// Write response (no payload)

	w.WriteHeader(http.StatusOK)
}

// createUserProfile is a Web request handler that registers a new user by creating a new user profile.
// Access to profile is controlled by a user-specified, secret passphrase, or its hash in specific.
func (c *Context) createUserProfile(w http.ResponseWriter, r *http.Request) {

	// Declare request and response data structures

	var request ploc.CreateUserProfileRequest
	var response ploc.CreateUserProfileResponse

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Check input types

	// TODO: Check request for malicious or malformed data values.
	// TODO: Better check entropy instead of password length.
	if utf8.RuneCountInString(request.Secret) < 16 {
		handleBadRequest(w, "Secret phrase is too short. Must contain at least 16 characters.")
		return
	}

	// Update database

	user, err := model.NewUserWithSecret(request.Secret)
	if err != nil {
		handleInternalError(w, "Could not create and initialize new user data structure.", err)
		return
	}

	err = c.db.CreateUser(&user)
	if err != nil {
		handleInternalError(w, "Could not add new user to database.", err)
		return
	}

	// Build response

	response.GUID = user.GUID

	// Respond

	writeResponse(w, response)
}

// deleteCollection is a Web request handler that deletes a user's collection and all related bookmarks in that collection.
func (c *Context) deleteCollection(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.DeleteCollectionRequest

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Update database

	if err := c.db.DeleteCollection(u.Id, request.CollectionId); err != nil {
		handleInternalError(w, "Could not delete collection.", err)
		return
	}

	// Write response (no payload)

	w.WriteHeader(http.StatusOK)
}

// deleteExpertBookmark is a Web request handler that deletes a bookmarked expert from a user's profile.
func (c *Context) deleteExpertBookmark(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.DeleteExpertBookmarkRequest

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Update database

	err := c.db.DeleteExpertBookmark(u.Id, request.ExpertId)
	if err != nil {
		handleInternalError(w, "Could not delete expert from bookmark list.", err)
		return
	}

	// Write response (no payload)

	w.WriteHeader(http.StatusOK)
}

// deleteExpertProfile is a Web request handler that withdraws a user's expert status.
func (c *Context) deleteExpertProfile(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Update database

	err := c.db.DeleteOrcId(u.Id)
	if err != nil {
		handleInternalError(w, "Could not delete user's ORCiD.", err)
		return
	}

	// Write response (no payload)

	w.WriteHeader(http.StatusOK)
}

// deleteInterest is a Web request handler that removes a single subject of interest from a user's profile.
func (c *Context) deleteInterest(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.DeleteInterestRequest
	var response ploc.DeleteInterestResponse

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Update database

	err := c.db.DeleteInterest(u.Id, request.SubjectId)
	if err != nil {
		handleInternalError(w, "Could not delete user's interest from database.", err)
		return
	}

	// Build response

	recordCount, err := c.db.ReadRecordFeedCount(u.Id)
	if err != nil {
		handleInternalError(w, "Could not read the number of records matching the user's interests from database.", err)
		return
	}

	response.RecordCount = recordCount

	// Respond

	writeResponse(w, response)
}

// deleteRecordBookmark is a Web request handler that removes a bookmarked record from a user's profile.
func (c *Context) deleteRecordBookmark(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.DeleteRecordBookmarkRequest

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Update database

	err := c.db.DeleteRecordBookmark(u.Id, request.RecordId)
	if err != nil {
		handleInternalError(w, "Could not delete record from bookmark list.", err)
		return
	}

	// Write response (no payload)

	w.WriteHeader(http.StatusOK)
}

// deleteUserProfile is a Web request handler that removes a user's profile and all user-related information.
// Feedback of that user is deleted locally but not in the public ledger.
func (c *Context) deleteUserProfile(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Delete authorized user

	err := c.db.DeleteUserById(u.Id)
	if err != nil {
		handleInternalError(w, "Internal database error while deleting user.", err)
		return
	}

	// Write response (no payload)

	w.WriteHeader(http.StatusOK)
}

// downloadPloc is a Web request handler that allows for downloading the ploc app.
// With this function a user can download the mobile app directly to its phone and install it without any app store.
func (c *Context) downloadPloc(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/vnd.android.package-archive")
	w.Header().Set("Content-Disposition", `attachment; filename="ploc.apk"`)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1
	w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0
	w.Header().Set("Expires", "0")

	fi, err := os.Stat(c.conf.PlocAPK)
	if err != nil {
		handleInternalError(w, "Could not read 'ploc.apk' file stats.", err)
		return
	} else {
		log.Printf("ploc.apk size: %d bytes", fi.Size())
	}

	http.ServeFile(w, r, c.conf.PlocAPK)
}

// readCollections is a Web request handler that returns all the bookmark collections that are stored in an user's profile.
func (c *Context) readCollections(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var response ploc.ReadCollectionsResponse

	// Build response

	collections, err := c.db.ReadCollections(u.Id)
	if err != nil {
		handleInternalError(w, "Could not read collections from database.", err)
		return
	}

	response.Collections = collections

	// Respond

	writeResponse(w, response)
}

// readExpertBookmarks is a Web request handler that returns all the experts that were bookmarked by a user.
func (c *Context) readExpertBookmarks(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var response ploc.ReadExpertBookmarksResponse

	// Build response

	rawBookmarks, err := c.db.ReadExpertBookmarks(u.Id)
	if err != nil {
		handleInternalError(w, "Could not read bookmarked experts from database.", err)
		return
	}

	response.RawBookmarks = rawBookmarks

	// Respond

	writeResponse(w, response)
}

// readExpertDetails is a Web request handler that returns a detailed profile about a specific expert.
// The profile includes information like name, ORCiD and publications.
func (c *Context) readExpertDetails(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.ReadExpertDetailsRequest
	var response ploc.ReadExpertDetailsResponse

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Read record details as precomputed JSON from the database

	rawDetails, err := c.db.ReadExpertDetails(request.ExpertId)
	if err != nil {
		handleInternalError(w, "Database error. Could not read expert details.", err)
		return
	}

	// Build response

	response.RawDetails = rawDetails

	// Respond

	writeResponse(w, response)
}

// readExpertFeed is a Web request handler that returns a specified subsegment of a list with experts that may be interesting to a user.
// The list is descendingly ordered by the relevance of the experts. The returned list depend on the user's defined interests and dislikes.
// Access to the full list is handled in subsegments via offset and limit, so that a client can read only the segments that are shown to the user.
func (c *Context) readExpertFeed(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.ReadExpertFeedRequest
	var response ploc.ReadExpertFeedResponse

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Read records from database

	rawExperts, err := c.db.ReadExpertFeed(u.Id, request.Offset, request.Limit)
	if err != nil {
		handleInternalError(w, "Database error. Could not read expert feed.", err)
		return
	}

	// Build response

	response.Offset = request.Offset
	response.Limit = request.Limit
	response.RawExperts = rawExperts

	// Respond

	writeResponse(w, response)
}

// readExpertProfile is a Web request handler that returns a short summary about a specific expert.
// The summary includes information like name, last year of publication, and subjects the expert has published about.
func (c *Context) readExpertProfile(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var response ploc.ReadExpertProfileResponse

	// Build response

	orcId, err := c.db.ReadOrcId(u.Id)
	if err != nil {
		handleInternalError(w, "Could not read the user's ORCiD from database.", err)
		return
	}

	response.OrcId = orcId

	// Respond

	writeResponse(w, response)
}

// readFeedback is a Web request handler that returns all feedback related to a specific publication record.
func (c *Context) readFeedback(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.ReadFeedbackRequest
	var response ploc.ReadFeedbackResponse

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Read records from database

	feedbacks, err := c.db.ReadFeedback(request.RecordId)
	if err != nil {
		handleInternalError(w, "Database error. Could not read the feedback for the specified publication.", err)
		return
	}

	// Build response

	response.Feedbacks = feedbacks

	// Respond

	writeResponse(w, response)
}

// readFeedbackFeed is a Web request handler that returns a list of publications that a user is expert of and that the user may provide feedback to.
// The list is descendingly ordered by the number of matching subjects of interest.
// Access to the full list is handled in subsegments, so that a client can read only the segments that are shown to the user, but not the whole list.
func (c *Context) readFeedbackFeed(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.ReadFeedbackFeedRequest
	var response ploc.ReadFeedbackFeedResponse

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Read records from database

	rawRecords, err := c.db.ReadFeedbackFeed(u.Id, request.Offset, request.Limit)
	if err != nil {
		handleInternalError(w, "Database error. Could not read record feed.", err)
		return
	}

	// Build response

	response.Offset = request.Offset
	response.Limit = request.Limit
	response.RawRecords = rawRecords

	// Respond

	writeResponse(w, response)
}

// readInterests is a Web request handler that returns a list of all the subjects that a user has specified as interesting.
func (c *Context) readInterests(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var response ploc.ReadInterestsResponse

	// Update database

	subjects, err := c.db.ReadUserInterests(u.Id)
	if err != nil {
		handleInternalError(w, "Could not read user's interests from database.", err)
		return
	}

	// Build response

	recordCount, err := c.db.ReadRecordFeedCount(u.Id)
	if err != nil {
		handleInternalError(w, "Could not read the number of records matching the user's interests from database.", err)
		return
	}

	response.Subjects = subjects
	response.RecordCount = recordCount

	// Respond

	writeResponse(w, response)
}

// readRecordBookmarks is a Web request handler that returns all records bookmarked by a user.
func (c *Context) readRecordBookmarks(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var response ploc.ReadRecordBookmarksResponse

	// Build response

	rawBookmarks, err := c.db.ReadRecordBookmarks(u.Id)
	if err != nil {
		handleInternalError(w, "Could not read bookmarked records from database.", err)
		return
	}

	response.RawBookmarks = rawBookmarks

	// Respond

	writeResponse(w, response)
}

// readRecordDetails is a Web request handler that returns detailed information about a specific publication record.
// The detailed information includes information like title, creator names, keywords, DOI and links.
func (c *Context) readRecordDetails(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.ReadRecordDetailsRequest
	var response ploc.ReadRecordDetailsResponse

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Read record details as precomputed JSON from the database

	rawDetails, err := c.db.ReadRecordDetails(u.Id, request.RecordId)
	if err != nil {
		handleInternalError(w, "Database error. Could not read record details.", err)
		return
	}

	// Build response

	response.RawDetails = rawDetails

	// Respond

	writeResponse(w, response)
}

// readRecordFeed is a Web request handler that returns a list of publications that match the user's subjects of interest.
// The list is descendingly ordered by year of publication and within that year by relevance.
// The list is accessed segment-wise, so that a client can read only the segments that are shown to the user, but not the whole list.
func (c *Context) readRecordFeed(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.ReadRecordFeedRequest
	var response ploc.ReadRecordFeedResponse

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Read records from database

	rawRecords, err := c.db.ReadRecordFeed(u.Id, request.Offset, request.Limit)
	if err != nil {
		handleInternalError(w, "Database error. Could not read record feed.", err)
		return
	}

	// Build response

	response.Offset = request.Offset
	response.Limit = request.Limit
	response.RawRecords = rawRecords

	// Respond

	writeResponse(w, response)
}

// readRecordTypes is a Web request handler that returns all supported types of publications and their numerical keys.
func (c *Context) readRecordTypes(w http.ResponseWriter, r *http.Request, u *model.User) {

	var response ploc.ReadRecordTypesResponse

	// Build response

	response.Types = []ploc.RecordType{
		ploc.RecordType{Id: int64(0), Keyword: "Article"},
		ploc.RecordType{Id: int64(1), Keyword: "Book"},
		ploc.RecordType{Id: int64(2), Keyword: "Other"},
		ploc.RecordType{Id: int64(3), Keyword: "Paper"},
		ploc.RecordType{Id: int64(4), Keyword: "Report"},
		ploc.RecordType{Id: int64(5), Keyword: "Thesis"},
	}

	// Respond

	writeResponse(w, response)
}

// readSubjects is a Web request handler that returns all supported subjects from the publication database.
// These subjects can be used to define the user's interests.
func (c *Context) readSubjects(w http.ResponseWriter, r *http.Request, u *model.User) {

	var response ploc.ReadSubjectsResponse

	// Setup response

	subs, err := c.db.ReadAllSubjects()
	if err != nil {
		handleInternalError(w, "Could not read subjects from database.", err)
		return
	}

	// Build response

	response.Subjects = subs

	// Respond

	writeResponse(w, response)
}

// searchExpertFeed is a Web request handler that makes a full text search within a user's expert feed and
// returns a summary for all the experts with matching textual content. The searched fields include name,
// publication titles and subjects.
func (c *Context) searchExpertFeed(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.SearchExpertFeedRequest
	var response ploc.SearchExpertFeedResponse

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Search database (this may take a while)
	// TODO: Filter search term and remove malicious character sequences

	rawExperts, err := c.db.SearchExpertFeed(u.Id, request.SearchTerm, request.Offset, request.Limit)
	if err != nil {
		handleInternalError(w, "Database error. Could not search expert feed.", err)
		return
	}

	// Build response

	response.Offset = request.Offset
	response.Limit = request.Limit
	response.RawExperts = rawExperts

	// Respond

	writeResponse(w, response)
}

// searchRecordFeed is a Web request handler that makes a full text search within a user's publication feed and
// returns a summary for all the publication records with matching textual content. The searched fields include title,
// abstract, subjects and author names.
func (c *Context) searchRecordFeed(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.SearchRecordFeedRequest
	var response ploc.SearchRecordFeedResponse

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Search database (this may take a while)
	// TODO: Filter search term and remove malicious character sequences

	rawRecords, err := c.db.SearchRecordFeed(u.Id, request.SearchTerm, request.Offset, request.Limit)
	if err != nil {
		handleInternalError(w, "Database error. Could not search record feed.", err)
		return
	}

	// Build response

	response.Offset = request.Offset
	response.Limit = request.Limit
	response.RawRecords = rawRecords

	// Respond

	writeResponse(w, response)
}

// updateCollection is a Web request handler that allows a user to change the title for one of its existing collections.
func (c *Context) updateCollection(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.UpdateCollectionRequest

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Update database

	if err := c.db.UpdateCollection(u.Id, request.CollectionId, request.Title); err != nil {
		handleInternalError(w, "Could not update title of collection.", err)
		return
	}

	// Write response (no payload)

	w.WriteHeader(http.StatusOK)
}

// updateRecordBookmarkCollections is a Web request handler that allows a user to specify to which of its bookmark collections
// a publication corresponds. With help of this function a publication can be added or removed from any of the user's
// collections.
func (c *Context) updateRecordBookmarkCollections(w http.ResponseWriter, r *http.Request, u *model.User) {

	// Declare request and response data structures

	var request ploc.UpdateRecordBookmarkCollectionsRequest

	// Unmarshal request

	if err := readRequest(w, r, &request); err != nil {
		handleInternalError(w, "Could not decode HTTP request. Body does not seem to contain the right JSON datastructure.", err)
		return
	}

	// Update database

	if err := c.db.UpdateCollectionsBookmarkLink(u.Id, request.RecordId, request.CollectionIds); err != nil {
		handleInternalError(w, "Could not update bookmark-collections-link.", err)
		return
	}

	// Write response (no payload)

	w.WriteHeader(http.StatusOK)
}

package webapi

import (
	"testing"
)

func TestCollections(t *testing.T) {

	// Setup database and service

	ts := NewTestService(t)
	defer ts.Close()

	// Setup some example data.

	profile := ts.CreateUserProfileWithData()

	// Perform test #1: test assumptions about profile

	if len(profile.Collections) != 2 {
		t.Errorf("Expected %d collections but got %d.", 2, len(profile.Collections))
		return
	}

	bookmarkA := profile.RecordBookmarks.SelectByTitle("Contagion dynamics in EMU government bond spreads")
	bookmarkB := profile.RecordBookmarks.SelectByTitle("Subprime mortgages, foreclosures, and urban neighborhoods")

	if len(bookmarkA.CollectionIds) != 2 {
		t.Errorf("Expected bookmark A to be in %d collections but it is in %d.", 2, len(bookmarkA.CollectionIds))
		return
	}

	if len(bookmarkB.CollectionIds) != 1 {
		t.Errorf("Expected bookmark B to be in %d collections but it is in %d.", 1, len(bookmarkB.CollectionIds))
		return
	}

	// Perform test #2: update collection title

	const newTitle = "Personal"

	collectionA := profile.Collections.SelectByTitle("Private")
	ts.UpdateCollection(collectionA.Id, newTitle)
	profile.Collections = ts.ReadCollections().Collections
	collectionA2 := profile.Collections.SelectByTitle(newTitle)

	if collectionA2.Title != newTitle {
		t.Errorf("Expected collection title to be '%s' but is '%s'.", newTitle, collectionA2.Title)
		return
	}

	// Perform test #3: delete a collection

	collectionB := profile.Collections.SelectByTitle("Work")
	ts.DeleteCollection(collectionB.Id)
	profile.Collections = ts.ReadCollections().Collections

	if len(profile.Collections) != 1 {
		t.Errorf("Expected %d collections but got %d.", 1, len(profile.Collections))
		return
	}
}

func TestDislike(t *testing.T) {

	// Setup database and service

	ts := NewTestService(t)
	defer ts.Close()

	// Setup some example data.

	_ = ts.CreateUserProfileWithData()

	// Perform test #1: test assumption about the user's test profile

	recordFeed := ts.ReadRecordFeed(0, 200).Records

	if len(recordFeed) != 90 {
		t.Errorf("Expected %d records in the feed but got %d.", 90, len(recordFeed))
		return
	}

	// Perform test #2: dislike another record

	recordA := recordFeed.SelectByTitle("The failure of supervisory stress testing: Fannie Mae, Freddie Mac, and OFHEO") // ID: 15708
	ts.CreateRecordDislike(recordA.Id)

	recordFeed = ts.ReadRecordFeed(0, 200).Records

	if len(recordFeed) != 89 {
		t.Errorf("Expected %d records in the feed but got %d.", 89, len(recordFeed))
		return
	}
}

func TestExpertBookmarks(t *testing.T) {

	// Setup database and service

	ts := NewTestService(t)
	defer ts.Close()

	// Setup some example data.

	profile := ts.CreateUserProfileWithData()

	// Perform test #1: test assumptions about the user's test profile

	if len(profile.ExpertBookmarks) != 2 {
		t.Errorf("Expected %d bookmarked experts but got %d.", 2, len(profile.ExpertBookmarks))
		return
	}

	// Perform test #2: delete bookmark

	someBookmark := profile.ExpertBookmarks[0]
	ts.DeleteExpertBookmark(someBookmark.Id)
	profile.ExpertBookmarks = ts.ReadExpertBookmarks().Bookmarks

	if len(profile.ExpertBookmarks) != 1 {
		t.Errorf("Expected %d bookmarked experts but got %d.", 1, len(profile.ExpertBookmarks))
		return
	}
}

func TestExpertProfile(t *testing.T) {

	// Setup database and service

	ts := NewTestService(t)
	defer ts.Close()

	// Setup some example data.

	_ = ts.CreateUserProfileWithData()
	someOrcId := "0000-0001-5393-1422"

	// Perform test #1: Read ORCiD from expert profile

	expertProfile := ts.ReadExpertProfile()

	if len(expertProfile.OrcId) != len(someOrcId) {
		t.Errorf("Expected an OrcId of length %d but length is %d.", len(someOrcId), len(expertProfile.OrcId))
		return
	}

	// Perform test #2: Delete expert profile

	ts.DeleteExpertProfile()

	expertProfile = ts.ReadExpertProfile()

	if len(expertProfile.OrcId) != 0 {
		t.Errorf("Expected an OrcId of length %d but length is %d.", 0, len(expertProfile.OrcId))
		return
	}
}

func TestExpertFeed(t *testing.T) {

	// Setup database and service

	ts := NewTestService(t)
	defer ts.Close()

	// Setup some example data.

	_ = ts.CreateUserProfileWithData()

	// Perform test #1

	respFeed := ts.ReadExpertFeed(0, 500)

	if len(respFeed.Experts) != 191 {
		t.Errorf("Expected %d experts but got %d.", 191, len(respFeed.Experts))
		return
	}

	respFeed = ts.ReadExpertFeed(10, 10)

	if len(respFeed.Experts) != 10 {
		t.Errorf("Expected %d experts but got %d.", 10, len(respFeed.Experts))
		return
	}

	if respFeed.Offset != 10 {
		t.Errorf("Expected %d as offset but got %d.", 10, respFeed.Offset)
		return
	}

	if respFeed.Limit != 10 {
		t.Errorf("Expected %d as limit but got %d.", 10, respFeed.Limit)
		return
	}

	// Perform test #2: load record details

	respFeed = ts.ReadExpertFeed(0, 200)

	previewNameA := "P. Abbassi"
	fullNameA := "Puriya Abbassi"
	expertA := respFeed.Experts.SelectByName(previewNameA)
	respDetails := ts.ReadExpertDetails(expertA.Id)

	if respDetails.Name != fullNameA {
		t.Errorf("Expected name '%s' but got '%s'.", fullNameA, respDetails.Name)
		return
	}

	if respDetails.LastPublicationYear != 2018 {
		t.Errorf("Expected year %d but got %d.", 2018, respDetails.LastPublicationYear)
		return
	}

	if len(respDetails.Subjects) != 47 {
		t.Errorf("Expected %d subjects but got %d.", 47, len(respDetails.Subjects))
		return
	}

	// Perform test #3: search feed

	searchTerm := "Podlich" // should match 1 expert

	respSearch := ts.SearchExpertFeed(searchTerm, 0, 100)

	if len(respSearch.Experts) != 1 {
		t.Errorf("Expected %d experts to match search term, but got %d.", 1, len(respSearch.Experts))
		return
	}

	expertA = respSearch.Experts[0]

	if expertA.Name != "N. Podlich" {
		t.Errorf("Expected name '%s', but got '%s'.", "N. Podlich", expertA.Name)
		return
	}
}

func TestFeedbackFeed(t *testing.T) {

	// Setup database and service

	ts := NewTestService(t)
	defer ts.Close()

	// Setup some example data.

	_ = ts.CreateUserProfileWithData()

	// Perform test #1: Read feedback feed.

	records := ts.ReadFeedbackFeed(0, 200).Records

	expRecordCount := 92

	if len(records) != expRecordCount {
		t.Errorf("Expected %d number of records in feedback feed but got %d.", expRecordCount, len(records))
		return
	}

	// Perform test #2: Provide and then read feedback for two records.

	recordA := records[0]
	recordB := records[1]

	ts.CreateFeedback(recordA.Id, 1, 0, 0) // Rate relevance, presentation, and methodology
	ts.CreateFeedback(recordB.Id, 1, 1, 1)

	records = ts.ReadFeedbackFeed(0, 200).Records // Rated records should not occur in record-feedback feed

	expRecordCount = expRecordCount - 2

	if len(records) != expRecordCount {
		t.Errorf("Expected %d number of records in feedback feed but got %d.", expRecordCount, len(records))
		return
	}

	// Perform test #3: Read feedback provided by user.

	feedbacks := ts.ReadFeedback(recordA.Id).Feedbacks

	expFeedbackCount := 1

	if len(feedbacks) != expFeedbackCount {
		t.Errorf("Expected %d number of feedbacks for record %d but got %d.", expFeedbackCount, recordA.Id, len(feedbacks))
		return
	}

	feedbackA := feedbacks[0]

	if feedbackA.RecordId != recordA.Id || feedbackA.Relevance != 1 {
		t.Errorf("Unexpected feeback field values.")
		return
	}
}

func TestInterests(t *testing.T) {

	// Setup database and service

	ts := NewTestService(t)
	defer ts.Close()

	// Setup some example data

	profile := ts.CreateUserProfileWithData()

	// Perform test #1: test assumptions about user's test profile

	if len(profile.Interests) != 2 {
		t.Errorf("Expected 2 interests but got %d.", len(profile.Interests))
		return
	}

	// Perform test #2: test reading call and returned number of matching records

	readResp := ts.ReadInterests()

	if len(readResp.Subjects) != 2 {
		t.Errorf("Expected 2 interests but got %d.", len(readResp.Subjects))
		return
	}

	if readResp.RecordCount != 90 {
		t.Errorf("Expected 90 matching records but got %d.", readResp.RecordCount)
		return
	}

	// Perform test #3: delete an interest

	interestA := profile.Interests.SelectByKeyword("Policy")
	deleteResp := ts.DeleteInterest(interestA.Id)

	if deleteResp.RecordCount != 90 {
		t.Errorf("Expected 90 matching records but got %d.", deleteResp.RecordCount)
		return
	}
}

func TestRecordBookmarks(t *testing.T) {

	// Setup database and service

	ts := NewTestService(t)
	defer ts.Close()

	// Setup some example data.

	profile := ts.CreateUserProfileWithData()

	// Perform test #1

	if len(profile.RecordBookmarks) != 2 {
		t.Errorf("Expected %d bookmarked records but got %d.", 2, len(profile.RecordBookmarks))
		return
	}

	// Perform test #2

	someBookmark := profile.RecordBookmarks[0]
	ts.DeleteRecordBookmark(someBookmark.Id)
	profile.RecordBookmarks = ts.ReadRecordBookmarks().Bookmarks

	if len(profile.RecordBookmarks) != 1 {
		t.Errorf("Expected %d bookmarked records but got %d.", 1, len(profile.RecordBookmarks))
		return
	}
}

func TestRecordFeed(t *testing.T) {

	// Setup database and service

	ts := NewTestService(t)
	defer ts.Close()

	// Setup some example data.

	_ = ts.CreateUserProfileWithData()

	// Perform test #1

	respFeed := ts.ReadRecordFeed(0, 10)

	if len(respFeed.Records) != 10 {
		t.Errorf("Expected %d records but got %d.", 10, len(respFeed.Records))
		return
	}

	respFeed = ts.ReadRecordFeed(85, 10)

	if len(respFeed.Records) != 5 {
		t.Errorf("Expected %d records but got %d.", 5, len(respFeed.Records))
		return
	}

	if respFeed.Offset != 85 {
		t.Errorf("Expected %d as offset but got %d.", 90, respFeed.Offset)
		return
	}

	if respFeed.Limit != 10 {
		t.Errorf("Expected %d as limit but got %d.", 10, respFeed.Limit)
		return
	}

	// Perform test #2: load record details

	respFeed = ts.ReadRecordFeed(0, 100)

	titleA := "Macroeconomic imbalances in the United States and their impact on the international financial system"
	recordA := respFeed.Records.SelectByTitle(titleA)
	respDetails := ts.ReadRecordDetails(recordA.Id)

	if respDetails.Title != titleA {
		t.Errorf("Expected title '%s' but got '%s'.", titleA, respDetails.Title)
		return
	}

	if respDetails.Year != 2009 {
		t.Errorf("Expected year %d but got %d.", 2009, respDetails.Year)
		return
	}

	if len(respDetails.Subjects) != 7 {
		t.Errorf("Expected %d subjects but got %d.", 7, len(respDetails.Subjects))
		return
	}

	if recordA.Visited {
		t.Errorf("Record must not be marked as 'visited' before record details were loaded.")
		return
	}

	respFeed = ts.ReadRecordFeed(0, 100)
	recordA = respFeed.Records.SelectByTitle(titleA)

	if !recordA.Visited {
		t.Errorf("Record must be marked as 'visited' after record details were loaded.")
		return
	}

	// Perform test #3: search feed

	searchTerm := "bank" // should match 3 records

	respSearch := ts.SearchRecordFeed(searchTerm, 0, 100)

	if len(respSearch.Records) != 77 {
		t.Errorf("Expected %d records to match search term, but got %d.", 77, len(respSearch.Records))
		return
	}
}

func TestRecordTypes(t *testing.T) {

	// Setup database and service

	ts := NewTestService(t)
	defer ts.Close()

	// Setup some example data.

	ts.CreateUserProfile()

	// Perform test

	resp := ts.ReadRecordTypes()

	if len(resp.Types) != 6 {
		t.Errorf("Expected %d different record types but got %d.", 6, len(resp.Types))
		return
	}
}

func TestSubjects(t *testing.T) {

	// Setup database and service

	ts := NewTestService(t)
	defer ts.Close()

	// Setup some example data.

	ts.CreateUserProfile()

	// Perform test

	resp := ts.ReadSubjects()

	if len(resp.Subjects) != 415 {
		t.Errorf("Expected %d different subjects but got %d.", 415, len(resp.Subjects))
		return
	}
}

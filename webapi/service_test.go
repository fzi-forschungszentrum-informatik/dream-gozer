package webapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

import (
	"github.com/fzi-forschungszentrum-informatik/gozer/config"
	"github.com/fzi-forschungszentrum-informatik/gozer/model/ploc"
	"github.com/fzi-forschungszentrum-informatik/gozer/storage"
	"github.com/fzi-forschungszentrum-informatik/gozer/storage/ledger"
)

type TestService struct {
	storage *storage.Storage
	server  *httptest.Server
	t       *testing.T
	guid    string
	secret  string
}

type TestUserProfile struct {
	RecordBookmarks ploc.RecordBookmarks
	ExpertBookmarks ploc.ExpertBookmarks
	Collections     ploc.Collections
	Interests       ploc.Subjects
}

func NewTestService(t *testing.T) *TestService {

	conf := config.DefaultConfiguration()
	conf.Storage.DBFilename = ":memory:"
	conf.Ledger.Enable = false // Enable to test with local Ethereum testbed (e.g. Ganache)
	conf.Ledger.RPCClient = "http://127.0.0.1:8545"
	conf.Ledger.ContractAddress = "17e91224c30c5b0b13ba2ef1e84fe880cb902352"                    // Adress for the open feedback storage contract in the Ganache testbed.
	conf.Ledger.PrivateKey = "6370fd033278c143179d81c5526140625662b8daa446c22ee2d73db3707e620c" // Private wallet key that is used to pay transaction fees in the Ganache testbed.

	storage := storage.Open(&conf.Storage)
	ledger := ledger.Open(&conf.Ledger)
	server := httptest.NewServer(newRouter(&conf.WebAPI, storage, ledger))

	storage.CreateTestPublications()
	storage.BuildSearchIndicies()

	return &TestService{
		storage: storage,
		server:  server,
		t:       t,
		guid:    "",
		secret:  "",
	}
}

func (ts *TestService) Close() {
	ts.storage.Close()
	ts.server.Close()
}

// postTestRequest() sends a JSON-encoded post request to a test server and returns a JSON-decoded response.
// This function is for test purposes only. If no JSON datastructure should be send or received 'nil' can be used.
// A user Id and passward can also be used to send HTTP basic authentication.
func (ts *TestService) PostRequest(urlPostfix string, request interface{}, response interface{}) (statusCode int, err error) {

	var jData []byte

	// Convert request data structure to JSON.

	if request != nil {
		jData, err = json.Marshal(request)
		if err != nil {
			return 0, fmt.Errorf("Could not convert request data structure to JSON. %s", err)
		}
	}

	// Send post request with JSON payload.

	url := ts.server.URL + "/plocapi/v1" + urlPostfix

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jData))
	if request != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if ts.guid != "" {
		req.SetBasicAuth(ts.guid, ts.secret)
	}

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	// Convert JSON body data to response data structure.

	if response != nil {
		jData, _ = ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(jData, response)
		if err != nil {
			return 0, fmt.Errorf("Could not JSON '%s' to response data structure. %s", string(jData), err)
		}
	}

	return resp.StatusCode, nil
}

func (ts *TestService) PostRequestOK(urlPostfix string, request interface{}, response interface{}) {

	statusCode, err := ts.PostRequest(urlPostfix, request, response)

	if err != nil {
		ts.t.Errorf("Unexpected error. %s", err)
		return
	}

	if statusCode != http.StatusOK {
		ts.t.Errorf("Expected HTTP.StatusOK but get: %d", statusCode)
		return
	}

	return
}

func (ts *TestService) CreateCollection(title string) (response ploc.CreateCollectionResponse) {
	request := ploc.CreateCollectionRequest{Title: title}
	ts.PostRequestOK("/collection/create", &request, &response)
	return
}

func (ts *TestService) CreateExpertBookmark(expertId int64) {
	request := ploc.CreateExpertBookmarkRequest{ExpertId: expertId}
	ts.PostRequestOK("/expert-bookmark/create", &request, nil)
	return
}

func (ts *TestService) CreateExpertProfile(orcId string) {
	request := ploc.CreateExpertProfileRequest{OrcId: orcId}
	ts.PostRequestOK("/expert-profile/create", &request, nil)
	return
}

func (ts *TestService) CreateFeedback(recordId int64, relevance int64, presentation int64, methodology int64) {

	request := ploc.CreateFeedbackRequest{
		RecordId:     recordId,
		Relevance:    relevance,
		Presentation: presentation,
		Methodology:  methodology,
	}

	ts.PostRequestOK("/feedback/create", &request, nil)
	return
}

func (ts *TestService) CreateInterest(subjectId int64) (response ploc.CreateInterestResponse) {

	request := ploc.CreateInterestRequest{SubjectId: subjectId}

	ts.PostRequestOK("/interest/create", &request, &response)
	return
}

func (ts *TestService) CreateRecordBookmark(recordId int64) {
	request := ploc.CreateRecordBookmarkRequest{RecordId: recordId}
	ts.PostRequestOK("/record-bookmark/create", &request, nil)
	return
}

func (ts *TestService) CreateRecordDislike(recordId int64) {
	request := ploc.CreateRecordDislikeRequest{RecordId: recordId}
	ts.PostRequestOK("/record-dislike/create", &request, nil)
	return
}

func (ts *TestService) CreateUserProfile() {

	var response ploc.CreateUserProfileResponse

	secret := "dsj738hFs3d:Kl67jdk"
	request := ploc.CreateUserProfileRequest{Secret: secret}

	ts.PostRequestOK("/user-profile/create", &request, &response)

	ts.guid = response.GUID
	ts.secret = secret

	return
}

func (ts *TestService) CreateUserProfileWithData() (user TestUserProfile) {

	ts.CreateUserProfile()

	subs := ts.ReadSubjects().Subjects

	// Define some interests

	interestA := subs.SelectByKeyword("Financial Economics")             // expected to match 78 records
	interestB := subs.SelectByKeyword("International Financial Markets") // expected to match 20 records (6 overlapping)

	ts.CreateInterest(interestA.Id)
	ts.CreateInterest(interestB.Id)

	user.Interests = ts.ReadInterests().Subjects

	// Create some collections

	collectionIdA := ts.CreateCollection("Private").CollectionId
	collectionIdB := ts.CreateCollection("Work").CollectionId

	user.Collections = ts.ReadCollections().Collections

	// Dislike some records

	recordFeed := ts.ReadRecordFeed(0, 200).Records

	recordA := recordFeed.SelectByTitle("The Crisis and Beyond") // ID:60135
	recordB := recordFeed.SelectByTitle("AIG in hindsight")      // ID:30407

	ts.CreateRecordDislike(recordA.Id)
	ts.CreateRecordDislike(recordB.Id)

	// Bookmark some records and add them to collections

	recordC := recordFeed.SelectByTitle("Contagion dynamics in EMU government bond spreads")         // ID:37012
	recordD := recordFeed.SelectByTitle("Subprime mortgages, foreclosures, and urban neighborhoods") // ID:65954

	ts.CreateRecordBookmark(recordC.Id)
	ts.CreateRecordBookmark(recordD.Id)

	ts.UpdateRecordBookmarkCollections(recordC.Id, []int64{collectionIdA, collectionIdB})
	ts.UpdateRecordBookmarkCollections(recordD.Id, []int64{collectionIdA})

	user.RecordBookmarks = ts.ReadRecordBookmarks().Bookmarks

	// Bookmark some experts

	expertFeed := ts.ReadExpertFeed(0, 500).Experts

	expertA := expertFeed.SelectByName("M. McAleer") // ID:802
	expertB := expertFeed.SelectByName("M. Koetter") // ID: 1783

	ts.CreateExpertBookmark(expertA.Id)
	ts.CreateExpertBookmark(expertB.Id)

	user.ExpertBookmarks = ts.ReadExpertBookmarks().Bookmarks

	// Register dummy ORCiD profile to gain access to feedback feed

	someOrcId := "0000-0001-5393-1421"
	ts.CreateExpertProfile(someOrcId)

	return
}

func (ts *TestService) DeleteCollection(collectionId int64) {
	request := ploc.DeleteCollectionRequest{CollectionId: collectionId}
	ts.PostRequestOK("/collection/delete", &request, nil)
	return
}

func (ts *TestService) DeleteExpertBookmark(expertId int64) {
	request := ploc.DeleteExpertBookmarkRequest{ExpertId: expertId}
	ts.PostRequestOK("/expert-bookmark/delete", &request, nil)
	return
}

func (ts *TestService) DeleteExpertProfile() {
	ts.PostRequestOK("/expert-profile/delete", nil, nil)
	return
}

func (ts *TestService) DeleteInterest(subjectId int64) (response ploc.DeleteInterestResponse) {
	request := ploc.DeleteInterestRequest{SubjectId: subjectId}
	ts.PostRequestOK("/interest/delete", &request, &response)
	return
}

func (ts *TestService) DeleteRecordBookmark(recordId int64) {
	request := ploc.DeleteRecordBookmarkRequest{RecordId: recordId}
	ts.PostRequestOK("/record-bookmark/delete", &request, nil)
	return
}

func (ts *TestService) DeleteUserProfile() {
	ts.PostRequestOK("/user-profile/delete", nil, nil)
	return
}

func (ts *TestService) ReadCollections() (response ploc.ReadCollectionsResponse) {
	ts.PostRequestOK("/collections/read", nil, &response)
	return
}

func (ts *TestService) ReadExpertBookmarks() (response ploc.ReadExpertBookmarksResponse) {
	ts.PostRequestOK("/expert-bookmarks/read", nil, &response)
	return
}

func (ts *TestService) ReadExpertDetails(expertId int64) (response ploc.ReadExpertDetailsResponse) {
	request := ploc.ReadExpertDetailsRequest{ExpertId: expertId}
	ts.PostRequestOK("/expert-details/read", &request, &response)
	return
}

func (ts *TestService) ReadExpertFeed(offset int64, limit int64) (response ploc.ReadExpertFeedResponse) {
	request := ploc.ReadExpertFeedRequest{Offset: offset, Limit: limit}
	ts.PostRequestOK("/expert-feed/read", &request, &response)
	return
}

func (ts *TestService) ReadExpertProfile() (response ploc.ReadExpertProfileResponse) {
	ts.PostRequestOK("/expert-profile/read", nil, &response)
	return
}

func (ts *TestService) ReadFeedback(recordId int64) (response ploc.ReadFeedbackResponse) {
	request := ploc.ReadFeedbackRequest{RecordId: recordId}
	ts.PostRequestOK("/feedback/read", &request, &response)
	return
}

func (ts *TestService) ReadFeedbackFeed(offset int64, limit int64) (response ploc.ReadFeedbackFeedResponse) {
	request := ploc.ReadFeedbackFeedRequest{Offset: offset, Limit: limit}
	ts.PostRequestOK("/feedback-feed/read", &request, &response)
	return
}

func (ts *TestService) ReadInterests() (response ploc.ReadInterestsResponse) {
	ts.PostRequestOK("/interests/read", nil, &response)
	return
}

func (ts *TestService) ReadRecordBookmarks() (response ploc.ReadRecordBookmarksResponse) {
	ts.PostRequestOK("/record-bookmarks/read", nil, &response)
	return
}

func (ts *TestService) ReadRecordDetails(recordId int64) (response ploc.ReadRecordDetailsResponse) {
	request := ploc.ReadRecordDetailsRequest{RecordId: recordId}
	ts.PostRequestOK("/record-details/read", &request, &response)
	return
}

func (ts *TestService) ReadRecordFeed(offset int64, limit int64) (response ploc.ReadRecordFeedResponse) {
	request := ploc.ReadRecordFeedRequest{Offset: offset, Limit: limit}
	ts.PostRequestOK("/record-feed/read", &request, &response)
	return
}

func (ts *TestService) ReadRecordTypes() (response ploc.ReadRecordTypesResponse) {
	ts.PostRequestOK("/record-types/read", nil, &response)
	return
}

func (ts *TestService) ReadSubjects() (response ploc.ReadSubjectsResponse) {
	ts.PostRequestOK("/subjects/read", nil, &response)
	return
}

func (ts *TestService) SearchExpertFeed(searchTerm string, offset int64, limit int64) (response ploc.SearchExpertFeedResponse) {
	request := ploc.SearchExpertFeedRequest{SearchTerm: searchTerm, Offset: offset, Limit: limit}
	ts.PostRequestOK("/expert-feed/search", &request, &response)
	return
}

func (ts *TestService) SearchRecordFeed(searchTerm string, offset int64, limit int64) (response ploc.SearchRecordFeedResponse) {
	request := ploc.SearchRecordFeedRequest{SearchTerm: searchTerm, Offset: offset, Limit: limit}
	ts.PostRequestOK("/record-feed/search", &request, &response)
	return
}

func (ts *TestService) UpdateCollection(collectionId int64, title string) {
	request := ploc.UpdateCollectionRequest{CollectionId: collectionId, Title: title}
	ts.PostRequestOK("/collection/update", &request, nil)
	return
}

func (ts *TestService) UpdateRecordBookmarkCollections(recordId int64, collectionIds []int64) {
	request := ploc.UpdateRecordBookmarkCollectionsRequest{RecordId: recordId, CollectionIds: collectionIds}
	ts.PostRequestOK("/record-bookmark/collections/update", &request, nil)
	return
}

package webapi

import (
	"net/http"
)

import (
	"github.com/fzi-forschungszentrum-informatik/gozer/config"
	"github.com/fzi-forschungszentrum-informatik/gozer/storage"
	"github.com/fzi-forschungszentrum-informatik/gozer/storage/ledger"
	"github.com/gorilla/mux"
)

// newRouter creates a HTTP request router and dispatcher that maps all incoming HTTP request to their
// corresponding request handlers.
func newRouter(conf *config.WebAPIConfiguration, st *storage.Storage, ledger *ledger.Ledger) (handler http.Handler) {

	router := mux.NewRouter()
	context := newContext(conf, st, ledger)

	// API request handler
	plocRouter := router.PathPrefix("/plocapi/v1/").Subrouter()

	// User profile
	plocRouter.HandleFunc("/user-profile/create", defaultHandler(context.createUserProfile)).Methods("POST")
	plocRouter.HandleFunc("/user-profile/delete", authorizationHandler(context.deleteUserProfile, st)).Methods("POST")

	// Expert profile
	plocRouter.HandleFunc("/expert-profile/create", authorizationHandler(context.createExpertProfile, st)).Methods("POST")
	plocRouter.HandleFunc("/expert-profile/delete", authorizationHandler(context.deleteExpertProfile, st)).Methods("POST")
	plocRouter.HandleFunc("/expert-profile/read", authorizationHandler(context.readExpertProfile, st)).Methods("POST")

	// Subjects
	plocRouter.HandleFunc("/subjects/read", authorizationHandler(context.readSubjects, st)).Methods("POST")

	// Interests
	plocRouter.HandleFunc("/interest/create", authorizationHandler(context.createInterest, st)).Methods("POST")
	plocRouter.HandleFunc("/interest/delete", authorizationHandler(context.deleteInterest, st)).Methods("POST")
	plocRouter.HandleFunc("/interests/read", authorizationHandler(context.readInterests, st)).Methods("POST")

	// Record-Types
	plocRouter.HandleFunc("/record-types/read", authorizationHandler(context.readRecordTypes, st)).Methods("POST")

	// Record-Feed
	plocRouter.HandleFunc("/record-feed/read", authorizationHandler(context.readRecordFeed, st)).Methods("POST")
	plocRouter.HandleFunc("/record-feed/search", authorizationHandler(context.searchRecordFeed, st)).Methods("POST")
	plocRouter.HandleFunc("/record-details/read", authorizationHandler(context.readRecordDetails, st)).Methods("POST")

	// Expert-Feed
	plocRouter.HandleFunc("/expert-feed/read", authorizationHandler(context.readExpertFeed, st)).Methods("POST")
	plocRouter.HandleFunc("/expert-feed/search", authorizationHandler(context.searchExpertFeed, st)).Methods("POST")
	plocRouter.HandleFunc("/expert-details/read", authorizationHandler(context.readExpertDetails, st)).Methods("POST")

	// Feedback-Feed
	plocRouter.HandleFunc("/feedback-feed/read", authorizationHandler(context.readFeedbackFeed, st)).Methods("POST")

	// Record-Bookmarks
	plocRouter.HandleFunc("/record-bookmark/collections/update", authorizationHandler(context.updateRecordBookmarkCollections, st)).Methods("POST")
	plocRouter.HandleFunc("/record-bookmark/create", authorizationHandler(context.createRecordBookmark, st)).Methods("POST")
	plocRouter.HandleFunc("/record-bookmark/delete", authorizationHandler(context.deleteRecordBookmark, st)).Methods("POST")
	plocRouter.HandleFunc("/record-bookmarks/read", authorizationHandler(context.readRecordBookmarks, st)).Methods("POST")

	// Expert-Bookmarks
	plocRouter.HandleFunc("/expert-bookmark/create", authorizationHandler(context.createExpertBookmark, st)).Methods("POST")
	plocRouter.HandleFunc("/expert-bookmark/delete", authorizationHandler(context.deleteExpertBookmark, st)).Methods("POST")
	plocRouter.HandleFunc("/expert-bookmarks/read", authorizationHandler(context.readExpertBookmarks, st)).Methods("POST")

	// Collections
	plocRouter.HandleFunc("/collection/create", authorizationHandler(context.createCollection, st)).Methods("POST")
	plocRouter.HandleFunc("/collection/delete", authorizationHandler(context.deleteCollection, st)).Methods("POST")
	plocRouter.HandleFunc("/collection/update", authorizationHandler(context.updateCollection, st)).Methods("POST")
	plocRouter.HandleFunc("/collections/read", authorizationHandler(context.readCollections, st)).Methods("POST")

	// Feedback
	plocRouter.HandleFunc("/feedback/create", authorizationHandler(context.createFeedback, st)).Methods("POST")
	plocRouter.HandleFunc("/feedback/read", authorizationHandler(context.readFeedback, st)).Methods("POST")

	// Personalization
	plocRouter.HandleFunc("/record-dislike/create", authorizationHandler(context.createRecordDislike, st)).Methods("POST")

	// Download request handler
	downloadRouter := router.PathPrefix("/download").Subrouter()

	// App download
	downloadRouter.HandleFunc("/ploc", defaultHandler(context.downloadPloc)).Methods("GET")

	return router
}

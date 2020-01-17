package storage

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"
)

import (
	"github.com/fzi-forschungszentrum-informatik/gozer/model"
	"github.com/fzi-forschungszentrum-informatik/gozer/model/ploc"
	_ "github.com/mattn/go-sqlite3"
)

// BuildSearchIndicies repopulates the search indices for publication and experts.
// Must be called after new records were added to the database.
func (st *Storage) BuildSearchIndicies() (err error) {

	const recordIndexQuery = `
		INSERT INTO vrecord 
			SELECT record.id, record.title, record.abstract, rs.subjects, ra.authors
			FROM record, 
				(SELECT creator.record_id AS record_id,
					GROUP_CONCAT(creator.first_name || ' ' || creator.last_name, ' ') AS authors
					FROM creator GROUP BY creator.record_id) AS ra,
				(SELECT rsl.record_id, GROUP_CONCAT(subject.keyword, ' ') AS subjects
					FROM record_subject_link AS rsl, subject
					WHERE rsl.subject_id=subject.id
					GROUP BY rsl.record_id) AS rs
			WHERE ra.record_id=record.id AND rs.record_id=record.id`

	const expertIndexQuery = `
		INSERT INTO vexpert 
			SELECT e.id, e.first_name || ' ' || e.last_name AS full_name, et.titles, es.subjects
			FROM expert AS e,
				(SELECT creator.expert_id, GROUP_CONCAT(record.title,' ') AS titles
					FROM creator, record
					WHERE creator.record_id=record.id
					GROUP BY creator.expert_id) AS et,
				(SELECT esl.expert_id, GROUP_CONCAT(subject.keyword, ' ') AS subjects
					FROM expert_subject_link AS esl, subject
					WHERE esl.subject_id=subject.id
					GROUP BY esl.expert_id) AS es
			WHERE e.id=et.expert_id
				AND e.id=es.expert_id`

	tx, err := st.db.Begin()

	if err != nil {
		log.Printf("Database error. Could not initialize transaction for building record and expert search index. %s", err)
		return
	}

	tx.Exec("DELETE from vrecord")
	tx.Exec(recordIndexQuery)
	tx.Exec("DELETE from vexpert")
	tx.Exec(expertIndexQuery)

	err = tx.Commit()

	if err != nil {
		log.Printf("Database error. Could not commit transaction for for building record and expert search index. %s", err)
		return
	}

	return
}

// CreateCollection creates a new named bookmark collection.
// The user-ID specifies the authorized user profile to which this operation is related.
func (st *Storage) CreateCollection(uid int64, title string) (collectionId int64, err error) {

	result, err := st.db.Exec("INSERT INTO collection (user_id,title) VALUES(?,?)", uid, title)
	if err != nil {
		log.Printf("Database error. Could not insert collection. %s", err)
		return
	}

	collectionId, err = result.LastInsertId()
	if err != nil {
		log.Printf("Database error. Could not get ID for inserted collection. %s", err)
		return
	}

	return
}

// CreateExpertBookmark bookmarks an expert.
// The user-ID specifies the authorized user profile to which this operation is related.
func (st *Storage) CreateExpertBookmark(uid int64, expertId int64) (err error) {

	_, err = st.db.Exec("INSERT OR IGNORE INTO expert_bookmark (user_id,expert_id) VALUES(?,?)", uid, expertId)
	if err != nil {
		log.Printf("Database error. Could not insert expert bookmark. %s", err)
		return
	}

	return
}

// CreateExpertProfile registers an user that has an ORCiD as an expert.
func (st *Storage) CreateExpertProfile(uid int64, orcId string) (err error) {

	_, err = st.db.Exec("UPDATE user SET orcid=? WHERE id=?", orcId, uid)
	if err != nil {
		log.Printf("Database error. Could not update ORCiD for user %d. %s", uid, err)
		return
	}

	return
}

// CreateFeedback adds a user's feedback to a record. The feedback consists of binary flags for relevance,
// quality of presentation and the soundness of the methodology (0 = false, 1 = true).
func (st *Storage) CreateFeedback(uid int64, recordId int64, relevance int64, presentation int64, methodology int64) (err error) {

	q := "INSERT OR IGNORE INTO feedback (user_id,record_id,orcid,relevance,presentation,methodology) VALUES(?,?,(SELECT orcid FROM user WHERE id=?),?,?,?)"

	_, err = st.db.Exec(q, uid, recordId, uid, relevance, presentation, methodology)
	if err != nil {
		log.Printf("Database error. Could not insert feedback. %s", err)
		return
	}

	return
}

// CreateInterest defines a user's interest in a subject.
func (st *Storage) CreateInterest(uid int64, subjectId int64) (err error) {

	tx, err := st.db.Begin()

	if err != nil {
		log.Printf("Database error. Could not initialize transaction for creating and updating interests. %s", err)
		return
	}

	tx.Exec("INSERT OR IGNORE INTO interest (user_id,subject_id) VALUES(?,?)", uid, subjectId)
	st.rebuildExpertAndRecordFeed(tx, uid)

	err = tx.Commit()

	if err != nil {
		log.Printf("Database error. Could not commit transaction for creating and updating interests. %s", err)
		return
	}

	return
}

// CreateRecordBookmark bookmarks a specific record.
// The user-ID specifies the user to which bookmark is related.
func (st *Storage) CreateRecordBookmark(uid int64, recordId int64) (err error) {

	_, err = st.db.Exec("INSERT OR IGNORE INTO record_bookmark (user_id,record_id) VALUES(?,?)", uid, recordId)
	if err != nil {
		log.Printf("Database error. Could not insert record bookmark. %s", err)
		return
	}

	return
}

// CreateRecordDislike defines a user's disinterest in a specific record.
// Such records are of no interest and should be avoided in future search results.
func (st *Storage) CreateRecordDislike(uid int64, recordId int64) (err error) {

	tx, err := st.db.Begin()

	if err != nil {
		log.Printf("Database error. Could not initialize transaction for disliking a record. %s", err)
		return
	}

	tx.Exec("INSERT OR IGNORE INTO record_dislike (user_id,record_id) VALUES(?,?)", uid, recordId)
	tx.Exec("DELETE FROM record_feed WHERE user_id=? AND record_id=?", uid, recordId)

	err = tx.Commit()

	if err != nil {
		log.Printf("Database error. Could not commit transaction for disliking a record. %s", err)
		return
	}

	return
}

// CreateUser registers a new user by creating a new user-ID.
// Access to profile is controlled by a user-specified, secret passphrase, or its hash in specific.
func (st *Storage) CreateUser(u *model.User) (err error) {

	result, err := st.db.Exec("INSERT INTO user (guid,hashed_secret) VALUES(?,?);",
		u.GUID,
		u.HashedSecret)

	if err != nil {
		log.Printf("Database error. Could not insert new user into database. %s", err)
		return
	}

	u.Id, err = result.LastInsertId()
	if err != nil {
		log.Printf("Database error. Could not get ID for last inserted user. %s", err)
		return
	}

	return
}

// CreateTestPublications populates the database with test entries like records, creators and subjects.
// The test entries are used for the purpose of unit tests only.
func (st *Storage) CreateTestPublications() (err error) {

	_, err = st.db.Exec(sqlTestRecords)

	if err != nil {
		log.Printf("Database error. Could not insert test records into database. %s", err)
		return
	}

	return
}

// DeleteCollection deletes a specific bookmark collection of a user and all the bookmarks in that collection.
func (st *Storage) DeleteCollection(uid int64, collectionId int64) (err error) {

	_, err = st.db.Exec("DELETE FROM collection WHERE user_id=? AND id=?", uid, collectionId)

	if err != nil {
		log.Printf("Database error. Could not delete collection. %s", err)
		return
	}

	_, err = st.db.Exec("DELETE FROM record_bookmark WHERE user_id=? AND collection_id=?", uid, collectionId)
	if err != nil {
		log.Printf("Database error. Could not delete bookmarked record from collection. %s", err)
		return
	}

	return
}

// DeleteExpertBookmark deletes a bookmarked expert from a user's profile.
func (st *Storage) DeleteExpertBookmark(uid int64, expertId int64) (err error) {

	_, err = st.db.Exec("DELETE FROM expert_bookmark WHERE user_id=? AND expert_id=?", uid, expertId)

	if err != nil {
		log.Printf("Database error. Could not delete expert from bookmark list. %s", err)
		return
	}

	return
}

// DeleteInterest removes a single subject from the user's list of interest.
func (st *Storage) DeleteInterest(uid int64, subjectId int64) (err error) {

	tx, err := st.db.Begin()

	if err != nil {
		log.Printf("Database error. Could not initialize transaction for deleting user interest. %s", err)
		return
	}

	tx.Exec("DELETE FROM interest WHERE user_id=? AND subject_id=?", uid, subjectId)
	st.rebuildExpertAndRecordFeed(tx, uid)

	err = tx.Commit()

	if err != nil {
		log.Printf("Database error. Could not commit transaction for deleting user interest. %s", err)
		return
	}

	return
}

// DeleteOrcId withdraws a user's expert status by removing its ORCiD identity.
func (st *Storage) DeleteOrcId(uid int64) (err error) {

	_, err = st.db.Exec("UPDATE user SET orcid=NULL WHERE id=?", uid)

	if err != nil {
		log.Printf("Database error. Could not delete ORCiD for user %d. %s", uid, err)
		return
	}

	return
}

// DeleteRecordBookmark removes a bookmarked record from a user's profile.
func (st *Storage) DeleteRecordBookmark(uid int64, recordId int64) (err error) {

	_, err = st.db.Exec("DELETE FROM record_bookmark WHERE user_id=? AND record_id=?", uid, recordId)

	if err != nil {
		log.Printf("Database error. Could not delete record from bookmark list. %s", err)
		return
	}

	return
}

// DeleteUserById removes a user's identity and all user-related information.
// Note that while even the user's feedback is removed from local database, it still remains in the public ledger.
func (st *Storage) DeleteUserById(uid int64) (err error) {

	tx, err := st.db.Begin()

	if err != nil {
		log.Printf("Database error. Could not initialize transaction for deleting user profile. %s", err)
		return
	}

	tx.Exec("DELETE FROM bookmark WHERE user_id=?", uid)
	tx.Exec("DELETE FROM record_feed WHERE user_id=?", uid)
	tx.Exec("DELETE FROM expert_feed WHERE user_id=?", uid)
	tx.Exec("DELETE FROM record_bookmark WHERE user_id=?", uid)
	tx.Exec("DELETE FROM expert_bookmark WHERE user_id=?", uid)
	tx.Exec("DELETE FROM collection WHERE user_id=?", uid)
	tx.Exec("DELETE FROM interest WHERE user_id=?", uid)
	tx.Exec("DELETE FROM feedback WHERE user_id=?", uid)
	tx.Exec("DELETE FROM record_dislike WHERE user_id=?", uid)
	tx.Exec("DELETE FROM record_visit WHERE user_id=?", uid)
	tx.Exec("DELETE FROM user WHERE id=?", uid)

	err = tx.Commit()

	if err != nil {
		log.Printf("Database error. Could not commit transaction for deleting user profile. %s", err)
		return
	}

	return nil
}

// ReadAllSubjects returns all supported subjects from the publication database.
// These subjects can be used to define the user's interests.
func (st *Storage) ReadAllSubjects() (subs ploc.Subjects, err error) {

	rows, err := st.db.Query("SELECT id,keyword FROM subject")
	if err != nil {
		log.Printf("Database error. Querying subjects failed. %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		var s ploc.Subject

		err = rows.Scan(&s.Id, &s.Keyword)
		if err != nil {
			log.Printf("Database error. Scanning subjects failed. %s", err)
			return
		}

		subs = append(subs, s)
	}

	return
}

// ReadBibHashByRecordId reads the bibliographic hash for the specified record.
// On error the returned bibliographic hash is undefinied.
func (st *Storage) ReadBibHashByRecordId(recordId int64) (bibHash string, err error) {

	err = st.db.QueryRow(`SELECT bib_hash FROM record WHERE id=?`, recordId).Scan(&bibHash)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("Database error. Could not read record's bibliographic hash. Record with ID %d seems not to exist.", recordId)
	case err != nil:
		log.Printf("Database error. Could not read record's bibliographic hash. %s", err)
	}

	return
}

// ReadCollections returns all the bookmark collections without the records for a user.
func (st *Storage) ReadCollections(uid int64) (collections ploc.Collections, err error) {

	rows, err := st.db.Query("SELECT id,title FROM collection WHERE user_id=?", uid)
	if err != nil {
		log.Printf("Database error. Reading collections failed. %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		var c ploc.Collection

		err = rows.Scan(&c.Id, &c.Title)
		if err != nil {
			log.Printf("Database error. Scanning collections failed. %s", err)
			return
		}

		collections = append(collections, c)
	}

	return
}

// ReadExpertBookmarks returns all the experts that a user has bookmarked.
func (st *Storage) ReadExpertBookmarks(uid int64) (rawBookmarks []json.RawMessage, err error) {

	const query = `
		SELECT e.json_bookmark
		FROM expert AS e, expert_bookmark AS b
		WHERE b.user_id = ? AND b.expert_id = e.id
		ORDER BY b.rowid DESC`

	rows, err := st.db.Query(query, uid)
	if err != nil {
		log.Printf("Database error. Could not read expert bookmarks. %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		var rawBookmark string

		err = rows.Scan(&rawBookmark)
		if err != nil {
			log.Printf("Database error. Scanning raw JSON expert bookmarks failed. %s", err)
			return
		}

		rawBookmarks = append(rawBookmarks, json.RawMessage(rawBookmark))
	}

	return
}

// ReadExpertDetails returns a detailed profile about a specific expert.
// The profile includes information like name, ORCiD and publications.
// The profile is returned as a precomputed JSON data structure for performance reasons.
func (st *Storage) ReadExpertDetails(expertId int64) (rawDetails json.RawMessage, err error) {

	var precomputedDetails string

	err = st.db.QueryRow(`SELECT json_detail FROM expert WHERE id=?`, expertId).Scan(&precomputedDetails)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("Database error. Could not read expert details. Expert with ID %d seems not to exist.", expertId)
	case err != nil:
		log.Printf("Database error. Could not read expert details. %s", err)
	}

	return json.RawMessage(precomputedDetails), err
}

// ReadExpertFeed returns the specified subsegment of a list with popular experts that match a user's interest.
// The experts' short profiles are returned as a precomputed JSON data structure for performance reasons.
// The list is descendingly ordered by the relevance of the experts. The returned list depend on the user's defined interests and dislikes.
// Access to the full list is handled in subsegments via offset and limit, so that a client can read only the segments that are shown to the user.
func (st *Storage) ReadExpertFeed(uid int64, offset int64, limit int64) (rawExperts []json.RawMessage, err error) {

	const query = `
		SELECT e.json_preview 
		FROM expert AS e, expert_feed AS f
		WHERE f.user_id=?
			AND f.expert_id=e.id 
			AND f.expert_id NOT IN (SELECT expert_id FROM expert_bookmark WHERE user_id=?)
		ORDER BY f.rowid ASC
		LIMIT ?
		OFFSET ?`

	rows, err := st.db.Query(query, uid, uid, limit, offset)
	if err != nil {
		log.Printf("Database error. Querying expert feed failed. %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		var rawExpert string

		err = rows.Scan(&rawExpert)
		if err != nil {
			log.Printf("Database error. Scanning raw JSON expert failed. %s", err)
			return
		}

		rawExperts = append(rawExperts, json.RawMessage(rawExpert))
	}

	return
}

// ReadFeedback returns all feedback related to a specific publication record.
func (st *Storage) ReadFeedback(recordId int64) (feedbacks ploc.Feedbacks, err error) {

	query := `SELECT orcid,relevance,presentation,methodology FROM feedback WHERE record_id=?`

	rows, err := st.db.Query(query, recordId)
	if err != nil {
		log.Printf("Database error. Reading a record's feedbacks has failed. %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		var f = ploc.Feedback{RecordId: recordId}

		err = rows.Scan(&f.OrcId, &f.Relevance, &f.Presentation, &f.Methodology)
		if err != nil {
			log.Printf("Database error. Scanning feedback failed. %s", err)
			return
		}

		feedbacks = append(feedbacks, f)
	}

	return
}

// ReadFeedbackFeed returns a list of publications that a user is expert of and that the user may provide feedback to.
// The records are returned as a precomputed JSON data structure for performance reasons.
// The list is descendingly ordered by the number of matching subjects of interest.
// Access to the full list is handled in subsegments via offset and limit, so that a client can read only the segments that are shown to the user.
func (st *Storage) ReadFeedbackFeed(uid int64, offset int64, limit int64) (rawRecords []json.RawMessage, err error) {

	// TODO: Order records by number of received feedbacks.
	const query = `
		 SELECT RTRIM(f.json_preview,'false}'), f.record_id=IFNULL(v.record_id,0)
		 FROM (SELECT r.id AS record_id, r.json_preview AS json_preview 
		 		FROM record AS r, record_feed AS f
		 		WHERE f.user_id=?
		 			AND f.record_id=r.id 
		 			AND f.record_id NOT IN (SELECT record_id FROM feedback WHERE user_id=?) 
		 		ORDER BY f.rowid ASC LIMIT ? OFFSET ?) AS f
		 LEFT JOIN (SELECT record_id FROM record_visit WHERE user_id=?) AS v 
		 ON f.record_id=v.record_id`

	rows, err := st.db.Query(query, uid, uid, limit, offset, uid)
	if err != nil {
		log.Printf("Database error. Querying record feed failed. %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		var rawRecord string
		var visited int64

		err = rows.Scan(&rawRecord, &visited)
		if err != nil {
			log.Printf("Database error. Scanning raw JSON records failed. %s", err)
			return
		}

		var postFix string

		if visited == 0 {
			postFix = "false}"
		} else {
			postFix = "true}"
		}

		rawRecords = append(rawRecords, json.RawMessage(rawRecord+postFix))
	}

	return
}

// ReadOrcId returns a user's associated ORCiD identifiert, which is empty in the case the user has none.
func (st *Storage) ReadOrcId(uid int64) (orcId string, err error) {

	var nullableOrcId sql.NullString

	err = st.db.QueryRow(`SELECT orcid FROM user WHERE id=?`, uid).Scan(&nullableOrcId)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("Database error. Could not read ORCiD of user. User with ID %d seems not to exist.", uid)
	case err != nil:
		log.Printf("Database error. Could not read ORCiD of user with ID %d. %s", uid, err)
	}

	return NullToString(nullableOrcId), err
}

// ReadRecordBookmarks returns all records that were bookmarked by a user.
func (st *Storage) ReadRecordBookmarks(uid int64) (rawBookmarks []json.RawMessage, err error) {

	const query = `
		SELECT b.json, b.collection_ids, b.record_id=IFNULL(v.record_id,0)
		FROM (SELECT r.id AS record_id, r.json_bookmark AS json, IFNULL(GROUP_CONCAT(CAST(b.collection_id AS TEXT)),"") AS collection_ids
				FROM record AS r, record_bookmark AS b
				WHERE b.user_id = ? AND b.record_id = r.id
				GROUP BY r.id
				ORDER BY b.rowid DESC) AS b
		LEFT JOIN (SELECT record_id FROM record_visit WHERE user_id=?) AS v
		ON b.record_id=v.record_id`

	// rows, err := st.db.Query(query, uid)
	rows, err := st.db.Query(query, uid, uid)
	if err != nil {
		log.Printf("Database error. Could not read record bookmarks. %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		var rawBookmark string
		var collectionIds string // the collection IDs are preconcatenated in the SQL query
		var visited int64

		err = rows.Scan(&rawBookmark, &collectionIds, &visited)
		if err != nil {
			log.Printf("Database error. Scanning raw JSON record bookmarks failed. %s", err)
			return
		}

		var inFix string

		if visited == 0 {
			inFix = "false,\"collection_ids\":["
		} else {
			inFix = "true,\"collection_ids\":["
		}

		rightTrimmedRawBookmark := rawBookmark[0 : len(rawBookmark)-26] // 26 == len(`false,"collection_ids":[]}`)

		rawBookmarks = append(rawBookmarks, json.RawMessage(rightTrimmedRawBookmark+inFix+collectionIds+"]}"))
	}

	return
}

// ReadRecordDetails returns detailed information about a specific publication record.
// The detailed information includes information like title, creator names, keywords, DOI and links.
// The details are returned as a precomputed JSON data structure for performance reasons.
func (st *Storage) ReadRecordDetails(uid int64, recordId int64) (rawDetails json.RawMessage, err error) {

	var precomputedDetails string

	// Read precomputed record details

	err = st.db.QueryRow(`SELECT json_detail FROM record WHERE id=?`, recordId).Scan(&precomputedDetails)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("Database error. Could not read details of record. Record with ID %d seems not to exist.", recordId)
	case err != nil:
		log.Printf("Database error. Could not read details of record. %s", err)
	}

	// Mark record as viewed

	_, err = st.db.Exec("INSERT OR IGNORE INTO record_visit (user_id,record_id) VALUES (?,?)", uid, recordId)

	if err != nil {
		log.Printf("Database error. Could not mark record as 'visited'. %s", err)
		return
	}

	return json.RawMessage(precomputedDetails), err
}

// ReadRecordFeedCount returns the total number of records that matches a user's interest.
func (st *Storage) ReadRecordFeedCount(uid int64) (recordCount int64, err error) {

	const query = `
		SELECT count(*) FROM record_feed
		WHERE user_id=? AND record_id NOT IN (SELECT record_id FROM record_bookmark WHERE user_id=?)`

	err = st.db.QueryRow(query, uid, uid).Scan(&recordCount)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("Database error. Could not read number of records matching users interests from database. User %d seems not to exist.", uid)
	case err != nil:
		log.Printf("Database error. Could not read number of records matching users interests from database. %s", err)
	}

	return
}

// ReadRecordFeed returns a list of publications that match the user's subjects of interest.
// The list is descendingly ordered by year of publication and within that year by relevance.
// The list is accessed segment-wise, so that a client can read only that segment that is shown to the user and not the whole list.
// The publications are returned as a precomputed JSON data structure for performance reasons.
func (st *Storage) ReadRecordFeed(uid int64, offset int64, limit int64) (rawRecords []json.RawMessage, err error) {

	// Query needs to exclude bookmarked records (NOT IN), and must determine each record's visited status (LEFT JOIN).
	// TODO: Check performance of left join to determine visited status
	const query = `
		 SELECT RTRIM(f.json_preview,'false}'), f.record_id=IFNULL(v.record_id,0)
		 FROM (SELECT r.id AS record_id, r.json_preview AS json_preview 
		 		FROM record AS r, record_feed AS f
		 		WHERE f.user_id=?
		 			AND f.record_id=r.id 
		 			AND f.record_id NOT IN (SELECT record_id FROM record_bookmark WHERE user_id=?) 
		 		ORDER BY f.rowid ASC LIMIT ? OFFSET ?) AS f
		 LEFT JOIN (SELECT record_id FROM record_visit WHERE user_id=?) AS v 
		 ON f.record_id=v.record_id`

	rows, err := st.db.Query(query, uid, uid, limit, offset, uid)
	if err != nil {
		log.Printf("Database error. Querying record feed failed. %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		var rawRecord string
		var visited int64

		err = rows.Scan(&rawRecord, &visited)
		if err != nil {
			log.Printf("Database error. Scanning raw JSON records failed. %s", err)
			return
		}

		var postFix string

		if visited == 0 {
			postFix = "false}"
		} else {
			postFix = "true}"
		}

		rawRecords = append(rawRecords, json.RawMessage(rawRecord+postFix))
	}

	return
}

// ReadUserInterests returns a list of all the subjects that a user has specified as interesting.
func (st *Storage) ReadUserInterests(uid int64) (subjects ploc.Subjects, err error) {

	query := `
		SELECT s.id, s.keyword 
		FROM subject AS s, interest AS i
		WHERE i.user_id=? AND i.subject_id=s.id`

	rows, err := st.db.Query(query, uid)
	if err != nil {
		log.Printf("Database error. Reading subjects with interest failed. %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		var s ploc.Subject

		err = rows.Scan(&s.Id, &s.Keyword)
		if err != nil {
			log.Printf("Database error. Scanning subjects with interest failed. %s", err)
			return
		}

		subjects = append(subjects, s)
	}

	return
}

// rebuildExpertAndRecordFeed precomputes the record and expert feed for the specified user, based on the subjects
// the user has specified as interesting. The feeds are precomputed because of performance reasons. The function
// needs to be called each time the subjects of interest change or if new records are added to the database.
func (st *Storage) rebuildExpertAndRecordFeed(tx *sql.Tx, uid int64) {

	const createUserExpertFeed = `
		INSERT OR IGNORE INTO expert_feed (user_id,expert_id)
			SELECT DISTINCT ?,esl.expert_id 
			FROM interest AS i, expert_subject_link AS esl 
			WHERE i.user_id=? AND i.subject_id=esl.subject_id 
			GROUP by esl.expert_id
			ORDER BY SUM(esl.record_count) DESC`

	const createUserRecordFeed = `
		INSERT OR IGNORE INTO record_feed (user_id,record_id)
		SELECT DISTINCT ?,rsl.record_id 
		FROM record AS r, interest AS i, record_subject_link AS rsl
		WHERE i.user_id=?
			AND i.subject_id=rsl.subject_id
			AND rsl.record_id NOT IN
				(SELECT record_id FROM record_dislike WHERE user_id=?)
			AND r.id=rsl.record_id
		ORDER BY r.year DESC;`

	tx.Exec("DELETE FROM expert_feed WHERE user_id=?;", uid)
	tx.Exec(createUserExpertFeed, uid, uid)
	tx.Exec("DELETE FROM record_feed WHERE user_id=?;", uid)
	tx.Exec(createUserRecordFeed, uid, uid, uid)

	return
}

// SearchExpertFeed makes a full text search within a user's expert feed and returns a summary for all the
// experts with matching textual content. The searched fields include name, publication titles and subjects.
// The experts are returned as a precomputed JSON data structure for performance reasons.
func (st *Storage) SearchExpertFeed(uid int64, searchTerm string, offset int64, limit int64) (rawExperts []json.RawMessage, err error) {

	const query = `
		SELECT json_preview 
			FROM expert 
			WHERE id IN 
				(SELECT expert_id 
					FROM vexpert 
					WHERE expert_id IN (SELECT expert_id FROM expert_feed WHERE user_id=?) 
						AND expert_id NOT IN (SELECT expert_id FROM expert_bookmark WHERE user_id=?) 
						AND vexpert MATCH '- {expert_id} : ' || ?
						LIMIT ? OFFSET ?)`

	// Quote search term with leading and trailing " to include special characters in FTS.
	// First remove leading and trailing " to avoid duplications.
	// TODO: Remove, as it hides SQLite search query language

	searchTerm = `"` + strings.Trim(searchTerm, `"`) + `"`

	rows, err := st.db.Query(query, uid, uid, searchTerm, limit, offset)
	if err != nil {
		log.Printf("Database error. Searching expert feed failed. %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		var rawExpert string

		err = rows.Scan(&rawExpert)
		if err != nil {
			log.Printf("Database error. Scanning raw JSON expert failed. %s", err)
			return
		}

		rawExperts = append(rawExperts, json.RawMessage(rawExpert))
	}

	return
}

// SearchRecordFeed makes a full text search within a user's publication feed and returns a summary for all the
// publication records with matching textual content. The searched fields include title, abstract, subjects and author names.
// The publications are returned as a precomputed JSON data structure for performance reasons.
func (st *Storage) SearchRecordFeed(uid int64, searchTerm string, offset int64, limit int64) (rawRecords []json.RawMessage, err error) {

	const query = `
		SELECT RTRIM(f.json_preview,'false}'), f.id=IFNULL(v.record_id,0)
		FROM (SELECT id, json_preview FROM record WHERE
			id IN (SELECT record_id FROM vrecord WHERE
				record_id IN (SELECT record_id FROM record_feed WHERE user_id=?) 
				AND record_id NOT IN (SELECT record_id FROM record_bookmark WHERE user_id=?)
			AND vrecord MATCH '- {record_id} : ' || ?
			LIMIT ? OFFSET ?)) AS f 
		LEFT JOIN (SELECT record_id FROM record_visit WHERE user_id=?) AS v 
		ON f.id=v.record_id`

	// Quote search term with leading and trailing " to include special characters in FTS.
	// First remove leading and trailing " to avoid duplications.
	// TODO: Remove, as it hides SQLite search query language

	searchTerm = `"` + strings.Trim(searchTerm, `"`) + `"`

	rows, err := st.db.Query(query, uid, uid, searchTerm, limit, offset, uid)
	if err != nil {
		log.Printf("Database error. Searching record feed failed. %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		var rawRecord string
		var visited int64

		err = rows.Scan(&rawRecord, &visited)
		if err != nil {
			log.Printf("Database error. Scanning raw JSON records failed. %s", err)
			return
		}

		var postFix string

		if visited == 0 {
			postFix = "false}"
		} else {
			postFix = "true}"
		}

		rawRecords = append(rawRecords, json.RawMessage(rawRecord+postFix))
	}

	return
}

// UpdateCollection allows a user to change the title for one of its existing collections.
func (st *Storage) UpdateCollection(uid int64, collectionId int64, title string) (err error) {

	_, err = st.db.Exec("UPDATE collection SET title=? WHERE user_id=? AND id=?", title, uid, collectionId)
	if err != nil {
		log.Printf("Database error. Could not update title of collection. %s", err)
		return
	}

	return
}

// UpdateRecordBookmarkCollections allows a user to specify to which of its bookmark collections a publiction corresponds.
// With help of this function a publication can be added or removed from any of the user's bookmark collections.
func (st *Storage) UpdateCollectionsBookmarkLink(uid int64, recordId int64, collectionIds []int64) (err error) {

	tx, err := st.db.Begin()

	if err != nil {
		log.Printf("Database error. Could not initialize transaction for creating and updating interests. %s", err)
		return
	}

	// 1. Delete all existing relations to collections for the specified record.

	tx.Exec("DELETE FROM record_bookmark WHERE user_id=? AND record_id=? AND (collection_id IS NOT NULL);", uid, recordId)

	// 2. Add all specified collections to each bookmark, but only if user is also the owner of that collection.
	//    Ownership is assured by a sub-select that fails if the 'NOT NULL' condition of 'record_bookmark.user_id' schema is not fullfilled.

	const query = `
		INSERT INTO record_bookmark (user_id,record_id,collection_id)
		VALUES ((SELECT user_id FROM collection WHERE collection.user_id=? AND collection.id=?), ?, ?)`

	for _, collectionId := range collectionIds {
		tx.Exec(query, uid, collectionId, recordId, collectionId)
	}

	err = tx.Commit()

	if err != nil {
		log.Printf("Database error. Could not commit transaction for updating links between record bookmark and collections. %s", err)
		return
	}

	return
}

// UserByGUID returns the major user information like database ID, ORCiD identifier and the hashed secret, based
// on the user's public GUID. This function is primarily used for handling authorization during the routing process.
func (st *Storage) UserByGUID(guid string) (user *model.User, err error) {

	var (
		id           int64
		hashedSecret string
		orcId        sql.NullString
	)

	err = st.db.QueryRow("SELECT id,hashed_secret,orcid FROM user WHERE guid = $1 LIMIT 1;", guid).Scan(
		&id, &hashedSecret, &orcId)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}

	u := model.User{
		Id:           id,
		GUID:         guid,
		HashedSecret: hashedSecret,
		OrcId:        NullToString(orcId),
	}

	return &u, nil
}

package ploc

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// MarshalJSON converts a ReadExpertBookmarksResponse into JSON format, paying respect to fields that already
// contain data in JSON format for performance reasons.
func (r ReadExpertBookmarksResponse) MarshalJSON() ([]byte, error) {

	var buffer bytes.Buffer

	prefix := fmt.Sprintf(`{"bookmarks":[`)
	postfix := `]}`

	buffer.WriteString(prefix)

	delimiter := ``

	for _, bm := range r.RawBookmarks {
		buffer.WriteString(delimiter)
		buffer.Write(bm)
		delimiter = `,`
	}

	buffer.WriteString(postfix)

	return []byte(buffer.String()), nil
}

// MarshalJSON converts a ReadExpertDetailsResponse into JSON format, paying respect to fields that already
// contain data in JSON format for performance reasons.
func (r ReadExpertDetailsResponse) MarshalJSON() ([]byte, error) {
	return []byte(r.RawDetails), nil
}

// MarshalJSON converts a ReadExpertFeedResponse into JSON format, paying respect to fields that already
// contain data in JSON format for performance reasons.
func (r ReadExpertFeedResponse) MarshalJSON() ([]byte, error) {
	return marshalRawExpertFeed(r.RawExperts, r.Offset, r.Limit)
}

// MarshalJSON converts a ReadFeedbackFeedResponse into JSON format, paying respect to fields that already
// contain data in JSON format for performance reasons.
func (r ReadFeedbackFeedResponse) MarshalJSON() ([]byte, error) {
	return marshalRawRecordFeed(r.RawRecords, r.Offset, r.Limit)
}

// MarshalJSON converts a ReadRecordBookmarksResponse into JSON format, paying respect to fields that already
// contain data in JSON format for performance reasons.
func (r ReadRecordBookmarksResponse) MarshalJSON() ([]byte, error) {

	var buffer bytes.Buffer

	prefix := fmt.Sprintf(`{"bookmarks":[`)
	postfix := `]}`

	buffer.WriteString(prefix)

	delimiter := ``

	for _, bm := range r.RawBookmarks {
		buffer.WriteString(delimiter)
		buffer.Write(bm)
		delimiter = `,`
	}

	buffer.WriteString(postfix)

	return []byte(buffer.String()), nil
}

// MarshalJSON converts a ReadRecordDetailsResponse into JSON format, paying respect to fields that already
// contain data in JSON format for performance reasons.
func (r ReadRecordDetailsResponse) MarshalJSON() ([]byte, error) {
	return []byte(r.RawDetails), nil
}

// MarshalJSON converts a ReadRecordFeedResponse into JSON format, paying respect to fields that already
// contain data in JSON format for performance reasons.
func (r ReadRecordFeedResponse) MarshalJSON() ([]byte, error) {
	return marshalRawRecordFeed(r.RawRecords, r.Offset, r.Limit)
}

// MarshalJSON converts a SearchExpertFeedResponse into JSON format, paying respect to fields that already
// contain data in JSON format for performance reasons.
func (r SearchExpertFeedResponse) MarshalJSON() ([]byte, error) {
	return marshalRawExpertFeed(r.RawExperts, r.Offset, r.Limit)
}

// MarshalJSON converts a SearchRecordFeedResponse into JSON format, paying respect to fields that already
// contain data in JSON format for performance reasons.
func (r SearchRecordFeedResponse) MarshalJSON() ([]byte, error) {
	return marshalRawRecordFeed(r.RawRecords, r.Offset, r.Limit)
}

// marshalRawExpertFeed assembles a list of expert in JSON format, using a list of raw experts, offset and limit
// parameters as input.
func marshalRawExpertFeed(rawExperts []json.RawMessage, offset int64, limit int64) ([]byte, error) {

	var buffer bytes.Buffer

	prefix := fmt.Sprintf(`{"offset":%d,"limit":%d,"experts":[`, offset, limit)
	postfix := `]}`

	buffer.WriteString(prefix)

	delimiter := ``

	for _, exp := range rawExperts {
		buffer.WriteString(delimiter)
		buffer.Write(exp)
		delimiter = `,`
	}

	buffer.WriteString(postfix)

	return []byte(buffer.String()), nil
}

// marshalRawRecordFeed assembles a list of records in JSON format, using a list of raw records, offset and limit
// parameters as input.
func marshalRawRecordFeed(rawRecords []json.RawMessage, offset int64, limit int64) ([]byte, error) {

	var buffer bytes.Buffer

	prefix := fmt.Sprintf(`{"offset":%d,"limit":%d,"records":[`, offset, limit)
	postfix := `]}`

	buffer.WriteString(prefix)

	delimiter := ``

	for _, rec := range rawRecords {
		buffer.WriteString(delimiter)
		buffer.Write(rec)
		delimiter = `,`
	}

	buffer.WriteString(postfix)

	return []byte(buffer.String()), nil
}

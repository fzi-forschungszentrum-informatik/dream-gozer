BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS record ( -- a publication record
  id INTEGER PRIMARY KEY, -- unique record ID
  abstract TEXT DEFAULT NULL, -- abstract (optional)
  source_id TEXT NOT NULL UNIQUE,
  doi TEXT DEFAULT NULL, -- document object identification (optional)
  oa INTEGER DEFAULT 0, -- published under open access (-1 = false; 0 = unknown; 1 = true)
  repository_link TEXT DEFAULT NULL, -- link to repository frontpage 
  pdf_link TEXT DEFAULT NULL, -- link to an publicaly available PDF document (optional)
  title TEXT NOT NULL, -- the full title of the publication
  type INTEGER NOT NULL, -- kind of record (e.g. book, conference article, PhD thesis, ...)
  year INTEGER NOT NULL, -- the year of publication
  json_bookmark TEXT NOT NULL, -- precomputed bookmark version of the record in JSON
  json_preview TEXT NOT NULL, -- precomputed preview version of the record in JSON
  json_detail TEXT NOT NULL, -- precomputed detailed version of the record in JSON
  bib_hash TEXT DEFAULT NULL -- bibliographic hash
);

CREATE TABLE IF NOT EXISTS creator ( -- an author's name that is related to a record
  id INTEGER PRIMARY KEY, -- unique creator ID
  expert_id INTEGER DEFAULT NULL, -- related expert profile
  first_name TEXT NOT NULL, -- the first name of the author
  last_name TEXT NOT NULL, -- the last name of the author
  record_id INTEGER NOT NULL -- the record this creator / author relates to
);

CREATE TABLE IF NOT EXISTS expert ( -- an author's name that is related to a record
  id INTEGER PRIMARY KEY, -- unique expert ID
  first_name TEXT NOT NULL, -- the first name of the expert
  last_name TEXT NOT NULL, -- the last name of the expert
  last_publication_year INTEGER DEFAULT NULL, -- the last known year of publication
  total_publication_count INTEGER DEFAULT NULL, -- the number of publications related to that expert
  orcid TEXT DEFAULT NULL, -- optional ORCiD that relates to that expert
  affiliation TEXT DEFAULT NULL, -- optional affiliation of that expert
  json_bookmark TEXT DEFAULT NULL, -- precomputed bookmark version of the expert profile in JSON
  json_preview TEXT DEFAULT NULL, -- precomputed preview version of the expert profile in JSON
  json_detail TEXT DEFAULT NULL -- precomputed detailed version of the expert profile in JSON
);

CREATE TABLE IF NOT EXISTS feedback ( -- open feedback provided by domain experts
  user_id INTEGER NOT NULL, -- the user that has provided the feedback
  record_id INTEGER NOT NULL, -- the record that has received the feedback
  orcid TEXT NOT NULL, -- the expert's ORCID under which the feedback was provided
  relevance INT NOT NULL, -- the rating if the work is relevant (1) or not (0)
  presentation INT NOT NULL, -- the rating if the presentation is convenient (1) or not (0)
  methodology INT NOT NULL, -- the rating if the methodology is conclusive (1) or not (0)
  UNIQUE(user_id,record_id)
);

CREATE TABLE IF NOT EXISTS subject ( -- a subject represents a topic or category of interest
  id INTEGER PRIMARY KEY, -- unique subject ID
  keyword TEXT NOT NULL UNIQUE -- the keyword itself
);

CREATE TABLE IF NOT EXISTS collection ( -- a named list of bookmarks
  id INTEGER PRIMARY KEY, -- unique bookmark collection ID
  user_id INTEGER NOT NULL, -- the owner of the bookmark collection
  title TEXT NOT NULL -- a user-given name
);

CREATE TABLE IF NOT EXISTS user ( -- information related to a user's identity
  id INTEGER PRIMARY KEY, -- we do not use the GUID as primary key for performance and security reasons
  guid TEXT NOT NULL UNIQUE, -- globaly unique identifier for the user (used for access identifikation)
  hashed_secret TEXT NOT NULL, -- the user's secret in its hashed form
  orcid TEXT DEFAULT NULL -- optional ORCiD (identification ID)
);

CREATE TABLE IF NOT EXISTS interest ( -- a number of disjunct subjects define a user's interest
  user_id INTEGER NOT NULL, -- user that has interest
  subject_id INTEGER NOT NULL, -- in a subject
  UNIQUE(user_id,subject_id)
);

CREATE TABLE IF NOT EXISTS record_dislike ( -- publications that are marked as less interesting by a user
  user_id INTEGER NOT NULL, -- the user
  record_id INTEGER NOT NULL, -- the record that is marked as less interesting
  UNIQUE(user_id,record_id)
);

CREATE TABLE IF NOT EXISTS record_bookmark ( -- records that have been bookmarked by a user
  user_id INTEGER NOT NULL, -- the user that has bookmarked a record
  record_id INTEGER NOT NULL, -- the record that is bookmarked
  collection_id INTEGER DEFAULT NULL -- the assiocated collection (optional)
  -- UNIQUE(user_id,record_id,collection_id) -- unfortunately tuples with NULL rows are ignored in SQLite (see index below)
);

CREATE UNIQUE INDEX IF NOT EXISTS record_bookmark_index ON record_bookmark ( -- implement unique constraint even if collection_id is NULL
  user_id,
  record_id,
  IFNULL(collection_id, 0)
);

CREATE TABLE IF NOT EXISTS expert_bookmark ( -- experts that have been bookmarked by a user
  user_id INTEGER NOT NULL, -- the user that has bookmarked an expert
  expert_id INTEGER NOT NULL, -- the expert that is bookmarked
  UNIQUE(user_id,expert_id)
);

CREATE TABLE IF NOT EXISTS record_feed ( -- precomputed record feeds for all the users (performance hack)
  user_id INTEGER NOT NULL, -- reference to the user that owns the feed
  record_id INTEGER NOT NULL, -- reference to the record that shows up in the feed
  UNIQUE(user_id,record_id)
);

CREATE TABLE IF NOT EXISTS expert_feed ( -- precomputed expert feeds for all the users (performance hack)
  user_id INTEGER NOT NULL, -- reference to the user that owns the feed
  expert_id INTEGER NOT NULL, -- reference to the expert that shows up in the feed
  UNIQUE(user_id,expert_id)
);

CREATE TABLE IF NOT EXISTS record_visit ( -- keeps track of the records a user has watched in detail 
  user_id INTEGER NOT NULL, -- reference to the user
  record_id INTEGER NOT NULL, -- reference to the record that has been visited
  UNIQUE(user_id,record_id)
);

-- Junction tables

CREATE TABLE IF NOT EXISTS record_subject_link (
  record_id INTEGER NOT NULL, -- a record
  subject_id INTEGER NOT NULL, -- a subject that is associated with the record
  weight REAL NOT NULL, -- how strong does the subject relate to that topic (0..1)
  UNIQUE(record_id,subject_id)
);

CREATE TABLE IF NOT EXISTS expert_subject_link (
  expert_id INTEGER NOT NULL, -- an expert
  subject_id INTEGER NOT NULL, -- a subject that represents a field of expertise
  record_count INTEGER NOT NULL, -- defines how many record written by that expert are related to that topic (TODO: remove performance hack)
  UNIQUE(expert_id,subject_id)
);

-- Virtual tables for full text search

CREATE VIRTUAL TABLE IF NOT EXISTS vrecord USING FTS5(record_id,title,abstract,subjects,authors,tokenize = 'porter ascii');
CREATE VIRTUAL TABLE IF NOT EXISTS vexpert USING FTS5(expert_id,full_name,titles,subjects,tokenize = 'porter ascii');

COMMIT TRANSACTION;

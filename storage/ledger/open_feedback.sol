pragma solidity ^0.5.11;

// A smart contract that allows to publish and read feedback to scientific publications transparently.
// Contributors that provided feedback are identified by their ORCiDs.
// Publications that are reviewed are identified by their unique bibliographic hashes (based on normalized title, author names and year of publication).
// Feedback given can be requested by anyone either by concrete publication (via bibliographic hash) or by a person (via its ORCiD).

contract open_feedback{

    struct Feedback{
        address serviceAddress; // The service address that has delivered a feedback. 
        bytes16 orcId; // Globaly unique identifiaction number for researchers.
        bytes16 bibHash; // Globaly unique identification number for publications.
        uint timestamp; // Date that defines when the feedback was delivered.
        uint8 relevance; // Is this work relevant in its field?
        uint8 presentation; // Does the presentation make a good impression? 
        uint8 methodology; // Does the methodology make a sound impression?
    }
    
    // Contract owner
    address owner;    
    
    mapping(bytes16 => Feedback[]) feedbackByOrcId;
    
    mapping(bytes16 => Feedback[]) feedbackByBibHash;
    
    constructor() public {owner = msg.sender;}

    function addFeedback(bytes16 _orcId, bytes16 _bibHash, uint8 _relevance, uint8 _presentation, uint8 _methodology) public {

        Feedback memory newFeedback;
        newFeedback.serviceAddress = msg.sender;
        newFeedback.orcId = _orcId;
        newFeedback.bibHash = _bibHash;
        newFeedback.timestamp = now;
        newFeedback.relevance = _relevance;
        newFeedback.presentation = _presentation;
        newFeedback.methodology = _methodology;
        
        feedbackByBibHash[_bibHash].push(newFeedback);
        feedbackByOrcId[_orcId].push(newFeedback);
    }
    
    // Total number of feedbacks for a publication. The publication is represented by its bibliographic hash.
    function getFeedbackCountByBibHash(bytes16 _bibHash) public view returns (uint count_) {
        count_ = feedbackByBibHash[_bibHash].length;
    }

    // Returns the n-th feedback for a publication. The publication is represented by its biliographic hash. The index specifies which of all the feedbacks for that publication is returned.
    function getFeedbackByBibHash(bytes16 _bibhash, uint _index) public view returns (address serviceAddress_, bytes16 orcId_, uint timestamp_, uint8 relevance_, uint8 presentation_, uint8 methodology_) {
        
        require(
            _index < feedbackByBibHash[_bibhash].length,
            "Specified index is out of range"
        );

        Feedback storage f = feedbackByBibHash[_bibhash][_index];
        
        serviceAddress_ = f.serviceAddress;
        orcId_ = f.orcId;
        timestamp_ = f.timestamp;
        relevance_ = f.relevance;
        presentation_ = f.presentation;
        methodology_ = f.methodology;
    }

    // Total number of feedbacks that an expert has provided. The expert is represented by its OrcId.
    function getFeedbackCountByOrcId(bytes16 _orcid) public view returns (uint count_) {
        count_ = feedbackByOrcId[_orcid].length;
    }

    // Returns the n-th feedback provided by an expert. The expert is represented by its OrcId. The index specifies which of all the feedbacks provided by that expert is returned.
    function getFeedbackByOrcId(bytes16 _orcid, uint _index) public view returns (address serviceAddress_, bytes16 bibHash_, uint timestamp_, uint8 relevance_, uint8 presentation_, uint8 methodology_) {

        require(
            _index < feedbackByOrcId[_orcid].length,
            "Specified index is out of range"
        );

        Feedback storage f = feedbackByOrcId[_orcid][_index];
        
        serviceAddress_ = f.serviceAddress;
        bibHash_ = f.bibHash;
        timestamp_ = f.timestamp;
        relevance_ = f.relevance;
        presentation_ = f.presentation;
        methodology_ = f.methodology;
    }

    // Returns the total feedback of an publication. The publication is referenced by its bibliographic hash.
    function getTotalFeedbackByBibHash(bytes16 _bibhash) public view returns (uint feedbackCount_, uint relevanceTotal_, uint presentationTotal_, uint methodologyTotal_) {
        feedbackCount_ = feedbackByBibHash[_bibhash].length;
        for (uint i = 0; i < feedbackCount_; i++) {
          if (feedbackByBibHash[_bibhash][i].relevance > 0) {
              relevanceTotal_++;
          }
          if (feedbackByBibHash[_bibhash][i].presentation > 0) {
              presentationTotal_++;
          }
          if (feedbackByBibHash[_bibhash][i].methodology > 0) {
              methodologyTotal_++;
          }
       }
    }
}  

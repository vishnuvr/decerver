// GithubHandler is an example of an abstract handler that is used to service
// webhook-postings (like those provided by the Github API). It uses a mapping
// of "X-github-event" types and functions to handle various types of posts.
// Each of those handler functions will create an object from the JSON
// data inside the post body.
//
// The event handler functions can be added and removed from the map during
// runtime. The GithubHandler itself has a mutex in order to prevent people
// from mucking around with the map while it's being used.
package ghhandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

// Function that handles a specific type of posting
type PostHandlerFunc func(http.ResponseWriter, *http.Request)

type GithubHandler struct {
	mutex    *sync.Mutex
	mappings map[string]PostHandlerFunc
}

// Create a new handler
func NewHandler() *GithubHandler {
	gh := &GithubHandler{}
	gh.mutex = &sync.Mutex{}
	gh.mappings = make(map[string]PostHandlerFunc)
	gh.mappings["issues"] = handleIssues
	gh.mappings["issue_comment"] = handleIssueComment
	return gh
}

// Function that handles incoming github webhook-postings. It delegates the actual
// handling based on the "X-github-event" type of the request.
func (gh *GithubHandler) Handle(res http.ResponseWriter, req *http.Request) {
	evt := req.Header.Get("X-github-event")
	if evt == "" {
		fmt.Errorf("Request is not a github event: %s\n", evt)
		return
	}

	gh.mutex.Lock()
	if gh.mappings[evt] == nil {
		fmt.Errorf("Request not supported: %s\n", evt)
	} else {
		gh.mappings[evt](res, req)
	}

	gh.mutex.Unlock()
}

// Add a new event handler
func (gh *GithubHandler) AddPostHandler(eventType string, postHandler PostHandlerFunc, replaceOld bool) error {
	gh.mutex.Lock()
	if gh.mappings[eventType] != nil {
		if !replaceOld {
			gh.mutex.Unlock()
			return errors.New("Tried to overwrite an already existing function mapping.")
		} else {
			fmt.Println("Overwriting old handler for '" + eventType + "'.")
		}

	}
	gh.mappings[eventType] = postHandler
	gh.mutex.Unlock()
	return nil
}

// Remove an event handler
func (gh *GithubHandler) RemovePostHandler(eventType string) {
	gh.mutex.Lock()
	if gh.mappings[eventType] == nil {
		fmt.Println("Removal failed. There is no handler for '" + eventType + "'.")
	} else {
		delete(gh.mappings, eventType)
	}
	gh.mutex.Unlock()
	return
}

// Create an object from the body of a "issues" post.
func handleIssues(res http.ResponseWriter, req *http.Request) {
	post := &IssuePost{}
	json.NewDecoder(req.Body).Decode(post)
	fmt.Printf("%+v\n", post.Action)
	fmt.Printf("%+v\n", post.Issue)
	fmt.Printf("%+v\n", post.Assignee)
	fmt.Printf("%+v\n", post.Repository)
	fmt.Printf("%+v\n", post.Sender)
}

// Create an object from the body of a "issue_comment" post.
func handleIssueComment(res http.ResponseWriter, req *http.Request) {
	post := &IssueCommentPost{}
	json.NewDecoder(req.Body).Decode(post)
	fmt.Printf("%+v\n", post.Action)
	fmt.Printf("%+v\n", post.Issue)
	fmt.Printf("%+v\n", post.Comment)
	fmt.Printf("%+v\n", post.Repository)
	fmt.Printf("%+v\n", post.Sender)
}

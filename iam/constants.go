package iam

import "encoding/json"

type IAM struct {
	Clipsynphony Clipsynphony `json:"clipsynphony"`
	Newapp []string `json:"newapp"`
	Age int `json:"age"`
}
type Clipsynphony struct {
	Root []string `json:"root"`
	Admin []string `json:"admin"`
	Commenter []string `json:"commenter"`
	Editor []string `json:"editor"`
	Publisher []string `json:"publisher"`
	Reader []string `json:"reader"`
}

const data = `{
  "clipsynphony" : {
    "root": [
      "CLIPSYNPHONY_CREATE_TEAMS",
      "CLIPSYNPHONY_CREATE_ORGANISATIONS",
      "CLIPSYNPHONY_DELETE_TEAMS",
      "CLIPSYNPHONY_DELETE_ORGANISATIONS"
    ],
    "admin": [
      "CLIPSYNPHONY_UPDATE_TEAMS",
      "CLIPSYNPHONY_ASSIGN_ROLES",
      "CLIPSYNPHONY_INVITE_USERS",
      "CLIPSYNPHONY_SENDMAIL"
    ],
    "commenter": [
      "CLIPSYNPHONY_COMMENT_ARTICLE",
      "CLIPSYNPHONY_REPORT_ARTICLE"
    ],
    "editor": [
      "CLIPSYNPHONY_WRITE_ARTICLE",
      "CLIPSYNPHONY_EDIT_ARTICLE",
      "CLIPSYNPHONY_CREATE_ARTICLE",
      "CLIPSYNPHONY_APPROVE_ARTICLE"
    ],
    "publisher": [
      "CLIPSYNPHONY_PUBLISH_ARTICLE",
      "CLIPSYNPHONY_REDACT_ARTICLE"
    ],
    "reader": [
      "CLIPSYNPHONY_READ_ARTICLE",
      "CLIPSYNPHONY_REVIEW_ARTICLE"
    ]
  },
  "newapp": ["just something"],
  "age": 17
}`

func AllPermissions() IAM {
	var permissions IAM
	_ = json.Unmarshal([]byte(data), &permissions)

	return permissions
}

func AllPermissionsMap() map[string]*json.RawMessage {
	var Objmap map[string]*json.RawMessage
	_ = json.Unmarshal([]byte(data), &Objmap)

	return Objmap
}

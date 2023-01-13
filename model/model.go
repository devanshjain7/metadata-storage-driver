// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type MetaBackend struct {
	Id          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	BackendName string        `json:"backendName,omitempty" bson:"backendName,omitempty"`
	Type        string        `json:"type,omitempty" bson:"type,omitempty"`
	Region      string        `json:"region" bson:"region"`
	Buckets     []MetaBucket  `json:"buckets,omitempty" bson:"buckets,omitempty"`
}

type MetaBucket struct {
	CreationDate    *time.Time   			`type:"timestamp"`
	Name            string       			`type:"string"`
	BucketType      string       			`json:"bucketType,omitempty" bson:"bucketType,omitempty"`
	Region          string       			`json:"region,omitempty" bson:"region,omitempty"`
	Access          string       			`json:"access,omitempty" bson:"access,omitempty"`
	NumberOfObjects int64        			`json:"numberOfObjects,omitempty" bson:"numberOfObjects,omitempty"`
	TotalSize       int64        			`json:"totalSize,omitempty" bson:"totalSize,omitempty"`
	BucketTags      map[string]string       `json:"tags,omitempty" bson:"tags,omitempty"`
	Objects         []MetaObject 			`json:"objects,omitempty" bson:"objects"`
}

type MetaObject struct {
	ObjectName                  string     `json:"objectName" bson:"objectName"`
	LastModifiedDate            *time.Time `type:"timestamp" json:"lastModifiedDate" bson:"lastModifiedDate"`
	Size                        int64      `json:"size" bson:"size"`
	ServerSideEncryption        string     `json:"serverSideEncryption" bson:"serverSideEncryption"`
	VersionId                   string     `json:"versionId,omitempty" bson:"versionId,omitempty"`
	StorageClass                string     `json:"storageClass,omitempty" bson:"storageClass,omitempty"`
	RedirectLocation            string     `json:"redirectLocation,omitempty" bson:"redirectLocation,omitempty"`
	ReplicationStatus           string     `json:"replicationStatus,omitempty" bson:"replicationStatus,omitempty"`
	ExpiresDate                 string     `json:"expiresDate,omitempty" bson:"expiresDate,omitempty"`
	GrantControl                string     `json:"grantControl,omitempty" bson:"grantControl,omitempty"`
	ObjectTags                  string     `json:"objectTags,omitempty" bson:"objectTags,omitempty"`
	Metadata                    string     `json:"metadata,omitempty" bson:"metadata,omitempty"`
	ObjectType                  string     `json:"objectType,omitempty" bson:"objectType,omitempty"`
}

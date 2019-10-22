/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TweetsSpec defines the desired state of Tweets
type TweetsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Message string `json:"message"`
}

// TweetsStatus defines the observed state of Tweets
type TweetsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	TweetId int64 `json:"tweetid"`
}

// +kubebuilder:object:root=true

// Tweets is the Schema for the tweets API
type Tweets struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec          TweetsSpec   `json:"spec,omitempty"`
	Status        TweetsStatus `json:"status,omitempty"`
	IsAutoCreated int          `json:"isautocreated"`
}

// +kubebuilder:object:root=true

// TweetsList contains a list of Tweets
type TweetsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Tweets `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Tweets{}, &TweetsList{})
}

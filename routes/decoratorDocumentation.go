//   Copyright 2017 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package routes

// Documentation decorates a handler to document it. Summary, Description and
// Tags will be set on the documentation if not zero.
type Documentation HandlerDocumentationBody

func (d Documentation) Decorate(h Handler) Handler {
	newDocs := h.Documentation
	if d.Summary != "" {
		newDocs.Summary = d.Summary
	}
	if d.Description != "" {
		newDocs.Description = d.Description
	}
	for k, v := range d.Tags {
		if newDocs.Tags == nil {
			newDocs.Tags = make(Tags)
		}
		t := newDocs.Tags[k]
		newDocs.Tags[k] = append(t, v...)
	}
	h.Documentation = newDocs
	return h
}

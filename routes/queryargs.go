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

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	maxUint = ^uint(0)
	maxInt  = int(maxUint >> 1)
	minInt  = -maxInt - 1
)

type (
	// QueryArgInt denotes an int query argument. It fulfills the QueryParser
	// interface.
	QueryArgInt struct{}

	// QueryArgUint denotes an uint query argument. It fulfills the QueryParser
	// interface.
	QueryArgUint struct{}

	// QueryArgString denotes a string query argument. It fulfills the QueryParser
	// interface.
	QueryArgString struct{}

	// QueryArgIntSlice denotes a []int query argument. It fulfills the QueryParser
	// interface.
	QueryArgIntSlice struct{}

	// QueryArgUintSlice denotes a []uint query argument. It fulfills the QueryParser
	// interface.
	QueryArgUintSlice struct{}

	// QueryParser parses a string and returns a typed value.
	// An error can be returned if the value could not be parse. The error's
	// message contains %s which has to be replaced by the argument's name
	// before being displayed.
	QueryParser interface {
		QueryParse(string) (interface{}, error)
	}

	// QueryArg defines an argument by its name and its type.
	QueryArg struct {
		Name string
		Type QueryParser
	}

	// WithQueryArg contains all the arguments to parse in the URL.
	// WithQueryArg has a method Decorate called to apply the decorators
	// on an endpoint.
	WithQueryArg []QueryArg
)

// QueryParse parses an int. A nil error indicates a success. With this func,
// QueryArgInt fulfills QueryArgType.
func (d QueryArgInt) QueryParse(val string) (interface{}, error) {
	if i, err := strconv.ParseInt(val, 10, 64); err == nil &&
		i <= int64(maxInt) && i >= int64(minInt) {
		return int(i), nil
	}
	return nil, errors.New("argument \"%s\" must be an int")
}

// QueryParse parses an uint. A nil error indicates a success. With this func,
// QueryArgInt fulfills QueryArgType.
func (d QueryArgUint) QueryParse(val string) (interface{}, error) {
	if i, err := strconv.ParseUint(val, 10, 64); err == nil &&
		i <= uint64(maxUint) {
		return uint(i), nil
	}
	return nil, errors.New("argument \"%s\" must be an uint")
}

// QueryParse parses a string. A nil error indicates a success. With this func,
// QueryArgInt fulfills QueryArgType.
func (d QueryArgString) QueryParse(val string) (interface{}, error) {
	return val, nil
}

// QueryParse parses an []int. A nil error indicates a success. With this func,
// QueryArgInt fulfills QueryArgType.
func (d QueryArgIntSlice) QueryParse(val string) (interface{}, error) {
	vals := strings.Split(val, ",")
	res := make([]int, 0, len(vals))
	for _, v := range vals {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil &&
			i <= int64(maxInt) && i >= int64(minInt) {
			res = append(res, int(i))
		} else {
			return nil, errors.New("argument \"%s\" must be a slice of int")
		}
	}
	return res, nil
}

// QueryParse parses an []uint. A nil error indicates a success. With this func,
// QueryArgInt fulfills QueryArgType.
func (d QueryArgUintSlice) QueryParse(val string) (interface{}, error) {
	vals := strings.Split(val, ",")
	res := make([]uint, 0, len(vals))
	for _, v := range vals {
		if i, err := strconv.ParseUint(v, 10, 64); err == nil &&
			i <= uint64(maxUint) {
			res = append(res, uint(i))
		} else {
			return nil, errors.New("argument \"%s\" must be a slice of uint")
		}
	}
	return res, nil
}

func parseArg(arg QueryArg, r *http.Request, a Arguments) (int, ErrorBody) {
	if rawVal := r.URL.Query().Get(arg.Name); rawVal != "" {
		if val, err := arg.Type.QueryParse(rawVal); err == nil {
			a[arg] = val
		} else {
			msg := fmt.Sprintf(err.Error(), arg.Name)
			return 400, ErrorBody{msg}
		}
	} else {
		msg := fmt.Sprintf("argument \"%s\" not found", arg.Name)
		return 400, ErrorBody{msg}
	}
	return 200, ErrorBody{}
}

// Decorate is the function called to apply the decorators to an endpoint. It returns
// a function. This function produces a 400 error code with a json error message or
// calls the next IntermediateHandler.
// The goal of this function is to get the URL parameters to store them in
// the Arguments.
func (d WithQueryArg) Decorate(h IntermediateHandler) IntermediateHandler {
	return func(w http.ResponseWriter, r *http.Request, a Arguments) (status int, output interface{}) {
		for _, arg := range d {
			if code, err := parseArg(arg, r, a); code != 200 {
				return code, err
			}
		}
		return h(w, r, a)
	}
}

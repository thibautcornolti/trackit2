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

package aws

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/models"
	"github.com/trackit/trackit2/routes"
	"github.com/trackit/trackit2/users"
)

var (
	errFailDeleteAccount = errors.New("Failed to delete account.")
)

// deleteAwsAccount is a route handler which lets the user delete AwsAccounts
// from their account.
func deleteAwsAccount(r *http.Request, a routes.Arguments) (int, interface{}) {
	ctx := r.Context()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	tx := a[db.Transaction].(*sql.Tx)
	u := a[users.AuthenticatedUser].(users.User)
	id := a[QueryArgAwsAccount].(uint)
	dbAwsAccount, err := models.AwsAccountByID(tx, int(id))
	if err != nil {
		logger.Error("Failed to get AWS Account.", err)
		return 500, errFailDeleteAccount
	} else if dbAwsAccount.UserID != u.Id {
		logger.Error("AWS Account does not belong to the user.", err)
		return 500, errFailDeleteAccount
	}
	if err := dbAwsAccount.Delete(tx); err != nil {
		logger.Error("Failed to delete AWS Account.", err)
		return 500, errFailDeleteAccount
	}
	return 200, fmt.Sprintf("Account %d deleted.", id)
}

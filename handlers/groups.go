// Statup
// Copyright (C) 2018.  Hunter Long and the project contributors
// Written by Hunter Long <info@socialeck.com> and the project contributors
//
// https://github.com/hunterlong/statup
//
// The licenses for most software and other practical works are designed
// to take away your freedom to share and change the works.  By contrast,
// the GNU General Public License is intended to guarantee your freedom to
// share and change all versions of a program--to make sure it remains free
// software for all its users.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/hunterlong/statping/core"
	"github.com/hunterlong/statping/utils"
	"net/http"
)

// apiAllGroupHandler will show all the groups
func apiAllGroupHandler(w http.ResponseWriter, r *http.Request) {
	if !IsReadAuthenticated(r) {
		sendUnauthorizedJson(w, r)
		return
	}
	auth := IsUser(r)
	groups := core.SelectGroups(false, auth)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

// apiGroupHandler will show a single group
func apiGroupHandler(w http.ResponseWriter, r *http.Request) {
	if !IsReadAuthenticated(r) {
		sendUnauthorizedJson(w, r)
		return
	}
	vars := mux.Vars(r)
	group := core.SelectGroup(utils.ToInt(vars["id"]))
	if group == nil {
		sendErrorJson(errors.New("group not found"), w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(group)
}

// apiCreateGroupHandler accepts a POST method to create new groups
func apiCreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	if !IsFullAuthenticated(r) {
		sendUnauthorizedJson(w, r)
		return
	}
	var group *core.Group
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&group)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	_, err = group.Create()
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	sendJsonAction(group, "create", w, r)
}

// apiGroupDeleteHandler accepts a DELETE method to delete groups
func apiGroupDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if !IsFullAuthenticated(r) {
		sendUnauthorizedJson(w, r)
		return
	}
	vars := mux.Vars(r)
	group := core.SelectGroup(utils.ToInt(vars["id"]))
	if group == nil {
		sendErrorJson(errors.New("group not found"), w, r)
		return
	}
	err := group.Delete()
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	sendJsonAction(group, "delete", w, r)
}

type groupOrder struct {
	Id    int64 `json:"group"`
	Order int   `json:"order"`
}

func apiGroupReorderHandler(w http.ResponseWriter, r *http.Request) {
	if !IsFullAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	var newOrder []*groupOrder
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&newOrder)
	for _, g := range newOrder {
		group := core.SelectGroup(g.Id)
		group.Order = g.Order
		group.Update()
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newOrder)
}

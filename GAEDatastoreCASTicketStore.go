package dsts

import (
	"github.com/OlympiaSchoolDistrict/cas"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

type DSAuthenticationResponse struct {
	User                string    // Users login name
	ProxyGrantingTicket string    `datastore:",noindex"` // Proxy Granting Ticket
	Proxies             []string  `datastore:",noindex"` // List of proxies
	AuthenticationDate  time.Time `datastore:",noindex"` // Time at which authentication was performed
	IsNewLogin          bool      `datastore:",noindex"` // Whether new authentication was used to grant the service ticket
	IsRememberedLogin   bool      `datastore:",noindex"` // Whether a long term token was used to grant the service ticket
	MemberOf            []string  // List of groups which the user is a member of
	// cas.AuthenticationResponse
	RecvdDate time.Time //Time this response was recieved
}

var (
	ticketStoreGAEDSType = "TicketStore2"
	ticketStoreDefaultID = "defaultTicketStore"
	authResponseKind     = "AuthResponse2"
	authResponseAttrKind = "AuthResponseAttr2"
)

type DatastoreTicketStore struct {
	TicketStoreID string
	ctx           context.Context
}

func NewDataTicketStore(ctx context.Context) *DatastoreTicketStore {
	return &DatastoreTicketStore{ctx: ctx}
}

func (s *DatastoreTicketStore) key() *datastore.Key {
	if s.TicketStoreID == "" {
		s.TicketStoreID = ticketStoreDefaultID
	}
	return datastore.NewKey(s.ctx, ticketStoreGAEDSType, s.TicketStoreID, 0, nil)
}

func (s *DatastoreTicketStore) Read(id string) (*cas.AuthenticationResponse, error) {
	// var resp cas.AuthenticationResponse
	var rresp DSAuthenticationResponse
	k := datastore.NewKey(s.ctx, authResponseKind, id, 0, s.key())
	err := datastore.Get(s.ctx, k, &rresp)
	if err != nil {
		return nil, err
	}
	resp := cas.AuthenticationResponse{
		User:                rresp.User,
		ProxyGrantingTicket: rresp.ProxyGrantingTicket,
		Proxies:             rresp.Proxies,
		AuthenticationDate:  rresp.AuthenticationDate,
		IsNewLogin:          rresp.IsNewLogin,
		IsRememberedLogin:   rresp.IsRememberedLogin,
		MemberOf:            rresp.MemberOf,
	}

	var pl datastore.PropertyList
	pk := datastore.NewKey(s.ctx, authResponseAttrKind, "Attributes", 0, k)
	err = datastore.Get(s.ctx, pk, &pl)

	ua := map[string][]string{}
	for _, p := range pl {
		ua[p.Name] = append(ua[p.Name], p.Value.(string))
	}
	resp.Attributes = ua
	return &resp, nil
}

func (s *DatastoreTicketStore) Write(id string, t *cas.AuthenticationResponse) error {
	pl := datastore.PropertyList{}
	for n, v := range t.Attributes {
		for _, vs := range v {
			pl = append(pl, datastore.Property{Name: n, Value: vs, Multiple: true})
		}
	}

	ticket := DSAuthenticationResponse{
		User:                t.User,
		ProxyGrantingTicket: t.ProxyGrantingTicket,
		Proxies:             t.Proxies,
		AuthenticationDate:  t.AuthenticationDate,
		IsNewLogin:          t.IsNewLogin,
		IsRememberedLogin:   t.IsRememberedLogin,
		MemberOf:            t.MemberOf,
		RecvdDate:           time.Now(),
	}

	k := datastore.NewKey(s.ctx, authResponseKind, id, 0, s.key())
	_, err := datastore.Put(s.ctx, k, &ticket)

	if err != nil {
		return err
	}

	a := datastore.NewKey(s.ctx, authResponseAttrKind, "Attributes", 0, k)
	_, err = datastore.Put(s.ctx, a, &pl)
	return err
}

func (s *DatastoreTicketStore) Delete(id string) error {
	k := datastore.NewKey(s.ctx, authResponseKind, id, 0, s.key())
	a := datastore.NewKey(s.ctx, authResponseAttrKind, "Attributes", 0, k)

	err := datastore.Delete(s.ctx, a)
	if err != nil {
		return err
	}

	err = datastore.Delete(s.ctx, k)
	return err

}

// Clear removes all ticket data
func (s *DatastoreTicketStore) Clear() error {
	ks, err := datastore.NewQuery("").Ancestor(s.key()).KeysOnly().GetAll(s.ctx, nil)
	if err != nil {
		return err
	}
	err = datastore.DeleteMulti(s.ctx, ks)
	return err
}

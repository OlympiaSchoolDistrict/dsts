package dsts

import (
	"github.com/OlympiaSchoolDistrict/cas"
	"google.golang.org/appengine/aetest"
	"reflect"
	"testing"
	"time"
)

var (
	Tdsts DatastoreTicketStore
)

func TestDSTS(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	Tdsts := DatastoreTicketStore{ticketStoreID: "testingID", ctx: ctx}
	tim, _ := time.Parse(time.Kitchen, "5:20AM")

	ticket := &cas.AuthenticationResponse{User: "exampleUser", AuthenticationDate: tim.Local(), Attributes: cas.UserAttributes{
		"singlething": []string{"justone"},
		"multithing":  []string{"one", "two", "three", "fifteen"},
	},
	}
	id := "reallyLongUnlikelyTicketNameThatWouldntLikelyOccurInRealLife"

	ti, err := Tdsts.Read(id)
	if err == nil {
		t.Errorf("Error: %v; want enitity not found error", err)
	}
	if ti != nil {
		t.Errorf("Error: %+v; bad read should return nil", ti)
	}

	err = Tdsts.Write(id, ticket)
	if err != nil {
		t.Fatalf("Write result: error: %+v", err)
	}
	// fmt.Fprintf(w, "<p>Write result: error: %+v</p>", err)

	ti, err = Tdsts.Read(id)
	if err != nil || !reflect.DeepEqual(ticket, ti) {
		t.Fatalf("Read result doesn't equal Write: error: %+v, Original: %+v, Read: %+v", err, ticket, ti)
	}

	Tdsts.Delete(id)
	if err != nil {
		t.Fatalf("Delete: error: %+v", err)
	}

	ti, err = Tdsts.Read(id)
	if err == nil {
		t.Errorf("Error: %v; want enitity not found error", err)
	}
	if ti != nil {
		t.Errorf("Error: %+v; bad read should return nil", ti)
	}

}

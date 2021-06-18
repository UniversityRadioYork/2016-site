package models

import (
	"testing"

	"github.com/UniversityRadioYork/myradio-go"
)

func TestIndexGetInvalidJSON(t *testing.T) {
	message := []byte("'invalid json'")

	// Error explicitly ignored as MockSession can never return an error
	session, _ := myradio.MockSession(message)
	model := NewIndexModel(session)
	_, _, _, _, _, err := model.Get()
	if err == nil {
		t.Errorf("Expected Error, Recieved %v", err)
	}
	t.Errorf("Expected Error, Recieved %v", err)
}

func TestIndexGetInvalidMessage(t *testing.T) {
	message := []byte("{\"message\": \"Invalid\"}")

	// Error explicitly ignored as MockSession can never return an error
	session, _ := myradio.MockSession(message)
	model := NewIndexModel(session)
	_, _, _, _, _, err := model.Get()
	if err == nil {
		t.Errorf("Expected Error, Recieved %v", err)
	}
	t.Errorf("Expected Error, Recieved %v", err)
}

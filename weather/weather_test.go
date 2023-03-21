package weather 


import (
	
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	
	
)


func TestGetLocationKey(t *testing.T) {
	testCases := []struct {
		city   string
		state  string
		status int
		key    string
		errMsg string
	}{
		{city: "New York", state: "NY", status: http.StatusOK, key: "349727", errMsg: ""},
		{city: "San Francisco", state: "CA", status: http.StatusOK, key: "347629", errMsg: ""},
		{city: "", state: "", status: http.StatusOK, key: "", errMsg: "no city/state found for , "},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("city=%s,state=%s", tc.city, tc.state), func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/locations/v1/cities/search" {
					t.Errorf("unexpected path %s", r.URL.Path)
				}
				if err := r.ParseForm(); err != nil {
					t.Errorf("failed to parse form: %v", err)
				}
				if apiKey := r.Form.Get("apikey"); apiKey != "9p3ELusMaMlbU4fJF79K97c4iQjp7Zq6" {
					t.Errorf("unexpected apikey %s", apiKey)
				}
				if q := r.Form.Get("q"); q != fmt.Sprintf("%s, %s", tc.city, tc.state) {
					t.Errorf("unexpected query %s", q)
				}

				w.Header().Set("Content-Type", "application/json")
				if tc.status == http.StatusOK {
					fmt.Fprintf(w, `[{"Key": "%s"}]`, tc.key)
				} else {
					http.Error(w, "not found", tc.status)
				}
			}

			server := httptest.NewServer(http.HandlerFunc(handler))
			defer server.Close()

			key, err := GetLocationKey( tc.city, tc.state)
			if err != nil {
				if tc.errMsg == "" {
					t.Errorf("unexpected error: %v", err)
				} else if errMsg := err.Error(); errMsg != tc.errMsg {
					t.Errorf("unexpected error message %s, expected %s", errMsg, tc.errMsg)
				}
			} else if key != tc.key {
				t.Errorf("unexpected location key %s, expected %s", key, tc.key)
			}
		})
	}
}

func TestGetCurrentConditions(t *testing.T) {
    // Create a test server
    server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        if req.URL.Path == "/currentconditions/v1/locationkey" {
            // Return a sample response for the AccuWeather API
            rw.Write([]byte(`[
                {
                    "WeatherText": "Sunny",
                    "Temperature": {
                        "Imperial": {
                            "Value": 72.0,
                            "Unit": "F"
                        }
                    }
                }
            ]`))
        } else {
            rw.WriteHeader(http.StatusNotFound)
        }
    }))
    defer server.Close()

    // Call the GetCurrentConditions function with a mocked location key
    cc, err := GetCurrentConditions(server.URL+"/currentconditions/v1/locationkey")

    // Check that there was no error
    if err != nil {
        t.Fatalf("expected no error, but got %v", err)
    }

    // Check that the current conditions were parsed correctly
    if cc.WeatherText != "Sunny" {
        t.Errorf("expected WeatherText to be 'Sunny', but got '%s'", cc.WeatherText)
    }

    if cc.Temperature.Imperial.Value != 72.0 {
        t.Errorf("expected Temperature.Imperial.Value to be 72.0, but got %f", cc.Temperature.Imperial.Value)
    }
}

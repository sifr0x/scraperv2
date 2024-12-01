package massapi

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/datahuys/scraperv2/domain"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestApiWozRepository_Store(t *testing.T) {
	// Start mocking http calls
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("Successful_Store", func(t *testing.T) {
		// Mock the response for the API endpoint
		httpmock.RegisterResponder("POST", "http://example.com/api",
			func(req *http.Request) (*http.Response, error) {
				// You can customize the response as needed for testing
				resp := httpmock.NewStringResponse(http.StatusOK, "success")
				return resp, nil
			},
		)

		// Create an instance of the repository
		repo := NewApiWozRepository("http://example.com/api", http.DefaultClient)

		// Prepare a dummy Woz object
		dummyWoz := domain.Woz{
			ID:        123,
			Status:    200, // Assuming this is a valid status
			ScrapedAt: time.Now(),
			Payload:   []byte(`{"wozObject":{"wozobjectnummer":11800039614,"woonplaatsnaam":"Hoogeveen","openbareruimtenaam":"Anjerstraat","straatnaam":"Anjerstraat","postcode":"7906LE","huisnummer":5,"huisletter":null,"huisnummertoevoeging":null,"locatieomschrijving":null,"gemeentecode":118,"grondoppervlakte":122,"adresseerbaarobjectid":118010900087577,"nummeraanduidingid":null,"verbondenAdresseerbareObjecten":[118010900087577],"ontleendeAdresseerbareObjecten":[118010900087577]},"wozWaarden":[],"panden":[{"bagpandidentificatie":118100001200409}],"kadastraleObjecten":[{"kadastraleGemeenteCode":"HGV00","kadastraleSectie":"I","kadastraalPerceelNummer":"4945"}]}`),
		}

		// Call the Store method
		err := repo.Store(dummyWoz)

		// Assert that no error occurred
		assert.NoError(t, err)
	})

	t.Run("Unsuccessful_Store", func(t *testing.T) {
		// Mock the response for the API endpoint
		httpmock.RegisterResponder("POST", "http://example.com/api",
			func(req *http.Request) (*http.Response, error) {
				// Return a non-OK status code
				resp := httpmock.NewStringResponse(http.StatusNotFound, "not found")
				return resp, nil
			},
		)

		// Create an instance of the repository
		repo := NewApiWozRepository("http://example.com/api", http.DefaultClient)

		// Prepare a dummy Woz object
		dummyWoz := domain.Woz{
			ID:        123,
			Status:    404, // Assuming this is a valid status
			ScrapedAt: time.Now(),
			Payload:   []byte(`{"wozObject":{"wozobjectnummer":11800039614,"woonplaatsnaam":"Hoogeveen","openbareruimtenaam":"Anjerstraat","straatnaam":"Anjerstraat","postcode":"7906LE","huisnummer":5,"huisletter":null,"huisnummertoevoeging":null,"locatieomschrijving":null,"gemeentecode":118,"grondoppervlakte":122,"adresseerbaarobjectid":118010900087577,"nummeraanduidingid":null,"verbondenAdresseerbareObjecten":[118010900087577],"ontleendeAdresseerbareObjecten":[118010900087577]},"wozWaarden":[],"panden":[{"bagpandidentificatie":118100001200409}],"kadastraleObjecten":[{"kadastraleGemeenteCode":"HGV00","kadastraleSectie":"I","kadastraalPerceelNummer":"4945"}]}`),
		}

		// Call the Store method
		err := repo.Store(dummyWoz)

		// Assert that an error occurred
		assert.Error(t, err)
	})

	t.Run("Reset_NummerAanduidingID", func(t *testing.T) {
		// Mock the response for the API endpoint
		httpmock.RegisterResponder("POST", "http://example.com/api",
			func(req *http.Request) (*http.Response, error) {
				// Verify that nummeraanduidingid has been reset
				body, _ := ioutil.ReadAll(req.Body)
				reqBody := string(body)
				expectedRequestBody := `{"wozObject":{"wozobjectnummer":11800039614,"woonplaatsnaam":"Hoogeveen","openbareruimtenaam":"Anjerstraat","straatnaam":"Anjerstraat","postcode":"7906LE","huisnummer":5,"huisletter":null,"huisnummertoevoeging":null,"locatieomschrijving":null,"gemeentecode":118,"grondoppervlakte":122,"adresseerbaarobjectid":118010900087577,"nummeraanduidingid":123,"verbondenAdresseerbareObjecten":[118010900087577],"ontleendeAdresseerbareObjecten":[118010900087577]},"wozWaarden":[],"panden":[{"bagpandidentificatie":118100001200409}],"kadastraleObjecten":[{"kadastraleGemeenteCode":"HGV00","kadastraleSectie":"I","kadastraalPerceelNummer":"4945"}]}`
				assert.Equal(t, expectedRequestBody, reqBody)

				// Return a success response
				resp := httpmock.NewStringResponse(http.StatusOK, "success")
				return resp, nil
			},
		)

		// Create an instance of the repository
		repo := NewApiWozRepository("http://example.com/api", http.DefaultClient)

		// Prepare a dummy Woz object with nummeraanduidingid as null
		dummyWoz := domain.Woz{
			ID:        123,
			Status:    200, // Assuming this is a valid status
			ScrapedAt: time.Now(),
			Payload:   []byte(`{"wozObject":{"wozobjectnummer":11800039614,"woonplaatsnaam":"Hoogeveen","openbareruimtenaam":"Anjerstraat","straatnaam":"Anjerstraat","postcode":"7906LE","huisnummer":5,"huisletter":null,"huisnummertoevoeging":null,"locatieomschrijving":null,"gemeentecode":118,"grondoppervlakte":122,"adresseerbaarobjectid":118010900087577,"nummeraanduidingid":null,"verbondenAdresseerbareObjecten":[118010900087577],"ontleendeAdresseerbareObjecten":[118010900087577]},"wozWaarden":[],"panden":[{"bagpandidentificatie":118100001200409}],"kadastraleObjecten":[{"kadastraleGemeenteCode":"HGV00","kadastraleSectie":"I","kadastraalPerceelNummer":"4945"}]}`),
		}

		// Call the Store method
		err := repo.Store(dummyWoz)

		// Assert that no error occurred
		assert.NoError(t, err)
	})
}

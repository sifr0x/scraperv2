package massapi

type WOZObject struct {
	WozObjectNummer           int64   `json:"wozobjectnummer"`
	WoonplaatsNaam            string  `json:"woonplaatsnaam"`
	OpenbareRuimteNaam        string  `json:"openbareruimtenaam"`
	StraatNaam                string  `json:"straatnaam"`
	Postcode                  string  `json:"postcode"`
	Huisnummer                int     `json:"huisnummer"`
	HuisLetter                *string `json:"huisletter"`
	HuisnummerToevoeging      *string `json:"huisnummertoevoeging"`
	LocatieOmschrijving       *string `json:"locatieomschrijving"`
	GemeenteCode              int     `json:"gemeentecode"`
	GrondOppervlakte          int     `json:"grondoppervlakte"`
	AdresseerbaarObjectID     int64   `json:"adresseerbaarobjectid"`
	NummerAanduidingID        int64   `json:"nummeraanduidingid"`
	VerbondenAdresseerbareObj []int64 `json:"verbondenAdresseerbareObjecten"`
	OntleendeAdresseerbareObj []int64 `json:"ontleendeAdresseerbareObjecten"`
}

type WOZWaarde struct {
	Peildatum          string `json:"peildatum"`
	VastgesteldeWaarde int    `json:"vastgesteldeWaarde"`
}

type KadastraalObject struct {
	KadastraleGemeenteCode  string `json:"kadastraleGemeenteCode"`
	KadastraleSectie        string `json:"kadastraleSectie"`
	KadastraalPerceelNummer string `json:"kadastraalPerceelNummer"`
}

type Data struct {
	WOZObject          WOZObject          `json:"wozObject"`
	WOZWaarden         []WOZWaarde        `json:"wozWaarden"`
	Panden             []interface{}      `json:"panden"` // Assuming type of Panden is unknown
	KadastraleObjecten []KadastraalObject `json:"kadastraleObjecten"`
}

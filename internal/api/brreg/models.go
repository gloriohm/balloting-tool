package brreg

type Hovedenhet struct {
	Orgnummer      string `json:"organisasjonsnummer"`
	Name           string `json:"navn"`
	Konkurs        bool   `json:"konkurs"`
	Avvikles       bool   `json:"underAvvikling"`
	Tvangsavvikles bool   `json:"underTvangsavviklingEllerTvangsopplosning"`
}

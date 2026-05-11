# SN-utils
SN-utils er et enkelt program bestående av flere verktøy som skal hjelpe med å gjøre kronglete oppgaver litt enklere. Om den kjøres uten noe preferanser, vil den kjøre standard jobb for å generere excel-liste med utestående avstemninger.
## De forskjellige verktøyene
SN-utils har i dag tre forskjellige verktøy:
1. ballots -- for å generere rapport for utestående avstemninger
2. standards -- for å telle antall standarder i et regneark
3. table -- for å konvertere tabeller fra XML til Excel (for øyeblikket ikke tilgjengelig)

For å bruke noe annet enn ballot-verktøyet, må man kjøre programmet fra en terminal (eks. PowerShell) og spesifisere verktøy med et såkalt flagg. Verktøy velges med flagget "-tool" etterfulgt av navnet på ønsket verktøy.

`./sn-utils.exe -tool standards` vil for eksempel kjøre standard-verktøyet med default innstillinger.
### Ballots
Om ikke noe blir spesifisert, forsøker programmet å kjøre ballots. Man kan også eksplisitt påkalle det med `./sn-utils.exe -tool ballots`.

Ballots-verktøyet krever to lister med ballots (en for CEN, en for ISO) og en liste over roller i input-mappen. Programmet vil samle de to listene med ballots til én liste og forsøke og matche dem med de som har rollen Voter og CEN Voter i rollelisten. Programmet bruker referansen på komité for å knytte ballot til Voter.

Hvis programmet har kjørt vellykket, vil den generere en Excel-fil med navn utestående_avstemninger-yyyy-mm-dd (der år-måned-år er dagens dato) i output-mappen. I tillegg vil den generere en fil som heter missing.txt som lister opp alle ballots der den ikke klarte å matche ballot med Voter. Om rollelisten som ble brukt til input er komplett og up-to-date, er dette et tegn på at komiteen der balloten er opprettet mangler Voter fra SN. Dersom missing.txt ender opp med å være veldig lang, er dette et tegn på at input-dataen ikke er korrekt formatert eller ufullstendig.
#### Centralized Voters
Centralized Voters må spesifiseres i config.json (se under for mer informasjon). Dette kan være en eller flere e-postadresser. Den kan også være tom.

Alle e-postadresser i centralizedVoters blir filtrert ut av hovedrapporten kalt utestående_avstemninger-yyyy-mm-dd. Assosierte ballots blir filtrert inn i egne ark (sheet på engelsk) ved navn Centralized1, Centralized2 osv.
### Standards
Standards-verktøyet forventer en eller flere lister i .csv eller .xlsx-format med en kolonne som heter References. Programmet henter disse fra input-mappen. Det vil samle all dataen fra denne kolonnen på tvers av filene og gir ut et tall på hvor mange treff det er. Operasjonene som benyttes kan configureres ved å bruke ytterligere flagg når programmet påkalles.

I dag vil programmet alltid:
+ Filtrere ut standarderer med lik referanse slik at et produkt bare teller én gang
+ Filtere ut standarder med språkkode eller tilleggskode i referanse (eks. NS 9401.E:1994 eller NS 6033.T:1981)
+ Filtrere ut alle tilleggsprodukter (eks. /NA, /AC, /A1, /G1...)
+ Filtrere ut alle språk som ikke har norsk (no, nb, nn) eller engelsk språkkode
#### Job
Det finnes fire forskjellige hovedfiltreringer standards-verktøyet kan utføre.
1. all (default) -- gjør ingen ytterligere filtreringer (obs: filterer i dag ut NORSOK)
2. national -- teller kun nasjonalt utviklede produkter
3. adoption -- teller kun adopsjoner
4. norsok -- teller kun norsok

Alle disse kan brukes med flagget -ns_only for å kunne telle referanser som er Norsk Standard -- altså begynner med NS.

eks:
`./ sn-utils.exe -tool standards -job national -ns_only` vil gi tall på alle unike, egenutviklede NSer som er med i kildedokumentet.
## config.json
config.json brukes til å overstyre programmets basis-innstillinger. Dersom ingenting er satt i config.json, brukes følgende defaults:
+ inputPath: /Users/ditt_brukernavn/downloads
+ outputPath: /Users/ditt_brukernavn/downloads/ballot_resultat
+ fileNames: 
	+ ballot1: iso_ballots.xlsx
	+ ballot2: cen_ballots.xlsx
	+ voterRoles: roles.csv

centralizedVoters må oppgis i en liste med verdiene i anførselstegn, separert med komma.
eks: `centralizedVoters: ["abc@standard.no", "xyz@standard.no"]`

Filnavn må inkludere filformat. .csv og .xlsx kan benyttes om hverandre, så lenge det er konsekvent med filtypen angitt under fileNames.
package ogame

import (
	"net/url"
	"sync/atomic"
	"time"
)

// Priorities
const (
	Low       = 1
	Normal    = 2
	Important = 3
	Critical  = 4
)

// Prioritize ...
type Prioritize struct {
	bot          *OGame
	name         string
	taskIsDoneCh chan struct{}
	isTx         int32
}

// Begin a new transaction. "Done" must be called to release the lock.
func (b *Prioritize) Begin() *Prioritize {
	return b.begin("Tx")
}

// Done terminate the transaction, release the lock.
func (b *Prioritize) Done() {
	b.done()
}

func (b *Prioritize) begin(name string) *Prioritize {
	atomic.AddInt32(&b.isTx, 1)
	if atomic.LoadInt32(&b.isTx) == 1 {
		b.name = name
		b.bot.botLock(name)
	}
	return b
}

func (b *Prioritize) done() {
	atomic.AddInt32(&b.isTx, -1)
	if atomic.LoadInt32(&b.isTx) == 0 {
		defer close(b.taskIsDoneCh)
		b.bot.botUnlock(b.name)
	}
}

// Tx locks the bot during the transaction and ensure the lock is released afterward
func (b *Prioritize) Tx(clb func(*Prioritize) error) error {
	tx := b.Begin()
	defer tx.Done()
	err := clb(tx)
	return err
}

// FakeCall used for debugging
func (b *Prioritize) FakeCall(name string, delay int) {
	b.begin("FakeCall")
	defer b.done()
	b.bot.fakeCall(name, delay)
}

// Login to ogame server
// Can fails with BadCredentialsError
func (b *Prioritize) Login() error {
	b.begin("Login")
	defer b.done()
	return b.bot.wrapLogin()
}

// Logout the bot from ogame server
func (b *Prioritize) Logout() {
	b.begin("Logout")
	defer b.done()
	b.bot.logout()
}

// GetAlliancePageContent gets the html for a specific ogame page
func (b *Prioritize) GetAlliancePageContent(vals url.Values) []byte {
	b.begin("GetAlliancePageContent")
	defer b.done()
	pageHTML, _ := b.bot.getAlliancePageContent(vals)
	return pageHTML
}

// GetPageContent gets the html for a specific ogame page
func (b *Prioritize) GetPageContent(vals url.Values) []byte {
	b.begin("GetPageContent")
	defer b.done()
	pageHTML, _ := b.bot.getPageContent(vals)
	return pageHTML
}

// PostPageContent make a post request to ogame server
// This is useful when simulating a web browser
func (b *Prioritize) PostPageContent(vals, payload url.Values) []byte {
	b.begin("PostPageContent")
	defer b.done()
	by, _ := b.bot.postPageContent(vals, payload)
	return by
}

// IsUnderAttack returns true if the user is under attack, false otherwise
func (b *Prioritize) IsUnderAttack() bool {
	b.begin("IsUnderAttack")
	defer b.done()
	return b.bot.isUnderAttack()
}

// GetPlanets returns the user planets
func (b *Prioritize) GetPlanets() []Planet {
	b.begin("GetPlanets")
	defer b.done()
	return b.bot.getPlanets()
}

// GetPlanet gets infos for planetID
// Fails if planetID is invalid
func (b *Prioritize) GetPlanet(v interface{}) (Planet, error) {
	b.begin("GetPlanet")
	defer b.done()
	return b.bot.getPlanet(v)
}

// GetMoons returns the user moons
func (b *Prioritize) GetMoons() []Moon {
	b.begin("GetMoons")
	defer b.done()
	return b.bot.getMoons()
}

// GetMoon gets infos for moonID
func (b *Prioritize) GetMoon(v interface{}) (Moon, error) {
	b.begin("GetMoon")
	defer b.done()
	return b.bot.getMoon(v)
}

// GetCelestials get the player's planets & moons
func (b *Prioritize) GetCelestials() ([]Celestial, error) {
	b.begin("GetCelestials")
	defer b.done()
	return b.bot.getCelestials()
}

// Abandon a planet. Warning: this is irreversible
func (b *Prioritize) Abandon(v interface{}) error {
	b.begin("Abandon")
	defer b.done()
	return b.bot.abandon(v)
}

// GetCelestial get the player's planet/moon using the coordinate
func (b *Prioritize) GetCelestial(v interface{}) (Celestial, error) {
	b.begin("GetCelestial")
	defer b.done()
	return b.bot.getCelestial(v)
}

// ServerTime returns server time
// Timezone is OGT (OGame Time zone)
func (b *Prioritize) ServerTime() time.Time {
	b.begin("ServerTime")
	defer b.done()
	return b.bot.serverTime()
}

// GetUserInfos gets the user information
func (b *Prioritize) GetUserInfos() UserInfos {
	b.begin("GetUserInfos")
	defer b.done()
	return b.bot.getUserInfos()
}

// SendMessage sends a message to playerID
func (b *Prioritize) SendMessage(playerID int, message string) error {
	b.begin("SendMessage")
	defer b.done()
	return b.bot.sendMessage(playerID, message)
}

// GetFleets get the player's own fleets activities
func (b *Prioritize) GetFleets() ([]Fleet, Slots) {
	b.begin("GetFleets")
	defer b.done()
	return b.bot.getFleets()
}

// GetFleetsFromEventList get the player's own fleets activities
func (b *Prioritize) GetFleetsFromEventList() []Fleet {
	b.begin("GetFleets")
	defer b.done()
	return b.bot.getFleetsFromEventList()
}

// CancelFleet cancel a fleet
func (b *Prioritize) CancelFleet(fleetID FleetID) error {
	b.begin("CancelFleet")
	defer b.done()
	return b.bot.cancelFleet(fleetID)
}

// GetAttacks get enemy fleets attacking you
func (b *Prioritize) GetAttacks() []AttackEvent {
	b.begin("GetAttacks")
	defer b.done()
	return b.bot.getAttacks()
}

// GalaxyInfos get information of all planets and moons of a solar system
func (b *Prioritize) GalaxyInfos(galaxy, system int) (SystemInfos, error) {
	b.begin("GalaxyInfos")
	defer b.done()
	return b.bot.galaxyInfos(galaxy, system)
}

// GetResourceSettings gets the resources settings for specified planetID
func (b *Prioritize) GetResourceSettings(planetID PlanetID) (ResourceSettings, error) {
	b.begin("GetResourceSettings")
	defer b.done()
	return b.bot.getResourceSettings(planetID)
}

// SetResourceSettings set the resources settings on a planet
func (b *Prioritize) SetResourceSettings(planetID PlanetID, settings ResourceSettings) error {
	b.begin("SetResourceSettings")
	defer b.done()
	return b.bot.setResourceSettings(planetID, settings)
}

// GetResourcesBuildings gets the resources buildings levels
func (b *Prioritize) GetResourcesBuildings(celestialID CelestialID) (ResourcesBuildings, error) {
	b.begin("GetResourcesBuildings")
	defer b.done()
	return b.bot.getResourcesBuildings(celestialID)
}

// GetDefense gets all the defenses units information of a planet
// Fails if planetID is invalid
func (b *Prioritize) GetDefense(celestialID CelestialID) (DefensesInfos, error) {
	b.begin("GetDefense")
	defer b.done()
	return b.bot.getDefense(celestialID)
}

// GetShips gets all ships units information of a planet
func (b *Prioritize) GetShips(celestialID CelestialID) (ShipsInfos, error) {
	b.begin("GetShips")
	defer b.done()
	return b.bot.getShips(celestialID)
}

// GetFacilities gets all facilities information of a planet
func (b *Prioritize) GetFacilities(celestialID CelestialID) (Facilities, error) {
	b.begin("GetFacilities")
	defer b.done()
	return b.bot.getFacilities(celestialID)
}

// GetProduction get what is in the production queue.
// (ships & defense being built)
func (b *Prioritize) GetProduction(celestialID CelestialID) ([]Quantifiable, error) {
	b.begin("GetProduction")
	defer b.done()
	return b.bot.getProduction(celestialID)
}

// GetResearch gets the player researches information
func (b *Prioritize) GetResearch() Researches {
	b.begin("GetResearch")
	defer b.done()
	return b.bot.getResearch()
}

// GetSlots gets the player current and total slots information
func (b *Prioritize) GetSlots() Slots {
	b.begin("GetSlots")
	defer b.done()
	return b.bot.getSlots()
}

// Build builds any ogame objects (building, technology, ship, defence)
func (b *Prioritize) Build(celestialID CelestialID, id ID, nbr int) error {
	b.begin("Build")
	defer b.done()
	return b.bot.build(celestialID, id, nbr)
}

// BuildCancelable builds any cancelable ogame objects (building, technology)
func (b *Prioritize) BuildCancelable(celestialID CelestialID, id ID) error {
	b.begin("BuildCancelable")
	defer b.done()
	return b.bot.buildCancelable(celestialID, id)
}

// BuildProduction builds any line production ogame objects (ship, defence)
func (b *Prioritize) BuildProduction(celestialID CelestialID, id ID, nbr int) error {
	b.begin("BuildProduction")
	defer b.done()
	return b.bot.buildProduction(celestialID, id, nbr)
}

// BuildBuilding ensure what is being built is a building
func (b *Prioritize) BuildBuilding(celestialID CelestialID, buildingID ID) error {
	b.begin("BuildBuilding")
	defer b.done()
	return b.bot.buildBuilding(celestialID, buildingID)
}

// BuildDefense builds a defense unit
func (b *Prioritize) BuildDefense(celestialID CelestialID, defenseID ID, nbr int) error {
	b.begin("BuildDefense")
	defer b.done()
	return b.bot.buildDefense(celestialID, defenseID, nbr)
}

// BuildShips builds a ship unit
func (b *Prioritize) BuildShips(celestialID CelestialID, shipID ID, nbr int) error {
	b.begin("BuildShips")
	defer b.done()
	return b.bot.buildShips(celestialID, shipID, nbr)
}

// ConstructionsBeingBuilt returns the building & research being built, and the time remaining (secs)
func (b *Prioritize) ConstructionsBeingBuilt(celestialID CelestialID) (ID, int, ID, int) {
	b.begin("ConstructionsBeingBuilt")
	defer b.done()
	return b.bot.constructionsBeingBuilt(celestialID)
}

// CancelBuilding cancel the construction of a building on a specified planet
func (b *Prioritize) CancelBuilding(celestialID CelestialID) error {
	b.begin("CancelBuilding")
	defer b.done()
	return b.bot.cancelBuilding(celestialID)
}

// CancelResearch cancel the research
func (b *Prioritize) CancelResearch(celestialID CelestialID) error {
	b.begin("CancelResearch")
	defer b.done()
	return b.bot.cancelResearch(celestialID)
}

// BuildTechnology ensure that we're trying to build a technology
func (b *Prioritize) BuildTechnology(celestialID CelestialID, technologyID ID) error {
	b.begin("BuildTechnology")
	defer b.done()
	return b.bot.buildTechnology(celestialID, technologyID)
}

// GetResources gets user resources
func (b *Prioritize) GetResources(celestialID CelestialID) (Resources, error) {
	b.begin("GetResources")
	defer b.done()
	return b.bot.getResources(celestialID)
}

// SendFleet sends a fleet
func (b *Prioritize) SendFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate,
	mission MissionID, resources Resources, expeditiontime int) (Fleet, error) {
	b.begin("SendFleet")
	defer b.done()
	return b.bot.sendFleet(celestialID, ships, speed, where, mission, resources, expeditiontime, false)
}

// EnsureFleet either sends all the requested ships or fail
func (b *Prioritize) EnsureFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate,
	mission MissionID, resources Resources, expeditiontime int) (Fleet, error) {
	b.begin("EnsureFleet")
	defer b.done()
	return b.bot.sendFleet(celestialID, ships, speed, where, mission, resources, expeditiontime, true)
}

// SendIPM sends IPM
func (b *Prioritize) SendIPM(planetID PlanetID, coord Coordinate, nbr int, priority ID) (int, error) {
	b.begin("SendIPM")
	defer b.done()
	return b.bot.sendIPM(planetID, coord, nbr, priority)
}

// GetCombatReportSummaryFor gets the latest combat report for a given coordinate
func (b *Prioritize) GetCombatReportSummaryFor(coord Coordinate) (CombatReportSummary, error) {
	b.begin("GetCombatReportSummaryFor")
	defer b.done()
	return b.bot.getCombatReportFor(coord)
}

// GetEspionageReportFor gets the latest espionage report for a given coordinate
func (b *Prioritize) GetEspionageReportFor(coord Coordinate) (EspionageReport, error) {
	b.begin("GetEspionageReportFor")
	defer b.done()
	return b.bot.getEspionageReportFor(coord)
}

// GetEspionageReportMessages gets the summary of each espionage reports
func (b *Prioritize) GetEspionageReportMessages() ([]EspionageReportSummary, error) {
	b.begin("GetEspionageReportMessages")
	defer b.done()
	return b.bot.getEspionageReportMessages()
}

// GetEspionageReport gets a detailed espionage report
func (b *Prioritize) GetEspionageReport(msgID int) (EspionageReport, error) {
	b.begin("GetEspionageReport")
	defer b.done()
	return b.bot.getEspionageReport(msgID)
}

// DeleteMessage deletes a message from the mail box
func (b *Prioritize) DeleteMessage(msgID int) error {
	b.begin("DeleteMessage")
	defer b.done()
	return b.bot.deleteMessage(msgID)
}

// GetResourcesProductions gets the planet resources production
func (b *Prioritize) GetResourcesProductions(planetID PlanetID) (Resources, error) {
	b.begin("GetResourcesProductions")
	defer b.done()
	return b.bot.getResourcesProductions(planetID)
}

// GetResourcesProductionsLight gets the planet resources production
func (b *Prioritize) GetResourcesProductionsLight(resBuildings ResourcesBuildings, researches Researches,
	resSettings ResourceSettings, temp Temperature) Resources {
	b.begin("GetResourcesProductionsLight")
	defer b.done()
	return b.bot.getResourcesProductionsLight(resBuildings, researches, resSettings, temp)
}

// FlightTime calculate flight time and fuel needed
func (b *Prioritize) FlightTime(origin, destination Coordinate, speed Speed, ships ShipsInfos) (secs, fuel int) {
	if b.bot.researches == nil {
		b.begin("FlightTime")
		b.bot.getResearch()
		b.done()
	} else {
		if atomic.LoadInt32(&b.isTx) == 0 {
			defer close(b.taskIsDoneCh)
		}
	}
	return calcFlightTime(origin, destination, b.bot.universeSize, b.bot.donutGalaxy, b.bot.donutSystem, b.bot.fleetDeutSaveFactor,
		float64(speed)/10, b.bot.universeSpeedFleet, ships, *b.bot.researches)
}

// Phalanx scan a coordinate from a moon to get fleets information
// IMPORTANT: My account was instantly banned when I scanned an invalid coordinate.
// IMPORTANT: This function DOES validate that the coordinate is a valid planet in range of phalanx
// 			  and that you have enough deuterium.
func (b *Prioritize) Phalanx(moonID MoonID, coord Coordinate) ([]Fleet, error) {
	b.begin("Phalanx")
	defer b.done()
	return b.bot.getPhalanx(moonID, coord)
}

// UnsafePhalanx same as Phalanx but does not perform any input validation.
func (b *Prioritize) UnsafePhalanx(moonID MoonID, coord Coordinate) ([]Fleet, error) {
	b.begin("Phalanx")
	defer b.done()
	return b.bot.getUnsafePhalanx(moonID, coord)
}

// JumpGate sends ships through a jump gate.
func (b *Prioritize) JumpGate(origin, dest MoonID, ships ShipsInfos) error {
	b.begin("JumpGate")
	defer b.done()
	return b.bot.executeJumpGate(origin, dest, ships)
}

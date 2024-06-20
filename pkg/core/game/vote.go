package game

import (
	"errors"
	"github.com/https-whoyan/MafiaBot/core/converter"
	"strconv"

	"github.com/https-whoyan/MafiaBot/core/channel"
	"github.com/https-whoyan/MafiaBot/core/player"
)

// This file contains everything about the voting mechanics.

// ___________________________________
// VoteProviderInterface
// ___________________________________

const (
	EmptyVoteStr = "-1"
	EmptyVoteInt = -1
)

// VoteProviderInterface Interface for voice reception (both daytime and nighttime).
//
// This allows you to implement the Vote command in your interpretation.
type VoteProviderInterface interface {
	// GetVotedPlayerID Provides 2 fields: information about the voting player.
	//
	// The 1st field provides the ID of the player who voted,
	// the second field is whether this ID is your server ID or in-game ID.
	GetVotedPlayerID() (votedUserID string, isUserServerID bool)
	// GetVote provide one field: ID of the player being voted for.
	// If you need empty vote, use the EmptyVoteStr constant.
	GetVote() (ID string)
}

// VoteProvider default implementation of VoteProviderInterface
type VoteProvider struct {
	votedPlayerID  string
	vote           string
	isServerUserID bool
}

func NewVoteProvider(votedPlayerID string, vote string, isServerUserID bool) *VoteProvider {
	return &VoteProvider{
		votedPlayerID:  votedPlayerID,
		vote:           vote,
		isServerUserID: isServerUserID,
	}
}

func (p *VoteProvider) GetVotedPlayerID() (votedUserID string, isUserServerID bool) {
	return p.votedPlayerID, p.isServerUserID
}
func (p *VoteProvider) GetVote() (ID string) { return p.vote }

// TwoVoteProviderInterface A special channel used only  for roles that specify 2 voices rather
// than one (example: detective)
//
// Its peculiarity is that instead of one voice it uses
// 2 voices - IDs of 2 players it wants to check, so I decided to make a separate interface for it
type TwoVoteProviderInterface interface {
	GetVotedPlayerID() (votedUserID string, isUserServerID bool)
	GetVote() (ID1 string, ID2 string)
}

// TwoVotesProvider default implementation of TwoVoteProviderInterface
type TwoVotesProvider struct {
	votedPlayerID  string
	vote1, vote2   string
	isServerUserID bool
}

func NewTwoVoteProvider(votedPlayerID string, vote1, vote2 string, isServerUserID bool) *TwoVotesProvider {
	return &TwoVotesProvider{
		votedPlayerID:  votedPlayerID,
		vote1:          vote1,
		vote2:          vote2,
		isServerUserID: isServerUserID,
	}
}

func (p *TwoVotesProvider) GetVotedPlayerID() (votedUserID string, isUserServerID bool) {
	return p.votedPlayerID, p.isServerUserID
}
func (p *TwoVotesProvider) GetVote() (ID1, ID2 string) { return p.vote1, p.vote2 }

// _______________________________
// Vote Validators
// _______________________________

var (
	NilValidatorErr      = errors.New("nil Validator")
	InVotePlayerNotFound = errors.New("voted player not found")
	IncorrectVoteType    = errors.New("incorrect vote type")
	IncorrectVoteChannel = errors.New("incorrect vote channel")
	IncorrectVotedPlayer = errors.New("incorrect voted player")
	VotePlayerNotFound   = errors.New("vote player not found")
	PlayerIsMutedErr     = errors.New("player is muted")
	VotePlayerIsNotAlive = errors.New("vote player is not alive")
)

// ________________________
// VoteProvider
// ________________________

// voteProviderValidator is validator for VoteProviderInterface
func (g *Game) voteProviderValidator(vP VoteProviderInterface) error {
	if vP == nil {
		return NilValidatorErr
	}
	votedPlayerID, isServerID := vP.GetVotedPlayerID()
	votedPlayer := player.SearchPlayerByID(g.Active, votedPlayerID, isServerID)
	if votedPlayer == nil {
		return InVotePlayerNotFound
	}
	vote := vP.GetVote()
	if vote == EmptyVoteStr {
		return nil
	}
	_, err := strconv.Atoi(vote)
	if err != nil {
		return IncorrectVoteType
	}
	if votedPlayer.LifeStatus == player.Alive {
		return VotePlayerIsNotAlive
	}
	toVotePlayer := player.SearchPlayerByID(g.Active, vote, false)
	if toVotePlayer == nil {
		return VotePlayerNotFound
	}
	return nil
}

// nightVoteValidatorByChannelIID performs the same validation as nightVoteValidator.
//
// Use it, if you want, that day vote should be in a particular channel.
func (g *Game) nightVoteValidatorByChannelIID(vP VoteProviderInterface, channelIID string) error {
	sliceChannels := converter.GetMapValues(g.RoleChannels)
	foundedChannel := channel.SearchRoleChannelByID(sliceChannels, channelIID)
	if foundedChannel == nil {
		return IncorrectVoteChannel
	}
	return g.nightVoteValidator(vP, foundedChannel)
}

// nightVoteValidator also check roleChannel.Role and vP.VotedPlayer role.
// Use nil if you don't need for this checking
func (g *Game) nightVoteValidator(vP VoteProviderInterface, roleChannel channel.RoleChannel) error {
	if err := g.voteProviderValidator(vP); err != nil {
		return err
	}

	votedPlayerID, isServerID := vP.GetVotedPlayerID()
	votedPlayer := player.SearchPlayerByID(g.Active, votedPlayerID, isServerID)
	if g.NightVoting != votedPlayer.Role {
		return IncorrectVotedPlayer
	}
	if roleChannel != nil && g.NightVoting != roleChannel.GetRole() {
		return IncorrectVoteChannel
	}
	if votedPlayer.InteractionStatus == player.Muted {
		return PlayerIsMutedErr
	}
	return nil
}

// dayVoteValidatorByChannelIID performs the same validation as dayVoteValidator
func (g *Game) dayVoteValidatorByChannelIID(vP VoteProviderInterface, channelIID string) error {
	var allChannels []channel.Channel
	allRoleChannels := converter.GetMapValues(g.RoleChannels)
	allChannels = append(allChannels, channel.RoleSliceToChannelSlice(allRoleChannels)...)

	allChannels = append(allChannels, g.MainChannel)

	channelVotedFrom := channel.SearchChannelByGameID(allChannels, channelIID)
	if channelVotedFrom == nil {
		return IncorrectVoteChannel
	}
	return g.dayVoteValidator(vP)
}

func (g *Game) dayVoteValidator(vP VoteProviderInterface) error {
	return g.voteProviderValidator(vP)
}

// ________________________
// TwoVoteProvider
// ________________________

// _______________________________
// Vote Functions
// _______________________________

// OptionalChannelIID Optional field Mechanism.
type OptionalChannelIID struct{ channelIID string }

func NewOptionalChannelIID(channelIID string) *OptionalChannelIID {
	return &OptionalChannelIID{channelIID}
}

// NightVote opt is OptionalChannelIID optional field Mechanism.
//
// If you not need it, pass nil to the field.
// If yes, use NewOptionalChannelIID
//
// Immediately puts all the right votes and changes the value of the fields if no error occurred.
func (g *Game) NightVote(vP VoteProviderInterface, opt *OptionalChannelIID) error {
	var err error
	if opt == nil {
		err = g.nightVoteValidator(vP, nil)
	} else {
		err = g.nightVoteValidatorByChannelIID(vP, opt.channelIID)
	}
	if err != nil {
		return err
	}

	votedPlayerID, isServerID := vP.GetVotedPlayerID()
	g.RLock()
	votedPlayer := player.SearchPlayerByID(g.Active, votedPlayerID, isServerID)
	g.RUnlock()
	vote := vP.GetVote()
	g.Lock()
	defer g.Lock()
	if vote == EmptyVoteStr {
		votedPlayer.Votes = append(votedPlayer.Votes, EmptyVoteInt)
	} else {
		// validated Before
		intVote, _ := strconv.Atoi(vote)
		votedPlayer.Votes = append(votedPlayer.Votes, intVote)
	}
	// Set empty votes to same role players
	sameRolePlayers := player.SearchAllPlayersWithRole(g.Active, votedPlayer.Role)
	for _, sameRolePlayer := range sameRolePlayers {
		if sameRolePlayer.ID != votedPlayer.ID {
			sameRolePlayer.Votes = append(sameRolePlayer.Votes, EmptyVoteInt)
		}
	}
	if votedPlayer.Role.UrgentCalculation {
		// Todo.
	}
	return nil
}

// DayVote opt is OptionalChannelIID optional field Mechanism.
//
// If you not need it, pass nil to the field.
// If yes, use NewOptionalChannelIID
//
// Immediately puts all the right votes and changes the value of the fields if no error occurred.
func (g *Game) DayVote(vP VoteProviderInterface, opt *OptionalChannelIID) error {
	var err error
	if opt == nil {
		err = g.dayVoteValidator(vP)
	} else {
		err = g.dayVoteValidatorByChannelIID(vP, opt.channelIID)
	}
	if err != nil {
		return err
	}

	votedPlayerID, isServerID := vP.GetVotedPlayerID()
	g.RLock()
	votedPlayer := player.SearchPlayerByID(g.Active, votedPlayerID, isServerID)
	g.RUnlock()
	vote := vP.GetVote()
	g.Lock()
	defer g.Lock()
	if vote == EmptyVoteStr {
		votedPlayer.DayVote = EmptyVoteInt
	}
	// validated Before
	votedPlayer.DayVote, _ = strconv.Atoi(vote)
	return nil
}

// ResetTheVotes use to reset all player votes
func (g *Game) ResetTheVotes() {
	g.Lock()
	defer g.Unlock()
	allPlayers := g.Active

	for _, activePlayer := range allPlayers {
		activePlayer.DayVote = EmptyVoteInt
	}
}

// ResetAllInteractionsStatuses use to reset all player interaction statuses
func (g *Game) ResetAllInteractionsStatuses() {
	g.Lock()
	allPlayers := g.Active
	defer g.Unlock()

	for _, activePlayer := range allPlayers {
		activePlayer.InteractionStatus = player.Passed
	}
}

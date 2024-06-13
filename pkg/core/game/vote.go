package game

import (
	"errors"
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

func (p *VoteProvider) GetVote() (ID string) {
	return p.vote
}

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
)

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
	foundedChannel := g.searchRoleChannelByIID(channelIID)
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
	return nil
}

// dayVoteValidatorByChannelIID performs the same validation as dayVoteValidator.
//
// Use it, if you want, that day vote should be in a particular channel.
func (g *Game) dayVoteValidatorByChannelIID(vP VoteProviderInterface, channelIID string) error {
	foundChannel := g.searchChannelByID(channelIID)
	if foundChannel == nil {
		return IncorrectVoteChannel
	}
	return g.dayVoteValidator(vP)
}

func (g *Game) dayVoteValidator(vP VoteProviderInterface) error {
	return g.voteProviderValidator(vP)
}

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
	g.Unlock()
	vote := vP.GetVote()
	g.RLock()
	defer g.RUnlock()
	if vote == EmptyVoteStr {
		votedPlayer.Vote = EmptyVoteInt
	}
	// validated Before
	votedPlayer.Vote, _ = strconv.Atoi(vote)
	return nil
}

// DayVote opt is OptionalChannelIID optional field Mechanism.
//
// If you not need it, pass nil to the field.
// If yes, use NewOptionalChannelIID
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
	g.RLock()
	defer g.RUnlock()
	if vote == EmptyVoteStr {
		votedPlayer.Vote = EmptyVoteInt
	}
	// validated Before
	votedPlayer.Vote, _ = strconv.Atoi(vote)
	return nil
}

// ResetTheVotes use to reset all player votes
func (g *Game) ResetTheVotes() {
	g.RLock()
	allPlayers := g.Active
	g.RUnlock()

	for _, activePlayer := range allPlayers {
		activePlayer.Vote = EmptyVoteInt
	}
}

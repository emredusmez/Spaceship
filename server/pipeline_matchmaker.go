package server

import (
	"log"
	"spaceship/socketapi"
)

func (p *Pipeline) matchmakerFind(session Session, envelope *socketapi.Envelope){
	incomingData := envelope.GetMatchFind()
	//TODO validate incomingData with game specs
	matchEntry, err := p.matchmaker.Find(session, incomingData.GameName, incomingData.QueueProperties)
	if err != nil {
		log.Println(err)
		session.Send(false, 0, &socketapi.Envelope{Cid: envelope.Cid, Message: &socketapi.Envelope_Error{Error: &socketapi.Error{
			Code:    int32(socketapi.Error_MATCH_JOIN_REJECTED),
			Message: "Could not find match.",
		}}})
	}

	session.Send(false, 0, &socketapi.Envelope{Cid: envelope.Cid, Message: &socketapi.Envelope_MatchEntry{MatchEntry:matchEntry}})
}

func (p *Pipeline) matchmakerJoin(session Session, envelope *socketapi.Envelope){
	incomingData := envelope.GetMatchJoin()

	//TODO validate incomingData with game specs
	game, err := p.matchmaker.Join(p, session, incomingData.MatchId)
	if err != nil {
		log.Println(err)
		session.Send(false, 0, &socketapi.Envelope{Cid: envelope.Cid, Message: &socketapi.Envelope_Error{Error: &socketapi.Error{
			Code:    int32(socketapi.Error_MATCH_JOIN_REJECTED),
			Message: "Could not join match.",
		}}})
	}

	ms := socketapi.MatchJoinResp{GameData: game}
	session.Send(false, 0, &socketapi.Envelope{Cid: envelope.Cid, Message: &socketapi.Envelope_MatchStart{MatchStart:&ms}})
}

func (p *Pipeline) matchmakerLeave(session Session, envelope *socketapi.Envelope){
	incomingData := envelope.GetMatchLeave()
	//TODO validate incomingData with game specs
	p.matchmaker.Leave(session, incomingData.MatchId)
	log.Println("MatchLeave received for game: ", incomingData.MatchId)
}
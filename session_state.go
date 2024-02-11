package al

/*
#include <abl_link.h>
*/
import "C"

import "time"

// SessionState is the representation of the current local state of a client in a Link Session.
//
// A session state represents a timeline and the start/stop
// state. The timeline is a representation of a mapping between time and
// beats for varying quanta. The start/stop state represents the user
// intention to start or stop transport at a specific time. Start stop
// synchronization is an optional feature that allows to share the user
// request to start or stop transport between a subgroup of peers in a
// Link session. When observing a change of start/stop state, audio
// playback of a peer should be started or stopped the same way it would
// have happened if the user had requested that change at the according
// time locally. The start/stop state can only be changed by the user.
// This means that the current local start/stop state persists when
// joining or leaving a Link session. After joining a Link session
// start/stop change requests will be communicated to all connected peers.
type SessionState struct {
	instance C.abl_link_session_state
}

// NewSessionState creates a new session_state instance.
//
// Thread-safe: yes
// Realtime-safe: no
//
// The session_state is to be used with the link.Capture[...]SessionState and
// link.Commit[...]SessionState functions to capture snapshots of the current link state and pass
// changes to the link session.
func NewSessionState() *SessionState {
	return &SessionState{
		instance: C.abl_link_create_session_state(),
	}
}

// Close deletes a SessionState instance.
//
// Thread-safe: yes
// Realtime-safe: no
func (s *SessionState) Close() {
	C.abl_link_destroy_session_state(s.instance)
}

// Tempo returns the tempo of the timeline, in Beats Per Minute.
//
// This is a stable value that is appropriate for display to the user. Beat
// time progress will not necessarily match this tempo exactly because of clock drift
// compensation.
func (s *SessionState) Tempo() float64 {
	return float64(C.abl_link_tempo(s.instance))
}

// SetTempo sets the timeline tempo to the given bpm value, taking effect at the given time.
func (s *SessionState) SetTempo(bpm float64, atTime time.Duration) {
	C.abl_link_set_tempo(s.instance, C.double(bpm), C.int64_t(atTime.Microseconds()))
}

// BeatAtTime gets the beat value corresponding to the given time for the given quantum.
//
// The magnitude of the resulting beat value is unique to this Link
// client, but its phase with respect to the provided quantum is shared among all
// session peers. For non-negative beat values, the following property holds:
//
// fmod(beatAtTime(t, q), q) == phaseAtTime(t, q)
func (s *SessionState) BeatAtTime(atTime time.Duration, quantum float64) float64 {
	return float64(C.abl_link_beat_at_time(s.instance, C.int64_t(atTime.Microseconds()), C.double(quantum)))
}

// PhaseAtTime gets the session phase at the given time for the given quantum.
//
// The result is in the interval [0, quantum). The result is equivalent to
// fmod(beatAtTime(t, q), q) for non-negative beat values. This function is convenient
// if the client application is only interested in the phase and not the beat
// magnitude. Also, unlike fmod, it handles negative beat values correctly.
func (s *SessionState) PhaseAtTime(atTime time.Duration, quantum float64) float64 {
	return float64(C.abl_link_phase_at_time(s.instance, C.int64_t(atTime.Microseconds()), C.double(quantum)))
}

// TimeAtBeat gets the time at which the given beat occurs for the given quantum.
//
// @discussion: The inverse of beatAtTime, assuming a constant tempo.
//
// beatAtTime(timeAtBeat(b, q), q) === b.
func (s *SessionState) TimeAtBeat(beat, quantum float64) float64 {
	return float64(C.abl_link_time_at_beat(s.instance, C.double(beat), C.double(quantum)))
}

// RequestBeatAtTime attempts to map the given beat to the given time in the context of the given quantum.
//
// This function behaves differently depending on the state of the
// session. If no other peers are connected, then this abl_link instance is in a
// session by itself and is free to re-map the beat/time relationship whenever it
// pleases. In this case, beatAtTime(time, quantum) == beat after this funtion has been
// called.
//
// If there are other peers in the session, this Link instance should not abruptly
// re-map the beat/time relationship in the session because that would lead to beat
// discontinuities among the other peers. In this case, the given beat will be mapped
// to the next time value greater than the given time with the same phase as the given
// beat.
//
// This function is specifically designed to enable the concept of "quantized launch"
// in client applications. If there are no other peers in the session, then an event
// (such as starting transport) happens immediately when it is requested. If there are
// other peers, however, we wait until the next time at which the session phase matches
// the phase of the event, thereby executing the event in-phase with the other peers in
// the session. The client application only needs to invoke this function to achieve
// this behavior and should not need to explicitly check the number of peers.
func (s *SessionState) RequestBeatAtTime(beat float64, atTime time.Duration, quantum float64) {
	C.abl_link_request_beat_at_time(s.instance, C.double(beat), C.int64_t(atTime.Microseconds()), C.double(quantum))
}

// ForceBeatAtTime rudely re-maps the beat/time relationship for all peers in a session.
//
// DANGER: This function should only be needed in certain special
// circumstances. Most applications should not use it. It is very similar to
// RequestBeatAtTime except that it does not fall back to the quantizing behavior when
// it is in a session with other peers. Calling this function will unconditionally map
// the given beat to the given time and broadcast the result to the session. This is
// very anti-social behavior and should be avoided.
//
// One of the few legitimate uses of this function is to synchronize a Link session
// with an external clock source. By periodically forcing the beat/time mapping
// according to an external clock source, a peer can effectively bridge that clock into
// a Link session. Much care must be taken at the application layer when implementing
// such a feature so that users do not accidentally disrupt Link sessions that they may
// join.
func (s *SessionState) ForceBeatAtTime(beat float64, atTime time.Duration, quantum float64) {
	C.abl_link_force_beat_at_time(s.instance, C.double(beat), C.uint64_t(atTime.Microseconds()), C.double(quantum))
}

// SetIsPlaying sets if transport should be playing or stopped, taking effect at the given time.
func (s *SessionState) SetIsPlaying(isPlaying bool, atTime time.Duration) {
	C.abl_link_set_is_playing(s.instance, C.bool(isPlaying), C.uint64_t(atTime.Microseconds()))
}

// IsPlaying returns if transport is playing.
func (s *SessionState) IsPlaying() bool {
	return bool(C.abl_link_is_playing(s.instance))
}

// TimeForIsPlaying gets the time at which a transport start/stop occurs.
func (s *SessionState) TimeForIsPlaying() time.Duration {
	return time.Duration((C.abl_link_time_for_is_playing(s.instance))) * time.Microsecond
}

// RequestBeatAtStartPlayingTime is a convenience function which attempts to map the given beat to the time
// when transport is starting to play in context of the given quantum. This function evaluates to a
// no-op if IsPlaying equals false.
func (s *SessionState) RequestBeatAtStartPlayingTime(beat, quantum float64) {
	(C.abl_link_request_beat_at_start_playing_time(s.instance, C.double(beat), C.double(quantum)))
}

// SetIsPlayingAndRequestBeatAtTime is a convenience function which starts or stop transport at a given time and attempt
// to map the given beat to this time in context of the given quantum.
func (s *SessionState) SetIsPlayingAndRequestBeatAtTime(isPlaying bool, atTime time.Duration, beat, quantum float64) {
	C.abl_link_set_is_playing_and_request_beat_at_time(s.instance, C.bool(isPlaying), C.uint64_t(atTime.Microseconds()), C.double(beat), C.double(quantum))
}

package al

/*
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS: -labl_link
#include "abl_link.h"
*/
import "C"

import "time"

type SessionState struct {
	instance C.abl_link_session_state
}

func NewSessionState() *SessionState {
	return &SessionState{
		instance: C.abl_link_create_session_state(),
	}
}

func (s *SessionState) Close() {
	C.abl_link_destroy_session_state(s.instance)
}

func (s *SessionState) Tempo() float64 {
	return float64(C.abl_link_tempo(s.instance))
}

func (s *SessionState) SetTempo(bpm float64, atTime time.Time) {
	C.abl_link_set_tempo(s.instance, C.double(bpm), C.int64_t(atTime.UnixMicro()))
}

func (s *SessionState) BeatAtTime(atTime time.Time, quantum float64) float64 {
	return float64(C.abl_link_beat_at_time(s.instance, C.int64_t(atTime.UnixMicro()), C.double(quantum)))
}

func (s *SessionState) PhaseAtTime(atTime time.Time, quantum float64) float64 {
	return float64(C.abl_link_phase_at_time(s.instance, C.int64_t(atTime.UnixMicro()), C.double(quantum)))
}

func (s *SessionState) TimeAtBeat(beat, quantum float64) float64 {
	return float64(C.abl_link_time_at_beat(s.instance, C.double(beat), C.double(quantum)))
}

func (s *SessionState) RequestBeatAtTime(beat float64, atTime time.Time, quantum float64) {
	C.abl_link_request_beat_at_time(s.instance, C.double(beat), C.int64_t(atTime.UnixMicro()), C.double(quantum))
}

func (s *SessionState) ForceBeatAtTime(beat float64, atTime time.Time, quantum float64) {
	C.abl_link_force_beat_at_time(s.instance, C.double(beat), C.uint64_t(atTime.UnixMicro()), C.double(quantum))
}

func (s *SessionState) SetIsPlaying(isPlaying bool, atTime time.Time) {
	C.abl_link_set_is_playing(s.instance, C.bool(isPlaying), C.uint64_t(atTime.UnixMicro()))
}

func (s *SessionState) IsPlaying() bool {
	return bool(C.abl_link_is_playing(s.instance))
}

func (s *SessionState) TimeForIsPlaying() uint64 {
	return uint64((C.abl_link_time_for_is_playing(s.instance)))
}

func (s *SessionState) RequestBeatAtStartPlayingTime(beat, quantum float64) {
	(C.abl_link_request_beat_at_start_playing_time(s.instance, C.double(beat), C.double(quantum)))
}

func (s *SessionState) SetIsPlayingAndRequestBeatAtTime(isPlaying bool, atTime time.Time, beat, quantum float64) {
	C.abl_link_set_is_playing_and_request_beat_at_time(s.instance, C.bool(isPlaying), C.uint64_t(atTime.UnixMicro()), C.double(beat), C.double(quantum))
}

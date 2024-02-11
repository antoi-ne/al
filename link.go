package al

/*
#cgo LDFLAGS: -labl_link
#include <abl_link.h>

extern void goNumPeersCallback(uint64_t num_peers, void* context);
extern void goTempoCallback(double tempo, void* context);
extern void goStartStopCallback(bool is_playing, void* context);
*/
import "C"

import (
	"time"
	"unsafe"
)

type NumPeersCallbackFn func(numPeers uint64)

type TempoCallbackFn func(tempo float64)

type StartStopCallbackFn func(isPlaying bool)

// Link is the representation of an abl_link instance.
type Link struct {
	instance C.abl_link
}

// NewLink Construct a new abl_link instance with an initial tempo.
//
// Thread-safe: yes
// Realtime-safe: no
func NewLink(bpm float64) *Link {
	return &Link{
		instance: C.abl_link_create(C.double(bpm)),
	}
}

// Close deletes an abl_link instance.
//
// Thread-safe: yes
// Realtime-safe: no
func (l *Link) Close() {
	C.abl_link_destroy(l.instance)
}

// Enabled returns if Link is currently enabled.
//
// Thread-safe: yes
// Realtime-safe: yes
func (l *Link) Enabled() bool {
	return bool(C.abl_link_is_enabled(l.instance))
}

// Enable / Disable Link.
//
// Thread-safe: yes
// Realtime-safe: no
func (l *Link) Enable(enabled bool) {
	C.abl_link_enable(l.instance, C.bool(enabled))
}

// StartStopSyncEnabled returns if start/stop synchronization is enabled.
//
// Thread-safe: yes
// Realtime-safe: no
func (l *Link) StartStopSyncEnabled() bool {
	return bool(C.abl_link_is_start_stop_sync_enabled(l.instance))
}

// EnableStartStopSync enables/disables start/stop synchronization.
//
// Thread-safe: yes
// Realtime-safe: no
func (l *Link) EnableStartStopSync(enabled bool) {
	C.abl_link_enable_start_stop_sync(l.instance, C.bool(enabled))
}

// NumPeers returns how many peers are currently connected in a Link session.
//
// Thread-safe: yes
// Realtime-safe: yes
func (l *Link) NumPeers() uint64 {
	return uint64(C.abl_link_num_peers(l.instance))
}

// SetNumPeersCallback registers a callback to be notified when the number of peers in the Link session changes.
//
// Thread-safe: yes
// Realtime-safe: no
//
// The callback is invoked on a Link-managed thread.
func (l *Link) SetNumPeersCallback(fn NumPeersCallbackFn) {
	fnPtr := unsafe.Pointer(&fn)
	cbPtr := C.goNumPeersCallback

	C.abl_link_set_num_peers_callback(l.instance, (*[0]byte)(cbPtr), fnPtr)
}

// SetNumPeersCallback registers a callback to be notified when the session tempo changes.
//
// Thread-safe: yes
// Realtime-safe: no
//
// The callback is invoked on a Link-managed thread.
func (l *Link) SetTempoCallback(fn TempoCallbackFn) {
	fnPtr := unsafe.Pointer(&fn)
	cbPtr := C.goTempoCallback

	C.abl_link_set_tempo_callback(l.instance, (*[0]byte)(cbPtr), fnPtr)
}

// SetNumPeersCallback registers a callback to be notified when the state of start/stop isPlaying changes.
//
// Thread-safe: yes
// Realtime-safe: no
//
// The callback is invoked on a Link-managed thread.
func (l *Link) SetStartStopCallback(fn StartStopCallbackFn) {
	fnPtr := unsafe.Pointer(&fn)
	cbPtr := C.goStartStopCallback

	C.abl_link_set_start_stop_callback(l.instance, (*[0]byte)(cbPtr), fnPtr)
}

// Clock returns the current link clock time.
//
// Thread-safe: yes
// Realtime-safe: yes
func (l *Link) Clock() time.Duration {
	return time.Duration(C.abl_link_clock_micros(l.instance)) * time.Microsecond
}

// CaptureAudioSessionState captures the current Link Session State from the audio thread.
//
// Thread-safe: no
// Realtime-safe: yes
//
// This function should ONLY be called in the audio thread and must not be
// accessed from any other threads. After capturing the SessionState holds a snapshot
// of the current Link SessionState, so it should be used in a local scope. The
// SessionState should not be created on the audio thread.
func (l *Link) CaptureAudioSessionState(state *SessionState) {
	C.abl_link_capture_audio_session_state(l.instance, state.instance)
}

// CommitAudioSessionState commits the given Session State to the Link session from the audio thread.
//
// Thread-safe: no
// Realtime-safe: yes
//
// This function should ONLY be called in the audio thread. The given
// SessionState will replace the current Link state. Modifications will be
// communicated to other peers in the session.
func (l *Link) CommitAudioSessionState(state *SessionState) {
	C.abl_link_commit_audio_session_state(l.instance, state.instance)
}

// CaptureAppSessionState capturep the current Link Session State from an application thread.
//
// Thread-safe: no
// Realtime-safe: yes
//
// Provides a mechanism for capturing the Link Session State from an
// application thread (other than the audio thread). After capturing the SessionState
// contains a snapshot of the current Link state, so it should be used in a local
// scope.
func (l *Link) CaptureAppSessionState(state *SessionState) {
	C.abl_link_capture_app_session_state(l.instance, state.instance)
}

// CommitAppSessionState commits the given Session State to the Link session from an application thread.
//
// Thread-safe: yes
// Realtime-safe: no
//
// The given session_state will replace the current Link SessionState.
// Modifications of the SessionState will be communicated to other peers in the
// session.
func (l *Link) CommitAppSessionState(state *SessionState) {
	C.abl_link_commit_app_session_state(l.instance, state.instance)
}

//export goNumPeersCallback
func goNumPeersCallback(numPeers C.uint64_t, context *C.void) {
	fn := (*NumPeersCallbackFn)(unsafe.Pointer(context))
	(*fn)(uint64(numPeers))
}

//export goTempoCallback
func goTempoCallback(tempo C.double, context *C.void) {
	fn := (*TempoCallbackFn)(unsafe.Pointer(context))
	(*fn)(float64(tempo))
}

//export goStartStopCallback
func goStartStopCallback(isPlaying C.bool, context *C.void) {
	fn := (*StartStopCallbackFn)(unsafe.Pointer(context))
	(*fn)(bool(isPlaying))
}

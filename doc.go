// Package al implements Go bindings for Ableton Link.
//
// Each Link instance has its own session state which
// represents a beat timeline and a transport start/stop state. The
// timeline starts running from beat 0 at the initial tempo when
// constructed. The timeline always advances at a speed defined by
// its current tempo, even if transport is stopped. Synchronizing to the
// transport start/stop state of Link is optional for every peer.
// The transport start/stop state is only shared with other peers when
// start/stop synchronization is enabled.
//
// An Link instance is initially disabled after construction, which
// means that it will not communicate on the network. Once enabled,
// an Link instance initiates network communication in an effort to
// discover other peers. When peers are discovered, they immediately
// become part of a shared Link session.
//
// Each function documents its thread-safety and
// realtime-safety properties. When a function is marked thread-safe,
// it means it is safe to call from multiple threads
// concurrently. When a function is marked realtime-safe, it means that
// it does not block and is appropriate for use in the thread that
// performs audio IO.
//
// One session state capture/commit function pair for use
// in the audio thread and one for all other application contexts is provided.
// In general, modifying the session state should be done in the audio
// thread for the most accurate timing results. The ability to modify
// the session state from application threads should only be used in
// cases where an application's audio thread is not actively running
// or if it doesn't generate audio at all. Modifying the Link session
// state from both the audio thread and an application thread
// concurrently is not advised and will potentially lead to unexpected
// behavior.
package al

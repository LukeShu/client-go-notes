// This is a more complete/helpful version of the Event documentation in
// <kubernetes.git/staging/src/k8s.io/api/core/v1/types.go>.
package v1 // import "k8s.io/api/core/v1"

// Event is a report of an event somewhere in the cluster.
//
// There is a pending enhancement to replace `Event.v1.` with what is currently
// `Event.v1beta1.events.k8s.io`; see <https://github.com/kubernetes/enhancements/issues/383> and
// <https://git.k8s.io/community/contributors/design-proposals/instrumentation/events-redesign.md>.
//
// `k8s.io/client-go/tools/record` is for `Event.v1.`
//
// `k8s.io/client-go/tools/events` is for `Event.v1beta1.events.k8s.io`
type Event struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	// The object that this event is about.
	InvolvedObject ObjectReference `json:"involvedObject" protobuf:"bytes,2,opt,name=involvedObject"`

	// This should be a short, machine understandable string that gives the reason
	// for the transition into the object's current status.
	// TODO: provide exact specification for format.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`

	// A human-readable description of the status of this operation.
	// TODO: decide on maximum length.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`

	// The component reporting this event. Should be a short machine understandable string.
	// +optional
	Source EventSource `json:"source,omitempty" protobuf:"bytes,5,opt,name=source"`

	// The time at which the event was first recorded. (Time of server receipt is in TypeMeta.)
	// +optional
	FirstTimestamp metav1.Time `json:"firstTimestamp,omitempty" protobuf:"bytes,6,opt,name=firstTimestamp"`

	// The time at which the most recent occurrence of this event was recorded.
	// +optional
	LastTimestamp metav1.Time `json:"lastTimestamp,omitempty" protobuf:"bytes,7,opt,name=lastTimestamp"`

	// The number of times this event has occurred.
	// +optional
	Count int32 `json:"count,omitempty" protobuf:"varint,8,opt,name=count"`

	// Type of this event (Normal, Warning), new types could be added in the future
	// +optional
	Type string `json:"type,omitempty" protobuf:"bytes,9,opt,name=type"`

	// Time when this Event was first observed.
	//
	// This exists to enable future migration with `Event.v1beta1.events.k8s.io`; prefer to use
	// `.firstTimestamp` for native `Event.v1.`.
	//
	// +optional
	EventTime metav1.MicroTime `json:"eventTime,omitempty" protobuf:"bytes,10,opt,name=eventTime"`

	// Data about the Event series this event represents or nil if it's a singleton Event.
	//
	// This exists to enable future migration with `Event.v1beta1.events.k8s.io`; prefer to use
	// `.count` and `.lastTimestamp` for native `Event.v1.` (there is no native `Event.v1.`
	// equivalent of `.series.state`, but `.series.state` is deprecated and slated for removal
	// in 1.18).
	//
	// +optional
	Series *EventSeries `json:"series,omitempty" protobuf:"bytes,11,opt,name=series"`

	// What action was taken/failed regarding the `.involvedObject` object.
	//
	// This exists to enable future migration with `Event.v1beta1.events.k8s.io`; there is no
	// native `Event.v1.` equivalent.
	//
	// +optional
	Action string `json:"action,omitempty" protobuf:"bytes,12,opt,name=action"`

	// Optional secondary object for more complex actions.
	//
	// This exists to enable future migration with `Event.v1beta1.events.k8s.io`; there is no
	// native `Event.v1.` equivalent.
	//
	// +optional
	Related *ObjectReference `json:"related,omitempty" protobuf:"bytes,13,opt,name=related"`

	// Name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`.
	//
	// This exists to enable future migration with `Event.v1beta1.events.k8s.io`; prefer to use
	// `.source.component` for native `Event.v1.`.
	//
	// +optional
	ReportingController string `json:"reportingComponent" protobuf:"bytes,14,opt,name=reportingComponent"`

	// ID of the controller instance, e.g. `kubelet-xyzf`.
	//
	// This exists to enable future migration with `Event.v1beta1.events.k8s.io`; prefer to use
	// `.source.host` for native `Event.v1.`.
	//
	// +optional
	ReportingInstance string `json:"reportingInstance" protobuf:"bytes,15,opt,name=reportingInstance"`
}

// EventSeries contain information on series of events, i.e. thing that was/is happening
// continuously for some time.
//
// This exists to enable future migration with `EventSeries.v1beta1.events.k8s.io`; prefer to not
// use it.
type EventSeries struct {
	// Number of occurrences in this series up to the last heartbeat time
	Count int32 `json:"count,omitempty" protobuf:"varint,1,name=count"`
	// Time of the last occurrence observed
	LastObservedTime metav1.MicroTime `json:"lastObservedTime,omitempty" protobuf:"bytes,2,name=lastObservedTime"`
	// State of this Series: Ongoing or Finished
	// Deprecated. Planned removal for 1.18
	State EventSeriesState `json:"state,omitempty" protobuf:"bytes,3,name=state"`
}

type EventSeriesState string

const (
	EventSeriesStateOngoing  EventSeriesState = "Ongoing"
	EventSeriesStateFinished EventSeriesState = "Finished"
	EventSeriesStateUnknown  EventSeriesState = "Unknown"
)

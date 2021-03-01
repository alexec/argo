package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type WorkflowTaskSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Spec              WorkflowTaskSetSpec    `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	Status            *WorkflowTaskSetStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

func (in *WorkflowTaskSet) GetTemplateByName(name string) *Template {
	for _, t := range in.Spec.Templates {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

type WorkflowTaskSetSpec struct {
	// +patchStrategy=merge
	// +patchMergeKey=name
	Templates []Template `json:"templates,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,1,opt,name=templates"`
	Nodes     Nodes      `json:"nodes,omitempty" protobuf:"bytes,2,rep,name=nodes,casttype=Nodes"`
}

type WorkflowTaskSetStatus struct {
	Nodes map[string]NodeResult `json:"nodes,omitempty" protobuf:"bytes,1,rep,name=nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type WorkflowTaskSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []WorkflowTaskSet `json:"items" protobuf:"bytes,2,opt,name=items"`
}

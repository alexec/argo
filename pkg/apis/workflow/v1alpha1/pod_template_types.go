package v1alpha1

import corev1 "k8s.io/api/core/v1"

type PodTemplate struct {
	Graph        Graph                `json:"graph,omitempty" protobuf:"bytes,1,rep,name=graph"`
	Sequence     []Task               `json:"sequence,omitempty" protobuf:"bytes,2,rep,name=sequence"`
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty" protobuf:"bytes,3,rep,name=volumeMounts"`
}

func (in *PodTemplate) GetContainers() []corev1.Container {
	var ctrs []corev1.Container
	for _, t := range in.GetGraph() {
		c := t.Container
		c.VolumeMounts = append(c.VolumeMounts, in.VolumeMounts...)
		ctrs = append(ctrs, c)
	}
	return ctrs
}

func (in *PodTemplate) HasContainerNamed(n string) bool {
	for _, c := range in.GetContainers() {
		if n == c.Name {
			return true
		}
	}
	return false
}

func (in *PodTemplate) GetGraph() Graph {
	if in == nil {
		return nil
	}
	out := in.Graph
	for i, t := range in.Sequence {
		var d []string
		if i > 0 {
			d = append(d, in.Sequence[i-1].Name)
		}
		out = append(out, Node{Task: t, Dependencies: d})
	}
	return out
}

type Graph []Node

type Node struct {
	Task         `json:",inline" protobuf:"bytes,1,opt,name=container"`
	Dependencies []string `json:"dependencies,omitempty" protobuf:"bytes,2,rep,name=dependencies"`
}

type Task struct {
	corev1.Container `json:",inline" protobuf:"bytes,1,opt,name=container"`
}

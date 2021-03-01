// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// WorkflowTaskLister helps list WorkflowTasks.
// All objects returned here must be treated as read-only.
type WorkflowTaskLister interface {
	// List lists all WorkflowTasks in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.WorkflowTask, err error)
	// WorkflowTasks returns an object that can list and get WorkflowTasks.
	WorkflowTasks(namespace string) WorkflowTaskNamespaceLister
	WorkflowTaskListerExpansion
}

// workflowTaskLister implements the WorkflowTaskLister interface.
type workflowTaskLister struct {
	indexer cache.Indexer
}

// NewWorkflowTaskLister returns a new WorkflowTaskLister.
func NewWorkflowTaskLister(indexer cache.Indexer) WorkflowTaskLister {
	return &workflowTaskLister{indexer: indexer}
}

// List lists all WorkflowTasks in the indexer.
func (s *workflowTaskLister) List(selector labels.Selector) (ret []*v1alpha1.WorkflowTask, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.WorkflowTask))
	})
	return ret, err
}

// WorkflowTasks returns an object that can list and get WorkflowTasks.
func (s *workflowTaskLister) WorkflowTasks(namespace string) WorkflowTaskNamespaceLister {
	return workflowTaskNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// WorkflowTaskNamespaceLister helps list and get WorkflowTasks.
// All objects returned here must be treated as read-only.
type WorkflowTaskNamespaceLister interface {
	// List lists all WorkflowTasks in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.WorkflowTask, err error)
	// Get retrieves the WorkflowTask from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.WorkflowTask, error)
	WorkflowTaskNamespaceListerExpansion
}

// workflowTaskNamespaceLister implements the WorkflowTaskNamespaceLister
// interface.
type workflowTaskNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all WorkflowTasks in the indexer for a given namespace.
func (s workflowTaskNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.WorkflowTask, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.WorkflowTask))
	})
	return ret, err
}

// Get retrieves the WorkflowTask from the indexer for a given namespace and name.
func (s workflowTaskNamespaceLister) Get(name string) (*v1alpha1.WorkflowTask, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("workflowtask"), name)
	}
	return obj.(*v1alpha1.WorkflowTask), nil
}

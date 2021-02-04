package syncer

import (
	"context"
	"fmt"

	"github.com/iancoleman/strcase"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Sync mutates the subject of the syncer interface via controller-runtime
// CreateOrUpdate method.
func Sync(ctx context.Context, syncer Interface, recorder record.EventRecorder) error {
	result, err := syncer.Sync(ctx)

	owner := syncer.ObjectOwner()
	if recorder != nil && owner != nil && result.EventType != "" && result.EventReason != "" && result.EventMessage != "" {
		if err != nil || result.Operation != controllerutil.OperationResultNone {
			recorder.Eventf(owner, result.EventType, result.EventReason, result.EventMessage)
		}
	}

	return err
}

// New return a syncer.Sync object
func New(name string, c client.Client, owner client.Object, obj client.Object, scheme *runtime.Scheme, fn controllerutil.MutateFn) Interface {
	return &objectSyncer{
		Name:   name,
		Owner:  owner,
		Self:   obj,
		SyncFn: fn,
		Client: c,
		Scheme: scheme}
}

func getKey(obj client.Object) (types.NamespacedName, error) {
	key := types.NamespacedName{}
	objMeta, ok := obj.(metav1.Object)
	if !ok {
		return key, fmt.Errorf("%T is not a metav1.Object", obj)
	}

	key.Name = objMeta.GetName()
	key.Namespace = objMeta.GetNamespace()

	return key, nil
}

func basicEventReason(objKindName string, err error) string {
	if err != nil {
		return fmt.Sprintf("%sSyncFailed", strcase.ToCamel(objKindName))
	}

	return fmt.Sprintf("%sSyncSuccessfull", strcase.ToCamel(objKindName))
}

// stripSecrets returns a copy for the secret without secret data in it
func stripSecrets(obj runtime.Object) runtime.Object {
	// if obj is secret, don't print secret data
	s, ok := obj.(*corev1.Secret)
	if ok {
		cObj := s.DeepCopyObject().(*corev1.Secret)
		cObj.Data = nil
		cObj.StringData = nil
		return cObj
	}

	return obj
}
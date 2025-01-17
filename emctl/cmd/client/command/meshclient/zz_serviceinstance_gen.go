/*
Copyright (c) 2021, MegaEase
All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// code generated by github.com/megaease/easemeshctl/cmd/generator, DO NOT EDIT.
package meshclient

import (
	"context"
	"encoding/json"
	"fmt"
	v1alpha1 "github.com/megaease/easemesh-api/v1alpha1"
	resource "github.com/megaease/easemeshctl/cmd/client/resource"
	client "github.com/megaease/easemeshctl/cmd/common/client"
	errors "github.com/pkg/errors"
	"net/http"
)

type serviceInstanceInterface struct {
	client *meshClient
}
type serviceInstanceGetter struct {
	client *meshClient
}

func (s *serviceInstanceGetter) ServiceInstance() ServiceInstanceInterface {
	return &serviceInstanceInterface{client: s.client}
}
func (s *serviceInstanceInterface) Get(args0 context.Context, args1 string, args2 string) (*resource.ServiceInstance, error) {
	url := fmt.Sprintf("http://"+s.client.server+apiURL+"/mesh/"+"serviceinstances/%s/%s", args1, args2)
	r0, err := client.NewHTTPJSON().GetByContext(args0, url, nil, nil).HandleResponse(func(buff []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "get ServiceInstance %s", args1)
		}
		if statusCode >= 300 {
			return nil, errors.Errorf("call %s failed, return status code %d text %+v", url, statusCode, string(buff))
		}
		ServiceInstance := &v1alpha1.ServiceInstance{}
		err := json.Unmarshal(buff, ServiceInstance)
		if err != nil {
			return nil, errors.Wrapf(err, "unmarshal data to v1alpha1.ServiceInstance")
		}
		return resource.ToServiceInstance(ServiceInstance), nil
	})
	if err != nil {
		return nil, err
	}
	return r0.(*resource.ServiceInstance), nil
}
func (s *serviceInstanceInterface) Delete(args0 context.Context, args1 string, args2 string) error {
	url := fmt.Sprintf("http://"+s.client.server+apiURL+"/mesh/"+"serviceinstances/%s/%s", args1, args2)
	_, err := client.NewHTTPJSON().DeleteByContext(args0, url, nil, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "Delete ServiceInstance %s", args1)
		}
		if statusCode < 300 && statusCode >= 200 {
			return nil, nil
		}
		return nil, errors.Errorf("call Delete %s failed, return statuscode %d text %+v", url, statusCode, string(b))
	})
	return err
}
func (s *serviceInstanceInterface) List(args0 context.Context) ([]*resource.ServiceInstance, error) {
	url := "http://" + s.client.server + apiURL + "/mesh/serviceinstances"
	result, err := client.NewHTTPJSON().GetByContext(args0, url, nil, nil).HandleResponse(func(b []byte, statusCode int) (interface{}, error) {
		if statusCode == http.StatusNotFound {
			return nil, errors.Wrapf(NotFoundError, "list service")
		}
		if statusCode >= 300 && statusCode < 200 {
			return nil, errors.Errorf("call GET %s failed, return statuscode %d text %+v", url, statusCode, b)
		}
		serviceInstance := []v1alpha1.ServiceInstance{}
		err := json.Unmarshal(b, &serviceInstance)
		if err != nil {
			return nil, errors.Wrapf(err, "unmarshal data to v1alpha1.")
		}
		results := []*resource.ServiceInstance{}
		for _, item := range serviceInstance {
			copy := item
			results = append(results, resource.ToServiceInstance(&copy))
		}
		return results, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*resource.ServiceInstance), nil
}

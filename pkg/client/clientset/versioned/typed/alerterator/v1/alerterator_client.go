// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/nais/alerterator/pkg/apis/alerterator/v1"
	"github.com/nais/alerterator/pkg/client/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type AlerteratorV1Interface interface {
	RESTClient() rest.Interface
	AlertsGetter
}

// AlerteratorV1Client is used to interact with features provided by the alerterator.nais.io group.
type AlerteratorV1Client struct {
	restClient rest.Interface
}

func (c *AlerteratorV1Client) Alerts(namespace string) AlertInterface {
	return newAlerts(c, namespace)
}

// NewForConfig creates a new AlerteratorV1Client for the given config.
func NewForConfig(c *rest.Config) (*AlerteratorV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &AlerteratorV1Client{client}, nil
}

// NewForConfigOrDie creates a new AlerteratorV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *AlerteratorV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new AlerteratorV1Client for the given RESTClient.
func New(c rest.Interface) *AlerteratorV1Client {
	return &AlerteratorV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *AlerteratorV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
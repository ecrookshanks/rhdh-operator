package model

import (
	"fmt"
	"os"
	"path/filepath"

	openshift "github.com/openshift/api/route/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	"k8s.io/utils/ptr"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/apimachinery/pkg/runtime"

	bsv1 "github.com/redhat-developer/rhdh-operator/api/v1alpha4"
)

// testBackstageObject it is a helper object to simplify testing model component allowing to customize and isolate testing configuration
// usual sequence of creating testBackstageObject contains such a steps:
// createBackstageTest(bsv1.Backstage).
// withDefaultConfig(useDef bool)
// addToDefaultConfig(key, fileName)
type testBackstageObject struct {
	backstage      bsv1.Backstage
	externalConfig ExternalConfig
	scheme         *runtime.Scheme
}

// initialises testBackstageObject object
func createBackstageTest(bs bsv1.Backstage) *testBackstageObject {
	ec := ExternalConfig{
		RawConfig: map[string]string{},
		//AppConfigs:          map[string]corev1.ConfigMap{},
		//ExtraFileConfigMaps: map[string]corev1.ConfigMap{},
		//ExtraEnvConfigMaps: map[string]corev1.ConfigMap{},
	}
	b := &testBackstageObject{backstage: bs, externalConfig: ec, scheme: runtime.NewScheme()}
	utilruntime.Must(bsv1.AddToScheme(b.scheme))
	utilruntime.Must(clientgoscheme.AddToScheme(b.scheme))
	utilruntime.Must(openshift.Install(b.scheme))
	return b
}

// enables LocalDB
func (b *testBackstageObject) withLocalDb() *testBackstageObject {
	b.backstage.Spec.Database.EnableLocalDb = ptr.To(true)
	return b
}

// tells if object should use default Backstage Deployment/Service configuration from ./testdata/default-config or not
func (b *testBackstageObject) withDefaultConfig(useDef bool) *testBackstageObject {
	if useDef {
		// here we have default-config folder
		_ = os.Setenv("LOCALBIN", "./testdata")
	} else {
		_ = os.Setenv("LOCALBIN", ".")
	}
	return b
}

// adds particular part of configuration pointing to configuration key
// where key is configuration key (such as "deployment.yaml" and fileName is a name of additional conf file in ./testdata
func (b *testBackstageObject) addToDefaultConfig(key string, fileName string) *testBackstageObject {

	yaml, err := readTestYamlFile(fileName)
	if err != nil {
		panic(err)
	}

	b.externalConfig.RawConfig[key] = string(yaml)

	return b
}

// reads file from ./testdata
func readTestYamlFile(name string) ([]byte, error) {

	b, err := os.ReadFile(filepath.Join("testdata", name)) // #nosec G304, path is constructed internally
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}
	return b, nil
}

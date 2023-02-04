package objects

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"
)

type ReleaseNames map[string]bool
type NsReleases map[string]ReleaseNames
type ContextNsReleases map[string]NsReleases

type NSHelmReleaseFilter = func(ns string, rel string) bool

func ReleaseFilterFactory(ctx string) (NSHelmReleaseFilter, error) {
	filterFile := os.Getenv("RELEASES_FILTER_FILE")
	defaultFilter := func(ns string, rel string) bool {
		return true
	}
	if filterFile == "" {
		return defaultFilter, nil
	}
	yamlFile, err := ioutil.ReadFile(filterFile)
	if err != nil {
		return nil, err
	}
	var cNRs ContextNsReleases
	err = yaml.Unmarshal(yamlFile, &cNRs)
	if err != nil {
		return nil, err
	}

	nRs, present := cNRs[ctx]
	if !present {
		return defaultFilter, nil
	}

	filterfunc := func(ns string, rel string) bool {
		log.Debugf("in Filterfunc %s %s", ns, nRs)
		rs, pres := nRs[ns]
		if !pres {
			log.Debugf("in Filterfunc not present1 %s %s", ns, nRs)
			return true
		}
		return rs[rel]
	}
	log.Debugf("ctx release filter found for cluster %s and %s", ctx, yamlFile)
	log.Debugf("ctx release filter MAP found for cluster %s and %s", ctx, nRs)

	return filterfunc, nil

}

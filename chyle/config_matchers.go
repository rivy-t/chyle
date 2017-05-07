package chyle

import (
	"github.com/antham/envh"
)

// matchersConfigurator validates jira config
// defined through environment variables
type matchersConfigurator struct {
	chyleConfig *CHYLE
	config      *envh.EnvTree
	definedKeys []string
}

func (m *matchersConfigurator) process() (bool, error) {
	if m.isDisabled() {
		return true, nil
	}

	for _, f := range []func() error{
		m.validateRegexpMatchers,
		m.validateTypeMatcher,
	} {
		if err := f(); err != nil {
			return true, err
		}
	}

	m.setMatchers()

	return true, nil
}

// isDisabled checks if matchers are enabled
func (m *matchersConfigurator) isDisabled() bool {
	return featureDisabled(m.config, [][]string{{"CHYLE", "MATCHERS"}})
}

// validateRegexpMatchers checks all config relying on valid regexp
func (m *matchersConfigurator) validateRegexpMatchers() error {
	for _, key := range []string{"MESSAGE", "COMMITTER", "AUTHOR"} {
		_, err := m.config.FindString("CHYLE", "MATCHERS", key)

		if err != nil {
			continue
		}

		if err := validateRegexp(m.config, []string{"CHYLE", "MATCHERS", key}); err != nil {
			return err
		}

		m.definedKeys = append(m.definedKeys, key)
	}

	return nil
}

// validateTypeMatcher checks custom field TYPE
func (m *matchersConfigurator) validateTypeMatcher() error {
	_, err := m.config.FindString("CHYLE", "MATCHERS", "TYPE")

	if err != nil {
		return nil
	}

	if err := validateOneOf(m.config, []string{"CHYLE", "MATCHERS", "TYPE"}, []string{regularTypeMatcher, mergeTypeMatcher}); err != nil {
		return err
	}

	m.definedKeys = append(m.definedKeys, "TYPE")

	return nil
}

// setMatchers update chyleConfig with extracted matchers
func (m *matchersConfigurator) setMatchers() {
	m.chyleConfig.MATCHERS = map[string]string{}

	for _, key := range m.definedKeys {
		m.chyleConfig.MATCHERS[key] = m.config.FindStringUnsecured("CHYLE", "MATCHERS", key)
	}
}

package config

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/tools"
)

// Test config files
const (
	specifiedConfig     string = "%s/config/test_configs/specified_config"
	badConfig           string = "%s/config/test_configs/bad_config"
	noPrivateKeysConfig string = "%s/config/test_configs/no_private_keys"
)

var _ = Describe("Config unit tests @unit", func() {

	Describe("Verify order of importance for environment variables and config files", func() {
		It("should load the default config file", func() {
			conf, err := NewConfig(LocalConfig, tools.ProjectRoot)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.KeepEnvironments).Should(Equal("Never"), "We either changed the default value in the config or it did not load correctly")
		})

		It("should load a specified file", func() {
			conf, err := NewConfig(LocalConfig, fmt.Sprintf(specifiedConfig, tools.ProjectRoot))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.KeepEnvironments).Should(Equal("Always"), "We did not load the correct config file")
		})

		It("should fail to load a bad file", func() {
			_, err := NewConfig(LocalConfig, fmt.Sprintf(badConfig, tools.ProjectRoot))
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("line 1: cannot unmarshal"))
		})

		It("should overwrite default values with ENV variables", func() {
			err := os.Setenv("KEEP_ENVIRONMENTS", "OnFail")
			Expect(err).ShouldNot(HaveOccurred())
			conf, err := NewConfig(LocalConfig, tools.ProjectRoot)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.KeepEnvironments).Should(Equal("OnFail"), "The environment variable should have been used to change the config value")
		})

		It("should overwrite specified file values with ENV variables", func() {
			err := os.Setenv("KEEP_ENVIRONMENTS", "OnFail")
			Expect(err).ShouldNot(HaveOccurred())
			conf, err := NewConfig(LocalConfig, fmt.Sprintf(specifiedConfig, tools.ProjectRoot))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.KeepEnvironments).Should(Equal("OnFail"), "The environment variable should have been used to change the config value")

		})

		It("should load the config if we specify a Secret Config type", func() {
			conf, err := NewConfig(SecretConfig, tools.ProjectRoot)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.KeepEnvironments).Should(Equal("Never"), "We either changed the default value in the config or it did not load correctly")
		})

		AfterEach(func() {
			err := os.Unsetenv("KEEP_ENVIRONMENTS")
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	It("should get the network config with a valid name", func() {
		conf, err := NewConfig(LocalConfig, fmt.Sprintf(specifiedConfig, tools.ProjectRoot))
		Expect(err).ShouldNot(HaveOccurred())
		network, err := conf.GetNetworkConfig("test_this_geth")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(network.Name).Should(Equal("Tester Ted"), "The network config was not loaded correctly")
	})

	It("should not get the network config with an invalid name", func() {
		conf, err := NewConfig(LocalConfig, fmt.Sprintf(specifiedConfig, tools.ProjectRoot))
		Expect(err).ShouldNot(HaveOccurred())
		_, err = conf.GetNetworkConfig("bad_name")
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring("no supported network"))
	})

	It("should fetch private keys when they exist", func() {
		conf, err := NewConfig(LocalConfig, fmt.Sprintf(specifiedConfig, tools.ProjectRoot))
		Expect(err).ShouldNot(HaveOccurred())
		network, err := conf.GetNetworkConfig("test_this_geth")
		Expect(err).ShouldNot(HaveOccurred())
		privateKeys := NewPrivateKeyStore(LocalConfig, network)
		keys, err := privateKeys.Fetch()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(len(keys)).Should(Equal(1), "The number of private keys was incorrect")
		Expect(keys[0]).Should(Equal("abcdefg"), "The private key did not get read correctly")
	})

	It("should not fetch private keys when they do not exist", func() {
		conf, err := NewConfig(LocalConfig, fmt.Sprintf(noPrivateKeysConfig, tools.ProjectRoot))
		Expect(err).ShouldNot(HaveOccurred())
		network, err := conf.GetNetworkConfig("test_this_geth")
		Expect(err).ShouldNot(HaveOccurred())
		privateKeys := NewPrivateKeyStore(LocalConfig, network)
		_, err = privateKeys.Fetch()
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring("no keys found"))
	})

	It("should fetch secret store, code still has a todo so fix this test when the code is updated", func() {
		conf, err := NewConfig(SecretConfig, fmt.Sprintf(specifiedConfig, tools.ProjectRoot))
		Expect(err).ShouldNot(HaveOccurred())
		network, err := conf.GetNetworkConfig("test_this_geth")
		Expect(err).ShouldNot(HaveOccurred())
		privateKeys := NewPrivateKeyStore(SecretConfig, network)
		keys, err := privateKeys.Fetch()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(len(keys)).Should(Equal(1), "The number of private keys was incorrect")
		Expect(keys[0]).Should(Equal(""), "The private key did not get read correctly")
	})

	It("should not fetch private keys when asking for a bad configuration type", func() {
		conf, err := NewConfig("badConfigurationType", fmt.Sprintf(noPrivateKeysConfig, tools.ProjectRoot))
		Expect(err).ShouldNot(HaveOccurred())
		network, err := conf.GetNetworkConfig("test_this_geth")
		Expect(err).ShouldNot(HaveOccurred())
		privateKeys := NewPrivateKeyStore("badConfigurationType", network)
		Expect(privateKeys).Should(BeNil(), "We should have no keystore for invalid configuration types")
	})

})

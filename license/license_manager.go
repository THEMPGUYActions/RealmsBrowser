package license

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/kitecyber/inputbox"
)

type Manager struct {
	client    *LicenseClient
	features  *FeatureManager
	cacheDir  string
	stop      chan struct{}
	onInvalid func(error)
	once      sync.Once
}

func NewManager(dataPath string, onInvalid func(error)) *Manager {
	return &Manager{
		client:    &LicenseClient{},
		features:  &FeatureManager{},
		cacheDir:  dataPath,
		onInvalid: onInvalid,
		stop:      make(chan struct{}),
	}
}

func (m *Manager) Features() *FeatureManager {
	return m.features
}

func (m *Manager) EnsureLicensed() error {
	cache, err := LoadCache(m.cacheDir)
	if err != nil {
		cache = &LicenseCache{}
	}

	key := strings.TrimSpace(cache.LicenseKey)
	if key == "" {
		key, err = requestLicenseKey()
		if err != nil {
			return err
		}
	}

	response, err := m.client.Validate(key)
	if err != nil {
		return err
	}

	cache.LicenseKey = key
	cache.Owner = response.Owner
	if err := SaveCache(m.cacheDir, cache); err != nil {
		return fmt.Errorf("license validated but could not save license cache: %w", err)
	}
	m.features.Enable(response.Owner)
	return nil
}

func (m *Manager) StartMonitoring(key string) {
	m.once.Do(func() {
		go func() {
			ticker := time.NewTicker(12 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if _, err := m.client.Validate(key); err != nil {
						m.features.Disable()
						if m.onInvalid != nil {
							m.onInvalid(err)
						}
						return
					}
				case <-m.stop:
					return
				}
			}
		}()
	})
}

func (m *Manager) StopMonitoring() {
	select {
	case <-m.stop:
	default:
		close(m.stop)
	}
}

func requestLicenseKey() (string, error) {
	message := "RealmsBrowser requires a valid license key to continue.\n\n" +
		"Your license key verifies ownership of this copy of RealmsBrowser and enables access to licensed features.\n\n" +
		"License Key:\n\n" +
		"Need a license? Add \"thempguy.\" on Discord!"
	key, ok := inputbox.InputBox("RealmsBrowser License Activation", message, "")
	if !ok {
		return "", fmt.Errorf("license activation was cancelled")
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return "", fmt.Errorf("no license key was entered")
	}
	return key, nil
}

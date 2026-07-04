import { useState, useEffect } from 'react';
import './Settings.css';
import { GetPreferences, SavePreferences, ListAllBrowsers, ClosePreferences } from '../../bindings/browski/browserservice';

export function Settings({ onClose }) {
    const [browsers, setBrowsers] = useState([]);
    const [preferences, setPreferences] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const loadData = async () => {
            try {
                const [browsersData, prefsData] = await Promise.all([
                    ListAllBrowsers(),
                    GetPreferences()
                ]);
                setBrowsers(browsersData);
                setPreferences(prefsData);
                setLoading(false);
            } catch (error) {
                console.error('Failed to load settings:', error);
            }
        };

        loadData();

        // Listen for window show event to reload data
        const handleVisibilityChange = () => {
            if (!document.hidden) {
                loadData();
            }
        };

        document.addEventListener('visibilitychange', handleVisibilityChange);
        return () => document.removeEventListener('visibilitychange', handleVisibilityChange);
    }, []);

    const getBrowserKey = (browserName, profileName) => {
        return profileName ? `${browserName}:${profileName}` : browserName;
    };

    const isBrowserHidden = (browserName, profileName) => {
        if (!preferences || !preferences.hiddenBrowsers) return false;
        const key = getBrowserKey(browserName, profileName);
        return preferences.hiddenBrowsers[key] === true;
    };

    const toggleBrowserVisibility = (browserName, profileName) => {
        const key = getBrowserKey(browserName, profileName);
        const newPrefs = { ...preferences };

        if (!newPrefs.hiddenBrowsers) {
            newPrefs.hiddenBrowsers = {};
        }

        if (newPrefs.hiddenBrowsers[key]) {
            delete newPrefs.hiddenBrowsers[key];
        } else {
            newPrefs.hiddenBrowsers[key] = true;
        }

        setPreferences(newPrefs);
    };

    const setDefaultBrowser = (browserName, profileName) => {
        const key = getBrowserKey(browserName, profileName);
        setPreferences({
            ...preferences,
            defaultBrowser: key
        });
    };

    const handleSave = async () => {
        try {
            await SavePreferences(preferences);
        } catch (error) {
            console.error('Failed to save preferences:', error);
            alert('Failed to save preferences');
        }
    };

    if (loading) {
        return (
            <div className="settings-window">
                <div className="settings-loading">Loading...</div>
            </div>
        );
    }

    // Get all browsers and profiles
    const allItems = [];

    browsers.forEach(browser => {
        const browserItem = {
            type: 'browser',
            name: browser.name,
            icon: browser.icon,
            browserName: browser.name,
            profileName: ''
        };

        allItems.push(browserItem);

        if (browser.profiles && browser.profiles.length > 0) {
            browser.profiles.forEach(profile => {
                const profileItem = {
                    type: 'profile',
                    name: profile.name,
                    icon: profile.icon,
                    browserName: browser.name,
                    profileName: profile.name
                };

                allItems.push(profileItem);
            });
        }
    });

    const handleRemoveItem = (item) => {
        toggleBrowserVisibility(item.browserName, item.profileName);
        setTimeout(() => handleSave(), 100);
    };

    return (
        <div className="settings-window">
            <div className="settings-header">
                <h2>Browser Visibility</h2>
                <p>Click to show/hide browsers in the selector</p>
            </div>
            <div className="browser-list">
                {allItems.map((item, index) => {
                    const isHidden = isBrowserHidden(item.browserName, item.profileName);
                    return (
                        <div
                            key={index}
                            className={`browser-item ${isHidden ? 'hidden' : 'visible'}`}
                            onClick={() => handleRemoveItem(item)}
                        >
                            <img src={item.icon} alt={item.name} className="browser-icon" />
                            <span className="browser-name">{item.name}</span>
                            <span className="visibility-badge">
                                {isHidden ? 'Hidden' : 'Visible'}
                            </span>
                        </div>
                    );
                })}
            </div>
        </div>
    );
}

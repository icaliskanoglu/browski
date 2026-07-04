import { BrowserCard } from './BrowserCard';
import './BrowserGrid.css';

export function BrowserGrid({ browsers, onSelect }) {
    const items = [];
    let index = 0;

    browsers.forEach(browser => {
        // Add the main browser only if showMainEntry is true
        if (browser.showMainEntry !== false) {
            items.push({
                browser,
                profile: null,
                index: index++,
            });
        }

        // Add profiles if they exist
        if (browser.profiles && browser.profiles.length > 0) {
            browser.profiles.forEach(profile => {
                items.push({
                    browser,
                    profile,
                    index: index++,
                });
            });
        }
    });

    return (
        <div className="browser-grid">
            {items.map((item, idx) => (
                <BrowserCard
                    key={`${item.browser.name}-${item.profile?.name || 'default'}-${idx}`}
                    browser={item.browser}
                    profile={item.profile}
                    index={item.index}
                    onSelect={onSelect}
                />
            ))}
        </div>
    );
}

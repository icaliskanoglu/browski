import { useState } from 'react';
import './BrowserCard.css';

export function BrowserCard({ browser, profile, index, onSelect }) {
    const [isHovered, setIsHovered] = useState(false);

    const displayName = profile ? profile.name : browser.name;
    const icon = profile ? profile.icon : browser.icon;
    const showKeyboardShortcut = index < 9;

    const handleClick = () => {
        onSelect(browser.type, browser.name, profile?.name);
    };

    return (
        <button
            className={`browser-card ${isHovered ? 'hovered' : ''}`}
            onClick={handleClick}
            onMouseEnter={() => setIsHovered(true)}
            onMouseLeave={() => setIsHovered(false)}
        >
            <div className="browser-card-icon-wrapper">
                <img
                    src={icon}
                    alt={displayName}
                    className="browser-card-icon"
                    loading="eager"
                    decoding="async"
                />
                {showKeyboardShortcut && (
                    <div className="browser-card-shortcut">{index + 1}</div>
                )}
            </div>
            <div className="browser-card-name">{displayName}</div>
        </button>
    );
}

import { useEffect, useState, useRef } from 'react';
import './Selector.css';
import { ListBrowsers, Open, HideWindow, ShowWindow, Resize, GetURL } from '../../bindings/browski/browserservice';
import { Events } from "@wailsio/runtime";
import { BrowserGrid } from './BrowserGrid';

export function Selector() {
    const [browsers, setBrowsers] = useState([]);
    const [url, setUrl] = useState('');
    const [isLoading, setIsLoading] = useState(true);
    const componentRef = useRef(null);
    const browserItemsRef = useRef([]);
    const isLoadedRef = useRef(false);

    useEffect(() => {
        const initialize = async () => {
            try {
                const [browsersData, urlData] = await Promise.all([
                    ListBrowsers(),
                    GetURL()
                ]);

                // Handle null or undefined browsers
                const validBrowsers = browsersData || [];
                setBrowsers(validBrowsers);
                setUrl(urlData);

                // Build items list for keyboard navigation
                const items = [];
                validBrowsers.forEach(browser => {
                    items.push({ browser, profile: null });
                    if (browser.profiles && browser.profiles.length > 0) {
                        browser.profiles.forEach(profile => {
                            items.push({ browser, profile });
                        });
                    }
                });
                browserItemsRef.current = items;

                // Preload images for faster rendering
                const images = [];
                browsersData.forEach(browser => {
                    images.push(browser.icon);
                    if (browser.profiles) {
                        browser.profiles.forEach(profile => {
                            images.push(profile.icon);
                        });
                    }
                });

                // Preload all images
                const imagePromises = images.map(src => {
                    return new Promise((resolve) => {
                        const img = new Image();
                        img.onload = resolve;
                        img.onerror = resolve; // Still resolve even on error
                        img.src = src;
                    });
                });

                await Promise.all(imagePromises);
                setIsLoading(false);
                isLoadedRef.current = true;

                // Calculate window dimensions based on content
                const browserCardWidth = 80; // min-width of each card
                const browserCardHeight = 100; // Card content height
                const gap = 12; // gap between cards
                const gridPadding = 32; // 16px top + 16px bottom
                const footerHeight = 45; // footer with URL (includes padding)

                const totalBrowsers = items.length;
                const calculatedWidth = (totalBrowsers * browserCardWidth) + ((totalBrowsers - 1) * gap) + gridPadding;
                const windowWidth = Math.max(300, calculatedWidth); // Minimum 300px
                const windowHeight = browserCardHeight + gridPadding + footerHeight; // Extra 6px buffer

                // Resize and show the window
                Resize(windowWidth, windowHeight);
                ShowWindow();
            } catch (error) {
                console.error('Failed to initialize:', error);
                setIsLoading(false);
                isLoadedRef.current = true;
                // Still show the window even if there was an error
                ShowWindow();
            }
        };

        initialize();

        // Listen for preferences changes
        const handleVisibilityChange = () => {
            if (!document.hidden) {
                // Reload browsers when window becomes visible
                ListBrowsers().then(browsersData => {
                    setBrowsers(browsersData);
                });
            }
        };

        document.addEventListener('visibilitychange', handleVisibilityChange);

        return () => {
            document.removeEventListener('visibilitychange', handleVisibilityChange);
        };
    }, []);

    // Separate effect for URL change event listener
    useEffect(() => {
        // Listen for URL changes (from IPC when another instance sends a URL)
        const unsubscribe = Events.On('url-changed', (event) => {
            // Wails events are wrapped in {name, data} format
            const newUrl = event.data || event;
            setUrl(newUrl);
            // If we're already loaded and receive a URL change (from IPC),
            // show the window
            if (isLoadedRef.current) {
                ShowWindow();
            }
        });

        return () => {
            if (unsubscribe) unsubscribe();
        };
    }, []);

    // Disable dynamic resizing - use fixed window size instead
    // useEffect(() => {
    //     if (!componentRef.current || isLoading) return;

    //     const resizeWindow = () => {
    //         const width = componentRef.current.offsetWidth;
    //         const height = componentRef.current.offsetHeight;
    //         Resize(width, height);
    //     };

    //     // Use requestAnimationFrame to ensure DOM is painted before resize
    //     requestAnimationFrame(() => {
    //         requestAnimationFrame(() => {
    //             resizeWindow();
    //         });
    //     });
    // }, [browsers, isLoading]);

    useEffect(() => {
        const handleKeyDown = (event) => {
            if (event.key === 'Escape') {
                HideWindow();
                return;
            }

            const num = parseInt(event.key);
            if (num >= 1 && num <= 9) {
                const index = num - 1;
                if (index < browserItemsRef.current.length) {
                    const item = browserItemsRef.current[index];
                    handleSelect(item.browser.type, item.browser.name, item.profile?.name);
                }
            }
        };

        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, []);

    const handleSelect = (type, browser, profile) => {
        Open({
            type,
            browser,
            profile: profile || ''
        });
    };

    if (isLoading) {
        return <div style={{color: 'white', padding: '20px'}}>Loading browsers...</div>;
    }

    if (!browsers || browsers.length === 0) {
        return (
            <div style={{color: 'white', padding: '20px'}}>
                <p>No browsers detected</p>
                <p style={{fontSize: '12px', marginTop: '10px'}}>
                    Please check if you have browsers installed
                </p>
            </div>
        );
    }

    return (
        <div className="selector" ref={componentRef}>
            <div className="selector-content">
                <div className="browser-grid-wrapper">
                    <BrowserGrid browsers={browsers} onSelect={handleSelect} />
                </div>
                <div className="selector-footer">
                    {url && <p className="selector-footer-url">{url}</p>}
                </div>
            </div>
        </div>
    );
}

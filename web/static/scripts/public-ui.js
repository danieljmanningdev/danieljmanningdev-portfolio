/* Progressive enhancement only. Native links and the details menu work without
   this script; it never intercepts navigation, sends a message or loads data. */
(() => {
    'use strict';

    const menu = document.querySelector('.mobile-nav');
    const summary = menu?.querySelector('summary');

    function closeMenu(restoreFocus = false) {
        if (!menu?.open) return;
        menu.open = false;
        if (restoreFocus) summary?.focus({ preventScroll: true });
    }

    function markHomeSection() {
        if (window.location.pathname !== '/') return;
        document.querySelectorAll('.site-header a[href]').forEach(link => {
            const destination = new URL(link.href, window.location.href);
            const current = destination.origin === window.location.origin &&
                destination.pathname === '/' && destination.hash &&
                destination.hash === window.location.hash;
            if (current) link.setAttribute('aria-current', 'location');
            else link.removeAttribute('aria-current');
        });
    }

    if (menu && summary) {
        document.addEventListener('keydown', event => {
            if (event.key === 'Escape' && menu.open) {
                event.preventDefault();
                closeMenu(true);
            }
        });
        document.addEventListener('click', event => {
            if (menu.open && !menu.contains(event.target)) {
                closeMenu(menu.contains(document.activeElement));
            }
        });
        menu.addEventListener('click', event => {
            const link = event.target.closest('a[href]');
            if (!link || event.defaultPrevented || event.button !== 0 ||
                event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;
            const destination = new URL(link.href, window.location.href);
            closeMenu();
            if (destination.origin !== window.location.origin ||
                destination.pathname !== window.location.pathname ||
                destination.search !== window.location.search || !destination.hash) return;
            let section;
            try { section = document.getElementById(decodeURIComponent(destination.hash.slice(1))); }
            catch { return; }
            if (!section) return;
            // Move keyboard focus out of the now-collapsed menu, without
            // replacing the browser's own anchor navigation or history.
            requestAnimationFrame(() => {
                const target = section.querySelector('h1, h2, h3') || section;
                const temporaryTabindex = !target.hasAttribute('tabindex');
                if (temporaryTabindex) target.setAttribute('tabindex', '-1');
                target.focus({ preventScroll: true });
                if (temporaryTabindex) target.addEventListener('blur', () => target.removeAttribute('tabindex'), { once: true });
            });
        });
        window.addEventListener('resize', () => {
            if (getComputedStyle(menu).display === 'none') closeMenu();
        });
    }
    window.addEventListener('hashchange', markHomeSection);
    markHomeSection();

    const copy = document.querySelector('[data-copy-email]');
    const status = document.getElementById('contact-copy-status');
    if (copy && status && navigator.clipboard?.writeText) {
        copy.hidden = false;
        copy.addEventListener('click', async () => {
            if (copy.getAttribute('aria-busy') === 'true') return;
            copy.setAttribute('aria-busy', 'true');
            status.textContent = '';
            try {
                await navigator.clipboard.writeText(copy.dataset.copyEmail);
                status.textContent = 'Email address copied.';
            } catch {
                status.textContent = 'Could not copy automatically. Select and copy the address above.';
            } finally {
                copy.removeAttribute('aria-busy');
            }
        });
    }
})();

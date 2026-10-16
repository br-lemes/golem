type Theme = "auto" | "dark" | "light";

const root = document.documentElement;
const themeKey = "tabler-theme";
const prefersDark = window.matchMedia("(prefers-color-scheme: dark)");

const storedTheme = (): Theme => {
	try {
		const theme = localStorage.getItem(themeKey);
		if (theme === "dark" || theme === "light") {
			return theme;
		}
	} catch {
		// Fall back to auto.
	}
	return "auto";
};

const updateSwitchers = (theme: Theme): void => {
	for (const item of document.querySelectorAll<HTMLButtonElement>(
		"[data-golem-theme]",
	)) {
		const active = item.dataset.golemTheme === theme;
		item.classList.toggle("active", active);
		item.setAttribute("aria-pressed", String(active));
	}
	for (const icon of document.querySelectorAll<HTMLElement>(
		"[data-golem-theme-icon]",
	)) {
		icon.classList.toggle("d-none", icon.dataset.golemThemeIcon !== theme);
	}
};

const applyTheme = (theme: Theme, persist = false): void => {
	if (persist) {
		try {
			localStorage.setItem(themeKey, theme);
		} catch {
			// Keep the current selection.
		}
	}
	root.dataset.bsTheme =
		theme === "auto" ? (prefersDark.matches ? "dark" : "light") : theme;
	updateSwitchers(theme);
};

applyTheme(storedTheme());

prefersDark.addEventListener("change", () => {
	if (storedTheme() === "auto") {
		applyTheme("auto");
	}
});

document.addEventListener("click", (event: MouseEvent) => {
	const item = (event.target as Element | null)?.closest<HTMLButtonElement>(
		"[data-golem-theme]",
	);
	const theme = item?.dataset.golemTheme;
	if (theme === "auto" || theme === "dark" || theme === "light") {
		applyTheme(theme, true);
	}
});

document.addEventListener("DOMContentLoaded", () => {
	updateSwitchers(storedTheme());
});

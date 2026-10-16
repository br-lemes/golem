import { Collapse, Dropdown } from "@tabler/core";

const SYNC = true;

const manuallyOpenDropdowns = new Set<HTMLElement>();
const manuallyClosedDropdowns = new Set<HTMLElement>();

document.addEventListener("click", (event: MouseEvent) => {
	const toggle = (event.target as Element | null)?.closest<HTMLButtonElement>(
		"[data-sidebar-fold]",
	);
	if (!toggle) {
		return;
	}
	const sidebar = document.getElementById("sidebar");
	if (!sidebar) {
		return;
	}
	const folded = sidebar.classList.toggle("navbar-folded");
	toggle.setAttribute("aria-expanded", String(folded));
	toggle.setAttribute("aria-pressed", String(folded));
	if (folded) {
		for (const dropdown of sidebarDropdowns()) {
			setDropdownOpen(dropdown, false);
		}
	} else if (SYNC) {
		synchronizeDropdowns(sidebarDropdowns(), window.location.pathname);
	}
});

document.addEventListener("show.bs.dropdown", (event: Event) => {
	if (isSidebarFolded()) {
		closeFoldedSiblings(event);
		return;
	}
	if (!SYNC) {
		return;
	}
	const toggle = event.target as HTMLElement | null;
	if (!toggle?.matches("#sidebar-menu [data-bs-toggle='dropdown']")) {
		return;
	}
	const selected = toggle.closest<HTMLElement>(".dropdown, .dropend");
	if (!selected) {
		return;
	}
	event.preventDefault();
	manuallyOpenDropdowns.clear();
	for (
		let dropdown: HTMLElement | null = selected;
		dropdown;
		dropdown =
			dropdown.parentElement?.closest<HTMLElement>(
				".dropdown, .dropend",
			) ?? null
	) {
		manuallyOpenDropdowns.add(dropdown);
		manuallyClosedDropdowns.delete(dropdown);
	}
	synchronizeDropdowns(sidebarDropdowns(), window.location.pathname);
});

document.addEventListener("hide.bs.dropdown", (event: Event) => {
	if (isSidebarFolded() || !SYNC) {
		return;
	}
	const toggle = event.target as HTMLElement | null;
	const selected = toggle?.closest<HTMLElement>(".dropdown, .dropend");
	if (!selected?.closest("#sidebar-menu")) {
		return;
	}
	manuallyClosedDropdowns.add(selected);
	for (const dropdown of manuallyOpenDropdowns) {
		if (selected.contains(dropdown)) {
			manuallyOpenDropdowns.delete(dropdown);
		}
	}
});

export const updateNavigation = (href: string): void => {
	const path = new URL(href).pathname;
	for (const link of document.querySelectorAll<HTMLAnchorElement>(
		".navbar-nav [data-swr]",
	)) {
		const active = matchesPath(link.href, path);
		link.classList.toggle("active", active);
		if (link.parentElement?.classList.contains("nav-item")) {
			link.parentElement.classList.toggle("active", active);
		}
	}
	const dropdowns = sidebarDropdowns();
	for (const dropdown of dropdowns) {
		dropdown.classList.toggle("active", containsPath(dropdown, path));
	}
	if (!isSidebarFolded() && SYNC) {
		synchronizeDropdowns(dropdowns, path);
	}
};

export const closeSidebarMenu = (source: HTMLElement): void => {
	const menu = document.getElementById("sidebar-menu");
	if (!menu || !source.closest("#sidebar-menu")) {
		return;
	}
	if (window.matchMedia("(max-width: 991.98px)").matches) {
		(Collapse.getOrCreateInstance(menu) as Collapse).hide();
		return;
	}
	if (!menu.closest(".navbar-folded")) {
		return;
	}
	const toggle = source
		.closest(".dropdown")
		?.querySelector<HTMLElement>("[data-bs-toggle='dropdown']");
	if (toggle) {
		(Dropdown.getOrCreateInstance(toggle) as Dropdown).hide();
	}
};

const synchronizeDropdowns = (dropdowns: HTMLElement[], path: string): void => {
	for (const dropdown of dropdowns) {
		setDropdownOpen(
			dropdown,
			!manuallyClosedDropdowns.has(dropdown) &&
				(containsPath(dropdown, path) ||
					manuallyOpenDropdowns.has(dropdown)),
		);
	}
};

const closeFoldedSiblings = (event: Event): void => {
	const toggle = event.target as HTMLElement | null;
	if (!toggle?.matches("#sidebar-menu [data-bs-toggle='dropdown']")) {
		return;
	}
	const selected = toggle.closest<HTMLElement>(".dropend");
	const menu = selected?.parentElement;
	if (!selected || !menu?.matches(".dropdown-menu")) {
		return;
	}
	for (const sibling of menu.querySelectorAll<HTMLElement>(
		":scope > .dropend",
	)) {
		if (sibling === selected) {
			continue;
		}
		closeDropdown(sibling);
	}
};

const closeDropdown = (dropdown: HTMLElement): void => {
	const dropdowns = [
		dropdown,
		...dropdown.querySelectorAll<HTMLElement>(".dropend"),
	].reverse();
	for (const nested of dropdowns) {
		const toggle = nested.querySelector<HTMLElement>(
			":scope > [data-bs-toggle='dropdown']",
		);
		if (toggle) {
			(Dropdown.getOrCreateInstance(toggle) as Dropdown).hide();
		}
	}
};

const sidebarDropdowns = (): HTMLElement[] => {
	return [
		...document.querySelectorAll<HTMLElement>(
			"#sidebar-menu .dropdown, #sidebar-menu .dropend",
		),
	];
};

const isSidebarFolded = (): boolean => {
	return (
		document
			.getElementById("sidebar")
			?.classList.contains("navbar-folded") ?? false
	);
};

const containsPath = (dropdown: HTMLElement, path: string): boolean => {
	return [
		...dropdown.querySelectorAll<HTMLElement>("[data-swr-href]"),
		...dropdown.querySelectorAll<HTMLAnchorElement>("a[data-swr]"),
	].some((link) => {
		const href =
			link.dataset.swrHref ??
			(link instanceof HTMLAnchorElement ? link.href : undefined);
		return href && matchesPath(href, path);
	});
};

const matchesPath = (href: string, path: string): boolean => {
	const target = new URL(href, window.location.href).pathname;
	return path === target || (target !== "/" && path.startsWith(`${target}/`));
};

const setDropdownOpen = (dropdown: HTMLElement, open: boolean): void => {
	const toggle = dropdown.querySelector<HTMLElement>(
		":scope > [data-bs-toggle='dropdown']",
	);
	if (!toggle) {
		return;
	}
	const menu = dropdown.querySelector<HTMLElement>(":scope > .dropdown-menu");
	menu?.classList.toggle("show", open);
	toggle.classList.toggle("show", open);
	toggle.setAttribute("aria-expanded", String(open));
};

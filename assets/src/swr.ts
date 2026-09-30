import type { Htmx } from "htmx.org";
import htmx from "htmx.org";

declare global {
	interface Window {
		htmx: Htmx;
		tui?: {
			toast?: {
				add(options: {
					type: "error";
					title: string;
					description: string;
				}): string | null;
			};
		};
	}
}

window.htmx = htmx;

const TARGET = "#main-content";
const SELECT = "#main-content";
const SWAP = "outerHTML";

const pages = new Map<string, string>();

type Operation = {
	id: number;
	href: string;
	cancelled: boolean;
	controller: AbortController;
};

let currentOperation: Operation | undefined;
let nextOperationID = 0;

document.addEventListener("click", (event: MouseEvent) => {
	const target = event.target as Element | null;
	const source = target?.closest("[data-swr]") as HTMLElement | null;
	const href = source ? hrefFor(source) : undefined;
	if (!target || !source || !href) {
		return;
	}
	if (
		isInteractive(target, source) ||
		!isPlainNavigation(event, source, href)
	) {
		return;
	}
	event.preventDefault();
	console.log("[swr] click", { url: href });
	navigate(href);
});

window.addEventListener("popstate", () => {
	console.log("[swr] popstate", { url: window.location.href });
	void navigate(window.location.href, false);
});

document.addEventListener("collapsible-open-change", (event) => {
	if (
		!(event instanceof CustomEvent) ||
		!(event.target instanceof HTMLElement) ||
		!event.detail?.open
	) {
		return;
	}

	const selected = event.target;
	const keepOpen = new Set<HTMLElement>([selected]);
	let ancestor = selected.parentElement?.closest<HTMLElement>(
		"[data-tui-collapsible]",
	);
	while (ancestor) {
		keepOpen.add(ancestor);
		ancestor = ancestor.parentElement?.closest<HTMLElement>(
			"[data-tui-collapsible]",
		);
	}
	const currentPath = window.location.pathname;

	for (const root of document.querySelectorAll<HTMLElement>(
		"[data-tui-collapsible]",
	)) {
		const containsCurrentPath = [
			...root.querySelectorAll<HTMLElement>("[data-swr-href]"),
			...root.querySelectorAll<HTMLAnchorElement>("a[href]"),
		].some((link) => {
			const href =
				link instanceof HTMLAnchorElement
					? link.href
					: link.dataset.swrHref;
			return href && new URL(href, window.location.href).pathname === currentPath;
		});
		if (keepOpen.has(root) || containsCurrentPath) {
			continue;
		}
		root.querySelector<HTMLElement>(
			"[data-tui-collapsible-trigger][aria-expanded='true']",
		)?.click();
	}
});

const navigate = async (
	href: string,
	pushHistory = true,
): Promise<void> => {
	cancel(currentOperation);
	const previousURL = window.location.href;
	const requestSource = document.getElementById("swr-loading-overlay");
	if (!requestSource) {
		console.error("[swr] loading overlay not found");
		return;
	}
	const operation = {
		id: ++nextOperationID,
		href,
		cancelled: false,
		controller: new AbortController(),
	};
	currentOperation = operation;
	if (pushHistory) {
		history.pushState({}, "", operation.href);
	}
	updateNavigation(operation.href);

	const cached = pages.get(operation.href);
	setLoading(operation, cached ? "cache" : "miss");
	console.log(`[swr] ${cached ? "cache hit" : "cache miss"}`, {
		id: operation.id,
		url: operation.href,
	});

	try {
		if (cached) {
			const target = document.querySelector(TARGET);
			if (target) {
				console.log("[swr] stale swap", {
					id: operation.id,
					url: operation.href,
				});
				await htmx.swap({
					text: cached,
					target,
					swap: SWAP,
					sourceElement: requestSource,
				});
			}
		}

		if (operation.cancelled || currentOperation !== operation) {
			return;
		}

		setLoading(operation, cached ? "cache" : "miss");
		console.log("[swr] request start", {
			id: operation.id,
			url: operation.href,
		});
		const response = await fetch(operation.href, {
			headers: {
				"HX-Current-URL": window.location.href,
				"HX-Request": "true",
				"HX-Target": TARGET,
			},
			signal: operation.controller.signal,
		});
		if (!response.ok) {
			throw new Error(`request failed with status ${response.status}`);
		}
		const text = await response.text();
		if (operation.cancelled || currentOperation !== operation) {
			return;
		}
		const selected = selectContent(text);
		pages.set(operation.href, selected);
		console.log("[swr] response cached", {
			id: operation.id,
			url: operation.href,
		});
		const target = document.querySelector(TARGET);
		if (!target) {
			throw new Error(`target not found: ${TARGET}`);
		}
		await htmx.swap({
			text: selected,
			target,
			swap: SWAP,
			sourceElement: requestSource,
		});
	} catch (error) {
		if (!operation.cancelled && pushHistory) {
			history.replaceState({}, "", previousURL);
			showError(
				error instanceof Error ? error.message : "request failed",
			);
			console.log("[swr] operation failed", {
				id: operation.id,
				url: operation.href,
				error,
			});
		}
	} finally {
		if (currentOperation === operation) {
			setLoading(operation, "idle");
			currentOperation = undefined;
		}
		console.log("[swr] operation finished", {
			id: operation.id,
			url: operation.href,
		});
	}
};

const showError = (message: string): void => {
	const toast = window.tui?.toast;
	if (!toast) {
		console.error("[swr] toast API not available", { message });
		return;
	}
	toast.add({ type: "error", title: "Request failed", description: message });
};

const cancel = (operation: Operation | undefined): void => {
	if (!operation || operation.cancelled) {
		return;
	}
	operation.cancelled = true;
	operation.controller.abort();
	console.log("[swr] operation cancelled", {
		id: operation.id,
		url: operation.href,
	});
};

const setLoading = (
	operation: Operation,
	state: "cache" | "miss" | "idle",
): void => {
	if (currentOperation !== operation) {
		return;
	}
	const bar = document.getElementById("swr-loading-bar");
	const overlay = document.getElementById("swr-loading-overlay");
	const wrapper = document.getElementById("swr-content-wrapper");
	bar?.classList.toggle("opacity-0", state !== "cache");
	bar?.setAttribute("aria-hidden", String(state !== "cache"));
	overlay?.classList.toggle("hidden", state !== "miss");
	overlay?.setAttribute("aria-hidden", String(state !== "miss"));
	wrapper?.setAttribute("aria-busy", String(state !== "idle"));
};

const selectContent = (text: string): string => {
	const document = new DOMParser().parseFromString(text, "text/html");
	return document.querySelector(SELECT)?.outerHTML || text;
};

const updateNavigation = (href: string): void => {
	const path = new URL(href).pathname;
	for (const link of document.querySelectorAll<HTMLElement>(
		'[data-swr][data-sidebar="menu-button"], [data-swr][data-sidebar="menu-sub-button"]',
	)) {
		const active =
			link instanceof HTMLAnchorElement &&
			new URL(link.href).pathname === path;
		link.toggleAttribute("data-active", active);
	}
	for (const trigger of document.querySelectorAll<HTMLElement>(
		"[data-swr-href]",
	)) {
		const triggerPath = trigger.dataset.swrHref;
		const active = triggerPath
			? new URL(triggerPath, window.location.href).pathname === path
			: false;
		trigger.toggleAttribute("data-active", active);
	}
	for (const root of document.querySelectorAll<HTMLElement>(
		"[data-tui-collapsible]",
	)) {
		const open =
			[...root.querySelectorAll<HTMLElement>("[data-swr-href]")].some(
				(trigger) =>
					new URL(trigger.dataset.swrHref ?? "", window.location.href)
						.pathname === path,
			) ||
			[...root.querySelectorAll<HTMLAnchorElement>("a[href]")].some(
				(link) => new URL(link.href).pathname === path,
			);
		root.toggleAttribute("data-open", open);
		root.toggleAttribute("data-closed", !open);
		const content = root.querySelector<HTMLElement>(
			"[data-tui-collapsible-content]",
		);
		content?.toggleAttribute("data-open", open);
		content?.toggleAttribute("data-closed", !open);
		if (content) {
			content.hidden = !open;
		}
		const trigger = root.querySelector<HTMLElement>(
			"[data-tui-collapsible-trigger]",
		);
		trigger?.setAttribute("aria-expanded", String(open));
		trigger?.toggleAttribute("data-panel-open", open);
	}
};

const hrefFor = (source: HTMLElement): string | undefined => {
	const value = source.dataset.swr;
	if (!value) {
		console.error("[swr] data-swr must not be empty", source);
		return;
	}
	if (value === "true") {
		if (!(source instanceof HTMLAnchorElement) || !source.href) {
			console.error(
				'[swr] data-swr="true" requires an anchor with an href',
				source,
			);
			return;
		}
		return source.href;
	}
	try {
		return new URL(value, window.location.href).href;
	} catch (error) {
		console.error("[swr] invalid data-swr URL", { source, value, error });
		return;
	}
};

const isInteractive = (target: Element, source: HTMLElement): boolean => {
	const interactive = target.closest(
		"a, button, input, select, textarea, summary, [role='button'], [contenteditable]",
	);
	return interactive !== null && interactive !== source;
};

const isPlainNavigation = (
	event: MouseEvent,
	source: HTMLElement,
	href: string,
): boolean => {
	return (
		event.button === 0 &&
		window.getSelection()?.isCollapsed !== false &&
		!event.defaultPrevented &&
		!event.metaKey &&
		!event.ctrlKey &&
		!event.shiftKey &&
		!event.altKey &&
		(!(source instanceof HTMLAnchorElement) || !source.target) &&
		new URL(href).origin === window.location.origin
	);
};

const cacheInitialPage = (): void => {
	const target = document.querySelector(TARGET);
	if (!target) {
		return;
	}
	pages.set(window.location.href, target.outerHTML);
	console.log("[swr] initial page cached", { url: window.location.href });
};

cacheInitialPage();

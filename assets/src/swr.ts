import { SELECTOR, selectContent, swapContent } from "./content";
import { closeSidebarMenu, updateNavigation } from "./navigation";
import * as notifications from "./notifications";

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
	closeSidebarMenu(source);
	console.log("[swr] click", { url: href });
	navigate(href);
});

window.addEventListener("popstate", () => {
	console.log("[swr] popstate", { url: window.location.href });
	void navigate(window.location.href, false);
});

const navigate = async (href: string, pushHistory = true): Promise<void> => {
	cancel(currentOperation);
	const previousURL = window.location.href;
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
			console.log("[swr] stale swap", {
				id: operation.id,
				url: operation.href,
			});
			swapContent(cached);
			if (pushHistory) {
				scrollToTop();
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
		swapContent(selected);
		if (pushHistory && !cached) {
			scrollToTop();
		}
	} catch (error) {
		if (!operation.cancelled && pushHistory) {
			history.replaceState({}, "", previousURL);
			notifications.error(
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
	bar?.classList.toggle("d-none", state === "idle");
	overlay?.classList.toggle("d-none", state !== "miss");
	wrapper?.setAttribute("aria-busy", String(state !== "idle"));
};

const scrollToTop = (): void => {
	window.scrollTo(0, 0);
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
	const target = document.querySelector(SELECTOR);
	if (!target) {
		return;
	}
	pages.set(window.location.href, target.outerHTML);
	console.log("[swr] initial page cached", { url: window.location.href });
};

cacheInitialPage();

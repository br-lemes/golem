import { SELECTOR, selectContent, swapContent } from "./content";
import { closeSidebarMenu, updateNavigation } from "./navigation";
import * as notifications from "./notifications";

const pages = new Map<string, string>();

const defaultRefreshPeriod = 5;
const refreshPeriodCookie = "swr-refresh-period";

const refreshPeriods = [0, 5, 10, 15, 30, 45, 60, 90];

type Operation = {
	id: number;
	href: string;
	cancelled: boolean;
	controller: AbortController;
};

let currentOperation: Operation | undefined;
let nextOperationID = 0;
let refreshTimer: number | undefined;

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
	void navigate(href);
});

document.addEventListener("change", (event: Event) => {
	const select = event.target as HTMLSelectElement | null;
	if (!select?.matches("[data-swr-refresh='auto']")) {
		return;
	}
	// biome-ignore lint/suspicious/noDocumentCookie: LAN HTTP, single write
	document.cookie = `${refreshPeriodCookie}=${select.value}; path=/; max-age=31536000; samesite=lax`;
	scheduleRefresh();
});

window.addEventListener("popstate", () => {
	console.log("[swr] popstate", { url: window.location.href });
	void navigate(window.location.href, false);
});

export const navigate = async (
	href: string,
	pushHistory = true,
): Promise<void> => {
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
		if (cached && pushHistory && !hasDiffContent()) {
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
		if (!diffContent(text)) {
			swapContent(selected);
		}
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
			scheduleRefresh();
		}
		console.log("[swr] operation finished", {
			id: operation.id,
			url: operation.href,
		});
	}
};

const diffContent = (text: string): boolean => {
	const current = document.querySelector<HTMLTableElement>(
		"table[data-swr-diff]",
	);
	const incomingDocument = new DOMParser().parseFromString(text, "text/html");
	const incoming = incomingDocument.querySelector<HTMLTableElement>(
		"table[data-swr-diff]",
	);
	if (!current || !incoming) {
		return false;
	}
	for (const region of ["summary", "pagination"]) {
		const currentRegion = document.querySelector<HTMLElement>(
			`[data-swr-region="${region}"]`,
		);
		const incomingRegion = incomingDocument.querySelector<HTMLElement>(
			`[data-swr-region="${region}"]`,
		);
		if (!currentRegion || !incomingRegion) {
			continue;
		}
		if (currentRegion.innerHTML !== incomingRegion.innerHTML) {
			currentRegion.replaceWith(
				document.importNode(incomingRegion, true),
			);
		}
	}
	const currentRows = new Map(
		Array.from(
			current.querySelectorAll<HTMLTableRowElement>(
				"tbody tr[data-swr-key]",
			),
		).map((row) => [row.dataset.swrKey, row]),
	);
	const incomingRows = Array.from(
		incoming.querySelectorAll<HTMLTableRowElement>(
			"tbody tr[data-swr-key]",
		),
	);
	if (currentRows.size !== incomingRows.length) {
		return false;
	}
	for (const incomingRow of incomingRows) {
		const currentRow = currentRows.get(incomingRow.dataset.swrKey);
		if (!currentRow) {
			return false;
		}
		for (const incomingCell of incomingRow.querySelectorAll<HTMLElement>(
			"[data-swr-field]",
		)) {
			const field = incomingCell.dataset.swrField;
			const currentCell = currentRow.querySelector<HTMLElement>(
				`[data-swr-field="${field}"]`,
			);
			if (
				!currentCell ||
				currentCell.innerHTML === incomingCell.innerHTML
			) {
				continue;
			}
			currentCell.innerHTML = incomingCell.innerHTML;
			currentCell.classList.add("table-warning");
			window.setTimeout(
				() => currentCell.classList.remove("table-warning"),
				2000,
			);
		}
	}
	return true;
};

const hasDiffContent = (): boolean =>
	document.querySelector("table[data-swr-diff]") !== null;

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

const scheduleRefresh = (): void => {
	clearRefreshTimer();
	const source = document.querySelector<HTMLSelectElement>(
		'select[data-swr-refresh="auto"]',
	);
	if (!source) {
		return;
	}
	const seconds = readRefreshPeriod();
	if (seconds <= 0) {
		return;
	}
	refreshTimer = window.setTimeout(() => {
		refreshTimer = undefined;
		if (
			document.visibilityState === "visible" &&
			currentOperation === undefined
		) {
			void navigate(window.location.href, false);
			return;
		}
		scheduleRefresh();
	}, seconds * 1000);
};

const clearRefreshTimer = (): void => {
	if (refreshTimer !== undefined) {
		window.clearTimeout(refreshTimer);
		refreshTimer = undefined;
	}
};

function readRefreshPeriod(): number {
	const match = document.cookie.match(
		new RegExp(`(?:^|;\\s*)${refreshPeriodCookie}=([^;]*)`),
	);
	if (!match) return defaultRefreshPeriod;
	const parsed = Number.parseInt(decodeURIComponent(match[1]), 10);
	return refreshPeriods.includes(parsed) ? parsed : defaultRefreshPeriod;
}

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
	if (source instanceof HTMLAnchorElement && source.target) {
		console.error(
			"[swr] anchors with data-swr must not set target; deferring to browser",
			source,
		);
		return false;
	}
	const sel = window.getSelection();
	const hasSelectionInSource =
		sel?.isCollapsed === false &&
		sel?.anchorNode != null &&
		source.contains(sel?.anchorNode);
	return (
		event.button === 0 &&
		!hasSelectionInSource &&
		!event.defaultPrevented &&
		!event.metaKey &&
		!event.ctrlKey &&
		!event.shiftKey &&
		!event.altKey &&
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
	scheduleRefresh();
};

cacheInitialPage();

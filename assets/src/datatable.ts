import { Modal } from "@tabler/core";
import { navigate } from "./swr";

// Keep this value in sync with ui/components/datatable/datatable.templ.
const defaultSize = "10";

const filterModalHash = "#filters";

document.addEventListener("submit", (event: SubmitEvent) => {
	const form = event.target;
	if (
		!(form instanceof HTMLFormElement) ||
		!form.hasAttribute("data-swr-navigate")
	) {
		return;
	}
	const effectFields = effectFilterFields(form);
	if (form.hasAttribute("data-effect-filter") && effectFields.length === 3) {
		const feedback = form.querySelector<HTMLElement>(
			"[data-effect-filter-feedback]",
		);
		const effectFilterValid = isEffectFilterValid(effectFields);
		effectFields.forEach((field) => {
			field.classList.toggle(
				"is-invalid",
				!effectFilterValid && field.value.trim() === "",
			);
		});
		feedback?.classList.toggle("d-none", effectFilterValid);
		feedback?.classList.toggle("d-block", !effectFilterValid);
		if (!effectFilterValid) {
			event.preventDefault();
			return;
		}
	}
	const modal = form.closest<HTMLElement>(".modal");
	if (modal) {
		(Modal.getOrCreateInstance(modal) as Modal).hide();
	}
	event.preventDefault();
	const url = new URL(form.action, window.location.href);
	const data = new FormData(form);
	const currentURL = new URL(window.location.href);
	const effect = data.get("filter[effect]")?.toString().trim() ?? "";
	const effectOperator =
		data.get("filter[effect_operator]")?.toString().trim() ?? "";
	url.searchParams.delete("page");
	data.forEach((value, key) => {
		if (typeof value !== "string") {
			return;
		}
		const normalized = value.trim();
		if (
			key === "filter[effect_value]" &&
			normalized === "0" &&
			effect === "" &&
			effectOperator === ""
		) {
			url.searchParams.delete(key);
		} else if (normalized === "") {
			url.searchParams.delete(key);
		} else {
			url.searchParams.set(key, normalized);
		}
	});
	if (!data.has("search")) {
		const search = currentURL.searchParams.get("search");
		if (search !== null && search !== "") {
			url.searchParams.set("search", search);
		}
	}
	if (!data.has("size")) {
		const size = currentURL.searchParams.get("size");
		if (size !== null && size !== "" && size !== defaultSize) {
			url.searchParams.set("size", size);
		}
	}
	if (url.searchParams.get("size") === defaultSize) {
		url.searchParams.delete("size");
	}
	navigateWithSearchFocus(form, url.href);
});

document.addEventListener("change", (event: Event) => {
	const select = event.target;
	if (
		!(select instanceof HTMLSelectElement) ||
		!select.hasAttribute("data-swr-navigate")
	) {
		return;
	}
	const url = new URL(window.location.href);
	if (select.value === defaultSize) {
		url.searchParams.delete("size");
	} else {
		url.searchParams.set("size", select.value);
	}
	url.searchParams.delete("page");
	void navigate(url.href);
});

const clearEffectFilterValidation = (event: Event): void => {
	const target = event.target;
	if (
		!(
			target instanceof HTMLInputElement ||
			target instanceof HTMLSelectElement
		) ||
		!target.name.startsWith("filter[")
	) {
		return;
	}
	const form = target.form;
	if (!form?.hasAttribute("data-effect-filter")) {
		return;
	}
	const fields = effectFilterFields(form);
	fields.forEach((field) => {
		field.classList.remove("is-invalid");
	});
	const feedback = form.querySelector<HTMLElement>(
		"[data-effect-filter-feedback]",
	);
	feedback?.classList.add("d-none");
	feedback?.classList.remove("d-block");
};

document.addEventListener("input", clearEffectFilterValidation);
document.addEventListener("change", clearEffectFilterValidation);

document.addEventListener("click", (event: MouseEvent) => {
	const target = event.target;
	if (!(target instanceof Element)) {
		return;
	}
	const button = target.closest<HTMLButtonElement>("[data-filter-reset]");
	const form = button?.form;
	if (!button || !form) {
		return;
	}
	form.querySelectorAll<HTMLInputElement | HTMLSelectElement>(
		"[name^='filter[']",
	).forEach((field) => {
		field.value = field.name === "filter[effect_value]" ? "0" : "";
		field.classList.remove("is-invalid");
	});
	const feedback = form.querySelector<HTMLElement>(
		"[data-effect-filter-feedback]",
	);
	feedback?.classList.add("d-none");
	feedback?.classList.remove("d-block");
});

document.addEventListener(
	"click",
	(event: MouseEvent) => {
		const target = event.target;
		if (!(target instanceof Element)) {
			return;
		}
		const button = target.closest<HTMLButtonElement>(
			"[data-effect-filter-submit]",
		);
		const form = button?.form;
		if (!button || !form?.hasAttribute("data-effect-filter")) {
			return;
		}
		const fields = effectFilterFields(form);
		if (!isEffectFilterValid(fields)) {
			const feedback = form.querySelector<HTMLElement>(
				"[data-effect-filter-feedback]",
			);
			fields.forEach((field) => {
				field.classList.toggle("is-invalid", !field.value.trim());
			});
			feedback?.classList.remove("d-none");
			feedback?.classList.add("d-block");
			event.preventDefault();
			event.stopImmediatePropagation();
		}
	},
	{ capture: true },
);

type EffectField = HTMLInputElement | HTMLSelectElement;

const effectFilterFields = (form: HTMLFormElement): EffectField[] =>
	[
		form.elements.namedItem("filter[effect]"),
		form.elements.namedItem("filter[effect_operator]"),
		form.elements.namedItem("filter[effect_value]"),
	].filter(
		(field): field is EffectField =>
			field instanceof HTMLInputElement ||
			field instanceof HTMLSelectElement,
	);

const isEffectFilterValid = (
	fields: Array<HTMLInputElement | HTMLSelectElement>,
): boolean => {
	const [effect, operator, value] = fields.map((field) => field.value.trim());
	if (effect === "" && operator === "" && (value === "" || value === "0")) {
		return true;
	}
	return effect !== "" && operator !== "" && value !== "";
};

const navigateWithSearchFocus = (form: HTMLFormElement, href: string): void => {
	const searchInput = form.elements.namedItem("search");
	const restoreFocus =
		searchInput instanceof HTMLInputElement &&
		document.activeElement === searchInput;
	const selectionStart =
		searchInput instanceof HTMLInputElement
			? searchInput.selectionStart
			: null;
	const selectionEnd =
		searchInput instanceof HTMLInputElement
			? searchInput.selectionEnd
			: null;
	void navigate(href).then(() => {
		if (!restoreFocus || window.location.href !== href) {
			return;
		}
		const input = document.querySelector<HTMLInputElement>(
			"input[name='search']",
		);
		if (!input) {
			return;
		}
		input.focus();
		if (selectionStart !== null && selectionEnd !== null) {
			input.setSelectionRange(selectionStart, selectionEnd);
		}
	});
};

const updateFilterModalHash = (hash: string): void => {
	const url = new URL(window.location.href);
	url.hash = hash;
	history.replaceState(null, "", url);
};

const openFilterModalFromHash = (): void => {
	if (window.location.hash !== filterModalHash) {
		return;
	}
	const button = document.querySelector<HTMLButtonElement>(
		'[data-filter-modal][data-bs-toggle="modal"][data-bs-target]',
	);
	if (button) {
		button.click();
	}
};

document.addEventListener("DOMContentLoaded", openFilterModalFromHash);
window.addEventListener("hashchange", openFilterModalFromHash);
new MutationObserver(openFilterModalFromHash).observe(document.body, {
	childList: true,
	subtree: true,
});

document.addEventListener("click", (event: MouseEvent) => {
	const button = event.target;
	if (
		!(button instanceof Element) ||
		!button.closest(
			'[data-filter-modal][data-bs-toggle="modal"][data-bs-target]',
		)
	) {
		return;
	}
	updateFilterModalHash(filterModalHash);
});

document.addEventListener("hidden.bs.modal", (event: Event) => {
	if (
		event.target instanceof Element &&
		window.location.hash === filterModalHash
	) {
		updateFilterModalHash("");
	}
});

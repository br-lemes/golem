export const SELECTOR = "#main-content";

export const selectContent = (text: string): string => {
	const page = new DOMParser().parseFromString(text, "text/html");
	const content = page.querySelector<HTMLElement>(SELECTOR);
	if (!content) {
		throw new Error(`content not found: ${SELECTOR}`);
	}
	return content.outerHTML;
};

export const swapContent = (text: string): void => {
	const target = document.querySelector<HTMLElement>(SELECTOR);
	if (!target) {
		throw new Error(`target not found: ${SELECTOR}`);
	}
	const page = new DOMParser().parseFromString(text, "text/html");
	const content = page.querySelector<HTMLElement>(SELECTOR);
	if (!content) {
		throw new Error(`content not found: ${SELECTOR}`);
	}
	target.replaceWith(document.importNode(content, true));
};

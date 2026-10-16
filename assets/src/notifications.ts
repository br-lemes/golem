import { Toast } from "@tabler/core";

type Variant = "success" | "danger";

const show = (variant: Variant, message: string): void => {
	const container = document.getElementById("notifications");
	const template = document.getElementById("notification-template");
	if (!container || !(template instanceof HTMLTemplateElement)) {
		console.error("[notifications] template or container not found");
		return;
	}
	const toast = template.content.firstElementChild?.cloneNode(true);
	if (!(toast instanceof HTMLElement)) {
		console.error("[notifications] template must contain an element");
		return;
	}
	const body = toast.querySelector<HTMLElement>(
		"[data-notification-message]",
	);
	if (!body) {
		console.error("[notifications] message element not found");
		return;
	}
	body.textContent = message;
	toast.classList.add(`text-bg-${variant}`);
	container.append(toast);
	toast.addEventListener("hidden.bs.toast", () => toast.remove());
	new Toast(toast).show();
};

export const error = (message: string): void => {
	show("danger", message);
};

export const success = (message: string): void => {
	show("success", message);
};

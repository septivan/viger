import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { Pagination } from "@/components/pagination";

describe("Pagination", () => {
  it("exposes the current page and navigation controls", async () => {
    const onPageChange = vi.fn();
    render(<Pagination label="Game pages" onPageChange={onPageChange} page={3} totalPages={8} />);
    expect(screen.getByRole("navigation", { name: "Game pages" })).toBeVisible();
    expect(screen.getByRole("button", { name: "3" })).toHaveAttribute("aria-current", "page");
    await userEvent.click(screen.getByRole("button", { name: "Next" }));
    expect(onPageChange).toHaveBeenCalledWith(4);
  });

  it("does not render for a single page", () => {
    const { container } = render(<Pagination label="Game pages" onPageChange={() => undefined} page={1} totalPages={1} />);
    expect(container).toBeEmptyDOMElement();
  });
});


import { expect, test } from "@playwright/test";

test("a visitor can search the catalog and open a game", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("heading", { name: /Find the game/i })).toBeVisible();
  await page.getByRole("searchbox", { name: "Search games" }).fill("Hades");
  await expect(page.getByRole("heading", { name: "Hades", exact: true })).toBeVisible();
  await page.getByRole("link", { name: "Read reviews for Hades" }).click();
  await expect(page.getByRole("heading", { name: "Hades", exact: true })).toBeVisible();
  await expect(page.getByText("The conversation")).toBeVisible();
});

test("catalog filters and ordering are reflected in visible card values", async ({ page }) => {
  await page.goto("/");
  await page.getByLabel("Platform", { exact: true }).selectOption("Mobile");

  const cards = page.locator(".game-card");
  await expect(cards.first()).toContainText("Mobile");
  for (const text of await cards.allTextContents()) {
    expect(text).toContain("Mobile");
  }

  await page.getByLabel("Sort by").selectOption("newest");
  const years = (await cards.allTextContents()).map((text) => {
    const year = text.match(/\b20\d{2}\b/)?.[0];
    expect(year).toBeDefined();
    return Number(year);
  });
  expect(years).toEqual([...years].sort((left, right) => right - left));
});

test("a new review appears immediately in another browser", async ({ browser }) => {
  const authorContext = await browser.newContext();
  const readerContext = await browser.newContext();
  const author = await authorContext.newPage();
  const reader = await readerContext.newPage();
  await Promise.all([author.goto("/games/game-002"), reader.goto("/games/game-002")]);
  await Promise.all([
    expect(author.getByText("Live reviews")).toBeVisible(),
    expect(reader.getByText("Live reviews")).toBeVisible(),
  ]);

  const uniqueReview = `Realtime review ${Date.now()} with excellent pacing and responsive combat.`;
  await author.getByRole("radio", { name: "5 out of 5" }).check();
  await author.getByLabel("Your name").fill("Playwright Author");
  await author.getByLabel("Your review").fill(uniqueReview);
  await author.getByRole("button", { name: /Publish review/i }).click();

  await expect(author.getByRole("status").filter({ hasText: "part of the conversation" })).toBeVisible();
  await expect(reader.getByText("A new review from Playwright Author just arrived.")).toBeVisible();
  await expect(reader.getByText(uniqueReview)).toBeVisible();
  await authorContext.close();
  await readerContext.close();
});
